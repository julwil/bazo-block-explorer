package main

import (
  "net"
  _ "os"
  _ "strings"
  "fmt"
  "log"
  "time"
  "errors"
  "bufio"
  _ "github.com/mchetelat/bazo_miner/miner"
	"github.com/mchetelat/bazo_miner/p2p"
	"github.com/mchetelat/bazo_miner/protocol"
)

var allBlockHeaders []*protocol.SPVHeader
var allBlocks []*protocol.Block
var logger *log.Logger

var block1 *protocol.Block

/*
func main()  {
  initState()
}
*/
func initState() {
  loadAllBlocks()
  /*
  duration := time.Minute * 2
  time.Sleep(duration)

  refreshState()

  duration = time.Second * 1
  time.Sleep(duration)
  */
  refreshState()
  //allBlocks = invertBlockArray(allBlocks)
}

//Update allBlockHeaders to the latest header
func refreshState() {
	var newBlocks []*protocol.Block
	newBlocks = getNewBlocks(reqBlock(nil), allBlocks[len(allBlocks)-1], newBlocks)

  newBlocks = invertBlockArray(newBlocks)

  for _, block := range newBlocks{
    //convert block
    convertedBlock := ConvertBlock(block)
    //write to db here
    fmt.Printf("Writing Block: %s\n", convertedBlock.Hash)
    WriteBlock(convertedBlock)
  }
	allBlocks = append(allBlocks, newBlocks...)
}

//Get new blockheaders recursively
func getNewBlocks(latest *protocol.Block, eldest *protocol.Block, list []*protocol.Block) []*protocol.Block {
	if latest.Hash != eldest.Hash {
		ancestor := reqBlock(latest.PrevHash[:])
		list = getNewBlocks(ancestor, eldest, list)
		list = append(list, latest)
	}
  return list
}

func loadAllBlocks() {
  //genesishash := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
  block := reqBlock(nil)
  allBlocks = append(allBlocks, block)
  prevHash := block.PrevHash

  for block.Hash != [32]byte{} {
    block = reqBlock(prevHash[:])
    allBlocks = append(allBlocks, block)
    prevHash = block.PrevHash
  }

  allBlocks = invertBlockArray(allBlocks)

  for _, oneBlock := range allBlocks{

    for _, fundsTxHash := range oneBlock.FundsTxData{
      fundsTx := reqTx(p2p.FUNDSTX_REQ, fundsTxHash)
      convertedTx := ConvertFundsTransaction(fundsTx.(*protocol.FundsTx), oneBlock.Hash, fundsTxHash)

      fmt.Printf("Writing Transaction: %s\n", convertedTx.Hash)
      WriteFundsTx(convertedTx)
    }

    for _, accTxHash := range oneBlock.AccTxData{
      accTx := reqTx(p2p.ACCTX_REQ, accTxHash)
      convertedTx := ConvertAccTransaction(accTx.(*protocol.AccTx), oneBlock.Hash, accTxHash)

      fmt.Printf("Writing Transaction: %s\n", convertedTx.Hash)
      WriteAccTx(convertedTx)
    }

    for _, configTxHash := range oneBlock.ConfigTxData{
      configTx := reqTx(p2p.CONFIGTX_REQ, configTxHash)
      convertedTx := ConvertConfigTransaction(configTx.(*protocol.ConfigTx), oneBlock.Hash, configTxHash)

      fmt.Printf("Writing Transaction: %s\n", convertedTx.Hash)
      WriteConfigTx(convertedTx)
    }

    convertedBlock := ConvertBlock(oneBlock)
    fmt.Printf("Writing Block: %s\n\n", convertedBlock.Hash)
    WriteBlock(convertedBlock)
  }
}

func Connect(connectionString string) (conn net.Conn) {
	conn, err := net.Dial("tcp", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetDeadline(time.Now().Add(20 * time.Second))

	return conn
}

func reqBlock(blockHash []byte) (block *protocol.Block) {

	conn := Connect(p2p.BOOTSTRAP_SERVER)

	packet := p2p.BuildPacket(p2p.BLOCK_REQ, blockHash[:])
	conn.Write(packet)

	header, payload, err := rcvData(conn)
	if err != nil {
		logger.Printf("Disconnected: %v\n", err)
		return
	}

	if header.TypeID == p2p.BLOCK_RES {
		block = block.Decode(payload)
	}

	conn.Close()

	return block
}

func reqTx(txType uint8, txHash [32]byte) interface{} {

	conn := Connect(p2p.BOOTSTRAP_SERVER)

	packet := p2p.BuildPacket(txType, txHash[:])
	conn.Write(packet)

	header, payload, err := rcvData(conn)
	if err != nil {
		logger.Printf("Disconnected: %v\n", err)
		panic(err)
	}
  defer conn.Close()

	switch header.TypeID {
	case p2p.ACCTX_RES:
		var accTx *protocol.AccTx
		accTx = accTx.Decode(payload)
		return accTx
	case p2p.CONFIGTX_RES:
		var configTx *protocol.ConfigTx
		configTx = configTx.Decode(payload)
		return configTx
	case p2p.FUNDSTX_RES:
		var fundsTx *protocol.FundsTx
		fundsTx = fundsTx.Decode(payload)
		return fundsTx
  default:
    panic(err)
	}
}

func rcvData(c net.Conn) (header *p2p.Header, payload []byte, err error) {
	reader := bufio.NewReader(c)
	header, err = p2p.ReadHeader(reader)
	if err != nil {
		c.Close()
		return nil, nil, errors.New(fmt.Sprintf("Connection to aborted: (%v)\n", err))
	}
	payload = make([]byte, header.Len)

	for cnt := 0; cnt < int(header.Len); cnt++ {
		payload[cnt], err = reader.ReadByte()
		if err != nil {
			c.Close()
			return nil, nil, errors.New(fmt.Sprintf("Connection to aborted: %v\n", err))
		}
	}

	//logger.Printf("Receive message:\nSender: %v\nType: %v\nPayload length: %v\n", p.getIPPort(), logMapping[header.TypeID], len(payload))
	return header, payload, nil
}

func invertBlockArray(array []*protocol.Block) []*protocol.Block {
	for i, j := 0, len(array)-1; i < j; i, j = i+1, j-1 {
		array[i], array[j] = array[j], array[i]
	}

	return array
}

func ConvertBlock(unconvertedBlock *protocol.Block) block {
  var convertedBlock block
  var convertedTxHash string
  //convertedBlock.Header = fmt.Sprintf("%x", unconvertedBlock.Header)
  convertedBlock.Hash = fmt.Sprintf("%x", unconvertedBlock.Hash)
  convertedBlock.PrevHash = fmt.Sprintf("%x", unconvertedBlock.PrevHash)
  //convertedBlock.Nonce = fmt.Sprintf("%x", unconvertedBlock.Nonce)
  convertedBlock.Timestamp = unconvertedBlock.Timestamp
  convertedBlock.MerkleRoot = fmt.Sprintf("%x", unconvertedBlock.MerkleRoot)
  convertedBlock.Beneficiary = fmt.Sprintf("%x", unconvertedBlock.Beneficiary)
  convertedBlock.NrFundsTx = unconvertedBlock.NrFundsTx
  convertedBlock.NrAccTx = unconvertedBlock.NrAccTx
  convertedBlock.NrConfigTx = unconvertedBlock.NrConfigTx
  for _, hash := range unconvertedBlock.FundsTxData {
    convertedTxHash = fmt.Sprintf("%x", hash)
    convertedBlock.FundsTxData = append(convertedBlock.FundsTxData, convertedTxHash)
  }
  fmt.Printf("%s\n", convertedBlock.FundsTxData)
  for _, hash := range unconvertedBlock.AccTxData {
    convertedTxHash = fmt.Sprintf("%x", hash)
    convertedBlock.AccTxData = append(convertedBlock.AccTxData, convertedTxHash)
  }
  fmt.Printf("%s\n", convertedBlock.AccTxData)
  for _, hash := range unconvertedBlock.ConfigTxData {
    convertedTxHash = fmt.Sprintf("%x", hash)
    convertedBlock.ConfigTxData = append(convertedBlock.ConfigTxData, convertedTxHash)
  }
  fmt.Printf("%s\n", convertedBlock.ConfigTxData)
  return convertedBlock
}

func ConvertFundsTransaction(unconvertedTx *protocol.FundsTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte) fundstx {
  var convertedTx fundstx
  //convertedTx.Header = fmt.Sprintf("%x", unconvertedTx.Header)
  convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
  convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
  convertedTx.Amount = unconvertedTx.Amount
  convertedTx.Fee = unconvertedTx.Fee
  convertedTx.TxCount = unconvertedTx.TxCnt
  convertedTx.From = fmt.Sprintf("%x", unconvertedTx.From)
  convertedTx.To = fmt.Sprintf("%x", unconvertedTx.To)
  convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)

  return convertedTx
}

func ConvertAccTransaction(unconvertedTx *protocol.AccTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte) acctx {
  var convertedTx acctx
  //convertedTx.Header = fmt.Sprintf("%x", unconvertedTx.Header)
  convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
  convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
  convertedTx.Fee = unconvertedTx.Fee
  convertedTx.Issuer = fmt.Sprintf("%x", unconvertedTx.Issuer)
  convertedTx.PubKey = fmt.Sprintf("%x", unconvertedTx.PubKey)
  convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)

  return convertedTx
}

func ConvertConfigTransaction(unconvertedTx *protocol.ConfigTx, unconvertedBlockHash [32]byte, unconvertedTxHash [32]byte) configtx {
  var convertedTx configtx
  //convertedTx.Header = fmt.Sprintf("%x", unconvertedTx.Header)
  convertedTx.Hash = fmt.Sprintf("%x", unconvertedTxHash)
  convertedTx.BlockHash = fmt.Sprintf("%x", unconvertedBlockHash)
  convertedTx.Id = unconvertedTx.Id
  convertedTx.Fee = unconvertedTx.Fee
  convertedTx.Payload = unconvertedTx.Payload
  convertedTx.TxCount = unconvertedTx.TxCnt
  convertedTx.Signature = fmt.Sprintf("%x", unconvertedTx.Sig)

  return convertedTx
}

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
  "math/big"
  _ "github.com/mchetelat/bazo_miner/miner"
	"github.com/mchetelat/bazo_miner/p2p"
	"github.com/mchetelat/bazo_miner/protocol"
)

var allBlockHeaders []*protocol.SPVHeader
var allBlocks []*protocol.Block
var logger *log.Logger

var block1 *protocol.Block

func runDB() {
  if !checkEmptyDB() {
    fmt.Println("The database is empty!")
    fmt.Println("Loading all blocks and saving them to the database...")
    loadAllBlocks()
    fmt.Println("All blocks loaded and saved to the database!")
  } else {
    fmt.Println("The database is not empty!")
    fmt.Println("Loading all blocks without saving them to the database...")
    loadAllBlocksWithoutWrite()
    fmt.Println("All blocks loaded!")
  }
  for 0 < 1 {
    time.Sleep(time.Second * 30)
    fmt.Println("Refreshing State...")
    refreshState()
    fmt.Println("State refreshed!")
  }
}

func loadAllBlocksWithoutWrite() {
  block := reqBlock(nil)
  allBlocks = append(allBlocks, block)
  prevHash := block.PrevHash

  for block.Hash != [32]byte{} {
    block = reqBlock(prevHash[:])
    allBlocks = append(allBlocks, block)
    prevHash = block.PrevHash
  }

  allBlocks = invertBlockArray(allBlocks)
}

func refreshState() {
	var newBlocks []*protocol.Block
	newBlocks = getNewBlocks(reqBlock(nil), allBlocks[len(allBlocks)-1], newBlocks)

  //newBlocks = invertBlockArray(newBlocks)

  for _, oneBlock := range newBlocks{

    for _, accTxHash := range oneBlock.AccTxData{
      accTx := reqTx(p2p.ACCTX_REQ, accTxHash)
      convertedTx := ConvertAccTransaction(accTx.(*protocol.AccTx), oneBlock.Hash, accTxHash)

      fmt.Printf("Writing Transaction: %s\n", convertedTx.Hash)
      WriteAccTx(convertedTx)
    }

    for _, fundsTxHash := range oneBlock.FundsTxData{
      fundsTx := reqTx(p2p.FUNDSTX_REQ, fundsTxHash)
      convertedTx := ConvertFundsTransaction(fundsTx.(*protocol.FundsTx), oneBlock.Hash, fundsTxHash)

      fmt.Printf("Writing Transaction: %s\n", convertedTx.Hash)
      WriteFundsTx(convertedTx)
    }

    for _, configTxHash := range oneBlock.ConfigTxData{
      configTx := reqTx(p2p.CONFIGTX_REQ, configTxHash)
      convertedTx := ConvertConfigTransaction(configTx.(*protocol.ConfigTx), oneBlock.Hash, configTxHash)

      fmt.Printf("Writing Transaction: %s\n", convertedTx.Hash)
      WriteConfigTx(convertedTx)
    }
    //convert block
    convertedBlock := ConvertBlock(oneBlock)
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

    for _, accTxHash := range oneBlock.AccTxData{
      accTx := reqTx(p2p.ACCTX_REQ, accTxHash)
      convertedTx := ConvertAccTransaction(accTx.(*protocol.AccTx), oneBlock.Hash, accTxHash)
      accountHashBytes := SerializeHashContent(accTx.(*protocol.AccTx).PubKey)
      accountHash := fmt.Sprintf("%x", accountHashBytes)

      fmt.Printf("Writing Transaction: %s\n", convertedTx.Hash)
      WriteAccountWithAddress(convertedTx, accountHash)
      WriteAccTx(convertedTx)
    }

    for _, fundsTxHash := range oneBlock.FundsTxData{
      fundsTx := reqTx(p2p.FUNDSTX_REQ, fundsTxHash)
      convertedTx := ConvertFundsTransaction(fundsTx.(*protocol.FundsTx), oneBlock.Hash, fundsTxHash)

      fmt.Printf("Writing Transaction: %s\n", convertedTx.Hash)
      UpdateAccountData(convertedTx)
      WriteFundsTx(convertedTx)
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

func FetchOpenTx(txHash string){
  //sampleHash := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
  var txByteHash [32]byte

  txHashInt, _ := new(big.Int).SetString(txHash, 16)
  copy(txByteHash[:], txHashInt.Bytes())

  openTx := reqTx(p2p.FUNDSTX_REQ, txByteHash)
  convertedTx := ConvertOpenFundsTransaction(openTx.(*protocol.FundsTx), txByteHash)
  WriteOpenFundsTx(convertedTx)
}

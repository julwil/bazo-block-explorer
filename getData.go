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
  //allBlocks = invertBlockArray(allBlocks)
  for _, block := range allBlocks{
    //convert block
    convertedBlock := ConvertBlock(block)
    //write to db here
    fmt.Printf("Writing Block: %s\n", convertedBlock.Hash)
    WriteBlock(convertedBlock)
  }

}

//Update allBlockHeaders to the latest header
func refreshState() {
	var newBlockHeaders []*protocol.SPVHeader
	newBlockHeaders = getNewBlockHeaders(reqSPVHeader(nil), allBlockHeaders[len(allBlockHeaders)-1], newBlockHeaders)
	allBlockHeaders = append(allBlockHeaders, newBlockHeaders...)
}

//Get new blockheaders recursively
func getNewBlockHeaders(latest *protocol.SPVHeader, eldest *protocol.SPVHeader, list []*protocol.SPVHeader) []*protocol.SPVHeader {
	if latest.Hash != eldest.Hash {
		ancestor := reqSPVHeader(latest.PrevHash[:])
		list = getNewBlockHeaders(ancestor, eldest, list)
		list = append(list, latest)
	}
  return list
}

func loadAllBlockHeaders() {
	//If no blockhash as param is defined, the last SPVHeader is given back
	spvHeader := reqSPVHeader(nil)
	allBlockHeaders = append(allBlockHeaders, spvHeader)
	prevHash := spvHeader.PrevHash

	for spvHeader.Hash != [32]byte{} {
		spvHeader = reqSPVHeader(prevHash[:])
		allBlockHeaders = append(allBlockHeaders, spvHeader)
		prevHash = spvHeader.PrevHash
	}
}

func loadAllBlocks() {
  // copy hash from /testheader to get most recent one
  goodhash := [32]byte{0, 115, 1, 136, 42, 47, 27, 96, 119, 30, 228, 141, 78, 1, 93, 196, 179, 199, 55, 118, 129, 131, 239, 133, 111, 216, 172, 55, 92, 82, 192, 206}
  //genesishash := [32]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, }
  block := reqBlock(goodhash)
  allBlocks = append(allBlocks, block)
  prevHash := block.PrevHash

  for block.Hash != [32]byte{} {
    block = reqBlock(prevHash)
    allBlocks = append(allBlocks, block)
    prevHash = block.PrevHash
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

func reqBlock(blockHash [32]byte) (block *protocol.Block) {

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

func reqTx(txType uint8, txHash [32]byte) (tx protocol.Transaction) {

	conn := Connect(p2p.BOOTSTRAP_SERVER)

	packet := p2p.BuildPacket(txType, txHash[:])
	conn.Write(packet)

	header, payload, err := rcvData(conn)
	if err != nil {
		logger.Printf("Disconnected: %v\n", err)
		return
	}

	switch header.TypeID {
	case p2p.ACCTX_RES:
		var accTx *protocol.AccTx
		accTx = accTx.Decode(payload)
		tx = accTx
	case p2p.CONFIGTX_RES:
		var configTx *protocol.ConfigTx
		configTx = configTx.Decode(payload)
		tx = configTx
	case p2p.FUNDSTX_RES:
		var fundsTx *protocol.FundsTx
		fundsTx = fundsTx.Decode(payload)
		tx = fundsTx
	}

	conn.Close()

	return tx
}

func reqSPVHeader(blockHash []byte) (spvHeader *protocol.SPVHeader) {

	conn := Connect(p2p.BOOTSTRAP_SERVER)

	packet := p2p.BuildPacket(p2p.BLOCK_HEADER_REQ, blockHash)
	conn.Write(packet)

	header, payload, err := rcvData(conn)
	if err != nil {
		logger.Printf("Disconnected: %v\n", err)
		return
	}

	if header.TypeID == p2p.BlOCK_HEADER_RES {
		spvHeader = spvHeader.SPVDecode(payload)
	}

	conn.Close()

	return spvHeader
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
  for _, hash := range unconvertedBlock.AccTxData {
    convertedTxHash = fmt.Sprintf("%x", hash)
    convertedBlock.AccTxData = append(convertedBlock.AccTxData, convertedTxHash)
  }
  for _, hash := range unconvertedBlock.FundsTxData {
    convertedTxHash = fmt.Sprintf("%x", hash)
    convertedBlock.ConfigTxData = append(convertedBlock.ConfigTxData, convertedTxHash)
  }
  return convertedBlock
}

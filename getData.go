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
  "github.com/mchetelat/bazo_miner/miner"
	"github.com/mchetelat/bazo_miner/p2p"
	"github.com/mchetelat/bazo_miner/protocol"
)

var allBlockHeaders []*protocol.SPVHeader
var logger *log.Logger

func main()  {
  initState()
}

func initState() {
	loadAllBlockHeaders()
	allBlockHeaders = miner.InvertSPVHeaderSlice(allBlockHeaders)
  for _, header := range allBlockHeaders{
    fmt.Printf("%x\n", header.Hash)
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

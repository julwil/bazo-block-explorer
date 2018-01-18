package data

import (
  "github.com/lucBoillat/BazoBlockExplorer/utilities"
  "net"
  "fmt"
  "log"
  "time"
  "errors"
  "bufio"
	"github.com/mchetelat/bazo_miner/p2p"
  "github.com/mchetelat/bazo_miner/miner"
	"github.com/mchetelat/bazo_miner/protocol"
)

var newestBlock *protocol.Block
var logger *log.Logger
var block1 *protocol.Block

func RunDB() {

  saveInitialParameters()
  loadAllBlocks()

  for 0 < 1 {
    time.Sleep(time.Second * 60)
    RefreshState()
  }
}

func saveInitialParameters()  {
  var convertedParameters utilities.Systemparams
  parameters := miner.NewDefaultParameters()

  convertedParameters.Timestamp = time.Now().Unix()
  convertedParameters.BlockSize = parameters.Block_size
  convertedParameters.DiffInterval = parameters.Diff_interval
  convertedParameters.MinFee = parameters.Fee_minimum
  convertedParameters.BlockInterval = parameters.Block_interval
  convertedParameters.BlockReward = parameters.Diff_interval

  WriteParameters(convertedParameters)
}

func loadAllBlocks() {
  block := reqBlock(nil)
  newestBlock = block
  SaveBlockAndTransactions(block)
  prevHash := block.PrevHash

  for block.Hash != [32]byte{} {
    block = reqBlock(prevHash[:])
    SaveBlockAndTransactions(block)
    prevHash = block.PrevHash
  }
  RemoveRootFromDB()
  UpdateTotals()
  fmt.Println("All Blocks Loaded!")
}

func RefreshState() {
  fmt.Println("Refreshing State...")
  block := reqBlock(nil)
  prevHash := block.PrevHash
  tempBlock := block

  if block.Hash == newestBlock.Hash {
    //No new Blocks
    RemoveRootFromDB()
    UpdateTotals()
    return

  } else if prevHash == newestBlock.Hash {
    //One new Block
    SaveBlockAndTransactions(block)
    newestBlock = block

    RemoveRootFromDB()
    UpdateTotals()
    return

  } else if block.Hash != newestBlock.Hash {
    //Multiple new Blocks
    SaveBlockAndTransactions(block)
  }

  for block.PrevHash != newestBlock.Hash {
    block = reqBlock(prevHash[:])
    prevHash = block.PrevHash
    SaveBlockAndTransactions(block)
  }
  newestBlock = tempBlock

  RemoveRootFromDB()
  UpdateTotals()
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

	return header, payload, nil
}

func SaveBlockAndTransactions(oneBlock *protocol.Block)  {
  for _, accTxHash := range oneBlock.AccTxData{
    accTx := reqTx(p2p.ACCTX_REQ, accTxHash)
    convertedTx := utilities.ConvertAccTransaction(accTx.(*protocol.AccTx), oneBlock.Hash, accTxHash)
    accountHashBytes := utilities.SerializeHashContent(accTx.(*protocol.AccTx).PubKey)
    accountHash := fmt.Sprintf("%x", accountHashBytes)

    WriteAccountWithAddress(convertedTx, accountHash)
    WriteAccTx(convertedTx)
  }

  for _, fundsTxHash := range oneBlock.FundsTxData{
    fundsTx := reqTx(p2p.FUNDSTX_REQ, fundsTxHash)
    convertedTx := utilities.ConvertFundsTransaction(fundsTx.(*protocol.FundsTx), oneBlock.Hash, fundsTxHash)

    UpdateAccountData(convertedTx)
    WriteFundsTx(convertedTx)
  }

  for _, configTxHash := range oneBlock.ConfigTxData{
    configTx := reqTx(p2p.CONFIGTX_REQ, configTxHash)
    convertedTx := utilities.ConvertConfigTransaction(configTx.(*protocol.ConfigTx), oneBlock.Hash, configTxHash)
    currentParams := ReturnNewestParameters()
    newParams := utilities.ExtractParameters(convertedTx, currentParams)

    WriteParameters(newParams)
    WriteConfigTx(convertedTx)
  }

  convertedBlock := utilities.ConvertBlock(oneBlock)
  WriteBlock(convertedBlock)
}

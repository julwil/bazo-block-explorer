package main

type block struct {
  Hash string
  PrevHash string
  Timestamp string
  MerkleRoot string
  Beneficiary string
  NrFundTx uint16
  NrAccTx uint16
  NrConfigTx uint8
  FundsTxDataString string
  FundsTxData []string
  AccTxData []string
  ConfigTxData []string
}

//Include bock hash maybe?
type fundstx struct {
  Hash string
  Amount uint64
  Fee uint64
  TxCount uint32
  From string
  To string
  Signature string
}

type acctx struct {
  Hash string
  Issuer string
  Fee uint64
  PubKey string
  Signature string
}

type configtx struct {
  Hash string
	Id uint8
	Payload uint64
	Fee uint64
	TxCount uint8
	Signature string
}
//Include tx hashes mybe?
type account struct {
  Address string
  Balance uint64
  TxCount uint32
  /*
  FundsTxData []string
  AccTxData []string
  ConfigTxData []string
  */
}
//change to map
type systemparams struct {
  BSName string
  BlockSize int
  DIName string
  DiffInterval int
  MFName string
  MinFee int
  BIName string
  BlockInterval int
  BRName string
  BlockReward int
}

type blocksandtx struct {
  Blocks []block
  Txs []fundstx
}

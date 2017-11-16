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

type fundstx struct {
  Hash string
  Amount uint64
  Fee uint64
  TxCount uint32
  From string
  To string
  Signature string
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

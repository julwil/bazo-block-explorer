package main

import (
  "database/sql"
)

type block struct {
  Header string
  Hash string
  PrevHash string
  Timestamp int64
  MerkleRoot string
  Beneficiary string
  NrFundsTx uint16
  NrAccTx uint16
  NrConfigTx uint8
  FundsTxDataString sql.NullString
  FundsTxData []string
  AccTxDataString sql.NullString
  AccTxData []string
  ConfigTxDataString sql.NullString
  ConfigTxData []string
}

//Include bock hash maybe?
type fundstx struct {
  Hash string
  BlockHash string
  Amount uint64
  Fee uint64
  TxCount uint32
  From string
  To string
  Signature string
}

type acctx struct {
  Hash string
  BlockHash string
  Issuer string
  Fee uint64
  PubKey string
  Signature string
}

type configtx struct {
  Hash string
  BlockHash string
	Id uint8
	Payload uint64
	Fee uint64
	TxCount uint8
	Signature string
}
//Include tx hashes mybe?
type account struct {
  Hash string
  Address string
  Balance int64
  TxCount int32
  FundsTxData []string
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

type accountwithtxs struct {
  Account account
  Txs []fundstx
}

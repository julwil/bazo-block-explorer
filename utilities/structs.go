package utilities

import (
  "database/sql"
)

type Block struct {
  Header string
  Hash string
  PrevHash string
  Timestamp int64
  TimeString string
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
  UrlLevel string
}

type Fundstx struct {
  Hash string
  BlockHash string
  Amount uint64
  Fee uint64
  TxCount uint32
  From string
  To string
  Timestamp int64
  Signature string
  UrlLevel string
}

type Acctx struct {
  Hash string
  BlockHash string
  Issuer string
  Fee uint64
  PubKey string
  Timestamp int64
  Signature string
  UrlLevel string
}

type Configtx struct {
  Hash string
  BlockHash string
	Id uint8
	Payload uint64
	Fee uint64
	TxCount uint8
  Timestamp int64
	Signature string
  UrlLevel string
}

type Account struct {
  Hash string
  Address string
  Balance int64
  TxCount int32
  FundsTxData []string
  UrlLevel string
}

type JSONAccount struct {
	Address       [64]byte `json:"-"`
	AddressString string   `json:"address"`
	Balance       uint64   `json:"balance"`
	TxCnt         uint32   `json:"txCnt"`
	IsCreated     bool     `json:"isCreated"`
	IsRoot        bool     `json:"isRoot"`
}

type Systemparams struct {
  Timestamp int64
  BlockHash string
  BlockSize uint64
  DiffInterval uint64
  MinFee uint64
  BlockInterval uint64
  BlockReward uint64
  UrlLevel string
}

type Blocksandtx struct {
  Blocks []Block
  Txs []Fundstx
  UrlLevel string
}

type Blocksandurl struct {
  Blocks []Block
  UrlLevel string
}

type Fundsandurl struct {
  Txs []Fundstx
  UrlLevel string
}

type Accsandurl struct {
  Txs []Acctx
  UrlLevel string
}

type Configssandurl struct {
  Txs []Configtx
  UrlLevel string
}

type Accountsandurl struct {
  Accounts []Account
  UrlLevel string
}

type Emptyresponse struct {
  value int
  UrlLevel string
}

type Accountwithtxs struct {
  Account Account
  Txs []Fundstx
  UrlLevel string
}

type Stats struct {
  ChartData string
  TotalSupply int
  TotalNrAccounts int
  Parameters Systemparams
  UrlLevel string
}

type Serie struct {
  Name  string `json:"x"`
	Value int    `json:"value"`
}

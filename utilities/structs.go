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
  Nonce string
  Seed string
  HashedSeed string
  Height uint32
  NrFundsTx uint16
  NrAccTx uint16
  NrConfigTx uint8
  NrStakeTx uint16
  FundsTxDataString sql.NullString
  FundsTxData []string
  AccTxDataString sql.NullString
  AccTxData []string
  ConfigTxDataString sql.NullString
  ConfigTxData []string
  StakeTxDataString sql.NullString
  StakeTxData []string
  UrlLevel string
}

type Fundstx struct {
	Header string
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
  Header string
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
  Header string
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

type Staketx struct {
	Header string
  Hash string
  BlockHash string
  Timestamp int64
  Fee uint64
  IsStaking bool
  Account string
  Signature string
  UrlLevel string
}

type Account struct {
  Hash string
  Address string
  Balance int64
  TxCount int32
  FundsTxData []string
  IsStaking bool
  UrlLevel string
}

type JSONAccount struct {
  Code          int      `json:"code"`
  Message       string   `json:"message"`
  Content       []Content `json:"content"`
}

type Content struct {
  Name        string   `json:"name"`
  Detail      struct {
    Address   string   `json:"address"`
    Balance   int   `json:"balance"`
    TxCnt     int   `json:"txCnt"`
    IsCreated bool   `json:"isCreated"`
    IsRoot    bool   `json:"isRoot"`
  } `json:"detail"`
}

type Systemparams struct {
  Timestamp int64
  BlockHash string
  BlockSize uint64
  DiffInterval uint64
  MinFee uint64
  BlockInterval uint64
  BlockReward uint64
  StakingMin uint64
  WaitingMin uint64
  AcceptanceTimeDiff uint64
  SlashingWindowSize uint64
  SlashingReward uint64
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

type Configsandurl struct {
  Txs []Configtx
  UrlLevel string
}

type Stakesandurl struct {
  Txs []Staketx
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

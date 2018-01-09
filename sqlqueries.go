package main

import (
  "fmt"
  "database/sql"
  "github.com/lib/pq"
  "strings"
)

const (
  host = "localhost"
  port = 5432
  user = "postgres"
  //password = ""
  //use blockexplorertest1 for dummy db
  dbname = "blockexplorerdb"
)

var sqlStatement string
var db *sql.DB
var err error

func connectToDB() {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err = sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  err = db.Ping()
  if err != nil {
    panic(err)
  }
}

func setupDB()  {
  fmt.Println("Setting up Database...")
  dropTables()
  createTables()
  fmt.Println("Setup Complete")
}

func dropTables() {
  connectToDB()
  defer db.Close()
  fmt.Println("Dropping Tables...")
  sqlStatement := `drop table blocks;
                   drop table fundstx;
                   drop table acctx;
                   drop table configtx;
                   drop table accounts;
                   drop table openfundstx;`
  db.Exec(sqlStatement)
  fmt.Println("Dropped Tables")

}

func ReturnOneBlock(UrlHash string) block {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, prevhash, timestamp, timestring, merkleroot, beneficiary, nrfundstx, nracctx, nrconfigtx, fundstxdata, acctxdata, configtxdata FROM blocks WHERE hash = $1;`
  var returnedblock block
  row := db.QueryRow(sqlStatement, UrlHash)
  switch err := row.Scan(&returnedblock.Hash, &returnedblock.PrevHash, &returnedblock.Timestamp, &returnedblock.TimeString, &returnedblock.MerkleRoot, &returnedblock.Beneficiary, &returnedblock.NrFundsTx, &returnedblock.NrAccTx, &returnedblock.NrConfigTx, &returnedblock.FundsTxDataString, &returnedblock.AccTxDataString, &returnedblock.ConfigTxDataString)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
  case nil:
    if len(returnedblock.FundsTxDataString.String) > 0 {
      returnedblock.FundsTxData = strings.Split(returnedblock.FundsTxDataString.String[1:len(returnedblock.FundsTxDataString.String)-1], ",")
    }
    if len(returnedblock.AccTxDataString.String) > 0 {
      returnedblock.AccTxData = strings.Split(returnedblock.AccTxDataString.String[1:len(returnedblock.AccTxDataString.String)-1], ",")
    }
    if len(returnedblock.ConfigTxDataString.String) > 0 {
      returnedblock.ConfigTxData = strings.Split(returnedblock.ConfigTxDataString.String[1:len(returnedblock.ConfigTxDataString.String)-1], ",")
    }
    return returnedblock
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var block1 block
  return block1
}

func ReturnAllBlocks(UrlHash string) []block {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, timestamp, timestring, beneficiary, nrFundsTx, nrAccTx, nrConfigTx FROM blocks ORDER BY timestamp DESC LIMIT 100`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedrows := make([]block, 0)
  for rows.Next() {
    var returnedrow block
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Timestamp, &returnedrow.TimeString, &returnedrow.Beneficiary, &returnedrow.NrFundsTx, &returnedrow.NrAccTx, &returnedrow.NrConfigTx)
    //returnedrow.Timestamp = returnedrow.Timestamp[:19]
    if err != nil {
      panic(err)
    }
    returnedrows = append(returnedrows, returnedrow)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  return returnedrows
}

func ReturnOneFundsTx(UrlHash string) fundstx {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, blockhash, amount, fee, txcount, sender, recipient, signature FROM fundstx WHERE hash = $1;`
  var returnedrow fundstx
  row := db.QueryRow(sqlStatement, UrlHash)
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.BlockHash, &returnedrow.Amount, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.From, &returnedrow.To, &returnedrow.Signature)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
  case nil:
    return returnedrow
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var tx1 fundstx
  return tx1
}

func ReturnAllFundsTx(UrlHash string) []fundstx {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, amount, fee, txcount, sender, recipient, signature FROM fundstx`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedrows := make([]fundstx, 0)
  for rows.Next() {
    var returnedrow fundstx
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Amount, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.From, &returnedrow.To, &returnedrow.Signature)
    if err != nil {
      panic(err)
    }
    returnedrows = append(returnedrows, returnedrow)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  return returnedrows
}

func ReturnOneAccTx(UrlHash string) acctx {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, blockhash, issuer, fee, pubkey, signature FROM acctx WHERE hash = $1;`
  var returnedrow acctx
  row := db.QueryRow(sqlStatement, UrlHash)
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.BlockHash, &returnedrow.Issuer, &returnedrow.Fee, &returnedrow.PubKey, &returnedrow.Signature)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
  case nil:
    return returnedrow
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var tx1 acctx
  return tx1
}

func ReturnAllAccTx(UrlHash string) []acctx {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, issuer, fee, pubkey FROM acctx`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedrows := make([]acctx, 0)
  for rows.Next() {
    var returnedrow acctx
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Issuer, &returnedrow.Fee, &returnedrow.PubKey)
    if err != nil {
      panic(err)
    }
    returnedrows = append(returnedrows, returnedrow)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  return returnedrows
}

func ReturnOneConfigTx(UrlHash string) configtx {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, blockhash, id, payload, fee, txcount, signature FROM configtx WHERE hash = $1;`
  var returnedrow configtx
  row := db.QueryRow(sqlStatement, UrlHash)
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.BlockHash, &returnedrow.Id, &returnedrow.Payload, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.Signature)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
  case nil:
    return returnedrow
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var tx1 configtx
  return tx1
}

func ReturnAllConfigTx(UrlHash string) []configtx {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, id, payload, fee, txcount FROM configtx`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedrows := make([]configtx, 0)
  for rows.Next() {
    var returnedrow configtx
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Id, &returnedrow.Payload, &returnedrow.Fee, &returnedrow.TxCount)
    if err != nil {
      panic(err)
    }
    returnedrows = append(returnedrows, returnedrow)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  return returnedrows
}

func ReturnBlocksAndTransactions(UrlHash string) blocksandtx {
  var returnedBlocksAndTxs blocksandtx
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, timestamp, timestring, beneficiary, nrFundsTx, nrAccTx, nrConfigTx FROM blocks ORDER BY timestamp DESC LIMIT 6`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedblocks := make([]block, 0)
  for rows.Next() {
    var returnedrow block
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Timestamp, &returnedrow.TimeString, &returnedrow.Beneficiary, &returnedrow.NrFundsTx, &returnedrow.NrAccTx, &returnedrow.NrConfigTx)
    //returnedrow.Timestamp = returnedrow.Timestamp[:19]
    if err != nil {
      panic(err)
    }
    returnedblocks = append(returnedblocks, returnedrow)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }

  sqlStatement = `SELECT hash, amount, fee, txcount, sender, recipient, signature FROM fundstx`
  rows, err = db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedtxs := make([]fundstx, 0)
  for rows.Next() {
    var returnedrow fundstx
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Amount, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.From, &returnedrow.To, &returnedrow.Signature)
    if err != nil {
      panic(err)
    }
    returnedtxs = append(returnedtxs, returnedrow)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  returnedBlocksAndTxs.Blocks = returnedblocks
  returnedBlocksAndTxs.Txs = returnedtxs

  return returnedBlocksAndTxs
}

func WriteBlock(block block)  {
  connectToDB()
  defer db.Close()

  sqlStatement = `
    INSERT INTO blocks (hash, prevhash, timestamp, timestring, merkleroot, beneficiary, nrfundstx, nracctx, nrconfigtx, fundstxdata, acctxdata, configtxdata)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
  _, err = db.Exec(sqlStatement, block.Hash, block.PrevHash, block.Timestamp, block.TimeString, block.MerkleRoot, block.Beneficiary, block.NrFundsTx, block.NrAccTx, block.NrConfigTx, pq.Array(block.FundsTxData), pq.Array(block.AccTxData), pq.Array(block.ConfigTxData))
  if err != nil {
    panic(err)
  }
}

func WriteFundsTx(tx fundstx) {
  connectToDB()
  defer db.Close()

  sqlStatement = `
    INSERT INTO fundstx (hash, blockhash, amount, fee, txcount, sender, recipient, signature)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
  _, err = db.Exec(sqlStatement, tx.Hash, tx.BlockHash, tx.Amount, tx.Fee, tx.TxCount, tx.From, tx.To, tx.Signature)
  if err != nil {
    panic(err)
  }
}

func WriteAccTx(tx acctx) {
  connectToDB()
  defer db.Close()

  sqlStatement = `
    INSERT INTO acctx (hash, blockhash, fee, issuer, pubkey, signature)
    VALUES ($1, $2, $3, $4, $5, $6)`
  _, err = db.Exec(sqlStatement, tx.Hash, tx.BlockHash, tx.Fee, tx.Issuer, tx.PubKey, tx.Signature)
  if err != nil {
    panic(err)
  }
}

func WriteConfigTx(tx configtx) {
  connectToDB()
  defer db.Close()

  sqlStatement = `
    INSERT INTO configtx (hash, blockhash, id, payload, fee, txcount, signature)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`
  _, err = db.Exec(sqlStatement, tx.Hash, tx.BlockHash, tx.Id, tx.Payload, tx.Fee, tx.TxCount, tx.Signature)
  if err != nil {
    panic(err)
  }
}

func checkEmptyDB() bool {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT CASE WHEN EXISTS (SELECT * FROM blocks LIMIT 1) THEN 1 ELSE 0 END`
  var notEmpty bool
  row := db.QueryRow(sqlStatement)
  switch err := row.Scan(&notEmpty)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("No rows returned!")
  case nil:
    return notEmpty
  default:
    //on website 500 error maybe.
    panic(err)
  }
  return true
}

func WriteOpenFundsTx(tx fundstx) {
  connectToDB()
  defer db.Close()

  sqlStatement = `
    INSERT INTO openfundstx (hash, amount, fee, txcount, sender, recipient, signature)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`
  _, err = db.Exec(sqlStatement, tx.Hash, tx.Amount, tx.Fee, tx.TxCount, tx.From, tx.To, tx.Signature)
  if err != nil {
    panic(err)
  }
}

func ReturnOpenFundsTx(UrlHash string) fundstx {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, amount, fee, txcount, sender, recipient, signature FROM openfundstx WHERE hash = $1;`
  var returnedrow fundstx
  row := db.QueryRow(sqlStatement, UrlHash)
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.Amount, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.From, &returnedrow.To, &returnedrow.Signature)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("Transaction could not be found!")
  case nil:
    return returnedrow
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var tx1 fundstx
  return tx1
}

func UpdateAccountData(tx fundstx) {
  connectToDB()
  defer db.Close()

  sqlStatement := `UPDATE accounts SET balance = accounts.balance - $2, txcount = accounts.txcount + 1 WHERE hash = $1`
  totalAmount := tx.Amount + tx.Fee
  //totalCount := tx.TxCount + 1
  _, err = db.Exec(sqlStatement, tx.From, totalAmount)
  if err != nil {
    panic(err)
  }
  sqlStatement = `UPDATE accounts SET balance = accounts.balance + $2 WHERE hash = $1`
  _, err = db.Exec(sqlStatement, tx.To, tx.Amount)
  if err != nil {
    panic(err)
  }
}

func WriteAccountWithAddress(tx acctx, accountHash string) {
  connectToDB()
  defer db.Close()

  sqlStatement := `INSERT INTO accounts (hash, address, balance, txcount)
                    VALUES ($1, $2, $3, $4)`

  _, err = db.Exec(sqlStatement, accountHash, tx.PubKey, 0, 0)
  if err != nil {
    panic(err)
  }
}

func ReturnOneAccount(UrlHash string) accountwithtxs {
  var returnedData accountwithtxs

  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, address, balance, txcount FROM accounts WHERE hash = $1 OR address = $1`
  var returnedaccount account
  row := db.QueryRow(sqlStatement, UrlHash)
  switch err = row.Scan(&returnedaccount.Hash, &returnedaccount.Address, &returnedaccount.Balance, &returnedaccount.TxCount)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    //return fitting type
    returnedData.Account = returnedaccount
    return returnedData
  case nil:
    sqlStatement = `SELECT hash, amount, fee, txcount, sender, recipient FROM fundstx WHERE sender = $1 OR recipient = $1`
    rows, err := db.Query(sqlStatement, returnedaccount.Hash)
    if err != nil {
      panic(err)
    }
    defer rows.Close()
    returnedrows := make([]fundstx, 0)
    for rows.Next() {
      var returnedrow fundstx
      err = rows.Scan(&returnedrow.Hash, &returnedrow.Amount, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.From, &returnedrow.To)
      if err != nil {
        panic(err)
      }
      returnedrows = append(returnedrows, returnedrow)
    }
    err = rows.Err()
    if err != nil {
      panic(err)
    }
    returnedData.Account = returnedaccount
    returnedData.Txs = returnedrows
    return returnedData
  default:
    //on website 500 error maybe.
    panic(err)
  }
  return returnedData
}

func ReturnTopAccounts(UrlHash string) []account {
  connectToDB()
  defer db.Close()

  sqlStatement := `SELECT hash, address, balance, txcount FROM accounts ORDER BY balance DESC LIMIT 10`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedrows := make([]account, 0)
  for rows.Next() {
    var returnedrow account
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Address, &returnedrow.Balance, &returnedrow.TxCount)
    if err != nil {
      panic(err)
    }
    returnedrows = append(returnedrows, returnedrow)
  }
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  return returnedrows
}

func createTables() {
  connectToDB()
  defer db.Close()
  fmt.Println("Creating Tables...")

  sqlStatement :=   `create table blocks (
                    header bit(8),
                    hash char(64) primary key,
                    prevHash char(64) not null,
                    nonce char(16),
                    timestamp bigint not null,
                    timestring varchar(100) not null,
                    merkleRoot char(64) not null,
                    beneficiary char(64) not null,
                    nrFundsTx smallint not null,
                    nrAccTx smallint not null,
                    nrConfigTx smallint not null,
                    fundsTxData varchar(100)[],
                    accTxData varchar(100)[],
                    configTxData varchar(100)[]
                    );

                    create table fundstx (
                    header bit(8),
                    hash char(64) primary key,
                    blockhash char(64) not null,
                    amount bigint not null,
                    fee bigint not null,
                    txcount int not null,
                    sender char(64) not null,
                    recipient char(64) not null,
                    signature char(128) not null
                    );

                    create table openfundstx (
                    header bit(8),
                    hash char(64) primary key,
                    amount bigint not null,
                    fee bigint not null,
                    txcount int not null,
                    sender char(64) not null,
                    recipient char(64) not null,
                    signature char(128) not null
                    );

                    create table acctx(
                    header bit(8),
                    hash char(64) primary key,
                    blockhash char(64),
                    issuer char(64) not null,
                    fee bigint not null,
                    pubkey char(128) not null,
                    signature char(128) not null
                    );

                    create table configtx(
                    header bit(8),
                    hash char(64) primary key,
                    blockhash char(64),
                    id int not null,
                    payload bigint not null,
                    fee bigint not null,
                    txcount int not null,
                    signature char(128) not null
                    );

                    create table accounts(
                    hash char(64) primary key,
                    address char(128),
                    balance bigint not null,
                    txcount int not null
                    );`
                    db.Exec(sqlStatement)
  fmt.Println("Created Tables Successfully")
}

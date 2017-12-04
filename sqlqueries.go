package main

import (
  _ "io"
  "fmt"
  "net/http"
  _ "html/template"
  "database/sql"
  _ "github.com/lib/pq"
  "github.com/julienschmidt/httprouter"
  _ "strconv"
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

func ReturnOneBlock(params httprouter.Params) block {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement = `SELECT hash, prevhash, timestamp, merkleroot, beneficiary, nrfundtx, nracctx, nrconfigtx FROM blocks WHERE hash = $1;`
  var returnedblock block
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err := row.Scan(&returnedblock.Hash, &returnedblock.PrevHash, &returnedblock.Timestamp, &returnedblock.MerkleRoot, &returnedblock.Beneficiary, &returnedblock.NrFundsTx, &returnedblock.NrAccTx, &returnedblock.NrConfigTx)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("No rows returned!")
  case nil:
    //returnedblock.FundsTxData = strings.Split(returnedblock.FundsTxDataString[1:len(returnedblock.FundsTxDataString)-1], ",")
    return returnedblock
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var block1 block
  return block1
}

func ReturnAllBlocks(params httprouter.Params) []block {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement := `SELECT hash, timestamp, beneficiary, nrFundTx, nrAccTx, nrConfigTx FROM blocks`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedrows := make([]block, 0)
  for rows.Next() {
    var returnedrow block
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Timestamp, &returnedrow.Beneficiary, &returnedrow.NrFundsTx, &returnedrow.NrAccTx, &returnedrow.NrConfigTx)
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

func ReturnOneFundsTx(params httprouter.Params) fundstx {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement := `SELECT hash, amount, fee, txcount, sender, recipient, signature FROM fundstx WHERE hash = $1;`
  var returnedrow fundstx
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.Amount, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.From, &returnedrow.To, &returnedrow.Signature)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("No rows returned!")
  case nil:
    return returnedrow
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var tx1 fundstx
  return tx1
}

func ReturnAllFundsTx(params httprouter.Params) []fundstx {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

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

func ReturnOneAccTx(params httprouter.Params) acctx {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement := `SELECT hash, issuer, fee, pubkey, signature FROM acctx WHERE hash = $1;`
  var returnedrow acctx
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.Issuer, &returnedrow.Fee, &returnedrow.PubKey, &returnedrow.Signature)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("No rows returned!")
  case nil:
    return returnedrow
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var tx1 acctx
  return tx1
}

func ReturnAllAccTx(params httprouter.Params) []acctx {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

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

func ReturnOneConfigTx(params httprouter.Params) configtx {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement := `SELECT hash, id, payload, fee, txcount, signature FROM configtx WHERE hash = $1;`
  var returnedrow configtx
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.Id, &returnedrow.Payload, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.Signature)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("No rows returned!")
  case nil:
    return returnedrow
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var tx1 configtx
  return tx1
}

func ReturnAllConfigTx(params httprouter.Params) []configtx {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

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

func ReturnAccount(params httprouter.Params) account {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement := `SELECT address, balance, txcount FROM accounts WHERE address = $1;`
  var returnedaccount account
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedaccount.Address, &returnedaccount.Balance, &returnedaccount.TxCount)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("No rows returned!")
  case nil:
    return returnedaccount
  default:
    //on website 500 error maybe.
    panic(err)
  }
  var account1 account
  return account1
}

func ReturnSearchResult(r *http.Request) (block, fundstx) {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement := `SELECT hash, prevhash, timestamp, merkleroot, beneficiary, nrfundtx, nracctx, nrconfigtx, fundstxdata FROM blocks WHERE hash = $1;`
  var returnedblock block
  var returnedtx fundstx
  row := db.QueryRow(sqlStatement, r.PostFormValue("search-value"))
  switch err = row.Scan(&returnedblock.Hash, &returnedblock.PrevHash, &returnedblock.Timestamp, &returnedblock.MerkleRoot, &returnedblock.Beneficiary, &returnedblock.NrFundsTx, &returnedblock.NrAccTx, &returnedblock.NrConfigTx, &returnedblock.FundsTxDataString)
  err {
  case sql.ErrNoRows:
    sqlStatement = `SELECT hash, amount, fee, txcount, sender, recipient, signature FROM fundstx WHERE hash = $1;`
    row2 := db.QueryRow(sqlStatement, r.PostFormValue("search-value"))
    err = row2.Scan(&returnedtx.Hash, &returnedtx.Amount, &returnedtx.Fee, &returnedtx.TxCount, &returnedtx.From, &returnedtx.To, &returnedtx.Signature)
    return returnedblock, returnedtx
  case nil:
    returnedblock.FundsTxData = strings.Split(returnedblock.FundsTxDataString[1:len(returnedblock.FundsTxDataString)-1], ",")
    return returnedblock, returnedtx
  default:
    panic(err)
  }
}

func ReturnBlocksAndTransactions(params httprouter.Params) blocksandtx {
  var returnedBlocksAndTxs blocksandtx
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement := `SELECT hash, timestamp, beneficiary, nrFundTx, nrAccTx, nrConfigTx FROM blocks`
  rows, err := db.Query(sqlStatement)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  returnedblocks := make([]block, 0)
  for rows.Next() {
    var returnedrow block
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Timestamp, &returnedrow.Beneficiary, &returnedrow.NrFundsTx, &returnedrow.NrAccTx, &returnedrow.NrConfigTx)
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
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
    host, port, user, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()
  err = db.Ping()
  if err != nil {
    panic(err)
  }

  sqlStatement = `
    INSERT INTO blocks (hash, prevhash, timestamp, merkleroot, beneficiary, nrfundtx, nracctx, nrconfigtx)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
  _, err = db.Exec(sqlStatement, block.Hash, block.PrevHash, block.Timestamp, block.MerkleRoot, block.Beneficiary, block.NrFundsTx, block.NrAccTx, block.NrConfigTx)
  if err != nil {
    panic(err)
  }
}

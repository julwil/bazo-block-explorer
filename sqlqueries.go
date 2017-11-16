package main

import (
  _ "io"
  "fmt"
  _ "net/http"
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
  dbname = "blockexplorertest1"
)

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

  sqlStatement := `SELECT hash, prevhash, timestamp, merkleroot, beneficiary, nrfundtx, nracctx, nrconfigtx, fundstxdata FROM blocks WHERE hash = $1;`
  var returnedblock block
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err := row.Scan(&returnedblock.Hash, &returnedblock.PrevHash, &returnedblock.Timestamp, &returnedblock.MerkleRoot, &returnedblock.Beneficiary, &returnedblock.NrFundTx, &returnedblock.NrAccTx, &returnedblock.NrConfigTx, &returnedblock.FundsTxDataString)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("No rows returned!")
  case nil:
    returnedblock.FundsTxData = strings.Split(returnedblock.FundsTxDataString[1:len(returnedblock.FundsTxDataString)-1], ",")
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
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Timestamp, &returnedrow.Beneficiary, &returnedrow.NrFundTx, &returnedrow.NrAccTx, &returnedrow.NrConfigTx)
    returnedrow.Timestamp = returnedrow.Timestamp[:19]
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

func ReturnOneTransaction(params httprouter.Params) fundstx {
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

func ReturnAllTransactions(params httprouter.Params) []fundstx {
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

func ReturnSearchResult(params httprouter.Params) (block, fundstx) {
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
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedblock.Hash, &returnedblock.PrevHash, &returnedblock.Timestamp, &returnedblock.MerkleRoot, &returnedblock.Beneficiary, &returnedblock.NrFundTx, &returnedblock.NrAccTx, &returnedblock.NrConfigTx, &returnedblock.FundsTxDataString)
  err {
  case sql.ErrNoRows:
    sqlStatement = `SELECT hash, amount, fee, txcount, sender, recipient, signature FROM fundstx WHERE hash = $1;`
    row2 := db.QueryRow(sqlStatement, params.ByName("hash"))
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
    err = rows.Scan(&returnedrow.Hash, &returnedrow.Timestamp, &returnedrow.Beneficiary, &returnedrow.NrFundTx, &returnedrow.NrAccTx, &returnedrow.NrConfigTx)
    returnedrow.Timestamp = returnedrow.Timestamp[:19]
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

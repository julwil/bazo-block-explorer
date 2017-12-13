package main

import (
  _ "io"
  "fmt"
  "net/http"
  _ "html/template"
  "database/sql"
  "github.com/lib/pq"
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

  sqlStatement := `SELECT hash, prevhash, timestamp, merkleroot, beneficiary, nrfundstx, nracctx, nrconfigtx, fundstxdata, acctxdata, configtxdata FROM blocks WHERE hash = $1;`
  var returnedblock block
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err := row.Scan(&returnedblock.Hash, &returnedblock.PrevHash, &returnedblock.Timestamp, &returnedblock.MerkleRoot, &returnedblock.Beneficiary, &returnedblock.NrFundsTx, &returnedblock.NrAccTx, &returnedblock.NrConfigTx, &returnedblock.FundsTxDataString, &returnedblock.AccTxDataString, &returnedblock.ConfigTxDataString)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    fmt.Printf("No rows returned!")
  case nil:
    if len(returnedblock.FundsTxDataString.String) > 0 {
      returnedblock.FundsTxData = strings.Split(returnedblock.FundsTxDataString.String[:len(returnedblock.FundsTxDataString.String)], ",")
    }
    if len(returnedblock.AccTxDataString.String) > 0 {
      returnedblock.AccTxData = strings.Split(returnedblock.AccTxDataString.String[:len(returnedblock.AccTxDataString.String)], ",")
    }
    if len(returnedblock.ConfigTxDataString.String) > 0 {
      returnedblock.ConfigTxData = strings.Split(returnedblock.ConfigTxDataString.String[:len(returnedblock.ConfigTxDataString.String)], ",")
    }
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

  sqlStatement := `SELECT hash, timestamp, beneficiary, nrFundsTx, nrAccTx, nrConfigTx FROM blocks ORDER BY timestamp DESC LIMIT 100`
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

  sqlStatement := `SELECT hash, blockhash, amount, fee, txcount, sender, recipient, signature FROM fundstx WHERE hash = $1;`
  var returnedrow fundstx
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.BlockHash, &returnedrow.Amount, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.From, &returnedrow.To, &returnedrow.Signature)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    return returnedrow
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

  sqlStatement := `SELECT hash, blockhash, issuer, fee, pubkey, signature FROM acctx WHERE hash = $1;`
  var returnedrow acctx
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.BlockHash, &returnedrow.Issuer, &returnedrow.Fee, &returnedrow.PubKey, &returnedrow.Signature)
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

  sqlStatement := `SELECT hash, blockhash, id, payload, fee, txcount, signature FROM configtx WHERE hash = $1;`
  var returnedrow configtx
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedrow.Hash, &returnedrow.BlockHash, &returnedrow.Id, &returnedrow.Payload, &returnedrow.Fee, &returnedrow.TxCount, &returnedrow.Signature)
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

  sqlStatement := `SELECT hash, prevhash, timestamp, merkleroot, beneficiary, nrfundstx, nracctx, nrconfigtx, fundstxdata FROM blocks WHERE hash = $1;`
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
    if len(returnedblock.FundsTxDataString.String) > 0 {
      returnedblock.FundsTxData = strings.Split(returnedblock.FundsTxDataString.String[1:len(returnedblock.FundsTxDataString.String)-1], ",")
    }
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

  sqlStatement := `SELECT hash, timestamp, beneficiary, nrFundsTx, nrAccTx, nrConfigTx FROM blocks`
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
    INSERT INTO blocks (hash, prevhash, timestamp, merkleroot, beneficiary, nrfundstx, nracctx, nrconfigtx, fundstxdata, acctxdata, configtxdata)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
  _, err = db.Exec(sqlStatement, block.Hash, block.PrevHash, block.Timestamp, block.MerkleRoot, block.Beneficiary, block.NrFundsTx, block.NrAccTx, block.NrConfigTx, pq.Array(block.FundsTxData), pq.Array(block.AccTxData), pq.Array(block.ConfigTxData))
  if err != nil {
    panic(err)
  }
}

func WriteFundsTx(tx fundstx) {
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
    INSERT INTO fundstx (hash, blockhash, amount, fee, txcount, sender, recipient, signature)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
  _, err = db.Exec(sqlStatement, tx.Hash, tx.BlockHash, tx.Amount, tx.Fee, tx.TxCount, tx.From, tx.To, tx.Signature)
  if err != nil {
    panic(err)
  }
}

func WriteAccTx(tx acctx) {
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
    INSERT INTO acctx (hash, blockhash, fee, issuer, pubkey, signature)
    VALUES ($1, $2, $3, $4, $5, $6)`
  _, err = db.Exec(sqlStatement, tx.Hash, tx.BlockHash, tx.Fee, tx.Issuer, tx.PubKey, tx.Signature)
  if err != nil {
    panic(err)
  }
}

func WriteConfigTx(tx configtx) {
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
    INSERT INTO configtx (hash, blockhash, id, payload, fee, txcount, signature)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`
  _, err = db.Exec(sqlStatement, tx.Hash, tx.BlockHash, tx.Id, tx.Payload, tx.Fee, tx.TxCount, tx.Signature)
  if err != nil {
    panic(err)
  }
}

func checkEmptyDB() bool {
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
    INSERT INTO openfundstx (hash, amount, fee, txcount, sender, recipient, signature)
    VALUES ($1, $2, $3, $4, $5, $6, $7)`
  _, err = db.Exec(sqlStatement, tx.Hash, tx.Amount, tx.Fee, tx.TxCount, tx.From, tx.To, tx.Signature)
  if err != nil {
    panic(err)
  }
}

func ReturnOpenFundsTx(params httprouter.Params) fundstx {
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

  sqlStatement := `SELECT hash, amount, fee, txcount, sender, recipient, signature FROM openfundstx WHERE hash = $1;`
  var returnedrow fundstx
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
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

func WriteAccountData(tx fundstx) {
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

  sqlStatement := `INSERT INTO accounts (hash, address, balance, txcount)
                    VALUES ($1, $2, $3, $4)
                    ON CONFLICT (hash) DO UPDATE SET balance = accounts.balance - $3, txcount = accounts.txcount + 1`
  emptyString := ""
  totalAmount := tx.Amount + tx.Fee
  totalCount := tx.TxCount + 1
  _, err = db.Exec(sqlStatement, tx.From, emptyString, totalAmount, totalCount)
  if err != nil {
    panic(err)
  }
  sqlStatement = `INSERT INTO accounts (hash, address, balance, txcount)
                    VALUES ($1, $2, $3, 0)
                    ON CONFLICT (hash) DO UPDATE SET balance = accounts.balance + $3`
  _, err = db.Exec(sqlStatement, tx.To, emptyString, tx.Amount)
  if err != nil {
    panic(err)
  }
}

func ReturnOneAccount(params httprouter.Params) accountwithtxs {
  var returnedData accountwithtxs

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

  sqlStatement := `SELECT hash, address, balance, txcount FROM accounts WHERE hash = $1;`
  var returnedaccount account
  row := db.QueryRow(sqlStatement, params.ByName("hash"))
  switch err = row.Scan(&returnedaccount.Hash, &returnedaccount.Address, &returnedaccount.Balance, &returnedaccount.TxCount)
  err {
  case sql.ErrNoRows:
    //on website 404 would be more suitable maybe
    //return fitting type
    returnedData.Account = returnedaccount
    return returnedData
  case nil:
    sqlStatement = `SELECT hash, amount, fee, txcount, sender, recipient FROM fundstx WHERE sender = $1 OR recipient = $1`
    rows, err := db.Query(sqlStatement, params.ByName("hash"))
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

func ReturnTopAccounts(params httprouter.Params) []account {
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

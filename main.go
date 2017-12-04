package main

import (
  _ "io"
  _ "fmt"
  _ "os/exec"
  "net/http"
  "html/template"
  _ "github.com/julienschmidt/httprouter"
  _ "github.com/dgrijalva/jwt-go"
  _ "strconv"
)

var tpl *template.Template

var values systemparams

const BLOCKSIZE string = "1"
const DIFF_INTERV string = "2"
const MIN_FEE string = "3"
const BLOCK_INTERV string = "4"
const BLOCK_REWARD string = "5"

func init() {
  tpl = template.Must(template.ParseGlob("static/src/*.gohtml"))
  values.BlockSize = 0
  values.BSName = "Block Size"
  values.DiffInterval = 0
  values.DIName = "Difficulty Interval"
  values.MinFee = 0
  values.MFName = "Minimum Fee"
  values.BlockInterval = 0
  values.BIName = "Block Interval"
  values.BlockReward = 0
  values.BRName = "Block Reward"
}

func main() {
  router := initializeRouter()
  //fill db with blocks, make sure to use empty table, otherwise you get an error
  initState()
  http.ListenAndServe(":8080", router)
}

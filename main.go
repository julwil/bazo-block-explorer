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

func init() {
  tpl = template.Must(template.ParseGlob("static/src/*.gohtml"))
  values.BlockSize = 1
  values.BSName = "Block Size"
  values.DiffInterval = 1
  values.DIName = "Difficulty Interval"
  values.MinFee = 1
  values.MFName = "Minimum Fee"
  values.BlockInterval = 1
  values.BIName = "Block Interval"
  values.BlockReward = 1
  values.BRName = "Block Reward"
}

func main() {
  router := initializeRouter()
  http.ListenAndServe(":8080", router)
}

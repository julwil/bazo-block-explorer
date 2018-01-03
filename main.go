package main

import (
  _"os"
  "fmt"
  "net/http"
  "html/template"
  _ "github.com/dgrijalva/jwt-go"
)

var tpl *template.Template
var values systemparams


func init() {
  tpl = template.Must(template.ParseGlob("static/src/*"))
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

  setupDB()
  //go runDB()
  fmt.Println("Listening...")
  http.ListenAndServe(":8080", router)
}

package main

import (
  "fmt"
  "net/http"
  "html/template"
  _ "github.com/dgrijalva/jwt-go"
)

var tpl *template.Template

func init() {
  tpl = template.Must(template.ParseGlob("static/src/*.gohtml"))
}

func main() {
  router := initializeRouter()

  //go runDB()
  //loadAllBlocks()
  fmt.Println("Listening...")
  http.ListenAndServe(":8080", router)
}

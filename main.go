package main

import (
  _"os"
  "fmt"
  "net/http"
  "html/template"
)

var tpl *template.Template

func init() {
  tpl = template.Must(template.ParseGlob("source/html/*"))
}

func main() {
  router := initializeRouter()

  //setupDB()
  //go runDB()
  fmt.Println("Listening...")
  http.ListenAndServe(":8080", router)
}

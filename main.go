package main

import (
  _"os"
  "github.com/lucBoillat/BazoBlockExplorer/router"
  "github.com/lucBoillat/BazoBlockExplorer/data"
  "fmt"
  "net/http"
)

func main() {
  requestRouter := router.InitializeRouter()

  data.SetupDB()
  //go data.RunDB()
  fmt.Println("Listening...")
  http.ListenAndServe(":8080", requestRouter)
}

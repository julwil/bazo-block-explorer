package main

import (
  _"os"
  "BazoBlockExplorer/router"
  "fmt"
  "net/http"
)

func main() {
  requestRouter := router.InitializeRouter()

  //setupDB()
  //go runDB()
  fmt.Println("Listening...")
  http.ListenAndServe(":8080", requestRouter)
}

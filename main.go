package main

import (
  "os"
  "github.com/lucBoillat/BazoBlockExplorer/router"
  "github.com/lucBoillat/BazoBlockExplorer/data"
  "fmt"
  "net/http"
)

func main() {
  if len(os.Args) >= 3 {
    requestRouter := router.InitializeRouter()

    //drops all tables in the database and creates them again
    data.SetupDB(os.Args[2], os.Args[3])
    //retrieves data from the blockchain and stores it in the database
    //go data.RunDB()
    //starts the router
    fmt.Println("Listening...")
    http.ListenAndServe(os.Args[1], requestRouter)

  } else {
    fmt.Println("Not enough arguments!")
    fmt.Println("./BazoBlockExplorer <<:PORT>> <<db_username>> <<password>>")
  }
}

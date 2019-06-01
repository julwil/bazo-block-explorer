package main

import (
	"fmt"
	"github.com/bazo-blockchain/bazo-block-explorer/data"
	"github.com/bazo-blockchain/bazo-block-explorer/router"
	"github.com/bazo-blockchain/bazo-client/network"
	"github.com/bazo-blockchain/bazo-client/util"
	"github.com/bazo-blockchain/bazo-miner/p2p"
	"net/http"
	"os"
)

func main() {
	p2p.InitLogging()
	util.Config = util.LoadConfiguration()
	network.Init(p2p.MINER_PING)

	if len(os.Args) == 5 {
		requestRouter := router.InitializeRouter()

		//drops all tables in the database and creates them again
		data.SetupDB(os.Args[3], os.Args[4])

		if os.Args[1] == "data" {
			//retrieves data from the blockchain and stores it in the database
			go data.RunDB()
		}
		//starts the router
		fmt.Println("Listening...")
		http.ListenAndServe(os.Args[2], requestRouter)

	} else {
		fmt.Println("Incorrect number of arguments")
		fmt.Println("./BazoBlockExplorer <<data or nodata>> <<:WEB_PORT>> <<db_username>> <<db_password>>")
	}
}

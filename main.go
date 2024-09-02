package main

import (
	"fmt"
	"github.com/way365/bazo-block-explorer/data"
	"github.com/way365/bazo-block-explorer/router"
	"net/http"
	"os"
)

func main() {
	//time.Sleep(10 * time.Second)
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
		err := http.ListenAndServe(os.Args[2], requestRouter)
		if err != nil {
			fmt.Printf("%v", err)
		}

	} else {
		fmt.Println("Incorrect number of arguments")
		fmt.Println("./BazoBlockExplorer <<data or nodata>> <<:WEB_PORT>> <<db_username>> <<db_password>>")
	}
}

package main

import (
	"log"

	"github.com/felipehfs/rpgapi/config"
	"github.com/gorilla/mux"
)

func main() {
	conn, err := config.SetupDatabase(config.Development)
	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter()
	server := config.NewServer(conn, router)
	defer conn.Close()
	server.Start("8083")
}

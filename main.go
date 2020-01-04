package main

import (
	"log"
	"net/http"

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
	server.Routes()
	defer conn.Close()

	http.ListenAndServe(":8083", server)
}

package config

import (
	"database/sql"
	"net/http"

	"github.com/felipehfs/rpgapi/controllers"
	"github.com/gorilla/mux"
)

// Server represents the resources of the api
type Server struct {
	DB  *sql.DB
	Mux *mux.Router
}

// NewServer instantiate the server
func NewServer(db *sql.DB, router *mux.Router) *Server {
	return &Server{
		DB:  db,
		Mux: router,
	}
}

// Start setups the server with routes and port
func (server *Server) Start(port string) {
	server.Mux.HandleFunc("/api/characters", controllers.CreateCharacter(server.DB)).Methods("POST")
	server.Mux.HandleFunc("/api/characters", controllers.ReadCharacter(server.DB)).Methods("GET")
	server.Mux.HandleFunc("/api/characters/{id}", controllers.UpdateCharacter(server.DB)).Methods("PUT")
	server.Mux.HandleFunc("/api/characters/{id}", controllers.RemoveCharacter(server.DB)).Methods("DELETE")
	server.Mux.HandleFunc("/api/characters/{id}", controllers.GetByIDCharacter(server.DB)).Methods("GET")

	http.ListenAndServe(":"+port, server)
}

func (server Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.Mux.ServeHTTP(w, r)
}

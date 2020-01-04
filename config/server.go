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

// Routes register all routes
func (server *Server) Routes() {
	server.Mux.HandleFunc("/api/characters", controllers.CreateCharacter(server.DB)).Methods("POST")
	server.Mux.HandleFunc("/api/characters", controllers.ReadCharacter(server.DB)).Methods("GET")
	server.Mux.HandleFunc("/api/characters/{id}", controllers.UpdateCharacter(server.DB)).Methods("PUT")
}

func (server Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.Mux.ServeHTTP(w, r)
}

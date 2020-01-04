package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/felipehfs/rpgapi/models"
	"github.com/felipehfs/rpgapi/repositories"
	"github.com/gorilla/mux"
)

// CreateCharacter saves the caracter
func CreateCharacter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var character models.Character

		if err := json.NewDecoder(r.Body).Decode(&character); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := character.IsValid(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		repository := repositories.NewCharacterRepository(db)
		id, err := repository.Create(character)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		character.ID = id
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(character)
	}
}

// ReadCharacter retrieves all characters for while
func ReadCharacter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repository := repositories.NewCharacterRepository(db)
		result, err := repository.Read()
		if err != nil {
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}

// UpdateCharacter changes the users by ID
func UpdateCharacter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		var character models.Character
		defer r.Body.Close()

		id, _ := strconv.ParseInt(vars["id"], 10, 64)

		if err := json.NewDecoder(r.Body).Decode(&character); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := character.IsValid(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		character.ID = id

		repo := repositories.NewCharacterRepository(db)
		_, err := repo.Update(character)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

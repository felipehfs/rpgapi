package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/felipehfs/rpgapi/models"
	"github.com/felipehfs/rpgapi/repositories"
	"github.com/gorilla/mux"
)

var (
	ErrQueryNotFound = errors.New("Query not found in endpoint")
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

func hasQueryParam(query string, r *http.Request) bool {
	return r.URL.Query().Get(query) != ""
}

func parseQueryInt64(query string, r *http.Request, fallback int64) (int64, error) {
	if hasQueryParam(query, r) {
		query, err := strconv.ParseInt(r.URL.Query().Get(query), 10, 64)
		if err != nil {
			return -1, err
		}
		return query, nil
	}

	return fallback, ErrQueryNotFound
}

// ReadCharacter retrieves all characters for while
func ReadCharacter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repository := repositories.NewCharacterRepository(db)
		total, err := repository.Count()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var limit int64 = 50
		var page int64 = 1

		userLimit, err := parseQueryInt64("limit", r, limit)

		if err != nil && err != ErrQueryNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		limit = userLimit

		userPage, err := parseQueryInt64("page", r, page)

		if err != nil && err != ErrQueryNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		page = userPage

		offset := (page - 1) * limit

		characters, err := repository.Read(limit, offset)
		if err != nil {
			http.Error(w, err.Error(), http.StatusFailedDependency)
			return
		}

		result := map[string]interface{}{
			"data":        characters,
			"page":        page,
			"limit":       limit,
			"total_pages": total / limit,
			"total":       total,
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

		id, err := strconv.ParseInt(vars["id"], 10, 64)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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
		_, err = repo.Update(character)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

// RemoveCharacter removes the character by ID
func RemoveCharacter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		response := make(map[string]interface{})
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		repo := repositories.NewCharacterRepository(db)

		affected, err := repo.Remove(id)

		if err != nil {
			http.Error(w, err.Error(), getCode(err))
			return
		}

		response["rows_affected"] = affected
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(response)
	}
}

func getCode(err error) int {
	if err == sql.ErrNoRows {
		return http.StatusNotFound
	} else {
		return http.StatusInternalServerError
	}
}

// GetByIDCharacter retrieves the character by ID
func GetByIDCharacter(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseInt(vars["id"], 10, 64)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		repo := repositories.NewCharacterRepository(db)
		character, err := repo.GetByID(id)

		if err != nil {
			http.Error(w, err.Error(), getCode(err))
			return
		}

		json.NewEncoder(w).Encode(character)
	}
}

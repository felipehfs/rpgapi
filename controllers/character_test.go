package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/felipehfs/rpgapi/config"
	"github.com/felipehfs/rpgapi/models"
	"github.com/gorilla/mux"

	"github.com/felipehfs/rpgapi/controllers"
)

type MockDB struct {
	DB *sql.DB
}

var mockDB *MockDB

func TestMain(m *testing.M) {
	db, err := config.SetupDatabase(config.Test)
	if err != nil {
		log.Fatal(err)
	}
	mockDB = &MockDB{
		DB: db,
	}
	defer db.Close()
	os.Exit(m.Run())
}

var body = []byte(`
	{
		"name": "Zeo (created)",
		"attack": 1230, 
		"defense": 3200, 
		"speed": 30,
		"life": 120
	}
`)

func createCharacter() *models.Character {
	var character models.Character
	req := httptest.NewRequest("POST", "/api/characters", bytes.NewBuffer(body))
	res := httptest.NewRecorder()

	createCaracter := controllers.CreateCharacter(mockDB.DB)
	createCaracter(res, req)
	result := res.Result()
	defer result.Body.Close()
	json.NewDecoder(result.Body).Decode(&character)
	return &character
}

func TestCreateCharacterHandler(t *testing.T) {
	testcases := []struct {
		Name           string
		Body           []byte
		ExpectedStatus int
	}{
		{"Expected Status created", []byte(`{"name": "Zeo", "attack": 1230, "defense": 3200, "speed": 30, "life": 120}`), 201},
		{"Request empty", []byte(``), http.StatusBadRequest},
		{"Required fields", []byte(`{"name": "Lucy" }`), http.StatusBadRequest},
	}
	for _, tt := range testcases {
		t.Run(tt.Name, func(t *testing.T) {

			req := httptest.NewRequest("POST", "/api/characters", bytes.NewBuffer(tt.Body))

			res := httptest.NewRecorder()

			createCaracter := controllers.CreateCharacter(mockDB.DB)
			createCaracter(res, req)
			result := res.Result()

			if result.StatusCode != tt.ExpectedStatus {
				body, _ := ioutil.ReadAll(result.Body)
				defer result.Body.Close()
				fmt.Println(string(body))
				t.Errorf("%s: Expected status %d, got %d", tt.Name, result.StatusCode, tt.ExpectedStatus)
			}
		})
	}
}

func TestReadCharacterHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/characters", nil)
	res := httptest.NewRecorder()
	readCharactersHandler := controllers.ReadCharacter(mockDB.DB)
	readCharactersHandler(res, req)
	result := res.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, result.StatusCode)
	}

}

func TestUpdateCharacterHandler(t *testing.T) {
	char := createCharacter()

	testcases := []struct {
		Name         string
		Body         []byte
		ID           int64
		ExpectedCode int
	}{
		{"Expected 200 OK", body, char.ID, http.StatusOK},
		{"Required Fields: 400 Status", nil, char.ID, http.StatusBadRequest},
	}
	updateHandler := controllers.UpdateCharacter(mockDB.DB)

	for _, test := range testcases {
		t.Run(test.Name, func(t *testing.T) {

			endpoint := "/api/characters/"
			req := httptest.NewRequest("PUT", endpoint, bytes.NewBuffer(body))
			res := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{
				"id": fmt.Sprintf("%v", test.ID),
			})

			updateHandler(res, req)
			result := res.Result()

			if result.StatusCode != http.StatusOK {
				t.Errorf("Expected status code %d, but got %d", http.StatusOK, result.StatusCode)
			}

		})
	}
}

func TestRemoveCharacterHandler(t *testing.T) {
	example := createCharacter()
	testcases := []struct {
		Name         string
		ID           int64
		ExpectedCode int
	}{
		{"Expected 201 status Code", example.ID, http.StatusCreated},
	}

	handler := controllers.RemoveCharacter(mockDB.DB)

	for _, test := range testcases {
		t.Run(test.Name, func(t *testing.T) {
			endpoint := "/api/characters/"
			req := httptest.NewRequest("DELETE", endpoint, nil)

			req = mux.SetURLVars(req, map[string]string{
				"id": fmt.Sprintf("%v", test.ID),
			})

			res := httptest.NewRecorder()

			handler(res, req)

			result := res.Result()

			if result.StatusCode != test.ExpectedCode {
				body, _ := ioutil.ReadAll(result.Body)
				defer result.Body.Close()
				fmt.Println(string(body))
				t.Errorf("Expected status code %d, but got %d", test.ExpectedCode, result.StatusCode)
			}
		})
	}
}

func TestGetByIDCharacterHandler(t *testing.T) {
	example := createCharacter()
	testCases := []struct {
		Name         string
		ID           int64
		ExpectedCode int
	}{
		{"Expected Status 200 - OK", example.ID, http.StatusOK},
		{"Expected Status 404 - Not Found", 999999, http.StatusNotFound},
	}

	handler := controllers.GetByIDCharacter(mockDB.DB)

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {

			req := httptest.NewRequest("GET", "/api/characters", nil)
			req = mux.SetURLVars(req, map[string]string{
				"id": fmt.Sprintf("%v", testCase.ID),
			})

			res := httptest.NewRecorder()

			handler(res, req)

			result := res.Result()

			if result.StatusCode != testCase.ExpectedCode {
				t.Errorf("Expected status code %v but got %v", testCase.ExpectedCode, result.StatusCode)
			}
		})
	}
}

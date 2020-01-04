package controllers_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/felipehfs/rpgapi/config"

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

			req := httptest.NewRequest("POST", "/characters", bytes.NewBuffer(tt.Body))

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
	req := httptest.NewRequest("GET", "/products", nil)
	res := httptest.NewRecorder()
	readCharactersHandler := controllers.ReadCharacter(mockDB.DB)
	readCharactersHandler(res, req)
	result := res.Result()

	if result.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, result.StatusCode)
	}

}

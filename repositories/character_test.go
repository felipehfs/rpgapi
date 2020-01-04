package repositories_test

import (
	"log"
	"os"
	"testing"

	"github.com/felipehfs/rpgapi/config"
	"github.com/felipehfs/rpgapi/models"
	"github.com/felipehfs/rpgapi/repositories"
)

type MockDB struct {
	Repo *repositories.CharacterRepository
}

var mockDB *MockDB

var data = models.Character{
	ID:      0,
	Name:    "Test",
	Attack:  10,
	Defense: 10,
	Speed:   10,
	Life:    20,
}

func TestMain(m *testing.M) {
	db, err := config.SetupDatabase(config.Test)

	if err != nil {
		log.Fatal(err)
	}
	repo := repositories.NewCharacterRepository(db)

	mockDB = &MockDB{
		Repo: repo,
	}

	defer db.Close()
	result := m.Run()
	os.Exit(result)
}
func TestRepositoryCreateCharacter(t *testing.T) {

	_, err := mockDB.Repo.Create(data)
	if err != nil {
		t.Error(err)
	}
}
func TestRepositoryUpdateCharacter(t *testing.T) {
	id, err := mockDB.Repo.Create(data)
	data.ID = id
	data.Name = "Test (updated)"
	var expectedChanges int64 = 1
	rows, err := mockDB.Repo.Update(data)

	if err != nil {
		t.Error(err)
	}
	if rows != expectedChanges {
		t.Errorf("Expected %d rows affected, got %d", expectedChanges, rows)
	}
}

func TestRepositoryReadCharacter(t *testing.T) {
	_, err := mockDB.Repo.Read()
	if err != nil {
		t.Error(err)
	}
}

func TestRepositoryDeleteCharacter(t *testing.T) {

	id, err := mockDB.Repo.Create(data)
	if err != nil {
		t.Error(err)
	}

	rowsAffected, err := mockDB.Repo.Remove(id)
	if err != nil {
		t.Error(err)
	}

	if rowsAffected != 1 {
		t.Errorf("Expected remove one row but got %d", rowsAffected)
	}
}

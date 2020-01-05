package repositories

import (
	"database/sql"

	"github.com/felipehfs/rpgapi/models"
)

// CharacterRepository makes the commom database operations
type CharacterRepository struct {
	DB *sql.DB
}

// NewCharacterRepository instantiate the repository
func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{DB: db}
}

// Create saves the new character
func (repo CharacterRepository) Create(character models.Character) (int64, error) {
	sql := `
		INSERT INTO characters(name, attack, defense, speed, life) 
			VALUES($1, $2, $3, $4, $5) RETURNING id
	`
	var id int64
	err := repo.DB.QueryRow(sql, character.Name,
		character.Attack, character.Defense,
		character.Speed, character.Life).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// GetByID retrieves one character if exists
func (repo CharacterRepository) GetByID(id int64) (models.Character, error) {
	var char models.Character
	query := "SELECT * FROM characters WHERE id=$1"
	err := repo.DB.QueryRow(query, id).Scan(&char.ID,
		&char.Name, &char.Attack,
		&char.Defense, &char.Speed,
		&char.Life)
	return char, err
}

func (repo CharacterRepository) Read() ([]models.Character, error) {
	var characters []models.Character
	sql := "SELECT * FROM characters"
	rows, err := repo.DB.Query(sql)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c models.Character
		rows.Scan(&c.ID, &c.Name, &c.Attack, &c.Defense, &c.Speed, &c.Life)
		characters = append(characters, c)
	}

	return characters, nil
}

// Update changes the character by ID
func (repo CharacterRepository) Update(character models.Character) (int64, error) {
	query := `UPDATE characters 
		SET name=$2, attack=$3, defense=$4, speed=$5, life=$6 
		WHERE id=$1`
	result, err := repo.DB.Exec(query, character.ID, character.Name,
		character.Attack, character.Defense, character.Speed, character.Life)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

// Remove makes what the method describe by id
func (repo CharacterRepository) Remove(id int64) (int64, error) {
	query := "DELETE FROM characters WHERE id=$1"

	result, err := repo.DB.Exec(query, id)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

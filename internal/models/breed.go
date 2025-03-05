// internal/models/breed.go
package models

import (
	"database/sql"
	"errors"
	"strconv"
)

type BreedStore struct {
	db *sql.DB
}

type Breed struct {
	ID                       int    `json:"id"`
	Species                  string `json:"species"`
	PetSize                  string `json:"pet_size"`
	Name                     string `json:"name"`
	AverageMaleAdultWeight   int    `json:"average_male_adult_weight"`
	AverageFemaleAdultWeight int    `json:"average_female_adult_weight"`
}

func NewBreedStore(db *sql.DB) *BreedStore {
	return &BreedStore{db: db}
}

// GetAllBreeds retrieves all breeds from the MySQL database
func (bs *BreedStore) GetAllBreeds() ([]Breed, error) {
	rows, err := bs.db.Query("SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight FROM breeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allBreeds []Breed
	for rows.Next() {
		var breed Breed
		if err := rows.Scan(&breed.ID, &breed.Species, &breed.PetSize, &breed.Name, &breed.AverageMaleAdultWeight, &breed.AverageFemaleAdultWeight); err != nil {
			return nil, err
		}
		allBreeds = append(allBreeds, breed)
	}
	return allBreeds, nil
}

// CreateBreed adds a new breed to the MySQL database
func (bs *BreedStore) CreateBreed(breed *Breed) error {
	result, err := bs.db.Exec("INSERT INTO breeds (species, pet_size, name, average_male_adult_weight, average_female_adult_weight) VALUES (?, ?, ?, ?, ?)",
		breed.Species, breed.PetSize, breed.Name, breed.AverageMaleAdultWeight, breed.AverageFemaleAdultWeight)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	breed.ID = int(id)
	return nil
}

// GetBreedByID retrieves a breed by its ID from the MySQL database
func (bs *BreedStore) GetBreedByID(id string) (*Breed, error) {
	breedID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid breed ID")
	}

	var breed Breed
	err = bs.db.QueryRow("SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight FROM breeds WHERE id = ?", breedID).
		Scan(&breed.ID, &breed.Species, &breed.PetSize, &breed.Name, &breed.AverageMaleAdultWeight, &breed.AverageFemaleAdultWeight)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("breed not found")
		}
		return nil, err
	}
	return &breed, nil
}

// UpdateBreed updates an existing breed in the MySQL database
func (bs *BreedStore) UpdateBreed(id string, breed *Breed) error {
	breedID, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid breed ID")
	}

	_, err = bs.db.Exec("UPDATE breeds SET species = ?, pet_size = ?, name = ?, average_male_adult_weight = ?, average_female_adult_weight = ? WHERE id = ?",
		breed.Species, breed.PetSize, breed.Name, breed.AverageMaleAdultWeight, breed.AverageFemaleAdultWeight, breedID)
	if err != nil {
		return err
	}
	breed.ID = breedID
	return nil
}

// DeleteBreed removes a breed from the MySQL database
func (bs *BreedStore) DeleteBreed(id string) error {
	breedID, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid breed ID")
	}

	_, err = bs.db.Exec("DELETE FROM breeds WHERE id = ?", breedID)
	if err != nil {
		return err
	}
	return nil
}

// SearchBreeds retrieves breeds from the MySQL database based on search criteria
func (bs *BreedStore) SearchBreeds(species string, minWeight, maxWeight int) ([]Breed, error) {
	query := "SELECT id, species, pet_size, name, average_male_adult_weight, average_female_adult_weight FROM breeds WHERE 1=1"
	args := []interface{}{}

	if species != "" {
		query += " AND species = ?"
		args = append(args, species)
	}
	if minWeight > 0 {
		query += " AND (average_male_adult_weight >= ? OR average_female_adult_weight >= ?)"
		args = append(args, minWeight, minWeight)
	}
	if maxWeight > 0 {
		query += " AND (average_male_adult_weight <= ? OR average_female_adult_weight <= ?)"
		args = append(args, maxWeight, maxWeight)
	}

	rows, err := bs.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var breeds []Breed
	for rows.Next() {
		var breed Breed
		if err := rows.Scan(&breed.ID, &breed.Species, &breed.PetSize, &breed.Name, &breed.AverageMaleAdultWeight, &breed.AverageFemaleAdultWeight); err != nil {
			return nil, err
		}
		breeds = append(breeds, breed)
	}

	// Return an empty slice instead of nil if no results found
	if len(breeds) == 0 {
		return []Breed{}, nil
	}

	return breeds, nil
}

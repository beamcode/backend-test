// internal/handlers/breeds.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/models"
)

var store *models.BreedStore

// InitHandlers initializes the breed store with the given database connection
func InitHandlers(db *sql.DB) {
	store = models.NewBreedStore(db)
}

// ListBreeds handles the GET request to list all breeds
func ListBreeds(w http.ResponseWriter, r *http.Request) {
	breeds, err := store.GetAllBreeds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breeds)
}

// CreateBreed handles the POST request to create a new breed
func CreateBreed(w http.ResponseWriter, r *http.Request) {
	var breed models.Breed
	if err := json.NewDecoder(r.Body).Decode(&breed); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := store.CreateBreed(&breed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(breed)
}

// GetBreed handles the GET request to retrieve a breed by ID
func GetBreed(w http.ResponseWriter, r *http.Request) {
	breed, err := store.GetBreedByID(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Breed not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breed)
}

// UpdateBreed handles the PUT request to update an existing breed
func UpdateBreed(w http.ResponseWriter, r *http.Request) {
	var breed models.Breed
	if err := json.NewDecoder(r.Body).Decode(&breed); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := store.UpdateBreed(mux.Vars(r)["id"], &breed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(breed)
}

// DeleteBreed handles the DELETE request to delete a breed
func DeleteBreed(w http.ResponseWriter, r *http.Request) {
	if err := store.DeleteBreed(mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// SearchBreeds handles the GET request to search breeds based on filters
func SearchBreeds(w http.ResponseWriter, r *http.Request) {
	species := r.URL.Query().Get("species")
	minWeightStr := r.URL.Query().Get("minWeight")
	maxWeightStr := r.URL.Query().Get("maxWeight")

	var minWeight, maxWeight int
	var err error

	if minWeightStr != "" {
		if minWeight, err = strconv.Atoi(minWeightStr); err != nil {
			http.Error(w, "Invalid minWeight", http.StatusBadRequest)
			return
		}
	}

	if maxWeightStr != "" {
		if maxWeight, err = strconv.Atoi(maxWeightStr); err != nil {
			http.Error(w, "Invalid maxWeight", http.StatusBadRequest)
			return
		}
	}

	// Ensure at least one parameter is provided
	if species == "" && minWeightStr == "" && maxWeightStr == "" {
		http.Error(w, "At least one parameter (species, minWeight, or maxWeight) is required", http.StatusBadRequest)
		return
	}

	// Call store function
	breeds, err := store.SearchBreeds(species, minWeight, maxWeight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Return message if no breeds found
	if len(breeds) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "No breeds found"})
		return
	}

	json.NewEncoder(w).Encode(breeds)
}

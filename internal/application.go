// internal/application.go
package internal

import (
	"net/http"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/handlers"
)

type App struct {
	logger *charmLog.Logger
}

func NewApp(logger *charmLog.Logger) *App {
	return &App{
		logger: logger,
	}
}

func (a *App) RegisterRoutes(r *mux.Router) {
	breedsRouter := r.PathPrefix("/breeds").Subrouter()
	breedsRouter.HandleFunc("", a.ListBreeds).Methods(http.MethodGet)
	breedsRouter.HandleFunc("", a.CreateBreed).Methods(http.MethodPost)
	// breedsRouter.HandleFunc("/{id}", a.GetBreed).Methods(http.MethodGet)
	// breedsRouter.HandleFunc("/{id}", a.UpdateBreed).Methods(http.MethodPut)
	// breedsRouter.HandleFunc("/{id}", a.DeleteBreed).Methods(http.MethodDelete)
	breedsRouter.HandleFunc("/search", a.SearchBreeds).Methods(http.MethodGet)
}

func (a *App) ListBreeds(w http.ResponseWriter, r *http.Request) {
	handlers.ListBreeds(w, r)
}

func (a *App) CreateBreed(w http.ResponseWriter, r *http.Request) {
	handlers.CreateBreed(w, r)
}

func (a *App) GetBreed(w http.ResponseWriter, r *http.Request) {
	handlers.GetBreed(w, r)
}

func (a *App) UpdateBreed(w http.ResponseWriter, r *http.Request) {
	handlers.UpdateBreed(w, r)
}

func (a *App) DeleteBreed(w http.ResponseWriter, r *http.Request) {
	handlers.DeleteBreed(w, r)
}

func (a *App) SearchBreeds(w http.ResponseWriter, r *http.Request) {
	handlers.SearchBreeds(w, r)
}

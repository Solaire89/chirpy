package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Solaire89/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (a *apiConfig) handlerPostChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	param := params{}
	err := decoder.Decode(&param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding request")
		return
	}

	if len(param.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	param.Body = censorBadWord(param.Body)
	chirp := database.CreateChirpParams{
		Body:   param.Body,
		UserID: param.UserID,
	}

	c, err := a.db.CreateChirp(r.Context(), chirp)
	if err != nil {
		log.Printf("Error: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		UserID:    c.UserID,
		Body:      c.Body,
	})
}

func (a *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := a.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error getting chirp")
		return
	}
	foundChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
		Body:      chirp.Body,
	}
	respondWithJSON(w, http.StatusOK, foundChirp)
}

func (a *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := a.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("error getting chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	chirpResponse := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		chirpResponse[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserID:    chirp.UserID,
			Body:      chirp.Body,
		}

	}
	respondWithJSON(w, http.StatusOK, chirpResponse)
}

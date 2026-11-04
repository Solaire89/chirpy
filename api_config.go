package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/Solaire89/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) handlerPrintRequests(w http.ResponseWriter, r *http.Request) {
	hits := a.fileserverHits.Load()
	fmt.Fprintf(w, adminHTML, hits)
}

func (a *apiConfig) handlerResetCount(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading environment: %s", err)
	}
	PLATFORM := os.Getenv("PLATFORM")
	if PLATFORM != "dev" {
		respondWithError(w, 403, "forbidden")
	}
	a.db.ClearTable(r.Context())
	a.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Resets counter to 0"))
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}
	returnErrorIns := returnError{
		Error: msg,
	}
	dat, err := json.Marshal(returnErrorIns)
	if err != nil {
		log.Printf("error marshalling data: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(dat)
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
		respondWithError(w, 400, "Bad response")
		return
	}

	if len(param.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	param.Body = censorBadWord(param.Body)

	type returnCleanedChirp struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	respondWithJSON(w, 201, returnCleanedChirp{
		Body:   param.Body,
		UserID: param.UserID,
	})
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type emailParams struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	param := emailParams{}
	if err := decoder.Decode(&param); err != nil {
		respondWithError(w, 400, "Bad response")
		return
	}
	user, err := a.db.CreateUser(r.Context(), param.Email)
	userParam := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, 500, "internal server error")
		return
	}

	respondWithJSON(w, 201, userParam)
}

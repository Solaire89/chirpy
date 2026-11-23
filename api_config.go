package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Solaire89/chirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	Secret         string
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
		return
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
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

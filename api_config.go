package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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

func (a *apiConfig) handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}

	type returnError struct {
		Error string `json:"error"`
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

	type returnValid struct {
		Valid bool `json:"valid"`
	}

	respondWithJSON(w, 200, returnValid{
		Valid: true,
	})
}

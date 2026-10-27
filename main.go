package main

import (
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
	fmt.Fprintf(w, "Hits: %d\n", hits)
}

func (a *apiConfig) handlerResetCount(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Resets counter to 0"))
}

func main() {
	const port = "8080"

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/metrics", apiCfg.handlerPrintRequests)
	mux.HandleFunc("/reset", apiCfg.handlerResetCount)
	mux.HandleFunc("/healthz", handlerServerCode)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}

func handlerServerCode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		log.Printf("Error writing body: %s", err)
	}
}

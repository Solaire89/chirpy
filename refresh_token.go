package main

import (
	"net/http"
	"time"

	"github.com/Solaire89/chirpy/internal/auth"
)

func (a *apiConfig) handlerGetRefreshToken(w http.ResponseWriter, r *http.Request) {
	type JWT struct {
		JWT string `json:"token"`
	}

	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting refresh token")
		return
	}
	user, err := a.db.GetUserFromRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not found")
		return
	}
	tokenJWT, err := auth.MakeJWT(user.ID, a.Secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error making token")
		return
	}

	respondWithJSON(w, http.StatusOK, JWT{JWT: tokenJWT})
}

func (a *apiConfig) handlerAPIRevoke(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting bearer token")
		return
	}

	_, err = a.db.RevokeRefreshToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

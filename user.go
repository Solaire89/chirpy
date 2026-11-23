package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Solaire89/chirpy/internal/auth"
	"github.com/Solaire89/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type userParams struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	param := userParams{}
	if err := decoder.Decode(&param); err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}
	hashPassword, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password")
		return
	}
	user, err := a.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          param.Email,
		HashedPassword: hashPassword,
	})

	const defaultSeconds = 60
	const maxSeconds = 3600
	seconds := defaultSeconds
	if param.ExpiresInSeconds > 0 && param.ExpiresInSeconds <= maxSeconds {
		seconds = param.ExpiresInSeconds
	} else {
		seconds = maxSeconds
	}

	timeOut := time.Duration(seconds) * time.Second

	token, err := auth.MakeJWT(user.ID, a.Secret, timeOut)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}

	userParam := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}

	if err != nil {
		log.Printf("error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusCreated, userParam)
}

func (a *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type loginParams struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	loginParam := loginParams{}
	if err := decoder.Decode(&loginParam); err != nil {
		respondWithError(w, http.StatusBadRequest, "bad request")
	}
	user, err := a.db.LoginEmail(r.Context(), loginParam.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting user")
	}
	authPassword, err := auth.CheckPasswordHash(loginParam.Password, user.HashedPassword)
	if err != nil || authPassword == false {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	const defaultSeconds = 60
	const maxSeconds = 3600
	seconds := defaultSeconds
	if loginParam.ExpiresInSeconds > 0 && loginParam.ExpiresInSeconds <= maxSeconds {
		seconds = loginParam.ExpiresInSeconds
	} else {
		seconds = maxSeconds
	}

	timeOut := time.Duration(seconds) * time.Second

	token, err := auth.MakeJWT(user.ID, a.Secret, timeOut)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}

	loginUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	}

	respondWithJSON(w, http.StatusOK, loginUser)
}

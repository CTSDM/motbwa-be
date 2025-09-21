package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/CTSDM/motbwa-be/internal/auth"
	"github.com/CTSDM/motbwa-be/internal/database"
	"github.com/google/uuid"
)

func (cfg *CfgAPI) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type responseVals struct {
		ID        uuid.UUID `json:"user_id"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if params.Username == "" || params.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong handling the password", err)
		return
	}

	// add the user to the database
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Username:       params.Username,
		HashedPassword: hashed_password,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong while creating the user", err)
		return
	}

	RespondWithJSON(w, http.StatusCreated, responseVals{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (cfg *CfgAPI) HandlerDeleteAllUsers(w http.ResponseWriter, r *http.Request) {
	err := cfg.DB.DeleteUsers(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong while deleting from the users table", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (cfg *CfgAPI) HandlerCheckUserExists(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		RespondWithError(w, http.StatusBadRequest, "user parameter required", nil)
		return
	}

	if _, err := cfg.DB.GetUserByUsername(r.Context(), username); err == sql.ErrNoRows {
		RespondWithError(w, http.StatusNotFound, "user not found", nil)
		return
	} else if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong while retrieving the user from the database", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

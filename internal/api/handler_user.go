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
	if err := decoder.Decode(&params); err != nil {
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
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/CTSDM/motbwa-be/internal/api"
	"github.com/CTSDM/motbwa-be/internal/database"
	"github.com/CTSDM/motbwa-be/internal/ws"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		return
	}

	portNumber := "8080"

	serveMux, err := setupAPI()
	if err != nil {
		log.Fatal(err)
		return
	}

	server := &http.Server{
		Addr:    ":" + portNumber,
		Handler: serveMux,
		// timeout to avoid slowloris attack
		ReadTimeout: 5 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

func setupAPI() (*http.ServeMux, error) {
	ctx := context.Background()
	manager := ws.NewManager(ctx)

	// let's load the DB
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	// let's load jwt env variables
	tokenSecret := os.Getenv("TOKEN_SECRET")
	tokenExpiration, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION_SECONDS"))
	if err != nil {
		return nil, err
	} else if tokenExpiration == 0 {
		tokenExpiration = 120
	}
	refreshTokenExpiration, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRATION_SECONDS"))
	if err != nil {
		return nil, err
	} else if refreshTokenExpiration == 0 {
		refreshTokenExpiration = 60 * 60 * 24
	}

	cfgApi := api.CfgAPI{
		DB:                     database.New(db),
		TokenSecret:            tokenSecret,
		TokenExpiration:        time.Second * time.Duration(tokenExpiration),
		RefreshTokenExpiration: time.Second * time.Duration(refreshTokenExpiration),
	}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/ws", cfgApi.HandlerMiddlewareLogin(manager.ServeWS))
	serveMux.HandleFunc("POST /api/users", cfgApi.HandlerCreateUser)
	serveMux.HandleFunc("POST /api/login", cfgApi.HandlerLogin)
	serveMux.HandleFunc("POST /admin/reset", cfgApi.HandlerDeleteAllUsers)

	return serveMux, nil
}

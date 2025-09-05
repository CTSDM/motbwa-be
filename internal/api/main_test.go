package api

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var db *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	connURL, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	db, err = sql.Open("postgres", connURL)
	if err != nil {
		log.Printf("failed to connect to database: %s", err)
		return
	}

	if err := db.Ping(); err != nil {
		log.Printf("Couldn't ping the db...: %s", err)
		return
	}

	if err := goose.UpContext(ctx, db, "../../sql/schema"); err != nil {
		log.Printf("Something went wrong while migrating the database...: %s", err)
		return
	}

	res := m.Run()

	os.Exit(res)
}

func cleanup() {
	db.Exec("DELETE FROM users;")
}

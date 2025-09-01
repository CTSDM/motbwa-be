package api

import (
	"time"

	"github.com/CTSDM/motbwa-be/internal/database"
)

type CfgAPI struct {
	DB                     *database.Queries
	TokenSecret            string
	TokenExpiration        time.Duration
	RefreshTokenExpiration time.Duration
}

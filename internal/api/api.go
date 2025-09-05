package api

import (
	"time"

	"github.com/CTSDM/motbwa-be/internal/database"
)

// type Querier interface {
// 	CreateMessage(ctx context.Context, arg database.CreateMessageParams) (database.Message, error)
// 	DeleteMessages(ctx context.Context) error
// 	GetMessagesByReceiver(ctx context.Context, receiverID uuid.UUID) ([]database.Message, error)
// 	GetMessagesBySender(ctx context.Context, senderID uuid.UUID) ([]database.Message, error)
// 	CreateRefreshToken(ctx context.Context, arg database.CreateRefreshTokenParams) (database.RefreshToken, error)
// 	DeleteRefreshToken(ctx context.Context, userID uuid.UUID) error
// 	GetRefreshToken(ctx context.Context, userID uuid.UUID) (database.RefreshToken, error)
// 	UpdateRefreshToken(ctx context.Context, arg database.UpdateRefreshTokenParams) (database.RefreshToken, error)
// 	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)
// 	DeleteUsers(ctx context.Context) error
// 	GetUser(ctx context.Context, id uuid.UUID) (database.User, error)
// 	GetUserByUsername(ctx context.Context, username string) (database.User, error)
// }

// var _ Querier = (*database.Queries)(nil)

type CfgAPI struct {
	DB                     *database.Queries
	TokenSecret            string
	TokenExpiration        time.Duration
	RefreshTokenExpiration time.Duration
}

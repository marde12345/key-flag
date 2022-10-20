package user

import (
	"context"
	"database/sql"

	// entity dependency
	userentity "github.com/tokopedia/kv-middleware/internal/entity/user"
)

//go:generate mockgen -source=repository.go -package=user -destination=repository_mock_test.go

type userRepository interface {
	GetDBTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	GetUser(ctx context.Context, username string) (userentity.User, error)
	GetUserAccess(ctx context.Context, userID int) ([]userentity.Role, error)
	VerifyUser(ctx context.Context, email, token string) (bool, error)
	MapUserAccess(ctx context.Context, tx *sql.Tx, userID int, roles []userentity.Role) error
	DeleteUserAccess(ctx context.Context, email string) error
	CreateUser(ctx context.Context, username, email, token string) error
	CreateRole(ctx context.Context, tx *sql.Tx, prefix, permission string, userID int) (int, error)
	GetAllRoles(ctx context.Context) ([]userentity.Role, error)
	GetRole(ctx context.Context, prefix, permission string) (userentity.Role, error)
	RevokeUserAccess(ctx context.Context, userID, roleID, requestedBy int) error
	SearchRole(ctx context.Context, prefix string) ([]userentity.Role, error)
}

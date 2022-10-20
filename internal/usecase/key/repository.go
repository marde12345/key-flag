package key

import (
	"context"
	"database/sql"

	// entity dependency
	keyentity "github.com/marde12345/key-flag/internal/entity/key"
	userentity "github.com/marde12345/key-flag/internal/entity/user"
)

//go:generate mockgen -source=repository.go -package=key -destination=repository_mock_test.go

type keyRepository interface {
	GetDBTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	GetKey(ctx context.Context, key string, status int) ([]keyentity.KV, error)
	GetKeyByID(ctx context.Context, keyID int) (keyentity.KV, error)
	GetKeyByPrefix(ctx context.Context, prefix string, status int) ([]keyentity.KV, error)
	GetKeyHistory(ctx context.Context, key string, isPrefix bool, limit int) ([]keyentity.KV, error)
	GetKeyListWithoutValue(prefix string) ([]string, error)
	CreateKeyEntry(ctx context.Context, tx *sql.Tx, kv keyentity.KV) error
	CreateKey(ctx context.Context, tx *sql.Tx, key string, value string, valType string, userID, status int) error
	ModifyKey(ctx context.Context, tx *sql.Tx, keyID int, kv keyentity.KV) error
	SetCache(ctx context.Context, key keyentity.KV) error
	GetCache(ctx context.Context, key string) (keyentity.KV, error)
	GetCaches(ctx context.Context, key string) ([]keyentity.KV, error)
	InvalidateCache(ctx context.Context, key string) error
	ModifyOldActiveKey(ctx context.Context, tx *sql.Tx, key string) error
	IsKeyExist(ctx context.Context, key string) bool
	RegisterCanaryDeployment(ctx context.Context, service string, nodesIP []string) error
	ReleaseCanaryIP(ctx context.Context, service string) error
	GetCanaryIP(ctx context.Context, service string) []string
	GetCanaryKV(ctx context.Context, ip string) ([]keyentity.KV, error)
	CreateCanaryKey(ctx context.Context, tx *sql.Tx, id int, ip string) error
	ModifyCanaryKey(ctx context.Context, tx *sql.Tx, id, status int) error
	GetCanaryKVByID(ctx context.Context, id int) ([]keyentity.CanaryKV, error)
}

type userRepository interface {
	GetDBTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	VerifyUser(ctx context.Context, email, token string) (bool, error)
	CreateRole(ctx context.Context, tx *sql.Tx, prefix, permission string, userID int) (int, error)
	MapUserAccess(ctx context.Context, tx *sql.Tx, userID int, roles []userentity.Role) error
	GetUser(ctx context.Context, username string) (userentity.User, error)
}

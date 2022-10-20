package user

import (
	"context"
	"database/sql"
	"fmt"

	// internal dependency
	userentity "github.com/tokopedia/kv-middleware/internal/entity/user"
)

type Usecase struct {
	userRepo userRepository
}

func New(user userRepository) *Usecase {
	return &Usecase{
		userRepo: user,
	}
}

func (u *Usecase) CreateUser(user userentity.User) error {
	ctx := context.Background()

	// check user is exist or not first
	// prevent double row
	userRecord, _ := u.userRepo.GetUser(ctx, user.Username)
	if userRecord.ID > 0 {
		// just return if user already created
		return nil
	}

	// create if not exist
	return u.userRepo.CreateUser(ctx, user.Username, user.Email, user.Token)
}

func (u *Usecase) GetUserDetails(username string) (userentity.UserDetails, error) {
	ctx := context.Background()

	user, err := u.userRepo.GetUser(ctx, username)
	if err != nil && err != sql.ErrNoRows {
		return userentity.UserDetails{}, err
	}

	// normalize error if sql no rows
	err = nil

	if user.ID == 0 {
		return userentity.UserDetails{}, fmt.Errorf("User %s is not found.", username)
	}

	roles, err := u.userRepo.GetUserAccess(ctx, user.ID)
	if err != nil {
		return userentity.UserDetails{}, err
	}

	return userentity.UserDetails{
		User:  user,
		Roles: roles,
	}, nil
}

func (u *Usecase) CreateRole(roles []userentity.Role, userID int) error {
	ctx := context.Background()

	tx, err := u.userRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, role := range roles {
		_, err := u.userRepo.CreateRole(ctx, tx, role.Prefix, role.Permission, userID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (u *Usecase) MapUserAccess(userID int, roles []userentity.Role) error {
	ctx := context.Background()

	tx, err := u.userRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = u.userRepo.MapUserAccess(ctx, tx, userID, roles)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Usecase) GetAllRoles() ([]userentity.Role, error) {
	ctx := context.Background()

	return u.userRepo.GetAllRoles(ctx)
}

func (u *Usecase) GetRole(prefix, permission string) (userentity.Role, error) {
	ctx := context.Background()

	return u.userRepo.GetRole(ctx, prefix, permission)
}

func (u *Usecase) RevokeUserAccess(userID, requestedBy int, roles []userentity.Role) error {
	ctx := context.Background()

	for _, role := range roles {
		if err := u.userRepo.RevokeUserAccess(ctx, userID, role.ID, requestedBy); err != nil {
			return err
		}
	}

	return nil
}

func (u *Usecase) SearchRole(prefix string) ([]userentity.Role, error) {
	ctx := context.Background()

	return u.userRepo.SearchRole(ctx, prefix)
}

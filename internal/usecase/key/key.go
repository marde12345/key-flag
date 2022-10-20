package key

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	keyentity "github.com/marde12345/key-flag/internal/entity/key"
	userentity "github.com/marde12345/key-flag/internal/entity/user"
)

type Usecase struct {
	keyRepo  keyRepository
	userRepo userRepository
}

func New(key keyRepository, user userRepository) *Usecase {
	return &Usecase{
		keyRepo:  key,
		userRepo: user,
	}
}

func (u *Usecase) UpdateKey(kv keyentity.KV) error {
	ctx := context.Background()

	//get requested update key first
	keyWaitingApprovalUpdate, err := u.keyRepo.GetKey(ctx, kv.Key, keyentity.PlacedKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if len(keyWaitingApprovalUpdate) >= 1 {
		return errors.New("Detected more than one pending approval key")
	}

	//get requested delete key
	keyWaitingApprovalDelete, err := u.keyRepo.GetKey(ctx, kv.Key, keyentity.PlacedDeleteKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if len(keyWaitingApprovalDelete) >= 1 {
		return errors.New("Detected more than one pending approval key")
	}

	keyInCanary, err := u.keyRepo.GetKey(ctx, kv.Key, keyentity.CanaryKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if len(keyInCanary) >= 1 {
		return errors.New("Can not change value in canary.")
	}

	tx, err := u.keyRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// all keys that newly updated will have placed status
	err = u.keyRepo.CreateKey(ctx, tx, kv.Key, kv.Value, kv.Type, kv.CreatedBy, keyentity.PlacedKey)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Usecase) CreateDeleteKey(kv keyentity.KV) error {
	ctx := context.Background()

	//get requested update key first
	keyWaitingApprovalUpdate, err := u.keyRepo.GetKey(ctx, kv.Key, keyentity.PlacedKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if len(keyWaitingApprovalUpdate) >= 1 {
		return errors.New("Detected more than one pending approval key")
	}

	//get requested delete key
	keyWaitingApprovalDelete, err := u.keyRepo.GetKey(ctx, kv.Key, keyentity.PlacedDeleteKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if len(keyWaitingApprovalDelete) >= 1 {
		return errors.New("Detected more than one pending approval key")
	}

	keyInCanary, err := u.keyRepo.GetKey(ctx, kv.Key, keyentity.CanaryKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if len(keyInCanary) >= 1 {
		return errors.New("Can not delete value in canary.")
	}

	tx, err := u.keyRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// all keys that newly updated will have placedDelete status
	err = u.keyRepo.CreateKey(ctx, tx, kv.Key, kv.Value, kv.Type, kv.CreatedBy, keyentity.PlacedDeleteKey)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Usecase) ApproveKeyWithTx(tx *sql.Tx, key string, userID, status int) error {
	ctx := context.Background()

	// check if keys placed if no keys placed return error
	keyPlaced, err := u.keyRepo.GetKey(ctx, key, keyentity.PlacedKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return errors.New("no keys pending approval")
		}
	}

	if len(keyPlaced) > 1 {
		return errors.New("no more than 1 placed key can available")
	}
	if len(keyPlaced) == 0 {
		keyCanary, err := u.keyRepo.GetKey(ctx, key, keyentity.CanaryKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return errors.New("no keys pending approval")
			}
		}

		if len(keyCanary) == 0 {
			return errors.New("no keys pending approval")
		}

		keyPlaced = keyCanary
	}

	modifiedKey := keyPlaced[0]
	modifiedKey.ApprovedBy = userID
	modifiedKey.UpdateTime = time.Now()

	// Destroy all canary ip if any
	if modifiedKey.Status == keyentity.CanaryKey {
		if err := u.keyRepo.ModifyCanaryKey(ctx, tx, modifiedKey.ID, keyentity.StatusInactive); err != nil {
			return err
		}
	}

	if status == keyentity.DissaprovedKey {
		modifiedKey.Status = keyentity.DissaprovedKey
		return u.keyRepo.ModifyKey(ctx, tx, modifiedKey.ID, modifiedKey)
	}

	// modify current key to approved status and create new approved and active status
	modifiedKey.Status = keyentity.ApprovedKey
	err = u.keyRepo.ModifyKey(ctx, tx, keyPlaced[0].ID, modifiedKey)
	if err != nil {
		return err
	}

	// Change all old approve and active to approve and expire
	err = u.keyRepo.ModifyOldActiveKey(ctx, tx, key)
	if err != nil {
		return err
	}

	modifiedKey.Status = keyentity.ApprovedAndActive
	err = u.keyRepo.CreateKeyEntry(ctx, tx, modifiedKey)
	if err != nil {
		return err
	}

	return u.keyRepo.SetCache(ctx, modifiedKey)
}

func (u *Usecase) ApproveKey(key string, userID, status int) error {
	ctx := context.Background()

	// check if keys placed if no keys placed return error
	keyPlaced, err := u.keyRepo.GetKey(ctx, key, keyentity.PlacedKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return errors.New("no keys pending approval")
		}
	}

	if len(keyPlaced) > 1 {
		return errors.New("no more than 1 placed key can available")
	}
	if len(keyPlaced) == 0 {
		keyCanary, err := u.keyRepo.GetKey(ctx, key, keyentity.CanaryKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return errors.New("no keys pending approval")
			}
		}

		if len(keyCanary) == 0 {
			return errors.New("no keys pending approval")
		}

		keyPlaced = keyCanary
	}

	modifiedKey := keyPlaced[0]
	modifiedKey.ApprovedBy = userID
	modifiedKey.UpdateTime = time.Now()

	tx, err := u.keyRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Destroy all canary ip if any
	if modifiedKey.Status == keyentity.CanaryKey {
		if err := u.keyRepo.ModifyCanaryKey(ctx, tx, modifiedKey.ID, keyentity.StatusInactive); err != nil {
			return err
		}
	}

	if status == keyentity.DissaprovedKey {
		modifiedKey.Status = keyentity.DissaprovedKey

		err = u.keyRepo.ModifyKey(ctx, tx, modifiedKey.ID, modifiedKey)
		if err != nil {
			return err
		}

		return tx.Commit()
	}

	// modify current key to approved status and create new approved and active status
	modifiedKey.Status = keyentity.ApprovedKey
	err = u.keyRepo.ModifyKey(ctx, tx, keyPlaced[0].ID, modifiedKey)
	if err != nil {
		return err
	}

	// Change all old approve and active to approve and expire
	err = u.keyRepo.ModifyOldActiveKey(ctx, tx, key)
	if err != nil {
		return err
	}

	modifiedKey.Status = keyentity.ApprovedAndActive
	err = u.keyRepo.CreateKeyEntry(ctx, tx, modifiedKey)
	if err != nil {
		return err
	}

	err = u.keyRepo.SetCache(ctx, modifiedKey)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Usecase) ApproveDeleteKey(key string, userID, status int) error {
	ctx := context.Background()

	// check if keys placed if no keys placed return error
	keyPlaced, err := u.keyRepo.GetKey(ctx, key, keyentity.PlacedDeleteKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return errors.New("no keys pending approval")
		}
	}

	if len(keyPlaced) > 1 {
		return errors.New("no more than 1 placed delete key can available")
	}

	if len(keyPlaced) < 1 {
		return errors.New("no keys pending approval")
	}

	modifiedKey := keyPlaced[0]
	modifiedKey.ApprovedBy = userID
	modifiedKey.UpdateTime = time.Now()

	tx, err := u.keyRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Destroy all canary ip if any
	if modifiedKey.Status == keyentity.CanaryKey {
		if err := u.keyRepo.ModifyCanaryKey(ctx, tx, modifiedKey.ID, keyentity.StatusInactive); err != nil {
			return err
		}
	}

	if status == keyentity.DissaprovedKey {
		modifiedKey.Status = keyentity.DissaprovedKey
		err = u.keyRepo.ModifyKey(ctx, tx, modifiedKey.ID, modifiedKey)
		if err != nil {
			return err
		}

		return tx.Commit()
	}

	// modify current key to approved status
	modifiedKey.Status = keyentity.ApprovedKey
	err = u.keyRepo.ModifyKey(ctx, tx, keyPlaced[0].ID, modifiedKey)
	if err != nil {
		return err
	}

	// Change all old approve and active to approve and expire
	err = u.keyRepo.ModifyOldActiveKey(ctx, tx, key)
	if err != nil {
		return err
	}

	modifiedKey.Status = keyentity.DeletedKey
	err = u.keyRepo.CreateKeyEntry(ctx, tx, modifiedKey)
	if err != nil {
		return err
	}

	err = u.keyRepo.SetCache(ctx, modifiedKey)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Usecase) ApproveKeyCanary(key string, userID, status int, nodesIP []string) error {
	ctx := context.Background()

	var isFirstTimeCanary bool

	// check if already in canary before
	keyCanary, err := u.keyRepo.GetKey(ctx, key, keyentity.CanaryKey)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	if len(keyCanary) == 0 {
		keyPlaced, err := u.keyRepo.GetKey(ctx, key, keyentity.PlacedKey)
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}
		}

		if len(keyPlaced) == 0 {
			return errors.New("no keys pending approval")
		}

		isFirstTimeCanary = true
		keyCanary = keyPlaced
	}

	// modify current key to approved status and create new approved and active status
	approvedKeyEntry := keyCanary[0]
	approvedKeyEntry.ApprovedBy = userID
	approvedKeyEntry.UpdateTime = time.Now()
	approvedKeyEntry.Status = keyentity.CanaryKey

	tx, err := u.keyRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, ip := range nodesIP {
		err = u.keyRepo.CreateCanaryKey(ctx, tx, approvedKeyEntry.ID, ip)
		if err != nil {
			return err
		}
	}

	if isFirstTimeCanary {
		// first time canary, modify old key to status canary
		err = u.keyRepo.ModifyKey(ctx, tx, keyCanary[0].ID, approvedKeyEntry)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (u *Usecase) DeleteKey(keyID, userID int) error {
	ctx := context.Background()

	keyFetched, err := u.keyRepo.GetKeyByID(ctx, keyID)
	if err != nil {
		return err
	}

	// TODO: get user by id and verify if the user can modify and then change status to approved
	keyFetched.Status = keyentity.ApprovedAndExpiredKey

	tx, err := u.keyRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = u.keyRepo.ModifyKey(ctx, tx, keyID, keyFetched)
	if err != nil {
		return err
	}

	err = u.keyRepo.InvalidateCache(ctx, keyFetched.Key)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Usecase) GetHistoryKey(key string, isPrefix bool, limit int) ([]keyentity.KV, error) {
	ctx := context.Background()

	return u.keyRepo.GetKeyHistory(ctx, key, isPrefix, limit)
}

func (u *Usecase) GetKey(key string) (keyentity.KV, error) {
	ctx := context.Background()

	// get from cache first and if failed get from db
	keyFromCache, err := u.keyRepo.GetCache(ctx, key)
	if err != nil {
		keysActive, err := u.keyRepo.GetKey(ctx, key, keyentity.ApprovedKey)
		if err != nil {
			return keyentity.KV{}, err
		}
		if len(keysActive) > 1 {
			return keyentity.KV{}, errors.New("Detected more than one active keys")
		} else if len(keysActive) == 0 {
			return keyentity.KV{}, errors.New("No key found")
		}

		return keysActive[0], nil
	}

	// get only approved key
	return keyFromCache, nil
}

func (u *Usecase) GetKeys(prefix, ip string) ([]keyentity.KV, error) {
	ctx := context.Background()

	// get only approved key
	approvedKeys, err := u.keyRepo.GetCaches(ctx, prefix)
	if err != nil {
		return nil, err
	}

	// get value for specific ip
	if ip != "" {
		canaryKeys, err := u.keyRepo.GetCanaryKV(ctx, ip)
		if err != nil {
			return nil, err
		}

		for i, aKey := range approvedKeys {
			for _, cKey := range canaryKeys {
				if aKey.Key == cKey.Key {
					approvedKeys[i] = cKey
				}
			}
		}
	}

	return approvedKeys, nil
}

func (u *Usecase) BrowseKeys(prefix string) ([]string, error) {
	return u.keyRepo.GetKeyListWithoutValue(prefix)
}

func (u *Usecase) PendingApprovalKey(prefix string) ([]keyentity.KV, error) {
	ctx := context.Background()

	// get key with placed status basen on the prefix
	placedKeys, err := u.keyRepo.GetKeyByPrefix(ctx, prefix, keyentity.PlacedKey)
	if err != nil {
		return []keyentity.KV{}, err
	}

	// if keys with placed status found, return
	if len(placedKeys) > 0 {
		return placedKeys, nil
	}

	// if keys with placed status not found, get placed delete status
	placedRejectKeys, err := u.keyRepo.GetKeyByPrefix(ctx, prefix, keyentity.PlacedDeleteKey)
	if err != nil {
		return []keyentity.KV{}, err
	}

	return placedRejectKeys, nil
}

// Create service will create key, role user, role admin, and mapping user as lead for that service
func (u *Usecase) CreateService(username, tribe, service string) error {
	ctx := context.Background()
	key := fmt.Sprintf("service/%s/%s/default", tribe, service)
	prefix := fmt.Sprintf("service/%s/%s", tribe, service)

	if u.keyRepo.IsKeyExist(ctx, key) {
		return errors.New("Service already exist.")
	}

	user, err := u.userRepo.GetUser(ctx, username)
	if err != nil {
		return err
	}

	tx, err := u.keyRepo.GetDBTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := u.keyRepo.CreateKey(ctx, tx, key, "false", "bool", user.ID, keyentity.PlacedKey); err != nil {
		return err
	}

	if err := u.ApproveKeyWithTx(tx, key, user.ID, keyentity.ApprovedAndActive); err != nil {
		return err
	}

	if _, err := u.userRepo.CreateRole(ctx, tx, prefix, userentity.RoleSuperUser, user.ID); err != nil {
		return err
	}

	if _, err := u.userRepo.CreateRole(ctx, tx, prefix, userentity.RoleUser, user.ID); err != nil {
		return err
	}

	id, err := u.userRepo.CreateRole(ctx, tx, prefix, userentity.RoleLead, user.ID)
	if err != nil {
		return err
	}

	var roles []userentity.Role
	roles = append(roles, userentity.Role{ID: id})

	if err := u.userRepo.MapUserAccess(ctx, tx, user.ID, roles); err != nil {
		return err
	}

	return tx.Commit()
}

func (u *Usecase) RegisterCanaryDeployment(service string, nodesIP []string) error {
	ctx := context.Background()

	return u.keyRepo.RegisterCanaryDeployment(ctx, service, nodesIP)
}

func (u *Usecase) ReleaseCanaryIP(service string) error {
	ctx := context.Background()

	return u.keyRepo.ReleaseCanaryIP(ctx, service)
}

func (u *Usecase) GetKeyCanaryIP(id int) ([]string, []string, error) {
	ctx := context.Background()
	var canaryIPs, recommendedIP []string

	// Get current canary IP from db
	canaryKVs, err := u.keyRepo.GetCanaryKVByID(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return canaryIPs, recommendedIP, err
	}
	for _, ckv := range canaryKVs {
		canaryIPs = append(canaryIPs, ckv.IP)
	}

	// Get recommended canary IP in redis from jenkins
	requestedKey, err := u.keyRepo.GetKeyByID(ctx, id)
	if err != nil && err != sql.ErrNoRows {
		return canaryIPs, recommendedIP, err
	}

	if requestedKey.Key != "" {
		path := strings.Split(requestedKey.Key, "/")

		var service string
		if len(path) > 3 {
			service = path[2]
		} else {
			service = path[len(path)-2]
		}

		ips := u.keyRepo.GetCanaryIP(ctx, service)
		for _, ip := range ips {
			recommendedIP = append(recommendedIP, ip)
		}
	}

	return canaryIPs, recommendedIP, nil
}

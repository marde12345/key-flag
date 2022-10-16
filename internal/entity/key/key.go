package key

import "time"

type KV struct {
	ID          int       `db:"id" json:"id"`
	Key         string    `db:"key" json:"key"`
	Value       string    `db:"value" json:"value"`
	Type        string    `db:"type" json:"type"`
	CreateTime  time.Time `db:"create_time" json:"create_time"`
	UpdateTime  time.Time `db:"update_time" json:"update_time"`
	CreatedBy   int       `db:"created_by" json:"created_by"`
	ApprovedBy  int       `db:"approved_by" json:"approved_by"`
	Status      int       `db:"status" json:"status"`
	CreateByStr string    `db:"created_by_str" json:"created_by_str"`
}

type CanaryKV struct {
	KeyID  int    `db:"key_id" json:"key_id"`
	IP     string `db:"ip" json:"ip"`
	Status int    `db:"status" json:"status"`
}

const (
	ApprovedAndExpiredKey = iota
	ApprovedKey
	ApprovedAndActive
	PlacedKey
	DissaprovedKey
	CanaryKey
	PlacedDeleteKey
	DeletedKey
)

const (
	RedisKeyCanaryDeployment = "deployment:canary"
)

const (
	StatusInactive = 0
	StatusActive   = 1
)

func (kv KV) StatusString() string {
	switch kv.Status {
	case ApprovedAndExpiredKey:
		return "expired"
	case PlacedKey:
		return "placed"
	case ApprovedKey:
		return "approved"
	case DissaprovedKey:
		return "disapproved"
	case ApprovedAndActive:
		return "active"
	case PlacedDeleteKey:
		return "placed delete"
	case DeletedKey:
		return "inactive"
	}

	return ""
}

package shareuser

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrShareNotFound = errors.New("share not found")
	ErrShareExpired  = errors.New("share expired")
	ErrInvalidShare  = errors.New("invalid share")
)

// ShareUserItem 定向分享实体
type ShareUserItem struct {
	ID                  string
	OwnerUserID         string
	OwnerUsername       string
	TargetUserID        string
	TargetWalletAddress string
	Name                string
	Path                string
	IsDir               bool
	Permissions         string
	ExpiresAt           *time.Time
	CreatedAt           time.Time
}

// NewShareUserItem 创建定向分享记录
func NewShareUserItem(ownerID, ownerUsername, targetID, targetWallet, path, name string, isDir bool, permissions string, expiresAt *time.Time) *ShareUserItem {
	now := time.Now()
	return &ShareUserItem{
		ID:                  uuid.NewString(),
		OwnerUserID:         ownerID,
		OwnerUsername:       ownerUsername,
		TargetUserID:        targetID,
		TargetWalletAddress: targetWallet,
		Name:                name,
		Path:                path,
		IsDir:               isDir,
		Permissions:         permissions,
		ExpiresAt:           expiresAt,
		CreatedAt:           now,
	}
}

// IsExpired 判断分享是否过期
func (s *ShareUserItem) IsExpired() bool {
	if s.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*s.ExpiresAt)
}

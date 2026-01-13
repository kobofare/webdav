package addressbook

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrContactNotFound    = errors.New("contact not found")
	ErrDuplicateGroupName = errors.New("group name already exists")
	ErrDuplicateWallet    = errors.New("wallet address already exists")
)

type Group struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt time.Time
}

type Contact struct {
	ID            string
	UserID        string
	GroupID       string
	Name          string
	WalletAddress string
	Tags          []string
	CreatedAt     time.Time
}

func NewGroup(userID, name string) (*Group, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("group name is required")
	}
	now := time.Now()
	return &Group{
		ID:        uuid.NewString(),
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
	}, nil
}

func NewContact(userID, groupID, name, walletAddress string, tags []string) (*Contact, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("contact name is required")
	}
	walletAddress = strings.TrimSpace(walletAddress)
	if walletAddress == "" {
		return nil, errors.New("wallet address is required")
	}
	now := time.Now()
	return &Contact{
		ID:            uuid.NewString(),
		UserID:        userID,
		GroupID:       groupID,
		Name:          name,
		WalletAddress: strings.ToLower(walletAddress),
		Tags:          tags,
		CreatedAt:     now,
	}, nil
}

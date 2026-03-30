package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByWalletAddress(ctx context.Context, walletAddress string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}

type User struct {
	Id            uuid.UUID            `gorm:"type:uuid;primary_key" json:"id"`
	Email         string               `gorm:"type:varchar(255);uniqueIndex:idx_email"    json:"-"`
	Password      []byte               `gorm:"type:varchar(255);"    json:"-"`
	IsVerified    bool                 `gorm:"type:bool;not null"    json:"is_verified"`
	AuthMethod    string               `gorm:"type:varchar(50);not null;default:'email'" json:"auth_method"`
	WalletAddress *string              `gorm:"type:varchar(255);uniqueIndex:idx_wallet_address" json:"wallet_address,omitempty"`
	Status        bool `gorm:"type:bool;default:true;not null" json:"status"`
	DeletedAt     time.Time            `gorm:"type:timestamp" json:"deleted_at"`
	CreatedAt     time.Time            `gorm:"type:timestamp;not null" json:"created_at"`
	UpdatedAt     time.Time            `gorm:"type:timestamp;not null" json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Id == uuid.Nil {
		u.Id = uuid.New()
	}

	return nil
}

type WalletNonce struct {
	ID        uint      `gorm:"type:uuid;primary_key" json:"id"`
	Address   string    `gorm:"type:varchar(42);index;not null"`
	Nonce     string    `gorm:"type:varchar(255);not null"`
	ExpiredAt time.Time `gorm:"type:timestamp" json:"expired_at"`
	Used      bool      `gorm:"default:false"`

	CreatedAt time.Time
}

type AuthenticationInfo struct {
	User      *User
	SessionId uuid.UUID
}

func NewAuthenticationInfo(user *User, sessionId uuid.UUID) *AuthenticationInfo {
	return &AuthenticationInfo{
		User:      user,
		SessionId: sessionId,
	}
}

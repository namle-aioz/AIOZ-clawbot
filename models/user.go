package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByGoogleSubject(ctx context.Context, googleSubject string) (*User, error)
	GetUserByWalletAddress(ctx context.Context, walletAddress string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
}

type User struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	FirstName     string    `gorm:"type:varchar(255);"    json:"first_name"`
	LastName      string    `gorm:"type:varchar(255);"    json:"last_name"`
	Email         string    `gorm:"type:varchar(255);"    json:"-"`
	DisplayEmail  string    `gorm:"type:varchar(255);"    json:"email"`
	Password      string    `gorm:"type:varchar(255);"    json:"-"`
	IsVerified    bool      `gorm:"type:bool;not null"    json:"is_verified"`
	AuthProvider  string    `gorm:"type:varchar(50);not null;default:'email'" json:"auth_provider"`
	GoogleSubject *string   `gorm:"type:varchar(255);uniqueIndex" json:"-"`
	WalletAddress *string   `gorm:"type:varchar(255);uniqueIndex" json:"wallet_address,omitempty"`
	AvatarURL     *string   `gorm:"type:text" json:"avatar_url,omitempty"`
	CreatedAt     time.Time `gorm:"type:timestamp;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"type:timestamp;not null" json:"updated_at"`
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

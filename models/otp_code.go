package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPCodeRepository interface {
	GetOTPCodeByEmail(ctx context.Context, email string) (*OTPCode, error)
	CreateOrUpdateOTPCode(ctx context.Context, otpCode *OTPCode) error
}

type OTPCode struct {
	Id          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Email       string    `gorm:"type:varchar(255);index;not null" json:"email"`
	Code        []byte    `gorm:"type:varchar(255);not null" json:"-"`
	IssuedAt    time.Time `gorm:"index;not null" json:"issued_at"`
	ExpiredAt   time.Time `gorm:"index;not null" json:"expired_at"`
	Attempts    int       `gorm:"not null;default:1"`
	MaxAttempts int       `gorm:"not null;default:3"`

	CreatedAt time.Time `gorm:"type:timestamp;not null" json:"created_at"`
}

func (o *OTPCode) BeforeCreate(tx *gorm.DB) error {
	if o.Id == uuid.Nil {
		o.Id = uuid.New()
	}

	return nil
}

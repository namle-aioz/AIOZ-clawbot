package db

import (
	"context"

	"gorm.io/gorm"

	"backend/models"
)

type otpCodeRepository struct {
	db *gorm.DB
}

func MustNewOTPCodeRepository(db *gorm.DB, init bool) models.OTPCodeRepository {
	if init {
		if err := db.AutoMigrate(&models.OTPCode{}); err != nil {
			panic(err)
		}
	}

	return &otpCodeRepository{db: db}
}

func (r *otpCodeRepository) GetOTPCodeByEmail(ctx context.Context, email string) (*models.OTPCode, error) {
	var otpCode models.OTPCode
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&otpCode).Error; err != nil {
		return nil, err
	}

	return &otpCode, nil
}

func (r *otpCodeRepository) CreateOrUpdateOTPCode(ctx context.Context, otpCode *models.OTPCode) error {
	result := r.db.WithContext(ctx).Save(otpCode)
	return result.Error
}

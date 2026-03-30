package db

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend/models"
)

type userRepository struct {
	db *gorm.DB
}

func MustNewUserRepository(db *gorm.DB, init bool) models.UserRepository {
	if init {
		if err := db.AutoMigrate(&models.User{}); err != nil {
			panic(err)
		}
	}

	return &userRepository{db: db}
}

func (r *userRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByWalletAddress(ctx context.Context, walletAddress string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("wallet_address = ?", strings.ToLower(walletAddress)).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

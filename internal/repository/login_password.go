package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"gorm.io/gorm"
)

type LoginPasswordRepo interface {
	GetLoginPasswordData(userID string) (*models.LoginPassword, error)
	SaveNewLoginPassword(*models.LoginPassword) error
}

type LPRepo struct {
	db *gorm.DB
}

func NewLPRepo(db *gorm.DB) *LPRepo {
	return &LPRepo{db: db}
}

func (lp *LPRepo) GetLoginPasswordData(userID string) (*models.LoginPassword, error) {
	var loginPassword models.LoginPassword
	if err := lp.db.Where("user_id = ?", userID).First(&loginPassword).Error; err != nil {
		return nil, fmt.Errorf("failed to get login-password list by username %s: %w", userID, err)
	}
	return &loginPassword, nil
}

func (lp *LPRepo) SaveNewLoginPassword(logPass *models.LoginPassword) error {
	if err := lp.db.Create(logPass).Error; err != nil {
		return fmt.Errorf("failed to add new login-password: %w", err)
	}
	return nil
}

package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"gorm.io/gorm"
)

type TextDataRepo interface {
	GetTextData(userID string) (*models.TextData, error)
	SaveNewTextData(*models.TextData) error
}

type TDRepo struct {
	db *gorm.DB
}

func NewTDRepo(db *gorm.DB) *TDRepo {
	return &TDRepo{db: db}
}

func (td *TDRepo) GetTextData(userID string) (*models.TextData, error) {
	var textData models.TextData
	if err := td.db.Where("user_id = ?", userID).First(&textData).Error; err != nil {
		return nil, fmt.Errorf("failed to get text data by username %s: %w", userID, err)
	}
	return &textData, nil
}

func (td *TDRepo) SaveNewTextData(textData *models.TextData) error {
	if err := td.db.Create(textData).Error; err != nil {
		return fmt.Errorf("failed to add new text data: %w", err)
	}
	return nil
}

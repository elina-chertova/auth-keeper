package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"gorm.io/gorm"
)

type BinaryDataRepo interface {
	GetBinaryData(userID string) ([]*models.BinaryData, error)
	SaveNewBinaryData(*models.BinaryData) error
}

type BDRepo struct {
	db *gorm.DB
}

func NewBDRepo(db *gorm.DB) *BDRepo {
	return &BDRepo{db: db}
}

func (bd *BDRepo) GetBinaryData(userID string) ([]*models.BinaryData, error) {
	var binaryData []*models.BinaryData
	if err := bd.db.Where("user_id = ?", userID).Find(&binaryData).Error; err != nil {
		return nil, fmt.Errorf("failed to get binary data by username %s: %w", userID, err)
	}
	return binaryData, nil
}

func (bd *BDRepo) SaveNewBinaryData(binaryData *models.BinaryData) error {
	if err := bd.db.Create(binaryData).Error; err != nil {
		return fmt.Errorf("failed to add new binary data: %w", err)
	}
	return nil
}

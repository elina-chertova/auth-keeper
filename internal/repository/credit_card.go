package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"gorm.io/gorm"
)

type CreditCardRepo interface {
	GetCreditCardList(userID string) ([]*models.CreditCard, error)
	SaveNewCreditCard(*models.CreditCard) error
}

type CCRepo struct {
	db *gorm.DB
}

func NewCCRepo(db *gorm.DB) *CCRepo {
	return &CCRepo{db: db}
}

func (cc *CCRepo) GetCreditCardList(userID string) ([]*models.CreditCard, error) {
	var creditCards []*models.CreditCard
	if err := cc.db.Where("user_id = ?", userID).Find(&creditCards).Error; err != nil {
		return nil, fmt.Errorf("failed to get credit cards by userID %s: %w", userID, err)
	}
	return creditCards, nil
}

func (cc *CCRepo) SaveNewCreditCard(creditCard *models.CreditCard) error {
	if err := cc.db.Create(creditCard).Error; err != nil {
		return fmt.Errorf("failed to add new credit card: %w", err)
	}
	return nil
}

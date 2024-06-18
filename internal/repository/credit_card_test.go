package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
	"testing"
)

func TestNewCCRepo(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *CCRepo
	}{
		{
			name: "Valid DB instance",
			args: args{
				db: &gorm.DB{},
			},
			want: &CCRepo{db: &gorm.DB{}},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, NewCCRepo(tt.args.db), "NewCCRepo(%v)", tt.args.db)
			},
		)
	}
}

func TestCCRepo_GetCreditCardList(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*models.CreditCard
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get credit card list for existing user",
			fields: fields{
				db: setupTestDBWithCreditCards(
					[]*models.CreditCard{
						{
							UserID:     "user1",
							CardNumber: "1234-5678-9876-5432",
							ExpiryDate: "12/24",
							CVV:        "123",
							CardHolder: "John Doe",
							Metadata:   "metadata1",
						},
						{
							UserID:     "user1",
							CardNumber: "4321-8765-6789-1234",
							ExpiryDate: "11/23",
							CVV:        "321",
							CardHolder: "Jane Doe",
							Metadata:   "metadata2",
						},
					},
				),
			},
			args: args{
				userID: "user1",
			},
			want: []*models.CreditCard{
				{
					UserID:     "user1",
					CardNumber: "1234-5678-9876-5432",
					ExpiryDate: "12/24",
					CVV:        "123",
					CardHolder: "John Doe",
					Metadata:   "metadata1",
				},
				{
					UserID:     "user1",
					CardNumber: "4321-8765-6789-1234",
					ExpiryDate: "11/23",
					CVV:        "321",
					CardHolder: "Jane Doe",
					Metadata:   "metadata2",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Get credit card list for non-existing user",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				userID: "nonexistentuser",
			},
			want:    []*models.CreditCard{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				cc := &CCRepo{
					db: tt.fields.db,
				}
				got, err := cc.GetCreditCardList(tt.args.userID)
				if !tt.wantErr(t, err, fmt.Sprintf("GetCreditCardList(%v)", tt.args.userID)) {
					return
				}
				for i := range got {
					tt.want[i].Model.ID = got[i].Model.ID
					tt.want[i].CreatedAt = got[i].CreatedAt
					tt.want[i].UpdatedAt = got[i].UpdatedAt
				}
				assert.Equalf(t, tt.want, got, "GetCreditCardList(%v)", tt.args.userID)
			},
		)
	}
}

func TestCCRepo_SaveNewCreditCard(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		creditCard *models.CreditCard
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Save new credit card",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				creditCard: &models.CreditCard{
					UserID:     "user1",
					CardNumber: "5678-1234-4321-8765",
					ExpiryDate: "10/25",
					CVV:        "456",
					CardHolder: "Alice Smith",
					Metadata:   "newmetadata",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				cc := &CCRepo{
					db: tt.fields.db,
				}
				tt.wantErr(
					t,
					cc.SaveNewCreditCard(tt.args.creditCard),
					fmt.Sprintf("SaveNewCreditCard(%v)", tt.args.creditCard),
				)
			},
		)
	}
}

func setupTestDBWithCreditCards(cards []*models.CreditCard) *gorm.DB {
	db := setupTestDB()
	for _, card := range cards {
		db.Create(card)
	}
	return db
}

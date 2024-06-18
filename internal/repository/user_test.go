package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestNewUserRepo(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *userRepo
	}{
		{
			name: "Valid DB instance",
			args: args{
				db: &gorm.DB{},
			},
			want: &userRepo{db: &gorm.DB{}},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, NewUserRepo(tt.args.db), "NewUserRepo(%v)", tt.args.db)
			},
		)
	}
}

func Test_userRepo_CreateUser(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Create valid user",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				user: &models.User{
					Username:    "testuser",
					Password:    "password123",
					Email:       "testuser@example.com",
					PersonalKey: []byte("personal_key"),
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ur := &userRepo{
					db: tt.fields.db,
				}
				tt.wantErr(
					t,
					ur.CreateUser(tt.args.user),
					fmt.Sprintf("CreateUser(%v)", tt.args.user),
				)
			},
		)
	}
}

func Test_userRepo_GetUserByUsername(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.User
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get existing user",
			fields: fields{
				db: setupTestDBWithUser(
					&models.User{
						Username:    "existinguser",
						Password:    "password123",
						Email:       "existinguser@example.com",
						PersonalKey: []byte("personal_key"),
					},
				),
			},
			args: args{
				username: "existinguser",
			},
			want: &models.User{
				Username:    "existinguser",
				Password:    "password123",
				Email:       "existinguser@example.com",
				PersonalKey: []byte("personal_key"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "Get non-existing user",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				username: "nonexistentuser",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ur := &userRepo{
					db: tt.fields.db,
				}
				got, err := ur.GetUserByUsername(tt.args.username)
				if !tt.wantErr(t, err, fmt.Sprintf("GetUserByUsername(%v)", tt.args.username)) {
					return
				}

				if got != nil {
					tt.want.Model.ID = got.Model.ID
					tt.want.CreatedAt = got.CreatedAt
					tt.want.UpdatedAt = got.UpdatedAt
				}
				assert.Equalf(t, tt.want, got, "GetUserByUsername(%v)", tt.args.username)
			},
		)
	}
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(
		&models.User{},
		&models.LoginPassword{},
		&models.TextData{},
		&models.CreditCard{},
		&models.BinaryData{},
	)
	return db
}

func setupTestDBWithUser(user *models.User) *gorm.DB {
	db := setupTestDB()
	db.Create(user)
	return db
}

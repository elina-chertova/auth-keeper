package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
	"testing"
)

func TestLPRepo_GetLoginPasswordData(t *testing.T) {
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
		want    []*models.LoginPassword
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get login-password data for existing user",
			fields: fields{
				db: setupTestDBWithLoginPasswords(
					[]*models.LoginPassword{
						{
							UserID:   "user1",
							Login:    "login1",
							Password: "password1",
							Metadata: "metadata1",
						},
						{
							UserID:   "user1",
							Login:    "login2",
							Password: "password2",
							Metadata: "metadata2",
						},
					},
				),
			},
			args: args{
				userID: "user1",
			},
			want: []*models.LoginPassword{
				{
					UserID:   "user1",
					Login:    "login1",
					Password: "password1",
					Metadata: "metadata1",
				},
				{
					UserID:   "user1",
					Login:    "login2",
					Password: "password2",
					Metadata: "metadata2",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Get login-password data for non-existing user",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				userID: "nonexistentuser",
			},
			want:    []*models.LoginPassword{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				lp := &LPRepo{
					db: tt.fields.db,
				}
				got, err := lp.GetLoginPasswordData(tt.args.userID)
				if !tt.wantErr(t, err, fmt.Sprintf("GetLoginPasswordData(%v)", tt.args.userID)) {
					return
				}
				for i := range got {
					tt.want[i].Model.ID = got[i].Model.ID
					tt.want[i].CreatedAt = got[i].CreatedAt
					tt.want[i].UpdatedAt = got[i].UpdatedAt
				}
				assert.Equalf(t, tt.want, got, "GetLoginPasswordData(%v)", tt.args.userID)
			},
		)
	}
}

func TestLPRepo_SaveNewLoginPassword(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		logPass *models.LoginPassword
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Save new login-password",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				logPass: &models.LoginPassword{
					UserID:   "user1",
					Login:    "newlogin",
					Password: "newpassword",
					Metadata: "newmetadata",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				lp := &LPRepo{
					db: tt.fields.db,
				}
				tt.wantErr(
					t,
					lp.SaveNewLoginPassword(tt.args.logPass),
					fmt.Sprintf("SaveNewLoginPassword(%v)", tt.args.logPass),
				)
			},
		)
	}
}

func TestNewLPRepo(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *LPRepo
	}{
		{
			name: "Valid DB instance",
			args: args{
				db: &gorm.DB{},
			},
			want: &LPRepo{db: &gorm.DB{}},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, NewLPRepo(tt.args.db), "NewLPRepo(%v)", tt.args.db)
			},
		)
	}
}

func setupTestDBWithLoginPasswords(logins []*models.LoginPassword) *gorm.DB {
	db := setupTestDB()
	for _, login := range logins {
		db.Create(login)
	}
	return db
}

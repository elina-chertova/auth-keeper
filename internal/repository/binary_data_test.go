package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
	"testing"
)

func setupTestDBWithBinaryData(data []*models.BinaryData) *gorm.DB {
	db := setupTestDB()
	for _, d := range data {
		db.Create(d)
	}
	return db
}

func TestNewBDRepo(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *BDRepo
	}{
		{
			name: "Valid DB instance",
			args: args{
				db: &gorm.DB{},
			},
			want: &BDRepo{db: &gorm.DB{}},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, NewBDRepo(tt.args.db), "NewBDRepo(%v)", tt.args.db)
			},
		)
	}
}

func TestBDRepo_GetBinaryData(t *testing.T) {
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
		want    []*models.BinaryData
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get binary data for existing user",
			fields: fields{
				db: setupTestDBWithBinaryData(
					[]*models.BinaryData{
						{
							UserID:   "user1",
							Content:  []byte("content1"),
							Metadata: "metadata1",
						},
						{
							UserID:   "user1",
							Content:  []byte("content2"),
							Metadata: "metadata2",
						},
					},
				),
			},
			args: args{
				userID: "user1",
			},
			want: []*models.BinaryData{
				{
					UserID:   "user1",
					Content:  []byte("content1"),
					Metadata: "metadata1",
				},
				{
					UserID:   "user1",
					Content:  []byte("content2"),
					Metadata: "metadata2",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Get binary data for non-existing user",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				userID: "nonexistentuser",
			},
			want:    []*models.BinaryData{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				bd := &BDRepo{
					db: tt.fields.db,
				}
				got, err := bd.GetBinaryData(tt.args.userID)
				if !tt.wantErr(t, err, fmt.Sprintf("GetBinaryData(%v)", tt.args.userID)) {
					return
				}
				for i := range got {
					tt.want[i].Model.ID = got[i].Model.ID
					tt.want[i].CreatedAt = got[i].CreatedAt
					tt.want[i].UpdatedAt = got[i].UpdatedAt
				}
				assert.Equalf(t, tt.want, got, "GetBinaryData(%v)", tt.args.userID)
			},
		)
	}
}

func TestBDRepo_SaveNewBinaryData(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		binaryData *models.BinaryData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Save new binary data",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				binaryData: &models.BinaryData{
					UserID:   "user1",
					Content:  []byte("new content"),
					Metadata: "newmetadata",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				bd := &BDRepo{
					db: tt.fields.db,
				}
				tt.wantErr(
					t,
					bd.SaveNewBinaryData(tt.args.binaryData),
					fmt.Sprintf("SaveNewBinaryData(%v)", tt.args.binaryData),
				)
			},
		)
	}
}

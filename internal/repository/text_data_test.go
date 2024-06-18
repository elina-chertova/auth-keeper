package repository

import (
	"fmt"
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"github.com/stretchr/testify/assert"

	"gorm.io/gorm"
	"testing"
)

func TestNewTDRepo(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want *TDRepo
	}{
		{
			name: "Valid DB instance",
			args: args{
				db: &gorm.DB{},
			},
			want: &TDRepo{db: &gorm.DB{}},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				assert.Equalf(t, tt.want, NewTDRepo(tt.args.db), "NewTDRepo(%v)", tt.args.db)
			},
		)
	}
}

func TestTDRepo_GetTextData(t *testing.T) {
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
		want    []*models.TextData
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Get text data for existing user",
			fields: fields{
				db: setupTestDBWithTextData(
					[]*models.TextData{
						{
							UserID:   "user1",
							Content:  "content1",
							Metadata: "metadata1",
						},
						{
							UserID:   "user1",
							Content:  "content2",
							Metadata: "metadata2",
						},
					},
				),
			},
			args: args{
				userID: "user1",
			},
			want: []*models.TextData{
				{
					UserID:   "user1",
					Content:  "content1",
					Metadata: "metadata1",
				},
				{
					UserID:   "user1",
					Content:  "content2",
					Metadata: "metadata2",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Get text data for non-existing user",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				userID: "nonexistentuser",
			},
			want:    []*models.TextData{},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				td := &TDRepo{
					db: tt.fields.db,
				}
				got, err := td.GetTextData(tt.args.userID)
				if !tt.wantErr(t, err, fmt.Sprintf("GetTextData(%v)", tt.args.userID)) {
					return
				}
				for i := range got {
					tt.want[i].Model.ID = got[i].Model.ID
					tt.want[i].CreatedAt = got[i].CreatedAt
					tt.want[i].UpdatedAt = got[i].UpdatedAt
				}
				assert.Equalf(t, tt.want, got, "GetTextData(%v)", tt.args.userID)
			},
		)
	}
}

func TestTDRepo_SaveNewTextData(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		textData *models.TextData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Save new text data",
			fields: fields{
				db: setupTestDB(),
			},
			args: args{
				textData: &models.TextData{
					UserID:   "user1",
					Content:  "newcontent",
					Metadata: "newmetadata",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Save duplicate text data",
			fields: fields{
				db: setupTestDBWithTextData(
					[]*models.TextData{
						{
							UserID:   "user1",
							Content:  "newcontent",
							Metadata: "newmetadata",
						},
					},
				),
			},
			args: args{
				textData: &models.TextData{
					UserID:   "user1",
					Content:  "newcontent",
					Metadata: "newmetadata",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				td := &TDRepo{
					db: tt.fields.db,
				}
				tt.wantErr(
					t,
					td.SaveNewTextData(tt.args.textData),
					fmt.Sprintf("SaveNewTextData(%v)", tt.args.textData),
				)
			},
		)
	}
}

func setupTestDBWithTextData(data []*models.TextData) *gorm.DB {
	db := setupTestDB()
	for _, d := range data {
		db.Create(d)
	}
	return db
}

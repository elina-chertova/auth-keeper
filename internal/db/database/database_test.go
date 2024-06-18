package database

import (
	"github.com/elina-chertova/auth-keeper.git/internal/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"reflect"
	"testing"
)

type OriginalInitDBType func(conf *DBConfig) *gorm.DB

var OriginalInitDB OriginalInitDBType = InitDB

func MockInitDB(conf *DBConfig) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error during test database initialization: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.CreditCard{},
		&models.BinaryData{},
		&models.TextData{},
		&models.LoginPassword{},
	)
	if err != nil {
		log.Fatalf("Error during test migration: %v", err)
	}
	return db
}

func Test_InitDB(t *testing.T) {
	type args struct {
		conf *DBConfig
	}
	tests := []struct {
		name string
		args args
		want *gorm.DB
	}{
		{
			name: "Valid configuration",
			args: args{
				conf: &DBConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "testuser",
					Password: "password",
					DBName:   "testdb",
				},
			},
			want: MockInitDB(&DBConfig{}),
		},
	}

	InitDB := MockInitDB
	defer func() { InitDB = OriginalInitDB }()
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := InitDB(tt.args.conf)
				if got == nil {
					t.Errorf("InitDB() returned nil")
					return
				}

				sqlDB, err := got.DB()
				if err != nil {
					t.Errorf("InitDB() error getting SQL DB: %v", err)
					return
				}
				if err = sqlDB.Ping(); err != nil {
					t.Errorf("InitDB() error pinging database: %v", err)
					return
				}

				expectedTables := []string{
					"users",
					"credit_cards",
					"binary_data",
					"text_data",
					"login_passwords",
				}
				for _, table := range expectedTables {
					if !got.Migrator().HasTable(table) {
						t.Errorf("InitDB() missing expected table: %s", table)
					}
				}

				wantTables, _ := tt.want.Migrator().GetTables()
				gotTables, _ := got.Migrator().GetTables()
				if !reflect.DeepEqual(gotTables, wantTables) {
					t.Errorf("InitDB() tables = %v, want %v", gotTables, wantTables)
				}
			},
		)
	}
}

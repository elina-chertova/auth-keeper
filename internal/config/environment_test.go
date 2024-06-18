package config

import (
	"github.com/elina-chertova/auth-keeper.git/internal/db/database"
	"github.com/joho/godotenv"
	"os"
	"reflect"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		want      database.DBConfig
		want1     AppConf
		secretKey string
	}{
		{
			name: "Valid environment variables",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "testuser",
				"DB_NAME":     "testdb",
				"DB_PASSWORD": "password",
				"APP_ADDRESS": "localhost:8080",
				"SECRET_KEY":  "mysecretkey",
			},
			want: database.DBConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "testuser",
				DBName:   "testdb",
				Password: "password",
			},
			want1: AppConf{
				Address: "localhost:8080",
			},
			secretKey: "mysecretkey",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				f, err := os.Create(".env")
				if err != nil {
					t.Fatalf("Error creating .env file: %v", err)
				}

				for key, value := range tt.env {
					_, err := f.WriteString(key + "=" + value + "\n")
					if err != nil {
						t.Fatalf("Error writing to .env file: %v", err)
					}
				}

				if err := godotenv.Load(); err != nil {
					t.Fatalf("Error loading .env file: %v", err)
				}

				got, got1 := LoadEnv()

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("LoadEnv() got = %v, want %v", got, tt.want)
				}
				if !reflect.DeepEqual(got1, tt.want1) {
					t.Errorf("LoadEnv() got1 = %v, want %v", got1, tt.want1)
				}

				if SecretKey != tt.secretKey {
					t.Errorf("LoadEnv() SecretKey = %v, want %v", SecretKey, tt.secretKey)
				}

				for key := range tt.env {
					os.Unsetenv(key)
				}
			},
		)
	}
}

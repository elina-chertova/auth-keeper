package security

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecryptData(t *testing.T) {
	key := []byte("a very very very very secret key")
	data := []byte("exampleplaintext")
	type args struct {
		ciphertext []byte
		key        []byte
	}
	encryptedData, err := EncryptData(data, key)
	if err != nil {
		t.Fatalf("Failed to encrypt data: %v", err)
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid decryption",
			args: args{
				ciphertext: encryptedData,
				key:        key,
			},
			want:    data,
			wantErr: assert.NoError,
		},
		{
			name: "Invalid key",
			args: args{
				ciphertext: encryptedData,
				key:        []byte("wrong key"),
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := DecryptData(tt.args.ciphertext, tt.args.key)
				fmt.Println("got", got, err)
				if !tt.wantErr(
					t,
					err,
					fmt.Sprintf("DecryptData(%v, %v)", tt.args.ciphertext, tt.args.key),
				) {
					return
				}
				assert.Equalf(
					t,
					tt.want,
					got,
					"DecryptData(%v, %v)",
					tt.args.ciphertext,
					tt.args.key,
				)
			},
		)
	}
}

func TestDecryptPersonalKey(t *testing.T) {
	type args struct {
		encryptedKey []byte
	}
	key := []byte("a very very very very secret key")

	encryptedKey, err := EncryptPersonalKey(key)
	if err != nil {
		t.Fatalf("Failed to encrypt personal key: %v", err)
	}

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid decryption",
			args: args{
				encryptedKey: encryptedKey,
			},
			want:    key,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := DecryptPersonalKey(tt.args.encryptedKey)
				if !tt.wantErr(
					t,
					err,
					fmt.Sprintf("DecryptPersonalKey(%v)", tt.args.encryptedKey),
				) {
					return
				}
				assert.Equalf(t, tt.want, got, "DecryptPersonalKey(%v)", tt.args.encryptedKey)
			},
		)
	}
}

func TestEncryptData(t *testing.T) {
	type args struct {
		data []byte
		key  []byte
	}
	key := []byte("a very very very very secret key")
	data := []byte("exampleplaintext")

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid encryption",
			args: args{
				data: data,
				key:  key,
			},
			wantErr: assert.NoError,
		},
		{
			name: "Invalid key",
			args: args{
				data: data,
				key:  []byte("short key"),
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := EncryptData(tt.args.data, tt.args.key)
				if !tt.wantErr(
					t,
					err,
					fmt.Sprintf("EncryptData(%v, %v)", tt.args.data, tt.args.key),
				) {
					return
				}
			},
		)
	}
}

func TestEncryptPersonalKey(t *testing.T) {
	type args struct {
		personalKey []byte
	}
	key := []byte("a very very very very secret key") // 32 bytes

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Valid encryption",
			args: args{
				personalKey: key,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				_, err := EncryptPersonalKey(tt.args.personalKey)
				if !tt.wantErr(t, err, fmt.Sprintf("EncryptPersonalKey(%v)", tt.args.personalKey)) {
					return
				}
			},
		)
	}
}

func TestGeneratePersonalKey(t *testing.T) {
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Generate key",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got, err := GeneratePersonalKey()
				if !tt.wantErr(t, err, "GeneratePersonalKey()") {
					return
				}
				assert.Equal(t, 32, len(got), "GeneratePersonalKey() length")
			},
		)
	}
}

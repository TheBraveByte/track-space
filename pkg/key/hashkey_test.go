package key

import (
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
)

func TestHashPassword(t *testing.T) {

	tests := []struct {
		name           string
		inputPassword  string
		expectedString string
	}{
		{"hash-password", "$2a$14$ST9BpDyKW0ZbQPDCmWfDZOv5cefJqz8jQcBtu7yeY.ydoLDlvS/de", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.inputPassword) >= 1 {

				Byte, err := bcrypt.GenerateFromPassword([]byte(tt.inputPassword), 14)
				if err != nil {
					panic(err)
				}
				tt.expectedString = string(Byte)
				if got := HashPassword(tt.inputPassword); reflect.DeepEqual(got, tt.expectedString) {
					t.Errorf("HashPassword() = %v, want %v", got, tt.expectedString)
				}
			}

		})
	}
}

func TestVerifyPassword(t *testing.T) {

	tests := []struct {
		name           string
		inputPassword  string
		hashedPassword string
		want           bool
		want1          string
	}{
		{"valid-password", "track-space", "$2a$14$ST9BpDyKW0ZbQPDCmWfDZOv5cefJqz8jQcBtu7yeY.ydoLDlvS/de", true, "password hashed successfully"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := VerifyPassword(tt.inputPassword, tt.hashedPassword)
			if got != tt.want {
				t.Errorf("VerifyPassword() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("VerifyPassword() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

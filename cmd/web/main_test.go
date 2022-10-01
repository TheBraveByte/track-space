package main

import (
	"github.com/go-playground/assert/v2"
	"os"
	"testing"
)

func TestEnvValue(t *testing.T) {

	err := os.Setenv("MONGODB_URI", "mongodb+srv://ayaaakinleye:2701Akin2000@cluster0.byrpjo8.mongodb.net/test")
	err = os.Setenv("PORT_NUMBER", ":8080")
	err = os.Setenv("TOKEN", "$2a$12$qOlXAyF.dtYFFBi/2UKA3uKN2y98lkHWP3X.LFvMOcAEFHbgC1r7.")
	err = os.Setenv("TOKEN_SCHEME", "BearerToken")
	if err != nil {
		t.Errorf("Track-space failed setup application tes")

	}
	defer func() {
		err := os.Unsetenv("MONGODB_URI")
		err = os.Unsetenv("PORT_NUMBER")
		err = os.Unsetenv("TOKEN")
		err = os.Unsetenv("TOKEN_SCHEME")

		if err != nil {
			t.Errorf("cannot unsetenv env values")
		}
	}()
	tests := []struct {
		key           string
		expectedValue string
		defaultValue  string
		wantErr       bool
	}{
		{"MONGODB_URI", os.Getenv("MONGODB_URI"), "", false},
		{"TOKEN", os.Getenv("TOKEN"), "", false},
		{"PORT_NUMBER", os.Getenv("PORT_NUMBER"), "", false},
		{"TOKEN_SCHEME", os.Getenv("TOKEN_SCHEME"), "", false},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			if (err != nil) != tt.wantErr {
				t.Errorf("setUpApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.expectedValue, tt.expectedValue)
		})
	}
}

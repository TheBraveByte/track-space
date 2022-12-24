package main

import (
	"github.com/yusuf/track-space/pkg/model"
	"testing"
)

func TestSendMailToUser(t *testing.T) {

	tests := []struct {
		name string
		m    model.Email
	}{
		{name: "track-spaceAdmin", m: model.Email{

			Subject:  "Email confirmation",
			Content:  "Hello world",
			Receiver: "trackspace@admin.com",
			Sender:   "trackspace@admin.com",
			Template: "",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendMail(tt.m)
		})
	}
}

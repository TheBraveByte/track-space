package main

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func TestIsAuthorized(t *testing.T) {

	var tests []struct {
		middlewareFunc string
		expectedType   gin.HandlerFunc
	}
	for _, tt := range tests {
		t.Run(tt.middlewareFunc, func(t *testing.T) {
			if got := IsAuthorized(); !reflect.DeepEqual(got, tt.expectedType) {
				t.Errorf("IsAuthorized() = %v, expectedType %v", got, tt.expectedType)
			}
		})
	}
}

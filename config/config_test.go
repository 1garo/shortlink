package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		filename []string
		e error
	}{
		{[]string{".env", ".env.test"}, errors.New("passing more than 1 filename.")},
		{[]string{".env.not.found"}, errors.New("No .env file found")},
		{[]string{".env.empty"}, errors.New("MONGODB_URI not set.")},
		{[]string{"../.env"}, nil},
	}

	for _, tt := range testCases {
		_, err := NewConfig(tt.filename...)

		assert.Equal(t, tt.e, err)
	}
}

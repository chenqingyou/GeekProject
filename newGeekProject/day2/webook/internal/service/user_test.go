package service

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	passwd := "hello#go"
	password, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	bcrypt.CompareHashAndPassword(password, []byte(passwd))
	assert.NoError(t, err)
}

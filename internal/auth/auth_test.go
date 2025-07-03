package auth

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	rdm := make([]byte, 10)
	rand.Read(rdm)

	result, err := HashPassword(string(rdm))
	if err != nil {
		t.Error(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(result), rdm) != nil {
		t.Errorf("the password and its hash do not match")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	rdm := make([]byte, 10)
	rand.Read(rdm)

	expected, err := bcrypt.GenerateFromPassword(rdm, 0)
	if err != nil {
		t.Error(err)
	}

	if bcrypt.CompareHashAndPassword(expected, rdm) != CheckPasswordHash(string(rdm), string(expected)) {
		t.Errorf("hash didn't match")
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	mockPwd := make([]byte, 10)
	rand.Read(mockPwd)
	expected := uuid.New()
	signedToken, err := MakeJWT(expected, string(mockPwd), 1*time.Hour)
	if err != nil {
		t.Error("got unexpected error with MakeJWT")
	}
	result, err := ValidateJWT(signedToken, string(mockPwd))
	if err != nil {
		t.Error(err)
	}
	if result != expected {
		t.Error("userId didn't match")
	}
}

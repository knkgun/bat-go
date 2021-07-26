package cryptography

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeLimitedCredentialVerification(t *testing.T) {

	timeLimitedSecret := NewTimeLimitedSecret([]byte("1cdd2af5-1c5a-47f6-bdad-89090acee31c"))

	verifyTime, err := time.Parse("2006-01-02", "2020-12-23")
	assert.NoError(t, err)
	expirationTime, err := time.Parse("2006-01-02", "2020-12-30")
	assert.NoError(t, err)

	result, err := timeLimitedSecret.Derive(verifyTime, expirationTime)
	assert.NoError(t, err, "Error deriving token")

	correctVerifyTime, err := time.Parse("2006-01-02", "2020-12-23")
	assert.NoError(t, err)
	incorrectVerifyTime, err := time.Parse("2006-01-02", "2020-12-24")
	assert.NoError(t, err)

	correctlyVerified, err := timeLimitedSecret.Verify(correctVerifyTime, expirationTime, result)
	assert.NoError(t, err, "Error verifying token")
	assert.True(t, correctlyVerified, "Error verifying with correct verify time")

	incorrectlyVerified, err := timeLimitedSecret.Verify(incorrectVerifyTime, expirationTime, result)
	assert.NoError(t, err, "Error verifying token")
	assert.False(t, incorrectlyVerified, "Token should not have verified with incorrect verify time")
}

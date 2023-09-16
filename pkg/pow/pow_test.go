package pow_test

import (
	"testing"

	"github.com/SergeiSkv/pow/pkg/pow"
	"github.com/stretchr/testify/assert"
)

func TestSolvePoW(t *testing.T) {
	challenge := "abc"
	targetPrefix := "00"
	nonce, data := pow.SolvePoW(challenge, targetPrefix)

	assert.NotEmpty(t, nonce)
	assert.NotEmpty(t, data)
	assert.True(t, pow.ValidatePoW(challenge, nonce, data, targetPrefix))
}

func TestIsValidHash(t *testing.T) {
	data := "fu9mqkou45311"
	targetPrefix := "0000"
	hash := pow.ComputeHash(data)

	assert.True(t, pow.IsValidHash(hash[:], targetPrefix))
}

func TestGenerateChallenge(t *testing.T) {
	length := 10
	challenge := pow.GenerateChallenge(length)

	assert.Len(t, challenge, length)
}

func TestValidatePoW(t *testing.T) {
	challenge := "fu9mqkou45311"
	targetPrefix := "0000"
	nonce, data := pow.SolvePoW(challenge, targetPrefix)

	assert.True(t, pow.ValidatePoW(challenge, nonce, data, targetPrefix))
	assert.False(t, pow.ValidatePoW(challenge, nonce, "fake_data", targetPrefix))
	assert.False(t, pow.ValidatePoW(challenge, "fake_nonce", data, targetPrefix))
}

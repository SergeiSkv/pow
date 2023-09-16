package pow

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
)

func SolvePoW(challenge, targetPrefix string) (string, string) {
	for nonce := 0; ; nonce++ {
		data := fmt.Sprintf("%s%d", challenge, nonce)
		hash := ComputeHash(data)

		if IsValidHash(hash[:], targetPrefix) {
			return fmt.Sprintf("%d", nonce), data
		}
	}
}

func IsValidHash(hash []byte, targetPrefix string) bool {
	res := hex.EncodeToString(hash)[:len(targetPrefix)]
	return res == targetPrefix
}

func ComputeHash(data string) [32]byte {
	return sha256.Sum256([]byte(data))
}

func GenerateChallenge(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	challenge := ""
	for i := 0; i < length; i++ {
		challenge += string(chars[rand.Intn(len(chars))])
	}
	return challenge
}

func ValidatePoW(challenge, nonce, powRequest, targetPrefix string) bool {
	data := fmt.Sprintf("%s%s", challenge, nonce)
	hash := ComputeHash(data)

	return powRequest == data && IsValidHash(hash[:], targetPrefix)
}

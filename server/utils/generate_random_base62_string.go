package utils

import "crypto/rand"

var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomBase62String(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i := range length {
		bytes[i] = charset[bytes[i]%62]
	}

	return string(bytes), nil
}

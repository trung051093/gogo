package common

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
)

type CryptoUtils interface {
	GenerateRandomString(length int) string
	GenerateMD5(str string) string
	GenerateSHA512(str string, salt string) string
	GenerateSHA256(str string, salt string) string
	GenerateSHA1(str string, salt string) string
	CreatePassword(password string, length int) string
}

type cryptoUtils struct{}

var letters = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateMD5(str string, salt string) string {
	h := md5.New()
	h.Write([]byte(str))
	hash := hex.EncodeToString(h.Sum([]byte(salt)))
	return hash
}

func GenerateSHA512(str string, salt string) string {
	h := sha512.New()
	h.Write([]byte(str))
	hash := hex.EncodeToString(h.Sum([]byte(salt)))
	return hash
}

func GenerateSHA256(str string, salt string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hash := hex.EncodeToString(h.Sum([]byte(salt)))
	return hash
}

func GenerateSHA1(str string, salt string) string {
	h := sha1.New()
	h.Write([]byte(str))
	hash := hex.EncodeToString(h.Sum([]byte(salt)))
	return hash
}

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(length)]
	}
	return string(b)
}

func CreatePassword(password string, length int) (salt string, sha512 string) {
	salt = GenerateRandomString(length)
	sha512 = GenerateSHA512(password, salt)
	return salt, sha512
}

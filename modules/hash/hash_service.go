package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"math/rand"
)

type hashService struct{}

var letters = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func NewHashService() *hashService {
	return &hashService{}
}

func (hashS *hashService) GenerateMD5(str string, salt string) string {
	h := md5.New()
	h.Write([]byte(str))
	hash := hex.EncodeToString(h.Sum([]byte(salt)))
	return hash
}

func (hashS *hashService) GenerateSHA512(str string, salt string) string {
	h := sha512.New()
	h.Write([]byte(str))
	hash := hex.EncodeToString(h.Sum([]byte(salt)))
	return hash
}

func (hashS *hashService) GenerateSHA256(str string, salt string) string {
	h := sha256.New()
	h.Write([]byte(str))
	hash := hex.EncodeToString(h.Sum([]byte(salt)))
	return hash
}

func (hashS *hashService) GenerateSHA1(str string, salt string) string {
	h := sha1.New()
	h.Write([]byte(str))
	hash := hex.EncodeToString(h.Sum([]byte(salt)))
	return hash
}

func (hashS *hashService) GenerateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(length)]
	}
	return string(b)
}

func (hashS *hashService) CreatePassword(password string, length int) (salt string, sha512 string) {
	salt = hashS.GenerateRandomString(length)
	sha512 = hashS.GenerateSHA512(password, salt)
	return salt, sha512
}

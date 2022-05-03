package hash

type HashService interface {
	GenerateRandomString(length int) string
	GenerateMD5(str string, salt string) string
	GenerateSHA512(str string, salt string) string
	GenerateSHA256(str string, salt string) string
	GenerateSHA1(str string, salt string) string
	CreatePassword(password string, length int) (salt string, sha512 string)
}

package security

import "golang.org/x/crypto/bcrypt"

type BcryptService struct {
}

func NewBcryptService() *BcryptService {
	return &BcryptService{}
}

func (b *BcryptService) HashPassword(password string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err

}

func (b *BcryptService) CheckPassword(password, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err
}


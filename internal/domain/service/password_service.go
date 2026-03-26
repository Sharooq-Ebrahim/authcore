package service

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPassword(password, hash string) error
}

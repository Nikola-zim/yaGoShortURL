package service

import "log"

type AuthService struct {
	authCash AuthUser
}

func (aS *AuthService) FindUser(idMsg string) (uint64, bool) {
	return aS.authCash.FindUser(idMsg)
}

func (aS *AuthService) AddUser() (string, error) {
	log.Println("In service")
	cookie, err := aS.authCash.AddUser()
	return cookie, err
}

func NewAuthService(authCash AuthUser) *AuthService {
	return &AuthService{
		authCash: authCash,
	}
}

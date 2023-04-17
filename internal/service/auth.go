package service

import "log"

type AuthService struct {
	authCash AuthUser
}

func (aS *AuthService) FindUser(idMsg string) (uint64, bool) {
	return aS.authCash.FindUser(idMsg)
}

func (aS *AuthService) AddUser() (string, uint64, error) {
	log.Println("In service")
	cookie, id, err := aS.authCash.AddUser()
	return cookie, id, err
}

func NewAuthService(authCash AuthUser) *AuthService {
	return &AuthService{
		authCash: authCash,
	}
}

package service

import (
	"github.com/arganaphangquestian/usermanagement/server/model"
)

type (
	// UserService interface
	UserService interface {
		Register(register model.InputUser) (*model.User, error)
		Login(login model.Login) (*model.UserToken, error)
		Users() ([]*model.User, error)
		CreateSubcriber(token string) error
		RefreshToken(refeshToken string) (*model.UserToken, error)
		Verify(tokenString string) error
	}
)

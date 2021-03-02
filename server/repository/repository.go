package repository

import (
	"github.com/arganaphangquestian/usermanagement/server/model"
)

type (
	// UserRepository interface
	UserRepository interface {
		Register(register model.InputUser) (*model.User, error)
		GetUserByUsername(username string) (*model.UserWithPassword, error)
		Users() ([]*model.User, error)
	}
)

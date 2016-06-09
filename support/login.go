package support

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"

	"errors"
)

func Login(login string, password string, u model.UserInterface) (loginOut string, err error) {
	if len(login) == 0 {
		return "", errors.New("user name cannot be blank")
	}

	if len(password) == 0 {
		return "", errors.New("password cannot be blank")
	}

	loginOut, passwordHash := u.GetUserIdPasswordHashByLogin(login)

	if auth.CompareHashAndPassword(password, passwordHash) != nil {
		return "", errors.New("incorrect username or password")
	}

	return loginOut, nil
}

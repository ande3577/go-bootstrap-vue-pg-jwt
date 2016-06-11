package support

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"

	"errors"
)

func CreateUser(user model.UserInterface, session model.SessionInterface, password string, passwordConfirmation string) error {
	u := user.Get()

	if len(u.Login) == 0 {
		return errors.New("login cannot be blank")
	}

	if len(u.Email) == 0 {
		return errors.New("e-mail cannot be blank")
	}

	if len(password) == 0 {
		return errors.New("password cannot be blank")
	}

	if password != passwordConfirmation {
		return errors.New("password confirmation does not match")
	}

	u.HashedPassword = auth.GenerateHashFromPassword(password)

	if err := user.Create(); err != nil {
		return err
	}

	s := session.Get()

	s.UserId = u.Id
	return session.Create()
}

package support

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"

	"errors"
)

func updateUserPassword(u *model.User, password string, passwordConfirmation string) error {
	if len(password) == 0 {
		return errors.New("password cannot be blank")
	}

	if password != passwordConfirmation {
		return errors.New("password confirmation does not match")
	}

	u.HashedPassword = auth.GenerateHashFromPassword(password)
	return nil
}

func CreateUser(user model.UserInterface, session model.SessionInterface, password string, passwordConfirmation string) error {
	u := user.Get()

	if len(u.Login) == 0 {
		return errors.New("login cannot be blank")
	}

	if len(u.Email) == 0 {
		return errors.New("e-mail cannot be blank")
	}

	if err := updateUserPassword(u, password, passwordConfirmation); err != nil {
		return err
	}

	if err := user.Create(); err != nil {
		return err
	}

	s := session.Get()

	s.UserId = u.Id
	return session.Create()
}

func UpdateUser(user model.UserInterface, session model.SessionInterface, password string, passwordConfirmation string, fromHttp bool, developmentMode bool) (passwordChanged bool, tokenData *auth.TokenData, err error) {
	tokenData = &auth.TokenData{}

	if len(password) != 0 || len(passwordConfirmation) != 0 {
		if err = updateUserPassword(user.Get(), password, passwordConfirmation); err != nil {
			return false, tokenData, err
		}
		passwordChanged = true
	}

	if err = user.Update(); err != nil {
		return false, tokenData, err
	}

	if passwordChanged {
		user.DestroySessions()
		tokenData, err = CreateUserSession(user.Get(), session, fromHttp, developmentMode)
	}

	return passwordChanged, tokenData, err
}

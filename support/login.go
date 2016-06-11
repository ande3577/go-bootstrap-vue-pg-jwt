package support

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"

	"errors"
)

func Login(login string, password string, fromHttp bool, developmentMode bool, user model.UserInterface, sess model.SessionInterface) (tokenData *auth.TokenData, err error) {
	if len(login) == 0 {
		return tokenData, errors.New("user name cannot be blank")
	}

	if len(password) == 0 {
		return tokenData, errors.New("password cannot be blank")
	}

	if err := user.FindByLogin(login); err != nil {
		return tokenData, err
	}
	u := user.Get()

	if auth.CompareHashAndPassword(password, u.HashedPassword) != nil {
		return tokenData, errors.New("incorrect username or password")
	}

	if tokenData, err = auth.Login(u.Login, fromHttp, developmentMode); err != nil {
		return tokenData, err
	}

	s := sess.Get()
	s.UserId = u.Id
	s.Session = tokenData.SessionIdentifier

	err = sess.Create()

	return tokenData, err
}

func Logout(login string, sessionIdentifier string, user model.UserInterface, sess model.SessionInterface) {
	if err := user.FindByLogin(login); err != nil {
		return
	}

	if err := sess.FindBySession(sessionIdentifier); err != nil {
		return
	}

	if sess.Get().UserId == user.Get().Id {
		sess.Destroy()
	}
}

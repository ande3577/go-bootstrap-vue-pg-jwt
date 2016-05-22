package support

import (
	"errors"
)

func Login(userId string, password string) error {
	if len(userId) == 0 {
		return errors.New("user name cannot be blank")
	}

	if len(password) == 0 {
		return errors.New("password cannot be blank")
	}

	// to get framework in place, just check for userId == password
	if userId != password {
		return errors.New("incorrect password")
	}

	return nil
}

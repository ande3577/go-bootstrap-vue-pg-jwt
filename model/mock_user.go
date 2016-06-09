package model

type MockUser struct {
	User
	Created bool
	Deleted bool
	Err     error
}

func (u *MockUser) Create() error {
	u.Created = true
	return u.Err
}

func (u *MockUser) Destroy() error {
	u.Deleted = true
	return u.Err
}

func (u *MockUser) Get() *User {
	return &u.User
}

func (u *MockUser) GetUserIdPasswordHashByLogin(login string) (userId string, passwordHash string) {
	if login != u.Login && login != u.Email {
		return "", ""
	} else {
		return u.Login, u.HashedPassword
	}
}

func (u *MockUser) FindByLogin(login string) error {
	return nil
}

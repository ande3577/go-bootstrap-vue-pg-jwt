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

func (u *MockUser) FindByLogin(login string) error {
	return nil
}

package model

type MockUser struct {
	User
	Created           bool
	Deleted           bool
	Updated           bool
	SessionsDestroyed bool
	Err               error
}

func (u *MockUser) Create() error {
	u.Created = true
	return u.Err
}

func (u *MockUser) Update() error {
	u.Updated = true
	return u.Err
}

func (u *MockUser) DestroySessions() error {
	u.SessionsDestroyed = true
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

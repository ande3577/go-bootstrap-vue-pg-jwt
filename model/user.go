package model

import (
	"time"
)

type User struct {
	Id             int
	Login          string
	Email          string
	HashedPassword string    `db:"hashed_password"`
	CreatedAt      time.Time `db:"created_at""`
}

type UserInterface interface {
	Create() error
	Destroy() error
	Get() *User
	GetUserIdPasswordHashByLogin(login string) (userId string, passwordHash string)
	FindByLogin(login string) error
}

func (u *User) Create() error {
	return dbMap.Insert(u)
}

func (u *User) Destroy() error {
	_, err := dbMap.Delete(u)
	return err
}

func (u *User) Get() *User {
	return u
}

func (u *User) GetUserIdPasswordHashByLogin(login string) (userId string, passwordHash string) {
	if err := dbMap.SelectOne(u, "select login, hashed_password from users WHERE login=$1 OR email=$1", login); err != nil {
		return "", ""
	}
	return u.Login, u.HashedPassword
}

func (u *User) FindByLogin(login string) error {
	return dbMap.SelectOne(u, "select * from users WHERE login=$1 OR email=$1", login)
}

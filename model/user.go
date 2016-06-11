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
	Update() error
	Destroy() error
	Get() *User
	FindByLogin(login string) error
	DestroySessions() error
}

func (u *User) Create() error {
	return dbMap.Insert(u)
}

func (u *User) DestroySessions() error {
	var sessions []Session
	if _, err := dbMap.Select(&sessions, "select id from sessions where user_id=$1", u.Id); err != nil {
		return err
	}

	for _, s := range sessions {
		if _, err := dbMap.Delete(&s); err != nil {
			return err
		}
	}

	return nil
}

func (u *User) Update() error {
	_, err := dbMap.Update(u)
	return err
}

func (u *User) Destroy() error {
	if err := u.DestroySessions(); err != nil {
		return err
	}

	_, err := dbMap.Delete(u)
	return err
}

func (u *User) Get() *User {
	return u
}

func (u *User) FindByLogin(login string) error {
	return dbMap.SelectOne(u, "select * from users WHERE login=$1 OR email=$1", login)
}

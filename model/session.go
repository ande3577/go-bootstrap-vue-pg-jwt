package model

import (
	"time"
)

type Session struct {
	Id        int
	UserId    int `db:"user_id"`
	Session   string
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}

type SessionInterface interface {
	Create() error
	Destroy() error
	FindBySession(session string) error
	Get() *Session
}

func (s *Session) Create() error {
	return dbMap.Insert(s)
}

func (s *Session) Destroy() error {
	_, err := dbMap.Delete(s)
	return err
}

func (s *Session) Get() *Session {
	return s
}

func (s *Session) FindBySession(session string) error {
	err := dbMap.SelectOne(s, "select * from sessions where session=$1", session)
	return err
}

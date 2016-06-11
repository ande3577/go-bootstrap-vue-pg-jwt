package model

import (
	"errors"
)

type MockSession struct {
	Session
	CreateCalled  bool
	DestroyCalled bool
	Err           error
}

func (s *MockSession) Create() error {
	s.CreateCalled = true
	return nil
}

func (s *MockSession) Destroy() error {
	s.DestroyCalled = true
	return nil
}

func (s *MockSession) Get() *Session {
	return &s.Session
}

func (s *MockSession) FindBySession(session string) error {
	if session == s.Session.Session {
		return nil
	} else {
		return errors.New("session mismatch")
	}
}

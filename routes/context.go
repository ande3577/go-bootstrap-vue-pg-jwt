package routes

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/gorilla/sessions"
	"net/http"
)

type Context struct {
	Session         *sessions.Session
	User            string
	DevelopmentMode bool
	XSRFToken       string
}

func NewContext(r *http.Request) (c *Context, err error) {
	s, userId, _, xsrfToken, err := auth.Authorize(r, settings.DevelopmentMode)

	return &Context{
		Session:         s,
		User:            userId,
		XSRFToken:       xsrfToken,
		DevelopmentMode: settings.DevelopmentMode,
	}, err
}

func NewContextFromJson(r *http.Request) (c *Context, err error) {
	userIdString, _, xsrfToken, err := auth.AuthorizeJSON(r, settings.DevelopmentMode)

	return &Context{
		User:            userIdString,
		XSRFToken:       xsrfToken,
		DevelopmentMode: settings.DevelopmentMode,
	}, err
}

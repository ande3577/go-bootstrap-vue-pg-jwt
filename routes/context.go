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
	s, tokenData, err := auth.Authorize(r, settings.DevelopmentMode)

	return &Context{
		Session:         s,
		User:            tokenData.UserId,
		XSRFToken:       tokenData.XsrfToken,
		DevelopmentMode: settings.DevelopmentMode,
	}, err
}

func NewContextFromJson(r *http.Request) (c *Context, err error) {
	tokenData, err := auth.AuthorizeJSON(r, settings.DevelopmentMode)

	return &Context{
		User:            tokenData.UserId,
		XSRFToken:       tokenData.XsrfToken,
		DevelopmentMode: settings.DevelopmentMode,
	}, err
}

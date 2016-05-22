package routes

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/support"

	"net/http"
)

func logout(c *Context) {
	c.XSRFToken = auth.Logout(c.Session, settings.DevelopmentMode)
	c.User = ""
}

func PostLogin(w http.ResponseWriter, r *http.Request, c *Context) error {
	userId := r.FormValue("user_id")
	password := r.FormValue("password")

	if err := support.Login(userId, password); err != nil {
		c.Session.AddFlash(err.Error())
		logout(c)
		getIndexTemplate().Execute(w, map[string]interface{}{"Context": c})
		return nil
	}

	c.User = userId
	if tokenString, _, err := auth.Login(userId, true, settings.DevelopmentMode); err == nil {
		c.Session.Values["token"] = tokenString
	} else {
		c.Session.AddFlash(err.Error())
		logout(c)
		getIndexTemplate().Execute(w, map[string]interface{}{"Context": c})
		return nil
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func PostLogout(w http.ResponseWriter, r *http.Request, c *Context) error {
	logout(c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

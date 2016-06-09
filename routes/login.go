package routes

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/support"

	"html/template"
	"net/http"
	"path/filepath"
)

func logout(c *Context) {
	c.XSRFToken = auth.Logout(c.Session, settings.DevelopmentMode)
	c.User = ""
}

func login(c *Context) error {
	if tokenString, _, err := auth.Login(c.User, true, settings.DevelopmentMode); err == nil {
		c.Session.Values["token"] = tokenString
	} else {
		return err
	}
	return nil
}

func PostLogin(w http.ResponseWriter, r *http.Request, c *Context) (err error) {
	userId := r.FormValue("user_id")
	password := r.FormValue("password")

	if c.User, err = support.Login(userId, password, &model.User{}); err != nil {
		c.Session.AddFlash(err.Error())
		logout(c)
		getIndexTemplate().Execute(w, map[string]interface{}{"Context": c})
		return nil
	}

	if err := login(c); err != nil {
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

func getRegisterTemplate() *template.Template {
	lockTemplateStore()
	defer unlockTemplateStore()

	registerTemplate, ok := getTemplateFromStore("register")
	if !ok {
		registerTemplate = template.Must(template.ParseFiles(filepath.Join(settings.RootDirectory, "templates/_base.html"),
			filepath.Join(settings.RootDirectory, "templates/register.html")))
		addTemplateToStore("register", registerTemplate)
	}
	return registerTemplate
}

func GetRegister(w http.ResponseWriter, r *http.Request, c *Context) error {
	getRegisterTemplate().Execute(w, map[string]interface{}{"Context": c})
	return nil
}

func PostRegister(w http.ResponseWriter, r *http.Request, c *Context) error {
	u := &model.User{Login: r.FormValue("login"),
		Email: r.FormValue("email")}

	if err := support.CreateUser(u, r.FormValue("password-main"), r.FormValue("password-confirmation")); err != nil {
		logout(c)
		c.Session.AddFlash(err.Error())
		getRegisterTemplate().Execute(w, map[string]interface{}{"Context": c, "Login": u.Login, "Email": u.Email})
		return nil
	}

	c.User = u.Login

	login(c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

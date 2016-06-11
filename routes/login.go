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
	support.Logout(c.User, c.SessionIdentifier, &model.User{}, &model.Session{})
	tokenData := auth.Logout(c.Session, settings.DevelopmentMode)

	c.XSRFToken = tokenData.XsrfToken
	c.User = ""
}

func login(c *Context, login string, password string, fromHttp bool) (err error) {
	if tokenData, err := support.Login(login, password, fromHttp, settings.DevelopmentMode, &model.User{}, &model.Session{}); err != nil {
		return err
	} else {
		c.Session.Values["token"] = tokenData.TokenString
	}

	return nil
}

func PostLogin(w http.ResponseWriter, r *http.Request, c *Context) (err error) {
	userId := r.FormValue("user_id")
	password := r.FormValue("password")

	if err := login(c, userId, password, true); err != nil {
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

	s := &model.Session{}

	password := r.FormValue("password-main")
	err := support.CreateUser(u, s, password, r.FormValue("password-confirmation"))
	if err == nil {
		err = login(c, u.Login, password, true)
	}

	if err != nil {
		logout(c)
		c.Session.AddFlash(err.Error())
		getRegisterTemplate().Execute(w, map[string]interface{}{"Context": c, "Login": u.Login, "Email": u.Email})
		return nil
	}

	c.User = u.Login

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

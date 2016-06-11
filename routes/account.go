package routes

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/support"

	"html/template"
	"net/http"
	"path/filepath"
)

func getAccountTemplate() *template.Template {
	lockTemplateStore()
	defer unlockTemplateStore()

	registerTemplate, ok := getTemplateFromStore("account")
	if !ok {
		registerTemplate = template.Must(template.ParseFiles(filepath.Join(settings.RootDirectory, "templates/_base.html"),
			filepath.Join(settings.RootDirectory, "templates/account.html"),
			filepath.Join(settings.RootDirectory, "templates/_account_info_form.html")))
		addTemplateToStore("account", registerTemplate)
	}
	return registerTemplate
}

func validateLoggedIn(c *Context) bool {
	return len(c.User) > 0
}

func validateSession(c *Context) bool {
	if !validateLoggedIn(c) {
		return false
	}

	s := &model.Session{}
	if len(c.SessionIdentifier) == 0 {
		return false
	}

	if s.FindBySession(c.SessionIdentifier) != nil {
		return false
	}

	return true
}

func GetAccount(w http.ResponseWriter, r *http.Request, c *Context) error {
	if !validateSession(c) {
		denyAccess(w, r, c)
		return nil
	}

	user := &model.User{}
	if err := user.FindByLogin(c.User); err != nil {
		return err
	}

	u := user.Get()

	getAccountTemplate().Execute(w, map[string]interface{}{"Context": c, "Email": u.Email})
	return nil
}

func PostAccount(w http.ResponseWriter, r *http.Request, c *Context) error {
	if !validateSession(c) {
		denyAccess(w, r, c)
		return nil
	}

	user := &model.User{}
	if err := user.FindByLogin(c.User); err != nil {
		return err
	}

	user.Email = r.FormValue("email")
	password := r.FormValue("password-main")
	passwordConfirmation := r.FormValue("password-confirmation")

	passwordChanged, tokenData, err := support.UpdateUser(user, &model.Session{}, password, passwordConfirmation, true, settings.DevelopmentMode)
	if err != nil {
		c.Session.AddFlash(err.Error())
		getAccountTemplate().Execute(w, map[string]interface{}{"Context": c, "Email": user.Email})
		return nil
	}

	if passwordChanged {
		addTokenToSession(tokenData, c)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

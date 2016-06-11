package routes

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/app"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"

	"errors"
	"fmt"
	"github.com/goods/httpbuf"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var settings *app.ApplicationSettings

type handlerWithContext func(http.ResponseWriter, *http.Request, *Context) error
type jsonHandlerWithContext handlerWithContext

func writeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func denyAccess(w http.ResponseWriter, r *http.Request, c *Context) {
	c.Session.AddFlash(errors.New("access denied").Error())
	logout(c)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h jsonHandlerWithContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context, err := NewContextFromJson(r)
	if err != nil {
		writeError(w, err)
		return
	}

	buf := &httpbuf.Buffer{}

	if err := h(buf, r, context); err != nil {
		writeError(w, err)
	}

	buf.Apply(w)
}

func (h handlerWithContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := NewContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := &httpbuf.Buffer{}

	err = h(buf, r, c)
	if err != nil {
		writeError(w, err)
		return
	}

	if c.Session != nil {
		if err = c.Session.Save(r, buf); err != nil {
			writeError(w, err)
		}
	}

	buf.Apply(w)
}

func getBooleanFormValue(r *http.Request, name string) bool {
	value := strings.ToLower(strings.TrimSpace(r.FormValue(name)))
	value = strings.ToLower(strings.TrimSpace(value))

	return value == "1" || value == "true"
}

func getFloatFormValue(r *http.Request, name string) (value float64, ok bool) {
	valueString := r.FormValue(name)
	if len(valueString) == 0 {
		return 0, false
	}

	value, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		return 0, false
	}

	return value, true
}

func buildTemplatePath(templateName string) string {
	return filepath.Join(settings.RootDirectory, "templates", templateName)
}

func addTokenToSession(tokenData *auth.TokenData, c *Context) {
	c.Session.Values["token"] = tokenData.TokenString
}

func SetupApplication(s *app.ApplicationSettings) chan int {
	settings = s

	Initialize(settings.RootDirectory)

	db, err := app.OpenDB(settings)
	if err != nil {
		panic(err)
	}
	model.Initialize(db)

	serverChannel := make(chan int) // Allocate a channel.
	go func() {
		http.ListenAndServe(":"+settings.Port, nil)
		serverChannel <- 1 // Send a signal; value does not matter.
	}()
	fmt.Println("Starting... on port: " + settings.Port)
	return serverChannel
}

func Initialize(root string) {
	r := mux.NewRouter()
	r.Handle("/", handlerWithContext(Index)).Methods("GET")
	r.Handle("/login", handlerWithContext(PostLogin)).Methods("POST")
	r.Handle("/logout", handlerWithContext(PostLogout)).Methods("POST")
	r.Handle("/register", handlerWithContext(GetRegister)).Methods("GET")
	r.Handle("/register", handlerWithContext(PostRegister)).Methods("POST")
	r.Handle("/account", handlerWithContext(GetAccount)).Methods("GET")
	r.Handle("/account", handlerWithContext(PostAccount)).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(filepath.Join(root, "static")))).Methods("GET")
	http.Handle("/", r)
}

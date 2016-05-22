package routes

import (
	"fmt"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/app"
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

func (h jsonHandlerWithContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context, err := NewContextFromJson(r)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h(w, r, context); err != nil {
		writeError(w, err)
	}
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

func SetupApplication(s *app.ApplicationSettings) chan int {
	settings = s

	Initialize(settings.RootDirectory)

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
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(filepath.Join(root, "static"))))
	http.Handle("/", r)
}

package routes

import (
	"html/template"
	"sync"
)

var templateStoreMutex sync.Mutex = sync.Mutex{}
var templateStore map[string]*template.Template = make(map[string]*template.Template)

func lockTemplateStore() {
	templateStoreMutex.Lock()
}

func unlockTemplateStore() {
	templateStoreMutex.Unlock()
}

func addTemplateToStore(name string, t *template.Template) {
	if settings.DevelopmentMode {
		return
	}
	templateStore[name] = t
}

func getTemplateFromStore(name string) (t *template.Template, ok bool) {
	if settings.DevelopmentMode {
		return nil, false
	}

	t, ok = templateStore[name]
	return t, ok
}

package routes

import (
	"html/template"
	"net/http"
)

type Page struct {
	Title string
}

func getIndexTemplate() (t *template.Template) {
	lockTemplateStore()
	defer unlockTemplateStore()
	indexTemplate, ok := getTemplateFromStore("index")
	if !ok {
		indexTemplate = template.Must(template.ParseFiles(buildTemplatePath("_base.html"),
			buildTemplatePath("index.html")))
		addTemplateToStore("index", indexTemplate)
	}
	return indexTemplate
}

func Index(w http.ResponseWriter, r *http.Request, c *Context) error {
	getIndexTemplate().Execute(w, map[string]interface{}{"Context": c})
	return nil
}

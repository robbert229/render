package render

import (
	"html/template"
	"net/http"
	"path/filepath"

	glog "github.com/go-kit/kit/log"
)

var (
	// LogError is used in the structured logging to denote that this is an error
	LogError = "error"
	// LogTemplate denotes a template
	LogTemplate = "template"
	LogLoading  = "loading"
	LogPath     = "path"
	LogName     = "name"
	LogDetails  = "details"
	LogLocation = "where"
	LogCausedBy = "who"
)

var (
	templates                = map[string]*template.Template{}
	sharedTemplatesDirectory = "./templates/shared/*.tmpl"
	logger                   glog.Logger
)

// SetLogger sets the logger that the render library should use to log.
func SetLogger(newLogger glog.Logger) {
	logger = newLogger
}

func getFiles(path string) []string {
	files, err := filepath.Glob(path)
	if err != nil {
		panic(err)
	}

	return files
}

func getTemplate(tmplPath string) *template.Template {
	var err error

	logger.Log(LogLoading, LogTemplate, LogPath, tmplPath)
	tmpl := template.New(tmplPath)
	tmpl, err = tmpl.ParseFiles(append(getFiles(sharedTemplatesDirectory), tmplPath)...)
	if err != nil {
		logger.Log(LogError, "unable to parse template", LogPath, tmplPath, LogLocation, "getTemplate", LogDetails, err.Error())
	}

	return tmpl
}

//View is called to render a template with the given path, and model
func View(w http.ResponseWriter, r *http.Request, tmplPath string, model map[interface{}]interface{}) {
	tmpl := getTemplate(tmplPath)
	if err := tmpl.ExecuteTemplate(w, "layout", model); err != nil {
		logger.Log(LogError, "unable to execute template", LogPath, tmplPath, LogLocation, "RenderView", LogDetails, err.Error())
	}
}

//Error is called to render a template called error in the templates directory.
func Error(w http.ResponseWriter, r *http.Request, err error, source string) {
	logger.Log(LogError, "rendering error", LogPath, source, LogDetails, err.Error())
	View(w, r, "templates\\error.tmpl", map[interface{}]interface{}{
		"error": err.Error(),
	})
}

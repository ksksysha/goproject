package handler

import (
	"html/template"
	"mygoproject/internal/model"
	templatefuncs "mygoproject/internal/template"
	"net/http"
	"path/filepath"
	"runtime"
)

var templatesDir string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	templatesDir = filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(currentFile))), "templates")
}

func getTemplatePath(name string) string {
	return filepath.Join(templatesDir, name)
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data *model.PageData, useLayout bool) {
	var tmplContent *template.Template
	var err error

	// Создаем новый шаблон с функциями-помощниками
	tmplContent = template.New("").Funcs(templatefuncs.FuncMap)

	if useLayout {
		layoutPath := getTemplatePath("layout.html")
		tmplPath := getTemplatePath(tmpl)
		tmplContent, err = tmplContent.ParseFiles(layoutPath, tmplPath)
	} else {
		tmplPath := getTemplatePath(tmpl)
		tmplContent, err = tmplContent.ParseFiles(tmplPath)
	}

	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if useLayout {
		err = tmplContent.ExecuteTemplate(w, "layout", data)
	} else {
		err = tmplContent.Execute(w, data)
	}

	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}

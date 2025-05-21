package http

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"myproject/internal/domain"

	"github.com/gorilla/sessions"
)

type Handler struct {
	authUC    *usecase.AuthUseCase
	bookingUC *usecase.BookingUseCase
	adminUC   *usecase.AdminUseCase
	store     *sessions.CookieStore
}

func NewHandler(repo domain.Repository, store *sessions.CookieStore) *Handler {
	return &Handler{
		authUC:    usecase.NewAuthUseCase(repo),
		bookingUC: usecase.NewBookingUseCase(repo),
		adminUC:   usecase.NewAdminUseCase(repo),
		store:     store,
	}
}

func (h *Handler) renderTemplate(w http.ResponseWriter, tmpl string, data *domain.PageData, useLayout bool) {
	var tmplContent *template.Template
	var err error

	if useLayout {
		layoutPath := filepath.Join("templates", "layout.html")
		tmplPath := filepath.Join("templates", tmpl)
		tmplContent, err = template.ParseFiles(layoutPath, tmplPath)
	} else {
		tmplPath := filepath.Join("templates", tmpl)
		tmplContent, err = template.ParseFiles(tmplPath)
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

func (h *Handler) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path
	if page == "/" {
		page = "/home"
	}

	pageFile := strings.TrimPrefix(page, "/") + ".html"
	fullPath := filepath.Join("templates", pageFile)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		h.NotFoundHandler(w, r)
		return
	}

	data := &domain.PageData{
		Title:   strings.Title(strings.TrimSuffix(filepath.Base(pageFile), ".html")),
		Content: template.HTML(content),
	}

	h.renderTemplate(w, pageFile, data, true)
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := &domain.PageData{
		Title:        "404 - Страница не найдена",
		ErrorMessage: "Запрашиваемая страница не существует.",
	}
	h.renderTemplate(w, "404.html", data, true)
}

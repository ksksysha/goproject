package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/sessions"
)

type Handler struct {
	store       *sessions.CookieStore
	templates   map[string]*template.Template
	projectRoot string
}

type PageData struct {
	Title   string
	Content template.HTML
}

func NewHandler(store *sessions.CookieStore) (*Handler, error) {
	// Более надежный способ определения корня проекта
	projectRoot, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Определяем пути к шаблонам
	templatesDir := filepath.Join(projectRoot, "templates")

	// Загружаем базовый шаблон
	layoutTmpl, err := template.ParseFiles(filepath.Join(templatesDir, "layout.html"))
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки layout.html: %v", err)
	}

	// Предварительно загружаем все шаблоны
	templates := make(map[string]*template.Template)

	// Основные шаблоны
	mainTemplates := []string{
		"index.html", "login.html", "register.html", "profile.html",
		"admin.html", "404.html", "about.html", "contacts.html",
		"services.html", "home.html", "book.html",
	}

	for _, tmpl := range mainTemplates {
		tmplPath := filepath.Join(templatesDir, tmpl)
		if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
			log.Printf("Предупреждение: шаблон %s не найден", tmplPath)
			continue
		}

		// Создаем новый шаблон на основе layout
		parsedTmpl, err := layoutTmpl.Clone()
		if err != nil {
			return nil, fmt.Errorf("ошибка клонирования layout для %s: %v", tmpl, err)
		}

		// Добавляем контент шаблона
		_, err = parsedTmpl.ParseFiles(tmplPath)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга %s: %v", tmpl, err)
		}

		templates[tmpl] = parsedTmpl
	}

	return &Handler{
		store:       store,
		templates:   templates,
		projectRoot: projectRoot,
	}, nil
}

func (h *Handler) renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := h.templates[name]
	if !ok {
		http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
		return
	}

	// Если data не передан, создаем базовую структуру
	if data == nil {
		data = PageData{
			Title:   "Главная страница",
			Content: template.HTML(""),
		}
	}

	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Ошибка рендеринга шаблона %s: %v", name, err)
		http.Error(w, "Ошибка отображения страницы", http.StatusInternalServerError)
	}
}

func (h *Handler) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:   "Главная страница",
		Content: template.HTML("{{template \"content\" .}}"),
	}
	h.renderTemplate(w, "index.html", data)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "login.html", nil)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "register.html", nil)
}

func (h *Handler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "profile.html", nil)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) BookServiceHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "book.html", nil)
}

func (h *Handler) AdminHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "admin.html", nil)
}

func (h *Handler) DeleteBookingHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (h *Handler) DeleteUserBookingHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *Handler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.renderTemplate(w, "404.html", nil)
}

func (h *Handler) AboutHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "about.html", nil)
}

func (h *Handler) ContactsHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "contacts.html", nil)
}

func (h *Handler) ServicesHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "services.html", nil)
}

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	h.renderTemplate(w, "home.html", nil)
}

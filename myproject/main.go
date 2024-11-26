package main

import (
	"html/template"
	"log"
	"net/http"
)

// Инициализация шаблонов
var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}
}

// Главная страница
func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html")
}

// О нас
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "about.html")
}

// Контакты
func contactsHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contacts.html")
}

// Услуги
func servicesHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "services.html")
}

// Функция для рендеринга шаблона
func renderTemplate(w http.ResponseWriter, page string) {
	err := templates.ExecuteTemplate(w, "layout.html", map[string]interface{}{
		"Content": page,
	})
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
	}
}

// Обработка статических файлов
func staticFileHandler() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
}

func main() {
	// Настройка маршрутов
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contacts", contactsHandler)
	http.HandleFunc("/services", servicesHandler)
	http.Handle("/static/", staticFileHandler())

	// Порт
	port := ":8080"
	log.Printf("Сервер запущен на http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Структура для передачи данных в шаблон
type PageData struct {
	Title   string
	Content template.HTML
}

// Функция для рендеринга шаблонов
func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	// Указываем путь к шаблонам
	layoutPath := filepath.Join("templates", "layout.html")
	tmplPath := filepath.Join("templates", tmpl)

	// Парсим оба шаблона: layout и нужную страницу
	tmplContent, err := template.ParseFiles(layoutPath, tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Рендерим основной шаблон (layout.html) с данными
	err = tmplContent.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Обработчик главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := &PageData{
		Title:   "Главная страница",
		Content: "<h1>Добро пожаловать в наш салон красоты!</h1><p>Мы предлагаем широкий спектр услуг для вашего ухода.</p>",
	}
	renderTemplate(w, "home.html", data)
}

// Обработчик страницы "О нас"
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := &PageData{
		Title:   "О нас",
		Content: "<h1>Мы - команда профессионалов!</h1><p>Наш салон предлагает только качественные и безопасные процедуры.</p>",
	}
	renderTemplate(w, "about.html", data)
}

// Обработчик страницы "Контакты"
func contactsHandler(w http.ResponseWriter, r *http.Request) {
	data := &PageData{
		Title:   "Контакты",
		Content: "<h1>Наши контакты</h1><p>Адрес: ул. Примерная 123, Телефон: +7 (123) 456-78-90</p>",
	}
	renderTemplate(w, "contacts.html", data)
}

// Обработчик страницы "Услуги"
func servicesHandler(w http.ResponseWriter, r *http.Request) {
	data := &PageData{
		Title:   "Услуги",
		Content: "<h1>Наши услуги</h1><p>Маникюр, педикюр, стрижки, укладки и многое другое.</p>",
	}
	renderTemplate(w, "services.html", data)
}

// Обработчик страницы 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	data := &PageData{
		Title:   "404 - Страница не найдена",
		Content: "<h1>Извините, страница не найдена!</h1>",
	}
	renderTemplate(w, "404.html", data)
}

func main() {
	// Регистрация маршрутов
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contacts", contactsHandler)
	http.HandleFunc("/services", servicesHandler)

	// Обработка 404 страницы
	http.HandleFunc("/404", notFoundHandler)

	// Статические файлы (например, CSS)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	// Запуск сервера на localhost:8080
	log.Println("Запуск сервера на http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

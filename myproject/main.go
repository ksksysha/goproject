package main

import (
	"html/template"
	"io/ioutil" // Исправлено здесь, убрали "so"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Структура для передачи данных в шаблон
type PageData struct {
	Title   string
	Content template.HTML
}

// Функция для рендеринга шаблонов
func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	log.Println("Rendering template:", tmpl)

	// Указываем путь к шаблонам
	layoutPath := filepath.Join("templates", "layout.html")
	tmplPath := filepath.Join("templates", tmpl)

	// Парсим оба шаблона: layout и нужную страницу
	tmplContent, err := template.ParseFiles(layoutPath, tmplPath, "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Рендерим основной шаблон с данными
	err = tmplContent.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Универсальный обработчик для всех страниц
func pageHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path

	// Если путь пустой, перенаправляем на главную страницу
	if page == "/" {
		page = "/home"
	}

	// Убираем начальный слэш и добавляем .html к имени файла
	pageFile := strings.TrimPrefix(page, "/") + ".html"
	fullPath := filepath.Join("templates", pageFile)

	// Проверка наличия HTML-файла
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		notFoundHandler(w, r) // Если файл не найден
		return
	}

	// Читаем содержимое HTML-файла
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "Ошибка при чтении страницы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Формируем данные для шаблона
	data := &PageData{
		Title:   strings.Title(strings.TrimSuffix(filepath.Base(pageFile), ".html")),
		Content: template.HTML(content), // Вставляем содержимое файла
	}

	renderTemplate(w, pageFile, data)
}

// Обработчик для 404 ошибки
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	data := &PageData{
		Title:   "404 - Страница не найдена",
		Content: template.HTML("<h1>Извините, страница не найдена!</h1>"),
	}
	w.WriteHeader(http.StatusNotFound)
	renderTemplate(w, "404.html", data)
}

// Отображение страницы входа
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Login handler invoked, method:", r.Method)

	data := &PageData{
		Title:   "Вход",
		Content: template.HTML(""), // В контенте ничего нет, всё на странице login.html
	}
	renderTemplate(w, "login.html", data) // Обратите внимание, мы вызываем 'renderTemplate' с 'login.html'
}

func main() {
	// Статические файлы (например, CSS)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Обработчики
	http.HandleFunc("/", pageHandler)
	http.HandleFunc("/login", loginPageHandler)

	// Запуск сервера на localhost:8080
	log.Println("Запуск сервера на http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

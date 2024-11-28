package main

import (
	"html/template"
	"io/ioutil"
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
	// Указываем путь к шаблонам
	layoutPath := filepath.Join("templates", "layout.html")
	tmplPath := filepath.Join("templates", tmpl)

	// Парсим оба шаблона: layout и нужную страницу
	tmplContent, err := template.ParseFiles(layoutPath, tmplPath, "templates/header.html", "templates/footer.html")
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

// Универсальный обработчик для всех страниц
func pageHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем имя страницы из URL
	page := r.URL.Path

	// Если путь пустой, перенаправляем на главную страницу
	if page == "/" {
		page = "/home" // По умолчанию главная страница
	}

	// Убираем начальный слэш и добавляем .html к имени файла
	pageFile := strings.TrimPrefix(page, "/") + ".html"
	fullPath := filepath.Join("templates", pageFile)

	// Проверка наличия HTML-файла
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Если файл не найден, перенаправляем на страницу 404
		notFoundHandler(w, r)
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

	// Рендерим шаблон
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

func main() {
	// Статические файлы (например, CSS)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Универсальный обработчик для всех страниц
	http.HandleFunc("/", pageHandler)

	// Запуск сервера на localhost:8080
	log.Println("Запуск сервера на http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

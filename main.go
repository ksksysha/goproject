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
	Title        string
	Content      template.HTML
	ErrorMessage string // Добавляем поле для сообщения об ошибке
}

// Здесь вы можете изменить имя пользователя и пароль
var (
	validUsername = "user"
	validPassword = "password"
)

// Функция для рендеринга шаблонов
func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	layoutPath := filepath.Join("templates", "layout.html")
	tmplPath := filepath.Join("templates", tmpl)

	// Парсим шаблоны
	tmplContent, err := template.ParseFiles(layoutPath, tmplPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmplContent.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Универсальный обработчик для всех страниц
func pageHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path

	if page == "/" {
		page = "/home"
	}

	pageFile := strings.TrimPrefix(page, "/") + ".html"
	fullPath := filepath.Join("templates", pageFile)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		notFoundHandler(w, r)
		return
	}

	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "Ошибка при чтении страницы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := &PageData{
		Title:   strings.Title(strings.TrimSuffix(filepath.Base(pageFile), ".html")),
		Content: template.HTML(content),
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

// Обработчик для страницы входа
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var data PageData
	data.Title = "Вход"

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == validUsername && password == validPassword {
			// Успешный вход
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		} else {
			// Неверный логин или пароль
			data.ErrorMessage = "Неверный логин или пароль"
		}
	}

	// Отображаем страницу входа с формой
	data.Content = template.HTML(`
<h2>Вход</h2>
<form action="/login" method="POST">
<label for="username">Имя пользователя:</label>
<input type="text" id="username" name="username" required>
<label for="password">Пароль:</label>
<input type="password" id="password" name="password" required>
<input type="submit" value="Войти">
</form>
<p>Нет аккаунта? <a href="/register">Зарегистрируйтесь</a></p>
`)

	renderTemplate(w, "login.html", &data)
}

// Обработчик для страницы регистрации
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Здесь можно обработать данные формы (например, сохранить в БД).
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Отображаем страницу регистрации
	data := &PageData{
		Title: "Регистрация",
		Content: template.HTML(`<h2>Регистрация</h2>
<form action="/register" method="POST">
<label for="username">Имя пользователя:</label>
<input type="text" id="username" name="username" required>
<label for="password">Пароль:</label>
<input type="password" id="password" name="password" required>
<input type="submit" value="Зарегистрироваться">
</form>
<p>Уже есть аккаунт? <a href="/login">Войдите</a></p>`),
	}
	renderTemplate(w, "register.html", data)
}

// Обработчик для главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := &PageData{
		Title:   "Главная",
		Content: template.HTML("<h1>Добро пожаловать на главную страницу!</h1>"),
	}
	renderTemplate(w, "home.html", data)

}

func main() {
	// Обработчики для статических файлов
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/login", loginHandler)       // Обработчик входа
	http.HandleFunc("/register", registerHandler) // Обработчик для регистрации
	http.HandleFunc("/home", homeHandler)         // Обработчик для главной страницы
	http.HandleFunc("/", pageHandler)             // Обработчик для других страниц

	log.Println("Запуск сервера на http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

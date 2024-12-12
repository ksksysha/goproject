package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Данные для подключения к базе данных PostgreSQL
	db *sql.DB
	// Данные для логина
	validUsername = "user"
	validPassword = "password"
)

type PageData struct {
	Title        string
	Content      template.HTML
	ErrorMessage string
	Username     string
}

func init() {
	var err error
	// Замените на ваши данные
	connStr := "host=localhost port=5432 user=postgres password=0000 dbname=myproject sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Проверка соединения
	if err = db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
	log.Println("Подключение к базе данных успешно!")
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	layoutPath := filepath.Join("templates", "layout.html")
	tmplPath := filepath.Join("templates", tmpl)

	// Загружаем layout и текущий шаблон
	content, err := os.ReadFile(tmplPath) // Читаем содержимое текущего шаблона
	if err != nil {
		http.Error(w, "Ошибка чтения шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data.Content = template.HTML(content) // Вставляем HTML-контент в структуру PageData

	tmplContent, err := template.ParseFiles(layoutPath)
	if err != nil {
		http.Error(w, "Ошибка загрузки layout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем рендеринг layout
	err = tmplContent.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Вход"}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Получаем хэш пароля из БД
		var dbPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				data.ErrorMessage = "Неверный логин или пароль"
			} else {
				data.ErrorMessage = "Ошибка при подключении к базе данных"
			}
		} else {
			// Сравниваем хэш с паролем
			err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
			if err == nil {
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			} else {
				data.ErrorMessage = "Неверный логин или пароль"
			}
		}
	}

	renderTemplate(w, "login.html", &data)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Регистрация"}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Хэшируем пароль
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			data.ErrorMessage = "Ошибка при хэшировании пароля"
			renderTemplate(w, "register.html", &data)
			return
		}

		// Сохраняем в базе
		_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
		if err != nil {
			data.ErrorMessage = "Ошибка при регистрации: пользователь уже существует"
		} else {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	renderTemplate(w, "register.html", &data)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:    "Личный кабинет",
		Username: "Ваше имя пользователя",
	}

	renderTemplate(w, "profile.html", &data)
}

func main() {
	// Обработчик для статических файлов
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Обработчики страниц
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Path
		if page == "/" {
			page = "/home"
		}

		pageFile := strings.TrimPrefix(page, "/") + ".html"
		fullPath := filepath.Join("templates", pageFile)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			http.Error(w, "Ошибка при чтении страницы: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := &PageData{
			Title:   strings.Title(strings.TrimSuffix(filepath.Base(pageFile), ".html")),
			Content: template.HTML(content),
		}

		renderTemplate(w, pageFile, data)
	})

	log.Println("Запуск сервера на http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

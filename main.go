package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	store = sessions.NewCookieStore([]byte("your-secret-key"))
	db    *sql.DB
)

type PageData struct {
	Title        string
	Content      template.HTML
	ErrorMessage string
	Username     string
	Services     []Service
	Bookings     []Booking
}

type Service struct {
	ID    int
	Name  string
	Price float64
}

type Booking struct {
	ServiceName string
	BookingTime string
}

func init() {
	var err error
	connStr := "host=localhost port=5432 user=postgres password=0000 dbname=myproject sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
	log.Println("Подключение к базе данных успешно!")

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData, useLayout bool) {
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Вход"}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				data.ErrorMessage = "Неверный логин или пароль"
			} else {
				data.ErrorMessage = "Ошибка при подключении к базе данных"
			}
		} else {
			err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
			if err == nil {
				session, _ := store.Get(r, "session-name")
				session.Values["username"] = username
				session.Save(r, w)

				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			} else {
				data.ErrorMessage = "Неверный логин или пароль"
			}
		}
	}

	// Встроенный HTML-код для страницы входа
	data.Content = template.HTML(`
	<div class="auth-container">
		<div class="auth-box">
			<h2>Вход</h2>
			<form method="POST" action="/login">
				<div class="input-group">
					<span class="input-icon"><i class="fas fa-user"></i></span>
					<input type="text" name="username" placeholder="Имя пользователя" required>
				</div>
				<div class="input-group">
					<span class="input-icon"><i class="fas fa-lock"></i></span>
					<input type="password" name="password" placeholder="Пароль" required>
				</div>
				<button type="submit" class="auth-button">Войти</button>
				<p>Нет аккаунта? <a href="/register">Зарегистрируйтесь</a></p>
			</form>
		</div>
	</div>
	`)
	renderTemplate(w, "login.html", &data, true)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Регистрация"}

	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			data.ErrorMessage = "Ошибка при хэшировании пароля"
			renderTemplate(w, "register.html", &data, false)
			return
		}

		_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
		if err != nil {
			data.ErrorMessage = "Ошибка при регистрации: пользователь уже существует"
		} else {
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
	}

	// Встроенный HTML-код для страницы регистрации
	data.Content = template.HTML(`
	<div class="auth-container">
		<div class="auth-box">
			<h2>Регистрация</h2>
			<form method="POST" action="/register">
				<div class="input-group">
					<span class="input-icon"><i class="fas fa-user"></i></span>
					<input type="text" name="username" placeholder="Имя пользователя" required>
				</div>
				<div class="input-group">
					<span class="input-icon"><i class="fas fa-lock"></i></span>
					<input type="password" name="password" placeholder="Пароль" required>
				</div>
				<button type="submit" class="auth-button">Зарегистрироваться</button>
				<p>Уже есть аккаунт? <a href="/login">Войдите</a></p>
			</form>
		</div>
	</div>
	`)
	renderTemplate(w, "register.html", &data, true)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Получаем список услуг
	rows, err := db.Query("SELECT id, name, price FROM services")
	if err != nil {
		http.Error(w, "Ошибка при получении услуг", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var service Service
		if err := rows.Scan(&service.ID, &service.Name, &service.Price); err != nil {
			http.Error(w, "Ошибка при получении данных услуг", http.StatusInternalServerError)
			return
		}
		services = append(services, service)
	}

	// Получаем список записей пользователя
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка при получении ID пользователя", http.StatusInternalServerError)
		return
	}

	bookingRows, err := db.Query(`
		SELECT s.name, b.booking_time
		FROM bookings b
		JOIN services s ON b.service_id = s.id
		WHERE b.user_id = $1
		ORDER BY b.booking_time DESC
	`, userID)
	if err != nil {
		http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
		return
	}
	defer bookingRows.Close()

	var bookings []Booking
	for bookingRows.Next() {
		var booking Booking
		var bookingTime time.Time
		if err := bookingRows.Scan(&booking.ServiceName, &bookingTime); err != nil {
			http.Error(w, "Ошибка при обработке записей", http.StatusInternalServerError)
			return
		}
		// Форматируем время записи для отображения
		booking.BookingTime = bookingTime.Format("2006-01-02 15:04")
		bookings = append(bookings, booking)
	}

	// Отправляем данные в шаблон
	data := PageData{
		Title:    "Личный кабинет",
		Username: username,
		Services: services,
		Bookings: bookings,
	}

	renderTemplate(w, "profile.html", &data, false)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	delete(session.Values, "username")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func bookServiceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	serviceID := r.FormValue("service_id")
	bookingTime := r.FormValue("booking_time")

	// Преобразуем serviceID в число
	serviceIDInt, err := strconv.Atoi(serviceID)
	if err != nil {
		http.Error(w, "Неверный формат ID услуги", http.StatusBadRequest)
		return
	}

	// Преобразуем bookingTime в формат time.Time
	bookingTimeParsed, err := time.Parse("2006-01-02T15:04", bookingTime)
	if err != nil {
		http.Error(w, "Неверный формат времени бронирования", http.StatusBadRequest)
		return
	}

	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&userID)
	if err != nil {
		http.Error(w, "Ошибка при получении данных пользователя", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO bookings (user_id, service_id, booking_time) VALUES ($1, $2, $3)", userID, serviceIDInt, bookingTimeParsed)
	if err != nil {
		http.Error(w, "Ошибка при записи на услугу", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/book", bookServiceHandler)
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

		renderTemplate(w, pageFile, data, true)
	})

	log.Println("Запуск сервера на http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

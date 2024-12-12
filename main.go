package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"        // Подключение к PostgreSQL
	"golang.org/x/crypto/bcrypt" // Для хэширования паролей
)

var (
	store = sessions.NewCookieStore([]byte("your-secret-key")) // Создаем хранилище сессий
	// Данные для подключения к базе данных PostgreSQL
	db *sql.DB
)

// Структуры данных для профиля
type PageData struct {
	Title        string        // Заголовок страницы
	Content      template.HTML // Контент страницы (HTML)
	ErrorMessage string        // Сообщение об ошибке
	Username     string        // Имя пользователя (если авторизован)
	Services     []Service     // Список доступных услуг
	Bookings     []Booking     // Список записей пользователя на услуги
}

type Service struct {
	ID    int     // ID услуги
	Name  string  // Название услуги
	Price float64 // Цена услуги
}

type Booking struct {
	ServiceName string // Название услуги
	BookingTime string // Время записи
}

func init() {
	var err error
	// Замените на ваши данные
	connStr := "host=localhost port=5432 user=postgres password=0000 dbname=myproject sslmode=disable"
	db, err = sql.Open("postgres", connStr) // Открытие соединения с базой данных
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	// Проверка соединения
	if err = db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}
	log.Println("Подключение к базе данных успешно!")

	// Настройки для cookies (сессий)
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,  // Время жизни сессии в секундах
		HttpOnly: true,  // Делает cookie доступным только через HTTP
		Secure:   false, // Для безопасных соединений (https) используйте true
	}
}

// Функция рендеринга шаблонов
func renderTemplate(w http.ResponseWriter, tmpl string, data *PageData) {
	layoutPath := filepath.Join("templates", "layout.html") // Путь к основному шаблону
	tmplPath := filepath.Join("templates", tmpl)            // Путь к конкретному шаблону

	// Загружаем layout и текущий шаблон
	content, err := os.ReadFile(tmplPath)
	if err != nil {
		http.Error(w, "Ошибка чтения шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data.Content = template.HTML(content) // Вставляем HTML-контент в структуру PageData

	// Загружаем layout (основной шаблон)
	tmplContent, err := template.ParseFiles(layoutPath)
	if err != nil {
		http.Error(w, "Ошибка загрузки layout: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Выполняем рендеринг layout с переданными данными
	err = tmplContent.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}

// Обработчик для страницы входа
func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Вход"}

	if r.Method == http.MethodPost {
		username := r.FormValue("username") // Получаем имя пользователя из формы
		password := r.FormValue("password") // Получаем пароль из формы

		log.Printf("Попытка входа: %s", username) // Логируем попытку входа

		var dbPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&dbPassword)
		if err != nil {
			if err == sql.ErrNoRows {
				data.ErrorMessage = "Неверный логин или пароль" // Ошибка: пользователь не найден
			} else {
				data.ErrorMessage = "Ошибка при подключении к базе данных"
			}
		} else {
			err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)) // Проверяем пароль
			if err == nil {
				// Сохранение имени пользователя в сессии
				session, _ := store.Get(r, "session-name")
				session.Values["username"] = username
				err = session.Save(r, w) // Обязательно сохраняем сессию
				if err != nil {
					log.Println("Ошибка сохранения сессии:", err)
				}

				// Переход на страницу личного кабинета
				http.Redirect(w, r, "/profile", http.StatusSeeOther)
				return
			} else {
				data.ErrorMessage = "Неверный логин или пароль"
			}
		}
	}

	renderTemplate(w, "login.html", &data) // Отображаем страницу входа
}

// Обработчик для страницы регистрации
func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{Title: "Регистрация"}

	if r.Method == http.MethodPost {
		username := r.FormValue("username") // Получаем имя пользователя
		password := r.FormValue("password") // Получаем пароль

		log.Printf("Попытка регистрации: %s", username)

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
			// Успешная регистрация: перенаправляем в личный кабинет
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
	}

	renderTemplate(w, "register.html", &data) // Отображаем страницу регистрации
}

// Обработчик для личного кабинета
func profileHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных из сессии и имени пользователя
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Если пользователь не авторизован, перенаправляем на страницу входа
		return
	}

	log.Println("Пользователь:", username)

	// Получаем список услуг из базы данных
	rows, err := db.Query("SELECT id, name, price FROM services")
	if err != nil {
		log.Println("Ошибка при выполнении запроса для получения услуг:", err)
		http.Error(w, "Ошибка при получении услуг", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var service Service
		if err := rows.Scan(&service.ID, &service.Name, &service.Price); err != nil {
			log.Println("Ошибка при сканировании данных услуги:", err)
			http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
			return
		}
		services = append(services, service)
	}

	// Получаем записи пользователя на услуги
	rows, err = db.Query(`
        SELECT s.name, b.booking_time
        FROM bookings b
        JOIN services s ON b.service_id = s.id
        JOIN users u ON b.user_id = u.id
        WHERE u.username = $1`, username)
	if err != nil {
		log.Println("Ошибка при выполнении запроса для получения записей:", err)
		http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		if err := rows.Scan(&booking.ServiceName, &booking.BookingTime); err != nil {
			log.Println("Ошибка при сканировании записей:", err)
			http.Error(w, "Ошибка при получении записей", http.StatusInternalServerError)
			return
		}
		bookings = append(bookings, booking)
	}

	// Передаем данные в шаблон
	data := PageData{
		Title:    "Личный кабинет",
		Username: username,
		Services: services,
		Bookings: bookings,
	}

	renderTemplate(w, "profile.html", &data) // Отображаем страницу личного кабинета
}

// Обработчик для выхода из системы
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	delete(session.Values, "username") // Удаляем имя пользователя из сессии
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther) // Перенаправляем на страницу входа
}

// Обработчик для записи на услугу
func bookServiceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session-name")
	username, ok := session.Values["username"].(string)
	if !ok || username == "" {
		log.Println("Ошибка: пользователь не авторизован")
		http.Redirect(w, r, "/login", http.StatusSeeOther) // Если пользователь не авторизован, перенаправляем на страницу входа
		return
	}

	// Получаем данные из формы
	serviceID := r.FormValue("service_id")
	bookingTime := r.FormValue("booking_time")

	log.Println("Получены данные формы для записи: service_id =", serviceID, "booking_time =", bookingTime)

	if serviceID == "" || bookingTime == "" {
		log.Println("Ошибка: пустые данные формы")
		http.Redirect(w, r, "/profile", http.StatusSeeOther) // Если данные пустые, перенаправляем в профиль
		return
	}

	// Получаем ID пользователя из базы данных
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username=$1", username).Scan(&userID)
	if err != nil {
		log.Println("Ошибка при получении ID пользователя:", err)
		http.Error(w, "Ошибка при получении данных пользователя", http.StatusInternalServerError)
		return
	}

	log.Println("ID пользователя:", userID)

	// Сохраняем запись
	_, err = db.Exec("INSERT INTO bookings (user_id, service_id, booking_time) VALUES ($1, $2, $3)", userID, serviceID, bookingTime)
	if err != nil {
		log.Println("Ошибка при записи на услугу:", err)
		http.Error(w, "Ошибка при записи на услугу", http.StatusInternalServerError)
		return
	}

	log.Println("Запись на услугу прошла успешно.")
	http.Redirect(w, r, "/profile", http.StatusSeeOther) // Перенаправляем в личный кабинет
}

// Основная функция для запуска сервера
func main() {
	// Обработчик для статических файлов (например, изображений и стилей)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Обработчики страниц
	http.HandleFunc("/login", loginHandler)       // Страница входа
	http.HandleFunc("/register", registerHandler) // Страница регистрации
	http.HandleFunc("/profile", profileHandler)   // Личный кабинет
	http.HandleFunc("/logout", logoutHandler)     // Выход
	http.HandleFunc("/book", bookServiceHandler)  // Запись на услугу
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Path
		if page == "/" {
			page = "/home"
		}

		pageFile := strings.TrimPrefix(page, "/") + ".html"
		fullPath := filepath.Join("templates", pageFile)

		content, err := os.ReadFile(fullPath) // Загружаем контент страницы
		if err != nil {
			http.Error(w, "Ошибка при чтении страницы: "+err.Error(), http.StatusInternalServerError)
			return
		}

		data := &PageData{
			Title:   strings.Title(strings.TrimSuffix(filepath.Base(pageFile), ".html")),
			Content: template.HTML(content),
		}

		renderTemplate(w, pageFile, data) // Отображаем страницу
	})

	log.Println("Запуск сервера на http://localhost:8080...")
	err := http.ListenAndServe(":8080", nil) // Запуск HTTP сервера на порту 8080
	if err != nil {
		log.Fatal("Ошибка запуска сервера: ", err)
	}
}

package session

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/sessions"
)

var (
	Store *sessions.CookieStore
	once  sync.Once
)

func Init() {
	once.Do(func() {
		// Используем фиксированный секретный ключ из переменной окружения
		secretKey := os.Getenv("SESSION_SECRET")
		if secretKey == "" {
			// Если переменная окружения не установлена, используем более надежный ключ
			secretKey = "myproject-session-secret-key-2024-development"
			log.Printf("ВНИМАНИЕ: Используется ключ разработки. Для продакшена установите SESSION_SECRET")
		}

		Store = sessions.NewCookieStore([]byte(secretKey))
		Store.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7, // 7 дней
			HttpOnly: true,
			Secure:   false, // установите в true для HTTPS
			SameSite: http.SameSiteLaxMode,
		}
		log.Printf("Инициализация хранилища сессий завершена")
	})
}

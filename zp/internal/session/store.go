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
			// Если переменная окружения не установлена, используем фиксированный ключ
			// В продакшене всегда должен быть установлен SESSION_SECRET
			secretKey = "default-secret-key-for-development-only"
			log.Printf("ВНИМАНИЕ: Используется фиксированный секретный ключ сессии. В продакшене установите SESSION_SECRET")
		}

		Store = sessions.NewCookieStore([]byte(secretKey))
		Store.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7, // 7 дней
			HttpOnly: true,
			Secure:   false, // установите в true, если используете HTTPS
			SameSite: http.SameSiteLaxMode,
		}
	})
}

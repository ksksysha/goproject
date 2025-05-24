package handler

import (
	"mygoproject/internal/model"
	"mygoproject/internal/session"
	"net/http"
	"path/filepath"
	"strings"
)

var protectedPages = map[string]bool{
	"/profile": true,
	"/book":    true,
	"/admin":   true,
}

func PageHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path
	if page == "/" {
		page = "/home"
	}

	// Проверяем авторизацию для защищенных страниц
	if protectedPages[page] {
		userID := session.GetUserID(r)
		if userID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	// Получаем имя файла шаблона
	pageFile := strings.TrimPrefix(page, "/")
	if strings.HasPrefix(pageFile, "services/description/") {
		parts := strings.Split(pageFile, "/")
		if len(parts) >= 3 {
			pageFile = filepath.Join("services", "description", parts[2]) + ".html"
		} else {
			NotFoundHandler(w, r)
			return
		}
	} else {
		parts := strings.Split(pageFile, "/")
		if len(parts) > 1 {
			pageFile = filepath.Join(parts...)
		}
		if !strings.HasSuffix(pageFile, ".html") {
			pageFile = pageFile + ".html"
		}
	}

	// Получаем имя страницы для заголовка
	pageName := strings.TrimSuffix(filepath.Base(pageFile), ".html")
	title := strings.Title(pageName)
	if pageName == "home" {
		title = "Главная"
	} else if strings.HasPrefix(page, "/services/") {
		parts := strings.Split(strings.TrimPrefix(page, "/services/"), "/")
		if len(parts) > 0 {
			switch parts[0] {
			case "nails":
				title = "Ногтевой сервис"
			case "hair":
				title = "Волосы"
			case "lashes":
				title = "Ресницы"
			case "brows":
				title = "Брови"
			}
		}
	}

	// Получаем данные сессии
	sess, _ := session.Store.Get(r, "session-name")
	username, _ := sess.Values["username"].(string)
	role, _ := sess.Values["role"].(string)

	// Создаем данные для шаблона
	data := &model.PageData{
		Title:    title + " - Салон красоты",
		Username: username,
		Role:     role,
		UserID:   session.GetUserID(r),
	}

	// Рендерим шаблон
	RenderTemplate(w, pageFile, data, true)
}

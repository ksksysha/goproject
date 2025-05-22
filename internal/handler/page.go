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
	pageFile := strings.TrimPrefix(page, "/") + ".html"
	if !strings.HasSuffix(pageFile, ".html") {
		pageFile = pageFile + ".html"
	}

	// Получаем имя страницы для заголовка
	pageName := strings.TrimSuffix(filepath.Base(pageFile), ".html")
	title := strings.Title(pageName)
	if pageName == "home" {
		title = "Главная"
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

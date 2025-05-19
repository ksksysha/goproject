package handler

import (
	"myproject/internal/model"
	"myproject/internal/session"
	"net/http"
	"os"
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

	// Проверяем, что запрашивается только HTML-страница
	if !strings.HasSuffix(page, ".html") {
		page = page + ".html"
	}

	pageFile := strings.TrimPrefix(page, "/")
	fullPath := getTemplatePath(pageFile)

	// Проверяем существование файла
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	// Читаем содержимое файла
	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	// Получаем имя страницы для заголовка
	pageName := strings.TrimSuffix(filepath.Base(pageFile), ".html")
	title := strings.Title(pageName)
	if pageName == "home" {
		title = "Главная"
	}

	// Создаем данные для шаблона
	data := &model.PageData{
		Title:   title + " - Салон красоты",
		Content: string(content),
		UserID:  session.GetUserID(r),
	}

	// Рендерим шаблон
	RenderTemplate(w, pageFile, data, true)
}

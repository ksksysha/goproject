package handler

import (
	"mygoproject/internal/model"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	data := &model.PageData{
		Title: "404 - Страница не найдена",
	}
	RenderTemplate(w, "404.html", data, true)
}

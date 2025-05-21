package handler

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

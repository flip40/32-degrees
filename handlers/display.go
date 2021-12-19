package handlers

import (
	"net/http"
)

func (h *Handler) ShowDisplay(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/display.html") // TODO: put file path in constants
}

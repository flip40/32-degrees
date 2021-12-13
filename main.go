package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/flip40/32-degrees/handlers"
)

func main() {
	// Setup Router and Handlers
	router := NewRouter()
	h := &handlers.Handler{} // TODO: add DB connection to handler

	router.AddAPI(http.MethodGet, "/add", h.AddDataHandler)
	router.AddAPI(http.MethodGet, "/get", h.GetDataHandler)

	// Start Server
	fmt.Printf("Started server on port %d\n", Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), router); err != nil {
		log.Fatal(err)
	}
}

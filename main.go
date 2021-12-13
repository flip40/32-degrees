package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Setup Handlers
	http.HandleFunc("/add", AddDataHandler)
	http.HandleFunc("/get", GetDataHandler)

	// Start Server
	fmt.Printf("Starter server on port %d\n", Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil); err != nil {
		log.Fatal(err)
	}
}

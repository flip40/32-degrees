package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/flip40/32-degrees/handlers"
	"github.com/flip40/32-degrees/mysql"
)

func main() {
	// Setup DB connection
	dsn := os.Getenv(MySQLURL)
	if dsn == "" {
		fmt.Println("using default DSN")
		dsn = DefaultDSN
	}
	db := mysql.NewConnection(dsn) // TODO: replace with docker env var
	defer db.Close()

	// Setup Router and Handlers
	router := NewRouter()
	h := &handlers.Handler{
		MySQL: db,
	}

	router.AddAPI(http.MethodGet, "/add", h.AddDataHandler)
	router.AddAPI(http.MethodGet, "/get", h.GetDataHandler)

	// Start Server
	fmt.Printf("Started server on port %d\n", Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), router); err != nil {
		log.Fatal(err)
	}
}

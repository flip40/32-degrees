package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	var database *mysql.DB
	for {
		db, err := mysql.NewConnection(dsn) // TODO: replace with docker env var
		if err != nil {
			fmt.Printf("failed to connect, %s\n", err)
			fmt.Println("trying again in 5 seconds")
			time.Sleep(5 * time.Second)
			continue
		}
		database = db
		defer database.Close()
		break
	}

	// Setup Router and Handlers
	router := NewRouter()
	h := &handlers.Handler{
		MySQL: database,
	}

	router.AddAPI(http.MethodGet, "/", h.ShowDisplay)
	router.AddAPI(http.MethodGet, "/add", h.AddDataHandler)
	router.AddAPI(http.MethodGet, "/get", h.GetDataHandler)
	router.AddAPI(http.MethodGet, "/plotdata", h.GetPlotDataHandler)

	// Start Server
	fmt.Printf("Started server on port %d\n", Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), router); err != nil {
		log.Fatal(err)
	}
}

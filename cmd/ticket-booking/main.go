package main

import (
	"database/sql"
	"log"

	// "os"

	"movie-ticket-booking/internal/api/routes"

	// "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	// gin.SetMode(gin.ReleaseMode)
	// os.Setenv("GIN_MODE", "release")
	// gin.SetMode(gin.ReleaseMode)
	// initial db
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(localhost:3306)/ticket_booking")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := routes.SetupRouter(db)
	// run server
	if err := router.Run(":8000"); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}

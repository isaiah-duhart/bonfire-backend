package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/isaiah-duhart/bonfire-backend/api/handlers"
	"github.com/isaiah-duhart/bonfire-backend/api/routes"
	"github.com/isaiah-duhart/bonfire-backend/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	h := handlers.Handler{}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		fmt.Println("Error creating db connection: ", err)
		return
	}

	h.Database = db
	h.Queries = database.New(db)
	h.Secret = os.Getenv("JWT_SECRET")

	if err = http.ListenAndServe(":8080", routes.GetRoutes(&h)); err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
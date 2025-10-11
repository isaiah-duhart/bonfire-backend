package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/isaiah-duhart/bonfire-backend/api/handlers"
	"github.com/isaiah-duhart/bonfire-backend/api/routes"
	"github.com/isaiah-duhart/bonfire-backend/internal/database"
	"github.com/isaiah-duhart/bonfire-backend/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()

	db_url := os.Getenv("DB_URL")
	jwt_secret := os.Getenv("JWT_SECRET")

	if utils.IsRunningOnEC2() {
		secrets := utils.GetAWSSecrets()
		db_url = secrets.DB_URL
		jwt_secret = secrets.JWT_SECRET
	}

	h := handlers.Handler{}

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatalf("Error creating db connection: %v", err)
	}

	h.Database = db
	h.Queries = database.New(db)
	h.Secret = jwt_secret

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(routes.GetRoutes(&h))

	if err = http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
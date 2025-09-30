package handlers

import (
	"database/sql"

	"github.com/isaiah-duhart/bonfire-backend/internal/database"
)

type Handler struct {
	Queries *database.Queries
	Database *sql.DB
}
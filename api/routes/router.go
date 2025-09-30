package routes

import (
	"net/http"

	"github.com/isaiah-duhart/bonfire-backend/api/handlers"
)

func GetRoutes(h *handlers.Handler) *http.ServeMux {
	serveMux := http.NewServeMux()

	// User Endpoints
	serveMux.HandleFunc("POST /api/users", h.CreateUser)
	serveMux.HandleFunc("DELETE /api/users", h.DeleteUser)

	// Group Endpoints
	serveMux.HandleFunc("GET /api/groups", h.GetGroups)
	serveMux.HandleFunc("POST /api/groups", h.CreateGroup)
	serveMux.HandleFunc("DELETE /api/groups", h.DeleteGroup)

	// Question Endpoints
	serveMux.HandleFunc("POST /api/questions", h.CreateQuestion)
	serveMux.HandleFunc("DELETE /api/questions", h.DeleteQuestion)
	
	// Group Question Endpoints
	serveMux.HandleFunc("POST /api/group-questions", h.GetGroupQuestions)
	serveMux.HandleFunc("DELETE /api/group-questions", h.DeleteGroupQuestions)

	// Group Response Endpoints
	serveMux.HandleFunc("GET /api/group-responses", h.GetGroupResponses)
	serveMux.HandleFunc("POST /api/group-responses", h.CreateGroupResponse)

	return serveMux
}
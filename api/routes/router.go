package routes

import (
	"net/http"

	"github.com/isaiah-duhart/bonfire-backend/api/handlers"
)

func GetRoutes(h *handlers.Handler) *http.ServeMux {
	serveMux := http.NewServeMux()

	// Login Endpoints
	serveMux.HandleFunc("POST /api/login", h.Login)

	// User Endpoints
	serveMux.HandleFunc("POST /api/users", h.CreateUser)
	serveMux.HandleFunc("DELETE /api/users", h.AuthMiddleware(h.DeleteUser))

	// Group Endpoints
	serveMux.HandleFunc("GET /api/groups", h.AuthMiddleware(h.GetGroups))
	serveMux.HandleFunc("POST /api/groups", h.AuthMiddleware(h.CreateGroup))
	serveMux.HandleFunc("DELETE /api/groups", h.DeleteGroup)

	// Question Endpoints
	serveMux.HandleFunc("POST /api/questions", h.CreateQuestion)
	serveMux.HandleFunc("DELETE /api/questions", h.DeleteQuestion)
	
	// Group Question Endpoints
	serveMux.HandleFunc("POST /api/group-questions", h.AuthMiddleware(h.GetGroupQuestions))
	serveMux.HandleFunc("DELETE /api/group-questions", h.DeleteGroupQuestions)

	// Group Response Endpoints
	serveMux.HandleFunc("GET /api/group-responses/{group_question_id}", h.AuthMiddleware(h.GetGroupResponses))
	serveMux.HandleFunc("POST /api/group-responses", h.AuthMiddleware(h.CreateGroupResponse))

	return serveMux
}
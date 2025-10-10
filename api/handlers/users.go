package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/isaiah-duhart/bonfire-backend/internal/database"
	"github.com/isaiah-duhart/bonfire-backend/utils"
)

type UserResponse struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	Name string `json:"name"`
	Birthday time.Time `json:"birthday"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		Name string `json:"name"`
		Birthday civil.Date `json:"birthday"`
	}

	params := parameters{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		fmt.Println("Error decoding params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	if params.Name == "" {
		fmt.Println("name is missing from request")
		utils.RespondWithError(w, 400, "name is required")
		return
	}

	hash, err := utils.HashPassword(params.Password)
	if err != nil {
		fmt.Println("Error hashing password: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	birthdaySqlTime := sql.NullTime{
		Time: params.Birthday.In(time.UTC),
		Valid: true,
	}
	if params.Birthday.IsZero() {
		birthdaySqlTime.Valid = false
	}

	user, err := h.Queries.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		Password: hash,
		Name: params.Name,
		Birthday: birthdaySqlTime,
	})
	if err != nil {
		fmt.Println("Error inserting user: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	token, err := utils.MakeJWT(user.ID, h.Secret, time.Hour * 24)
	if err != nil {
		fmt.Println("Error generating jwt: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	utils.RespondWithJson(w, 200, authResponse{
		Token: token,
		UserID: user.ID.String(),
	})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		fmt.Println("Couldn't get userID from context: ", r.Context().Value(userIDKey))
		utils.RespondWithError(w, 403, "invalid jwt")
		return
	}

	if err := h.Queries.DeleteUser(r.Context(), userID); err != nil {
		fmt.Println("Error deleting user: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	w.WriteHeader(204)
}
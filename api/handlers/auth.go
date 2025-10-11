package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/isaiah-duhart/bonfire-backend/utils"
)

type authResponse struct {
	Token string `json:"token"`
	UserID string `json:"user_id"`
}

type uuidKey string
const userIDKey uuidKey = "UserId"

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println("Error decoding params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	user, err := h.Queries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Println("Error getting user: ", err)
		utils.RespondWithError(w, 400, "invalid email or password")
		return
	}

	if err = utils.ComparePassword(user.Password, params.Password); err != nil {
		log.Println("Error commparing password: ", err)
		utils.RespondWithError(w, 400, "invalid email or password")
		return
	}

	token, err := utils.MakeJWT(user.ID, h.Secret, time.Hour * 24)
	if err != nil {
		log.Println("Error generating jwt: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	utils.RespondWithJson(w, 200, authResponse{
		Token: token,
		UserID: user.ID.String(),
	})
}

func (h *Handler) AuthMiddleware(next func (http.ResponseWriter, *http.Request)) (func (http.ResponseWriter, *http.Request)) {
	return func (w http.ResponseWriter, r *http.Request) {
		token, err := utils.GetBearerToken(r.Header)
		if err != nil {
			log.Println("Error getting token from auth header: ", err)
			utils.RespondWithError(w, 403, "missing auth header")
			return
		}

		userID, err := utils.ValidateJWT(token, h.Secret)
		if err != nil {
			log.Println("Error validating: ", err)
			utils.RespondWithError(w, 403, "invalid bearer token")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}
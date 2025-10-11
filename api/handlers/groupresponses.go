package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/isaiah-duhart/bonfire-backend/internal/database"
	"github.com/isaiah-duhart/bonfire-backend/utils"
)

type GroupResponse struct {
	ID uuid.UUID `json:"id"`
	GroupQuestionID uuid.UUID `json:"group_question_id"`
	Response string `json:"response"`
	Author string `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	AuthorID uuid.UUID `json:"author_id"`
}

func (h *Handler) CreateGroupResponse(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupQuestionId uuid.UUID `json:"group_question_id"`
		Response string `json:"response"`
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		log.Println("Couldn't get userID from context: ", r.Context().Value(userIDKey))
		utils.RespondWithError(w, 403, "invalid jwt")
		return
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println("Error decoding params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}


	exists, err := h.Queries.IsUserInGroupByGroupQuestionID(r.Context(), database.IsUserInGroupByGroupQuestionIDParams{
		ID: params.GroupQuestionId,
		UserID: userID,
	})
	if err != nil {
		log.Println("Error checking if user is in group: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}
	if !exists {
		log.Printf("Error user %v is not in group with questionID %v\n", userID, params.GroupQuestionId)
		utils.RespondWithError(w, 403, "user is not in group")
		return
	}

	groupResponse, err := h.Queries.CreateGroupResponse(r.Context(), database.CreateGroupResponseParams{
		GroupQuestionID: params.GroupQuestionId,
		Response: params.Response,
		AuthorID: userID,
	})
	if err != nil {
		log.Println("Error inserting group response: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	utils.RespondWithJson(w, 200, GroupResponse{
		ID: groupResponse.ID,
		GroupQuestionID: groupResponse.GroupQuestionID,
		Response: groupResponse.Response,
		Author: groupResponse.Name,
		CreatedAt: groupResponse.CreatedAt,
		AuthorID: groupResponse.ID_2,
	})
}

func (h *Handler) GetGroupResponses(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		log.Println("Couldn't get userID from context: ", r.Context().Value(userIDKey))
		utils.RespondWithError(w, 403, "invalid jwt")
		return
	}

	groupQuestionID, err := uuid.Parse(r.PathValue("group_question_id"))
	if err != nil {
		log.Println("Error parsing group_question_id: ", err)
		utils.RespondWithError(w, 400, "group_question_id is not a uuid")
		return
	}
	
	groupResponses, err := h.Queries.GetGroupResponses(r.Context(), database.GetGroupResponsesParams{
		GroupQuestionID: groupQuestionID,
		AuthorID: userID,
	})
	if err != nil {
		log.Println("Error getting group answers: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	resp := []GroupResponse{}
	for _, groupResponse := range groupResponses {
		resp = append(resp, GroupResponse{
			ID: groupResponse.ID,
			GroupQuestionID: groupResponse.GroupQuestionID,
			Response: groupResponse.Response,
			Author: groupResponse.Name,
			CreatedAt: groupResponse.CreatedAt,
			AuthorID: groupResponse.ID_2,
		})
	}

	utils.RespondWithJson(w, 200, resp)
}
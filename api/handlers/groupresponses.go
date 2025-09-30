package handlers

import (
	"encoding/json"
	"fmt"
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
	AuthorID uuid.UUID `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) CreateGroupResponse(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupQuestionId uuid.UUID `json:"group_question_id"`
		Response string `json:"response"`
		AuthorId uuid.UUID `json:"author_id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		fmt.Println("Error decoding params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	groupResponse, err := h.Queries.CreateGroupResponse(r.Context(), database.CreateGroupResponseParams{
		GroupQuestionID: params.GroupQuestionId,
		Response: params.Response,
		AuthorID: params.AuthorId,
	})
	if err != nil {
		fmt.Println("Error inserting group response: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	utils.RespondWithJson(w, 200, GroupResponse{
		ID: groupResponse.ID,
		GroupQuestionID: groupResponse.GroupQuestionID,
		Response: groupResponse.Response,
		AuthorID: groupResponse.AuthorID,
	})
}

func (h *Handler) GetGroupResponses(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if query.Get("group_question_id") == "" || query.Get("author_id") == "" {
		fmt.Println("Params missing from query: ", query)
		utils.RespondWithError(w, 400, "missing required params group_question_id and author_id")
		return
	}

	groupQuestionID, err := uuid.Parse(query.Get("group_question_id"))
	if err != nil {
		fmt.Println("Error parsing group_question_id: ", err)
		utils.RespondWithError(w, 400, "group_question_id is not a uuid")
		return
	}

	authorID, err := uuid.Parse(query.Get("author_id"))
	if err != nil {
		fmt.Println("Error parsing author_id: ", err)
		utils.RespondWithError(w, 400, "author_id is not a uuid")
		return
	}
	

	groupResponses, err := h.Queries.GetGroupResponses(r.Context(), database.GetGroupResponsesParams{
		GroupQuestionID: groupQuestionID,
		AuthorID: authorID,
	})
	if err != nil {
		fmt.Println("Error getting group answers: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	resp := []GroupResponse{}
	for _, groupResponse := range groupResponses {
		resp = append(resp, GroupResponse{
			ID: groupResponse.ID,
			GroupQuestionID: groupResponse.GroupQuestionID,
			Response: groupResponse.Response,
			AuthorID: groupResponse.AuthorID,
			CreatedAt: groupResponse.CreatedAt,
		})
	}

	utils.RespondWithJson(w, 200, resp)
}
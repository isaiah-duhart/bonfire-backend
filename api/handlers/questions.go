package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/isaiah-duhart/bonfire-backend/utils"
)

type QuestionResponse struct {
	ID uuid.UUID `json:"id"`
	Text string `json:"text"`
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Text string `json:"text"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		fmt.Println("Error decoding params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	question, err := h.Queries.CreateQuestion(r.Context(), params.Text)
	if err != nil {
		fmt.Println("Error creating question: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	utils.RespondWithJson(w, 200, QuestionResponse{
		ID: question.ID,
		Text: question.Text,
	})
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID uuid.UUID `json:"id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		fmt.Println("Error decoding params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	err := h.Queries.DeleteQuestion(r.Context(), params.ID)
	if err != nil {
		fmt.Println("Error deleting question: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	w.WriteHeader(204)
}
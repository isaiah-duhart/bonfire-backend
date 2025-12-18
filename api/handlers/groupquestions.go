package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/isaiah-duhart/bonfire-backend/internal/database"
	"github.com/isaiah-duhart/bonfire-backend/utils"
)

type GroupQuestionsResponse struct {
	ID uuid.UUID `json:"id"`
	GroupID uuid.UUID `json:"group_id"`
	Date civil.Date `json:"date"`
	QuestionText string`json:"question_text"`
	CreatedBy uuid.UUID `json:"created_by"`
}

const dailyLimit = 3

func (h *Handler) DeleteGroupQuestions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupID uuid.UUID `json:"group_id"`
		Date civil.Date `json:"date"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println("Error decoding params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	err := h.Queries.DeleteGroupQuestions(r.Context(), database.DeleteGroupQuestionsParams{
		GroupID: params.GroupID,
		Date: params.Date,
	})
	if err != nil {
		log.Println("Error deleting group_questions: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	w.WriteHeader(204)
}

func (h *Handler) GetGroupQuestions(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupID uuid.UUID `json:"group_id"`
		Date civil.Date `json:"date"`
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

	tx, err := h.Database.BeginTx(r.Context(), nil)
	if err != nil {
		log.Println("Error creating params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}
	defer tx.Rollback()

	qtx := h.Queries.WithTx(tx)

	exists, err := qtx.IsUserInGroup(r.Context(), database.IsUserInGroupParams{
		GroupID: params.GroupID,
		UserID: userID,
	})
	if err != nil {
		log.Println("Error checking if user is in group: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}
	if !exists {
		log.Printf("Error user %v is not in group %v\n", userID, params.GroupID)
		utils.RespondWithError(w, 403, "user is not in group")
		return
	}
	
	resp := []GroupQuestionsResponse{}

	numUsersGroupQuestions, err := qtx.CountGroupQuestions(r.Context(), database.CountGroupQuestionsParams{
		Date: params.Date,
		CreatedBy: userID,
		GroupID: params.GroupID,
	})
	if err != nil {
		log.Println("Error getting group questions: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	if numUsersGroupQuestions < dailyLimit {
		groupQuestions, err := qtx.CreateGroupQuestions(r.Context(), database.CreateGroupQuestionsParams{
			GroupID: params.GroupID,
			CreatedBy: userID,
			Date: params.Date,
			Limit: int32(dailyLimit - numUsersGroupQuestions),
		})
		if err != nil {
			log.Println("Error getting group questions: ", err)
			utils.RespondWithError(w, 500, "something went wrong")
			return
		}

		if len(groupQuestions) < dailyLimit {
			log.Println("Not enough questions to create group questions: ", groupQuestions)
			utils.RespondWithError(w, 500, "tell Isaiah to add more questions")
			return
		}
	}

	groupQuestions, err := qtx.GetGroupQuestions(
		r.Context(), 
		database.GetGroupQuestionsParams{
			GroupID: params.GroupID,
			CreatedBy: userID,
			Date: params.Date,
	})
	if err != nil {
		log.Println("Error getting group questions: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	for _, groupQuestion := range(groupQuestions) {
		resp = append(resp, GroupQuestionsResponse{
			ID: groupQuestion.ID,
			GroupID: groupQuestion.GroupID,
			Date: groupQuestion.Date,
			QuestionText: groupQuestion.Text,
			CreatedBy: groupQuestion.CreatedBy,
		})
	}
	utils.RespondWithJson(w, 200, resp)
}

func (h *Handler) GetAllGroupQuestions(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		log.Println("Couldn't get userID from context: ", r.Context().Value(userIDKey))
		utils.RespondWithError(w, 403, "invalid jwt")
		return
	}

	groupID, err := uuid.Parse(r.PathValue("group_id"))
	if err != nil {
		log.Println("Error getting group_id from path: ", err)
		utils.RespondWithError(w, 400, fmt.Sprint("Invalid group_id in path: ", r.PathValue("group_id")))
		return
	}

	tx, err := h.Database.BeginTx(r.Context(), nil)
	if err != nil {
		log.Panicln("Error creating database tx: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}
	defer tx.Rollback()

	qtx := h.Queries.WithTx(tx)

	exists, err := qtx.IsUserInGroup(r.Context(), database.IsUserInGroupParams{
		GroupID: groupID,
		UserID: userID,
	})
	if err != nil {
		log.Println("Error checking if user is in group: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	if !exists {
		log.Printf("Error user: %v is not in group %v\n", userID, groupID)
		utils.RespondWithError(w, 403, "invalid jwt")
		return
	}

	resp := []GroupQuestionsResponse{}
	groupQuestions, err := qtx.GetAllGroupQuestions(r.Context(), database.GetAllGroupQuestionsParams{
		GroupID: groupID,
		CreatedBy: userID,
	})
	if err != nil {
		log.Println("Error getting all group questions: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error committing transaction: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	for _, groupQuestion := range(groupQuestions) {
		resp = append(resp, GroupQuestionsResponse{
			ID: groupQuestion.ID,
			GroupID: groupQuestion.GroupID,
			Date: groupQuestion.Date,
			QuestionText: groupQuestion.Text,
			CreatedBy: groupQuestion.CreatedBy,
		})
	}

	utils.RespondWithJson(w, 200, resp)
}


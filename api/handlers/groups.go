package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/isaiah-duhart/bonfire-backend/internal/database"
	"github.com/isaiah-duhart/bonfire-backend/utils"
)

type Group struct {
	ID uuid.UUID `json:"id"`
	GroupID uuid.UUID `json:"group_id"`
	GroupName string `json:"group_name"`
	UserID uuid.UUID `json:"user_id"`
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		GroupName string `json:"group_name"`
		MemberEmails []string `json:"member_emails"`
	}

	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		log.Println("Invalid jwt: ", userID)
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
		log.Println("Error creating context: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}
	defer tx.Rollback()

	qtx := h.Queries.WithTx(tx)
	
	groups := []Group{}
	groupId := uuid.New()

	// Create row for group creator
	row, err := qtx.CreateGroupWithUserID(r.Context(), database.CreateGroupWithUserIDParams{
		GroupID: groupId,
		GroupName: params.GroupName,
		UserID: userID,
	})
	if err != nil {
		log.Println("Error creating group: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	groups = append(groups, Group{
		ID: row.ID,
		GroupID: row.GroupID,
		GroupName: row.GroupName,
		UserID: row.UserID,
	})

	// Create group requests
	for _, memberEmail := range params.MemberEmails {
		row, err := qtx.CreateGroup(r.Context(), database.CreateGroupParams{
			GroupID: groupId,
			GroupName: params.GroupName,
			Email: memberEmail,
		})
		if err != nil {
			log.Println("Error creating group: ", err)
			utils.RespondWithError(w, 500, "something went wrong")
			return
		}

		groups = append(groups, Group{
			ID: row.ID,
			GroupID: row.GroupID,
			GroupName: row.GroupName,
			UserID: row.UserID,
		})
	}

	if err := tx.Commit(); err != nil {
		log.Println("Error committing context: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	utils.RespondWithJson(w, 200, groups)
}

func (h *Handler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ID uuid.UUID `json:"id"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println("Error decoding params: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	if err := h.Queries.DeleteGroup(r.Context(), params.ID); err != nil {
		log.Println("Error deleting group: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	w.WriteHeader(204)
}

func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request){
	userID, ok := r.Context().Value(userIDKey).(uuid.UUID)
	if !ok {
		log.Println("Couldn't get userID from context: ", r.Context().Value(userIDKey))
		utils.RespondWithError(w, 403, "invalid jwt")
		return
	}

	groups, err := h.Queries.GetGroupsByUserID(r.Context(), userID)
	if err != nil {
		log.Println("Error getting groups: ", err)
		utils.RespondWithError(w, 500, "something went wrong")
		return
	}

	resp := []Group{}
	for _, group := range groups {
		resp = append(resp, Group{
			ID: group.ID,
			GroupID: group.GroupID,
			GroupName: group.GroupName,
			UserID: group.UserID,
		})
	}

	utils.RespondWithJson(w, 200, resp)
}
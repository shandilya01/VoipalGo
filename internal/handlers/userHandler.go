package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shandilya01/VoipalGo/internal/services"
)

type UserHandler struct {
	db *pgxpool.Pool
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var reqBody map[string]interface{}

	json.NewDecoder(r.Body).Decode(&reqBody)

	s := services.NewUserService(h.db)
	userObj, err := s.UserLogin(ctx, reqBody)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userObj)
	}

}

func (h *UserHandler) HandleUserSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var reqBody map[string]interface{}
	json.NewDecoder(r.Body).Decode(&reqBody)

	s := services.NewUserService(h.db)
	userObj, err := s.UserSignUp(ctx, reqBody)

	log.Print(userObj)
	log.Print(err)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userObj)
	}

}

func (h *UserHandler) HandleContacts(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	userId := ""
	if len(r.URL.Query()["id"]) > 0 {
		userId = r.URL.Query()["id"][0]
	}

	s := services.NewUserService(h.db)
	contactsArr, err := s.GetUserContactsById(ctx, userId)

	log.Print(contactsArr)
	log.Print(err)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contactsArr)
	}

}

func (h *UserHandler) HandlePushNotification(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	peerId := ""
	if len(r.URL.Query()["peerId"]) > 0 {
		peerId = r.URL.Query()["peerId"][0]
	}

	userId := ""
	if len(r.URL.Query()["userId"]) > 0 {
		userId = r.URL.Query()["userId"][0]
	}

	s := services.NewUserService(h.db)
	err := s.CallPushNotification(ctx, userId, peerId)

	log.Print("HandlePushNotification", err)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/shandilya01/VoipalGo/internal/services"
)

type UserHandler struct {
	db *pgx.Conn
}

func NewUserHandler(db *pgx.Conn) *UserHandler {
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

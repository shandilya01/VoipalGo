package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shandilya01/VoipalGo/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(db *pgxpool.Pool) *UserHandler {
	return &UserHandler{service: services.NewUserService(db)}
}

func (h *UserHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var reqBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	userObj, err := h.service.UserLogin(ctx, reqBody)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(err)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(userObj)
	}

}

func (h *UserHandler) HandleUserSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var reqBody map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	userObj, err := h.service.UserSignUp(ctx, reqBody)

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
	ctx := r.Context()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	userId := ""
	if len(r.URL.Query()["id"]) > 0 {
		userId = r.URL.Query()["id"][0]
	}

	contactsArr, err := h.service.GetUserContactsById(ctx, userId)

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
	ctx := r.Context()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
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

	roomId := ""
	if len(r.URL.Query()["roomId"]) > 0 {
		roomId = r.URL.Query()["roomId"][0]
	}

	err := h.service.CallPushNotification(ctx, userId, peerId, roomId)

	log.Print("HandlePushNotification", err)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		w.WriteHeader(http.StatusOK)
	}

}

func (h *UserHandler) HandleVoipId(w http.ResponseWriter, r *http.Request) {
	// ctx := r.Context()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// to be developed later
	w.WriteHeader(http.StatusUnprocessableEntity)
	///

	// userId := ""
	// if len(r.URL.Query()["userId"]) > 0 {
	// 	userId = r.URL.Query()["userId"][0]
	// }

	// inc, err := h.service.HandleVoipId(ctx, userId)

	// log.Print("HandleVoipId", err)

	// if err != nil {
	// 	w.WriteHeader(http.StatusUnprocessableEntity)
	// } else {
	// 	w.WriteHeader(http.StatusOK)
	// 	json.NewEncoder(w).Encode(inc)
	// }
}

func (h *UserHandler) HandleWordList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	wordList, err := h.service.HandleWordList(ctx)

	log.Print("HandleWordList error", err)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(wordList)
	}
}

func (h *UserHandler) HandleUserByVoipId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	voipId := ""
	if len(r.URL.Query()["voipId"]) > 0 {
		voipId = r.URL.Query()["voipId"][0]
	}

	user, err := h.service.HandleUserByVoipId(ctx, voipId)

	log.Print("HandleUserByVoipId error", err)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

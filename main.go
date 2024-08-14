package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shandilya01/VoipalGo/config"
	"github.com/shandilya01/VoipalGo/internal/handlers"
	"github.com/shandilya01/VoipalGo/pkg/db"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()
	db, err := db.NewPgxConn(ctx, cfg.DatabaseUrl)

	if err != nil {
		log.Fatalf("Could not connect to postgres %v \n", err)
	} else {
		log.Print("Connected to Postgres")
	}

	defer db.Close()

	server := mux.NewRouter()

	userHandler := handlers.NewUserHandler(db)
	signallingHandler := handlers.NewSignallingHandler()

	// creates multiple goroutines for each api call and each goroutines for each user
	server.HandleFunc("/login", userHandler.HandleUserLogin).Methods("POST", "OPTIONS")
	server.HandleFunc("/signup", userHandler.HandleUserSignUp).Methods("POST", "OPTIONS")
	server.HandleFunc("/contactsbyid", userHandler.HandleContacts).Methods("GET", "OPTIONS")
	server.HandleFunc("/refreshVoipId", userHandler.HandleVoipId).Methods("GET", "OPTIONS")
	server.HandleFunc("/pushCallNotification", userHandler.HandlePushNotification).Methods("GET", "POST", "OPTIONS")
	server.HandleFunc("/getWordList", userHandler.HandleWordList).Methods("GET", "OPTIONS")
	server.HandleFunc("/getUserByVoipId", userHandler.HandleUserByVoipId).Methods("GET", "OPTIONS")

	server.HandleFunc("/ws", signallingHandler.HandleNewSocket)

	log.Printf("Server listening on %s", cfg.ServerUrl)
	if err := http.ListenAndServe(cfg.ServerUrl, server); err != nil {
		log.Fatal(err)
	}
}

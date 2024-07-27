package services

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type SignallingService struct {
	Upgrader *websocket.Upgrader
	Clients  map[*websocket.Conn]string
	Rooms    map[string]map[*websocket.Conn]bool
	Mutex    sync.Mutex
}

// to be sent bw sockets
type Message struct {
	RoomId string
	Event  string
	Data   map[string]interface{}
}

func NewSignallingService() *SignallingService {
	return &SignallingService{
		Clients: make(map[*websocket.Conn]string),
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Rooms: make(map[string]map[*websocket.Conn]bool),
	}
}

func (s *SignallingService) HandleNewSocketConnection(w http.ResponseWriter, r *http.Request) error {
	ws, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Error in socket connection", err)
		return errors.New("could not upgrade to a websocket connection, bad request")
	}
	defer ws.Close()

	log.Print("Client Connected")

	s.Clients[ws] = "" // no room joined as of now thus empty string
	log.Print("Active Clients", s.Clients)
	log.Print("Active Clients in room", len(s.Rooms["123"]))
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			// clear the room from the ws WARNING : might need a mutex lock in delete operation
			delete(s.Rooms[s.Clients[ws]], ws)
			delete(s.Clients, ws)
			break
		}

		switch msg.Event {
		case "join":
			s.handleJoin(&msg, ws)
		case "ready":
			s.handleReady(ws, msg)
		case "candidate":
			s.handleCandidate(ws, msg)
		case "offer":
			s.handleOffer(ws, msg)
		case "answer":
			s.handleAnswer(ws, msg)
		}
	}

	log.Print("Client Disconnected")
	log.Print("Active Clients", s.Clients)
	log.Print("Active Clients in room", len(s.Rooms["123"]))

	return nil
}

func (s *SignallingService) handleJoin(msg *Message, ws *websocket.Conn) {
	// room is a map thus pass by reference
	room, exists := s.Rooms[msg.RoomId]
	if !exists {
		room = make(map[*websocket.Conn]bool)
		s.Rooms[msg.RoomId] = room
	}
	room[ws] = true
	s.Clients[ws] = msg.RoomId
	if len(room) == 1 {
		ws.WriteJSON(Message{Event: "created"})
	} else if len(room) == 2 {
		ws.WriteJSON(Message{Event: "joined"})
	} else {
		ws.WriteJSON(Message{Event: "full"})
	}
}

func (s *SignallingService) handleReady(ws *websocket.Conn, msg Message) {
	s.broadcastToRoom(msg.RoomId, Message{Event: "ready"}, ws)
}

func (s *SignallingService) handleCandidate(ws *websocket.Conn, msg Message) {
	s.broadcastToRoom(msg.RoomId, Message{Event: "candidate", Data: msg.Data}, ws)
}

func (s *SignallingService) handleOffer(ws *websocket.Conn, msg Message) {
	s.broadcastToRoom(msg.RoomId, Message{Event: "offer", Data: msg.Data}, ws)
}

func (s *SignallingService) handleAnswer(ws *websocket.Conn, msg Message) {
	s.broadcastToRoom(msg.RoomId, Message{Event: "answer", Data: msg.Data}, ws)
}

// broadcasting to all clients except the broadcaster in a particular room
func (s *SignallingService) broadcastToRoom(roomName string, msg Message, ignore *websocket.Conn) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	for client := range s.Rooms[roomName] {
		if client != ignore {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Print("broadcast error for a client, closing connection:", err)
				client.Close()
				delete(s.Rooms[roomName], client)
			}
		}
	}
}

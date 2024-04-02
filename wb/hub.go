package wb

import (
	"Hygieia/worker"
	"context"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"log"
)

var hub = newHub()

type Hub struct {
	clients      map[uint64]*Client
	Broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	EegData      chan []byte
	FatigueLevel chan []byte
	distributor  worker.TaskDistributor
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint64]*Client),
		EegData:    make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.SessionId] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.SessionId]; ok {
				close(client.send)
				delete(h.clients, client.SessionId)

			}
		case message := <-h.EegData:
			err := h.distributor.DistributeTaskSendEEGData(context.Background(), message, nil)
			if err != nil {
				log.Printf("failed to distribute task")
			}
		case message := <-h.Broadcast:
			for _, client := range h.clients {
				select {
				case client.send <- message:
				}
			}
		case message := <-h.FatigueLevel:
			var payload worker.PayloadSendEEGFatigueLevel
			if err := json.Unmarshal(message, &payload); err != nil {
				log.Printf("cannot marshal json payload")
			}
			level, err := proto.Marshal(&payload.FatigueLevel)
			if err != nil {
				log.Printf("cannot marshal proto payload")
			}
			sessionId := payload.SessionId
			if _, ok := h.clients[sessionId]; ok {
				h.clients[sessionId].send <- level
			} else {
				log.Printf("client has gone away")
			}
		}
	}
}

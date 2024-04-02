package sse

import (
	"Hygieia/middleware"
	"Hygieia/token"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"time"
)

const patience time.Duration = time.Second * 1

type (
	NotificationEvent struct {
		EventName  string
		ReceiveUid []uint64
		Payload    interface{}
	}

	NotifierChan chan NotificationEvent

	Broker struct {
		Broadcast NotifierChan
		// Events are pushed to this channel by the main events-gathering routine
		Notifier NotifierChan

		// New client connections
		newClients chan *client

		// Closed client connections
		closingClients chan *client

		// Client connections registry
		clients map[uint64]*client
	}
)

func NewBroker() (broker *Broker) {
	// Instantiate a broker
	return &Broker{
		Broadcast:      make(NotifierChan, 1),
		Notifier:       make(NotifierChan, 1),
		newClients:     make(chan *client),
		closingClients: make(chan *client),
		clients:        make(map[uint64]*client),
	}
}

func (broker *Broker) ServeHTTP(c *gin.Context) {
	authPayload := c.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	client := NewClient(authPayload.UID)

	broker.newClients <- client

	defer func() {
		broker.closingClients <- client
	}()

	c.Stream(func(w io.Writer) bool {
		// Emit Server Sent Events compatible
		event := <-client.notificationChannel
		c.SSEvent(event.EventName, event.Payload)
		c.Writer.Flush()
		return true
	})
}

func (broker *Broker) Listen() {
	for {
		select {
		case s := <-broker.newClients:

			broker.clients[s.uid] = s
			log.Printf("Client added. %d registered clients", len(broker.clients))
		case s := <-broker.closingClients:

			delete(broker.clients, s.uid)
			log.Printf("Removed client. %d registered clients", len(broker.clients))
		case event := <-broker.Notifier:
			for _, uid := range event.ReceiveUid {
				if client, ok := broker.clients[uid]; ok {
					client.notificationChannel <- event
				}
			}
		case event := <-broker.Broadcast:
			for _, client := range broker.clients {
				client.notificationChannel <- event
			}
		}
	}
}

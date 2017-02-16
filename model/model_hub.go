package model

import "log"

type Hub struct {
	clients map[string][]*Client
	Broadcast chan *HubMessage
	register chan *Client
	unregister chan *Client
}

type HubMessage struct {
	Key string
	Message []byte
}
func NewHub() *Hub {
	return &Hub{
		Broadcast: make(chan *HubMessage),
		register: make(chan *Client),
		unregister: make(chan *Client),
		clients: make(map[string][]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Println("Registering new client")
			h.clients[client.key] = append(h.clients[client.key], client)
		case client := <-h.unregister:
			log.Println("Unregistering new client")
			if _, ok := h.clients[client.key]; ok {
				//delete(h.clients[client.key], client)
				close(client.send)
			}
		case hubMessage := <-h.Broadcast:
			log.Println("sending", string(hubMessage.Message))
			log.Println("to", hubMessage.Key)
			if _, ok := h.clients[hubMessage.Key]; ok {
				for _, client := range h.clients[hubMessage.Key] {
					select {
					case client.send <- hubMessage.Message:
					default:
						close(client.send)
					//delete(h.clients, client)
					}
				}
			}
		}
	}
}
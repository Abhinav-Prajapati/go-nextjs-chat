package ws

import "fmt"

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Brodcast   chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Brodcast:   make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			// check if room with given id exists
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				// check if client already exists
				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
					newClient := h.Rooms[cl.RoomID].Clients[cl.ID].Username
					fmt.Println("new registere client name ", newClient)

				}
			}
		case cl := <-h.Unregister:
			// check if room exists
			if _, ok := h.Rooms[cl.RoomID]; ok {
				// send notification to other client if there are any
				if len(h.Rooms[cl.RoomID].Clients) > 0 {
					h.Brodcast <- &Message{
						Content:  "user left the chat",
						RoomID:   cl.RoomID,
						Username: cl.Username,
					}
				}
				delete(h.Rooms[cl.RoomID].Clients, cl.ID)
				close(cl.Message)
			}

		case m := <-h.Brodcast:
			if _, ok := h.Rooms[m.RoomID]; ok {
				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}

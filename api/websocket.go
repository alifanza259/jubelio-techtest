package api

import (
	"fmt"
	"log"

	"github.com/alifanza259/jubelio-interview/token"
	"github.com/gofiber/websocket/v2"
)

func (server *Server) handleWs(c *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		server.unregister <- c
		c.Close()
	}()

	// Register the client
	server.register <- c

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			server.broadcast <- string(message)
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}

func (server *Server) RunHub() {
	for {
		select {
		case connection := <-server.register:
			user := connection.Locals("user").(*token.Payload)
			client := Client{
				Conn:   connection,
				UserID: user.ID,
			}
			server.clients[connection] = client
			log.Println("connection registered")

		case message := <-server.broadcast:
			log.Println("message received:", message)

			// Send the message to all clients
			for connection := range server.clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					server.unregister <- connection
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
				}
			}

		case connection := <-server.unregister:
			// Remove the client from the hub
			delete(server.clients, connection)

			log.Println("connection unregistered")
		}
	}
}

func (server *Server) BroadcastSupabaseMessage(b []byte, senderID int, receiverID int) {
	for ws, val := range server.clients {
		if val.UserID == senderID || val.UserID == receiverID {
			go func(ws *websocket.Conn) {
				if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
					fmt.Println(err)
					delete(server.clients, ws)
				}
			}(ws)
		}
	}
}

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alifanza259/jubelio-interview/token"
	"github.com/gofiber/fiber/v2"
)

type sendMessageRequest struct {
	Content    string `json:"content"`
	ReceiverID int    `json:"receiver_id"`
}

func (server *Server) sendMessage(c *fiber.Ctx) error {
	user := c.Locals("user").(*token.Payload)

	var req sendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	row := server.db.QueryRow("SELECT * FROM users WHERE id = $1", req.ReceiverID)

	var receiver User
	if err := row.Scan(&receiver.ID, &receiver.Email, &receiver.Password, &receiver.Name, &receiver.ImageUrl, &receiver.CreatedAt); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("receiver not found")
	}

	type Message struct {
		Content    string `json:"content"`
		SenderID   int    `json:"sender_id"`
		ReceiverID int    `json:"receiver_id"`
	}

	message := Message{
		Content:    req.Content,
		ReceiverID: req.ReceiverID,
		SenderID:   user.ID,
	}
	jsonStr, _ := json.Marshal(message)
	reqSupabase, err := http.NewRequest("POST", server.config.SupabaseHost+"/rest/v1/messages?columns=%22content%22%2C%22receiver_id%22%2C%22sender_id%22&select=*", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
		return err
	}
	reqSupabase.Header.Set("Apikey", server.config.SupabaseApiKey)
	resp, err := http.DefaultClient.Do(reqSupabase)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return err
	}

	return c.JSON("Success")
}

type Message struct {
	ID            int       `json:"id"`
	Content       string    `json:"content"`
	SenderID      string    `json:"sender_id"`
	ReceiverID    string    `json:"receiver_id"`
	CreatedAt     time.Time `json:"created_at"`
	SenderName    string    `json:"sender_name"`
	ReceiverName  string    `json:"receiver_name"`
	SenderImage   string    `json:"sender_image_url"`
	ReceiverImage string    `json:"receiver_image_url"`
}

func (server *Server) searchHistory(c *fiber.Ctx) error {
	q := c.Query("q")

	rows, err := server.db.Query(`SELECT 
		messages.id,
		messages.sender_id,
		messages.receiver_id,
		messages.content,
		messages.created_at,
		sender.name,
		receiver.name,
		sender.image_url,
		receiver.image_url
	FROM messages 
	JOIN users sender ON sender.id = messages.sender_id 
	JOIN users receiver ON receiver.id = messages.receiver_id 
	WHERE content LIKE $1
	ORDER BY messages.created_at DESC`, "%"+q+"%")

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt, &msg.SenderName, &msg.ReceiverName, &msg.SenderImage, &msg.ReceiverImage); err != nil {
			fmt.Println(err)
			return err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(messages)
}

func (server *Server) getMessages(c *fiber.Ctx) error {
	user := c.Locals("user").(*token.Payload)
	receiverId := c.Query("receiver_id")

	rows, err := server.db.Query(`SELECT 
		messages.id,
		messages.sender_id,
		messages.receiver_id,
		messages.content,
		messages.created_at,
		sender.name,
		receiver.name,
		sender.image_url,
		receiver.image_url
	FROM messages 
	JOIN users sender ON sender.id = messages.sender_id 
	JOIN users receiver ON receiver.id = messages.receiver_id  
	WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
	ORDER BY messages.created_at DESC`, user.ID, receiverId)

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.CreatedAt, &msg.SenderName, &msg.ReceiverName, &msg.SenderImage, &msg.ReceiverImage); err != nil {
			return err
		}
		messages = append(messages, msg)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(messages)
}

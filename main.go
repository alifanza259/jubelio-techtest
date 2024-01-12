package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alifanza259/jubelio-interview/api"
	"github.com/alifanza259/jubelio-interview/token"
	"github.com/alifanza259/jubelio-interview/util"
	_ "github.com/lib/pq"
	realtimego "github.com/overseedio/realtime-go"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load env file: %s", err)
	}

	pool, err := sql.Open(config.DBDriver, config.DBUrl)
	if err != nil {
		log.Fatalf("cannot connect to database: %s", err)
	}

	tokenMaker := token.NewJWTMaker(config.JwtSecretKey)

	server, err := api.NewServer(pool, config, tokenMaker)
	if err != nil {
		log.Fatalf("cannot create server: %s", err)
	}

	// handle incoming messages from supabase and save to local db. And broadcast to websocket clients
	go handleIncomingMessages(server, config, pool)

	go server.RunHub() // register and unregister websocket clients

	err = server.Start(config.Host)
	if err != nil {
		log.Fatalf("cannot start server: %s", err)
	}
}

func handleIncomingMessages(server *api.Server, config util.Config, db *sql.DB) {
	// create client
	c, err := realtimego.NewClient(config.SupabaseHost, config.SupabaseApiKey)
	if err != nil {
		log.Fatal(err)
	}

	// connect to server
	err = c.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// create and subscribe to channel
	database := "realtime"
	schema := "public"
	table := "messages"
	ch, err := c.Channel(realtimego.WithTable(&database, &schema, &table))
	if err != nil {
		log.Fatal(err)
	}

	// setup hooks
	ch.OnInsert = func(m realtimego.Message) {
		type Column struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}
		type Record struct {
			Content    string    `json:"content"`
			SenderId   int       `json:"sender_id"`
			ReceiverId int       `json:"receiver_id"`
			CreatedAt  time.Time `json:"created_at"`
		}
		type Payload struct {
			Columns         []Column  `json:"columns"`
			CommitTimestamp time.Time `json:"commit_timestamp"`
			Errors          error     `json:"errors"`
			Record          Record    `json:"record"`
			Schema          string    `json:"schema"`
			Table           string    `json:"table"`
			OperationType   string    `json:"type"`
		}

		var payload Payload

		jsonStr, _ := json.Marshal(m.Payload)
		json.Unmarshal(jsonStr, &payload)

		rows, err := db.Query("INSERT INTO messages (content, sender_id, receiver_id, created_at) VALUES ($1, $2, $3, $4)",
			payload.Record.Content, payload.Record.SenderId, payload.Record.ReceiverId, time.Now())
		if err != nil {
			fmt.Println(err)
			return
		}

		defer rows.Close()

		type User struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			ImageUrl string `json:"image_url"`
		}
		var user User

		row := db.QueryRow("SELECT id,name,image_url FROM users WHERE id = $1", payload.Record.SenderId)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := row.Scan(&user.ID, &user.Name, &user.ImageUrl); err != nil {
			fmt.Println(err)
			return
		}

		broadcastMessage := struct {
			Payload Record `json:"payload"`
			Sender  User   `json:"sender"`
			Type    string `json:"type"`
		}{
			Payload: payload.Record,
			Type:    "insert",
			Sender:  user,
		}
		broadcastMessageJson, _ := json.Marshal(broadcastMessage)

		server.BroadcastSupabaseMessage(broadcastMessageJson, payload.Record.SenderId, payload.Record.ReceiverId)

	}

	// subscribe to channel
	err = ch.Subscribe()
	if err != nil {
		log.Fatal(err)
	}
}

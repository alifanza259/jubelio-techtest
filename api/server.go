package api

import (
	"database/sql"

	"github.com/alifanza259/jubelio-interview/token"
	"github.com/alifanza259/jubelio-interview/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Server struct {
	db         *sql.DB
	tokenMaker token.Maker
	config     util.Config
	router     *fiber.App
	clients    map[*websocket.Conn]Client
	register   chan *websocket.Conn
	broadcast  chan string
	unregister chan *websocket.Conn
}

type Client struct {
	UserID int
	Conn   *websocket.Conn
}

func NewServer(db *sql.DB, config util.Config, tokenMaker token.Maker) (*Server, error) {
	server := &Server{config: config,
		db:         db,
		tokenMaker: tokenMaker,
		clients:    make(map[*websocket.Conn]Client),
		register:   make(chan *websocket.Conn),
		broadcast:  make(chan string),
		unregister: make(chan *websocket.Conn),
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	app := fiber.New()
	app.Use(CORSMiddleware())

	app.Static("/public", "./public")

	app.Post("/login", server.login)

	app.Get("/messages", authMiddleware(server.tokenMaker), server.getMessages)
	app.Get("/search", authMiddleware(server.tokenMaker), server.searchHistory)
	app.Post("/messages", authMiddleware(server.tokenMaker), server.sendMessage)

	app.Get("/ws", websocketAuthMiddleware(server.tokenMaker), websocket.New(server.handleWs))

	server.router = app
}

func (server *Server) Start(address string) error {
	return server.router.Listen(address)
}

func CORSMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Response().Header.Set("Access-Control-Allow-Origin", "*")
		c.Response().Header.Set("Access-Control-Allow-Methods", "*")
		c.Response().Header.Set("Access-Control-Allow-Headers", "*")
		if c.Method() == "OPTIONS" {
			return c.SendString("ok!")
		}
		return c.Next()
	}
}

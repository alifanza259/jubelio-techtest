package api

import (
	"time"

	"github.com/alifanza259/jubelio-interview/util"
	"github.com/gofiber/fiber/v2"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"password,omitempty"`
	ImageUrl  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

func (server *Server) login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	row := server.db.QueryRow("SELECT * FROM users WHERE email = $1", req.Email)

	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.ImageUrl, &user.CreatedAt); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("invalid credentials")
	}

	err := util.CheckPassword(req.Password, user.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON("invalid credentials")
	}

	token, expiresAt, err := server.tokenMaker.CreateToken(user.Email, user.ID, 24*time.Hour)
	if err != nil {
		return err
	}

	user.Password = ""
	return c.JSON(fiber.Map{
		"token":      token,
		"expires_at": expiresAt,
		"user":       user,
	})
}

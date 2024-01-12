package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alifanza259/jubelio-interview/token"
	"github.com/gofiber/fiber/v2"
)

func handler(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authorizationHeader := ctx.GetReqHeaders()["Authorization"]
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}

		fields := strings.Split(authorizationHeader[0], " ")
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != "bearer" {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			fmt.Println(err)
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}

		ctx.Locals("user", payload)
		return ctx.Next()
	}
}
func authMiddleware(tokenMaker token.Maker) fiber.Handler {
	return handler(tokenMaker)
}

func websocketAuthMiddleware(tokenMaker token.Maker) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		accessToken := ctx.Query("access_token")

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			return ctx.Status(http.StatusUnauthorized).JSON(err)
		}

		ctx.Locals("user", payload)
		return ctx.Next()
	}
}

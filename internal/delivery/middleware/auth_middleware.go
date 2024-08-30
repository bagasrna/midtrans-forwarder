// middleware/auth_middleware.go
package middleware

import (
	"midtrans-forwarder/pkg/auth"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func BearerAuth(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if header == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fiber.StatusUnauthorized})
	}
	bearer := strings.SplitN(header, " ", 2)
	if bearer[0] != "Bearer" || len(bearer) != 2 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fiber.StatusBadRequest})
	}
	token, err := paseto.Decode(bearer[1])
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fiber.StatusUnauthorized})
	}
	sub, err := token.GetSubject()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": fiber.StatusUnauthorized})
	}
	id, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return err
	}

	c.Locals("clientID", id)
	return c.Next()
}
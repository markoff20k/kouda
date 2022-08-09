package middlewares

import (
	"net"

	"github.com/gofiber/fiber/v2"
)

func ParseIP(c *fiber.Ctx) error {
	if c.Locals("remote_ip") == nil {
		c.Locals("remote_ip", net.ParseIP(c.IP()))
	}

	return c.Next()
}

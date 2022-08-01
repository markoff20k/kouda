package public

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetTimestamp Get current timestamp
func GetTimestamp(c *fiber.Ctx) error {
	return c.JSON(time.Now())
}

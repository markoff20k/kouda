package public

import "github.com/gofiber/fiber/v2"

func NewRouter(router fiber.Router) {
	router.Get("/timestamp", GetTimestamp)
}

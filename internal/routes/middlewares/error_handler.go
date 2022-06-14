package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg"

	"github.com/zsmartex/kouda/log"
	"github.com/zsmartex/kouda/params"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		// Override status code if fiber.Error type
		code = e.Code

		switch code {
		case fiber.StatusNotFound:
			return c.Status(code).JSON("404 Not Found")
		}
	} else if e, ok := err.(*pkg.Error); ok {
		// Override status code if pkg.Error type
		code = e.Code

		return c.Status(code).JSON(pkg.Error{
			Errors: e.Errors,
		})
	}

	log.Error(err)
	return c.Status(code).JSON(pkg.Error{
		Errors: params.ErrServerInternal.Errors,
	})
}

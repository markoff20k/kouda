package params

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg/v2"
)

var (
	ErrServerInternal = pkg.NewError(fiber.StatusInternalServerError, "server.internal_error")

	ErrServerInvalidQuery = pkg.NewError(fiber.StatusBadRequest, "server.method.invalid_message_query")

	ErrServerInvalidBody = pkg.NewError(fiber.StatusBadRequest, "server.method.invalid_message_body")

	ErrRecordNotFound = pkg.NewError(fiber.StatusNotFound, "record.not_found")

	ErrJWTDecodeAndVerify = pkg.NewError(fiber.StatusUnauthorized, "jwt.decode_and_verify")
)

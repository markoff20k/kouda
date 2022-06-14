package middlewares

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg/jwt"

	"github.com/zsmartex/kouda/config"
	"github.com/zsmartex/kouda/params"
)

func Authorization() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		jwtHeader := c.Get("Authorization")
		jwtToken := strings.TrimPrefix(jwtHeader, "Bearer ")

		key_store := jwt.KeyStore{}
		if err := key_store.LoadPublicKeyFromString(config.Env.JWTPublicKey); err != nil {
			log.Panicf("Failed to load public key: %v", err)
		}

		auth, err := jwt.ParseAndValidate(jwtToken, key_store.PublicKey)
		if err != nil {
			return params.ErrJWTDecodeAndVerify
		}

		c.Locals("current_member", auth)

		return c.Next()
	}
}

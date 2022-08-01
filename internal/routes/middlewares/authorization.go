package middlewares

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/volatiletech/null/v9"
	"github.com/zsmartex/pkg/v2/gpa/filters"
	"github.com/zsmartex/pkg/v2/jwt"

	"github.com/zsmartex/kouda/config"
	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/kouda/params"
	"github.com/zsmartex/kouda/usecases"
)

func Authorization(memberUsecase usecases.MemberUsecase) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		jwtHeader := c.Get("Authorization")
		jwtToken := strings.TrimPrefix(jwtHeader, "Bearer ")

		keyStore := jwt.KeyStore{}
		if err := keyStore.LoadPublicKeyFromString(config.Env.JWTPublicKey); err != nil {
			log.Panicf("Failed to load public key: %v", err)
		}

		auth, err := jwt.ParseAndValidate(jwtToken, keyStore.PublicKey)
		if err != nil {
			return params.ErrJWTDecodeAndVerify
		}

		member := &models.Member{
			UID:   auth.UID,
			Email: auth.Email,
		}
		memberUsecase.FirstOrCreate(
			c.Context(),
			&member,
			filters.WithFieldEqual("uid", auth.UID),
			filters.WithAssign(models.Member{
				Username: null.NewString(auth.Username, auth.Username != ""),
				Level:    auth.Level,
				Role:     auth.Role,
				State:    models.MemberState(auth.State),
			}),
		)

		c.Locals("current_member", member)

		return c.Next()
	}
}

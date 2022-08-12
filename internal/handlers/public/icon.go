package public

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/gpa/filters"
)

var ErrIconDoesntExist = pkg.NewError(fiber.StatusNotFound, "public.icon.doesnt_exist")

func (h Handler) GetIconImage(c *fiber.Ctx) error {
	code := c.Params("code")
	ctx := c.Context()

	icon, err := h.iconUsecase.First(ctx, filters.WithFieldEqual("code", code), filters.WithFieldEqual("state", models.IconStateEnabled))
	if err != nil {
		return ErrIconDoesntExist
	}

	body, err := h.Uploader.GetBodyContent(ctx, fmt.Sprintf("icons/%s", icon.Code))
	if err != nil {
		return err
	}

	c.Set("Content-Type", "image/jpeg")

	return c.Send(body)
}

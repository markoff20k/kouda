package public

import (
	"fmt"
	"net/http"
	"strings"

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

	// get type of image
	mimeType := http.DetectContentType(body)
	typeFile := strings.Replace(mimeType, "image/", "", -1)

	c.Set("Content-Type", fmt.Sprintf("image/%s", typeFile))

	return c.Send(body)
}

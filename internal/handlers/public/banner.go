package public

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/gpa/filters"
	"github.com/zsmartex/pkg/v2/queries"

	"github.com/zsmartex/kouda/internal/handlers/admin/entities"
	"github.com/zsmartex/kouda/internal/handlers/helpers"
	"github.com/zsmartex/kouda/internal/models"
)

var (
	ErrBannerNotFound = pkg.NewError(fiber.StatusNotFound, "public.banner.not_found")
)

func (h Handler) GetBanners(c *fiber.Ctx) error {
	type Params struct {
		queries.Order
		queries.Period
		queries.Pagination
	}

	params := new(Params)
	if err := helpers.QueryParser(c, params, "public.banner"); err != nil {
		return err
	}

	q := make([]gpa.Filter, 0)
	q = append(
		q,
		filters.WithFieldEqual("state", models.BannerStateEnabled),
		filters.WithPageable(params.Page, params.Limit),
		filters.WithOrder(fmt.Sprintf("%s %s", params.OrderBy, params.Ordering)),
	)

	banners := h.bannerUsecase.Find(q...)

	bannerEntities := make([]*entities.Banner, 0)
	for _, banner := range banners {
		bannerEntities = append(bannerEntities, entities.BannerToEntity(banner))
	}

	return c.JSON(bannerEntities)
}

func (h Handler) GetBannerImage(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	banner, err := h.bannerUsecase.First(filters.WithFieldEqual("uuid", uuid))
	if err != nil {
		return ErrBannerNotFound
	}

	body, err := h.Uploader.GetBodyContent(fmt.Sprintf("banners/%s.%s", banner.UUID.String(), banner.Type))
	if err != nil {
		return err
	}

	c.Set("Content-Type", "image/jpeg")

	return c.Send(body)
}

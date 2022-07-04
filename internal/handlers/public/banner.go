package public

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/strutil"
	"github.com/zsmartex/pkg/gpa"
	"github.com/zsmartex/pkg/gpa/filters"
	"github.com/zsmartex/pkg/queries"

	"github.com/zsmartex/kouda/internal/handlers/admin/entities"
	"github.com/zsmartex/kouda/internal/handlers/helpers"
	"github.com/zsmartex/kouda/internal/models"
)

func (h Handler) GetBanners(c *fiber.Ctx) error {
	type Params struct {
		Tag string `query:"tag"`
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

	if strutil.IsNotBlank(params.Tag) {
		q = append(q, filters.WithFieldEqual("tag", params.Tag))
	}

	banners := h.bannerUsecase.Find(q...)

	bannerEntities := make([]*entities.Banner, 0)
	for _, banner := range banners {
		bannerEntities = append(bannerEntities, entities.BannerToEntity(banner))
	}

	for _, e := range bannerEntities {
		image_url, err := h.Uploader.GetURL(fmt.Sprintf("banners/%s.%s", e.UUID.String(), e.Type))
		if err != nil {
			return err
		}

		e.ImageURL = image_url
	}

	return c.JSON(bannerEntities)
}

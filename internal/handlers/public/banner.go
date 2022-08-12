package public

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
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
	ErrBannerOverSize = pkg.NewError(fiber.StatusNotFound, "public.banner.over_size")
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

	ctx := c.Context()

	q := make([]gpa.Filter, 0)
	q = append(
		q,
		filters.WithFieldEqual("state", models.BannerStateEnabled),
		filters.WithPageable(params.Page, params.Limit),
		filters.WithOrder(fmt.Sprintf("%s %s", params.OrderBy, params.Ordering)),
	)

	banners := h.bannerUsecase.Find(ctx, q...)

	bannerEntities := make([]*entities.Banner, 0)
	for _, banner := range banners {
		bannerEntities = append(bannerEntities, entities.BannerToEntity(banner))
	}

	return c.JSON(bannerEntities)
}

func (h Handler) GetBannerImage(c *fiber.Ctx) error {
	type Params struct {
		D string `query:"d" validate:"sizeBanner"`
	}

	params := new(Params)
	if err := helpers.QueryParser(c, params, "public.banner"); err != nil {
		return err
	}

	uuid := c.Params("uuid")
	ctx := c.Context()

	banner, err := h.bannerUsecase.First(ctx, filters.WithFieldEqual("uuid", uuid), filters.WithFieldEqual("state", models.BannerStateEnabled))
	if err != nil {
		return ErrBannerNotFound
	}

	body, err := h.Uploader.GetBodyContent(ctx, fmt.Sprintf("banners/%s.%s", banner.UUID.String(), banner.Type))
	if err != nil {
		return err
	}

	if len(params.D) > 0 {
		image, err := vips.NewImageFromBuffer(body)
		if err != nil {
			return err
		}

		width, err := strconv.Atoi(params.D[:strings.Index(params.D, "x")])
		if err != nil {
			return err
		}

		height, err := strconv.Atoi(params.D[strings.Index(params.D, "x")+1:])
		if err != nil {
			return err
		}

		if image.Width() < width || image.Height() < height {
			return ErrBannerOverSize
		}

		if err := image.ThumbnailWithSize(width, height, vips.Interesting(vips.ImageTypeMagick), vips.SizeForce); err != nil {
			return err
		}

		ep := vips.NewDefaultJPEGExportParams()
		imageBytes, _, err := image.Export(ep)
		if err != nil {
			return err
		}

		body = imageBytes
	}

	c.Set("Content-Type", "image/jpeg")

	return c.Send(body)
}

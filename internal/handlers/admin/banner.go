package admin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gookit/goutil/strutil"
	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/gpa/filters"
	"github.com/zsmartex/pkg/v2/queries"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/internal/handlers/admin/entities"
	"github.com/zsmartex/kouda/internal/handlers/helpers"
	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/kouda/utils"
)

var (
	ErrBannerMissingUUID  = pkg.NewError(fiber.StatusBadRequest, "admin.banner.missing_uuid")
	ErrBannerMissingImage = pkg.NewError(fiber.StatusUnprocessableEntity, "resource.banner.missing_image")
	ErrBannerInvalidImage = pkg.NewError(fiber.StatusUnprocessableEntity, "resource.banner.invalid_image")
	ErrBannerDoesntExist  = pkg.NewError(fiber.StatusUnprocessableEntity, "resource.banner.doesnt_exist")
)

func (h Handler) GetBanners(c *fiber.Ctx) error {
	type Params struct {
		State models.BannerState `query:"state" validate:"bannerState"`
		queries.Order
		queries.Period
		queries.Pagination
	}

	params := new(Params)
	if err := helpers.QueryParser(c, params, "admin.banner"); err != nil {
		return err
	}

	q := make([]gpa.Filter, 0)
	q = append(
		q,
		filters.WithPageable(params.Page, params.Limit),
		filters.WithOrder(fmt.Sprintf("%s %s", params.OrderBy, params.Ordering)),
	)

	if strutil.IsNotBlank(string(params.State)) {
		q = append(q, filters.WithFieldEqual("state", params.State))
	}

	banners := h.bannerUsecase.Find(q...)

	bannerEntities := make([]*entities.Banner, 0)
	for _, banner := range banners {
		bannerEntities = append(bannerEntities, entities.BannerToEntity(banner))
	}

	for _, e := range bannerEntities {
		image_url, err := h.uploader.GetURL(fmt.Sprintf("banners/%s.%s", e.UUID.String(), e.Type))
		if err != nil {
			return err
		}

		e.ImageURL = image_url
	}

	return c.JSON(bannerEntities)
}

func (h Handler) CreateBanner(c *fiber.Ctx) error {
	type Params struct {
		URL   string             `json:"url" validate:"required"`
		State models.BannerState `json:"state" validate:"required|bannerState"`
	}

	params := new(Params)
	if err := helpers.BodyParser(c, params, "admin.banner"); err != nil {
		return err
	}

	img, err := c.FormFile("image")
	if err != nil {
		return ErrBannerMissingImage
	}

	file, err := img.Open()
	if err != nil {
		return ErrBannerInvalidImage
	}

	if !utils.ValidateImageFile(file) {
		return ErrBannerInvalidImage
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return err
	}

	file_bytes := buf.Bytes()

	// get type of image
	mime_type := http.DetectContentType(file_bytes)
	type_file := strings.Replace(mime_type, "image/", "", -1)

	var banner *models.Banner

	if err := h.bannerUsecase.Transaction(func(tx *gorm.DB) error {
		banner = &models.Banner{
			UUID:  uuid.New(),
			State: params.State,
			Type:  type_file,
			URL:   params.URL,
		}

		h.bannerUsecase.WithTrx(tx).Create(&banner)

		url := fmt.Sprintf("banners/%s.%s", banner.UUID.String(), type_file)

		if _, err := h.uploader.Upload(url, file_bytes); err != nil {
			return err
		}

		h.bannerUsecase.WithTrx(tx).Updates(&banner, models.Banner{
			URL: url,
		})

		return nil
	}); err != nil {
		panic(err)
	}

	return c.Status(201).JSON(entities.BannerToEntity(banner))
}

func (h Handler) UpdateBanner(c *fiber.Ctx) error {
	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return ErrBannerMissingUUID
	}

	type Params struct {
		URL   string             `json:"url" validate:"required"`
		State models.BannerState `json:"state" validate:"required|bannerState"`
	}

	params := new(Params)
	if err := helpers.BodyParser(c, params, "admin.banner"); err != nil {
		return err
	}

	banner, err := h.bannerUsecase.First(filters.WithAssign(&models.Banner{UUID: uuid}))
	if err != nil {
		return err
	}

	h.bannerUsecase.Updates(banner, models.Banner{State: params.State, URL: params.URL})

	return c.JSON(201)
}

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
	"github.com/zsmartex/pkg"
	"github.com/zsmartex/pkg/gpa"
	"github.com/zsmartex/pkg/gpa/filters"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/internal/handlers/admin/entities"
	"github.com/zsmartex/kouda/internal/handlers/helpers"
	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/kouda/pkg/queries"
	"github.com/zsmartex/kouda/utils"
)

var (
	ErrBannerMissingUUID  = pkg.NewError(fiber.StatusBadRequest, "admin.banner.missing_uuid")
	ErrBannerMissingImage = pkg.NewError(fiber.StatusUnprocessableEntity, "resource.banner.missing_image")
	ErrBannerInvalidImage = pkg.NewError(fiber.StatusUnprocessableEntity, "resource.banner.invalid_image")
)

func (h Handler) GetBanners(c *fiber.Ctx) error {
	type Params struct {
		Tag   string             `query:"tag"`
		State models.BannerState `query:"state"`
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

	if strutil.IsNotBlank(params.Tag) {
		q = append(q, filters.WithFieldEqual("tag", params.Tag))
	}

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
		Tag   string             `json:"tag" validate:"required"`
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

	tx := h.bannerUsecase.DoTrx()

	if err := h.bannerUsecase.HandleTrx(tx, func(tx *gorm.DB) error {
		banner = &models.Banner{
			UUID:  uuid.New(),
			Tag:   params.Tag,
			State: params.State,
			Type:  type_file,
			URL:   "",
		}

		h.bannerUsecase.Create(&banner)

		url := fmt.Sprintf("banners/%s.%s", banner.UUID.String(), type_file)

		if _, err := h.uploader.Upload(url, file_bytes); err != nil {
			return err
		}

		h.bannerUsecase.Updates(&banner, models.Banner{
			URL: url,
		})

		return nil
	}); err != nil {
		panic(err)
	}

	return c.Status(201).JSON("aaaa")
}

func (h Handler) UpdateBanner(c *fiber.Ctx) error {
	uuid, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return ErrBannerMissingUUID
	}

	type Params struct {
		Tag   string             `json:"tag"`
		State models.BannerState `json:"state"`
	}

	params := new(Params)
	if err := helpers.BodyParser(c, params, "admin.banner"); err != nil {
		return err
	}

	update_params := &models.Banner{
		UUID: uuid,
	}

	h.bannerUsecase.FirstOrCreate(update_params, filters.WithAssign(&models.Banner{
		Tag:   params.Tag,
		State: params.State,
	}))

	return c.JSON(201)
}

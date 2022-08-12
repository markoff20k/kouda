package admin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
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
	"github.com/zsmartex/kouda/types"
	"github.com/zsmartex/kouda/utils"
)

var (
	ErrBannerMissingUUID  = pkg.NewError(fiber.StatusBadRequest, "admin.banner.missing_uuid")
	ErrBannerMissingImage = pkg.NewError(fiber.StatusUnprocessableEntity, "admin.banner.missing_image")
	ErrBannerInvalidImage = pkg.NewError(fiber.StatusUnprocessableEntity, "admin.banner.invalid_image")
	ErrBannerDoesntExist  = pkg.NewError(fiber.StatusNotFound, "admin.banner.doesnt_exist")
	ErrBannerOverSize     = pkg.NewError(fiber.StatusNotFound, "admin.banner.over_size")
	ErrBannerNotFound     = pkg.NewError(fiber.StatusNotFound, "admin.banner.not_found")
)

func (h Handler) GetBanners(c *fiber.Ctx) error {
	h.adminAuthorize(c, types.AbilityAdminPermission(AbilityAdminPermissionRead), "Banner")

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

	ctx := c.Context()

	q := make([]gpa.Filter, 0)
	q = append(
		q,
		filters.WithPageable(params.Page, params.Limit),
		filters.WithOrder(fmt.Sprintf("%s %s", params.OrderBy, params.Ordering)),
	)

	if strutil.IsNotBlank(string(params.State)) {
		q = append(q, filters.WithFieldEqual("state", params.State))
	}

	banners := h.bannerUsecase.Find(ctx, q...)
	total := h.bannerUsecase.Count(ctx, q...)

	bannerEntities := make([]*entities.Banner, 0)
	for _, banner := range banners {
		bannerEntities = append(bannerEntities, entities.BannerToEntity(banner))
	}

	c.Set("Page", fmt.Sprint(params.Page))
	c.Set("Per-Size", fmt.Sprint(params.Limit))
	c.Set("Total", fmt.Sprint(total))

	return c.JSON(bannerEntities)
}

func (h Handler) CreateBanner(c *fiber.Ctx) error {
	h.adminAuthorize(c, types.AbilityAdminPermission(AbilityAdminPermissionManage), "Banner")

	type Params struct {
		URL   string             `json:"url" validate:"required"`
		State models.BannerState `json:"state" validate:"bannerState" default:"disabled"`
	}

	params := new(Params)
	if err := helpers.BodyParser(c, params, "admin.banner"); err != nil {
		return err
	}

	ctx := c.Context()

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

	fileBytes := buf.Bytes()

	// get type of image
	mimeType := http.DetectContentType(fileBytes)
	typeFile := strings.Replace(mimeType, "image/", "", -1)

	var banner *models.Banner

	if err := h.bannerUsecase.Transaction(func(tx *gorm.DB) error {
		banner = &models.Banner{
			UUID:  uuid.New(),
			State: params.State,
			Type:  typeFile,
			URL:   params.URL,
		}

		h.bannerUsecase.WithTrx(tx).Create(ctx, &banner)

		key := fmt.Sprintf("banners/%s.%s", banner.UUID.String(), typeFile)

		if err := h.uploader.Upload(ctx, key, fileBytes); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return c.Status(201).JSON(entities.BannerToEntity(banner))
}

func (h Handler) UpdateBanner(c *fiber.Ctx) error {
	h.adminAuthorize(c, types.AbilityAdminPermission(AbilityAdminPermissionManage), "Banner")

	type Params struct {
		UUID  uuid.UUID          `json:"uuid"`
		URL   string             `json:"url"`
		State models.BannerState `json:"state" validate:"bannerState"`
	}

	params := new(Params)
	if err := helpers.BodyParser(c, params, "admin.banner"); err != nil {
		return err
	}

	ctx := c.Context()

	banner, err := h.bannerUsecase.First(ctx, filters.WithFieldEqual("uuid", params.UUID))
	if err != nil {
		return ErrBannerDoesntExist
	}

	targetBanner := models.Banner{}

	if len(params.State) > 0 {
		targetBanner.State = params.State
	}

	if len(params.URL) > 0 {
		targetBanner.URL = params.URL
	}

	img, _ := c.FormFile("image")
	if img != nil {
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

		fileBytes := buf.Bytes()

		// get type of image
		mimeType := http.DetectContentType(fileBytes)
		typeFile := strings.Replace(mimeType, "image/", "", -1)

		key := fmt.Sprintf("banners/%s.%s", banner.UUID.String(), typeFile)

		if h.uploader.Delete(ctx, key) != nil {
			return err
		}

		if err := h.uploader.Upload(ctx, key, fileBytes); err != nil {
			return err
		}

		targetBanner.Type = typeFile
	}

	h.bannerUsecase.Updates(ctx, &banner, targetBanner)

	return c.JSON(201)
}

func (h Handler) DeleteBanner(c *fiber.Ctx) error {
	h.adminAuthorize(c, types.AbilityAdminPermission(AbilityAdminPermissionManage), "Banner")

	uuidBanner, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return ErrBannerMissingUUID
	}

	ctx := c.Context()

	banner, err := h.bannerUsecase.First(ctx, filters.WithFieldEqual("uuid", uuidBanner))
	if err != nil {
		return ErrBannerDoesntExist
	}

	key := fmt.Sprintf("banners/%s.%s", banner.UUID.String(), banner.Type)

	if err := h.uploader.Delete(ctx, key); err != nil {
		return err
	}

	h.bannerUsecase.Delete(ctx, models.Banner{}, filters.WithFieldEqual("uuid", uuidBanner))

	return c.JSON(201)
}

func (h Handler) GetBannerImage(c *fiber.Ctx) error {
	type Params struct {
		D string `query:"d"`
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

	body, err := h.uploader.GetBodyContent(ctx, fmt.Sprintf("banners/%s.%s", banner.UUID.String(), banner.Type))
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

	c.Set("Content-Type", fmt.Sprintf("image/%s", banner.Type))

	return c.Send(body)
}

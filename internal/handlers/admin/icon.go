package admin

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/strutil"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/internal/handlers/admin/entities"
	"github.com/zsmartex/kouda/internal/handlers/helpers"
	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/kouda/types"
	"github.com/zsmartex/kouda/utils"
	"github.com/zsmartex/pkg/v2"
	"github.com/zsmartex/pkg/v2/gpa"
	"github.com/zsmartex/pkg/v2/gpa/filters"
	"github.com/zsmartex/pkg/v2/queries"
)

var (
	ErrIconMissingCode  = pkg.NewError(fiber.StatusBadRequest, "admin.icon.missing_code")
	ErrIconMissingImage = pkg.NewError(fiber.StatusUnprocessableEntity, "admin.icon.missing_image")
	ErrIconInvalidImage = pkg.NewError(fiber.StatusUnprocessableEntity, "admin.icon.invalid_image")
	ErrIconDoesntExist  = pkg.NewError(fiber.StatusNotFound, "admin.icon.doesnt_exist")
	ErrIconOverSize     = pkg.NewError(fiber.StatusNotFound, "admin.icon.over_size")
	ErrIconNotFound     = pkg.NewError(fiber.StatusNotFound, "admin.icon.not_found")
)

func (h Handler) GetIcons(c *fiber.Ctx) error {
	h.adminAuthorize(c, types.AbilityAdminPermission(AbilityAdminPermissionRead), "Icon")

	type Params struct {
		State models.IconState `query:"state" validate:"iconState"`
		queries.Order
		queries.Period
		queries.Pagination
	}

	params := new(Params)
	if err := helpers.QueryParser(c, params, "admin.icon"); err != nil {
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

	icons := h.iconUsecase.Find(ctx, q...)
	total := h.iconUsecase.Count(ctx, q...)

	iconEntities := make([]*entities.Icon, 0)
	for _, icon := range icons {
		iconEntities = append(iconEntities, entities.IconToEntity(icon))
	}

	c.Set("Page", fmt.Sprint(params.Page))
	c.Set("Per-Size", fmt.Sprint(params.Limit))
	c.Set("Total", fmt.Sprint(total))

	return c.JSON(iconEntities)
}

func (h Handler) CreateIcon(c *fiber.Ctx) error {
	h.adminAuthorize(c, types.AbilityAdminPermission(AbilityAdminPermissionManage), "Icon")

	type Params struct {
		Code  string           `json:"code" validate:"required"`
		URL   string           `json:"url" validate:"required"`
		State models.IconState `json:"state" validate:"iconState" default:"disabled"`
	}

	params := new(Params)
	if err := helpers.BodyParser(c, params, "admin.icon"); err != nil {
		return err
	}

	ctx := c.Context()

	img, err := c.FormFile("image")
	if err != nil {
		return ErrIconMissingImage
	}

	file, err := img.Open()
	if err != nil {
		return ErrIconInvalidImage
	}

	if !utils.ValidateImageFile(file) {
		return ErrIconInvalidImage
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return err
	}

	fileBytes := buf.Bytes()

	var icon *models.Icon

	if err := h.iconUsecase.Transaction(func(tx *gorm.DB) error {
		icon = &models.Icon{
			Code:  params.Code,
			State: params.State,
			URL:   params.URL,
		}

		h.iconUsecase.WithTrx(tx).Create(ctx, &icon)

		key := fmt.Sprintf("icons/%s", icon.Code)

		if err := h.uploader.Upload(ctx, key, fileBytes); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return c.Status(201).JSON(entities.IconToEntity(icon))
}

func (h Handler) GetIconImage(c *fiber.Ctx) error {
	code := c.Params("code")
	ctx := c.Context()

	icon, err := h.iconUsecase.First(ctx, filters.WithFieldEqual("code", code))
	if err != nil {
		return ErrIconDoesntExist
	}

	body, err := h.uploader.GetBodyContent(ctx, fmt.Sprintf("icons/%s", icon.Code))
	if err != nil {
		return err
	}

	c.Set("Content-Type", "image/jpeg")

	return c.Send(body)
}

func (h Handler) UpdateIcon(c *fiber.Ctx) error {
	h.adminAuthorize(c, types.AbilityAdminPermission(AbilityAdminPermissionManage), "Icon")

	type Params struct {
		Code  string           `json:"code"`
		URL   string           `json:"url"`
		State models.IconState `json:"state" validate:"iconState"`
	}

	params := new(Params)
	if err := helpers.BodyParser(c, params, "admin.icon"); err != nil {
		return err
	}

	ctx := c.Context()

	icon, err := h.iconUsecase.First(ctx, filters.WithFieldEqual("code", params.Code))
	if err != nil {
		return ErrIconDoesntExist
	}

	targetIcon := models.Icon{}

	if len(params.State) > 0 {
		targetIcon.State = params.State
	}

	if len(params.URL) > 0 {
		targetIcon.URL = params.URL
	}

	img, _ := c.FormFile("image")
	if img != nil {
		file, err := img.Open()
		if err != nil {
			return ErrIconInvalidImage
		}

		if !utils.ValidateImageFile(file) {
			return ErrIconInvalidImage
		}

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return err
		}

		fileBytes := buf.Bytes()

		key := fmt.Sprintf("icons/%s", icon.Code)

		if h.uploader.Delete(ctx, key) != nil {
			return err
		}

		if err := h.uploader.Upload(ctx, key, fileBytes); err != nil {
			return err
		}
	}

	h.iconUsecase.Updates(ctx, &icon, targetIcon)

	return c.JSON(201)
}

func (h Handler) DeleteIcon(c *fiber.Ctx) error {
	h.adminAuthorize(c, types.AbilityAdminPermission(AbilityAdminPermissionManage), "Icon")

	code := c.Params("code")
	if len(code) == 0 {
		return ErrIconMissingCode
	}

	ctx := c.Context()

	icon, err := h.iconUsecase.First(ctx, filters.WithFieldEqual("code", code))
	if err != nil {
		return ErrIconDoesntExist
	}

	key := fmt.Sprintf("icons/%s", icon.Code)

	if err := h.uploader.Delete(ctx, key); err != nil {
		return err
	}

	h.iconUsecase.Delete(ctx, models.Icon{}, filters.WithFieldEqual("code", code))

	return c.JSON(201)
}

package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/fsutil"
	"github.com/zsmartex/pkg/infrastucture/uploader"
	"gopkg.in/yaml.v2"

	"github.com/zsmartex/kouda/types"
	"github.com/zsmartex/kouda/usecases"
)

type Handler struct {
	bannerUsecase usecases.BannerUsecase
	uploader      *uploader.Uploader
	abilities     *types.Abilities
}

func NewRouter(
	router fiber.Router,
	banner_usecase usecases.BannerUsecase,
	uploader *uploader.Uploader,
	abilities *types.Abilities,
) {
	bytes := fsutil.MustReadFile("config/abilities.yml")
	yaml.Unmarshal(bytes, &abilities)

	handler := Handler{
		bannerUsecase: banner_usecase,
		uploader:      uploader,
		abilities:     abilities,
	}

	router.Post("/banner", handler.CreateBanner)
	router.Get("/banners", handler.GetBanners)
	router.Patch("/banner/:uuid", handler.UpdateBanner)
}

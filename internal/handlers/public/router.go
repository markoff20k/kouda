package public

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg/v2/infrastucture/uploader"

	"github.com/zsmartex/kouda/usecases"
)

type Handler struct {
	bannerUsecase usecases.BannerUsecase

	Uploader *uploader.Uploader
}

func NewRouter(
	router fiber.Router,
	bannerUsecase usecases.BannerUsecase,
	uploader *uploader.Uploader,
) {

	handler := Handler{
		bannerUsecase: bannerUsecase,
		Uploader:      uploader,
	}

	router.Get("/timestamp", GetTimestamp)

	router.Get("/banners", handler.GetBanners)
	router.Get("/banners/:uuid", handler.GetBannerImage)
}

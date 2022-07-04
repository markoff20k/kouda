package public

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zsmartex/pkg/infrastucture/uploader"

	"github.com/zsmartex/kouda/usecases"
)

type Handler struct {
	bannerUsecase usecases.BannerUsecase

	Uploader *uploader.Uploader
}

func NewRouter(
	router fiber.Router,
	banner_usecase usecases.BannerUsecase,
	uploader *uploader.Uploader,
) {

	handler := Handler{
		bannerUsecase: banner_usecase,
		Uploader:      uploader,
	}

	router.Get("/banners", handler.GetBanners)
}

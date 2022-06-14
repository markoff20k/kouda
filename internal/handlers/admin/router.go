package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/fsutil"
	"gopkg.in/yaml.v2"

	"github.com/zsmartex/kouda/types"
	"github.com/zsmartex/kouda/usecases"
)

type Handler struct {
	bannerUsecase usecases.BannerUsecase
	abilities     *types.Abilities
}

func NewRouter(
	router fiber.Router,
	banner_usecase usecases.BannerUsecase,
	abilities *types.Abilities,
) {
	bytes := fsutil.MustReadFile("config/abilities.yml")
	yaml.Unmarshal(bytes, &abilities)

	// handler := Handler{
	// 	bannerUsecase: banner_usecase,
	// 	abilities:     abilities,
	// }
}

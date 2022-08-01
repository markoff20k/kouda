package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/fsutil"
	"github.com/zsmartex/pkg/v2/infrastucture/uploader"
	"gopkg.in/yaml.v2"

	"github.com/zsmartex/kouda/types"
	"github.com/zsmartex/kouda/usecases"
)

type Handler struct {
	bannerUsecase usecases.BannerUsecase
	uploader      *uploader.Uploader
	abilities     *types.Abilities
}

type Abilities struct {
	Roles            []AbilityRole                                       `yaml:"roles"`
	AdminPermissions map[AbilityRole]map[AbilityAdminPermission][]string `yaml:"admin_permissions"`
}

type AbilityRole string
type AbilityAdminPermission string

const (
	AbilityAdminPermissionRead   AbilityAdminPermission = "read"
	AbilityAdminPermissionManage AbilityAdminPermission = "manage"
)

func NewRouter(
	router fiber.Router,
	bannerUsecase usecases.BannerUsecase,
	uploader *uploader.Uploader,
	abilities *types.Abilities,
) {
	bytes := fsutil.MustReadFile("config/abilities.yml")
	if err := yaml.Unmarshal(bytes, &abilities); err != nil {
		return
	}

	handler := Handler{
		bannerUsecase: bannerUsecase,
		uploader:      uploader,
		abilities:     abilities,
	}

	router.Get("/banners", handler.GetBanners)
	router.Post("/banners", handler.CreateBanner)
	router.Put("/banners", handler.UpdateBanner)
	router.Delete("/banners/:uuid", handler.DeleteBanner)
}

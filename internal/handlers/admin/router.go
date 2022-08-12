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
	iconUsecase   usecases.IconUsecase
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
	iconUsecase usecases.IconUsecase,
	uploader *uploader.Uploader,
	abilities *types.Abilities,
) {
	bytes := fsutil.MustReadFile("config/abilities.yml")
	if err := yaml.Unmarshal(bytes, &abilities); err != nil {
		return
	}

	handler := Handler{
		bannerUsecase: bannerUsecase,
		iconUsecase:   iconUsecase,
		uploader:      uploader,
		abilities:     abilities,
	}

	router.Get("/banners", handler.GetBanners)
	router.Get("/banners/:uuid", handler.GetBannerImage)
	router.Post("/banners", handler.CreateBanner)
	router.Put("/banners", handler.UpdateBanner)
	router.Delete("/banners/:uuid", handler.DeleteBanner)

	router.Get("/icons", handler.GetIcons)
	router.Get("/icons/:code", handler.GetIconImage)
	router.Post("/icons", handler.CreateIcon)
	router.Put("/icons", handler.UpdateIcon)
	router.Delete("/icons/:code", handler.DeleteIcon)
}

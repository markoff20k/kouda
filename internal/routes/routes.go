package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/zsmartex/pkg/infrastucture/uploader"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/config"
	"github.com/zsmartex/kouda/infrastucture/repository"
	"github.com/zsmartex/kouda/internal/handlers/admin"
	"github.com/zsmartex/kouda/internal/handlers/public"
	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/kouda/internal/routes/middlewares"
	"github.com/zsmartex/kouda/internal/routes/middlewares/logger"
	"github.com/zsmartex/kouda/types"
	"github.com/zsmartex/kouda/usecases"
)

func InitializeRoutes(
	db *gorm.DB,
	uploader *uploader.Uploader,
	abilities *types.Abilities,
) *fiber.App {
	config := fiber.Config{
		BodyLimit:               10 * 1024 * 1024, // this is the default limit of 10MB
		EnableTrustedProxyCheck: true,
		ProxyHeader:             "X-Forwarded-For",
		TrustedProxies:          []string{},
		AppName:                 config.Env.ApplicationName,
		ErrorHandler:            middlewares.ErrorHandler,
	}

	app := fiber.New(config)

	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	bannerRepository := repository.New(db, models.Banner{})

	bannerUsecase := usecases.NewBannerUsecase(bannerRepository)

	api_v2 := app.Group("/api/v2")

	public.NewRouter(api_v2.Group("/public"),
		bannerUsecase,
		uploader,
	)

	admin.NewRouter(api_v2.Group("/admin"),
		bannerUsecase,
		uploader,
		abilities,
	)

	return app
}

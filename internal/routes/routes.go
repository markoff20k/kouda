package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/helmet/v2"
	"github.com/zsmartex/pkg/v2/infrastucture/uploader"
	"github.com/zsmartex/pkg/v2/log"
	"gorm.io/gorm"

	"github.com/zsmartex/kouda/config"
	"github.com/zsmartex/kouda/internal/handlers/admin"
	"github.com/zsmartex/kouda/internal/handlers/public"
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

	app.Use(middlewares.ParseIP)
	app.Use(compress.New())
	app.Use(helmet.New())
	app.Use(requestid.New())
	app.Use(logger.New(log.Logger))
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	bannerUsecase := usecases.NewBannerUsecase(db)
	iconUsecase := usecases.NewIconUsecase(db)
	memberUsecase := usecases.NewMemberUsecase(db)

	apiV2 := app.Group("/api/v2")

	public.NewRouter(apiV2.Group("/public"),
		bannerUsecase,
		iconUsecase,
		uploader,
	)

	admin.NewRouter(apiV2.Group("/admin", middlewares.Authorization(memberUsecase)),
		bannerUsecase,
		iconUsecase,
		uploader,
		abilities,
	)

	return app
}

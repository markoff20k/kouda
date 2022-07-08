package main

import (
	"fmt"
	"os"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gookit/validate"
	"github.com/labstack/gommon/log"
	"github.com/mkideal/cli"
	"github.com/zsmartex/pkg/v2/infrastucture/uploader"

	"github.com/zsmartex/kouda/config"
	"github.com/zsmartex/kouda/infrastucture/database"
	"github.com/zsmartex/kouda/internal/models"
	"github.com/zsmartex/kouda/internal/routes"
	"github.com/zsmartex/kouda/migrates"
	"github.com/zsmartex/kouda/types"
)

func InitValidation() {
	validate.AddValidators(map[string]interface{}{
		"bannerState": func(val interface{}) bool {
			state := val.(models.BannerState)

			for _, s := range models.BannerStates {
				if state == s {
					return true
				}
			}

			return false
		},
	})
	validate.Config(func(opt *validate.GlobalOption) {
		opt.SkipOnEmpty = true
	})
	validate.AddGlobalMessages(map[string]string{
		"uint":     "non_positive_{field}",
		"int":      "non_integer_{field}",
		"state":    "invalid_{field}",
		"role":     "invalid_{field}",
		"email":    "invalid_{field}",
		"password": "invalid_{field}",
		"required": "missing_{field}",
	})
}

func main() {
	if err := config.Initialize(); err != nil {
		panic(err)
	}

	if err := cli.Root(root,
		cli.Tree(api),
		cli.Tree(migration),
	).Run(os.Args[1:]); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

// root command
type rootT struct {
	cli.Helper
}

var root = &cli.Command{
	Desc: "this is root command",
	Argv: func() interface{} { return new(rootT) },
}

// child command
type apiArgs struct {
	cli.Helper
	Port int `cli:"p,port" usage:"Kouda api will listen on address running" dft:"8000"`
}

var api = &cli.Command{
	Name: "api",
	Desc: "This command will run kouda api",
	Argv: func() interface{} { return new(apiArgs) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*apiArgs)

		InitValidation()

		database, err := database.NewDatabase(config.Env.DatabaseHost, config.Env.DatabasePort, config.Env.DatabaseUser, config.Env.DatabasePass, config.Env.DatabaseName)
		if err != nil {
			return err
		}

		uploader := uploader.New(config.Env.ObjectStorageBucket, config.Env.ObjectStorageAccessKey, config.Env.ObjectStorageAccessSecret, config.Env.ObjectStorageRegion)
		if err != nil {
			return err
		}

		ablities := &types.Abilities{}

		app := routes.InitializeRoutes(
			database,
			uploader,
			ablities,
		)

		app.Listen(fmt.Sprintf(":%d", argv.Port))

		return nil
	},
}

var migration = &cli.Command{
	Name: "migration",
	Desc: "this is migration command",
	Fn: func(ctx *cli.Context) error {
		database, err := database.NewDatabase(config.Env.DatabaseHost, config.Env.DatabasePort, config.Env.DatabaseUser, config.Env.DatabasePass, config.Env.DatabaseName)
		if err != nil {
			return err
		}

		migrate := gormigrate.New(database, gormigrate.DefaultOptions, migrates.ModelSchemaList)

		return migrate.Migrate()
	},
}

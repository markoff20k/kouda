package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/gookit/validate"
	"github.com/labstack/gommon/log"
	"github.com/mkideal/cli"
	"github.com/zsmartex/pkg/v2/infrastucture/uploader"

	"github.com/zsmartex/pkg/v2/infrastucture/database"

	"github.com/zsmartex/kouda/config"
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
		"iconState": func(val interface{}) bool {
			state := val.(models.IconState)

			for _, s := range models.IconStates {
				if state == s {
					return true
				}
			}

			return false
		},
		"sizeBanner": func(val interface{}) bool {
			str := val.(string)
			if !strings.Contains(str, "x") {
				return false
			}

			stringBeforeX, err := strconv.Atoi(str[:strings.Index(str, "x")])
			if err != nil || stringBeforeX <= 0 {
				return false
			}

			stringAfterX, err := strconv.Atoi(str[strings.Index(str, "x")+1:])
			if err != nil || stringAfterX <= 0 {
				return false
			}

			return true
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

		"bannerState": "invalid_{field}",
		"iconState":   "invalid_{field}",
		"sizeBanner":  "invalid_width_or_height",
	})
}

func main() {
	if err := config.Initialize(); err != nil {
		panic(err)
	}

	if err := cli.Root(root,
		cli.Tree(api),
		cli.Tree(migration),
		cli.Tree(createDB),
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

		db, err := database.New(&database.Config{
			Host:     config.Env.DatabaseHost,
			Port:     config.Env.DatabasePort,
			User:     config.Env.DatabaseUser,
			Password: config.Env.DatabasePass,
			DBName:   config.Env.DatabaseName,
		})
		if err != nil {
			panic(err)
		}

		configUploader := &uploader.Config{
			Bucket:       config.Env.ObjectStorageBucket,
			AccessKey:    config.Env.ObjectStorageAccessKey,
			AccessSecret: config.Env.ObjectStorageAccessSecret,
			Region:       config.Env.ObjectStorageRegion,
			Enpoint:      config.Env.ObjectStorageEnpoint,
			Version:      int64(config.Env.ObjectStorageVersion),
		}

		uploader := uploader.New(configUploader)

		abilities := &types.Abilities{}

		app := routes.InitializeRoutes(
			db,
			uploader,
			abilities,
		)

		if err := app.Listen(fmt.Sprintf(":%d", argv.Port)); err != nil {
			return err
		}

		return nil
	},
}

var migration = &cli.Command{
	Name: "migration",
	Desc: "this is migration command",
	Fn: func(ctx *cli.Context) error {
		db, err := database.New(&database.Config{
			Host:     config.Env.DatabaseHost,
			Port:     config.Env.DatabasePort,
			User:     config.Env.DatabaseUser,
			Password: config.Env.DatabasePass,
			DBName:   config.Env.DatabaseName,
		})
		if err != nil {
			panic(err)
		}

		migrate := gormigrate.New(db, gormigrate.DefaultOptions, migrates.ModelSchemaList)

		return migrate.Migrate()
	},
}

var createDB = &cli.Command{
	Name: "createdb",
	Desc: "this is createdb command",
	Fn: func(ctx *cli.Context) error {
		if err := database.CreateDatabase(config.Env.DatabaseHost, config.Env.DatabasePort, config.Env.DatabaseUser, config.Env.DatabasePass, config.Env.DatabaseName); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		} else {
			log.Infof("Database: %s created successfully", config.Env.DatabaseName)
		}

		return nil
	},
}

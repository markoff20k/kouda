package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/zsmartex/pkg/v2/log"
	"github.com/zsmartex/pkg/v2/validate"

	"github.com/zsmartex/kouda/types"
)

var Env types.ENV

func Initialize() error {
	if err := env.Parse(&Env); err != nil {
		return err
	}

	log.New(Env.ApplicationName)
	validate.InitValidation()

	return nil
}

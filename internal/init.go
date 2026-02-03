package internal

import (
	"go.uber.org/fx"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/data"
	"github.com/go-kratos/kratos-layout/internal/driver"
	"github.com/go-kratos/kratos-layout/internal/service"
)

// ConfigsParams provides the config split
type ConfigsParams struct {
	fx.In

	Bootstrap *conf.Bootstrap
}

// ConfigsResult provides separate server and data configs
type ConfigsResult struct {
	fx.Out

	Server *conf.Server
	Data   *conf.Data
}

// ProvideConfigs splits the bootstrap config into server and data configs
func ProvideConfigs(p ConfigsParams) ConfigsResult {
	return ConfigsResult{
		Server: p.Bootstrap.Server,
		Data:   p.Bootstrap.Data,
	}
}

// Modules returns the complete FX module for the application
func Modules() fx.Option {
	return fx.Options(
		fx.Provide(
			ProvideConfigs,
			data.NewData,
			data.NewGreeterRepo,
			biz.NewGreeterUsecase,
			service.NewGreeterService,
			driver.NewGRPCServer,
			driver.NewHTTPServer,
		),
	)
}

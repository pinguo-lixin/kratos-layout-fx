package internal

import (
	"go.uber.org/fx"

	"github.com/pinguo-lixin/kratos-layout-fx/internal/biz"
	"github.com/pinguo-lixin/kratos-layout-fx/internal/conf"
	"github.com/pinguo-lixin/kratos-layout-fx/internal/data"
	"github.com/pinguo-lixin/kratos-layout-fx/internal/driver"
	"github.com/pinguo-lixin/kratos-layout-fx/internal/service"
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

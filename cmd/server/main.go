package main

import (
	"context"
	"flag"
	"os"

	"github.com/pinguo-lixin/kratos-layout-fx/internal"
	"github.com/pinguo-lixin/kratos-layout-fx/internal/conf"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.uber.org/fx"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	// Initialize the application with FX
	app := fx.New(
		// Load configuration
		fx.Provide(loadConfig),

		// Set up logging
		fx.Provide(setupLogger),

		// Register application modules
		internal.Modules(),

		// Provide the app creator function
		fx.Provide(newApp),

		// startup
		fx.Invoke(func(lifecycle fx.Lifecycle, kapp *kratos.App) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						if err := kapp.Run(); err != nil {
							panic(err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return kapp.Stop()
				},
			})
		}),
	)

	// Start the application
	app.Run()
}

// newApp creates a Kratos application using global variables from main
func newApp(logger log.Logger, httpServer *http.Server, grpcServer *grpc.Server) *kratos.App {
	id, _ := os.Hostname()
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			httpServer,
			grpcServer,
		),
	)
}

// loadConfig loads the application configuration
func loadConfig() *conf.Bootstrap {
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)

	defer func() {
		_ = c.Close()
	}()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	return &bc
}

// setupLogger sets up the application logger
func setupLogger() log.Logger {
	logger := log.NewStdLogger(os.Stdout)
	return log.With(logger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
}

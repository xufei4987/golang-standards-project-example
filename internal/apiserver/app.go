package apiserver

import (
	"golang-standards-project-example/internal/apiserver/config"
	"golang-standards-project-example/internal/apiserver/options"
	"golang-standards-project-example/pkg/app"
)

const commandDesc = `The API server contains a simple http server for study`

/**
IAM API Server：应用的简短描述。
basename：应用的二进制文件名。
opts：应用的命令行选项。
commandDesc：应用的详细描述。
run(opts)：应用的启动函数，初始化应用，并最终启动 HTTP 服务。
*/

func NewApp(basename string) *app.App {
	opts := options.NewOptions()
	application := app.NewApp("user api server",
		basename,
		app.WithOptions(opts),
		app.WithDescription(commandDesc),
		app.WithNoConfig(true),
		app.WithRunFunc(run(opts)),
	)
	return application
}

func run(opts *options.Options) app.RunFunc {
	return func(basename string) error {
		cfg, err := config.CreateConfigFromOptions(opts)
		if err != nil {
			return err
		}
		return Run(cfg)
	}
}

// Run runs the specified APIServer. This should never exit.
func Run(cfg *config.Config) error {
	server, err := NewApiServer(cfg)
	if err != nil {
		return err
	}

	return server.PrepareRun().Run()
}

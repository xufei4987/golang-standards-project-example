package apiserver

import (
	"golang-standards-project-example/internal/apiserver/config"
	"golang-standards-project-example/internal/pkg/server"
	"golang-standards-project-example/pkg/shutdown"
	"golang-standards-project-example/pkg/shutdown/posixsignal"
	"log"
)

type apiServer struct {
	gs                *shutdown.GracefulShutdown
	genericHttpServer *server.GenericHttpServer
	//gRPCAPIServer    *grpcAPIServer
}

type preparedApiServer struct {
	*apiServer
}

func NewApiServer(cfg *config.Config) (*apiServer, error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())
	serverConfig, err := buildApiServerConfig(cfg)
	if err != nil {
		return nil, err
	}
	genericHttpServer, err := serverConfig.Complete().New()
	if err != nil {
		return nil, err
	}
	server := &apiServer{
		gs:                gs,
		genericHttpServer: genericHttpServer,
	}
	return server, nil
}

func buildApiServerConfig(cfg *config.Config) (*server.Config, error) {
	httpConfig := server.NewConfig()
	cfg.HttpServingOptions.ApplyTo(httpConfig)
	return httpConfig, nil
}

func (s *apiServer) PrepareRun() preparedApiServer {
	initRouter(s.genericHttpServer.Engine)
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.genericHttpServer.Close()
		return nil
	}))
	return preparedApiServer{s}
}

func (s preparedApiServer) Run() error {
	//go s.gRPCAPIServer.Run()
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	return s.genericHttpServer.Run()
}

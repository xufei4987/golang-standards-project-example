package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-standards-project-example/internal/pkg/middleware"
	"golang-standards-project-example/pkg/core"
	"golang-standards-project-example/pkg/version"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"strings"
	"time"
)

type GenericHttpServer struct {
	middlewares []string
	// SecureServingInfo holds configuration of the TLS server.
	HttpServingInfo *HttpServingInfo

	// ShutdownTimeout is the timeout used for server shutdown. This specifies the timeout before server
	// gracefully shutdown returns.
	ShutdownTimeout time.Duration

	*gin.Engine
	healthz bool

	httpServer *http.Server
}

func initGenericHttpServer(s *GenericHttpServer) {
	// do some setup
	// s.GET(path, ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.Setup()
	s.InstallMiddlewares()
	s.InstallAPIs()
}

// Setup do some setup work for gin engine.
func (s *GenericHttpServer) Setup() {
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("%-6s %-s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
}

// InstallMiddlewares install generic middlewares.
func (s *GenericHttpServer) InstallMiddlewares() {
	// necessary middlewares
	s.Use(middleware.RequestID())
	s.Use(middleware.Context())

	// install custom middlewares
	for _, m := range s.middlewares {
		mw, ok := middleware.Middlewares[m]
		if !ok {
			log.Printf("can not find middleware: %s\n", m)
			continue
		}

		log.Printf("install middleware: %s\n", m)
		s.Use(mw)
	}
}

// InstallAPIs install generic apis.
func (s *GenericHttpServer) InstallAPIs() {
	// install healthz handler
	if s.healthz {
		s.GET("/healthz", func(c *gin.Context) {
			core.WriteResponse(c, nil, map[string]string{"status": "ok"})
		})
	}

	s.GET("/version", func(c *gin.Context) {
		core.WriteResponse(c, nil, version.Get())
	})
}

func (s *GenericHttpServer) Run() error {
	s.httpServer = &http.Server{
		Addr:    s.HttpServingInfo.Address,
		Handler: s,
	}
	var eg errgroup.Group

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	eg.Go(func() error {
		log.Printf("Start to listening the incoming requests on http address: %s\n", s.HttpServingInfo.Address)

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err.Error())

			return err
		}

		log.Printf("Server on %s stopped\n", s.HttpServingInfo.Address)

		return nil
	})

	// Ping the server to make sure the router is working.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if s.healthz {
		if err := s.ping(ctx); err != nil {
			return err
		}
	}

	if err := eg.Wait(); err != nil {
		log.Fatal(err.Error())
	}

	return nil
}

func (s *GenericHttpServer) Close() {
	// The context is used to inform the server it has 10 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Shutdown http server failed: %s\n", err.Error())
	}
}

// ping pings the http server to make sure the router is working.
func (s *GenericHttpServer) ping(ctx context.Context) error {
	url := fmt.Sprintf("http://%s/healthz", s.HttpServingInfo.Address)
	if strings.Contains(s.HttpServingInfo.Address, "0.0.0.0") {
		url = fmt.Sprintf("http://127.0.0.1:%s/healthz", strings.Split(s.HttpServingInfo.Address, ":")[1])
	}

	for {
		// Change NewRequest to NewRequestWithContext and pass context it
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		// Ping the server by sending a GET request to `/healthz`.

		resp, err := http.DefaultClient.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Printf("The router has been deployed successfully.\n")

			resp.Body.Close()

			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Printf("Waiting for the router, retry in 1 second.\n")
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			log.Fatal("can not ping http server within the specified time interval.")
		default:
		}
	}
	// return fmt.Errorf("the router has no response, or it might took too long to start up")
}

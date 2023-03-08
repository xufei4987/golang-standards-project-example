package apiserver

import (
	"github.com/gin-gonic/gin"
	v1_example "golang-standards-project-example/internal/apiserver/controller/v1/user"
	v2_example "golang-standards-project-example/internal/apiserver/controller/v2/user"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {
	v1 := g.Group("/v1")
	{
		userController := v1_example.NewUserController()
		v1.GET("/hello", userController.Hello)
	}

	v2 := g.Group("/v2")
	{
		userController := v2_example.NewUserController()
		v2.GET("/hello", userController.Hello)
	}
	return g
}

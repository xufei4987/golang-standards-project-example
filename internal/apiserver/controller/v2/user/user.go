package user

import (
	"github.com/gin-gonic/gin"
	"golang-standards-project-example/internal/apiserver/model"
	"golang-standards-project-example/pkg/core"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

func (h *UserController) Hello(c *gin.Context) {
	core.WriteResponse(c, nil, &model.User{
		Nickname: "a2",
		Email:    "a2@email.com",
		Phone:    "13511235123",
	})
}

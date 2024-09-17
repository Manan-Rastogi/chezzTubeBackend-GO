package routers

import (
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/controllers"
	"github.com/gin-gonic/gin"
)

type UserRouter interface {
	InitializeUserRouter()
}

type userRouter struct {
	router *gin.Engine
}

var (
	userController controllers.UserController = controllers.NewUserController()
)

func NewUserRouter(router *gin.Engine) UserRouter {
	return &userRouter{
		router: router,
	}
}

func (ur userRouter) InitializeUserRouter() {
	user := ur.router.Group("api/v1/user")

	// All Endpoints of user
	user.GET("/id/:id", userController.GetUserById)
}

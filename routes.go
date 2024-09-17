package main

import (
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/routers"
	"github.com/gin-gonic/gin"
)

func initializeRoutes(router *gin.Engine) {
	userRouter := routers.NewUserRouter(router)
	userRouter.InitializeUserRouter()
}

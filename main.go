package main

import (
	"runtime/debug"
	"strings"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/db"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	defer db.DisconnectDB()
	defer utils.Logger.Sync()
	utils.Logger.Info("Service Started!!")

	router := gin.New()

	router.Use(gin.Logger(), gin.CustomRecovery(customRecover))

	// initializing routers
	initializeRoutes(router)

	router.Run(configs.ENV.AppPort)
}

func customRecover(ctx *gin.Context, recovered any) {
	if err := recover(); err != nil {
		utils.Logger.Errorf("Panic Occured: %v", err)

		panicError := strings.ReplaceAll(string(debug.Stack()), "/n", "   >>   ")
		panicError = strings.ReplaceAll(panicError, "/t", "  --  ")

		utils.Logger.Errorf("StackTrace: %v", panicError)

	}
}

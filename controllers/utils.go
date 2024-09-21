package controllers

import (
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/gin-gonic/gin"
)

func respondErr(ctx *gin.Context, httpcode, code int) {
	ctx.JSON(httpcode, models.Response{
		ErrorResponse: models.ErrorResponse{
			Code:    code,
			Message: configs.ServiceCodes[code],
		},
	})
}

func respond(ctx *gin.Context, data interface{}, httpcode, code int) {
	ctx.JSON(httpcode, models.Response{
		SuccessResponse: models.SuccessResponse{
			Code: code,
			Data: data,
		},
	})
}

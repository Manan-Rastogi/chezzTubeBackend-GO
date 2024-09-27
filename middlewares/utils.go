package middlewares

import (
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/gin-gonic/gin"
)

func respondErr(ctx *gin.Context, httpcode, code int) {
	ctx.AbortWithStatusJSON(httpcode, models.Response{
		ErrorResponse: &models.ErrorResponse{
			Code:    code,
			Message: configs.ServiceCodes[code],
		},
		SuccessResponse: nil,
	})
}

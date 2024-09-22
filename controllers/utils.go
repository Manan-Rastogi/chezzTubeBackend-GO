package controllers

import (
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func respondErr(ctx *gin.Context, httpcode, code int) {
	ctx.JSON(httpcode, models.Response{
		ErrorResponse: &models.ErrorResponse{
			Code:    code,
			Message: configs.ServiceCodes[code],
		},
		SuccessResponse: nil,
	})
}

func respond(ctx *gin.Context, data interface{}, httpcode, code int) {
	ctx.JSON(httpcode, models.Response{
		ErrorResponse: nil,
		SuccessResponse: &models.SuccessResponse{
			Code: code,
			Data: data,
		},
	})
}

func encryptPassword(password string) (encPass []byte, err error){
	encPass, err = bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil{
		utils.Logger.Error(err.Error())
	}
	return
}
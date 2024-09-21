package controllers

import (
	"net/http"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	GetUserById(ctx *gin.Context)
}

type userController struct {
}

func NewUserController() UserController {
	return &userController{}
}

func (uc *userController) RegisterUser(ctx *gin.Context) {
	input, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			ErrorResponse: models.ErrorResponse{
				Code:    1001,
				Message: configs.ServiceCodes[1001],
			},
		})
		return
	}

	data := input.Value
	// files := input.File
	var userInput = models.Users{}

	
	if utils.IsFormKeyPresent("username", data) {
		userInput.UserName = data["username"][0]
		if userInput.UserName == "" {
			ctx.JSON(http.StatusBadRequest, models.Response{
				ErrorResponse: models.ErrorResponse{
					Code:    1002,
					Message: configs.ServiceCodes[1002],
				},
			})
			return
		}
	} else {
		ctx.JSON(http.StatusBadRequest, models.Response{
			ErrorResponse: models.ErrorResponse{
				Code:    1003,
				Message: configs.ServiceCodes[1003],
			},
		})
		return
	}

	// Check in DB for unique username


}

func (uc *userController) GetUserById(ctx *gin.Context) {

}

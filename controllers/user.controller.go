package controllers

import "github.com/gin-gonic/gin"

type UserController interface {
	GetUserById(ctx *gin.Context)
}

type userController struct {
}

func NewUserController() UserController {
	return &userController{}
}

func (uc *userController) GetUserById(ctx *gin.Context) {
	
}

package controllers

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/validators"
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
	// 0. No Login Checks
	// 1. get data
	// 2. validate data
	// 2b. chk user/email already exists
	// 3. routine upload for cloudinary
	// 4. prepare data for mongo upload
	// 5. return success

	input, err := ctx.MultipartForm()
	if err != nil {
		respondErr(ctx, http.StatusBadRequest, 1001)
		return
	}

	data := input.Value
	// files := input.File
	var userInput = models.Users{}

	if utils.IsFormKeyPresent("username", data) {
		userInput.UserName = data["username"][0]
		if userInput.UserName == "" {
			respondErr(ctx, http.StatusBadRequest, 1002)
			return
		}
	} else {
		respondErr(ctx, http.StatusBadRequest, 1003)
		return
	}

	if utils.IsFormKeyPresent("email", data) {
		userInput.Email = data["email"][0]
		if userInput.Email == "" {
			respondErr(ctx, http.StatusBadRequest, 1004)
			return
		} else if !validators.IsValidEmail(userInput.Email) {
			respondErr(ctx, http.StatusBadRequest, 1005)
			return
		}
	} else {
		respondErr(ctx, http.StatusBadRequest, 1006)
		return
	}

	// email username exist check
	userEmailExist := models.UserEmailCheck{}
	userEmailExistChannel := make(chan models.UserEmailCheck)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go validators.CheckUsernameAndEmailExists(userInput.UserName, userInput.Email, &wg, ctx, userEmailExistChannel)

	/////////////////////////////////////////// Do remaining validations

	wg.Wait()

	userEmailExist = <-userEmailExistChannel

	if userEmailExist.Err != nil {
		if os.IsTimeout(userEmailExist.Err) {
			respondErr(ctx, http.StatusInternalServerError, 1007)
			return
		}else{
			respondErr(ctx, http.StatusServiceUnavailable, 5000)
			return
		}
	} else if strings.EqualFold(userEmailExist.UserData.UserName, userInput.UserName) {
		respondErr(ctx, http.StatusConflict, 1008)
		return
	} else if strings.EqualFold(userEmailExist.UserData.Email, userInput.Email) {
		respondErr(ctx, http.StatusConflict, 1009)
	}
}

func (uc *userController) GetUserById(ctx *gin.Context) {

}

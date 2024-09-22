package controllers

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/cloudinary"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/validators"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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
	files := input.File
	var userInput = models.Users{}

	if utils.IsFormKeyPresent("username", data) {
		userInput.UserName = strings.ToLower(data["username"][0])
		if userInput.UserName == "" {
			respondErr(ctx, http.StatusBadRequest, 1002)
			return
		}
	} else {
		respondErr(ctx, http.StatusBadRequest, 1003)
		return
	}

	if utils.IsFormKeyPresent("email", data) {
		userInput.Email = strings.ToLower(data["email"][0])
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
	defer close(userEmailExistChannel)
	wg := sync.WaitGroup{}

	wg.Add(1)
	go validators.CheckUsernameAndEmailExists(userInput.UserName, userInput.Email, &wg, userEmailExistChannel)

	// Fullname is a non essential field
	if utils.IsFormKeyPresent("fullName", data) {
		if len(data["fullName"][0]) > 64 {
			respondErr(ctx, http.StatusBadRequest, 1019)
			return
		}
		userInput.FullName = data["fullName"][0]
	}

	wg.Wait()

	userEmailExist = <-userEmailExistChannel

	if userEmailExist.Err != nil {
		if os.IsTimeout(userEmailExist.Err) {
			respondErr(ctx, http.StatusInternalServerError, 1007)
			return
		} else if userEmailExist.Err != mongo.ErrNoDocuments {
			respondErr(ctx, http.StatusServiceUnavailable, 5000)
			return
		}
	} else if strings.EqualFold(userEmailExist.UserData.UserName, userInput.UserName) {
		respondErr(ctx, http.StatusConflict, 1008)
		return
	} else if strings.EqualFold(userEmailExist.UserData.Email, userInput.Email) {
		respondErr(ctx, http.StatusConflict, 1009)
	}

	// User is Validated.. we can register the user now - uploadimages , create tokens and update DB

	avatarChan := make(chan models.ImageUploadChan)
	defer close(avatarChan)
	if utils.IsFormFileKeyPresent("avatar", files) {
		avatar := files["avatar"][0]

		if !configs.AllowedImagesExt[utils.FileExtension(avatar.Filename)] {
			respondErr(ctx, http.StatusBadRequest, 1010)
			return
		} else if !utils.IsImageFile(avatar.Filename) {
			respondErr(ctx, http.StatusBadRequest, 1011)
			return
		} else if avatar.Size < 0 && avatar.Size > configs.ENV.AvatarMaxSize {
			respondErr(ctx, http.StatusBadRequest, 1012)
			return
		}

		// as the file is present we can now upload the file to cloud and get a url in return
		avatarFile, err := avatar.Open()
		if err != nil {
			respondErr(ctx, http.StatusBadRequest, 1013)
			return
		}
		wg.Add(1)
		go cloudinary.UploadImage(cloudinary.CLOUDINARY, &wg, 30*time.Second, avatarFile, userInput.UserName, avatarChan)

	}
	// avatar is non essential

	coverImageChan := make(chan models.ImageUploadChan)
	defer close(coverImageChan)
	if utils.IsFormFileKeyPresent("coverImage", files) {
		coverImage := files["coverImage"][0]

		if !configs.AllowedImagesExt[utils.FileExtension(coverImage.Filename)] {
			respondErr(ctx, http.StatusBadRequest, 1014)
			return
		} else if !utils.IsImageFile(coverImage.Filename) {
			respondErr(ctx, http.StatusBadRequest, 1015)
			return
		} else if coverImage.Size < 0 || coverImage.Size > configs.ENV.CoverImageMaxSize {
			respondErr(ctx, http.StatusBadRequest, 1016)
			return
		}

		// as the file is present we can now upload the file to cloud and get a url in return
		coverImageFile, err := coverImage.Open()
		if err != nil {
			respondErr(ctx, http.StatusBadRequest, 1017)
			return
		}
		wg.Add(1)
		go cloudinary.UploadImage(cloudinary.CLOUDINARY, &wg, 30*time.Second, coverImageFile, userInput.UserName, coverImageChan)
	} else {
		respondErr(ctx, http.StatusBadRequest, 1018)
		return
	}

	// Create Access and Refresh Token in Meantime
	

}

func (uc *userController) GetUserById(ctx *gin.Context) {

}

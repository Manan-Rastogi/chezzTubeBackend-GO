package controllers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/auth"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/cloudinary"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/db"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/validators"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	///////////////////////////////// Password Check and Encrypt Password.

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

	// with urls of images register user
	wg.Done()
	avatarUpload := <-avatarChan
	coverImageUpload := <-coverImageChan

	if utils.IsFormFileKeyPresent("avatar", files) {
		if avatarUpload.Err != nil {
			respondErr(ctx, http.StatusInternalServerError, 1020)
			return
		}
		userInput.Avatar = avatarUpload.SecureUrl
	}

	if coverImageUpload.Err != nil {
		respondErr(ctx, http.StatusInternalServerError, 1021)
		return
	} else {
		userInput.CoverImage = &coverImageUpload.SecureUrl
	}

	// now send to DB
	userInput.CreatedAt = time.Now().Local()
	userCollection := db.Client.Database(configs.DB_NAME).Collection("users")

	userResult, err := userCollection.InsertOne(context.Background(), userInput)
	if err != nil {
		utils.Logger.Error(err.Error())
		respondErr(ctx, 500, 2001)
		return
	}

	userInput.Id, _ = userResult.InsertedID.(primitive.ObjectID) // assuming entry is done. else handle another error

	// Create Access and Refresh Token in after registering user
	accessTokenChan := make(chan auth.JwtToken)
	refreshTokenChan := make(chan auth.JwtToken)

	wg.Add(1)
	accessContext, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	go auth.CreateNewToken(accessContext, userInput.Id.String(), userInput.UserName, userInput.Email, 24*time.Hour, "access", accessTokenChan, &wg)

	wg.Add(1)
	refreshContext, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	go auth.CreateNewToken(refreshContext, userInput.Id.String(), userInput.UserName, userInput.Email, 7*24*time.Hour, "refresh", refreshTokenChan, &wg)

	wg.Done()
	accessToken := <-accessTokenChan
	refreshToken := <-refreshTokenChan

	if accessToken.Err != nil || refreshToken.Err != nil {
		// Delete user from DB
		_, err := userCollection.DeleteOne(context.Background(), bson.D{
			{Key: "_id", Value: userInput.Id},
		})
		if err != nil {
			// Handle the error, e.g., trigger a mail to admin or queue the delete operation
			// You can also log the error or take appropriate action based on your application's requirements
			utils.Logger.Error(err.Error())
		}
		respondErr(ctx, http.StatusInternalServerError, 1022)
		return
	}

	// Update Refresh Token Channel in DB
	userInput.RefreshToken = &refreshToken.Token
	userInput.UpdatedAt = time.Now().Local()
	// Assuming userCollection is an instance of a collection from the database
	filter := bson.M{"_id": userInput.Id}
	update := bson.M{
		"$set": bson.M{
			"refreshToken": userInput.RefreshToken,
			"updatedAt":    userInput.UpdatedAt,
		},
	}
	_, err = userCollection.UpdateByID(context.Background(), filter, update)
	if err != nil {
		// Delete user from DB
		_, err := userCollection.DeleteOne(context.Background(), bson.D{
			{Key: "_id", Value: userInput.Id},
		})
		if err != nil {
			// Handle the error, e.g., trigger a mail to admin or queue the delete operation
			// You can also log the error or take appropriate action based on your application's requirements
			utils.Logger.Error(err.Error())
		}
		respondErr(ctx, http.StatusInternalServerError, 1022)
		return
	}

	// Set Cookies
	// Set Bearer token in Cookies
	ctx.SetCookie("accessToken", "Bearer "+accessToken.Token, -1, "", "", true, true)
	ctx.SetCookie("refreshToken", "Bearer "+refreshToken.Token, -1, "", "", true, true)

	// Resonse to User.
	response := map[string]interface{}{
		"user":         userInput,
		"msg":          "user registered successfully.",
		"accessToken":  "Bearer " + accessToken.Token,  // Tokens are sned here as well for mobile applns as they cannot access cookies.
		"refreshToken": "Bearer " + refreshToken.Token,
	}
	respond(ctx, response, http.StatusCreated, 1)

}


// Login API


func (uc *userController) GetUserById(ctx *gin.Context) {

}

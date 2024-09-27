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
	"golang.org/x/crypto/bcrypt"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	Login(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	UpdateAvatarOrCoverImage(ctx *gin.Context)
	GetUserById(ctx *gin.Context)

	Logout(ctx *gin.Context)
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

	utils.Logger.Info("User Registration begins.")
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

	wg := sync.WaitGroup{}

	wg.Add(1)
	go validators.CheckUsernameAndEmailExists(userInput.UserName, userInput.Email, &wg, userEmailExistChannel)

	///////////////////////////////// Password Check and Encrypt Password.
	if utils.IsFormKeyPresent("password", data) {
		if !validators.ValidateUserPassword(data["password"][0]) {
			respondErr(ctx, http.StatusBadRequest, 1024)
			return
		}
	} else {
		respondErr(ctx, http.StatusBadRequest, 1023)
		return
	}

	encPass, err := encryptPassword(data["password"][0])
	if err != nil {
		respondErr(ctx, http.StatusInternalServerError, 1025)
		return
	}

	userInput.Password = string(encPass)

	// Fullname is a non essential field
	if utils.IsFormKeyPresent("fullName", data) {
		if len(data["fullName"][0]) > 64 {
			respondErr(ctx, http.StatusBadRequest, 1019)
			return
		}
		userInput.FullName = data["fullName"][0]
	}

	userEmailExist = <-userEmailExistChannel
	wg.Wait()

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
		return
	}

	// User is Validated.. we can register the user now - uploadimages , create tokens and update DB

	avatarChan := make(chan models.ImageUploadChan)

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
		go cloudinary.UploadImage(cloudinary.CLOUDINARY, &wg, 30*time.Second, avatarFile, userInput.UserName+"_avatar", avatarChan)
	} else {
		go func() { // need to send something to channel to avoid deadlock
			avatarChan <- models.ImageUploadChan{
				SecureUrl: "",
				Err:       nil,
			}

			close(avatarChan)
		}()
	}
	// avatar is non essential

	coverImageChan := make(chan models.ImageUploadChan)

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
		go cloudinary.UploadImage(cloudinary.CLOUDINARY, &wg, 30*time.Second, coverImageFile, userInput.UserName+"_cover", coverImageChan)
	} else {
		respondErr(ctx, http.StatusBadRequest, 1018)
		return
	}

	// with urls of images register user
	avatarUpload := <-avatarChan
	coverImageUpload := <-coverImageChan
	wg.Wait()

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

	// Resonse to User.
	response := map[string]interface{}{
		"user": userInput,
		"msg":  "user registered successfully.",
	}
	respond(ctx, response, http.StatusCreated, 1)

}

// Login API

func (c *userController) Login(ctx *gin.Context) {
	// 1. Validate username/email + passowrd
	// 2. On Success - create access and refresh token
	// 3. Set refresh token in DB
	// 4. Set Cookies and return
	// Cases - already logged in? access token expire? refresh token expire?

	user := models.LoginUser{}

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		utils.Logger.Error(err.Error())
		respondErr(ctx, http.StatusBadRequest, 1026)
		return
	}

	if user.UserName == "" || user.Password == "" {
		respondErr(ctx, http.StatusBadRequest, 1027)
		return
	}

	user.UserName = strings.ToLower(user.UserName)

	// verify valid user.. then verify password
	userCollection := db.Client.Database(configs.DB_NAME).Collection("users")

	contextMongoSearch, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	userResult := userCollection.FindOne(contextMongoSearch, bson.D{{Key: "username", Value: user.UserName}})

	if userResult.Err() != nil {
		utils.Logger.Error(userResult.Err().Error())
		if userResult.Err() == mongo.ErrNoDocuments {
			respondErr(ctx, http.StatusNotFound, 1028)
			return
		} else if os.IsTimeout(userResult.Err()) {
			respondErr(ctx, http.StatusGatewayTimeout, 5001)
			return
		} else {
			respondErr(ctx, http.StatusInternalServerError, 5000)
			return
		}
	}

	userData := models.Users{}
	err = userResult.Decode(&userData)
	if err != nil {
		utils.Logger.Error(err.Error())
		respondErr(ctx, http.StatusInternalServerError, 5000)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password))
	if err != nil {
		utils.Logger.Error(err.Error())
		respondErr(ctx, http.StatusUnauthorized, 1029)
		return
	}

	// Create Access and Refresh Token in after registering user
	wg := sync.WaitGroup{}
	accessTokenChan := make(chan auth.JwtToken)
	refreshTokenChan := make(chan auth.JwtToken)

	wg.Add(1)
	accessContext, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	go auth.CreateNewToken(accessContext, userData.Id.String(), userData.UserName, userData.Email, 24*time.Hour, "access", accessTokenChan, &wg)

	wg.Add(1)
	refreshContext, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	go auth.CreateNewToken(refreshContext, userData.Id.String(), userData.UserName, userData.Email, 7*24*time.Hour, "refresh", refreshTokenChan, &wg)

	accessToken := <-accessTokenChan
	refreshToken := <-refreshTokenChan
	wg.Wait()

	if accessToken.Err != nil || refreshToken.Err != nil {
		utils.Logger.Error(accessToken.Err)
		utils.Logger.Error(accessToken.Err)
		respondErr(ctx, http.StatusInternalServerError, 1022)
		return
	}

	// Update Refresh Token Channel in DB
	userData.RefreshToken = &refreshToken.Token
	userData.UpdatedAt = time.Now().Local()
	// Assuming userCollection is an instance of a collection from the database
	filter := bson.M{"_id": userData.Id}
	update := bson.M{
		"$set": bson.M{
			"refreshToken": userData.RefreshToken,
			"updatedAt":    userData.UpdatedAt,
		},
	}
	_, err = userCollection.UpdateByID(context.Background(), filter, update)
	if err != nil {
		respondErr(ctx, http.StatusInternalServerError, 1022)
		return
	}

	// Set Cookies
	// Set Bearer token in Cookies
	ctx.SetCookie("accessToken", "Bearer "+accessToken.Token, -1, "", "", true, true)
	ctx.SetCookie("refreshToken", "Bearer "+refreshToken.Token, -1, "", "", true, true)

	response := map[string]interface{}{
		"user":         userData,
		"msg":          "log in successful!",
		"accessToken":  "Bearer " + accessToken.Token,
		"refreshToken": "Bearer " + refreshToken.Token,
	}
	respond(ctx, response, http.StatusOK, 1)
}

func (c *userController) ChangePassword(ctx *gin.Context) {
	// 1. User is already loggedin - check via middleware [DONE]
	// 2. Email oldPassword NewPassword ConfirmNewPassword inps required
	// 3. check newPass == confirmPass
	// 3b. Check in DB - email and oldPassword
	// 4. In meantime we can match newPass && confirmNewPass && hash it
	// 5. Password is retrieved, now match password first. If success
	// 6. Update NewPassword in DB and return success
	// NOTE: 4th point can be a waste of resources if old password didn't match but will save us time in othercases. So this decision is dependent on us, what we want to use. 
}

func (c *userController) UpdateAvatarOrCoverImage(ctx *gin.Context){
	// 1. User is already loggedin - check via middleware [DONE]
	// 2. which files are present - cover or avatar or both
	// 3. get user data from DB.
	// 4. In meantime upload files to cloudinary
	// 5. Update the files urls in DB
	// 6. In meantime Delete oldurls from cloudinary.
	// 7. Return Updated user Data
	// NOTE: Point 6 is optional as some companies might want to keep data for ML/AI training Or user want to use old image again [altough we are not storing thee in arrays and our naming convention use username as publicId]. This is also a decision company will take instead of individual developer.
}

func (uc *userController) GetUserById(ctx *gin.Context) {

}


func (c *userController) Logout(ctx *gin.Context) {
	// 1. User is already loggedin - check via middleware [DONE]
	// 2. unset cookies
	// 3. update expiretime of accessToken && refreshToken to currTime
	// 4. unset refreshToken from db
	// 5. return logout success with No userdetails
}
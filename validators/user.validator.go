package validators

import (
	"context"
	"net/mail"
	"sync"
	"time"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/db"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// IsValidEmail verifies if the provided email address is valid.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		utils.Logger.Error(err.Error())
		return false
	}

	return true
}

// Check in DB for unique email  AND  Check in DB for unique username
func CheckUsernameAndEmailExists(username, email string, wg *sync.WaitGroup, ctx *gin.Context, output chan models.UserEmailCheck) {
	defer wg.Done()

	userCollection := db.Client.Database(configs.DB_NAME).Collection("users")
	
	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "username", Value: username}},
			bson.D{{Key: "email", Value: email}},
		}},
	}

	ctxx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	var userExists models.Users
	err := userCollection.FindOne(ctxx, filter).Decode(&userExists)
	if err != nil {
		utils.Logger.Error(err.Error())
	}

	output <- models.UserEmailCheck{
		UserData: userExists,
		Err: err,
	}
}

package validators

import (
	"context"
	"net/mail"
	"sync"
	"time"
	"unicode"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/db"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/models"
	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
func CheckUsernameAndEmailExists(username, email string, wg *sync.WaitGroup, output chan models.UserEmailCheck) {
	defer wg.Done()
	userCollection := db.Client.Database(configs.DB_NAME).Collection("users")

	filter := bson.D{
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "username", Value: username}},
			bson.D{{Key: "email", Value: email}},
		}},
	}

	ctxx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var userExists models.Users
	err := userCollection.FindOne(ctxx, filter).Decode(&userExists)
	if err != nil && err != mongo.ErrNoDocuments {
		utils.Logger.Error(err.Error())
	}

	output <- models.UserEmailCheck{
		UserData: userExists,
		Err:      err,
	}
	close(output)
}

// Password must have a spl character, 1 lower, 1 uppercase and must be greater than 8 characters.
func ValidateUserPassword(password string) bool {
	if len(password) < 8 || len(password) > 25 {
		return false
	}

	var hasUpper, hasLower, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasSpecial
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	Id           primitive.ObjectID   `bson:"_id" json:"id"`
	UserName     string               `bson:"username" json:"username"`
	WatchHistory []primitive.ObjectID `bson:"watchHistory" json:"watchHistory"`
	Email        string               `bson:"email" json:"email"`
	FullName     string               `bson:"fullName" json:"fullName"`
	Avatar       string               `bson:"avatar" json:"avatar" `
	CoverImage   *string              `bson:"coverImage" json:"coverImage"`
	Password     string               `bson:"password" json:"-"`
	RefreshToken *string              `bson:"refreshToken" json:"-"`
	CreatedAt    *time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt    *time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type UserEmailCheck struct {
	UserData Users
	Err      error
}

type ImageUploadChan struct {
	SecureUrl string
	Err       error
}

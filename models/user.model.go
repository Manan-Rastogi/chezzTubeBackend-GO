package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Users struct {
	Id           primitive.ObjectID   `bson:"_id" json:"id"`
	UserName     string               `bson:"username" json:"username" binding:"required,min=3,max=100"`
	WatchHistory []primitive.ObjectID `bson:"watchHistory" json:"watchHistory"`
	Email        string               `bson:"email" json:"email" binding:"required,email"`
	FullName     string               `bson:"fullName" json:"fullName" binding:"required"`
	Avatar       string               `bson:"avatar" json:"avatar" `
	CoverImage   *string              `bson:"coverImage" json:"coverImage"`
	Password     string               `bson:"password" json:"-" binding:"required"`
	RefreshToken *string              `bson:"refreshToken" json:"-"`
	CreatedAt    *time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt    *time.Time           `bson:"updatedAt" json:"updatedAt"`
}


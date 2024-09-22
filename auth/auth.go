package auth

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/configs"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id       string `json:"_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type JwtToken struct {
	Token string
	Err   error
}

func CreateNewToken(ctx context.Context, id, username, email string, expirationTime time.Duration, subject string, tokenChan chan JwtToken, wg *sync.WaitGroup) {
	defer wg.Done()

	claims := Claims{
		id, username, email,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "admin",
			Subject:   subject,
			ID:        id,
			Audience:  []string{username},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var ss string
	var err error
	select {
	case <-ctx.Done():
		// Context canceled or timed out
		tokenChan <- JwtToken{
			Token: "",
			Err:   ctx.Err(),
		}
		return
	default:
		if strings.EqualFold(subject, "refresh") {
			ss, err = token.SignedString(configs.ENV.JwtRefreshKey)
		} else {
			ss, err = token.SignedString(configs.ENV.JwtAccessKey)
		}
	}

	tokenChan <- JwtToken{
		Token: ss,
		Err:   err,
	}
}

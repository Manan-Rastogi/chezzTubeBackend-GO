package middlewares

import (
	"net/http"
	"strings"

	"github.com/Manan-Rastogi/chezzTubeBackend-GO/utils"
	"github.com/gin-gonic/gin"
)

func ValidateSession() gin.HandlerFunc {
	// 1. Read cookies
	// 2. check validity of access and refresh token
	//2b. if accesstoken has expired but refresh token hasn't we need to replace access token after 4th step
	// 3. From _id set in claims in jwt, we read user from DB
	// 4. compare refresh token
	// 5. user validated. set user details in ctx and pass it to next middleware.
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie("accessToken")
		if err != nil {
			utils.Logger.Error(err.Error())
			respondErr(ctx, http.StatusUnauthorized, 1030)
			return
		}

		refreshToken, err := ctx.Cookie("refreshToken")
		if err != nil {
			utils.Logger.Error(err.Error())
			respondErr(ctx, http.StatusUnauthorized, 1031)
			return
		}

		if strings.Contains(accessToken, "Bearer ") {
			accessTokensParts := strings.Split(accessToken, "Bearer ")
			if len(accessTokensParts) > 1 {
				accessToken = accessTokensParts[1]
			} else {
				utils.Logger.Error("accessTokenParts in blank.")
				respondErr(ctx, http.StatusUnauthorized, 1030)
				return
			}
		} else {
			utils.Logger.Error("not a Bearer token.")
			respondErr(ctx, http.StatusUnauthorized, 1030)
			return
		}

		if strings.Contains(refreshToken, "Bearer ") {
			refreshTokensParts := strings.Split(refreshToken, "Bearer ")
			if len(refreshTokensParts) > 1 {
				refreshToken = refreshTokensParts[1]
			} else {
				utils.Logger.Error("refreshTokenParts in blank.")
				respondErr(ctx, http.StatusUnauthorized, 1030)
				return
			}
		} else {
			utils.Logger.Error("not a Bearer token.")
			respondErr(ctx, http.StatusUnauthorized, 1030)
			return
		}

		// 2

		

	}
}

package middleware

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"manga-bookmarker-backend/services"
	"strings"
)

func Auth(ctx iris.Context) {
	// Extract the Authorization header
	authHeader := ctx.GetHeader("Authorization")
	// Extract the token by trimming the "Bearer " prefix
	token := strings.TrimPrefix(authHeader, "Bearer ")

	//User ID
	userId, err, code := services.GetUserIdFromClaims(token)
	if err != nil {
		fmt.Println("Error obtaining claims: ", err)
		if code == 1 {
			ctx.StatusCode(iris.StatusUnauthorized)
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
		}
		return
	}

	validUser, err := services.ValidUser(userId)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": err.Error()})
		return
	}

	if !validUser {
		ctx.StatusCode(iris.StatusUnauthorized)
		return
	}

	ctx.Values().Set("userId", userId)

	ctx.Next()
}

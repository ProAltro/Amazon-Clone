package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProAltro/Amazon-Clone/http"
	"github.com/gin-gonic/gin"
)

func AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			ctx.JSON(403, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		email, uid, err := http.GetSession(sessionID)
		if err != nil {
			fmt.Println("error", err)
			ctx.JSON(403, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Set("email", email)
		ctx.Set("uid", uid)
		ctx.Set("session_id", sessionID)
		ctx.Next()
	}
}

func AuthenticateAdmin() gin.HandlerFunc {

	adminUsers := os.Getenv("ADMIN_USERS") //comma separated list of admin users
	adminUsersList := strings.Split(adminUsers, ",")

	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			ctx.JSON(403, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		email, uid, err := http.GetSession(sessionID)
		if err != nil {
			fmt.Println("error", err)
			ctx.JSON(403, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		//check if user is admin
		isAdmin := false
		for _, adminUser := range adminUsersList {
			if adminUser == email {
				isAdmin = true
			}
		}
		if !isAdmin {
			ctx.JSON(403, gin.H{
				"error": "not authorised",
			})
			ctx.Abort()
			return
		}
		ctx.Set("email", email)
		ctx.Set("uid", uid)
		ctx.Set("session_id", sessionID)
		ctx.Next()
	}
}

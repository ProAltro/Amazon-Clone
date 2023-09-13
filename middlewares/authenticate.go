package middlewares

import (
	"fmt"

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
		email, err := http.GetSession(sessionID)
		if err != nil {
			fmt.Println("error", err)
			ctx.JSON(403, gin.H{
				"error": err.Error(),
			})
			ctx.Abort()
			return
		}
		ctx.Set("email", email)
		ctx.Next()
	}
}

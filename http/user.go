package http

import (
	"github.com/ProAltro/Amazon-Clone/entity"
	"github.com/gin-gonic/gin"
)

func (http HTTPService) UserSignup(ctx *gin.Context) {
	var user entity.User
	userService := http.UserService
	//err := json.NewDecoder(ctx.Request.Body).Decode(&u)
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = userService.CreateUser(&user)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "user created successfully",
	})
}

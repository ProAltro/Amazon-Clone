package http

import (
	"strconv"
	"time"

	"github.com/ProAltro/Amazon-Clone/entity"
	"github.com/gin-gonic/gin"
)

func (http HTTPService) UserSignup(ctx *gin.Context) {
	var user entity.User
	userService := http.UserService
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

func (http HTTPService) UserLogin(ctx *gin.Context) {
	var user entity.User
	userService := http.UserService
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = userService.AuthenticateUser(user.Email, user.Password)
	if err != nil {
		ctx.JSON(403, gin.H{
			"error": err.Error(),
		})
		return
	}
	//create session
	sessionID, err := CreateSession(user.Email, time.Now().Add(24*time.Hour))
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.SetCookie("session_id", sessionID, 3600, "/", "localhost", false, true)
	ctx.JSON(200, gin.H{
		"message": "user logged in successfully",
	})
}

func (http HTTPService) FetchUser(ctx *gin.Context) {

	userService := http.UserService
	id, err := strconv.Atoi(ctx.Query("id"))
	email := ctx.Query("email")
	if email == "" && err != nil {
		ctx.JSON(400, gin.H{
			"error": "empty id",
		})
		return
	}
	if email != "" {
		user, err := userService.FindUserByEmail(email)
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"user": user,
		})
		return
	} else if id != 0 {
		user, err := userService.FindUserByID(id)
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"user": user,
		})

	}
}

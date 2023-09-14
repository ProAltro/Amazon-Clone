package http

import (
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
	_, err = userService.CreateUser(ctx, user.Name, user.Email, user.Password)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "user created successfully",
	})
}

func (http HTTPService) UserLogin(ctx *gin.Context) {
	type resp struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var u resp
	userService := http.UserService
	err := ctx.BindJSON(&u)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	user, err := userService.AuthenticateUser(ctx, u.Email, u.Password)

	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	sessionID, err := CreateSession(user.Id, user.Email, time.Now().Add(time.Hour))
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

func (http HTTPService) GetUser(ctx *gin.Context) {

	userService := http.UserService
	id := ctx.Value("uid").(int)
	email := ctx.Value("email").(string)
	if id == -1 || email == "" {
		ctx.JSON(403, gin.H{
			"error": "user not logged in",
		})
		return
	}

	var user *entity.User
	var err error
	if email != "" {
		user, err = userService.GetUserByEmail(ctx, email)
	} else {
		user, err = userService.GetUserByID(ctx, id)
	}
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"user": user,
	})
}

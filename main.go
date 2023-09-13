package main

import (
	"log"

	"github.com/ProAltro/Amazon-Clone/http"
	"github.com/ProAltro/Amazon-Clone/middlewares"
	"github.com/ProAltro/Amazon-Clone/mysql"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := mysql.NewDB()
	err = db.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	httpServ := http.NewHTTPService(mysql.NewUserService(db))

	router := gin.Default()
	superGroup := router.Group("/api/v1")
	{
		userGroup := superGroup.Group("/user")
		{
			userGroup.POST("/signup", httpServ.UserSignup)
			userGroup.POST("/login", httpServ.UserLogin)
		}
	}
	authorisedGroup := router.Group("/api/v1")
	authorisedGroup.Use(middlewares.AuthenticateUser())
	{
		authorisedGroup.GET("/get", httpServ.FetchUser)
	}

	router.Run(":8080")
}

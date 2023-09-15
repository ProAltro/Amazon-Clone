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
	httpServ := http.NewHTTPService(mysql.NewUserService(db), mysql.NewProductService(db), mysql.NewCartService(db), mysql.NewOrderService(db), mysql.NewInventoryService(db))
	router := gin.Default()
	superGroup := router.Group("/api/v1")
	{
		superGroup.POST("/signup", httpServ.UserSignup)
		superGroup.POST("/login", httpServ.UserLogin)
		authorisedGroup := superGroup.Group("/")
		authorisedGroup.Use(middlewares.AuthenticateUser())
		{
			authorisedGroup.GET("/user", httpServ.GetUser)
			cartGroup := authorisedGroup.Group("/cart")
			{
				cartGroup.GET("/", httpServ.GetCart)
				cartGroup.POST("/add", httpServ.AddProductToCart)
				cartGroup.POST("/remove", httpServ.RemoveProductFromCart)
				cartGroup.POST("/modify", httpServ.ModifyCart)
				cartGroup.POST("/clear", httpServ.ClearCart)
				cartGroup.POST("/checkout", httpServ.Checkout)
			}
			orderGroup := authorisedGroup.Group("/order")
			{
				orderGroup.GET("/", httpServ.GetOrders)
				orderGroup.GET("/:id", httpServ.GetOrder)
			}
			inventoryGroup := authorisedGroup.Group("/products")
			{
				inventoryGroup.GET("/", httpServ.GetAllProducts)
				inventoryGroup.GET("/:id", httpServ.GetProduct)
			}
		}
		adminGroup := superGroup.Group("/admin")
		adminGroup.Use(middlewares.AuthenticateAdmin())
		{
			productGroup := adminGroup.Group("/products")
			{
				productGroup.GET("/", httpServ.GetAllProducts)
				productGroup.POST("/create", httpServ.CreateProduct)
				productGroup.POST("/delete", httpServ.DeleteProduct)
			}
			inventoryGroup := adminGroup.Group("/inventory")
			{
				inventoryGroup.GET("/", httpServ.GetAllStocksFromInventory)
				inventoryGroup.GET("/:id", httpServ.GetStockFromInventory)
				inventoryGroup.POST("/add", httpServ.AddStockToInventory)
				inventoryGroup.POST("/remove", httpServ.RemoveStockFromInventory)
				inventoryGroup.POST("/modify", httpServ.UpdateStockInInventory)
			}
		}
	}

	router.Run(":8080")
}

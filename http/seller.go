package http

import (
	"strconv"

	"github.com/ProAltro/Amazon-Clone/entity"
	"github.com/gin-gonic/gin"
)

func (http HTTPService) SellerSignup(ctx *gin.Context) {
	var seller entity.Seller
	sellerService := http.SellerService
	err := ctx.BindJSON(&seller)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = sellerService.CreateSeller(&seller)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "seller created successfully",
	})
}

func (http HTTPService) SellerLogin(ctx *gin.Context) {
	var seller entity.Seller
	sellerService := http.SellerService
	err := ctx.BindJSON(&seller)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = sellerService.AuthenticateSeller(seller.Email, seller.Password)
	if err != nil {
		ctx.JSON(403, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "seller logged in successfully",
	})
}

func (http HTTPService) FetchSeller(ctx *gin.Context) {
	sellerService := http.SellerService
	id, err := strconv.Atoi(ctx.Query("id"))
	email := ctx.Query("email")
	if email == "" && err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	if email != "" {
		seller, err := sellerService.FindSellerByEmail(email)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{"seller": seller})
	} else {
		seller, err := sellerService.FindSellerByID(id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{"seller": seller})
	}

}

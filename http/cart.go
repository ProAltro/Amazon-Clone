package http

import (
	"github.com/ProAltro/Amazon-Clone/entity"
	"github.com/gin-gonic/gin"
)

func (http HTTPService) GetCart(ctx *gin.Context) {
	cartService := http.CartService
	cart, err := cartService.GetCart(ctx)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, cart)
}

func (http HTTPService) AddProductToCart(ctx *gin.Context) {
	type req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	var r req
	cartService := http.CartService
	err := ctx.BindJSON(&r)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = cartService.AddProductToCart(ctx, r.ProductID, r.Quantity)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "product added to cart",
	})
}

func (http HTTPService) RemoveProductFromCart(ctx *gin.Context) {
	type req struct {
		ProductID int `json:"product_id"`
	}
	var r req
	cartService := http.CartService
	err := ctx.BindJSON(&r)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = cartService.RemoveProductFromCart(ctx, r.ProductID)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "product removed from cart",
	})
}

func (http HTTPService) ModifyCart(ctx *gin.Context) {
	type req struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	var r req
	cartService := http.CartService
	err := ctx.BindJSON(&r)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	cart, err := cartService.ModifyCart(ctx, r.ProductID, r.Quantity)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, cart)
}

func (http HTTPService) ClearCart(ctx *gin.Context) {
	cartService := http.CartService
	err := cartService.ClearCart(ctx)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "cart cleared",
	})
}

func (http HTTPService) Checkout(ctx *gin.Context) {
	cartService := http.CartService
	err := cartService.Checkout(ctx)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "cart checked out",
	})
}

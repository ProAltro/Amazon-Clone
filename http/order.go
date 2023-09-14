package http

import (
	"github.com/ProAltro/Amazon-Clone/entity"
	"github.com/gin-gonic/gin"
)

func (http HTTPService) GetOrders(ctx *gin.Context) {
	orderService := http.OrderService
	orders, err := orderService.GetOrdersOfUser(ctx, ctx.Value("uid").(int))
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, orders)
}

func (http HTTPService) GetOrder(ctx *gin.Context) {
	orderService := http.OrderService
	order, err := orderService.GetOrder(ctx, ctx.GetInt("id"))
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, order)
}

package http

import (
	"strconv"

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
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "id must be an integer",
		})
		return
	}
	order, err := orderService.GetOrder(ctx, id)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, order)
}

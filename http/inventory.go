package http

import (
	"strconv"

	"github.com/ProAltro/Amazon-Clone/entity"
	"github.com/gin-gonic/gin"
)

func (http HTTPService) AddStockToInventory(ctx *gin.Context) {
	inventoryService := http.InventoryService
	type req struct {
		ID       int `json:"product_id" binding:"required"`
		Quantity int `json:"quantity" binding:"required"`
	}
	var resp req
	err := ctx.BindJSON(&resp)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = inventoryService.AddStockToInventory(ctx, resp.ID, resp.Quantity)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{"product with id " + strconv.Itoa(resp.ID): "added to inventory"})
}

func (http HTTPService) GetAllStocksFromInventory(ctx *gin.Context) {
	inventoryService := http.InventoryService
	stocks, err := inventoryService.GetAllStocksFromInventory(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"stocks": stocks,
	})
}

func (http HTTPService) GetStockFromInventory(ctx *gin.Context) {
	inventoryService := http.InventoryService
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "id must be an integer",
		})
		return
	}
	stock, err := inventoryService.GetStockFromInventory(ctx, id)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, stock)
}

func (http HTTPService) UpdateStockInInventory(ctx *gin.Context) {
	inventoryService := http.InventoryService
	type req struct {
		ID       int `json:"product_id" binding:"required"`
		Quantity int `json:"quantity" binding:"required"`
	}
	var resp req
	err := ctx.BindJSON(&resp)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = inventoryService.UpdateStockInInventory(ctx, resp.ID, resp.Quantity)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{"product with id " + strconv.Itoa(resp.ID): "updated in inventory"})
}

func (http HTTPService) RemoveStockFromInventory(ctx *gin.Context) {
	inventoryService := http.InventoryService
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "id must be an integer",
		})
		return
	}
	err = inventoryService.RemoveStockFromInventory(ctx, id)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{"product with id " + strconv.Itoa(id): "removed from inventory"})
}

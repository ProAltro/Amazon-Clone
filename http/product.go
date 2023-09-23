package http

import (
	"strconv"

	"github.com/ProAltro/Amazon-Clone/entity"
	"github.com/gin-gonic/gin"
)

func (http HTTPService) GetAllProducts(ctx *gin.Context) {
	productService := http.ProductService
	products, err := productService.GetAllProducts(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"products": products,
	})
}

func (http HTTPService) GetProduct(ctx *gin.Context) {
	productService := http.ProductService
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "id must be an integer",
		})
		return
	}
	product, err := productService.GetProduct(ctx, id)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, product)
}

func (http HTTPService) CreateProduct(ctx *gin.Context) {
	productService := http.ProductService
	var prod entity.Product
	err := ctx.BindJSON(&prod)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	product, err := productService.CreateProduct(ctx, prod.Name, prod.Description, prod.Price, prod.Seller)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "product created successfully",
		"product": product,
	})
}

func (http HTTPService) DeleteProduct(ctx *gin.Context) {
	productService := http.ProductService
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "id must be an integer",
		})
		return
	}
	err = productService.DeleteProduct(ctx, id)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "product deleted successfully",
	})
}

func (http HTTPService) GetProducts(ctx *gin.Context) {
	productService := http.ProductService
	ids := ctx.QueryArray("ids")
	intids := []int{}
	for _, id := range ids {
		intid, err := strconv.Atoi(id)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": "ids must be integers",
			})
			return
		}
		intids = append(intids, intid)
	}

	products, err := productService.GetProducts(ctx, intids)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, products)
}

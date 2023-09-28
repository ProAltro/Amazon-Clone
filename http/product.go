package http

import (
	"strconv"
	"strings"

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

	url_name := strings.ToLower(prod.Name) + strconv.Itoa(prod.ID)
	url_name = strings.ReplaceAll(url_name, " ", "_")
	url_name = strings.ToLower(url_name)
	image_urls, err := entity.CreateImages(prod.Images, url_name)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	product, err := productService.CreateProduct(ctx, prod.Name, prod.Description, prod.Price, prod.Seller, image_urls)
	if err != nil {
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}
	product, err = productService.GetProduct(ctx, product.ID)
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

func (http HTTPService) UpdateProduct(ctx *gin.Context) {
	productService := http.ProductService

	type req struct {
		ID          int      `json:"id" binding:"required"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Price       int      `json:"price" `
		Seller      string   `json:"seller"`
		Images      []string `json:"images"`
	}

	var prod req
	err := ctx.BindJSON(&prod)
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	if prod.Images != nil {
		url_name := strings.ToLower(prod.Name) + strconv.Itoa(prod.ID)
		url_name = strings.ReplaceAll(url_name, " ", "_")
		url_name = strings.ToLower(url_name)
		image_urls, err := entity.CreateImages(prod.Images, url_name)
		if err != nil {
			ctx.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}
		prod.Images = image_urls
	}

	product, err := productService.UpdateProduct(ctx, entity.Product{
		ID:          prod.ID,
		Name:        prod.Name,
		Description: prod.Description,
		Price:       prod.Price,
		Seller:      prod.Seller,
		Images:      prod.Images,
	})

	if err != nil {
		entity.DeleteImages(prod.Images)
		ctx.JSON(entity.GetStatusCode(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "product updated successfully",
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

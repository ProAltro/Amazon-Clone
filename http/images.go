package http

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func (http HTTPService) ServImage(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println(id)
	image, err := os.ReadFile("./images/" + strings.ToLower(id))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ctx.JSON(404, gin.H{
				"error": "image not found",
			})
			return
		}
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Data(200, "image/jpeg", image)
}

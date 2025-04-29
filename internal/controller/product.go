package controller

import (
	"net/http"

	"barcode-checker/internal/model"
	"barcode-checker/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	productService service.ProductService
}

func NewProductController(productService service.ProductService) *ProductController {
	return &ProductController{productService: productService}
}

func (c *ProductController) CheckProduct(ctx *gin.Context) {
	var req model.CheckProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid barcode format",
			"details": err.Error(),
		})
		return
	}

	userID, _ := ctx.Get("userID")
	var uid uint
	if userID != nil {
		uid = userID.(uint)
	}

	result, err := c.productService.CheckProduct(req.Barcode, uid, ctx.ClientIP())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Barcode check failed",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

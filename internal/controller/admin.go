package controller

import (
	"net/http"
	"strconv"

	"barcode-checker/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	authService service.AuthService
}

func NewAdminController(authService service.AuthService) *AdminController {
	return &AdminController{authService: authService}
}

func (c *AdminController) ListUsers(ctx *gin.Context) {
	users, err := c.authService.ListUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get users list",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"users": users,
		},
	})
}

func (c *AdminController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid user ID",
			"details": err.Error(),
		})
		return
	}

	if err := c.authService.DeleteUser(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

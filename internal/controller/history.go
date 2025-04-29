package controller

import (
	"net/http"
	"strconv"

	"barcode-checker/internal/model"
	"barcode-checker/internal/service"

	"github.com/gin-gonic/gin"
)

type HistoryController struct {
	historyService service.HistoryService
}

func NewHistoryController(historyService service.HistoryService) *HistoryController {
	return &HistoryController{historyService: historyService}
}

func (c *HistoryController) GetHistory(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	history, total, err := c.historyService.GetUserHistory(userID, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to get history",
			"details": err.Error(),
		})
		return
	}

	var response []model.ProductCheckResult
	for _, h := range history {
		response = append(response, model.ProductCheckResult{
			ID:             h.ID.Hex(),
			Barcode:        h.Barcode,
			IsOriginal:     h.IsOriginal,
			BarcodeCountry: h.BarcodeCountry,
			IPCountry:      h.IPCountry,
			CheckedAt:      h.CheckedAt,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"items": response,
			"meta": gin.H{
				"total":      total,
				"page":       page,
				"limit":      limit,
				"totalPages": (int(total) + limit - 1) / limit,
			},
		},
	})
}

func (c *HistoryController) DeleteHistoryItem(ctx *gin.Context) {
	userID := ctx.MustGet("userID").(uint)
	id := ctx.Param("id")

	if err := c.historyService.DeleteHistoryItem(userID, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete history item",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "History item deleted successfully",
	})
}

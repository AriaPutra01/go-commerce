package product

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service}
}

func (h *handler) UploadProductImage(c *gin.Context) {
	req := new(UploadProductImageRequest)
	if err := c.ShouldBind(req); err != nil {
		c.Error(err)
		return
	}

	res, err := h.service.UploadProductImage(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	uploadDir := "./uploads/public/products"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.Error(err)
		return
	}

	dst := filepath.Join(uploadDir, res.SavedFileName)
	if err := c.SaveUploadedFile(req.File, dst); err != nil {
		c.Error(err)
		return
	}

	seoPreviewUrl := fmt.Sprintf("/products/preview/%s", res.SavedFileName)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Upload product image Successfully",
		"data": gin.H{
			"id":       res.ID,
			"url_path": seoPreviewUrl,
		},
	})
}

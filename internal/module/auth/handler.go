package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service}
}

func (h *handler) Login(c *gin.Context) {
	req := new(LoginRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	atMaxAge := int((15 * time.Minute).Seconds())
	rtMaxAge := int((7 * 24 * time.Hour).Seconds())

	c.SetSameSite(http.SameSiteLaxMode)

	// ! Set secure to true in production
	c.SetCookie("at", result.AccessToken, atMaxAge, "/", "", false, true)
	c.SetCookie("rt", result.RefreshToken, rtMaxAge, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login Successfully",
	})
}

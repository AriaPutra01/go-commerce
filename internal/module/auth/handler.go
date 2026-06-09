package auth

import (
	"net/http"
	"time"

	"github.com/AriaPutra01/go-commerce/internal/constant"
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
		c.Error(err)
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
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

func (h *handler) Register(c *gin.Context) {
	req := new(RegisterRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		c.Error(err)
		return
	}

	err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Register user Successfully",
	})
}

func (h *handler) Refresh(c *gin.Context) {
	refreshCookie, err := c.Cookie("rt")
	if err != nil {
		c.Error(constant.ErrInvalidRefreshToken)
		return
	}

	result, err := h.service.Refresh(c.Request.Context(), refreshCookie)
	if err != nil {
		c.Error(err)
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

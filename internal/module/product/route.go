package product

import (
	"github.com/AriaPutra01/go-commerce/internal/constant"
	"github.com/AriaPutra01/go-commerce/internal/middleware"
	"github.com/gin-gonic/gin"
)

type route struct {
	app        *gin.Engine
	middleware *middleware.Middleware
	handler    *handler
}

func NewRoute(app *gin.Engine, middleware *middleware.Middleware, handler *handler) *route {
	return &route{
		app:        app,
		middleware: middleware,
		handler:    handler,
	}
}

func (r *route) RegisterRoute() {
	product := r.app.Group("/api/v1/product")

	productAuth := product.Use(r.middleware.WithAuth())
	productAuth.POST("/image-upload", r.middleware.WithPerm(constant.ProductCreate), r.handler.UploadProductImage)
}

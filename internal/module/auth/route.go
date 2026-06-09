package auth

import "github.com/gin-gonic/gin"

type route struct {
	app     *gin.Engine
	handler *handler
}

func NewRoute(app *gin.Engine, handler *handler) *route {
	return &route{
		app:     app,
		handler: handler,
	}
}

func (r *route) RegisterRoute() {
	auth := r.app.Group("/api/v1/auth")
	auth.POST("/login", r.handler.Login)
	auth.POST("/register", r.handler.Register)
	auth.POST("/refresh", r.handler.Refresh)
}

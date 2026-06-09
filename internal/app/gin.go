package app

import (
	"errors"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/AriaPutra01/go-commerce/internal/config"
	"github.com/AriaPutra01/go-commerce/internal/exception"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	ginSlog "github.com/gin-contrib/slog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewGin(cfg *config.Config, log *slog.Logger) *gin.Engine {
	app := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	app.Use(requestid.New())
	app.Use(ginSlog.SetLogger(
		ginSlog.WithWriter(os.Stdout),
		ginSlog.WithDefaultLevel(slog.Level(cfg.LogLevel)),
	))
	app.Use(MaxAllowed(20))
	app.Use(ErrorHandler(log))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Timezone"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return app
}

func MaxAllowed(n int) gin.HandlerFunc {
	sem := make(chan struct{}, n)
	return func(c *gin.Context) {
		select {
		case sem <- struct{}{}:
			defer func() { <-sem }()
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "server busy"})
		}
	}
}

func ErrorHandler(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var valErrors validator.ValidationErrors
			if errors.As(err, &valErrors) {
				errorMessages := make(map[string]string)
				for _, fieldErr := range valErrors {
					errorMessages[fieldErr.Field()] = fieldErr.Tag()
				}
				c.JSON(http.StatusBadRequest, gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "Not valid input data",
					"errors":  errorMessages,
				})
				return
			}

			var appErr *exception.AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.StatusCode, gin.H{
					"code":    appErr.ErrCode,
					"message": appErr.Message,
				})
				return
			}

			log.Error(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
				"success": false,
			})
		}
	}
}

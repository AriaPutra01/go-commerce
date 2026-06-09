package app

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/AriaPutra01/go-commerce/internal/exception"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewGin() *gin.Engine {
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

	app.Use(ErrorHandler())

	return app
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request first

		// Check if any errors were added to the context
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

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "An unexpected error occurred",
			})
		}
	}
}

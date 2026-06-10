package middleware

import (
	"slices"

	"github.com/AriaPutra01/go-commerce/internal/constant"
	"github.com/AriaPutra01/go-commerce/internal/token"
	"github.com/gin-gonic/gin"
)

type UserKey struct {
	ID    string
	Email string
	Role  string
}

type Middleware struct {
	jwt *token.JWTMaker
}

func NewMiddleware(jwt *token.JWTMaker) *Middleware {
	return &Middleware{jwt}
}

func (m *Middleware) WithAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessCookie, err := c.Cookie("at")
		if err != nil {
			c.AbortWithError(constant.ErrInvalidAccessToken.StatusCode, constant.ErrInvalidAccessToken)
			return
		}

		claims, err := m.jwt.VerifyToken(accessCookie)
		if err != nil {
			c.AbortWithError(constant.ErrInvalidAccessToken.StatusCode, constant.ErrInvalidAccessToken)
			return
		}

		c.Set(UserKey{}, &UserKey{
			ID:    claims.ID,
			Email: claims.Email,
			Role:  claims.Role,
		})

		c.Next()
	}
}

func (m *Middleware) WithPerm(perm constant.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, ok := c.Get(UserKey{})
		if !ok {
			c.AbortWithError(constant.ErrInvalidAccessToken.StatusCode, constant.ErrInvalidAccessToken)
			return
		}

		user, ok := userVal.(*UserKey)
		if !ok {
			c.AbortWithError(constant.ErrInvalidAccessToken.StatusCode, constant.ErrInvalidAccessToken)
			return
		}

		perms := constant.GetPermissionByRole(constant.Role(user.Role))
		if slices.Contains(perms, constant.Permission(perm)) {
			c.AbortWithError(constant.ErrAccessForbidden.StatusCode, constant.ErrAccessForbidden)
			return
		}

		c.Next()
	}
}

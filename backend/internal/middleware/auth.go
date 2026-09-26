package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/gbstudyapply/gbstudyapply/internal/config"
	"github.com/gbstudyapply/gbstudyapply/internal/constants"
	"github.com/gbstudyapply/gbstudyapply/internal/dto"
	"github.com/gbstudyapply/gbstudyapply/internal/util"
)

// UserKey is the gin context key for authenticated claims.
const UserKey = "user"

// AuthRequired validates JWT and injects claims.
func AuthRequired(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail(constants.CodeUnauthorized, constants.MsgUnauthorized))
			return
		}
		claims, err := util.ParseToken(strings.TrimPrefix(header, "Bearer "), cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Fail(constants.CodeUnauthorized, constants.MsgUnauthorized))
			return
		}
		c.Set(UserKey, claims)
		c.Next()
	}
}

// GetUserID extracts the authenticated user id.
func GetUserID(c *gin.Context) uint {
	v, _ := c.Get(UserKey)
	claims, ok := v.(*util.Claims)
	if !ok {
		return 0
	}
	return claims.UserID
}

// GetUserRole extracts the authenticated role.
func GetUserRole(c *gin.Context) string {
	v, _ := c.Get(UserKey)
	claims, ok := v.(*util.Claims)
	if !ok {
		return ""
	}
	return claims.Role
}

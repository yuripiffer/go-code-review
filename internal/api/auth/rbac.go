package auth

import (
	"net/http"

	"coupon_service/pkg"
	"github.com/gin-gonic/gin"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

func RequireRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get roles from context
		rolesInterface, exists := c.Get("roles")
		if !exists {
			c.AbortWithError(http.StatusForbidden, pkg.Errorf(pkg.EFORBIDDEN, "no roles found", nil))
			return
		}

		userRoles, ok := rolesInterface.([]string)
		if !ok {
			c.AbortWithError(http.StatusForbidden, pkg.Errorf(pkg.EFORBIDDEN, "invalid roles format", nil))
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range requiredRoles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
		}
		if !hasRole {
			c.AbortWithError(http.StatusForbidden, pkg.Errorf(pkg.EFORBIDDEN, "insufficient permissions", nil))
			return
		}
		c.Next()
	}
}

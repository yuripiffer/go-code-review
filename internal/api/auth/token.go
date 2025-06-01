package auth

import (
	"fmt"
	"net/http"
	"strings"

	"coupon_service/pkg"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TokenMiddleware(JWTSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithError(http.StatusUnauthorized, pkg.Errorf(pkg.EUNAUTHORIZED, "missing authorization header", nil))
			return
		}

		bearerToken := strings.Split(authHeader, "Bearer ")
		if len(bearerToken) != 2 {
			c.AbortWithError(http.StatusUnauthorized, pkg.Errorf(pkg.EUNAUTHORIZED, "invalid token format", nil))
			return
		}

		token, err := jwt.ParseWithClaims(bearerToken[1], &jwt.MapClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(JWTSecret), nil
		})
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, pkg.Errorf(pkg.EUNAUTHORIZED, "invalid token", err))
			return
		}

		if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
			r := (*claims)["roles"]
			roles := make([]string, len(r.([]interface{})))
			for i, v := range r.([]interface{}) {
				roles[i] = fmt.Sprint(v)
			}
			c.Set("roles", roles)

			c.Set("user_id", (*claims)["sub"])

			c.Next()
		} else {
			c.AbortWithError(http.StatusUnauthorized, pkg.Errorf(pkg.EUNAUTHORIZED, "invalid token claims", nil))
		}
	}
}

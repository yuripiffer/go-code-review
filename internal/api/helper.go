package api

import (
	"coupon_service/pkg"
	"github.com/gin-gonic/gin"
)

// WebErr should be used specially when the error status was set before the web package
func WebErr(c *gin.Context, err error) {
	c.JSON(pkg.ErrorStatusCode(pkg.ErrorCode(err)), err)
}

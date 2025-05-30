package api

import (
	"net/http"

	"coupon_service/pkg"

	"github.com/gin-gonic/gin"
)

type ApplyCouponRequest struct {
	Code  string `json:"code"`
	Value int    `json:"value"`
}

type ApplyCouponResponse struct {
	Value                 int  `json:"value"`
	AppliedDiscount       int  `json:"applied_discount"`
	ApplicationSuccessful bool `json:"application_successful"`
}

// ApplyCoupon godoc
// @Summary      Apply a coupon to a basket
// @Description  Applies a coupon code to a given basket value and returns the result or error
// @Tags         coupons
// @Accept       json
// @Produce      json
// @Param        request body ApplyCouponRequest true "Coupon code and basket value"
// @Success      200 {object} ApplyCouponResponse
// @Failure      400 {object} pkg.Error
// @Failure      404 {object} pkg.Error
// @Router       /apply [post]
func (a *API) ApplyCoupon(c *gin.Context) {
	input := ApplyCouponRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "invalid request body", err))
		return
	}

	if input.Code == "" {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "code cannot be empty", nil))
		return
	}

	if input.Value <= 0 {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "value should be positive", nil))
		return
	}

	basket, err := a.svc.ApplyCoupon(input.Code, input.Value)
	if err != nil {
		WebErr(c, err)
		return
	}

	c.JSON(http.StatusOK, ApplyCouponResponse{
		Value:                 basket.Value,
		AppliedDiscount:       basket.AppliedDiscount,
		ApplicationSuccessful: basket.ApplicationSuccessful})
}

type CreateCouponRequest struct {
	Code           string `json:"code"`
	Discount       int    `json:"discount"`
	MinBasketValue int    `json:"minimum_basket_value"`
}

// CreateCoupon godoc
// @Summary      Create a new coupon
// @Description  Creates a new coupon with the specified code, discount, and minimum basket value
// @Tags         coupons
// @Accept       json
// @Param        request body CreateCouponRequest true "Coupon details"
// @Success      201
// @Failure      400 {object} pkg.Error
// @Failure      409 {object} pkg.Error
// @Router       /create [post]
func (a *API) CreateCoupon(c *gin.Context) {
	input := CreateCouponRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "invalid request body", err))
		return
	}

	if input.Code == "" {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "code cannot be empty", nil))
		return
	}

	if input.Discount <= 0 {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "discount should be positive", nil))
		return
	}

	if input.MinBasketValue <= 0 {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "minimum_basket_value should be positive", nil))
		return
	}

	err := a.svc.CreateCoupon(input.Code, input.Discount, input.MinBasketValue)
	if err != nil {
		WebErr(c, err)
		return
	}

	c.Status(http.StatusCreated)
}

type GetCouponsRequest struct {
	Codes []string `json:"codes"`
}

type GetCouponsResponse struct {
	Coupons []CouponResponse `json:"coupons"`
}

type CouponResponse struct {
	ID             string `json:"id"`
	Code           string `json:"code"`
	Discount       int    `json:"discount"`
	MinBasketValue int    `json:"minimum_basket_value"`
}

// GetCoupons godoc
// @Summary      Get coupons by codes
// @Description  Retrieves coupon details for the provided list of coupons if they are all existent
// @Tags         coupons
// @Accept       json
// @Produce      json
// @Param        request body GetCouponsRequest true "List of coupon codes"
// @Success      200 {array} CouponResponse
// @Failure      400 {object} pkg.Error
// @Failure      404 {object} pkg.Error
// @Router       /coupons [get]
func (a *API) GetCoupons(c *gin.Context) {
	input := GetCouponsRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "invalid request body", err))
		return
	}

	if len(input.Codes) == 0 {
		WebErr(c, pkg.Errorf(pkg.EINVALID, "minimum of one coupon code required", nil))
		return
	}

	coupons, err := a.svc.GetCoupons(input.Codes)
	if err != nil {
		WebErr(c, err)
		return
	}

	couponsResponse := make([]CouponResponse, len(coupons))
	for i, c := range coupons {
		couponsResponse[i] = CouponResponse{
			ID:             c.ID,
			Code:           c.Code,
			Discount:       c.Discount,
			MinBasketValue: c.MinBasketValue,
		}
	}

	c.JSON(http.StatusOK, couponsResponse)
}

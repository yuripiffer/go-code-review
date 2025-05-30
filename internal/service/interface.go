package service

import (
	"coupon_service/internal/entity"
)

//go:generate go run github.com/matryer/moq -out service_mock.go -stub . CouponService
type CouponService interface {
	ApplyCoupon(code string, value int) (entity.Basket, error)
	CreateCoupon(code string, discount, minBasketValue int) error
	GetCoupons([]string) ([]entity.Coupon, error)
}

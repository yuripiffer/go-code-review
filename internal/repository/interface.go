package repository

import (
	"coupon_service/internal/entity"
)

//go:generate go run github.com/matryer/moq -out repository_mock.go -stub . CouponRepository
type CouponRepository interface {
	FindByCode(string) (entity.Coupon, error)
	Save(entity.Coupon) error
}

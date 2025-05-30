package service

import (
	"strings"

	"coupon_service/internal/entity"
	"coupon_service/pkg"

	"github.com/google/uuid"
)

func (s Service) ApplyCoupon(code string, value int) (entity.Basket, error) {
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return entity.Basket{}, err
	}

	if value < coupon.MinBasketValue {
		return entity.Basket{}, pkg.Errorf(pkg.EINVALID, "basket value below minimum required for coupon", nil)
	}

	return entity.Basket{
		Value:                 value,
		AppliedDiscount:       coupon.Discount,
		ApplicationSuccessful: true,
	}, nil
}

func (s Service) CreateCoupon(code string, discount, minBasketValue int) error {
	if len(code) < 6 {
		return pkg.Errorf(pkg.EINVALID, "minimum length of code is 6 characters", nil)
	}

	if !pkg.IsAlphaNumericOnly(code) {
		return pkg.Errorf(pkg.EINVALID, "code must contain only number and letters", nil)
	}

	if discount > minBasketValue {
		return pkg.Errorf(pkg.EINVALID, "discount bigger than minimum basket value", nil)
	}

	_, err := s.repo.FindByCode(code)
	if err == nil {
		return pkg.Errorf(pkg.ECONFLICT, "coupon already exists", nil)
	}

	coupon := entity.Coupon{
		ID:             uuid.New().String(),
		Code:           strings.ToUpper(code),
		Discount:       discount,
		MinBasketValue: minBasketValue,
	}
	if err := s.repo.Save(coupon); err != nil {
		return err
	}
	return nil
}

func (s Service) GetCoupons(codes []string) ([]entity.Coupon, error) {
	coupons := make([]entity.Coupon, len(codes))

	for i, code := range codes {
		coupon, err := s.repo.FindByCode(code)
		if err != nil {
			return nil, err
		}
		coupons[i] = coupon
	}
	return coupons, nil
}

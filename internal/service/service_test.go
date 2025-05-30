package service

import (
	"errors"
	"testing"

	"coupon_service/internal/entity"
	"coupon_service/internal/repository"
	"coupon_service/pkg"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_ApplyCoupon(t *testing.T) {
	tests := []struct {
		name         string
		code         string
		value        int
		findCoupon   entity.Coupon
		findErr      error
		expectedErr  error
		expectedResp entity.Basket
	}{
		{
			name:  "success",
			code:  "ABC123",
			value: 200,
			findCoupon: entity.Coupon{
				Code:           "ABC123",
				Discount:       20,
				MinBasketValue: 100,
			},
			expectedResp: entity.Basket{
				Value:                 200,
				AppliedDiscount:       20,
				ApplicationSuccessful: true,
			},
		},
		{
			name:        "coupon not found",
			code:        "ABC123",
			value:       100,
			findErr:     pkg.Errorf(pkg.ENOTFOUND, "coupon not found", nil),
			expectedErr: pkg.Errorf(pkg.ENOTFOUND, "coupon not found", nil),
		},
		{
			name:  "basket value below coupon minimum value",
			code:  "ABC123",
			value: 50,
			findCoupon: entity.Coupon{
				Code:           "ABC123",
				Discount:       10,
				MinBasketValue: 100,
			},
			expectedErr: pkg.Errorf(pkg.EINVALID, "basket value below minimum required for coupon", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := &repository.CouponRepositoryMock{
				FindByCodeFunc: func(code string) (entity.Coupon, error) {
					assert.Equal(t, tt.code, code)
					return tt.findCoupon, tt.findErr
				},
			}
			svc := New(repoMock)
			basket, err := svc.ApplyCoupon(tt.code, tt.value)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResp, basket)
			}
		})
	}
}

func TestService_CreateCoupon(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		discount       int
		minBasketValue int
		findErr        error
		saveErr        error
		expectedErr    error
	}{
		{
			name:           "success",
			code:           "ABC123",
			discount:       10,
			minBasketValue: 100,
			findErr:        pkg.Errorf(pkg.ECONFLICT, "coupon already exists", nil),
		},
		{
			name:           "code too short",
			code:           "ABC",
			discount:       10,
			minBasketValue: 100,
			expectedErr:    pkg.Errorf(pkg.EINVALID, "minimum length of code is 6 characters", nil),
		},
		{
			name:           "code not alphanumeric",
			code:           "ABCD$12",
			discount:       10,
			minBasketValue: 100,
			expectedErr:    pkg.Errorf(pkg.EINVALID, "code must contain only number and letters", nil),
		},
		{
			name:           "code with space",
			code:           "ABCD 12",
			discount:       10,
			minBasketValue: 100,
			expectedErr:    pkg.Errorf(pkg.EINVALID, "code must contain only number and letters", nil),
		},
		{
			name:           "discount greater than minBasketValue",
			code:           "ABC123",
			discount:       200,
			minBasketValue: 100,
			expectedErr:    pkg.Errorf(pkg.EINVALID, "discount bigger than minimum basket value", nil),
		},
		{
			name:           "coupon already exists",
			code:           "ABC123",
			discount:       10,
			minBasketValue: 100,
			findErr:        nil,
			expectedErr:    pkg.Errorf(pkg.ECONFLICT, "coupon already exists", nil),
		},
		{
			name:           "save error",
			code:           "ABC123",
			discount:       10,
			minBasketValue: 100,
			findErr:        pkg.Errorf(pkg.ECONFLICT, "coupon already exists", nil),
			saveErr:        pkg.Errorf(pkg.EINTERNAL, "save coupon error", errors.New("db error")),
			expectedErr:    pkg.Errorf(pkg.EINTERNAL, "save coupon error", errors.New("db error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := &repository.CouponRepositoryMock{
				FindByCodeFunc: func(code string) (entity.Coupon, error) {
					assert.Equal(t, tt.code, code)
					return entity.Coupon{}, tt.findErr
				},
				SaveFunc: func(coupon entity.Coupon) error {
					_, err := uuid.Parse(coupon.ID)
					assert.NoError(t, err)
					assert.Equal(t, tt.code, coupon.Code)
					assert.Equal(t, tt.discount, coupon.Discount)
					assert.Equal(t, tt.minBasketValue, coupon.MinBasketValue)
					return tt.saveErr
				},
			}
			svc := New(repoMock)
			err := svc.CreateCoupon(tt.code, tt.discount, tt.minBasketValue)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetCoupons(t *testing.T) {
	tests := []struct {
		name        string
		codes       []string
		findCoupons []entity.Coupon
		findErrs    []error
		expectedErr error
	}{
		{
			name:  "all found",
			codes: []string{"ABC123", "DEF456"},
			findCoupons: []entity.Coupon{
				{
					ID:             "uuid-123",
					Code:           "ABC123",
					Discount:       10,
					MinBasketValue: 100,
				},
				{ID: "uuid-456",
					Code:           "DEF456",
					Discount:       10,
					MinBasketValue: 100,
				},
			},
			findErrs: []error{nil, nil},
		},
		{
			name:  "error in the second coupon",
			codes: []string{"ABC123", ""},
			findCoupons: []entity.Coupon{
				{
					ID:             "uuid-123",
					Code:           "ABC123",
					Discount:       10,
					MinBasketValue: 100,
				},
				{},
			},
			findErrs: []error{
				nil,
				pkg.Errorf(pkg.ENOTFOUND, "coupon not found", nil),
			},
			expectedErr: pkg.Errorf(pkg.ENOTFOUND, "coupon not found", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call := 0
			repoMock := &repository.CouponRepositoryMock{
				FindByCodeFunc: func(code string) (entity.Coupon, error) {
					defer func() { call++ }()
					assert.Equal(t, tt.codes[call], code)
					return tt.findCoupons[call], tt.findErrs[call]
				},
			}
			svc := New(repoMock)
			coupons, err := svc.GetCoupons(tt.codes)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.findCoupons, coupons)
			}
		})
	}
}

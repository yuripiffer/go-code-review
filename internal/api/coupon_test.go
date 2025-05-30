package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"coupon_service/internal/entity"
	"coupon_service/internal/service"
	"coupon_service/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(api *API) *gin.Engine {
	r := gin.Default()
	r.POST("/coupon/apply", api.ApplyCoupon)
	r.POST("/coupon/create", api.CreateCoupon)
	r.POST("/coupon/get", api.GetCoupons)
	return r
}

func TestAPI_ApplyCoupon(t *testing.T) {
	tests := []struct {
		name         string
		input        ApplyCouponRequest
		mockBasket   entity.Basket
		mockSvcError error
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			input: ApplyCouponRequest{
				Code:  "ABCDEF123",
				Value: 100,
			},
			mockBasket: entity.Basket{
				Value:                 100,
				AppliedDiscount:       10,
				ApplicationSuccessful: true,
			},
			mockSvcError: nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid: code is empty",
			input: ApplyCouponRequest{
				Code:  "",
				Value: 100,
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "code cannot be empty",
		},
		{
			name: "invalid: value is 0",
			input: ApplyCouponRequest{
				Code:  "ABCDEF123",
				Value: 0,
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "value should be positive",
		},
		{
			name: "invalid: value is negative",
			input: ApplyCouponRequest{
				Code:  "ABCDEF123",
				Value: -20,
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "value should be positive",
		},
		{
			name: "service error",
			input: ApplyCouponRequest{
				Code:  "ABCDEF",
				Value: 100,
			},
			mockSvcError: pkg.Errorf(pkg.EINVALID, "abc test error", nil),
			expectedCode: http.StatusBadRequest,
			expectedBody: "abc test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svcMock := &service.CouponServiceMock{
				ApplyCouponFunc: func(code string, value int) (entity.Basket, error) {
					assert.Equal(t, tt.input.Code, code)
					assert.Equal(t, tt.input.Value, value)
					return tt.mockBasket, tt.mockSvcError
				},
			}
			api := &API{svc: svcMock}
			router := setupRouter(api)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/coupon/apply", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestAPI_CreateCoupon(t *testing.T) {
	tests := []struct {
		name         string
		input        CreateCouponRequest
		mockSvcError error
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			input: CreateCouponRequest{
				Code:           "ABCDEF123",
				Discount:       10,
				MinBasketValue: 20,
			},
			mockSvcError: nil,
			expectedCode: http.StatusCreated,
		},
		{
			name: "empty code",
			input: CreateCouponRequest{
				Code:           "",
				Discount:       10,
				MinBasketValue: 20,
			},
			mockSvcError: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "code cannot be empty",
		},
		{
			name: "invalid, discount is 0",
			input: CreateCouponRequest{
				Code:           "ABCDEF123",
				Discount:       0,
				MinBasketValue: 20,
			},
			mockSvcError: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "discount should be positive",
		},
		{
			name: "invalid, discount is negative",
			input: CreateCouponRequest{
				Code:           "ABCDEF123",
				Discount:       -20,
				MinBasketValue: 20,
			},
			mockSvcError: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "discount should be positive",
		},
		{
			name: "invalid, minimum_basket_value is 0",
			input: CreateCouponRequest{
				Code:           "ABCDEF123",
				Discount:       10,
				MinBasketValue: 0,
			},
			mockSvcError: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "minimum_basket_value should be positive",
		},
		{
			name: "invalid, minimum_basket_value is negative",
			input: CreateCouponRequest{
				Code:           "ABCDEF123",
				Discount:       10,
				MinBasketValue: -20,
			},
			mockSvcError: nil,
			expectedCode: http.StatusBadRequest,
			expectedBody: "minimum_basket_value should be positive",
		},
		{
			name: "service error",
			input: CreateCouponRequest{
				Code:           "ABCDEF",
				Discount:       10,
				MinBasketValue: 20,
			},
			mockSvcError: pkg.Errorf(pkg.ECONFLICT, "coupon already exists", nil),
			expectedCode: http.StatusConflict,
			expectedBody: "coupon already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svcMock := &service.CouponServiceMock{
				CreateCouponFunc: func(code string, discount, minBasketValue int) error {
					assert.Equal(t, tt.input.Code, code)
					assert.Equal(t, tt.input.Discount, discount)
					assert.Equal(t, tt.input.MinBasketValue, minBasketValue)
					return tt.mockSvcError
				},
			}
			api := &API{svc: svcMock}
			router := setupRouter(api)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/coupon/create", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestAPI_GetCoupons(t *testing.T) {
	tests := []struct {
		name         string
		input        GetCouponsRequest
		mockCoupons  []entity.Coupon
		mockSvcError error
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			input: GetCouponsRequest{
				Codes: []string{"ABCDEF123", "XYZ987COUPON"},
			},
			mockCoupons: []entity.Coupon{
				{
					ID:             "1",
					Code:           "ABCDEF123",
					Discount:       10,
					MinBasketValue: 20,
				},
				{
					ID:             "2",
					Code:           "XYZ987COUPON",
					Discount:       10,
					MinBasketValue: 20,
				},
			},
			mockSvcError: nil,
			expectedCode: http.StatusOK,
		},
		{
			name: "invalid: no coupon is provided",
			input: GetCouponsRequest{
				Codes: []string{},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "minimum of one coupon code required",
		},
		{
			name: "service error",
			input: GetCouponsRequest{
				Codes: []string{"ABCDEF"},
			},
			mockSvcError: pkg.Errorf(pkg.ECONFLICT, "service error", nil),
			expectedCode: http.StatusConflict,
			expectedBody: "service error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svcMock := &service.CouponServiceMock{
				GetCouponsFunc: func(codes []string) ([]entity.Coupon, error) {
					assert.Equal(t, tt.input.Codes, codes)
					return tt.mockCoupons, tt.mockSvcError
				},
			}
			api := &API{svc: svcMock}
			router := setupRouter(api)

			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/coupon/get", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

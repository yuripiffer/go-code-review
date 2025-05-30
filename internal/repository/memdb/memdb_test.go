package memdb

import (
	"testing"

	"coupon_service/internal/entity"
	"coupon_service/pkg"
	"github.com/stretchr/testify/assert"
)

func TestRepository_FindByCode(t *testing.T) {
	tests := []struct {
		name    string
		entries map[string]entity.Coupon
		code    string
		want    entity.Coupon
		wantErr error
	}{
		{
			name: "coupon found",
			entries: map[string]entity.Coupon{
				"ABC123": {Code: "ABC123", Discount: 10, MinBasketValue: 100},
			},
			code:    "ABC123",
			want:    entity.Coupon{Code: "ABC123", Discount: 10, MinBasketValue: 100},
			wantErr: nil,
		},
		{
			name:    "coupon not found",
			entries: map[string]entity.Coupon{},
			code:    "NOTFOUND",
			want:    entity.Coupon{},
			wantErr: pkg.Errorf(pkg.ENOTFOUND, "coupon not found", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{entries: tt.entries}
			got, err := r.FindByCode(tt.code)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestRepository_Save(t *testing.T) {
	tests := []struct {
		name    string
		initial map[string]entity.Coupon
		coupon  entity.Coupon
		want    map[string]entity.Coupon
	}{
		{
			name: "save new coupon",
			initial: map[string]entity.Coupon{
				"OLDCOUPON10": {Code: "OLDCOUPON10", Discount: 10, MinBasketValue: 100},
			},
			coupon: entity.Coupon{Code: "NEW123", Discount: 15, MinBasketValue: 200},
			want: map[string]entity.Coupon{
				"OLDCOUPON10": {Code: "OLDCOUPON10", Discount: 10, MinBasketValue: 100},
				"NEW123":      {Code: "NEW123", Discount: 15, MinBasketValue: 200},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{entries: tt.initial}
			err := r.Save(tt.coupon)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, r.entries)
		})
	}
}

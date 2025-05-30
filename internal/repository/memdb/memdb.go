package memdb

import (
	"coupon_service/internal/entity"
	"coupon_service/pkg"
)

func (r *Repository) FindByCode(code string) (entity.Coupon, error) {
	coupon, ok := r.entries[code]
	if !ok {
		return entity.Coupon{}, pkg.Errorf(pkg.ENOTFOUND, "coupon not found", nil)
	}
	return coupon, nil
}

func (r *Repository) Save(coupon entity.Coupon) error {
	r.entries[coupon.Code] = coupon
	return nil
}

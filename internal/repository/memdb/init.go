package memdb

import (
	"coupon_service/internal/entity"
)

type Repository struct {
	entries map[string]entity.Coupon
}

func NewRepository() *Repository {
	return &Repository{
		entries: make(map[string]entity.Coupon),
	}
}

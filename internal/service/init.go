package service

import "coupon_service/internal/repository"

type Service struct {
	repo repository.CouponRepository
}

func New(repo repository.CouponRepository) Service {
	return Service{
		repo: repo,
	}
}

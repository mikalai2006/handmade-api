package service

import (
	"github.com/mikalai2006/handmade/internal/domain"
	"github.com/mikalai2006/handmade/internal/repository"
)

type ShopService struct {
	repo repository.Shop
}

func NewShopService(repo repository.Shop) *ShopService  {
	return &ShopService{repo: repo}
}


func (s *ShopService) GetAllShops() (domain.Response, error) {
	return s.repo.GetAllShops()
}

func (s *ShopService) CreateShop(userId string, shop domain.Shop) (*domain.Shop, error)  {
	return s.repo.CreateShop(userId, shop)
}

package service

import "github.com/s21platform/advert-service/internal/model"

type DBRepo interface {
	GetAdvert(UUID string) (*model.AdvertInfoList, error)
}

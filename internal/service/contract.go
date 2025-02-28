//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	advert "github.com/s21platform/advert-proto/advert-proto"
	"github.com/s21platform/advert-service/internal/model"
)

type DBRepo interface {
	CreateAdvert(ctx context.Context, UUID string, in *advert.CreateAdvertIn) error
	GetAdverts(UUID string) (*model.AdvertInfoList, error)
	CancelAdvert(ctx context.Context, in *advert.CancelAdvertIn) error
}

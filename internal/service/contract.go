//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"
	"time"

	"github.com/s21platform/advert-service/internal/model"
	"github.com/s21platform/advert-service/pkg/advert"
)

type DBRepo interface {
	CreateAdvert(ctx context.Context, UUID string, in *advert.CreateAdvertIn) error
	GetAdverts(UUID string) (*model.AdvertInfoList, error)
	CancelAdvert(ctx context.Context, in *advert.CancelAdvertIn) error
	GetAdvertCancelExpiry(ctx context.Context, ID int64) (*model.AdvertCancelExpiry, error)
	RestoreAdvert(ctx context.Context, ID int64, newExpiredAt time.Time) error
}

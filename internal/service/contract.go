//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import "github.com/s21platform/advert-service/internal/model"

type DBRepo interface {
	GetAdverts(UUID string) (*model.AdvertInfoList, error)
}

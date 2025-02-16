//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import advert "github.com/s21platform/advert-proto/advert-proto"

type DBRepo interface {
	CreateAdvert(UUID string, in *advert.CreateAdvertIn) error
}

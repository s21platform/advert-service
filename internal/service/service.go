package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	advert "github.com/s21platform/advert-proto/advert-proto"
	"github.com/s21platform/advert-service/internal/config"
)

type Service struct {
	advert.UnimplementedAdvertServiceServer
	dbR DBRepo
}

func New(dbR DBRepo) *Service {
	return &Service{dbR: dbR}
}

func (s *Service) CreateAdvert(ctx context.Context, in *advert.CreateAdvertIn) (*advert.AdvertEmpty, error) {
	ownerUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to retrieve uuid")
	}

	err := s.dbR.CreateAdvert(ownerUUID, in)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to create advert: %v", err)
	}

	return &advert.AdvertEmpty{}, nil
}

package service

import (
	"context"

	advert "github.com/s21platform/advert-proto/advert-proto"
	"github.com/s21platform/advert-service/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	advert.UnimplementedAdvertServiceServer
	dbR DBRepo
}

func New(dbR DBRepo) *Service {
	return &Service{dbR: dbR}
}

func (s *Service) GetAdvert(ctx context.Context, in *advert.AdvertEmpty) (*advert.GetAdvertOut, error) {
	ownerUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	adverts, err := s.dbR.GetAdvert(ownerUUID)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to find adverts: %v", err)
	}

	return &advert.GetAdvertOut{
		Adverts: adverts.FromDTO(),
	}, nil
}

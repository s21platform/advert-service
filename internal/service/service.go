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

func (s *Service) GetAdverts(ctx context.Context, _ *advert.AdvertEmpty) (*advert.GetAdvertsOut, error) {
	ownerUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	adverts, err := s.dbR.GetAdverts(ownerUUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find adverts: %v", err)
	}

	return &advert.GetAdvertsOut{
		Adverts: adverts.FromDTO(),
	}, nil
}

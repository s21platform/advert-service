package service

import (
	"context"
	"time"

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

	err := s.dbR.CreateAdvert(ctx, ownerUUID, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create advert: %v", err)
	}

	return &advert.AdvertEmpty{}, nil
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

func (s *Service) CancelAdvert(ctx context.Context, in *advert.CancelAdvertIn) (*advert.AdvertEmpty, error) {
	err := s.dbR.CancelAdvert(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to cancel advert: %v", err)
	}

	return &advert.AdvertEmpty{}, nil
}

func (s *Service) RestoreAdvert(ctx context.Context, in *advert.RestoreAdvertIn) (*advert.AdvertEmpty, error) {
	cancelExpiry, err := s.dbR.GetAdvertCancelExpiry(ctx, in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get advert cancel info: %v", err)
	}

	if cancelExpiry.CanceledAt == nil {
		return nil, status.Errorf(codes.Internal, "failed to advert wasn't canceled")
	}

	timeDiff := time.Since(*cancelExpiry.CanceledAt)
	newExpiredAt := cancelExpiry.ExpiredAt.Add(timeDiff)

	err = s.dbR.RestoreAdvert(ctx, in.Id, newExpiredAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to restore advert: %v", err)
	}

	return &advert.AdvertEmpty{}, nil
}

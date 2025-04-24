package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/advert-service/internal/config"
	"github.com/s21platform/advert-service/internal/model"
	"github.com/s21platform/advert-service/pkg/advert"
)

type Service struct {
	advert.UnimplementedAdvertServiceServer
	dbR DBRepo
}

func New(dbR DBRepo) *Service {
	return &Service{dbR: dbR}
}

func (s *Service) CreateAdvert(ctx context.Context, in *advert.CreateAdvertIn) (*advert.AdvertEmpty, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("CreateAdvert")

	ownerUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to retrieve uuid")
	}

	err := s.dbR.CreateAdvert(ctx, ownerUUID, in)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create advert: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to create advert: %v", err)
	}

	return &advert.AdvertEmpty{}, nil
}

func (s *Service) GetAdverts(ctx context.Context, _ *advert.AdvertEmpty) (*advert.GetAdvertsOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetAdverts")

	ownerUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	adverts, err := s.dbR.GetAdverts(ownerUUID)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to find adverts: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to find adverts: %v", err)
	}

	return &advert.GetAdvertsOut{
		Adverts: adverts.FromDTO(),
	}, nil
}

func (s *Service) CancelAdvert(ctx context.Context, in *advert.CancelAdvertIn) (*advert.AdvertEmpty, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("CancelAdvert")

	err := s.dbR.CancelAdvert(ctx, in)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to cancel advert: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to cancel advert: %v", err)
	}

	return &advert.AdvertEmpty{}, nil
}

func (s *Service) RestoreAdvert(ctx context.Context, in *advert.RestoreAdvertIn) (*advert.AdvertEmpty, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("RestoreAdvert")

	cancelExpiry, err := s.dbR.GetAdvertCancelExpiry(ctx, in.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get advert cancel info: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get advert cancel info: %v", err)
	}

	if !cancelExpiry.IsCanceled {
		logger.Error("failed to restore the advert due to a missing cancellation record")
		return nil, status.Errorf(codes.Internal, "failed to restore the advert due to a missing cancellation record")
	}

	timeDiff := time.Since(*cancelExpiry.CanceledAt)
	newExpiredAt := cancelExpiry.ExpiredAt.Add(timeDiff)

	err = s.dbR.RestoreAdvert(ctx, in.Id, newExpiredAt)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to restore advert: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to restore advert: %v", err)
	}

	return &advert.AdvertEmpty{}, nil
}

func (s *Service) EditAdvert(ctx context.Context, in *advert.EditAdvertIn) (*advert.AdvertEmpty, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("EditAdvert")

	isActive, err := s.dbR.IsAdvertActive(ctx, int(in.Id))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to advert is not active: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to advert is not active: %v", err)
	}

	if !isActive {
		logger.Error("failed to advert is not active")
		return nil, status.Errorf(codes.Unavailable, "failed to advert is not active")
	}

	uuid, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		logger.Error("failed to find uuid")
		return nil, status.Errorf(codes.Unauthenticated, "failed to find uuid")
	}

	ownerUUID, err := s.dbR.GetOwnerUUID(ctx, int(in.Id))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to get owner uuid: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get owner uuid: %v", err)
	}

	if ownerUUID != uuid {
		logger.Error("failed to user is not advert owner")
		return nil, status.Errorf(codes.PermissionDenied, "failed to user is not advert owner")
	}

	newAdvertData := &model.EditAdvert{}
	newAdvertData.ToDTO(in)
	err = s.dbR.EditAdvert(ctx, newAdvertData)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to edit advert: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to edit advert: %v", err)
	}

	return &advert.AdvertEmpty{}, nil
}

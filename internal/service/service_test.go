package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	advertproto "github.com/s21platform/advert-proto/advert-proto"

	"github.com/s21platform/advert-service/internal/config"
	"github.com/s21platform/advert-service/internal/model"
)

func TestServer_GetAdverts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)

	t.Run("get_ok", func(t *testing.T) {
		expectedAdverts := &model.AdvertInfoList{
			{
				Content:   "деревянные изделия",
				ExpiredAt: time.Now(),
			},
			{
				Content:   "деревянные изделия ручной работы",
				ExpiredAt: time.Now(),
			},
		}

		mockRepo.EXPECT().GetAdverts(uuid).Return(expectedAdverts, nil)

		s := New(mockRepo)
		adverts, err := s.GetAdverts(ctx, &advertproto.AdvertEmpty{})
		assert.NoError(t, err)
		assert.Equal(t, adverts, &advertproto.GetAdvertsOut{Adverts: adverts.Adverts})
	})

	t.Run("get_no_uuid", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockDBRepo(ctrl)

		s := New(mockRepo)
		_, err := s.GetAdverts(ctx, &advertproto.AdvertEmpty{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("get_err", func(t *testing.T) {
		expectedAdverts := &model.AdvertInfoList{}
		expectedErr := errors.New("get err")

		mockRepo.EXPECT().GetAdverts(uuid).Return(expectedAdverts, expectedErr)

		s := New(mockRepo)
		_, err := s.GetAdverts(ctx, &advertproto.AdvertEmpty{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to find adverts: get err")
	})
}

func TestServer_CreateAdverts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)

	t.Run("create_ok", func(t *testing.T) {
		mockRepo.EXPECT().CreateAdvert(ctx, uuid, gomock.Any()).Return(nil)

		s := New(mockRepo)
		_, err := s.CreateAdvert(ctx, &advertproto.CreateAdvertIn{})
		assert.NoError(t, err)
	})

	t.Run("create_no_uuid", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockDBRepo(ctrl)

		s := New(mockRepo)
		_, err := s.CreateAdvert(ctx, &advertproto.CreateAdvertIn{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to retrieve uuid")
	})

	t.Run("create_err", func(t *testing.T) {
		expectedErr := errors.New("get err")

		mockRepo.EXPECT().CreateAdvert(ctx, uuid, &advertproto.CreateAdvertIn{}).Return(expectedErr)

		s := New(mockRepo)
		_, err := s.CreateAdvert(ctx, &advertproto.CreateAdvertIn{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to create advert: get err")
	})
}

func TestServer_RestoreAdvert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ID := int64(123)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)

	t.Run("should_return_ok", func(t *testing.T) {
		baseTime := time.Date(2025, 3, 4, 21, 0, 0, 0, time.UTC)
		canceledAt := baseTime
		expiredAt := baseTime.Add(1 * time.Hour)

		expectedCancelExpiry := model.AdvertCancelExpiry{
			CanceledAt: &canceledAt,
			ExpiredAt:  &expiredAt,
		}

		mockRepo.EXPECT().GetAdvertCancelExpiry(ctx, ID).Return(&expectedCancelExpiry, nil)
		mockRepo.EXPECT().RestoreAdvert(ctx, ID, gomock.Any()).Return(nil)

		s := New(mockRepo)
		result, err := s.RestoreAdvert(ctx, &advertproto.RestoreAdvertIn{Id: ID})

		assert.NoError(t, err)
		assert.Equal(t, result, &advertproto.AdvertEmpty{})
	})

	t.Run("should_return_err_cancel_expiry", func(t *testing.T) {
		expectedErr := errors.New("err")
		mockRepo.EXPECT().GetAdvertCancelExpiry(ctx, ID).Return(nil, expectedErr)

		s := New(mockRepo)
		_, err := s.RestoreAdvert(ctx, &advertproto.RestoreAdvertIn{Id: ID})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "err")
	})

	t.Run("should_return_err_nil_canceled_at", func(t *testing.T) {
		baseTime := time.Date(2025, 3, 4, 21, 0, 0, 0, time.UTC)
		expiredAt := baseTime.Add(1 * time.Hour)

		expectedCancelExpiry := model.AdvertCancelExpiry{
			CanceledAt: nil,
			ExpiredAt:  &expiredAt,
		}

		mockRepo.EXPECT().GetAdvertCancelExpiry(ctx, ID).Return(&expectedCancelExpiry, nil)

		s := New(mockRepo)
		_, err := s.RestoreAdvert(ctx, &advertproto.RestoreAdvertIn{Id: ID})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to advert wasn't canceled")
	})

	t.Run("should_return_err_restore_advert", func(t *testing.T) {
		expectedErr := errors.New("err")

		baseTime := time.Date(2025, 3, 4, 21, 0, 0, 0, time.UTC)
		canceledAt := baseTime
		expiredAt := baseTime.Add(1 * time.Hour)

		expectedCancelExpiry := model.AdvertCancelExpiry{
			CanceledAt: &canceledAt,
			ExpiredAt:  &expiredAt,
		}

		mockRepo.EXPECT().GetAdvertCancelExpiry(ctx, ID).Return(&expectedCancelExpiry, nil)

		mockRepo.EXPECT().RestoreAdvert(ctx, ID, gomock.Any()).Return(expectedErr)

		s := New(mockRepo)
		_, err := s.RestoreAdvert(ctx, &advertproto.RestoreAdvertIn{Id: ID})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "err")
	})
}

func TestService_CancelAdvert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)

	t.Run("cancel_ok", func(t *testing.T) {
		mockRepo.EXPECT().CancelAdvert(ctx, gomock.Any()).Return(nil)

		s := New(mockRepo)
		_, err := s.CancelAdvert(ctx, &advertproto.CancelAdvertIn{})
		assert.NoError(t, err)
	})

	t.Run("cancel_error", func(t *testing.T) {
		expectedErr := errors.New("cancel err")

		mockRepo.EXPECT().CancelAdvert(ctx, gomock.Any()).Return(expectedErr)

		s := New(mockRepo)
		_, err := s.CancelAdvert(ctx, &advertproto.CancelAdvertIn{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to cancel advert: cancel err")
	})
}

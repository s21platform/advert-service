package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	logger_lib "github.com/s21platform/logger-lib"

	"github.com/s21platform/advert-service/internal/config"
	"github.com/s21platform/advert-service/internal/model"
	advertproto "github.com/s21platform/advert-service/pkg/advert"
)

func TestServer_GetAdverts(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)
	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("get_ok", func(t *testing.T) {
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

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

		mockLogger.EXPECT().AddFuncName("GetAdverts")
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

		mockLogger.EXPECT().AddFuncName("GetAdverts")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		mockLogger.EXPECT().Error("failed to find uuid")

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

		mockLogger.EXPECT().AddFuncName("GetAdverts")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to find adverts: %v", expectedErr))

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

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("create_ok", func(t *testing.T) {
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockRepo.EXPECT().CreateAdvert(ctx, uuid, gomock.Any()).Return(nil)
		mockLogger.EXPECT().AddFuncName("CreateAdvert")

		s := New(mockRepo)
		_, err := s.CreateAdvert(ctx, &advertproto.CreateAdvertIn{})
		assert.NoError(t, err)
	})

	t.Run("create_no_uuid", func(t *testing.T) {
		ctx := context.Background()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := NewMockDBRepo(ctrl)

		mockLogger.EXPECT().AddFuncName("CreateAdvert")
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)
		mockLogger.EXPECT().Error("failed to find uuid")

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
		mockLogger.EXPECT().AddFuncName("CreateAdvert")
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to create advert: %v", expectedErr))

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

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_return_ok", func(t *testing.T) {
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		baseTime := time.Date(2025, 3, 4, 21, 0, 0, 0, time.UTC)
		canceledAt := baseTime
		expiredAt := baseTime.Add(1 * time.Hour)

		expectedCancelExpiry := model.AdvertCancelExpiry{
			IsCanceled: true,
			CanceledAt: &canceledAt,
			ExpiredAt:  &expiredAt,
		}

		mockLogger.EXPECT().AddFuncName("RestoreAdvert").Times(1)
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

		mockLogger.EXPECT().AddFuncName("RestoreAdvert").Times(1)
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to get advert cancel info: %v", expectedErr))

		s := New(mockRepo)
		_, err := s.RestoreAdvert(ctx, &advertproto.RestoreAdvertIn{Id: ID})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "err")
	})

	t.Run("should_return_err_was_not_canceled", func(t *testing.T) {
		baseTime := time.Date(2025, 3, 4, 21, 0, 0, 0, time.UTC)
		expiredAt := baseTime.Add(1 * time.Hour)

		expectedCancelExpiry := model.AdvertCancelExpiry{
			IsCanceled: false,
			CanceledAt: nil,
			ExpiredAt:  &expiredAt,
		}

		mockRepo.EXPECT().GetAdvertCancelExpiry(ctx, ID).Return(&expectedCancelExpiry, nil)

		mockLogger.EXPECT().AddFuncName("RestoreAdvert").Times(1)
		mockLogger.EXPECT().Error("failed to restore the advert due to a missing cancellation record")

		s := New(mockRepo)
		_, err := s.RestoreAdvert(ctx, &advertproto.RestoreAdvertIn{Id: ID})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to restore the advert due to a missing cancellation record")
	})

	t.Run("should_return_err_restore_advert", func(t *testing.T) {
		expectedErr := errors.New("err")

		baseTime := time.Date(2025, 3, 4, 21, 0, 0, 0, time.UTC)
		canceledAt := baseTime
		expiredAt := baseTime.Add(1 * time.Hour)

		expectedCancelExpiry := model.AdvertCancelExpiry{
			IsCanceled: true,
			CanceledAt: &canceledAt,
			ExpiredAt:  &expiredAt,
		}

		mockRepo.EXPECT().GetAdvertCancelExpiry(ctx, ID).Return(&expectedCancelExpiry, nil)
		mockRepo.EXPECT().RestoreAdvert(ctx, ID, gomock.Any()).Return(expectedErr)

		mockLogger.EXPECT().AddFuncName("RestoreAdvert").Times(1)
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to restore advert: %v", expectedErr))

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

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("cancel_ok", func(t *testing.T) {
		ctx = context.WithValue(ctx, config.KeyLogger, mockLogger)

		mockLogger.EXPECT().AddFuncName("CancelAdvert").Times(1)
		mockRepo.EXPECT().CancelAdvert(ctx, gomock.Any()).Return(nil)

		s := New(mockRepo)
		_, err := s.CancelAdvert(ctx, &advertproto.CancelAdvertIn{})
		assert.NoError(t, err)
	})

	t.Run("cancel_error", func(t *testing.T) {
		expectedErr := errors.New("cancel err")

		mockRepo.EXPECT().CancelAdvert(ctx, gomock.Any()).Return(expectedErr)

		mockLogger.EXPECT().AddFuncName("CancelAdvert").Times(1)
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to cancel advert: %v", expectedErr))

		s := New(mockRepo)
		_, err := s.CancelAdvert(ctx, &advertproto.CancelAdvertIn{})

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), "failed to cancel advert: cancel err")
	})
}

func TestServer_EditAdvert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ID := int32(123)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)

	mockLogger := logger_lib.NewMockLoggerInterface(ctrl)

	t.Run("should_return_ok", func(t *testing.T) {
		testCtx := context.WithValue(ctx, config.KeyLogger, mockLogger)
		testCtx = context.WithValue(testCtx, config.KeyUUID, "user123")

		input := &advertproto.EditAdvertIn{
			Id:          ID,
			TextContent: "updated content",
			UserFilter:  &advertproto.UserFilter{Os: []int64{22}},
		}

		mockLogger.EXPECT().AddFuncName("EditAdvert").Times(1)
		mockRepo.EXPECT().IsAdvertActive(testCtx, int(ID)).Return(true, nil)
		mockRepo.EXPECT().GetOwnerUUID(testCtx, int(ID)).Return("user123", nil)
		mockRepo.EXPECT().EditAdvert(testCtx, gomock.Any()).DoAndReturn(func(_ context.Context, advert *model.EditAdvert) error {
			assert.Equal(t, int(ID), advert.ID)
			assert.Equal(t, "updated content", advert.TextContent)
			assert.Equal(t, []int64{22}, advert.UserFilter.Os)
			return nil
		})

		s := New(mockRepo)
		result, err := s.EditAdvert(testCtx, input)

		assert.NoError(t, err)
		assert.Equal(t, &advertproto.AdvertEmpty{}, result)
	})

	t.Run("should_return_err_advert_active_check", func(t *testing.T) {
		testCtx := context.WithValue(ctx, config.KeyLogger, mockLogger)
		expectedErr := errors.New("db error")

		input := &advertproto.EditAdvertIn{Id: ID}

		mockLogger.EXPECT().AddFuncName("EditAdvert").Times(1)
		mockRepo.EXPECT().IsAdvertActive(testCtx, int(ID)).Return(false, expectedErr)
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to advert is not active: %v", expectedErr))

		s := New(mockRepo)
		_, err := s.EditAdvert(testCtx, input)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), expectedErr.Error())
	})

	t.Run("should_return_err_advert_not_active", func(t *testing.T) {
		testCtx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		input := &advertproto.EditAdvertIn{Id: ID}

		mockLogger.EXPECT().AddFuncName("EditAdvert").Times(1)
		mockRepo.EXPECT().IsAdvertActive(testCtx, int(ID)).Return(false, nil)
		mockLogger.EXPECT().Error("failed to advert is not active")

		s := New(mockRepo)
		_, err := s.EditAdvert(testCtx, input)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unavailable, st.Code())
		assert.Contains(t, st.Message(), "failed to advert is not active")
	})

	t.Run("should_return_err_missing_uuid", func(t *testing.T) {
		testCtx := context.WithValue(ctx, config.KeyLogger, mockLogger)

		input := &advertproto.EditAdvertIn{Id: ID}

		mockLogger.EXPECT().AddFuncName("EditAdvert").Times(1)
		mockRepo.EXPECT().IsAdvertActive(testCtx, int(ID)).Return(true, nil)
		mockLogger.EXPECT().Error("failed to find uuid")

		s := New(mockRepo)
		_, err := s.EditAdvert(testCtx, input)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Unauthenticated, st.Code())
		assert.Contains(t, st.Message(), "failed to find uuid")
	})

	t.Run("should_return_err_get_owner_uuid", func(t *testing.T) {
		testCtx := context.WithValue(ctx, config.KeyLogger, mockLogger)
		testCtx = context.WithValue(testCtx, config.KeyUUID, "user123")
		expectedErr := errors.New("db error")

		input := &advertproto.EditAdvertIn{Id: ID}

		mockLogger.EXPECT().AddFuncName("EditAdvert").Times(1)
		mockRepo.EXPECT().IsAdvertActive(testCtx, int(ID)).Return(true, nil)
		mockRepo.EXPECT().GetOwnerUUID(testCtx, int(ID)).Return("", expectedErr)
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to get owner uuid: %v", expectedErr))

		s := New(mockRepo)
		_, err := s.EditAdvert(testCtx, input)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), expectedErr.Error())
	})

	t.Run("should_return_err_not_owner", func(t *testing.T) {
		testCtx := context.WithValue(ctx, config.KeyLogger, mockLogger)
		testCtx = context.WithValue(testCtx, config.KeyUUID, "user123")

		input := &advertproto.EditAdvertIn{Id: ID}

		mockLogger.EXPECT().AddFuncName("EditAdvert").Times(1)
		mockRepo.EXPECT().IsAdvertActive(testCtx, int(ID)).Return(true, nil)
		mockRepo.EXPECT().GetOwnerUUID(testCtx, int(ID)).Return("different_user", nil)
		mockLogger.EXPECT().Error("failed to user is not advert owner")

		s := New(mockRepo)
		_, err := s.EditAdvert(testCtx, input)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, st.Code())
		assert.Contains(t, st.Message(), "failed to user is not advert owner")
	})

	t.Run("should_return_err_edit_advert", func(t *testing.T) {
		testCtx := context.WithValue(ctx, config.KeyLogger, mockLogger)
		testCtx = context.WithValue(testCtx, config.KeyUUID, "user123")
		expectedErr := errors.New("db error")

		input := &advertproto.EditAdvertIn{
			Id:          ID,
			TextContent: "updated content",
			UserFilter:  &advertproto.UserFilter{Os: []int64{22}},
		}

		mockLogger.EXPECT().AddFuncName("EditAdvert").Times(1)
		mockRepo.EXPECT().IsAdvertActive(testCtx, int(ID)).Return(true, nil)
		mockRepo.EXPECT().GetOwnerUUID(testCtx, int(ID)).Return("user123", nil)
		mockRepo.EXPECT().EditAdvert(testCtx, gomock.Any()).Return(expectedErr)
		mockLogger.EXPECT().Error(fmt.Sprintf("failed to edit advert: %v", expectedErr))

		s := New(mockRepo)
		_, err := s.EditAdvert(testCtx, input)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Contains(t, st.Message(), expectedErr.Error())
	})
}

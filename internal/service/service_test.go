package service

import (
	"context"
	"github.com/golang/mock/gomock"
	advertproto "github.com/s21platform/advert-proto/advert-proto"
	"github.com/s21platform/advert-service/internal/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestServer_CreateAdvert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	uuid := "test-uuid"
	ctx = context.WithValue(ctx, config.KeyUUID, uuid)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepo := NewMockDBRepo(ctrl)

	t.Run("create_ok", func(t *testing.T) {
		mockRepo.EXPECT().CreateAdvert(uuid, gomock.Any()).Return(nil)

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
}

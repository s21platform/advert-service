package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	advert_proto "github.com/s21platform/advert-service/pkg/advert"
)

type AdvertInfoList []*AdvertInfo

type AdvertInfo struct {
	ID        int64     `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"text_content"`
	ExpiredAt time.Time `db:"expired_at"`
}

func (a *AdvertInfo) FromDTO() *advert_proto.AdvertText {
	return &advert_proto.AdvertText{
		Id:          a.ID,
		Title:       a.Title,
		TextContent: a.Content,
		ExpiredAt:   timestamppb.New(a.ExpiredAt),
	}
}

func (a *AdvertInfoList) ListFromDTO() []*advert_proto.AdvertText {
	result := make([]*advert_proto.AdvertText, 0, len(*a))

	for _, advert := range *a {
		result = append(result, advert.FromDTO())
	}

	return result
}

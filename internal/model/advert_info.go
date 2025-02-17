package model

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	advert_proto "github.com/s21platform/advert-proto/advert-proto"
)

type AdvertInfoList []*AdvertInfo

type AdvertInfo struct {
	Content   string    `db:"text_content"`
	ExpiredAt time.Time `db:"expired_at"`
}

func (a *AdvertInfoList) FromDTO() []*advert_proto.AdvertText {
	result := make([]*advert_proto.AdvertText, 0, len(*a))

	for _, advert := range *a {
		result = append(result, &advert_proto.AdvertText{
			TextContent: advert.Content,
			ExpiredAt:   timestamppb.New(advert.ExpiredAt),
		})
	}

	return result
}

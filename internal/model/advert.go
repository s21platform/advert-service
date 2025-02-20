package model

import (
	"encoding/json"
	"time"

	advert "github.com/s21platform/advert-proto/advert-proto"
)

type Advert struct {
	OwnerUUID   string    `db:"owner_uuid"`
	TextContent string    `db:"text_content"`
	UserFilter  []byte    `db:"filter"`
	ExpiresAt   time.Time `db:"expired_at"`
}

type UserFilter struct {
	Os string `json:"os"`
}

func (a *Advert) ToDTO(UUID string, in *advert.CreateAdvertIn) (Advert, error) {
	userFilterJSON, err := json.Marshal(in.User)
	if err != nil {
		userFilterJSON = nil
	}

	result := Advert{
		OwnerUUID:   UUID,
		TextContent: in.Text,
		UserFilter:  []byte(userFilterJSON),
		ExpiresAt:   in.ExpiredAt.AsTime(),
	}

	return result, nil
}

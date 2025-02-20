package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	advert "github.com/s21platform/advert-proto/advert-proto"
)

type Advert struct {
	OwnerUUID   string     `db:"owner_uuid"`
	TextContent string     `db:"text_content"`
	UserFilter  UserFilter `db:"filter"`
	ExpiresAt   time.Time  `db:"expired_at"`
}

type UserFilter struct {
	Os string `json:"os"`
}

func (uf UserFilter) Value() (driver.Value, error) {
	j, err := json.Marshal(uf)
	if err != nil {
		return "", err
	}
	return string(j), nil
}

func (uf UserFilter) Scan(value interface{}) error {
	b, isBytes := value.([]byte)
	if !isBytes {
		s, isString := value.(string)
		if !isString {
			return errors.New("failed to Scan lot.data field, supported types: `string` or `[]byte`")
		}
		b = []byte(s)
	}

	return json.Unmarshal(b, &uf)
}

func (a *Advert) ToDTO(UUID string, in *advert.CreateAdvertIn) (Advert, error) {
	result := Advert{
		OwnerUUID:   UUID,
		TextContent: in.Text,
		UserFilter:  UserFilter{Os: in.User.Os},
		ExpiresAt:   in.ExpiredAt.AsTime(),
	}

	return result, nil
}

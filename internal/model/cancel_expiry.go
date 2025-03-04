package model

import "time"

type AdvertCancelExpiry struct {
	CanceledAt *time.Time `db:"canceled_at"`
	ExpiredAt  *time.Time `db:"expired_at"`
}

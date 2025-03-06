package model

import "time"

type AdvertCancelExpiry struct {
	IsCanceled bool       `db:"is_canceled"`
	CanceledAt *time.Time `db:"canceled_at"`
	ExpiredAt  *time.Time `db:"expired_at"`
}

package model

import "database/sql"

type AdvertState struct {
	IsCanceled bool         `db:"is_canceled"`
	IsBanned   bool         `db:"is_banned"`
	ExpiredAt  sql.NullTime `db:"expired_at"`
}

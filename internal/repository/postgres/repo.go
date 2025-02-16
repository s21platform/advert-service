package postgres

import (
	"fmt"
	advert "github.com/s21platform/advert-proto/advert-proto"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	connection *sqlx.DB
}

func (r *Repository) CreateAdvert(UUID string, in *advert.CreateAdvertIn) error {
	query := `
		INSERT INTO advert_text (owner_uuid, text_content, filter, expired_at)
		VALUES ($1, $2, $3, $4)`
	_, err := r.connection.Exec(query, UUID, in.Text, in.User, in.ExpiredAt)
	if err != nil {
		return fmt.Errorf("failed to create advert: %v", err)
	}

	return nil
}

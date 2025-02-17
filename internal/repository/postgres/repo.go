package postgres

import (
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/s21platform/advert-service/internal/config"
	"github.com/s21platform/advert-service/internal/model"
)

type Repository struct {
	connection *sqlx.DB
}

func New(cfg *config.Config) *Repository {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	conn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Fatal("error connect: ", err)
	}

	return &Repository{
		connection: conn,
	}
}

func (r *Repository) Close() {
	_ = r.connection.Close()
}

func (r *Repository) GetAdverts(UUID string) (*model.AdvertInfoList, error) {
	var adverts model.AdvertInfoList

	query := squirrel.Select("text_content", "expired_at").
		From("advert_text").
		Where(squirrel.Eq{"owner_uuid": UUID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	err = r.connection.Select(&adverts, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get adverts from db: %v", err)
	}

	return &adverts, nil
}

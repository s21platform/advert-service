package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	advert "github.com/s21platform/advert-proto/advert-proto"

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

func (r *Repository) CreateAdvert(ctx context.Context, UUID string, in *advert.CreateAdvertIn) error {
	var advertObj model.Advert

	advertObj, err := advertObj.ToDTO(UUID, in)
	if err != nil {
		return fmt.Errorf("failed toconvert grpc message to dto: %v", err)
	}

	query := squirrel.Insert("advert_text").
		Columns("owner_uuid", "text_content", "filter", "expired_at").
		Values(advertObj.OwnerUUID, advertObj.TextContent, advertObj.UserFilter, advertObj.ExpiresAt).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()

	if err != nil {
		return fmt.Errorf("failed to build SQL query: %v", err)
	}

	_, err = r.connection.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to create advert: %v", err)
	}

	return nil
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

func (r *Repository) RestoreAdvert(ID int64) error {
	var canceledAt time.Time
	var expiredAt time.Time

	selectQuery := squirrel.
		Select("canceled_at", "expired_at").
		From("advert_text").
		Where(squirrel.Eq{"id": ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := selectQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build select query: %v", err)
	}

	err = r.connection.QueryRow(sql, args...).Scan(&canceledAt, &expiredAt)
	if err != nil {
		return fmt.Errorf("failed to get advert data: %v", err)
	}

	timeDiff := time.Since(canceledAt)
	newExpiresAt := expiredAt.Add(timeDiff)

	updateQuery := squirrel.
		Update("advert_text").
		Set("is_canceled", false).
		Set("canceled_at", nil).
		Set("expired_at", newExpiresAt).
		Where(squirrel.Eq{"id": ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err = updateQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	_, err = r.connection.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update advert: %v", err)
	}

	return nil
}

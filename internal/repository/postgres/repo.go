package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/s21platform/advert-service/internal/config"
	"github.com/s21platform/advert-service/internal/model"
	advert_api "github.com/s21platform/advert-service/pkg/advert"
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

func (r *Repository) CreateAdvert(ctx context.Context, UUID string, in *advert_api.CreateAdvertIn) error {
	var advertObj model.Advert

	advertObj, err := advertObj.AdvertToDTO(UUID, in)
	if err != nil {
		return fmt.Errorf("failed toconvert grpc message to dto: %v", err)
	}

	query := squirrel.Insert("advert_text").
		Columns("owner_uuid", "title", "text_content", "filter", "expired_at").
		Values(advertObj.OwnerUUID, advertObj.Title, advertObj.TextContent, advertObj.UserFilter, advertObj.ExpiresAt).
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

func (r *Repository) GetAdvert(ctx context.Context, in *advert_api.GetAdvertIn) (*model.AdvertInfo, error) {
	var advert model.AdvertInfo

	query, args, err := squirrel.Select("id", "title", "text_content", "expired_at").
		From("advert_text").
		Where(squirrel.Eq{"id": in.Id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	err = r.connection.SelectContext(ctx, &advert, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get advert from db: %v", err)
	}

	return &advert, nil
}

func (r *Repository) GetAdverts(UUID string) (*model.AdvertInfoList, error) {
	var adverts model.AdvertInfoList

	query := squirrel.Select("id", "title", "text_content", "expired_at").
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

func (r *Repository) CancelAdvert(ctx context.Context, in *advert_api.CancelAdvertIn) error {
	updateQuery := squirrel.Update("advert_text").
		Set("is_canceled", true).
		Set("canceled_at", time.Now()).
		Where(squirrel.And{
			squirrel.Eq{"id": in.Id},
			squirrel.Eq{"is_canceled": false},
			squirrel.Gt{"expired_at": time.Now().Unix()},
		}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := updateQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	_, err = r.connection.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to set cancel status in data: %v", err)
	}

	return nil
}

func (r *Repository) GetAdvertCancelExpiry(ctx context.Context, ID int64) (*model.AdvertCancelExpiry, error) {
	var cancelExpiry model.AdvertCancelExpiry

	query := squirrel.
		Select("is_canceled", "canceled_at", "expired_at").
		From("advert_text").
		Where(squirrel.Eq{"id": ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %v", err)
	}

	err = r.connection.GetContext(ctx, &cancelExpiry, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get advert data: %v", err)
	}

	return &cancelExpiry, nil
}

func (r *Repository) RestoreAdvert(ctx context.Context, ID int64, newExpiredAt time.Time) error {
	query := squirrel.
		Update("advert_text").
		Set("is_canceled", false).
		Set("canceled_at", nil).
		Set("expired_at", newExpiredAt).
		Where(squirrel.Eq{"id": ID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	_, err = r.connection.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update advert: %v", err)
	}

	return nil
}

func (r *Repository) IsAdvertActive(ctx context.Context, ID int) (bool, error) {
	query, args, err := squirrel.
		Select("is_canceled", "is_banned", "expired_at").
		From("advert_text").
		Where(squirrel.Eq{"id": ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build select query: %v", err)
	}

	result := model.AdvertState{}
	err = r.connection.GetContext(ctx, &result, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to query advert: %v", err)
	}

	isExpired := result.ExpiredAt.Valid && result.ExpiredAt.Time.Before(time.Now())
	return !result.IsCanceled && !result.IsBanned && !isExpired, nil
}

func (r *Repository) GetOwnerUUID(ctx context.Context, ID int) (string, error) {
	query, args, err := squirrel.
		Select("owner_uuid").
		From("advert_text").
		Where(squirrel.Eq{"id": ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("failed to build select query: %v", err)
	}

	var ownerUUID string
	err = r.connection.GetContext(ctx, &ownerUUID, query, args...)
	if err != nil {
		return "", fmt.Errorf("failed to query owner uuid: %v", err)
	}

	return ownerUUID, nil
}

func (r *Repository) EditAdvert(ctx context.Context, info *model.EditAdvert) error {
	query, args, err := squirrel.
		Update("advert_text").
		Set("text_content", info.TextContent).
		Set("title", info.Title).
		Set("filter", info.UserFilter).
		Where(squirrel.Eq{"id": info.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build update query: %v", err)
	}

	_, err = r.connection.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update advert: %v", err)
	}

	return nil
}

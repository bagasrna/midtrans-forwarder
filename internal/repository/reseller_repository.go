package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"midtrans-forwarder/internal/domain"
    redisHelper "midtrans-forwarder/pkg/redis"

	"github.com/redis/go-redis/v9"
)

type ResellerRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewResellerRepository(db *sql.DB, redis *redis.Client) *ResellerRepository {
	return &ResellerRepository{db: db, redis: redis}
}

func (r *ResellerRepository) CreateReseller(ctx context.Context, reseller *domain.Reseller) error {
	query := "INSERT INTO resellers (name, code, url, token) VALUES (?, ?, ?)"
	result, err := r.db.ExecContext(ctx, query, reseller.Name, reseller.Code, reseller.URL, reseller.Token)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	reseller.ID = id
	r.redis.Del(ctx, "resellers")
	return nil
}

func (r *ResellerRepository) GetResellerByID(ctx context.Context, id int64) (*domain.Reseller, error) {
	cacheKey := fmt.Sprintf("reseller:%d", id)
	cachedData, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var reseller domain.Reseller
		err = json.Unmarshal([]byte(cachedData), &reseller)
		if err == nil {
			return &reseller, nil
		}
	}

	query := "SELECT id, name, code, url, token FROM resellers WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var reseller domain.Reseller
	err = row.Scan(&reseller.ID, &reseller.Name, &reseller.Code, &reseller.URL, &reseller.Token)
	if err != nil {
		return nil, err
	}

	return &reseller, nil
}

func (r *ResellerRepository) GetAllResellers(ctx context.Context) ([]domain.Reseller, error) {
	// Cek cache terlebih dahulu
    cacheName := redisHelper.BuildCacheKey("resellers")
	cachedData, err := r.redis.Get(ctx, cacheName).Result()
	if err == nil {
		var resellers []domain.Reseller
		err = json.Unmarshal([]byte(cachedData), &resellers)
		if err == nil {
			return resellers, nil
		}
	}

	// Jika data tidak ada di cache, lakukan query ke database
	query := "SELECT id, name, code, url, token FROM resellers"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resellers []domain.Reseller
	for rows.Next() {
		var reseller domain.Reseller
		if err := rows.Scan(&reseller.ID, &reseller.Name, &reseller.Code, &reseller.URL, &reseller.Token); err != nil {
			return nil, err
		}
		resellers = append(resellers, reseller)
	}

	// Simpan hasil query ke cache
	if len(resellers) > 0 {
		data, err := json.Marshal(resellers)
		if err == nil {
			err = r.redis.Set(ctx, cacheName, data, 0).Err()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return resellers, nil
}

func (r *ResellerRepository) UpdateReseller(ctx context.Context, reseller *domain.Reseller) error {
	query := "UPDATE resellers SET name = ?, code = ?, url = ?, token = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, &reseller.Name, &reseller.Code, &reseller.URL, &reseller.Token, reseller.ID)
	if err != nil {
		return err
	}

	r.redis.Del(ctx, fmt.Sprintf("reseller:%d", reseller.ID))
	r.redis.Del(ctx, "resellers")

	return err
}

func (r *ResellerRepository) DeleteReseller(ctx context.Context, id int64) error {
	query := "DELETE FROM resellers WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	r.redis.Del(ctx, fmt.Sprintf("reseller:%d", id))
	r.redis.Del(ctx, "resellers")

	return err
}

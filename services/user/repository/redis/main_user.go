package redis

import (
	"time"
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/shellrean/extraordinary-raport/domain"
)

type redisUserRepository struct {
	Conn *redis.Client
}

func NewRedisUserRepository(Conn *redis.Client) domain.UserCacheRepository {
	return &redisUserRepository {
		Conn,
	}
}

func (m *redisUserRepository) StoreAuth(ctx context.Context, u domain.User, td *domain.TokenDetails) (err error) {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	err = m.Conn.Set(ctx, td.AccessUuid, u.ID, at.Sub(now)).Err()
	if err != nil {
		return domain.ErrServerError
	}
	err = m.Conn.Set(ctx, td.RefreshUuid, u.ID, rt.Sub(now)).Err()
	if err != nil {
		return domain.ErrServerError
	}

	return
}

func (m *redisUserRepository) DeleteAuth(ctx context.Context, uuid string) (err error) {
	_, err = m.Conn.Del(ctx, uuid).Result()
	if err != nil {
		return domain.ErrServerError
	}
	return
}
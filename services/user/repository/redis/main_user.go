package redis

import (
	"time"
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/shellrean/extraordinary-raport/domain"
)

type repository struct {
	Conn *redis.Client
}

func New(Conn *redis.Client) domain.UserCacheRepository {
	return &repository {
		Conn,
	}
}

func (m *repository) StoreAuth(ctx context.Context, u domain.User, td *domain.Token) (err error) {
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

func (m *repository) DeleteAuth(ctx context.Context, uuid string) (err error) {
	_, err = m.Conn.Del(ctx, uuid).Result()
	if err != nil {
		return domain.ErrServerError
	}
	return
}
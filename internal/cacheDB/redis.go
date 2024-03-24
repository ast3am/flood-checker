package cacheDB

import (
	"context"
	"errors"
	"github.com/ast3am/flood-checker/internal/config"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

type CahceDB struct {
	client *redis.Client
}

func NeWRedisClient(ctx context.Context, cfg *config.RedisConfig) (*CahceDB, error) {
	client := CahceDB{}
	client.client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.Password,
		DB:       cfg.DBName,
	})
	_, err := client.client.Ping(ctx).Result()
	if err != nil {
		return nil, errors.New("fail connecting to Redis")
	}
	return &client, nil
}

func (db *CahceDB) GetRequestCount(ctx context.Context, userID int64, interval time.Duration) (int64, error) {
	now := time.Now()
	key := "flood_control:" + strconv.FormatInt(userID, 10)

	pipe := db.client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(now.Add(interval).UnixNano(), 10))
	pipe.ZAdd(ctx, key, &redis.Z{Score: float64(now.UnixNano()), Member: now.UnixNano()})
	pipe.Expire(ctx, key, interval)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	count, err := pipe.ZCard(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return count, nil
}

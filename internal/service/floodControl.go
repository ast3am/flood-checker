package service

import (
	"context"
	"github.com/ast3am/flood-checker/internal/config"
	"sync"
	"time"
)

type db interface {
	GetRequestCount(ctx context.Context, userID int64, interval time.Duration) (int64, error)
}

type FloodControl struct {
	cacheDB     db
	maxRequests int
	interval    time.Duration
	mu          sync.Mutex
}

func NewFloodControl(client db, cfg *config.FloodControlCfg) *FloodControl {
	return &FloodControl{
		cacheDB:     client,
		maxRequests: cfg.CheckCount,
		interval:    time.Duration(cfg.CheckTime) * time.Second,
	}
}

func (fc *FloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	count, err := fc.cacheDB.GetRequestCount(ctx, userID, fc.interval)
	if err != nil {
		return false, err
	}
	if count > int64(fc.maxRequests) {
		return false, nil
	}
	return true, nil
}

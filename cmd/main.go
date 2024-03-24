package main

import (
	"context"
	"fmt"
	"github.com/ast3am/flood-checker/internal/cacheDB"
	"github.com/ast3am/flood-checker/internal/config"
	"github.com/ast3am/flood-checker/internal/service"
	"log"
)

func main() {
	ctx := context.Background()
	cacheCfg := config.GetRedisConfig("config/redisConfig.yml")
	fcCfg := config.GetFloodControlConfig("config/floodControlConfig.yml")
	db, err := cacheDB.NeWRedisClient(ctx, cacheCfg)
	if err != nil {
		log.Fatal("", err)
	}
	floodControl := service.NewFloodControl(db, fcCfg)
	banned, err := floodControl.Check(ctx, 15)
	if err != nil {
		log.Printf("chekker error: %s\n", err.Error())
	}
	if !banned {
		fmt.Println("User flooder")
	}
}

// FloodControl интерфейс, который нужно реализовать.
// Рекомендуем создать директорию-пакет, в которой будет находиться реализация.
type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

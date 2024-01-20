package main

import (
	"fmt"

	"github.com/SaidovZohid/test-task-crud/api"
	_ "github.com/SaidovZohid/test-task-crud/api/docs"
	"github.com/SaidovZohid/test-task-crud/config"
	"github.com/SaidovZohid/test-task-crud/pkg/logger"
	"github.com/SaidovZohid/test-task-crud/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func main() {
	logger.Init()
	log := logger.GetLogger()
	log.Info("logger initialized")

	cfg := config.Load(".")

	psqlUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
	)

	psqlConn, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println("--> Successfully connected to PostgreSQL <--")

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis,
	})

	strg := storage.NewStoragePg(psqlConn)
	inMemory := storage.NewInMemoryStorage(rdb)

	app := api.New(&api.RoutetOptions{
		Cfg:      &cfg,
		Log:      log,
		Storage:  strg,
		InMemory: inMemory,
	})

	log.Info("HTTP running in PORT -> ", cfg.HttpPort)
	log.Fatal("Error while listening http port:", app.Run(cfg.HttpPort))
}

package main

import (
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"

	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			Gopher Social Network API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		petstore.swagger.io
// @BasePath	/v1
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8081"),
		db: dbConfig{
			addr:            env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns:    env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns:    env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:     env.GetString("DB_MAX_IDLE_TIME", "15m"),
			maxConnLifetime: env.GetString("DB_MAX_CONN_LIFETIME", "15m"),
		},
		environment: env.GetString("ENV", "development"),
		apiURL:      env.GetString("API_URL", "http://localhost:8081"),
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := db.New(cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
		cfg.db.maxConnLifetime)

	if err != nil {
		logger.Fatal(err.Error())
	}

	defer db.Close()
	logger.Info("Database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger.Sugar(),
	}

	mux := app.mount()
	logger.Fatal(app.run(mux).Error())
}

package main

import (
	"log"
	"social/internal/db"
	"social/internal/env"
	"social/internal/store"
)

const version = "0.0.1"

func main() {

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:            env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns:    env.GetInt("DB_MAX_OPEN_CONNS", 25),
			maxIdleConns:    env.GetInt("DB_MAX_IDLE_CONNS", 25),
			maxIdleTime:     env.GetString("DB_MAX_IDLE_TIME", "15m"),
			maxConnLifetime: env.GetString("DB_MAX_CONN_LIFETIME", "15m"),
		},
		environment: env.GetString("ENV", "development"),
	}

	db, err := db.New(cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
		cfg.db.maxConnLifetime)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("Database connection established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}

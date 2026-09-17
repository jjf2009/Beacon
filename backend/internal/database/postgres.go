package database

import (
    "database/sql"
    "fmt"
    "time"
     _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/jjf2009/beacon/backend/internal/config"
)

func New(cfg *config.Config) (*sql.DB, error) {
    // Step 1: build DSN string from cfg.Database fields
    // "host=X port=X user=X password=X dbname=X sslmode=disable"
  dsn := fmt.Sprintf(
    "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    cfg.Database.Host,
    cfg.Database.Port,
    cfg.Database.User,
    cfg.Database.Password,
    cfg.Database.Name,
)


    // Step 2: sql.Open — validates DSN, no real connection yet
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        return nil, err
    }

    // Step 3: db.Ping — actually connects, fails if DB unreachable
    if err := db.Ping(); err != nil {
        return nil, err
    }

    // Step 4: set pool limits
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}

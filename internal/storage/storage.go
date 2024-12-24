package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/wavinamayola/user-management/internal/config"
	"github.com/wavinamayola/user-management/internal/utils"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Config) (*Storage, error) {
	dataSource, err := utils.NewDBStringFromDBConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to set data source: %+v", err)
	}

	db, err := sql.Open(cfg.Database.Driver, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed initialize mysql storage: %+v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to mysql storage, check db config")
	}

	return &Storage{db: db}, nil
}

func (s *Storage) GetDB() *sql.DB {
	return s.db
}

package db

import (
	"fmt"
	"time"

	"arvanch/config"

	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres driver should have blank import

	"github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

const HealthCheckPeriod = 1 * time.Second

func WithRetry(fn func(cfg config.Postgres) (*gorm.DB, error), cfg config.Postgres) *gorm.DB {
	const maxAttempts = 60
	for range maxAttempts {
		db, err := fn(cfg)
		if err == nil {
			return db
		}

		logrus.Errorf("Could not connect to DB. Waiting 1 second. Reason is => %s", err.Error())
		<-time.After(HealthCheckPeriod)
	}

	panic(fmt.Sprintf("Could not connect to postgres after %d attempts", maxAttempts))
}

func Create(cfg config.Postgres) (*gorm.DB, error) {
	url := "host=127.0.0.1 port=5432 dbname=arvanch password=postgres sslmode=disable"

	db, err := gorm.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	db.DB().SetConnMaxLifetime(cfg.ConnectionLifetime)
	db.DB().SetMaxOpenConns(cfg.MaxOpenConnections)
	db.DB().SetMaxIdleConns(cfg.MaxIdleConnections)

	return db, nil
}

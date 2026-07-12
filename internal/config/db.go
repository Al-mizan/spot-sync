package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	// db, err := gorm.Open(postgres.New(postgres.Config{
	// 	DSN:                  dsn,
	// 	PreferSimpleProtocol: true, // disables implicit prepared statement usage, fixing PgBouncer transaction pooling issues
	// }), &gorm.Config{
	// 	TranslateError: false,
	// })
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: cfg.Dsn,
	}), &gorm.Config{
		TranslateError: false,
	})

	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(15 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Println("Database connection successful")
	return db, nil
}

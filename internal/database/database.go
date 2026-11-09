package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/config"
	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/domain"
)

func Connect(cfg config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("підключення до PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("отримання sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("перевірка з'єднання з PostgreSQL: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{}, &domain.Incident{})
}

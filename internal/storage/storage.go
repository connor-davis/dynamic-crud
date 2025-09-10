package storage

import (
	"github.com/connor-davis/dynamic-crud/common"
	"github.com/connor-davis/dynamic-crud/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	Database() *gorm.DB
	Migrate() error
}

type storage struct {
	db *gorm.DB
}

func NewStorage() Storage {
	dsn := common.EnvString("APP_DSN", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database: %s" + err.Error())
	}

	return &storage{
		db: db,
	}
}

func (s *storage) Database() *gorm.DB {
	return s.db
}

func (s *storage) Migrate() error {
	return s.db.AutoMigrate(
		&models.User{},
	)
}

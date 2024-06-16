package repo

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnection(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Address()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, err
}

func CloseConnection(db *gorm.DB) error {
	_db, err := db.DB()
	if err != nil {
		return err
	}

	return _db.Close()
}

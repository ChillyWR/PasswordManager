package repo

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenConnection(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Address()), &gorm.Config{})
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

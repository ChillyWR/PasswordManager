package repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/okutsen/PasswordManager/model/db"
)

func NewCredentialRecordRepository(db *gorm.DB) (*CredentialRecordRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	
	return &CredentialRecordRepository{db: db}, nil
}

type CredentialRecordRepository struct {
	db *gorm.DB
}

func (r *CredentialRecordRepository) GetAll() ([]db.CredentialRecord, error) {
	var records []db.CredentialRecord
	result := r.db.Find(&records)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all records from db: %w", err)
	}

	return records, nil
}
func (r *CredentialRecordRepository) Get(id uuid.UUID) (*db.CredentialRecord, error) {
	var record db.CredentialRecord
	result := r.db.First(&record, id)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to get record from db: %w", err)
	}

	return &record, nil
}

func (r *CredentialRecordRepository) Create(record *db.CredentialRecord) (*db.CredentialRecord, error) {
	record.ID = uuid.New()
	result := r.db.Create(record)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to create record in db: %w", err)
	}

	return record, nil
}

func (r *CredentialRecordRepository) Update(record *db.CredentialRecord) (*db.CredentialRecord, error) {
	result := r.db.Model(record).Clauses(clause.Returning{}).Updates(record)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to update record in db: %w", err)
	}

	return record, nil
}

func (r *CredentialRecordRepository) Delete(id uuid.UUID) (*db.CredentialRecord, error) {
	var record db.CredentialRecord
	result := r.db.Model(&record).Clauses(clause.Returning{}).Delete(&record, id)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to remove record from db: %w", err)
	}

	return &record, nil
}

package repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/okutsen/PasswordManager/model"
)

func NewRecordRepository(db *gorm.DB) (*RecordRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &RecordRepository{db: db}, nil
}

type RecordRepository struct {
	db *gorm.DB
}

func (r *RecordRepository) GetAll() ([]model.CredentialRecord, error) {
	var records []model.CredentialRecord
	result := r.db.Find(&records)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all records from db: %w", err)
	}

	return records, nil
}

func (r *RecordRepository) Get(id uuid.UUID) (*model.CredentialRecord, error) {
	var record model.CredentialRecord
	result := r.db.First(&record, id)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to get record from db: %w", err)
	}

	return &record, nil
}

func (r *RecordRepository) Create(record *model.CredentialRecord) (*model.CredentialRecord, error) {
	result := r.db.Create(record)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to create record in db: %w", err)
	}

	return record, nil
}

func (r *RecordRepository) Update(record *model.CredentialRecord) (*model.CredentialRecord, error) {
	// TODO: validate updates
	result := r.db.Model(record).Clauses(clause.Returning{}).Updates(record)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to update record in db: %w", err)
	}

	return record, nil
}

func (r *RecordRepository) Delete(id uuid.UUID) (*model.CredentialRecord, error) {
	var record model.CredentialRecord
	result := r.db.Model(&record).Clauses(clause.Returning{}).Delete(&record, id)
	err := result.Error
	if err != nil {
		return nil, fmt.Errorf("failed to remove record from db: %w", err)
	}

	return &record, nil
}

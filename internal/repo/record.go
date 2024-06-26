package repo

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/ChillyWR/PasswordManager/model"
	"github.com/ChillyWR/PasswordManager/pkg/pmerror"
)

type CredentialRecord model.CredentialRecord

func (CredentialRecord) TableName() string {
	return "credential_record"
}

type LoginRecord struct {
	ID       uuid.UUID
	Username *string `json:"username"`
	Password *string `json:"password"`
	URL      *string `json:"url"`
}

func (LoginRecord) TableName() string {
	return "login"
}

type CardRecord struct {
	ID              uuid.UUID
	Brand           *string `json:"brand"`
	Number          *string `json:"number"`
	ExpirationMonth *string `json:"expiration_month"`
	ExpirationYear  *string `json:"expiration_year"`
	CVV             *string `json:"cvv"`
}

func (CardRecord) TableName() string {
	return "card"
}

type IdentityRecord struct {
	ID             uuid.UUID
	FirstName      *string `json:"first_name"`
	MiddleName     *string `json:"middle_name"`
	LastName       *string `json:"last_name"`
	Address        *string `json:"address"`
	Email          *string `json:"email"`
	PhoneNumber    *string `json:"phone_number"`
	PassportNumber *string `json:"passport_number"`
	Country        *string `json:"country"`
}

func (IdentityRecord) TableName() string {
	return "identity"
}

func NewRecordRepository(db *gorm.DB) (*RecordRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}

	return &RecordRepository{db: db}, nil
}

type RecordRepository struct {
	db *gorm.DB
}

func (r *RecordRepository) GetAll(userID uuid.UUID) ([]model.CredentialRecord, []model.LoginRecord, []model.CardRecord, []model.IdentityRecord, error) {
	var credentialRecords []CredentialRecord
	if err := r.db.Where("created_by = ?", userID).Order("created_on, name, id").Find(&credentialRecords).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get credential records: %w", convertError(err))
	}

	var loginRecords []LoginRecord
	if err := r.db.Model(&LoginRecord{}).Order("created_on, name, id").
		Joins("INNER JOIN credential_record cd ON cd.id = login.id AND cd.created_by = ?", userID).
		Scan(&loginRecords).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get logins: %w", convertError(err))
	}

	var cardRecords []CardRecord
	if err := r.db.Model(&CardRecord{}).Order("created_on, name, id").
		Joins("INNER JOIN credential_record cd ON cd.id = card.id AND cd.created_by = ?", userID).
		Scan(&cardRecords).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get cards: %w", convertError(err))
	}

	var identityRecords []IdentityRecord
	if err := r.db.Model(&IdentityRecord{}).Order("created_on, name, id").
		Joins("INNER JOIN credential_record cd ON cd.id = identity.id AND cd.created_by = ?", userID).
		Scan(&identityRecords).Error; err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get identities: %w", convertError(err))
	}

	core := make(map[uuid.UUID]CredentialRecord, len(credentialRecords))
	for _, record := range credentialRecords {
		core[record.ID] = record
	}

	logins := make([]model.LoginRecord, len(loginRecords))
	for i, record := range loginRecords {
		logins[i] = model.LoginRecord{
			CredentialRecord: model.CredentialRecord(core[record.ID]),
			Username:         record.Username,
			Password:         record.Password,
			URL:              record.URL,
		}

		delete(core, record.ID)
	}

	cards := make([]model.CardRecord, len(cardRecords))
	for i, record := range cardRecords {
		cards[i] = model.CardRecord{
			CredentialRecord: model.CredentialRecord(core[record.ID]),
			Brand:            record.Brand,
			Number:           record.Number,
			ExpirationMonth:  record.ExpirationMonth,
			ExpirationYear:   record.ExpirationYear,
			CVV:              record.CVV,
		}

		delete(core, record.ID)
	}

	identities := make([]model.IdentityRecord, len(identityRecords))
	for i, record := range identityRecords {
		identities[i] = model.IdentityRecord{
			CredentialRecord: model.CredentialRecord(core[record.ID]),
			FirstName:        record.FirstName,
			MiddleName:       record.MiddleName,
			LastName:         record.LastName,
			Address:          record.Address,
			Email:            record.Email,
			PhoneNumber:      record.PhoneNumber,
			PassportNumber:   record.PassportNumber,
			Country:          record.Country,
		}

		delete(core, record.ID)
	}

	secureNotes := make([]model.CredentialRecord, 0, len(core))
	// preserve order
	for _, record := range credentialRecords {
		if _, ok := core[record.ID]; ok {
			secureNotes = append(secureNotes, model.CredentialRecord(record))
		}
	}

	return secureNotes, logins, cards, identities, nil
}

func (r *RecordRepository) GetCredentialRecord(id uuid.UUID) (*model.CredentialRecord, error) {
	var record CredentialRecord
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("get record: %w", convertError(err))
	}

	return (*model.CredentialRecord)(&record), nil
}

func (r *RecordRepository) GetLogin(id uuid.UUID) (*model.LoginRecord, error) {
	var core CredentialRecord
	if err := r.db.First(&core, id).Error; err != nil {
		return nil, fmt.Errorf("get core: %w", convertError(err))
	}

	var record LoginRecord
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("get login: %w", convertError(err))
	}

	return &model.LoginRecord{
		CredentialRecord: model.CredentialRecord(core),
		Username:         record.Username,
		Password:         record.Password,
		URL:              record.URL,
	}, nil
}

func (r *RecordRepository) GetCard(id uuid.UUID) (*model.CardRecord, error) {
	var core CredentialRecord
	if err := r.db.First(&core, id).Error; err != nil {
		return nil, fmt.Errorf("get core: %w", convertError(err))
	}

	var record CardRecord
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("get login: %w", convertError(err))
	}

	return &model.CardRecord{
		CredentialRecord: model.CredentialRecord(core),
		Brand:            record.Brand,
		Number:           record.Number,
		ExpirationMonth:  record.ExpirationMonth,
		ExpirationYear:   record.ExpirationYear,
		CVV:              record.CVV,
	}, nil
}

func (r *RecordRepository) GetIdentity(id uuid.UUID) (*model.IdentityRecord, error) {
	var core CredentialRecord
	if err := r.db.First(&core, id).Error; err != nil {
		return nil, fmt.Errorf("get core: %w", convertError(err))
	}

	var record IdentityRecord
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("get login: %w", convertError(err))
	}

	return &model.IdentityRecord{
		CredentialRecord: model.CredentialRecord(core),
		FirstName:        record.FirstName,
		MiddleName:       record.MiddleName,
		LastName:         record.LastName,
		Address:          record.Address,
		Email:            record.Email,
		PhoneNumber:      record.PhoneNumber,
		PassportNumber:   record.PassportNumber,
		Country:          record.Country,
	}, nil
}

func (r *RecordRepository) CreateCredentialRecord(record *model.CredentialRecord) (*model.CredentialRecord, error) {
	credentialRecord := CredentialRecord(*record)
	if err := r.db.Create(&credentialRecord).Error; err != nil {
		return nil, fmt.Errorf("create: %w", convertError(err))
	}

	record.ID = credentialRecord.ID

	return record, nil
}

func (r *RecordRepository) CreateLogin(record *model.LoginRecord) (*model.LoginRecord, error) {
	core := CredentialRecord(record.CredentialRecord)
	if err := r.db.Create(&core).Error; err != nil {
		return nil, fmt.Errorf("create core: %w", convertError(err))
	}

	login := r.buildLogin(core.ID, record)

	if err := r.db.Create(login).Error; err != nil {
		return nil, fmt.Errorf("create login: %w", convertError(err))
	}

	record.ID = login.ID

	return record, nil
}

func (r *RecordRepository) CreateCard(record *model.CardRecord) (*model.CardRecord, error) {
	core := CredentialRecord(record.CredentialRecord)
	if err := r.db.Create(&core).Error; err != nil {
		return nil, fmt.Errorf("create core: %w", convertError(err))
	}

	card := r.buildCard(core.ID, record)

	if err := r.db.Create(card).Error; err != nil {
		return nil, fmt.Errorf("create card: %w", convertError(err))
	}

	record.ID = card.ID

	return record, nil
}

func (r *RecordRepository) CreateIdentity(record *model.IdentityRecord) (*model.IdentityRecord, error) {
	core := CredentialRecord(record.CredentialRecord)
	if err := r.db.Create(&core).Error; err != nil {
		return nil, fmt.Errorf("create core: %w", convertError(err))
	}

	identity := r.buildIdentity(core.ID, record)

	if err := r.db.Create(identity).Error; err != nil {
		return nil, fmt.Errorf("create identity: %w", convertError(err))
	}

	record.ID = identity.ID

	return record, nil
}

func (r *RecordRepository) UpdateCredentialRecord(record *model.CredentialRecord) (*model.CredentialRecord, error) {
	credentialRecord := CredentialRecord(*record)
	result := r.db.Model(credentialRecord).Clauses(clause.Returning{}).Updates(credentialRecord)
	if result.Error != nil {
		return nil, fmt.Errorf("update: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	return record, nil
}

func (r *RecordRepository) UpdateLogin(record *model.LoginRecord) (*model.LoginRecord, error) {
	// FIXME: empty cred
	core := CredentialRecord(record.CredentialRecord)
	result := r.db.Model(core).Clauses(clause.Returning{}).Updates(core)
	if result.Error != nil {
		return nil, fmt.Errorf("update core: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	login := r.buildLogin(core.ID, record)

	result = r.db.Model(login).Clauses(clause.Returning{}).Updates(login)
	if result.Error != nil {
		return nil, fmt.Errorf("update login: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	return record, nil
}

func (r *RecordRepository) UpdateCard(record *model.CardRecord) (*model.CardRecord, error) {
	core := CredentialRecord(record.CredentialRecord)
	result := r.db.Model(core).Clauses(clause.Returning{}).Updates(core)
	if result.Error != nil {
		return nil, fmt.Errorf("update core: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	card := r.buildCard(core.ID, record)

	result = r.db.Model(card).Clauses(clause.Returning{}).Updates(card)
	if result.Error != nil {
		return nil, fmt.Errorf("update card: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	return record, nil
}

func (r *RecordRepository) UpdateIdentity(record *model.IdentityRecord) (*model.IdentityRecord, error) {
	core := CredentialRecord(record.CredentialRecord)
	result := r.db.Model(core).Clauses(clause.Returning{}).Updates(core)
	if result.Error != nil {
		return nil, fmt.Errorf("update core: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	identity := r.buildIdentity(core.ID, record)

	result = r.db.Model(identity).Clauses(clause.Returning{}).Updates(identity)
	if result.Error != nil {
		return nil, fmt.Errorf("update card: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	return record, nil
}

func (r *RecordRepository) Delete(id uuid.UUID) (*model.CredentialRecord, error) {
	var record CredentialRecord
	result := r.db.Clauses(clause.Returning{}).Where("id = ?", id).Delete(&record)
	if result.Error != nil {
		return nil, fmt.Errorf("delete: %w", convertError(result.Error))
	}

	if result.RowsAffected == 0 {
		return nil, pmerror.ErrNotFound
	}

	return (*model.CredentialRecord)(&record), nil
}

func (r *RecordRepository) buildLogin(id uuid.UUID, record *model.LoginRecord) *LoginRecord {
	return &LoginRecord{
		ID:       id,
		Username: record.Username,
		Password: record.Password,
		URL:      record.URL,
	}
}

func (r *RecordRepository) buildCard(id uuid.UUID, record *model.CardRecord) *CardRecord {
	return &CardRecord{
		ID:              id,
		Brand:           record.Brand,
		Number:          record.Number,
		ExpirationMonth: record.ExpirationMonth,
		ExpirationYear:  record.ExpirationYear,
		CVV:             record.CVV,
	}
}

func (r *RecordRepository) buildIdentity(id uuid.UUID, record *model.IdentityRecord) *IdentityRecord {
	return &IdentityRecord{
		ID:             id,
		FirstName:      record.FirstName,
		MiddleName:     record.MiddleName,
		LastName:       record.LastName,
		Address:        record.Address,
		Email:          record.Email,
		PhoneNumber:    record.PhoneNumber,
		PassportNumber: record.PassportNumber,
		Country:        record.Country,
	}
}

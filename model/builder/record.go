package builder

import (
	"github.com/okutsen/PasswordManager/model/api"
	"github.com/okutsen/PasswordManager/model/controller"
	"github.com/okutsen/PasswordManager/model/db"
)

func BuildControllerRecordFromDBRecord(record *db.CredentialRecord) controller.CredentialRecord {
	return controller.CredentialRecord{
		ID:          record.ID,
		Name:        record.Name,
		Login:       record.Login,
		Password:    record.Password,
		URL:         record.URL,
		Description: record.Description,
		CreatedBy:   record.CreatedBy,
		UpdatedBy:   record.UpdatedBy,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}
}

func BuildControllerRecordsFromDBRecords(records []db.CredentialRecord) []controller.CredentialRecord {
	recordsController := make([]controller.CredentialRecord, len(records))
	for i, v := range records {
		recordsController[i] = BuildControllerRecordFromDBRecord(&v)
	}

	return recordsController
}

func BuildAPIRecordFromControllerRecord(record *controller.CredentialRecord) api.CredentialRecord {
	return api.CredentialRecord{
		ID:          record.ID,
		Name:        record.Name,
		Login:       record.Login,
		Password:    record.Password,
		URL:         record.URL,
		Description: record.Description,
		CreatedBy:   record.CreatedBy,
		UpdatedBy:   record.UpdatedBy,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}
}

func BuildAPIRecordsFromControllerRecords(records []controller.CredentialRecord) []api.CredentialRecord {
	recordsController := make([]api.CredentialRecord, len(records))
	for i, v := range records {
		recordsController[i] = BuildAPIRecordFromControllerRecord(&v)
	}

	return recordsController
}

func BuildControllerRecordFromAPIRecord(record *api.CredentialRecord) controller.CredentialRecord {
	return controller.CredentialRecord{
		ID:          record.ID,
		Name:        record.Name,
		Login:       record.Login,
		Password:    record.Password,
		URL:         record.URL,
		Description: record.Description,
		CreatedBy:   record.CreatedBy,
		UpdatedBy:   record.UpdatedBy,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}
}

func BuildDBRecordFromControllerRecord(record *controller.CredentialRecord) db.CredentialRecord {
	return db.CredentialRecord{
		ID:          record.ID,
		Name:        record.Name,
		Login:       record.Login,
		Password:    record.Password,
		URL:         record.URL,
		Description: record.Description,
		CreatedBy:   record.CreatedBy,
		UpdatedBy:   record.UpdatedBy,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
	}
}

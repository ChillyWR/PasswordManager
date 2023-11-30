package builder

import (
	"github.com/okutsen/PasswordManager/model/api"
	"github.com/okutsen/PasswordManager/model/controller"
	"github.com/okutsen/PasswordManager/model/db"
)

func BuildControllerRecordFromDBRecord(record *db.Record) controller.Record {
	return controller.Record{
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

func BuildControllerRecordsFromDBRecords(records []db.Record) []controller.Record {
	recordsController := make([]controller.Record, len(records))
	for i, v := range records {
		recordsController[i] = BuildControllerRecordFromDBRecord(&v)
	}

	return recordsController
}

func BuildAPIRecordFromControllerRecord(record *controller.Record) api.Record {
	return api.Record{
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

func BuildAPIRecordsFromControllerRecords(records []controller.Record) []api.Record {
	recordsController := make([]api.Record, len(records))
	for i, v := range records {
		recordsController[i] = BuildAPIRecordFromControllerRecord(&v)
	}

	return recordsController
}

func BuildControllerRecordFromAPIRecord(record *api.Record) controller.Record {
	return controller.Record{
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

func BuildDBRecordFromControllerRecord(record *controller.Record) db.Record {
	return db.Record{
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

package controller

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	pmlogger "github.com/ChillyWR/PasswordManager/internal/logger"
	"github.com/ChillyWR/PasswordManager/internal/mock"
	"github.com/ChillyWR/PasswordManager/model"
	"github.com/ChillyWR/PasswordManager/pkg/pmcrypto"
	"github.com/ChillyWR/PasswordManager/pkg/pmerror"
)

type controllerMocks struct {
	RecordRepository *mock.MockRecordRepository
	UserRepository   *mock.MockUserRepository
}

type controllerTestCase struct {
	Name string
	Run  func(t *testing.T, c *Controller, mocks *controllerMocks)
}

func (tc controllerTestCase) runTests(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mocks := &controllerMocks{
		RecordRepository: mock.NewMockRecordRepository(ctrl),
		UserRepository:   mock.NewMockUserRepository(ctrl),
	}

	logger := pmlogger.New()

	c, err := New(mocks.UserRepository, mocks.RecordRepository, logger)
	require.NoError(t, err)

	tc.Run(t, c, mocks)
}

func TestController_GetRecord(t *testing.T) {
	userID, err := uuid.NewUUID()
	require.NoError(t, err)

	testCases := []controllerTestCase{
		{
			Name: "success_get_credential_record",
			Run: func(t *testing.T, c *Controller, mocks *controllerMocks) {
				id, err := uuid.NewUUID()
				require.NoError(t, err)

				notes := "Test Record Notes"
				encryptedNotes, err := pmcrypto.Encrypt(notes, Salt)
				require.NoError(t, err)

				record := &model.CredentialRecord{
					ID:        id,
					Name:      "Test Record Name",
					Notes:     &encryptedNotes,
					CreatedOn: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedOn: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					CreatedBy: userID,
					UpdatedBy: userID,
				}

				mocks.RecordRepository.EXPECT().
					GetCredentialRecord(id).
					Return(record, nil)

				mocks.RecordRepository.EXPECT().
					GetCredentialRecord(id).
					Return(record, nil)

				mocks.RecordRepository.EXPECT().
					GetLogin(id).
					Return(nil, pmerror.ErrNotFound)

				mocks.RecordRepository.EXPECT().
					GetCard(id).
					Return(nil, pmerror.ErrNotFound)

				mocks.RecordRepository.EXPECT().
					GetIdentity(id).
					Return(nil, pmerror.ErrNotFound)

				actual, err := c.GetRecord(id, userID)
				require.NoError(t, err)

				expected := record
				expected.Notes = &notes

				require.Equal(t, record, actual)
			},
		},
		{
			Name: "error_not_found",
			Run: func(t *testing.T, c *Controller, mocks *controllerMocks) {
				id, err := uuid.NewUUID()
				require.NoError(t, err)

				mocks.RecordRepository.EXPECT().
					GetCredentialRecord(id).
					Return(nil, pmerror.ErrInternal)

				_, err = c.GetRecord(id, userID)
				require.True(t, errors.Is(err, pmerror.ErrInternal))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			tc.runTests(t)
		})
	}
}

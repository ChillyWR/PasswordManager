package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/require"

	"github.com/ChillyWR/PasswordManager/internal/controller"
	pmlogger "github.com/ChillyWR/PasswordManager/internal/logger"
	"github.com/ChillyWR/PasswordManager/internal/repo"
	"github.com/ChillyWR/PasswordManager/model"
	"github.com/ChillyWR/PasswordManager/pkg/pmcrypto"
	"github.com/ChillyWR/PasswordManager/pkg/pmpointer"
)

type TableTest struct {
	testName           string
	order              int
	handle             httprouter.Handle
	httpMethod         string
	path               string
	ps                 httprouter.Params
	expectedHTTPStatus int
	expectedBody       string
}

type TableTests struct {
	tt []*TableTest
}

func setup() *APIContext {
	logger := pmlogger.New()

	userRepo, err := repo.NewUserRepository(testDB)
	if err != nil {
		logger.Fatalf("Failed to initialize DB: %s", err.Error())
	}

	recordRepo, err := repo.NewRecordRepository(testDB)
	if err != nil {
		logger.Fatalf("Failed to initialize DB: %s", err.Error())
	}

	ctrl, err := controller.New(userRepo, recordRepo, logger)
	if err != nil {
		logger.Fatalf("Failed to init ctrl: %s", err.Error())
	}

	return &APIContext{ctrl, logger}
}

func TestGet(t *testing.T) {
	apictx := setup()

	testUser1, err := apictx.ctrl.CreateUser(&model.UserForm{
		Name:     pmpointer.String("Test User Name 1"),
		Password: pmpointer.String("Test User Password 1"),
	})
	require.NoError(t, err)

	rawTestRecord1, err := apictx.ctrl.CreateRecord(model.SecureNoteRecordType, json.RawMessage(`{
		"name": "Test Record Name 1",
		"notes": "Test Record Notes 1"
	}`), testUser1.ID)
	require.NoError(t, err)

	rawTestRecord2, err := apictx.ctrl.CreateRecord(model.SecureNoteRecordType, json.RawMessage(`{
		"name": "Test Record Name 2",
		"notes": "Test Record Notes 2"
	}`), testUser1.ID)
	require.NoError(t, err)

	testRecord1, ok := rawTestRecord1.(*model.CredentialRecord)
	require.True(t, ok)

	testRecord2, ok := rawTestRecord2.(*model.CredentialRecord)
	require.True(t, ok)

	defer cleanup(t, apictx.ctrl, []*model.CredentialRecord{testRecord1, testRecord2}, testUser1)

	val := *testRecord1
	decryptedTestRecord1 := val // copy
	v, err := pmcrypto.Decrypt(*decryptedTestRecord1.Notes, controller.Salt)
	require.NoError(t, err)
	decryptedTestRecord1.Notes = &v

	tts := TableTests{
		tt: []*TableTest{
			{
				testName: "success_get_all_records",
				order:    1,
				handle: ContextSetter(apictx.logger, testUserAuthentication(apictx.logger, testUser1.ID,
					Dispatch(NewListRecordsHandler(apictx)))),
				httpMethod:         http.MethodGet,
				path:               "/records",
				expectedHTTPStatus: http.StatusOK,
				expectedBody: fmt.Sprintf(
					`{"secure_notes":[%s, %s],"logins":[],"cards":[],"identities":[]}`,
					toJSONString(t, testRecord1), toJSONString(t, testRecord2),
				),
			},
			{
				testName: "success_get_decrypted_record_by_id",
				order:    2,
				handle: ContextSetter(apictx.logger, testUserAuthentication(apictx.logger, testUser1.ID,
					Dispatch(NewGetRecordHandler(apictx)))),
				httpMethod: http.MethodGet,
				path:       fmt.Sprintf("/records/%s", testRecord1.ID.String()),
				ps: httprouter.Params{
					httprouter.Param{Key: IDPPN, Value: testRecord1.ID.String()},
				},
				expectedHTTPStatus: http.StatusOK,
				expectedBody:       toJSONString(t, decryptedTestRecord1),
			},
			{
				testName: "error_invalid_record_id",
				order:    3,
				handle: ContextSetter(apictx.logger, testUserAuthentication(apictx.logger, testUser1.ID,
					Dispatch(NewGetRecordHandler(apictx)))),
				httpMethod: http.MethodGet,
				path:       "/records/a",
				ps: httprouter.Params{
					httprouter.Param{Key: IDPPN, Value: "a"},
				},
				expectedHTTPStatus: http.StatusBadRequest,
				expectedBody:       fmt.Sprintf(`{"message":"%s"}`, InvalidRecordIDMessage),
			},
		},
	}

	TableTestRunner(t, tts)
}

func TableTestRunner(t *testing.T, tts TableTests) {
	t.Helper()
	for _, test := range tts.tt {
		t.Run(test.testName, func(t *testing.T) {
			request := httptest.NewRequest(test.httpMethod, test.path, nil)
			response := httptest.NewRecorder()
			test.handle(response, request, test.ps)
			require.Equal(t, test.expectedHTTPStatus, response.Code)
			require.JSONEq(t, test.expectedBody, response.Body.String())
		})
	}
}

func OrderedTableTestRunner(t *testing.T, tts TableTests) {
	t.Helper()

	sort.Slice(tts.tt, func(i, j int) bool { return tts.tt[i].order < tts.tt[j].order })

	for _, test := range tts.tt {
		fmt.Printf("RUN %s\n", path.Join(t.Name(), test.testName))
		start := time.Now()
		request := httptest.NewRequest(test.httpMethod, test.path, nil)
		response := httptest.NewRecorder()
		test.handle(response, request, test.ps)
		require.Equal(t, test.expectedHTTPStatus, response.Code)
		require.JSONEq(t, test.expectedBody, response.Body.String())
		fmt.Printf("PASS %s (%s)\n", path.Join(t.Name(), test.testName), time.Since(start).String())
	}
}

func cleanup(t *testing.T, ctrl Controller, records []*model.CredentialRecord, user *model.User) {
	t.Helper()
	var err error
	for _, record := range records {
		_, err = ctrl.DeleteRecord(record.ID, user.ID)
		require.NoError(t, err)
	}
	_, err = ctrl.DeleteUser(user.ID)
	require.NoError(t, err)
}

func testUserAuthentication(logger pmlogger.Logger, userID uuid.UUID, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		rctx := unpackRequestContext(r.Context(), logger)
		rctx.userID = userID
		ctx := context.WithValue(r.Context(), RequestContextName, rctx)

		next(w, r.WithContext(ctx), ps)
	}
}

func toJSONString(t *testing.T, v any) string {
	t.Helper()
	raw, err := json.Marshal(v)
	require.NoError(t, err)
	return string(raw)
}

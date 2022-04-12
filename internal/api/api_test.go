package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/sirupsen/logrus"
)

const (
	testSeparator string = "\n---------------------\n"
)

type TableTest struct {
	// Create request for specified
	testName           string
	handle            HandlerFunc
	httpMethod         string
	httpPath           string
	ps httprouter.Params
	expectedHTTPStatus int
	expectedBody       string
}

type TableTests struct {
	tt         []*TableTest // by pointer or not?
	httpServer *API
}

type HandlerFunc func(*API, http.ResponseWriter, *http.Request, httprouter.Params)

func TestGetRecords(t *testing.T) {
	var config *Config = NewConfig()
	var logger log.Logger = logrus.New()
	tests := TableTests{
		tt: []*TableTest{
			{
				testName:           "Get all records",
				handle:            (*API).getAllRecords,
				httpMethod:         http.MethodGet,
				httpPath:           "/records/",
				expectedHTTPStatus: http.StatusOK,
				expectedBody:       "Records:\n0,1,2,3,4,5",
			},
			{
				testName:           "Get record by id 1",
				handle:            (*API).getRecord,
				httpMethod:         http.MethodGet,
				httpPath:           "/records/0",
				ps: httprouter.Params{
					httprouter.Param{Key: IDParamName, Value: "0"},
				},
				expectedHTTPStatus: http.StatusOK,
				expectedBody:       "Records:\n0",
			},
			{
				testName:           "Get record by id 5",
				handle:            (*API).getRecord,
				httpMethod:         http.MethodGet,
				httpPath:           "/records/5",
				ps: httprouter.Params{
					httprouter.Param{Key: IDParamName, Value: "5"},
				},
				expectedHTTPStatus: http.StatusOK,
				expectedBody:       "Records:\n5",
			},
			{
				testName:           "Returns 404 on missing record",
				handle:            (*API).getRecord,
				httpMethod:         http.MethodGet,
				httpPath:           "/records/6",
				ps: httprouter.Params{
					httprouter.Param{Key: IDParamName, Value: "6"},
				},
				expectedHTTPStatus: http.StatusBadRequest,
				expectedBody:       "Record not found",
			}},

		// TODO: use mocks
		httpServer: New(config, logger),
	}
	TableTestRunner(t, tests)
}

func TestPostRecords(t *testing.T) {
	var config *Config = NewConfig()
	var logger log.Logger = logrus.New() 
	tests := TableTests{
		tt: []*TableTest{
			{
				testName:           "Post record",
				expectedHTTPStatus: http.StatusAccepted,
				expectedBody:       "New Record created",
				handle:            (*API).createRecords,
				httpMethod:         http.MethodPost,
				httpPath:           "/records/",
			}},

		// TODO: use mocks
		httpServer: New(config, logger),
	}
	TableTestRunner(t, tests)
}

func TableTestRunner(t *testing.T, tt TableTests) {
	t.Helper()
	for _, test := range tt.tt {
		t.Run(test.testName, func(t *testing.T) {
			request := httptest.NewRequest(test.httpMethod, test.httpPath, nil)
			response := httptest.NewRecorder()
			test.handle(tt.httpServer, response, request, test.ps)

			assert(t, response.Code, test.expectedHTTPStatus, "Wrong status")
			assert(t, response.Body.String(), test.expectedBody, "Wrong body")
		})
	}
}

func assert[T comparable](t *testing.T, got, want T, errorMessage string) {
	t.Helper()
	if got != want {
		t.Errorf("%s\nGot:%s%v%sWant:%s%v%s",
			errorMessage,
			testSeparator, got, testSeparator,
			testSeparator, want, testSeparator)
	}
}
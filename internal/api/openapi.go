package api

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model/api"
)

func generateSchemas(logger log.Logger) openapi3.Schemas {
	schemas := make(openapi3.Schemas)
	gen := openapi3gen.NewGenerator()
	RecordRef, err := gen.NewSchemaRefForValue(&api.Record{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from Record")
	}
	UserRef, err := gen.NewSchemaRefForValue(&api.User{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from User")
	}
	ErrorRef, err := gen.NewSchemaRefForValue(&api.Error{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from Error")
	}
	resultSchema := openapi3.Schemas{
		"Record": RecordRef,
		"User":   UserRef,
		"Error":  ErrorRef,
	}
	return resultSchema
}

// NewOpenAPIv3 instantiates the OpenAPI specification
func NewOpenAPIv3(cfg *Config, logger log.Logger) *openapi3.T {
	spec := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:       "Password Manager",
			Description: "",
			Version:     "0.0.0",
			Contact: &openapi3.Contact{
				URL: "https://github.com/okutsen/PasswordManager",
			},
		},
		Servers: openapi3.Servers{
			&openapi3.Server{
				Description: "Local development",
				URL:         "http://" + cfg.LocalAddress(),
			},
		},
	}
	spec.Components.Schemas = generateSchemas(logger)
	spec.Components.SecuritySchemes = openapi3.SecuritySchemes{
		"AuthorizationToken": &openapi3.SecuritySchemeRef{
			Value: openapi3.NewJWTSecurityScheme(),
		},
	}
	spec.Components.Parameters = openapi3.ParametersMap{
		"IDPPN": &openapi3.ParameterRef{
			Value: openapi3.NewPathParameter(IDPPN).
				WithRequired(true).
				WithSchema(openapi3.NewUUIDSchema()),
		},
		"CorrelationIDHPN": &openapi3.ParameterRef{
			Value: openapi3.NewHeaderParameter(CorrelationIDHPN).
				WithDescription("Correlation id").
				WithSchema(openapi3.NewUUIDSchema()),
		},
		"AuthorizationTokenHPN": &openapi3.ParameterRef{
			Value: openapi3.NewHeaderParameter(AuthorizationTokenHPN).
				WithDescription("Correlation id").
				WithSchema(openapi3.NewUUIDSchema()),
		},
	}
	spec.Components.RequestBodies = openapi3.RequestBodies{
		"CreateRecordRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for creating a record.").
				WithRequired(true).
				WithJSONSchemaRef(&openapi3.SchemaRef{
					Ref: "#/components/schemas/Record",
				}),
		},
		"UpdateRecordRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for updating a record.").
				WithRequired(true).
				WithJSONSchemaRef(&openapi3.SchemaRef{
					Ref: "#/components/schemas/Record",
				}),
		},
	}
	spec.Components.Responses = openapi3.ResponseBodies{
		"ListRecordsResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returns back all records.").
				WithJSONSchema(&openapi3.Schema{
					Type:  openapi3.TypeArray,
					Items: openapi3.NewSchemaRef("#/components/schemas/Record", nil),
				}),
		},
		"RecordResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response returns back successfully found or created record.").
				WithJSONSchemaRef(&openapi3.SchemaRef{
					Ref: "#/components/schemas/Record",
				}),
		},
		"ErrorResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Response when errors happen.").
				WithJSONSchemaRef(&openapi3.SchemaRef{
					Ref: "#/components/schemas/Error",
				}),
		},
	}

	spec.Paths = openapi3.NewPaths(
		openapi3.WithPath("/records", &openapi3.PathItem{
			Get: &openapi3.Operation{
				OperationID: "ListRecords",
				Responses: openapi3.NewResponses(
					openapi3.WithStatus(200, &openapi3.ResponseRef{
						Ref: "#/components/responses/ListRecordsResponse",
					}),
					openapi3.WithStatus(500, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
				),
			},
			Post: &openapi3.Operation{
				OperationID: "CreateRecord",
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/CreateRecordRequest",
				},
				Responses: openapi3.NewResponses(
					openapi3.WithStatus(201, &openapi3.ResponseRef{
						Ref: "#/components/responses/RecordResponse",
					}),
					openapi3.WithStatus(400, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
					openapi3.WithStatus(500, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
				),
			},
		}),
		openapi3.WithPath(fmt.Sprintf("/records/{%s}", IDPPN), &openapi3.PathItem{
			Get: &openapi3.Operation{
				OperationID: "GetRecord",
				Parameters: []*openapi3.ParameterRef{{
					Ref: "#/components/parameters/IDPPN",
				}},
				Responses: openapi3.NewResponses(
					openapi3.WithStatus(200, &openapi3.ResponseRef{
						Ref: "#/components/responses/RecordResponse",
					}),
					openapi3.WithStatus(400, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
					openapi3.WithStatus(500, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
				),
			},
			Put: &openapi3.Operation{
				OperationID: "UpdateRecord",
				Parameters: []*openapi3.ParameterRef{{
					Ref: "#/components/parameters/IDPPN",
				}},
				RequestBody: &openapi3.RequestBodyRef{
					Ref: "#/components/requestBodies/UpdateRecordRequest",
				},
				Responses: openapi3.NewResponses(
					openapi3.WithStatus(202, &openapi3.ResponseRef{
						Ref: "#/components/responses/RecordResponse",
					}),
					openapi3.WithStatus(400, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
					openapi3.WithStatus(500, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
				),
			},
			Delete: &openapi3.Operation{
				OperationID: "DeleteRecord",
				Parameters: []*openapi3.ParameterRef{{
					Ref: "#/components/parameters/IDPPN",
				}},
				Responses: openapi3.NewResponses(
					openapi3.WithStatus(200, &openapi3.ResponseRef{
						Value: openapi3.NewResponse().WithDescription("Record deleted"),
					}),
					openapi3.WithStatus(400, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
					openapi3.WithStatus(500, &openapi3.ResponseRef{
						Ref: "#/components/responses/ErrorResponse",
					}),
				),
			},
		}),
	)

	return &spec
}

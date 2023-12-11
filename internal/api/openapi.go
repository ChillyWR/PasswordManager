package api

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"

	"github.com/okutsen/PasswordManager/internal/log"
	"github.com/okutsen/PasswordManager/model"
)

func generateSchemas(logger log.Logger) openapi3.Schemas {
	schemas := make(openapi3.Schemas)
	gen := openapi3gen.NewGenerator()

	UserRef, err := gen.NewSchemaRefForValue(&model.User{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from User")
	}

	CredentialRecordRef, err := gen.NewSchemaRefForValue(&model.CredentialRecord{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from CredentialRecord")
	}

	CredentialRecordFormRef, err := gen.NewSchemaRefForValue(&model.CredentialRecordForm{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from CredentialRecordForm")
	}

	LoginRecordRef, err := gen.NewSchemaRefForValue(&model.LoginRecord{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from LoginRecord")
	}

	LoginRecordFormRef, err := gen.NewSchemaRefForValue(&model.LoginRecordForm{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from LoginRecordForm")
	}

	CardRecordRef, err := gen.NewSchemaRefForValue(&model.CardRecord{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from CardRecord")
	}

	CardRecordFormRef, err := gen.NewSchemaRefForValue(&model.CardRecordForm{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from CardRecordForm")
	}

	IdentityRecordRef, err := gen.NewSchemaRefForValue(&model.IdentityRecord{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from IdentityRecord")
	}

	IdentityRecordFormRef, err := gen.NewSchemaRefForValue(&model.IdentityRecordForm{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from IdentityRecordForm")
	}

	ErrorRef, err := gen.NewSchemaRefForValue(&Error{}, schemas)
	if err != nil {
		logger.Fatal("Failed to generate schema from Error")
	}

	resultSchema := openapi3.Schemas{
		"User":                 UserRef,
		"CredentialRecord":     CredentialRecordRef,
		"CredentialRecordForm": CredentialRecordFormRef,
		"LoginRecord":          LoginRecordRef,
		"LoginRecordForm":      LoginRecordFormRef,
		"CardRecord":           CardRecordRef,
		"CardRecordForm":       CardRecordFormRef,
		"IdentityRecord":       IdentityRecordRef,
		"IdentityRecordForm":   IdentityRecordFormRef,
		"Error":                ErrorRef,
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
		Components: &openapi3.Components{
			Schemas: generateSchemas(logger),
		},
	}
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
				WithDescription("JWT Token").
				WithSchema(openapi3.NewUUIDSchema()),
		},
	}
	spec.Components.RequestBodies = openapi3.RequestBodies{
		"CreateRecordRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Record creation request body. Form should correspond to the record type.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
					"type": openapi3.NewStringSchema().WithEnum(
						"secure_note",
						"login",
						"card",
						"identity",
					),
					"form": openapi3.NewAnyOfSchema(
						openapi3.NewSchemaRef("#/components/schemas/CredentialRecordForm", nil).Value,
						openapi3.NewSchemaRef("#/components/schemas/LoginRecordForm", nil).Value,
						openapi3.NewSchemaRef("#/components/schemas/CardRecordForm", nil).Value,
						openapi3.NewSchemaRef("#/components/schemas/IdentityRecordForm", nil).Value,
					),
				})),
		},
		"UpdateRecordRequest": &openapi3.RequestBodyRef{
			Value: openapi3.NewRequestBody().
				WithDescription("Request used for updating a record. Form should correspond to the record type. Fields of the form are optional, filled ones will be used to update record.").
				WithRequired(true).
				WithJSONSchema(openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
					"type": openapi3.NewStringSchema().WithEnum(
						"secure_note",
						"login",
						"card",
						"identity",
					),
					"form": openapi3.NewAnyOfSchema(
						openapi3.NewSchemaRef("#/components/schemas/CredentialRecordForm", nil).Value,
						openapi3.NewSchemaRef("#/components/schemas/LoginRecordForm", nil).Value,
						openapi3.NewSchemaRef("#/components/schemas/CardRecordForm", nil).Value,
						openapi3.NewSchemaRef("#/components/schemas/IdentityRecordForm", nil).Value,
					),
				})),
		},
	}
	spec.Components.Responses = openapi3.ResponseBodies{
		"ListRecordsResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("All user's records response.").
				WithJSONSchema(openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
					"secure_notes": openapi3.NewArraySchema().WithItems(
						openapi3.NewSchemaRef("#/components/schemas/CredentialRecord", nil).Value),
					"logins": openapi3.NewArraySchema().WithItems(
						openapi3.NewSchemaRef("#/components/schemas/LoginRecord", nil).Value),
					"card": openapi3.NewArraySchema().WithItems(
						openapi3.NewSchemaRef("#/components/schemas/CardRecord", nil).Value),
					"identity": openapi3.NewArraySchema().WithItems(
						openapi3.NewSchemaRef("#/components/schemas/IdentityRecord", nil).Value),
				})),
		},
		"RecordResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Record by id response.").
				WithJSONSchema(openapi3.NewAnyOfSchema(
					openapi3.NewSchemaRef("#/components/schemas/CredentialRecord", nil).Value,
					openapi3.NewSchemaRef("#/components/schemas/LoginRecord", nil).Value,
					openapi3.NewSchemaRef("#/components/schemas/CardRecord", nil).Value,
					openapi3.NewSchemaRef("#/components/schemas/IdentityRecord", nil).Value,
				)),
		},
		"RecordDeletionResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Record by id response.").
				WithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/CredentialRecord", nil)),
		},
		"ErrorResponse": &openapi3.ResponseRef{
			Value: openapi3.NewResponse().
				WithDescription("Error response.").
				WithJSONSchemaRef(openapi3.NewSchemaRef("#/components/schemas/Error", nil)),
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
			Patch: &openapi3.Operation{
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
						Ref: "#/components/responses/RecordDeletionResponse",
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

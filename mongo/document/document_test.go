package document

import (
	"context"
	"testing"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2"
)

type unvalidatedDocument struct {
	Base `bson:",inline"`
}

const testDocuments = "test_documents"

type validatedDocument struct {
	Base  `bson:",inline"`
	Valid bool `bson:"valid" json:"valid"`
}

func (d *validatedDocument) Validate(ctx context.Context) *errors.ValidationErrors {
	verr := errors.NewValidationErrorsBuilder()
	if !d.Valid {
		verr.Set("valid", "must be true")
	}
	return verr.Build()
}

func buildValidatedDocument(valid bool) *validatedDocument {
	return &validatedDocument{
		Valid: valid,
	}
}

type validatedWithInternalErrorDocument struct {
	validatedDocument `bson:",inline"`
	InternalError     string `bson:"internal_error" json:"internal_error"`
}

func (d *validatedWithInternalErrorDocument) ValidateWithInternalError(ctx context.Context) (*errors.ValidationErrors, error) {
	if d.InternalError != "" {
		return nil, errors.New(ctx, d.InternalError)
	}
	return d.validatedDocument.Validate(ctx), nil
}

func buildValidatedWithInternalErrorDocument(valid bool, internalError string) *validatedWithInternalErrorDocument {
	return &validatedWithInternalErrorDocument{
		validatedDocument: *buildValidatedDocument(valid),
		InternalError:     internalError,
	}
}

func TestDocument_Create(t *testing.T) {
	t.Cleanup(func() {
		coll := mongo.Session(logger.Default()).Clone().DB("").C(testDocuments)
		err := coll.DropCollection()
		// Handle case when collection does not exist
		var queryErr *mgo.QueryError
		if !errors.As(err, &queryErr) || queryErr.Message != "ns not found" {
			require.NoError(t, err)
		}
	})
	t.Run("without validation", func(t *testing.T) {
		t.Run("it should create the document", func(t *testing.T) {
			d := &unvalidatedDocument{}
			err := Create(context.Background(), testDocuments, d)
			require.NoError(t, err)
			assert.NotEmpty(t, d.ID)
			assert.NotEmpty(t, d.CreatedAt)
		})
	})

	t.Run("with simple validation", func(t *testing.T) {
		t.Run("with a valid document, it should create it", func(t *testing.T) {
			d := buildValidatedDocument(true)
			err := Create(context.Background(), testDocuments, d)
			require.NoError(t, err)
			assert.NotEmpty(t, d.ID)
			assert.NotEmpty(t, d.CreatedAt)
		})

		t.Run("with an invalid document, it should return a validation error", func(t *testing.T) {
			d := buildValidatedDocument(false)
			err := Create(context.Background(), testDocuments, d)
			require.Error(t, err)
			require.IsType(t, &errors.ValidationErrors{}, err)
			assert.Empty(t, d.ID)
			assert.Empty(t, d.CreatedAt)
		})
	})

	t.Run("with internal error validation", func(t *testing.T) {
		t.Run("with a valid document, it should create it", func(t *testing.T) {
			d := buildValidatedWithInternalErrorDocument(true, "")
			err := Create(context.Background(), testDocuments, d)
			require.NoError(t, err)
			assert.NotEmpty(t, d.ID)
			assert.NotEmpty(t, d.CreatedAt)
		})
		t.Run("with an invalid document, it should return a validation error", func(t *testing.T) {
			d := buildValidatedWithInternalErrorDocument(false, "")
			err := Create(context.Background(), testDocuments, d)
			require.Error(t, err)
			require.IsType(t, &errors.ValidationErrors{}, err)
			assert.Empty(t, d.ID)
			assert.Empty(t, d.CreatedAt)
		})

		t.Run("with a validation returning an internal error, it should forward internal error", func(t *testing.T) {
			d := buildValidatedWithInternalErrorDocument(true, "internal error when validating")
			err := Create(context.Background(), testDocuments, d)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "internal error when validating")
			assert.Empty(t, d.ID)
			assert.Empty(t, d.CreatedAt)
		})

		t.Run("with a validation returning an internal error and a validation error, it should return the internal error", func(t *testing.T) {
			d := buildValidatedWithInternalErrorDocument(false, "internal error when validating")
			err := Create(context.Background(), testDocuments, d)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "internal error when validating")
			assert.Empty(t, d.ID)
			assert.Empty(t, d.CreatedAt)
		})
	})
}

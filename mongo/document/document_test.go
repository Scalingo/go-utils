package document

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/mongo"
)

type unvalidatedDocument struct {
	Base `bson:",inline"`
}

const testDocuments = "test_documents"

type validatedDocument struct {
	Base  `bson:",inline"`
	Valid bool `bson:"valid" json:"valid"`
}

func (d *validatedDocument) Validate(_ context.Context) error {
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
}

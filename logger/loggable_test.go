package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type StructWithTags struct {
	Field1 string `log:"field1"`
	Field2 string `log:"field2"`
	Field3 string
}

type StructWithoutTagsButWithStringer struct {
	Field1 string
	Field2 string
}

func (s StructWithoutTagsButWithStringer) String() string {
	return "My Stringer"
}

type StructWithoutTags struct {
	Field1 string
	Field2 string
}

func TestFieldsFor(t *testing.T) {
	t.Run("when the struct has some tags", func(t *testing.T) {
		// Given a struct with tags
		s := StructWithTags{
			Field1: "value1",
			Field2: "value2",
			Field3: "value3",
		}

		// When we try to add it to a logger
		fields := FieldsFor(s, "prefix")

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_field1": "value1",
			"prefix_field2": "value2",
		}, fields)
	})

	t.Run("when the struct has no tags but has a stringer", func(t *testing.T) {
		// Given a struct without tags but with a stringer
		s := StructWithoutTagsButWithStringer{
			Field1: "value1",
			Field2: "value2",
		}

		// When we try to add it to a logger
		fields := FieldsFor(s, "prefix")

		// Then it should be added as a single field
		assert.Equal(t, logrus.Fields{
			"prefix": "My Stringer",
		}, fields)
	})

	t.Run("when the struct has no tags and no stringer", func(t *testing.T) {
		// Given a struct without tags
		s := StructWithoutTags{
			Field1: "value1",
			Field2: "value2",
		}

		// When we try to add it to a logger
		fields := FieldsFor(s, "prefix")

		// Then it should be added as a single field
		assert.Equal(t, logrus.Fields{
			"prefix": "failed to use FieldsFor on struct: Invalid type",
		}, fields)
	})
}

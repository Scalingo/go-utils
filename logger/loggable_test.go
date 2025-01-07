package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type structWithTags struct {
	Field1 string `log:"field1"`
	Field2 string `log:"field2"`
	Field3 string
}

type structWithTagsAndLoggable struct {
	Field1 string `log:"field1"`
	Field2 string `log:"field2"`
	Field3 string
}

func (s structWithTagsAndLoggable) ToLogrusFields() logrus.Fields {
	return logrus.Fields{
		"another": "test",
	}
}

type structWithoutTagsButWithStringer struct {
	Field1 string
	Field2 string
}

func (s structWithoutTagsButWithStringer) String() string {
	return "My Stringer"
}

type structWithoutTags struct {
	Field1 string
	Field2 string
}

func TestFieldsFor(t *testing.T) {
	t.Run("when the struct has some tags", func(t *testing.T) {
		// Given a struct with tags
		s := structWithTags{
			Field1: "value1",
			Field2: "value2",
			Field3: "value3",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_field1": "value1",
			"prefix_field2": "value2",
		}, fields)
	})

	t.Run("when we get a pointer to a struct with some tags", func(t *testing.T) {
		// Given a pointer to a struct with tags
		s := &structWithTags{
			Field1: "value1",
			Field2: "value2",
			Field3: "value3",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_field1": "value1",
			"prefix_field2": "value2",
		}, fields)
	})

	t.Run("when the struct has some tags and implements Loggable", func(t *testing.T) {
		// Given a struct with tags and that implements Loggable
		s := structWithTagsAndLoggable{
			Field1: "value1",
			Field2: "value2",
			Field3: "value3",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_another": "test",
		}, fields)
	})

	t.Run("when the struct has no tags but has a stringer", func(t *testing.T) {
		// Given a struct without tags but with a stringer
		s := structWithoutTagsButWithStringer{
			Field1: "value1",
			Field2: "value2",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as a single field
		assert.Equal(t, logrus.Fields{
			"prefix": "My Stringer",
		}, fields)
	})

	t.Run("when the struct has no tags and no stringer", func(t *testing.T) {
		// Given a struct without tags
		s := structWithoutTags{
			Field1: "value1",
			Field2: "value2",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as a single field
		assert.Equal(t, logrus.Fields{
			"prefix": "failed to use FieldsFor on struct: invalid type",
		}, fields)
	})

	t.Run("It should not panic on non struct types", func(t *testing.T) {
		assert.Equal(t, logrus.Fields{
			"prefix": "failed to use FieldsFor on struct: invalid type",
		}, FieldsFor("prefix", "test"))

		assert.Equal(t, logrus.Fields{
			"prefix": "failed to use FieldsFor on struct: invalid type",
		}, FieldsFor("prefix", 10.45))

		assert.Equal(t, logrus.Fields{
			"prefix": "failed to use FieldsFor on struct: invalid type",
		}, FieldsFor("prefix", nil))

		assert.Equal(t, logrus.Fields{
			"prefix": "failed to use FieldsFor on struct: invalid type",
		}, FieldsFor("prefix", true))
	})
}

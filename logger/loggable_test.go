package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const zeroValue = "this is zero"

type zeroable string

func (z zeroable) IsZero() bool {
	return z == zeroValue
}

type structWithTags struct {
	Field1                 string   `log:"field1"`
	Field2                 string   `log:"field2"`
	Field3Omitempty        string   `log:"field3,omitempty"`
	Field4OmitzeroString   string   `log:"field4,omitzero"`
	Field5OmitzeroZeroable zeroable `log:"field5,omitzero"`
	Field6                 string
}

type structWithTagsAndLoggable struct {
	Field1 string `log:"field1"`
	Field2 string `log:"field2"`
	Field3 string
}

func (s structWithTagsAndLoggable) LogFields() logrus.Fields {
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

type structInStructWithTags struct {
	Field1 structWithoutTags `log:"field1"`
}

type structInStructWithoutTagsButWithStringer struct {
	Field1 structWithoutTags
}

func (s structInStructWithoutTagsButWithStringer) String() string {
	return "Embedded Field2 is " + s.Field1.Field2
}

type structInStructWithoutTags struct {
	Field1 structWithoutTags
}

func TestFieldsFor(t *testing.T) {
	t.Run("when the struct has some tags", func(t *testing.T) {
		// Given a struct with tags
		s := structWithTags{
			Field1:                 "value1",
			Field2:                 "",
			Field3Omitempty:        "value3",
			Field4OmitzeroString:   "value4",
			Field5OmitzeroZeroable: "value5",
			Field6:                 "value6",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_field1": "value1",
			"prefix_field2": "",
			"prefix_field3": "value3",
			"prefix_field4": "value4",
			"prefix_field5": zeroable("value5"),
		}, fields)
	})

	t.Run("when the struct has some tags and an omitempty field with an empty value", func(t *testing.T) {
		// Given a struct with tags
		s := structWithTags{
			Field1:                 "value1",
			Field2:                 "",
			Field3Omitempty:        "",
			Field4OmitzeroString:   "value4",
			Field5OmitzeroZeroable: "value5",
			Field6:                 "value6",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_field1": "value1",
			"prefix_field2": "",
			"prefix_field4": "value4",
			"prefix_field5": zeroable("value5"),
		}, fields)
	})

	t.Run("when the struct has some tags and an omitzero field with an empty value", func(t *testing.T) {
		// Given a struct with tags
		s := structWithTags{
			Field1:                 "value1",
			Field2:                 "",
			Field3Omitempty:        "",
			Field4OmitzeroString:   "",
			Field5OmitzeroZeroable: "",
			Field6:                 "value6",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_field1": "value1",
			"prefix_field2": "",
			"prefix_field5": zeroable(""),
		}, fields)
	})

	t.Run("when the struct has some tags and an omitzero field with a zero value", func(t *testing.T) {
		// Given a struct with tags
		s := structWithTags{
			Field1:                 "value1",
			Field2:                 "",
			Field3Omitempty:        "",
			Field4OmitzeroString:   "",
			Field5OmitzeroZeroable: zeroValue,
			Field6:                 "value6",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_field1": "value1",
			"prefix_field2": "",
		}, fields)
	})

	t.Run("when we get a pointer to a struct with some tags", func(t *testing.T) {
		// Given a pointer to a struct with tags
		s := &structWithTags{
			Field1:                 "value1",
			Field2:                 "",
			Field3Omitempty:        "value3",
			Field4OmitzeroString:   "value4",
			Field5OmitzeroZeroable: "value5",
			Field6:                 "value6",
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then it should be added as separate fields
		assert.Equal(t, logrus.Fields{
			"prefix_field1": "value1",
			"prefix_field2": "",
			"prefix_field3": "value3",
			"prefix_field4": "value4",
			"prefix_field5": zeroable("value5"),
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

	// We didn't implement yet the support for struct in struct.
	// Hence the output is uggly.
	t.Run("when a struct in struct has a tag", func(t *testing.T) {
		// Given a struct with tags
		s := structInStructWithTags{
			Field1: structWithoutTags{
				Field2: "value2",
			},
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then
		assert.Equal(t, logrus.Fields{
			"prefix_field1": structWithoutTags{
				Field2: "value2",
			},
		}, fields)
	})

	t.Run("when a struct in struct has no tag but implements Stringer", func(t *testing.T) {
		// Given a struct with tags
		s := structInStructWithoutTagsButWithStringer{
			Field1: structWithoutTags{
				Field2: "value2",
			},
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then
		assert.Equal(t, logrus.Fields{
			"prefix": "Embedded Field2 is value2",
		}, fields)
	})

	t.Run("when a struct in struct has no tag", func(t *testing.T) {
		// Given a struct with tags
		s := structInStructWithoutTags{
			Field1: structWithoutTags{
				Field2: "value2",
			},
		}

		// When we try to add it to a logger
		fields := FieldsFor("prefix", s)

		// Then
		assert.Equal(t, logrus.Fields{
			"prefix": "failed to use FieldsFor on struct: invalid type",
		}, fields)
	})
}

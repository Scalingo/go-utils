package logger

import (
	"regexp"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestRedactingFormatter_Format(t *testing.T) {
	t.Run("empty redacting formatter", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{}
		entry := &logrus.Entry{}

		// When
		got, err := f.Format(entry)

		// Then
		require.Error(t, err)
		require.Empty(t, got)
	})

	t.Run("nil fields in formatter and an empty entry produces a basic log entry", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields:    nil,
		}
		entry := &logrus.Entry{}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic\n", string(got))
	})

	t.Run("nil fields and an empty entry produces a basic log entry", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields:    nil,
		}
		entry := &logrus.Entry{}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic\n", string(got))
	})

	t.Run("nil fields and an empty entry produces a basic log entry", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields:    nil,
		}
		entry := &logrus.Entry{}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic\n", string(got))
	})

	t.Run("nil fields and a single field in the entry produces a basic log entry", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields:    nil,
		}
		entry := &logrus.Entry{
			Data: logrus.Fields{
				"example": "value",
			},
		}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic example=value\n", string(got))
	})

	t.Run("single non matching field and a single non matching field in the entry redacts nothing", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields: []*RedactionOption{
				{
					Field:  "password",
					Regexp: nil,
				},
			},
		}
		entry := &logrus.Entry{
			Data: logrus.Fields{
				"example": "value",
			},
		}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic example=value\n", string(got))
	})

	t.Run("single matching is redacted", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields: []*RedactionOption{
				{
					Field:  "example",
					Regexp: nil,
				},
			},
		}
		entry := &logrus.Entry{
			Data: logrus.Fields{
				"example": "value",
			},
		}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic example=\"[REDACTED]\"\n", string(got))
	})

	t.Run("single matching field is redacted with replacement text", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields: []*RedactionOption{
				{
					Field:       "example",
					Regexp:      nil,
					ReplaceWith: "HIDDEN*!£$)*",
				},
			},
		}
		entry := &logrus.Entry{
			Data: logrus.Fields{
				"example": "value",
			},
		}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic example=\"HIDDEN*!£$)*\"\n", string(got))
	})

	t.Run("single matching field is redacted with replacement text", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields: []*RedactionOption{
				{
					Field:       "example",
					Regexp:      nil,
					ReplaceWith: "HIDDEN*!£$)*",
				},
			},
		}
		entry := &logrus.Entry{
			Data: logrus.Fields{
				"example": "value",
			},
		}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic example=\"HIDDEN*!£$)*\"\n", string(got))
	})

	t.Run("single matching field is fully redacted using a regular expression", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields: []*RedactionOption{
				{
					Field:  "example",
					Regexp: regexp.MustCompile(`^redact-[^-]+-this$`),
				},
			},
		}
		entry := &logrus.Entry{
			Data: logrus.Fields{
				"example": "redact-This/is+very#very_secret-this",
			},
		}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic example=\"[REDACTED]\"\n", string(got))
	})

	t.Run("single matching field is NOT redacted using a regular expression that does not match", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields: []*RedactionOption{
				{
					Field:  "example",
					Regexp: regexp.MustCompile(`^redact-[^-]+-what$`),
				},
			},
		}
		entry := &logrus.Entry{
			Data: logrus.Fields{
				"example": "redact-This/is+very#very_secret-this",
			},
		}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic example=\"redact-This/is+very#very_secret-this\"\n", string(got))
	})

	t.Run("single matching field is partially redacted using a regular expression", func(t *testing.T) {
		// Given
		f := &RedactingFormatter{
			Formatter: &logrus.TextFormatter{},
			fields: []*RedactionOption{
				{
					Field:       "example",
					Regexp:      regexp.MustCompile(`^redact-[^-]+-this$`),
					ReplaceWith: "redact-[REDACTED]-this",
				},
			},
		}
		entry := &logrus.Entry{
			Data: logrus.Fields{
				"example": "redact-This/is+very#very_secret-this",
			},
		}

		// When
		got, err := f.Format(entry)

		// Then
		require.NoError(t, err)
		require.Equal(t, "time=\"0001-01-01T00:00:00Z\" level=panic example=\"redact-[REDACTED]-this\"\n", string(got))
	})
}

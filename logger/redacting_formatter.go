package logger

import (
	"github.com/sirupsen/logrus"
	"regexp"
)

type RedactionOption struct {
	Regexp      *regexp.Regexp
	ReplaceWith string
}

type RedactingFormatter struct {
	logrus.Formatter
	fields map[string]*RedactionOption
}

func (f *RedactingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	for field, redactionOption := range f.fields {
		if _, ok := entry.Data[field]; ok {
			replaceWith := "REDACTED"

			if redactionOption == nil {
				entry.Data[field] = replaceWith
				continue
			}

			if redactionOption.ReplaceWith != "" {
				replaceWith = redactionOption.ReplaceWith
			}

			if redactionOption.Regexp == nil {
				entry.Data[field] = replaceWith
				continue
			}

			entry.Data[field] = redactionOption.Regexp.ReplaceAllString(entry.Data[field].(string), replaceWith)
		}
	}
	return f.Formatter.Format(entry)
}

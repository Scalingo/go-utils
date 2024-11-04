package logger

import "github.com/sirupsen/logrus"

type RedactingFormatter struct {
	logrus.Formatter
	fields []string
}

func (f *RedactingFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	for _, field := range f.fields {
		if _, ok := entry.Data[field]; ok {
			entry.Data[field] = "REDACTED"
		}
	}
	return f.Formatter.Format(entry)
}

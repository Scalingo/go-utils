package document

import "bytes"

type ValidationError struct {
	Errors map[string][]string `json:"errors"`
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors: make(map[string][]string),
	}
}

func (v *ValidationError) Error() string {
	var buffer bytes.Buffer

	for field, errors := range v.Errors {
		buffer.WriteString(field)
		buffer.WriteString("=")
		for _, err := range errors {
			buffer.WriteString(err)
			buffer.WriteString(", ")
		}
	}
	return buffer.String()
}

func (v *ValidationError) Set(field, err string) {
	v.Errors[field] = append(v.Errors[field], err)
}

func (v *ValidationError) Get(field string) []string {
	return v.Errors[field]
}

func (v *ValidationError) Build() *ValidationError {
	if len(v.Errors) == 0 {
		return nil
	}

	return v
}

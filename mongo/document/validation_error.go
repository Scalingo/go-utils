package document

import "bytes"

// ValidationErrors is used to store all errors related to the model validation. The typical usecase is:
//	func (m *MyModel) Validate(ctx context.Context) *ValidationErrors {
//		validations := document.NewValidationErrors()
//
//		if m.Name == "" {
//			validations.Set("name", "should not be empty")
//		}
//
//		if m.Email == "" {
//			validations.Set("email", "should not be empty")
//		}
//
//		return validations.Build()
//	}
type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

// NewValidationErrors return an empty ValidationErrors struct
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make(map[string][]string),
	}
}

func (v *ValidationErrors) Error() string {
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

// Set will add an error on a specific field, if the field already contains an error, it will just add it to the current errors list
func (v *ValidationErrors) Set(field, err string) {
	v.Errors[field] = append(v.Errors[field], err)
}

// Get will return all errors set for a specific field
func (v *ValidationErrors) Get(field string) []string {
	return v.Errors[field]
}

// Build will send a ValidationErrors struct if there is some errors or nil if no errors has been defined
func (v *ValidationErrors) Build() *ValidationErrors {
	if len(v.Errors) == 0 {
		return nil
	}

	return v
}

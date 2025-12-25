package conv

import (
	"fmt"
	"strings"
)

func ValidateRequired(requiredFilled map[string]bool) error {
	unfilled := make([]string, 0)
	for fieldName, filled := range requiredFilled {
		if !filled {
			unfilled = append(unfilled, fieldName)
		}
	}

	if len(unfilled) == 0 {
		return nil
	}

	return &RequiredFieldsError{unfilled}
}

type RequiredFieldsError struct {
	fieldsNames []string
}

func (r *RequiredFieldsError) Error() string {
	if len(r.fieldsNames) == 1 {
		return fmt.Sprintf("required field is not filled: %s", r.fieldsNames[0])
	}
	return fmt.Sprintf("required fields are not filled: %s", strings.Join(r.fieldsNames, ", "))
}

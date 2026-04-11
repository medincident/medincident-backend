package organization

import (
	"strings"
	"unicode/utf8"

	"github.com/samber/oops"
)

// Invariant limits for the Organization aggregate's text fields.
const (
	maxNameLen        = 255
	maxDescriptionLen = 2000
)

// Error codes emitted by the Organization field validators.
const (
	ErrCodeOrganizationNameEmpty          = "name_empty"
	ErrCodeOrganizationNameTooLong        = "name_too_long"
	ErrCodeOrganizationDescriptionTooLong = "description_too_long"
)

func validateName(name string) error {
	if name == "" {
		return oops.In("orgstructure.organization").
			Code(ErrCodeOrganizationNameEmpty).
			Public("Organization name is required.").
			With("field", "name").
			Errorf("name is empty")
	}
	if n := utf8.RuneCountInString(name); n > maxNameLen {
		return oops.In("orgstructure.organization").
			Code(ErrCodeOrganizationNameTooLong).
			Public("Organization name is too long.").
			With("field", "name").
			With("actual_length", n).
			With("max_length", maxNameLen).
			Errorf("name too long")
	}
	return nil
}

func validateDescription(desc string) error {
	if desc == "" {
		return nil // optional
	}
	if n := utf8.RuneCountInString(desc); n > maxDescriptionLen {
		return oops.In("orgstructure.organization").
			Code(ErrCodeOrganizationDescriptionTooLong).
			Public("Organization description is too long.").
			With("field", "description").
			With("actual_length", n).
			With("max_length", maxDescriptionLen).
			Errorf("description too long")
	}
	return nil
}

func trim(s string) string { return strings.TrimSpace(s) }

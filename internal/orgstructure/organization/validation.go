package organization

import (
	"unicode/utf8"

	"github.com/samber/oops"
)

// Invariant limits for the Organization aggregate's text fields.
const (
	minNameLen        = 4
	maxNameLen        = 255
	minDescriptionLen = 4
	maxDescriptionLen = 2000
)

// Error codes emitted by the Organization field validators.
const (
	ErrCodeOrganizationNameEmpty           = "name_empty"
	ErrCodeOrganizationNameTooShort        = "name_too_short"
	ErrCodeOrganizationNameTooLong         = "name_too_long"
	ErrCodeOrganizationDescriptionTooShort = "description_too_short"
	ErrCodeOrganizationDescriptionTooLong  = "description_too_long"
)

func validateName(name string) error {
	if name == "" {
		return oops.In("orgstructure.organization").
			Code(ErrCodeOrganizationNameEmpty).
			Public("Organization name is required.").
			With("field", "name").
			Errorf("name is empty")
	}
	n := utf8.RuneCountInString(name)
	if n < minNameLen {
		return oops.In("orgstructure.organization").
			Code(ErrCodeOrganizationNameTooShort).
			Public("Organization name is too short.").
			With("field", "name").
			With("actual_length", n).
			With("min_length", minNameLen).
			Errorf("name too short")
	}
	if n > maxNameLen {
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
	n := utf8.RuneCountInString(desc)
	if n < minDescriptionLen {
		return oops.In("orgstructure.organization").
			Code(ErrCodeOrganizationDescriptionTooShort).
			Public("Organization description is too short.").
			With("field", "description").
			With("actual_length", n).
			With("min_length", minDescriptionLen).
			Errorf("description too short")
	}
	if n > maxDescriptionLen {
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

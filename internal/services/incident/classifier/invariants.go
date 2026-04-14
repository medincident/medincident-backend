package classifier

import (
	"strings"
	"unicode/utf8"

	"github.com/samber/oops"
)

// Length and depth limits for the classifier. Every invariant is a
// named constant.
const (
	incidentCategoryMinNameLen = 2
	incidentCategoryMaxNameLen = 256
	incidentCategoryMinDescLen = 8
	incidentCategoryMaxDescLen = 2048

	incidentTypeMinNameLen = 2
	incidentTypeMaxNameLen = 256
	incidentTypeMinDescLen = 8
	incidentTypeMaxDescLen = 2048

	incidentClassifierMaxDepth = 5
)

// Invariant error codes emitted by the validators below.
const (
	ErrCodeIncidentCategoryNameEmpty           = "incident_category_name_empty"
	ErrCodeIncidentCategoryNameTooShort        = "incident_category_name_too_short"
	ErrCodeIncidentCategoryNameTooLong         = "incident_category_name_too_long"
	ErrCodeIncidentCategoryDescriptionTooShort = "incident_category_description_too_short"
	ErrCodeIncidentCategoryDescriptionTooLong  = "incident_category_description_too_long"

	ErrCodeIncidentTypeNameEmpty           = "incident_type_name_empty"
	ErrCodeIncidentTypeNameTooShort        = "incident_type_name_too_short"
	ErrCodeIncidentTypeNameTooLong         = "incident_type_name_too_long"
	ErrCodeIncidentTypeDescriptionTooShort = "incident_type_description_too_short"
	ErrCodeIncidentTypeDescriptionTooLong  = "incident_type_description_too_long"
)

func validateIncidentCategoryName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return oops.In("services.incident.classifier.category").
			Code(ErrCodeIncidentCategoryNameEmpty).
			Public("Incident category name is required.").
			With("field", "name").
			Errorf("name is empty")
	}
	n := utf8.RuneCountInString(trimmed)
	if n < incidentCategoryMinNameLen {
		return oops.In("services.incident.classifier.category").
			Code(ErrCodeIncidentCategoryNameTooShort).
			Public("Incident category name is too short.").
			With("field", "name").
			With("actual_length", n).
			With("min_length", incidentCategoryMinNameLen).
			Errorf("name too short")
	}
	if n > incidentCategoryMaxNameLen {
		return oops.In("services.incident.classifier.category").
			Code(ErrCodeIncidentCategoryNameTooLong).
			Public("Incident category name is too long.").
			With("field", "name").
			With("actual_length", n).
			With("max_length", incidentCategoryMaxNameLen).
			Errorf("name too long")
	}
	return nil
}

func validateIncidentCategoryDescription(desc *string) error {
	if desc == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*desc)
	if trimmed == "" {
		return nil
	}
	n := utf8.RuneCountInString(trimmed)
	if n < incidentCategoryMinDescLen {
		return oops.In("services.incident.classifier.category").
			Code(ErrCodeIncidentCategoryDescriptionTooShort).
			Public("Incident category description is too short.").
			With("field", "description").
			With("actual_length", n).
			With("min_length", incidentCategoryMinDescLen).
			Errorf("description too short")
	}
	if n > incidentCategoryMaxDescLen {
		return oops.In("services.incident.classifier.category").
			Code(ErrCodeIncidentCategoryDescriptionTooLong).
			Public("Incident category description is too long.").
			With("field", "description").
			With("actual_length", n).
			With("max_length", incidentCategoryMaxDescLen).
			Errorf("description too long")
	}
	return nil
}

func validateIncidentTypeName(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeNameEmpty).
			Public("Incident type name is required.").
			With("field", "name").
			Errorf("name is empty")
	}
	n := utf8.RuneCountInString(trimmed)
	if n < incidentTypeMinNameLen {
		return oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeNameTooShort).
			Public("Incident type name is too short.").
			With("field", "name").
			With("actual_length", n).
			With("min_length", incidentTypeMinNameLen).
			Errorf("name too short")
	}
	if n > incidentTypeMaxNameLen {
		return oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeNameTooLong).
			Public("Incident type name is too long.").
			With("field", "name").
			With("actual_length", n).
			With("max_length", incidentTypeMaxNameLen).
			Errorf("name too long")
	}
	return nil
}

func validateIncidentTypeDescription(desc *string) error {
	if desc == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*desc)
	if trimmed == "" {
		return nil
	}
	n := utf8.RuneCountInString(trimmed)
	if n < incidentTypeMinDescLen {
		return oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeDescriptionTooShort).
			Public("Incident type description is too short.").
			With("field", "description").
			With("actual_length", n).
			With("min_length", incidentTypeMinDescLen).
			Errorf("description too short")
	}
	if n > incidentTypeMaxDescLen {
		return oops.In("services.incident.classifier.type").
			Code(ErrCodeIncidentTypeDescriptionTooLong).
			Public("Incident type description is too long.").
			With("field", "description").
			With("actual_length", n).
			With("max_length", incidentTypeMaxDescLen).
			Errorf("description too long")
	}
	return nil
}

// trimmedStringPtr returns a pointer to the trimmed input iff the
// trimmed value is non-empty; otherwise nil. Used to collapse
// whitespace-only descriptions to absence.
func trimmedStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}

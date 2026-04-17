package classifier

import (
	"errors"
	"strings"
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
)

func oopsCode(t *testing.T, err error) string {
	t.Helper()
	var oe oops.OopsError
	if !errors.As(err, &oe) {
		t.Fatalf("expected oops error, got %T: %v", err, err)
	}
	code, ok := oe.Code().(string)
	if !ok {
		t.Fatalf("oops.Code() returned non-string: %T %v", oe.Code(), oe.Code())
	}
	return code
}

func TestValidateIncidentCategoryName(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"empty", "", ErrCodeIncidentCategoryNameEmpty},
		{"whitespace", "   \t\n", ErrCodeIncidentCategoryNameEmpty},
		{"too short", "A", ErrCodeIncidentCategoryNameTooShort},
		{"ok min", "AB", ""},
		{"ok", "Surgical complications", ""},
		{"ok max", strings.Repeat("A", incidentCategoryMaxNameLen), ""},
		{"too long", strings.Repeat("A", incidentCategoryMaxNameLen+1), ErrCodeIncidentCategoryNameTooLong},
		{"trimmed ok", "  valid  ", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateIncidentCategoryName(tc.input)
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.wantErr, oopsCode(t, err))
			}
		})
	}
}

func TestValidateIncidentCategoryDescription(t *testing.T) {
	desc := func(s string) *string { return &s }
	cases := []struct {
		name    string
		input   *string
		wantErr string
	}{
		{"nil is fine", nil, ""},
		{"whitespace only is rejected", desc("   "), ErrCodeIncidentCategoryDescriptionTooShort},
		{"too short", desc("short"), ErrCodeIncidentCategoryDescriptionTooShort},
		{"ok", desc("long enough description"), ""},
		{"ok max", desc(strings.Repeat("a", incidentCategoryMaxDescLen)), ""},
		{"too long", desc(strings.Repeat("a", incidentCategoryMaxDescLen+1)), ErrCodeIncidentCategoryDescriptionTooLong},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateIncidentCategoryDescription(tc.input)
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.wantErr, oopsCode(t, err))
			}
		})
	}
}

func TestValidateIncidentTypeName(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"empty", "", ErrCodeIncidentTypeNameEmpty},
		{"too short", "X", ErrCodeIncidentTypeNameTooShort},
		{"ok", "Fall with injury", ""},
		{"too long", strings.Repeat("x", incidentTypeMaxNameLen+1), ErrCodeIncidentTypeNameTooLong},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateIncidentTypeName(tc.input)
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.wantErr, oopsCode(t, err))
			}
		})
	}
}

func TestValidateIncidentTypeDescription(t *testing.T) {
	desc := func(s string) *string { return &s }
	cases := []struct {
		name    string
		input   *string
		wantErr string
	}{
		{"nil", nil, ""},
		{"whitespace only is rejected", desc("  "), ErrCodeIncidentTypeDescriptionTooShort},
		{"too short", desc("tiny"), ErrCodeIncidentTypeDescriptionTooShort},
		{"ok", desc("this is a fine description"), ""},
		{"too long", desc(strings.Repeat("y", incidentTypeMaxDescLen+1)), ErrCodeIncidentTypeDescriptionTooLong},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateIncidentTypeDescription(tc.input)
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.wantErr, oopsCode(t, err))
			}
		})
	}
}

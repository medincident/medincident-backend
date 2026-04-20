// Package validation runs struct-tag based validation over command
// structs using go-playground/validator, and translates the library's
// ValidationErrors into oops multi-errors that the gRPC error
// interceptor already knows how to flatten into BadRequest /
// FieldViolation details.
//
// Services call Struct(cmd) as the first thing in every write
// method. The validator instance is a package-level singleton seeded
// at init with a snake_case field-name alias so the "field" key on
// every violation matches the proto field name.
package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/samber/oops"
)

// Oops codes emitted by translate. These are the suffixes the gRPC
// error interceptor already recognises (string_*, uuid_*, time_*,
// int_out_of_range, float_out_of_range) so the mapping to
// codes.InvalidArgument / BadRequest.FieldViolation stays intact.
const (
	CodeStringRequired  = "string_required"
	CodeStringTooShort  = "string_too_short"
	CodeStringTooLong   = "string_too_long"
	CodeUUIDRequired    = "uuid_required"
	CodeTimeRequired    = "time_required"
	CodeIntOutOfRange   = "int_out_of_range"
	CodeFloatOutOfRange = "float_out_of_range"
)

// v is shared across every caller. validator.Validate is documented
// as safe for concurrent Struct calls.
var v = newValidator()

// uuidType caches reflect.TypeFor[uuid.UUID]() so the translator does
// not allocate on every field error.
var uuidType = reflect.TypeFor[uuid.UUID]()

// timeType caches reflect.TypeFor[time.Time]() for the same reason.
var timeType = reflect.TypeFor[time.Time]()

func newValidator() *validator.Validate {
	inst := validator.New(validator.WithRequiredStructEnabled())
	// Expose field names as snake_case in FieldError.Namespace() so
	// clients see "legal_address.point.longitude" instead of
	// "CreateOrganizationCommand.LegalAddress.Point.Longitude".
	inst.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return toSnakeCase(fld.Name)
	})
	return inst
}

// Struct validates cmd against its struct tags. Every string field
// (including pointer-to-string and nested structs) is trimmed in a
// local copy first so that whitespace-only input is treated as empty
// and caught by required / min rules. Every violation is returned as
// errors.Join of oops leaves, each carrying a generic oops Code
// (string_required, string_too_short, float_out_of_range, …) plus a
// dotted "field" context so the existing error interceptor maps them
// to codes.InvalidArgument / BadRequest.FieldViolation without
// changes.
func Struct(cmd any) error {
	rv := reflect.ValueOf(cmd)
	if rv.Kind() == reflect.Pointer {
		trimStrings(rv.Elem())
	} else {
		cp := reflect.New(rv.Type())
		cp.Elem().Set(rv)
		trimStrings(cp.Elem())
		cmd = cp.Elem().Interface()
	}
	if err := v.Struct(cmd); err != nil {
		return translate(err)
	}
	return nil
}

// trimStrings walks v and strings.TrimSpace every settable string,
// including pointer-to-string and fields of nested structs and
// struct pointers. Non-addressable values are silently skipped.
func trimStrings(v reflect.Value) {
	if !v.IsValid() {
		return
	}
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	//nolint:exhaustive // only struct and string need handling; other kinds are intentionally no-op.
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			trimStrings(v.Field(i))
		}
	case reflect.String:
		if v.CanSet() {
			v.SetString(strings.TrimSpace(v.String()))
		}
	}
}

// translate converts a validator error into the oops multi-error
// contract. Non-ValidationErrors (e.g. InvalidValidationError) are
// wrapped as internal faults so they surface as codes.Internal.
func translate(err error) error {
	var ves validator.ValidationErrors
	if !errors.As(err, &ves) {
		return oops.In("validation").
			Code("validate_failed").
			Wrap(err)
	}
	leaves := make([]error, 0, len(ves))
	for _, fe := range ves {
		leaves = append(leaves, fieldErrorToOops(fe))
	}
	return errors.Join(leaves...)
}

func fieldErrorToOops(fe validator.FieldError) error {
	field := fieldPath(fe)
	code, public := classify(fe)
	b := oops.In("validation").
		Code(code).
		With("field", field).
		With("rule", fe.Tag()).
		Public(public)
	if param := fe.Param(); param != "" {
		b = b.With("param", param)
	}
	return b.Errorf("field %q failed rule %q", field, fe.Tag())
}

// fieldPath trims the root struct name ("CreateClinicCommand.") and
// returns a dotted snake_case path suitable for a FieldViolation.
func fieldPath(fe validator.FieldError) string {
	ns := fe.Namespace()
	if _, rest, ok := strings.Cut(ns, "."); ok {
		return rest
	}
	return ns
}

// classify maps (tag, type/kind) → (oops code, public message). The
// codes match the suffixes already configured in the error
// interceptor so the wire format stays intact.
func classify(fe validator.FieldError) (code, public string) {
	switch fe.Tag() {
	case "required":
		switch fe.Type() {
		case uuidType:
			return CodeUUIDRequired, "required"
		case timeType:
			return CodeTimeRequired, "required"
		}
		return CodeStringRequired, "required"
	case "min":
		//nolint:exhaustive // only the kinds that carry a "min" rule matter; others fall through to the generic fallback below.
		switch fe.Kind() {
		case reflect.String:
			return CodeStringTooShort, fmt.Sprintf("length must be at least %s characters", fe.Param())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return CodeIntOutOfRange, fmt.Sprintf("must be at least %s", fe.Param())
		case reflect.Float32, reflect.Float64:
			return CodeFloatOutOfRange, fmt.Sprintf("must be at least %s", fe.Param())
		}
	case "max":
		//nolint:exhaustive // only the kinds that carry a "max" rule matter; others fall through to the generic fallback below.
		switch fe.Kind() {
		case reflect.String:
			return CodeStringTooLong, fmt.Sprintf("length must be at most %s characters", fe.Param())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return CodeIntOutOfRange, fmt.Sprintf("must be at most %s", fe.Param())
		case reflect.Float32, reflect.Float64:
			return CodeFloatOutOfRange, fmt.Sprintf("must be at most %s", fe.Param())
		}
	}
	// Unknown tag. Fall back to a descriptive code so the interceptor
	// still classifies it as InvalidArgument via the _invalid suffix.
	return "rule_" + fe.Tag() + "_invalid", "invalid"
}

// toSnakeCase converts a Go UpperCamel identifier to snake_case while
// keeping runs of uppercase letters together (so "OrganizationID"
// becomes "organization_id" instead of "organization_i_d").
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 4)
	runes := []rune(s)
	for i, r := range runes {
		if i > 0 && unicode.IsUpper(r) {
			prev := runes[i-1]
			var nextLower bool
			if i+1 < len(runes) {
				nextLower = unicode.IsLower(runes[i+1])
			}
			if unicode.IsLower(prev) || nextLower {
				b.WriteByte('_')
			}
		}
		b.WriteRune(unicode.ToLower(r))
	}
	return b.String()
}

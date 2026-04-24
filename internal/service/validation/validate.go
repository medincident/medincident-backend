// Package validation runs struct-tag based validation over command
// structs using go-playground/validator, and translates the library's
// ValidationErrors into a single oops error with code
// CodeValidationFailed whose context carries a typed []Violation slice.
//
// The gRPC error interceptor recognises that code and unpacks the
// violations into a codes.InvalidArgument status with one
// BadRequest.FieldViolation per violation. The violation Rule is the
// raw validator tag ("required", "min", "max", "uuid", …) — the field
// path already implies the type, so no type-prefixed codes are needed.
//
// Services call Struct(cmd) as the first thing in every write method.
// The validator instance is a package-level singleton seeded at init
// with a snake_case field-name alias so the Field on every violation
// matches the proto field name.
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

// CodeValidationFailed is the single oops code emitted for every
// struct-tag validation failure. The gRPC error interceptor maps it to
// codes.InvalidArgument and reads ContextKeyViolations to build
// BadRequest.FieldViolation details.
const CodeValidationFailed = "validation_failed"

// CodeValidatorInvocationFailed is emitted when go-playground/validator
// itself refuses the input (e.g. validator.InvalidValidationError when
// a non-struct is passed to Struct). This is a developer bug, never a
// client input error, so the interceptor leaves it to fall through to
// the default codes.Internal mapping — no explicit override is needed.
const CodeValidatorInvocationFailed = "validator_invocation_failed"

// ContextKeyViolations is the oops-context key under which translate
// stashes the []Violation produced from a validator.ValidationErrors.
// The interceptor looks this up by name to stay decoupled from the
// exact slice element type at the call site.
const ContextKeyViolations = "violations"

// Violation describes one struct-tag rule failure on a specific field.
// Field is a dotted snake_case path (e.g. "legal_address.point.longitude"),
// Rule is the raw validator tag ("required", "min", "max", "uuid", …),
// Param is the tag parameter when present (e.g. "4" for min=4), and
// Message is a short client-facing description derived from the rule
// and field kind.
type Violation struct {
	Field   string
	Rule    string
	Param   string
	Message string
}

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
// and caught by required / min rules. A failure is returned as a
// single oops error with CodeValidationFailed whose context carries
// ContextKeyViolations → []Violation.
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
		for _, f := range v.Fields() {
			trimStrings(f)
		}
	case reflect.String:
		if v.CanSet() {
			v.SetString(strings.TrimSpace(v.String()))
		}
	}
}

// translate converts a validator error into a single oops error with
// code CodeValidationFailed for client-side field violations. Anything
// else returned by the validator (e.g. InvalidValidationError when a
// non-struct is passed to Struct) is a developer misuse: it is tagged
// with CodeValidatorInvocationFailed so the interceptor surfaces it as
// codes.Internal rather than leaking a spurious InvalidArgument.
func translate(err error) error {
	var ves validator.ValidationErrors
	if !errors.As(err, &ves) {
		return oops.In("validation").
			Code(CodeValidatorInvocationFailed).
			Wrap(err)
	}
	violations := make([]Violation, 0, len(ves))
	for _, fe := range ves {
		violations = append(violations, toViolation(fe))
	}
	return oops.In("validation").
		Code(CodeValidationFailed).
		With(ContextKeyViolations, violations).
		Public("request is invalid").
		Errorf("validation failed: %d violation(s)", len(violations))
}

func toViolation(fe validator.FieldError) Violation {
	return Violation{
		Field:   fieldPath(fe),
		Rule:    fe.Tag(),
		Param:   fe.Param(),
		Message: message(fe),
	}
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

// message returns a short client-facing description for fe. The
// caller-visible rule name is already carried separately; this string
// only adds bounds or format hints a human needs.
func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "uuid":
		return "must be a valid uuid"
	case "min":
		//nolint:exhaustive // only the kinds that carry a "min" rule matter; others fall through to the generic fallback below.
		switch fe.Kind() {
		case reflect.String:
			return fmt.Sprintf("length must be at least %s characters", fe.Param())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return fmt.Sprintf("must be at least %s", fe.Param())
		}
	case "max":
		//nolint:exhaustive // only the kinds that carry a "max" rule matter; others fall through to the generic fallback below.
		switch fe.Kind() {
		case reflect.String:
			return fmt.Sprintf("length must be at most %s characters", fe.Param())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return fmt.Sprintf("must be at most %s", fe.Param())
		}
	}
	// Unknown tag or tag on an unsupported kind: typed presence rules
	// on uuid.UUID / time.Time surface here as "required" because
	// validator.WithRequiredStructEnabled reports them with a non-"required"
	// tag. Falling back to the tag itself keeps the message honest.
	if fe.Type() == uuidType || fe.Type() == timeType {
		return "required"
	}
	return fe.Tag()
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

// Package zitadel provides a thin abstraction around Zitadel user
// verification. The command-service calls it from HireEmployee to
// confirm that a Zitadel user exists before persisting an employee
// row. Production uses JWTProfileVerifier; tests can use
// StubAlwaysExists for isolation.
package zitadel

import (
	"context"
	"errors"
)

// ErrUserNotFound is returned by UserVerifier.Verify when the Zitadel
// instance reports that the given user ID does not exist. Any other
// error (network, auth, unexpected status) is wrapped and returned
// as-is by the caller.
var ErrUserNotFound = errors.New("zitadel: user not found")

// UserVerifier checks whether a Zitadel user exists. Implementations
// must be safe for concurrent use.
type UserVerifier interface {
	Verify(ctx context.Context, zitadelUserID string) error
}

// StubAlwaysExists returns a UserVerifier that always reports the user
// as existing. Intended for unit tests of non-Hire commands where the
// Zitadel call is irrelevant.
func StubAlwaysExists() UserVerifier { return stubAlwaysExists{} }

type stubAlwaysExists struct{}

func (stubAlwaysExists) Verify(context.Context, string) error { return nil }

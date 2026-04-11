// Package organization holds integration tests for the organization
// aggregate's persistence path. All actual test files live in the
// external test package organization_test and are guarded by
// //go:build integration; this file exists (without a build tag) so
// that `go test ./...` without the integration tag can still load the
// directory instead of failing with "build constraints exclude all Go
// files".
package organization

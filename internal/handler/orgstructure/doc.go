// Package orgstructure is the gRPC handler for OrgStructureService.
// Each RPC method lives in its own file (<aggregate>_<method>.go).
// Handlers are strictly thin: unpack proto request, call a service,
// pack proto response. No validation, no error mapping — service-layer
// errors propagate as-is for now and will be translated to gRPC status
// codes by a middleware added in a follow-up spec.
package orgstructure

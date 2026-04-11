// Package organizationinfra holds the infrastructure adapters for the
// Organization aggregate. Today that's just the outbox EventInfo
// registrations: each domain event struct is mapped to its proto wire
// form and registered with outbox.Registry at DI startup.
//
// This is the ONLY file in the orgstructure/organization/ subtree that
// imports gen/... proto packages. The domain, the application service,
// and the postgres repository all know nothing about proto.
package organizationinfra

import (
	"reflect"

	"google.golang.org/protobuf/proto"

	orgstructurev1 "github.com/medincident/medincident-command-service/gen/medincident/orgstructure/v1"
	geov1 "github.com/medincident/medincident-command-service/gen/medincident/shared/geo/v1"

	"github.com/medincident/medincident-command-service/internal/orgstructure/organization"
	"github.com/medincident/medincident-command-service/internal/outbox"
	"github.com/medincident/medincident-command-service/internal/shared/geo"
)

// --- per-event mappers ---

func createdToProto(ev any) (proto.Message, error) {
	e := ev.(*organization.Created)
	out := &orgstructurev1.OrganizationCreated{
		Name:         e.Name,
		LegalAddress: addressToProto(e.LegalAddress),
	}
	if e.Description != "" {
		desc := e.Description
		out.Description = &desc
	}
	return out, nil
}

func renamedToProto(ev any) (proto.Message, error) {
	e := ev.(*organization.Renamed)
	return &orgstructurev1.OrganizationRenamed{Name: e.Name}, nil
}

func descriptionUpdatedToProto(ev any) (proto.Message, error) {
	e := ev.(*organization.DescriptionUpdated)
	out := &orgstructurev1.OrganizationDescriptionUpdated{}
	if e.Description != "" {
		desc := e.Description
		out.Description = &desc
	}
	return out, nil
}

func legalAddressRelocatedToProto(ev any) (proto.Message, error) {
	e := ev.(*organization.LegalAddressRelocated)
	return &orgstructurev1.OrganizationLegalAddressRelocated{
		LegalAddress: addressToProto(e.LegalAddress),
	}, nil
}

// addressToProto converts the shared VO to its proto form. nil domain
// address → nil proto (proto3 optional absent).
func addressToProto(a *geo.Address) *geov1.Address {
	if a == nil {
		return nil
	}
	out := &geov1.Address{Text: a.Text}
	if a.Point != nil {
		out.Point = &geov1.Point{Lng: a.Point.Longitude, Lat: a.Point.Latitude}
	}
	return out
}

// --- registry entries ---

// eventInfos is the canonical list of Organization events with their
// outbox registry metadata. Keeping it as a var slice (rather than
// inlining each Register call) lets the test suite iterate over it to
// assert coverage.
var eventInfos = []struct {
	Type reflect.Type
	Info outbox.EventInfo
}{
	{
		Type: reflect.TypeOf(&organization.Created{}),
		Info: outbox.EventInfo{
			TypeName: "medincident.orgstructure.v1.OrganizationCreated",
			Zero:     func() any { return &organization.Created{} },
			ToProto:  createdToProto,
		},
	},
	{
		Type: reflect.TypeOf(&organization.Renamed{}),
		Info: outbox.EventInfo{
			TypeName: "medincident.orgstructure.v1.OrganizationRenamed",
			Zero:     func() any { return &organization.Renamed{} },
			ToProto:  renamedToProto,
		},
	},
	{
		Type: reflect.TypeOf(&organization.DescriptionUpdated{}),
		Info: outbox.EventInfo{
			TypeName: "medincident.orgstructure.v1.OrganizationDescriptionUpdated",
			Zero:     func() any { return &organization.DescriptionUpdated{} },
			ToProto:  descriptionUpdatedToProto,
		},
	},
	{
		Type: reflect.TypeOf(&organization.LegalAddressRelocated{}),
		Info: outbox.EventInfo{
			TypeName: "medincident.orgstructure.v1.OrganizationLegalAddressRelocated",
			Zero:     func() any { return &organization.LegalAddressRelocated{} },
			ToProto:  legalAddressRelocatedToProto,
		},
	},
}

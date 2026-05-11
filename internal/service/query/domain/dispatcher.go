package domain

import (
	"context"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
	"github.com/samber/oops"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"

	clinicv1 "github.com/medincident/medincident-backend/pkg/event/clinic/v1"
	deptv1 "github.com/medincident/medincident-backend/pkg/event/department/v1"
	empv1 "github.com/medincident/medincident-backend/pkg/event/employee/v1"
	bufferv1 "github.com/medincident/medincident-backend/pkg/event/incident/buffer/v1"
	classifierv1 "github.com/medincident/medincident-backend/pkg/event/incident/classifier/v1"
	incidentv1 "github.com/medincident/medincident-backend/pkg/event/incident/v1"
	orgv1 "github.com/medincident/medincident-backend/pkg/event/organization/v1"
	requesttypev1 "github.com/medincident/medincident-backend/pkg/event/request_type/v1"
	servicerequestv1 "github.com/medincident/medincident-backend/pkg/event/service_request/v1"
	sav1 "github.com/medincident/medincident-backend/pkg/event/system_admin/v1"
	eventv1 "github.com/medincident/medincident-backend/pkg/event/v1"
	vacv1 "github.com/medincident/medincident-backend/pkg/event/vacation/v1"

	"github.com/medincident/medincident-backend/internal/service/query/projector"
)

// Dispatcher decodes an Envelope, resolves the event type from the NATS
// subject, and calls the matching projector function inside a DB transaction.
type Dispatcher struct {
	db   *gorm.DB
	proj *projector.Projectors
	log  *zerolog.Logger
}

// NewDispatcher returns a Dispatcher wired to the given DB and projector set.
func NewDispatcher(db *gorm.DB, proj *projector.Projectors, log *zerolog.Logger) *Dispatcher {
	return &Dispatcher{db: db, proj: proj, log: log}
}

// Dispatch decodes msg, applies the projector inside a transaction, and
// returns nil on success. Returns a malformed error (_malformed code) for
// permanent decode failures; transient DB errors are returned unwrapped so
// the caller naks with backoff.
func (d *Dispatcher) Dispatch(ctx context.Context, msg jetstream.Msg) error {
	env := new(eventv1.Envelope)
	if err := proto.Unmarshal(msg.Data(), env); err != nil {
		return oops.In("consumer.domain").
			Code(ErrCodeEnvelopeUnmarshalFailed).
			With("subject", msg.Subject()).
			Wrap(err)
	}
	aggregateID := env.GetAggregateId()
	occurredAt := envelopeTime(env)

	payload := env.GetPayload()
	if payload == nil {
		return oops.In("consumer.domain").
			Code(ErrCodePayloadUnmarshalFailed).
			With("subject", msg.Subject()).
			Errorf("nil payload")
	}
	evMsg, err := payload.UnmarshalNew()
	if err != nil {
		return oops.In("consumer.domain").
			Code(ErrCodePayloadUnmarshalFailed).
			With("subject", msg.Subject()).
			Wrap(err)
	}

	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		switch m := evMsg.(type) {
		// ── Organization ──────────────────────────────────────────────────
		case *orgv1.OrganizationCreated:
			return d.proj.OrganizationCreated(tx, aggregateID, occurredAt, m)
		case *orgv1.OrganizationDetailsChanged:
			return d.proj.OrganizationDetailsChanged(tx, aggregateID, occurredAt, m)
		case *orgv1.OrganizationLegalAddressChanged:
			return d.proj.OrganizationLegalAddressChanged(tx, aggregateID, occurredAt, m)

		// ── Clinic ────────────────────────────────────────────────────────
		case *clinicv1.ClinicCreated:
			return d.proj.ClinicCreated(tx, aggregateID, occurredAt, m)
		case *clinicv1.ClinicDetailsChanged:
			return d.proj.ClinicDetailsChanged(tx, aggregateID, occurredAt, m)
		case *clinicv1.ClinicPhysicalAddressChanged:
			return d.proj.ClinicPhysicalAddressChanged(tx, aggregateID, occurredAt, m)

		// ── Department ────────────────────────────────────────────────────
		case *deptv1.DepartmentCreated:
			return d.proj.DepartmentCreated(tx, aggregateID, occurredAt, m)
		case *deptv1.DepartmentDetailsChanged:
			return d.proj.DepartmentDetailsChanged(tx, aggregateID, occurredAt, m)

		// ── Employee ──────────────────────────────────────────────────────
		case *empv1.EmployeeHired:
			return d.proj.EmployeeHired(tx, aggregateID, occurredAt, m)
		case *empv1.EmployeeTerminated:
			return d.proj.EmployeeTerminated(tx, aggregateID, occurredAt, m)
		case *empv1.EmployeeDepartmentChanged:
			return d.proj.EmployeeDepartmentChanged(tx, aggregateID, occurredAt, m)
		case *empv1.EmployeePositionChanged:
			return d.proj.EmployeePositionChanged(tx, aggregateID, occurredAt, m)

		// ── Vacation ──────────────────────────────────────────────────────
		case *vacv1.VacationScheduled:
			return d.proj.VacationScheduled(tx, aggregateID, occurredAt, m)
		case *vacv1.VacationStarted:
			return d.proj.VacationStarted(tx, aggregateID, occurredAt, m)
		case *vacv1.VacationEnded:
			return d.proj.VacationEnded(tx, aggregateID, occurredAt, m)
		case *vacv1.VacationCancelled:
			return d.proj.VacationCancelled(tx, aggregateID, occurredAt, m)
		case *vacv1.VacationEndDateChanged:
			return d.proj.VacationEndDateChanged(tx, aggregateID, occurredAt, m)

		// ── System Admin ──────────────────────────────────────────────────
		case *sav1.SystemAdminGranted:
			return d.proj.SystemAdminGranted(tx, aggregateID, occurredAt, m)
		case *sav1.SystemAdminRevoked:
			return d.proj.SystemAdminRevoked(tx, aggregateID, occurredAt, m)

		// ── Org roles ─────────────────────────────────────────────────────
		case *orgv1.OrgAdminAssigned:
			return d.proj.OrgAdminAssigned(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgAdminDeputyAssigned:
			return d.proj.OrgAdminDeputyAssigned(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgAdminDeputyRemoved:
			return d.proj.OrgAdminDeputyRemoved(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgAdminRevoked:
			return d.proj.OrgAdminRevoked(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgHeadAssigned:
			return d.proj.OrgHeadAssigned(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgHeadDeputyAssigned:
			return d.proj.OrgHeadDeputyAssigned(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgHeadDeputyRemoved:
			return d.proj.OrgHeadDeputyRemoved(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgHeadRevoked:
			return d.proj.OrgHeadRevoked(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgDispatcherAssigned:
			return d.proj.OrgDispatcherAssigned(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgDispatcherDeputyAssigned:
			return d.proj.OrgDispatcherDeputyAssigned(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgDispatcherDeputyRemoved:
			return d.proj.OrgDispatcherDeputyRemoved(tx, aggregateID, occurredAt, m)
		case *orgv1.OrgDispatcherRevoked:
			return d.proj.OrgDispatcherRevoked(tx, aggregateID, occurredAt, m)

		// ── Clinic Head ───────────────────────────────────────────────────
		case *clinicv1.ClinicHeadAssigned:
			return d.proj.ClinicHeadAssigned(tx, aggregateID, occurredAt, m)
		case *clinicv1.ClinicHeadDeputyAssigned:
			return d.proj.ClinicHeadDeputyAssigned(tx, aggregateID, occurredAt, m)
		case *clinicv1.ClinicHeadDeputyRemoved:
			return d.proj.ClinicHeadDeputyRemoved(tx, aggregateID, occurredAt, m)
		case *clinicv1.ClinicHeadRevoked:
			return d.proj.ClinicHeadRevoked(tx, aggregateID, occurredAt, m)

		// ── Dept Responsible ──────────────────────────────────────────────
		case *deptv1.DeptResponsibleAssigned:
			return d.proj.DeptResponsibleAssigned(tx, aggregateID, occurredAt, m)
		case *deptv1.DeptResponsibleDeputyAssigned:
			return d.proj.DeptResponsibleDeputyAssigned(tx, aggregateID, occurredAt, m)
		case *deptv1.DeptResponsibleDeputyRemoved:
			return d.proj.DeptResponsibleDeputyRemoved(tx, aggregateID, occurredAt, m)
		case *deptv1.DeptResponsibleRevoked:
			return d.proj.DeptResponsibleRevoked(tx, aggregateID, occurredAt, m)

		// ── Incident Category ─────────────────────────────────────────────────
		case *classifierv1.IncidentCategoryCreated:
			return d.proj.CategoryCreated(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentCategoryDetailsUpdated:
			return d.proj.CategoryDetailsUpdated(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentCategoryMoved:
			return d.proj.CategoryMoved(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentCategoryDeactivated:
			return d.proj.CategoryDeactivated(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentCategoryReactivated:
			return d.proj.CategoryReactivated(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentCategoryDeleted:
			return d.proj.CategoryDeleted(tx, aggregateID, occurredAt, m)

		// ── Incident Type ─────────────────────────────────────────────────────
		case *classifierv1.IncidentTypeCreated:
			return d.proj.TypeCreated(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentTypeDetailsUpdated:
			return d.proj.TypeDetailsUpdated(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentTypeMoved:
			return d.proj.TypeMoved(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentTypeDeactivated:
			return d.proj.TypeDeactivated(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentTypeReactivated:
			return d.proj.TypeReactivated(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentTypeAllowedForPatients:
			return d.proj.TypeAllowedForPatients(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentTypeDisallowedForPatients:
			return d.proj.TypeDisallowedForPatients(tx, aggregateID, occurredAt, m)
		case *classifierv1.IncidentTypeDeleted:
			return d.proj.TypeDeleted(tx, aggregateID, occurredAt, m)

		// ── Incident ──────────────────────────────────────────────────────────
		case *incidentv1.IncidentCreated:
			return d.proj.IncidentCreated(tx, aggregateID, occurredAt, m)
		case *incidentv1.IncidentStatusChanged:
			return d.proj.IncidentStatusChanged(tx, aggregateID, occurredAt, m)
		case *incidentv1.IncidentPriorityChanged:
			return d.proj.IncidentPriorityChanged(tx, aggregateID, occurredAt, m)
		case *incidentv1.IncidentDescriptionUpdated:
			return d.proj.IncidentDescriptionUpdated(tx, aggregateID, occurredAt, m)

		// ── Patient Incident Buffer ───────────────────────────────────────────
		case *bufferv1.PatientIncidentBufferCreated:
			return d.proj.PatientIncidentBufferCreated(tx, aggregateID, occurredAt, m)
		case *bufferv1.PatientIncidentBufferUpdated:
			return d.proj.PatientIncidentBufferUpdated(tx, aggregateID, occurredAt, m)

		// ── Request Type ──────────────────────────────────────────────────────
		case *requesttypev1.RequestTypeCreated:
			return d.proj.RequestTypeCreated(tx, aggregateID, occurredAt, m)
		case *requesttypev1.RequestTypeDetailsUpdated:
			return d.proj.RequestTypeDetailsUpdated(tx, aggregateID, occurredAt, m)
		case *requesttypev1.RequestTypeDeactivated:
			return d.proj.RequestTypeDeactivated(tx, aggregateID, occurredAt, m)
		case *requesttypev1.RequestTypeReactivated:
			return d.proj.RequestTypeReactivated(tx, aggregateID, occurredAt, m)
		case *requesttypev1.RequestTypeDeleted:
			return d.proj.RequestTypeDeleted(tx, aggregateID, occurredAt, m)

		// ── Service Request ───────────────────────────────────────────────────
		case *servicerequestv1.ServiceRequestCreated:
			return d.proj.ServiceRequestCreated(tx, aggregateID, occurredAt, m)
		case *servicerequestv1.ServiceRequestStatusChanged:
			return d.proj.ServiceRequestStatusChanged(tx, aggregateID, occurredAt, m)
		case *servicerequestv1.ServiceRequestDescriptionUpdated:
			return d.proj.ServiceRequestDescriptionUpdated(tx, aggregateID, occurredAt, m)
		case *servicerequestv1.ServiceRequestExecutorAssigned:
			return d.proj.ServiceRequestExecutorAssigned(tx, aggregateID, occurredAt, m)
		case *servicerequestv1.ServiceRequestExecutorRemoved:
			return d.proj.ServiceRequestExecutorRemoved(tx, aggregateID, occurredAt, m)

		default:
			d.log.Warn().
				Str("subject", msg.Subject()).
				Str("type", payload.GetTypeUrl()).
				Msg("unknown domain event type; skipping")
			return nil
		}
	})
}

// envelopeTime converts the envelope's occurred_at timestamp to time.Time,
// defaulting to time.Now().UTC() when the publisher forgot to set one.
func envelopeTime(env *eventv1.Envelope) time.Time {
	ts := env.GetOccurredAt()
	if ts == nil || !ts.IsValid() {
		return time.Now().UTC()
	}
	return ts.AsTime().UTC()
}

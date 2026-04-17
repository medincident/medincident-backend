package classifier

// NATS subjects for every classifier event.
const (
	SubjectIncidentCategoryCreated        = "medincident.event.incident.category.v1.created"
	SubjectIncidentCategoryDetailsChanged = "medincident.event.incident.category.v1.details_changed"
	SubjectIncidentCategoryMoved          = "medincident.event.incident.category.v1.moved"
	SubjectIncidentCategoryDeactivated    = "medincident.event.incident.category.v1.deactivated"
	SubjectIncidentCategoryReactivated    = "medincident.event.incident.category.v1.reactivated"
	SubjectIncidentCategoryDeleted        = "medincident.event.incident.category.v1.deleted"

	SubjectIncidentTypeCreated        = "medincident.event.incident.type.v1.created"
	SubjectIncidentTypeDetailsChanged = "medincident.event.incident.type.v1.details_changed"
	SubjectIncidentTypeMoved          = "medincident.event.incident.type.v1.moved"
	SubjectIncidentTypeDeactivated    = "medincident.event.incident.type.v1.deactivated"
	SubjectIncidentTypeReactivated    = "medincident.event.incident.type.v1.reactivated"
	SubjectIncidentTypeDeleted        = "medincident.event.incident.type.v1.deleted"
)

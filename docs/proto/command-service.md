# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [event/clinic/v1/events.proto](#event_clinic_v1_events-proto)
    - [Address](#event-clinic-v1-Address)
    - [ClinicCreated](#event-clinic-v1-ClinicCreated)
    - [ClinicDetailsChanged](#event-clinic-v1-ClinicDetailsChanged)
    - [ClinicHeadAssigned](#event-clinic-v1-ClinicHeadAssigned)
    - [ClinicHeadDeputyAssigned](#event-clinic-v1-ClinicHeadDeputyAssigned)
    - [ClinicHeadDeputyRemoved](#event-clinic-v1-ClinicHeadDeputyRemoved)
    - [ClinicHeadRevoked](#event-clinic-v1-ClinicHeadRevoked)
    - [ClinicPhysicalAddressChanged](#event-clinic-v1-ClinicPhysicalAddressChanged)
    - [Point](#event-clinic-v1-Point)

- [event/department/v1/events.proto](#event_department_v1_events-proto)
    - [DepartmentCreated](#event-department-v1-DepartmentCreated)
    - [DepartmentDetailsChanged](#event-department-v1-DepartmentDetailsChanged)
    - [DepartmentResponsibleAssigned](#event-department-v1-DepartmentResponsibleAssigned)
    - [DepartmentResponsibleDeputyAssigned](#event-department-v1-DepartmentResponsibleDeputyAssigned)
    - [DepartmentResponsibleDeputyRemoved](#event-department-v1-DepartmentResponsibleDeputyRemoved)
    - [DepartmentResponsibleRevoked](#event-department-v1-DepartmentResponsibleRevoked)

- [event/employee/v1/events.proto](#event_employee_v1_events-proto)
    - [EmployeeDepartmentChanged](#event-employee-v1-EmployeeDepartmentChanged)
    - [EmployeeHired](#event-employee-v1-EmployeeHired)
    - [EmployeePositionChanged](#event-employee-v1-EmployeePositionChanged)
    - [EmployeeTerminated](#event-employee-v1-EmployeeTerminated)
    - [VacationCancelled](#event-employee-v1-VacationCancelled)
    - [VacationEndDateChanged](#event-employee-v1-VacationEndDateChanged)
    - [VacationEnded](#event-employee-v1-VacationEnded)
    - [VacationScheduled](#event-employee-v1-VacationScheduled)
    - [VacationStarted](#event-employee-v1-VacationStarted)

- [event/incident/category/v1/events.proto](#event_incident_category_v1_events-proto)
    - [IncidentCategoryCreated](#event-incident-category-v1-IncidentCategoryCreated)
    - [IncidentCategoryDeactivated](#event-incident-category-v1-IncidentCategoryDeactivated)
    - [IncidentCategoryDeleted](#event-incident-category-v1-IncidentCategoryDeleted)
    - [IncidentCategoryDetailsChanged](#event-incident-category-v1-IncidentCategoryDetailsChanged)
    - [IncidentCategoryMoved](#event-incident-category-v1-IncidentCategoryMoved)
    - [IncidentCategoryReactivated](#event-incident-category-v1-IncidentCategoryReactivated)

- [event/incident/type/v1/events.proto](#event_incident_type_v1_events-proto)
    - [IncidentTypeCreated](#event-incident-type-v1-IncidentTypeCreated)
    - [IncidentTypeDeactivated](#event-incident-type-v1-IncidentTypeDeactivated)
    - [IncidentTypeDeleted](#event-incident-type-v1-IncidentTypeDeleted)
    - [IncidentTypeDetailsChanged](#event-incident-type-v1-IncidentTypeDetailsChanged)
    - [IncidentTypeMoved](#event-incident-type-v1-IncidentTypeMoved)
    - [IncidentTypeReactivated](#event-incident-type-v1-IncidentTypeReactivated)

- [event/organization/v1/events.proto](#event_organization_v1_events-proto)
    - [Address](#event-organization-v1-Address)
    - [OrganizationAdminAssigned](#event-organization-v1-OrganizationAdminAssigned)
    - [OrganizationAdminDeputyAssigned](#event-organization-v1-OrganizationAdminDeputyAssigned)
    - [OrganizationAdminDeputyRemoved](#event-organization-v1-OrganizationAdminDeputyRemoved)
    - [OrganizationAdminRevoked](#event-organization-v1-OrganizationAdminRevoked)
    - [OrganizationCreated](#event-organization-v1-OrganizationCreated)
    - [OrganizationDetailsChanged](#event-organization-v1-OrganizationDetailsChanged)
    - [OrganizationDispatcherAssigned](#event-organization-v1-OrganizationDispatcherAssigned)
    - [OrganizationDispatcherDeputyAssigned](#event-organization-v1-OrganizationDispatcherDeputyAssigned)
    - [OrganizationDispatcherDeputyRemoved](#event-organization-v1-OrganizationDispatcherDeputyRemoved)
    - [OrganizationDispatcherRevoked](#event-organization-v1-OrganizationDispatcherRevoked)
    - [OrganizationHeadAssigned](#event-organization-v1-OrganizationHeadAssigned)
    - [OrganizationHeadDeputyAssigned](#event-organization-v1-OrganizationHeadDeputyAssigned)
    - [OrganizationHeadDeputyRemoved](#event-organization-v1-OrganizationHeadDeputyRemoved)
    - [OrganizationHeadRevoked](#event-organization-v1-OrganizationHeadRevoked)
    - [OrganizationLegalAddressChanged](#event-organization-v1-OrganizationLegalAddressChanged)
    - [Point](#event-organization-v1-Point)

- [event/system_admin/v1/events.proto](#event_system_admin_v1_events-proto)
    - [SystemAdminGranted](#event-system_admin-v1-SystemAdminGranted)
    - [SystemAdminRevoked](#event-system_admin-v1-SystemAdminRevoked)

- [event/v1/envelope.proto](#event_v1_envelope-proto)
    - [Envelope](#event-v1-Envelope)

- [service/incident/classifier/v1/incident_classifier.proto](#service_incident_classifier_v1_incident_classifier-proto)
    - [CreateIncidentCategoryRequest](#service-incident-classifier-v1-CreateIncidentCategoryRequest)
    - [CreateIncidentCategoryResponse](#service-incident-classifier-v1-CreateIncidentCategoryResponse)
    - [CreateIncidentTypeRequest](#service-incident-classifier-v1-CreateIncidentTypeRequest)
    - [CreateIncidentTypeResponse](#service-incident-classifier-v1-CreateIncidentTypeResponse)
    - [DeactivateIncidentCategoryRequest](#service-incident-classifier-v1-DeactivateIncidentCategoryRequest)
    - [DeactivateIncidentCategoryResponse](#service-incident-classifier-v1-DeactivateIncidentCategoryResponse)
    - [DeactivateIncidentTypeRequest](#service-incident-classifier-v1-DeactivateIncidentTypeRequest)
    - [DeactivateIncidentTypeResponse](#service-incident-classifier-v1-DeactivateIncidentTypeResponse)
    - [DeleteIncidentCategoryRequest](#service-incident-classifier-v1-DeleteIncidentCategoryRequest)
    - [DeleteIncidentCategoryResponse](#service-incident-classifier-v1-DeleteIncidentCategoryResponse)
    - [DeleteIncidentTypeRequest](#service-incident-classifier-v1-DeleteIncidentTypeRequest)
    - [DeleteIncidentTypeResponse](#service-incident-classifier-v1-DeleteIncidentTypeResponse)
    - [MoveIncidentCategoryRequest](#service-incident-classifier-v1-MoveIncidentCategoryRequest)
    - [MoveIncidentCategoryResponse](#service-incident-classifier-v1-MoveIncidentCategoryResponse)
    - [MoveIncidentTypeRequest](#service-incident-classifier-v1-MoveIncidentTypeRequest)
    - [MoveIncidentTypeResponse](#service-incident-classifier-v1-MoveIncidentTypeResponse)
    - [ReactivateIncidentCategoryRequest](#service-incident-classifier-v1-ReactivateIncidentCategoryRequest)
    - [ReactivateIncidentCategoryResponse](#service-incident-classifier-v1-ReactivateIncidentCategoryResponse)
    - [ReactivateIncidentTypeRequest](#service-incident-classifier-v1-ReactivateIncidentTypeRequest)
    - [ReactivateIncidentTypeResponse](#service-incident-classifier-v1-ReactivateIncidentTypeResponse)
    - [UpdateIncidentCategoryDetailsRequest](#service-incident-classifier-v1-UpdateIncidentCategoryDetailsRequest)
    - [UpdateIncidentCategoryDetailsResponse](#service-incident-classifier-v1-UpdateIncidentCategoryDetailsResponse)
    - [UpdateIncidentTypeDetailsRequest](#service-incident-classifier-v1-UpdateIncidentTypeDetailsRequest)
    - [UpdateIncidentTypeDetailsResponse](#service-incident-classifier-v1-UpdateIncidentTypeDetailsResponse)

    - [IncidentClassifierService](#service-incident-classifier-v1-IncidentClassifierService)

- [service/membership/v1/membership.proto](#service_membership_v1_membership-proto)
    - [AssignClinicHeadDeputyRequest](#service-membership-v1-AssignClinicHeadDeputyRequest)
    - [AssignClinicHeadDeputyResponse](#service-membership-v1-AssignClinicHeadDeputyResponse)
    - [AssignClinicHeadRequest](#service-membership-v1-AssignClinicHeadRequest)
    - [AssignClinicHeadResponse](#service-membership-v1-AssignClinicHeadResponse)
    - [AssignDepartmentResponsibleDeputyRequest](#service-membership-v1-AssignDepartmentResponsibleDeputyRequest)
    - [AssignDepartmentResponsibleDeputyResponse](#service-membership-v1-AssignDepartmentResponsibleDeputyResponse)
    - [AssignDepartmentResponsibleRequest](#service-membership-v1-AssignDepartmentResponsibleRequest)
    - [AssignDepartmentResponsibleResponse](#service-membership-v1-AssignDepartmentResponsibleResponse)
    - [AssignOrganizationAdminDeputyRequest](#service-membership-v1-AssignOrganizationAdminDeputyRequest)
    - [AssignOrganizationAdminDeputyResponse](#service-membership-v1-AssignOrganizationAdminDeputyResponse)
    - [AssignOrganizationAdminRequest](#service-membership-v1-AssignOrganizationAdminRequest)
    - [AssignOrganizationAdminResponse](#service-membership-v1-AssignOrganizationAdminResponse)
    - [AssignOrganizationDispatcherDeputyRequest](#service-membership-v1-AssignOrganizationDispatcherDeputyRequest)
    - [AssignOrganizationDispatcherDeputyResponse](#service-membership-v1-AssignOrganizationDispatcherDeputyResponse)
    - [AssignOrganizationDispatcherRequest](#service-membership-v1-AssignOrganizationDispatcherRequest)
    - [AssignOrganizationDispatcherResponse](#service-membership-v1-AssignOrganizationDispatcherResponse)
    - [AssignOrganizationHeadDeputyRequest](#service-membership-v1-AssignOrganizationHeadDeputyRequest)
    - [AssignOrganizationHeadDeputyResponse](#service-membership-v1-AssignOrganizationHeadDeputyResponse)
    - [AssignOrganizationHeadRequest](#service-membership-v1-AssignOrganizationHeadRequest)
    - [AssignOrganizationHeadResponse](#service-membership-v1-AssignOrganizationHeadResponse)
    - [CancelScheduledVacationRequest](#service-membership-v1-CancelScheduledVacationRequest)
    - [CancelScheduledVacationResponse](#service-membership-v1-CancelScheduledVacationResponse)
    - [ForceEndVacationRequest](#service-membership-v1-ForceEndVacationRequest)
    - [ForceEndVacationResponse](#service-membership-v1-ForceEndVacationResponse)
    - [GrantSystemAdminRequest](#service-membership-v1-GrantSystemAdminRequest)
    - [GrantSystemAdminResponse](#service-membership-v1-GrantSystemAdminResponse)
    - [HireEmployeeRequest](#service-membership-v1-HireEmployeeRequest)
    - [HireEmployeeResponse](#service-membership-v1-HireEmployeeResponse)
    - [RemoveClinicHeadDeputyRequest](#service-membership-v1-RemoveClinicHeadDeputyRequest)
    - [RemoveClinicHeadDeputyResponse](#service-membership-v1-RemoveClinicHeadDeputyResponse)
    - [RemoveDepartmentResponsibleDeputyRequest](#service-membership-v1-RemoveDepartmentResponsibleDeputyRequest)
    - [RemoveDepartmentResponsibleDeputyResponse](#service-membership-v1-RemoveDepartmentResponsibleDeputyResponse)
    - [RemoveOrganizationAdminDeputyRequest](#service-membership-v1-RemoveOrganizationAdminDeputyRequest)
    - [RemoveOrganizationAdminDeputyResponse](#service-membership-v1-RemoveOrganizationAdminDeputyResponse)
    - [RemoveOrganizationDispatcherDeputyRequest](#service-membership-v1-RemoveOrganizationDispatcherDeputyRequest)
    - [RemoveOrganizationDispatcherDeputyResponse](#service-membership-v1-RemoveOrganizationDispatcherDeputyResponse)
    - [RemoveOrganizationHeadDeputyRequest](#service-membership-v1-RemoveOrganizationHeadDeputyRequest)
    - [RemoveOrganizationHeadDeputyResponse](#service-membership-v1-RemoveOrganizationHeadDeputyResponse)
    - [RevokeClinicHeadRequest](#service-membership-v1-RevokeClinicHeadRequest)
    - [RevokeClinicHeadResponse](#service-membership-v1-RevokeClinicHeadResponse)
    - [RevokeDepartmentResponsibleRequest](#service-membership-v1-RevokeDepartmentResponsibleRequest)
    - [RevokeDepartmentResponsibleResponse](#service-membership-v1-RevokeDepartmentResponsibleResponse)
    - [RevokeOrganizationAdminRequest](#service-membership-v1-RevokeOrganizationAdminRequest)
    - [RevokeOrganizationAdminResponse](#service-membership-v1-RevokeOrganizationAdminResponse)
    - [RevokeOrganizationDispatcherRequest](#service-membership-v1-RevokeOrganizationDispatcherRequest)
    - [RevokeOrganizationDispatcherResponse](#service-membership-v1-RevokeOrganizationDispatcherResponse)
    - [RevokeOrganizationHeadRequest](#service-membership-v1-RevokeOrganizationHeadRequest)
    - [RevokeOrganizationHeadResponse](#service-membership-v1-RevokeOrganizationHeadResponse)
    - [RevokeSystemAdminRequest](#service-membership-v1-RevokeSystemAdminRequest)
    - [RevokeSystemAdminResponse](#service-membership-v1-RevokeSystemAdminResponse)
    - [ScheduleVacationRequest](#service-membership-v1-ScheduleVacationRequest)
    - [ScheduleVacationResponse](#service-membership-v1-ScheduleVacationResponse)
    - [StartVacationNowRequest](#service-membership-v1-StartVacationNowRequest)
    - [StartVacationNowResponse](#service-membership-v1-StartVacationNowResponse)
    - [TerminateEmployeeRequest](#service-membership-v1-TerminateEmployeeRequest)
    - [TerminateEmployeeResponse](#service-membership-v1-TerminateEmployeeResponse)
    - [UpdateEmployeeDepartmentRequest](#service-membership-v1-UpdateEmployeeDepartmentRequest)
    - [UpdateEmployeeDepartmentResponse](#service-membership-v1-UpdateEmployeeDepartmentResponse)
    - [UpdateEmployeePositionRequest](#service-membership-v1-UpdateEmployeePositionRequest)
    - [UpdateEmployeePositionResponse](#service-membership-v1-UpdateEmployeePositionResponse)
    - [UpdateVacationEndDateRequest](#service-membership-v1-UpdateVacationEndDateRequest)
    - [UpdateVacationEndDateResponse](#service-membership-v1-UpdateVacationEndDateResponse)

    - [MembershipService](#service-membership-v1-MembershipService)

- [service/orgstructure/v1/orgstructure.proto](#service_orgstructure_v1_orgstructure-proto)
    - [AddressInput](#service-orgstructure-v1-AddressInput)
    - [CreateClinicRequest](#service-orgstructure-v1-CreateClinicRequest)
    - [CreateClinicResponse](#service-orgstructure-v1-CreateClinicResponse)
    - [CreateDepartmentRequest](#service-orgstructure-v1-CreateDepartmentRequest)
    - [CreateDepartmentResponse](#service-orgstructure-v1-CreateDepartmentResponse)
    - [CreateOrganizationRequest](#service-orgstructure-v1-CreateOrganizationRequest)
    - [CreateOrganizationResponse](#service-orgstructure-v1-CreateOrganizationResponse)
    - [PointInput](#service-orgstructure-v1-PointInput)
    - [UpdateClinicDetailsRequest](#service-orgstructure-v1-UpdateClinicDetailsRequest)
    - [UpdateClinicDetailsResponse](#service-orgstructure-v1-UpdateClinicDetailsResponse)
    - [UpdateClinicPhysicalAddressRequest](#service-orgstructure-v1-UpdateClinicPhysicalAddressRequest)
    - [UpdateClinicPhysicalAddressResponse](#service-orgstructure-v1-UpdateClinicPhysicalAddressResponse)
    - [UpdateDepartmentDetailsRequest](#service-orgstructure-v1-UpdateDepartmentDetailsRequest)
    - [UpdateDepartmentDetailsResponse](#service-orgstructure-v1-UpdateDepartmentDetailsResponse)
    - [UpdateOrganizationDetailsRequest](#service-orgstructure-v1-UpdateOrganizationDetailsRequest)
    - [UpdateOrganizationDetailsResponse](#service-orgstructure-v1-UpdateOrganizationDetailsResponse)
    - [UpdateOrganizationLegalAddressRequest](#service-orgstructure-v1-UpdateOrganizationLegalAddressRequest)
    - [UpdateOrganizationLegalAddressResponse](#service-orgstructure-v1-UpdateOrganizationLegalAddressResponse)

    - [OrgStructureService](#service-orgstructure-v1-OrgStructureService)

- [Scalar Value Types](#scalar-value-types)



<a name="event_clinic_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/clinic/v1/events.proto



<a name="event-clinic-v1-Address"></a>

### Address



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  |  |
| point | [Point](#event-clinic-v1-Point) | optional |  |






<a name="event-clinic-v1-ClinicCreated"></a>

### ClinicCreated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| physical_address | [Address](#event-clinic-v1-Address) |  |  |






<a name="event-clinic-v1-ClinicDetailsChanged"></a>

### ClinicDetailsChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="event-clinic-v1-ClinicHeadAssigned"></a>

### ClinicHeadAssigned
ClinicHeadAssigned — an employee became the head of this clinic.
Payload carries the holder; aggregate_id in the envelope is
clinic_id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-clinic-v1-ClinicHeadDeputyAssigned"></a>

### ClinicHeadDeputyAssigned
ClinicHeadDeputyAssigned — a deputy was attached. The deputy must
currently work in any department of this clinic.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="event-clinic-v1-ClinicHeadDeputyRemoved"></a>

### ClinicHeadDeputyRemoved
ClinicHeadDeputyRemoved — the deputy slot was cleared.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-clinic-v1-ClinicHeadRevoked"></a>

### ClinicHeadRevoked
ClinicHeadRevoked — role revoked (explicit or cascade from
cross-clinic transfer / termination).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-clinic-v1-ClinicPhysicalAddressChanged"></a>

### ClinicPhysicalAddressChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| physical_address | [Address](#event-clinic-v1-Address) |  |  |






<a name="event-clinic-v1-Point"></a>

### Point



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| longitude | [double](#double) |  |  |
| latitude | [double](#double) |  |  |















<a name="event_department_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/department/v1/events.proto



<a name="event-department-v1-DepartmentCreated"></a>

### DepartmentCreated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="event-department-v1-DepartmentDetailsChanged"></a>

### DepartmentDetailsChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="event-department-v1-DepartmentResponsibleAssigned"></a>

### DepartmentResponsibleAssigned
DepartmentResponsibleAssigned — an employee became the &#34;responsible&#34;
of this department. Payload carries the holder; aggregate_id in the
envelope is department_id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-department-v1-DepartmentResponsibleDeputyAssigned"></a>

### DepartmentResponsibleDeputyAssigned
DepartmentResponsibleDeputyAssigned — a deputy was attached to the
role held by employee_id. The deputy must currently work in the
same department.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="event-department-v1-DepartmentResponsibleDeputyRemoved"></a>

### DepartmentResponsibleDeputyRemoved
DepartmentResponsibleDeputyRemoved — the deputy slot was cleared
(explicit removal, cascade from role revocation, or cascade from
deputy employee&#39;s termination or department transfer).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-department-v1-DepartmentResponsibleRevoked"></a>

### DepartmentResponsibleRevoked
DepartmentResponsibleRevoked — role revoked (either explicitly or
as a cascade from transfer/termination).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |















<a name="event_employee_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/employee/v1/events.proto



<a name="event-employee-v1-EmployeeDepartmentChanged"></a>

### EmployeeDepartmentChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |






<a name="event-employee-v1-EmployeeHired"></a>

### EmployeeHired



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| position | [string](#string) | optional |  |






<a name="event-employee-v1-EmployeePositionChanged"></a>

### EmployeePositionChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| position | [string](#string) | optional |  |






<a name="event-employee-v1-EmployeeTerminated"></a>

### EmployeeTerminated







<a name="event-employee-v1-VacationCancelled"></a>

### VacationCancelled



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="event-employee-v1-VacationEndDateChanged"></a>

### VacationEndDateChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-employee-v1-VacationEnded"></a>

### VacationEnded



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-employee-v1-VacationScheduled"></a>

### VacationScheduled



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| starts_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |






<a name="event-employee-v1-VacationStarted"></a>

### VacationStarted



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| starts_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |















<a name="event_incident_category_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/incident/category/v1/events.proto



<a name="event-incident-category-v1-IncidentCategoryCreated"></a>

### IncidentCategoryCreated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| parent_category_id | [string](#string) | optional | unset = created at root |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="event-incident-category-v1-IncidentCategoryDeactivated"></a>

### IncidentCategoryDeactivated







<a name="event-incident-category-v1-IncidentCategoryDeleted"></a>

### IncidentCategoryDeleted







<a name="event-incident-category-v1-IncidentCategoryDetailsChanged"></a>

### IncidentCategoryDetailsChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="event-incident-category-v1-IncidentCategoryMoved"></a>

### IncidentCategoryMoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| new_parent_category_id | [string](#string) | optional | unset = moved to root |






<a name="event-incident-category-v1-IncidentCategoryReactivated"></a>

### IncidentCategoryReactivated
















<a name="event_incident_type_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/incident/type/v1/events.proto



<a name="event-incident-type-v1-IncidentTypeCreated"></a>

### IncidentTypeCreated



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| category_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="event-incident-type-v1-IncidentTypeDeactivated"></a>

### IncidentTypeDeactivated







<a name="event-incident-type-v1-IncidentTypeDeleted"></a>

### IncidentTypeDeleted







<a name="event-incident-type-v1-IncidentTypeDetailsChanged"></a>

### IncidentTypeDetailsChanged



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="event-incident-type-v1-IncidentTypeMoved"></a>

### IncidentTypeMoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| new_category_id | [string](#string) |  |  |






<a name="event-incident-type-v1-IncidentTypeReactivated"></a>

### IncidentTypeReactivated
















<a name="event_organization_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/organization/v1/events.proto



<a name="event-organization-v1-Address"></a>

### Address



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  |  |
| point | [Point](#event-organization-v1-Point) | optional |  |






<a name="event-organization-v1-OrganizationAdminAssigned"></a>

### OrganizationAdminAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationAdminDeputyAssigned"></a>

### OrganizationAdminDeputyAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationAdminDeputyRemoved"></a>

### OrganizationAdminDeputyRemoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationAdminRevoked"></a>

### OrganizationAdminRevoked



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationCreated"></a>

### OrganizationCreated
OrganizationCreated — the aggregate entered this state on creation.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| legal_address | [Address](#event-organization-v1-Address) |  |  |






<a name="event-organization-v1-OrganizationDetailsChanged"></a>

### OrganizationDetailsChanged
OrganizationDetailsChanged — name and description became these.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="event-organization-v1-OrganizationDispatcherAssigned"></a>

### OrganizationDispatcherAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationDispatcherDeputyAssigned"></a>

### OrganizationDispatcherDeputyAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationDispatcherDeputyRemoved"></a>

### OrganizationDispatcherDeputyRemoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationDispatcherRevoked"></a>

### OrganizationDispatcherRevoked



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationHeadAssigned"></a>

### OrganizationHeadAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationHeadDeputyAssigned"></a>

### OrganizationHeadDeputyAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationHeadDeputyRemoved"></a>

### OrganizationHeadDeputyRemoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationHeadRevoked"></a>

### OrganizationHeadRevoked



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationLegalAddressChanged"></a>

### OrganizationLegalAddressChanged
OrganizationLegalAddressChanged — legal address became this.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| legal_address | [Address](#event-organization-v1-Address) |  |  |






<a name="event-organization-v1-Point"></a>

### Point



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| longitude | [double](#double) |  |  |
| latitude | [double](#double) |  |  |















<a name="event_system_admin_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/system_admin/v1/events.proto



<a name="event-system_admin-v1-SystemAdminGranted"></a>

### SystemAdminGranted







<a name="event-system_admin-v1-SystemAdminRevoked"></a>

### SystemAdminRevoked
















<a name="event_v1_envelope-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/v1/envelope.proto



<a name="event-v1-Envelope"></a>

### Envelope
Envelope is the transport wrapper for every medincident domain event
published on the bus. The publisher service (a separate drainer) reads
outbox rows, unmarshals the envelope, and ships it to NATS JetStream.
Consumers project events by inspecting aggregate_type and unpacking
payload via its Any type URL.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| occurred_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| aggregate_type | [string](#string) |  | Aggregate type the event belongs to — a short lowercase token identifying the aggregate root (e.g. &#34;organization&#34;, &#34;clinic&#34;, &#34;department&#34;, &#34;employee&#34;, &#34;system_admin&#34;, &#34;incident_category&#34;, &#34;incident_type&#34;). Consumers use this to pick a projection without unpacking payload. New aggregates are added as new values; this field is intentionally an open string rather than an enum so adding an aggregate does not require a proto revision. |
| aggregate_id | [string](#string) |  | Aggregate id in text form. For most aggregates this is a UUID string; for system_admin it is the Zitadel user id (opaque). |
| payload | [google.protobuf.Any](#google-protobuf-Any) |  |  |















<a name="service_incident_classifier_v1_incident_classifier-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## service/incident/classifier/v1/incident_classifier.proto



<a name="service-incident-classifier-v1-CreateIncidentCategoryRequest"></a>

### CreateIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| parent_category_id | [string](#string) | optional |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="service-incident-classifier-v1-CreateIncidentCategoryResponse"></a>

### CreateIncidentCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-CreateIncidentTypeRequest"></a>

### CreateIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="service-incident-classifier-v1-CreateIncidentTypeResponse"></a>

### CreateIncidentTypeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-DeactivateIncidentCategoryRequest"></a>

### DeactivateIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-DeactivateIncidentCategoryResponse"></a>

### DeactivateIncidentCategoryResponse







<a name="service-incident-classifier-v1-DeactivateIncidentTypeRequest"></a>

### DeactivateIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-DeactivateIncidentTypeResponse"></a>

### DeactivateIncidentTypeResponse







<a name="service-incident-classifier-v1-DeleteIncidentCategoryRequest"></a>

### DeleteIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-DeleteIncidentCategoryResponse"></a>

### DeleteIncidentCategoryResponse







<a name="service-incident-classifier-v1-DeleteIncidentTypeRequest"></a>

### DeleteIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-DeleteIncidentTypeResponse"></a>

### DeleteIncidentTypeResponse







<a name="service-incident-classifier-v1-MoveIncidentCategoryRequest"></a>

### MoveIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| new_parent_category_id | [string](#string) | optional |  |






<a name="service-incident-classifier-v1-MoveIncidentCategoryResponse"></a>

### MoveIncidentCategoryResponse







<a name="service-incident-classifier-v1-MoveIncidentTypeRequest"></a>

### MoveIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| new_category_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-MoveIncidentTypeResponse"></a>

### MoveIncidentTypeResponse







<a name="service-incident-classifier-v1-ReactivateIncidentCategoryRequest"></a>

### ReactivateIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-ReactivateIncidentCategoryResponse"></a>

### ReactivateIncidentCategoryResponse







<a name="service-incident-classifier-v1-ReactivateIncidentTypeRequest"></a>

### ReactivateIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="service-incident-classifier-v1-ReactivateIncidentTypeResponse"></a>

### ReactivateIncidentTypeResponse







<a name="service-incident-classifier-v1-UpdateIncidentCategoryDetailsRequest"></a>

### UpdateIncidentCategoryDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="service-incident-classifier-v1-UpdateIncidentCategoryDetailsResponse"></a>

### UpdateIncidentCategoryDetailsResponse







<a name="service-incident-classifier-v1-UpdateIncidentTypeDetailsRequest"></a>

### UpdateIncidentTypeDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="service-incident-classifier-v1-UpdateIncidentTypeDetailsResponse"></a>

### UpdateIncidentTypeDetailsResponse













<a name="service-incident-classifier-v1-IncidentClassifierService"></a>

### IncidentClassifierService
IncidentClassifierService is the command-side contract for a single
administrative surface that manages an organisation&#39;s incident
classifier: a tree of incident categories (max depth 5) with incident
types living inside any category. Every mutation returns an id (for
Create) or an empty response and publishes one domain event per
touched row.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateIncidentCategory | [CreateIncidentCategoryRequest](#service-incident-classifier-v1-CreateIncidentCategoryRequest) | [CreateIncidentCategoryResponse](#service-incident-classifier-v1-CreateIncidentCategoryResponse) | --- Categories --- |
| UpdateIncidentCategoryDetails | [UpdateIncidentCategoryDetailsRequest](#service-incident-classifier-v1-UpdateIncidentCategoryDetailsRequest) | [UpdateIncidentCategoryDetailsResponse](#service-incident-classifier-v1-UpdateIncidentCategoryDetailsResponse) |  |
| MoveIncidentCategory | [MoveIncidentCategoryRequest](#service-incident-classifier-v1-MoveIncidentCategoryRequest) | [MoveIncidentCategoryResponse](#service-incident-classifier-v1-MoveIncidentCategoryResponse) |  |
| DeactivateIncidentCategory | [DeactivateIncidentCategoryRequest](#service-incident-classifier-v1-DeactivateIncidentCategoryRequest) | [DeactivateIncidentCategoryResponse](#service-incident-classifier-v1-DeactivateIncidentCategoryResponse) |  |
| ReactivateIncidentCategory | [ReactivateIncidentCategoryRequest](#service-incident-classifier-v1-ReactivateIncidentCategoryRequest) | [ReactivateIncidentCategoryResponse](#service-incident-classifier-v1-ReactivateIncidentCategoryResponse) |  |
| DeleteIncidentCategory | [DeleteIncidentCategoryRequest](#service-incident-classifier-v1-DeleteIncidentCategoryRequest) | [DeleteIncidentCategoryResponse](#service-incident-classifier-v1-DeleteIncidentCategoryResponse) |  |
| CreateIncidentType | [CreateIncidentTypeRequest](#service-incident-classifier-v1-CreateIncidentTypeRequest) | [CreateIncidentTypeResponse](#service-incident-classifier-v1-CreateIncidentTypeResponse) | --- Types --- |
| UpdateIncidentTypeDetails | [UpdateIncidentTypeDetailsRequest](#service-incident-classifier-v1-UpdateIncidentTypeDetailsRequest) | [UpdateIncidentTypeDetailsResponse](#service-incident-classifier-v1-UpdateIncidentTypeDetailsResponse) |  |
| MoveIncidentType | [MoveIncidentTypeRequest](#service-incident-classifier-v1-MoveIncidentTypeRequest) | [MoveIncidentTypeResponse](#service-incident-classifier-v1-MoveIncidentTypeResponse) |  |
| DeactivateIncidentType | [DeactivateIncidentTypeRequest](#service-incident-classifier-v1-DeactivateIncidentTypeRequest) | [DeactivateIncidentTypeResponse](#service-incident-classifier-v1-DeactivateIncidentTypeResponse) |  |
| ReactivateIncidentType | [ReactivateIncidentTypeRequest](#service-incident-classifier-v1-ReactivateIncidentTypeRequest) | [ReactivateIncidentTypeResponse](#service-incident-classifier-v1-ReactivateIncidentTypeResponse) |  |
| DeleteIncidentType | [DeleteIncidentTypeRequest](#service-incident-classifier-v1-DeleteIncidentTypeRequest) | [DeleteIncidentTypeResponse](#service-incident-classifier-v1-DeleteIncidentTypeResponse) |  |





<a name="service_membership_v1_membership-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## service/membership/v1/membership.proto



<a name="service-membership-v1-AssignClinicHeadDeputyRequest"></a>

### AssignClinicHeadDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignClinicHeadDeputyResponse"></a>

### AssignClinicHeadDeputyResponse







<a name="service-membership-v1-AssignClinicHeadRequest"></a>

### AssignClinicHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignClinicHeadResponse"></a>

### AssignClinicHeadResponse







<a name="service-membership-v1-AssignDepartmentResponsibleDeputyRequest"></a>

### AssignDepartmentResponsibleDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignDepartmentResponsibleDeputyResponse"></a>

### AssignDepartmentResponsibleDeputyResponse







<a name="service-membership-v1-AssignDepartmentResponsibleRequest"></a>

### AssignDepartmentResponsibleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignDepartmentResponsibleResponse"></a>

### AssignDepartmentResponsibleResponse







<a name="service-membership-v1-AssignOrganizationAdminDeputyRequest"></a>

### AssignOrganizationAdminDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignOrganizationAdminDeputyResponse"></a>

### AssignOrganizationAdminDeputyResponse







<a name="service-membership-v1-AssignOrganizationAdminRequest"></a>

### AssignOrganizationAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignOrganizationAdminResponse"></a>

### AssignOrganizationAdminResponse







<a name="service-membership-v1-AssignOrganizationDispatcherDeputyRequest"></a>

### AssignOrganizationDispatcherDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignOrganizationDispatcherDeputyResponse"></a>

### AssignOrganizationDispatcherDeputyResponse







<a name="service-membership-v1-AssignOrganizationDispatcherRequest"></a>

### AssignOrganizationDispatcherRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignOrganizationDispatcherResponse"></a>

### AssignOrganizationDispatcherResponse







<a name="service-membership-v1-AssignOrganizationHeadDeputyRequest"></a>

### AssignOrganizationHeadDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignOrganizationHeadDeputyResponse"></a>

### AssignOrganizationHeadDeputyResponse







<a name="service-membership-v1-AssignOrganizationHeadRequest"></a>

### AssignOrganizationHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-AssignOrganizationHeadResponse"></a>

### AssignOrganizationHeadResponse







<a name="service-membership-v1-CancelScheduledVacationRequest"></a>

### CancelScheduledVacationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="service-membership-v1-CancelScheduledVacationResponse"></a>

### CancelScheduledVacationResponse







<a name="service-membership-v1-ForceEndVacationRequest"></a>

### ForceEndVacationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="service-membership-v1-ForceEndVacationResponse"></a>

### ForceEndVacationResponse







<a name="service-membership-v1-GrantSystemAdminRequest"></a>

### GrantSystemAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |






<a name="service-membership-v1-GrantSystemAdminResponse"></a>

### GrantSystemAdminResponse







<a name="service-membership-v1-HireEmployeeRequest"></a>

### HireEmployeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| position | [string](#string) | optional |  |






<a name="service-membership-v1-HireEmployeeResponse"></a>

### HireEmployeeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RemoveClinicHeadDeputyRequest"></a>

### RemoveClinicHeadDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RemoveClinicHeadDeputyResponse"></a>

### RemoveClinicHeadDeputyResponse







<a name="service-membership-v1-RemoveDepartmentResponsibleDeputyRequest"></a>

### RemoveDepartmentResponsibleDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RemoveDepartmentResponsibleDeputyResponse"></a>

### RemoveDepartmentResponsibleDeputyResponse







<a name="service-membership-v1-RemoveOrganizationAdminDeputyRequest"></a>

### RemoveOrganizationAdminDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RemoveOrganizationAdminDeputyResponse"></a>

### RemoveOrganizationAdminDeputyResponse







<a name="service-membership-v1-RemoveOrganizationDispatcherDeputyRequest"></a>

### RemoveOrganizationDispatcherDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RemoveOrganizationDispatcherDeputyResponse"></a>

### RemoveOrganizationDispatcherDeputyResponse







<a name="service-membership-v1-RemoveOrganizationHeadDeputyRequest"></a>

### RemoveOrganizationHeadDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RemoveOrganizationHeadDeputyResponse"></a>

### RemoveOrganizationHeadDeputyResponse







<a name="service-membership-v1-RevokeClinicHeadRequest"></a>

### RevokeClinicHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RevokeClinicHeadResponse"></a>

### RevokeClinicHeadResponse







<a name="service-membership-v1-RevokeDepartmentResponsibleRequest"></a>

### RevokeDepartmentResponsibleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RevokeDepartmentResponsibleResponse"></a>

### RevokeDepartmentResponsibleResponse







<a name="service-membership-v1-RevokeOrganizationAdminRequest"></a>

### RevokeOrganizationAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RevokeOrganizationAdminResponse"></a>

### RevokeOrganizationAdminResponse







<a name="service-membership-v1-RevokeOrganizationDispatcherRequest"></a>

### RevokeOrganizationDispatcherRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RevokeOrganizationDispatcherResponse"></a>

### RevokeOrganizationDispatcherResponse







<a name="service-membership-v1-RevokeOrganizationHeadRequest"></a>

### RevokeOrganizationHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-RevokeOrganizationHeadResponse"></a>

### RevokeOrganizationHeadResponse







<a name="service-membership-v1-RevokeSystemAdminRequest"></a>

### RevokeSystemAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |






<a name="service-membership-v1-RevokeSystemAdminResponse"></a>

### RevokeSystemAdminResponse







<a name="service-membership-v1-ScheduleVacationRequest"></a>

### ScheduleVacationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| starts_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |






<a name="service-membership-v1-ScheduleVacationResponse"></a>

### ScheduleVacationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="service-membership-v1-StartVacationNowRequest"></a>

### StartVacationNowRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |






<a name="service-membership-v1-StartVacationNowResponse"></a>

### StartVacationNowResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="service-membership-v1-TerminateEmployeeRequest"></a>

### TerminateEmployeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="service-membership-v1-TerminateEmployeeResponse"></a>

### TerminateEmployeeResponse







<a name="service-membership-v1-UpdateEmployeeDepartmentRequest"></a>

### UpdateEmployeeDepartmentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |






<a name="service-membership-v1-UpdateEmployeeDepartmentResponse"></a>

### UpdateEmployeeDepartmentResponse







<a name="service-membership-v1-UpdateEmployeePositionRequest"></a>

### UpdateEmployeePositionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| position | [string](#string) | optional |  |






<a name="service-membership-v1-UpdateEmployeePositionResponse"></a>

### UpdateEmployeePositionResponse







<a name="service-membership-v1-UpdateVacationEndDateRequest"></a>

### UpdateVacationEndDateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="service-membership-v1-UpdateVacationEndDateResponse"></a>

### UpdateVacationEndDateResponse













<a name="service-membership-v1-MembershipService"></a>

### MembershipService
MembershipService is the command-side contract for who-works-where
data. It covers the Employee lifecycle (hire / update / terminate),
Vacation lifecycle (start / schedule / end / cancel / change),
role grants with deputy management for DepartmentResponsible,
ClinicHead, OrganizationAdmin, OrganizationHead, and
OrganizationDispatcher, plus the global SystemAdmin grant that
operates directly on Zitadel user identifiers.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| HireEmployee | [HireEmployeeRequest](#service-membership-v1-HireEmployeeRequest) | [HireEmployeeResponse](#service-membership-v1-HireEmployeeResponse) | Employee lifecycle |
| UpdateEmployeePosition | [UpdateEmployeePositionRequest](#service-membership-v1-UpdateEmployeePositionRequest) | [UpdateEmployeePositionResponse](#service-membership-v1-UpdateEmployeePositionResponse) |  |
| UpdateEmployeeDepartment | [UpdateEmployeeDepartmentRequest](#service-membership-v1-UpdateEmployeeDepartmentRequest) | [UpdateEmployeeDepartmentResponse](#service-membership-v1-UpdateEmployeeDepartmentResponse) |  |
| TerminateEmployee | [TerminateEmployeeRequest](#service-membership-v1-TerminateEmployeeRequest) | [TerminateEmployeeResponse](#service-membership-v1-TerminateEmployeeResponse) |  |
| StartVacationNow | [StartVacationNowRequest](#service-membership-v1-StartVacationNowRequest) | [StartVacationNowResponse](#service-membership-v1-StartVacationNowResponse) | Vacation lifecycle |
| ScheduleVacation | [ScheduleVacationRequest](#service-membership-v1-ScheduleVacationRequest) | [ScheduleVacationResponse](#service-membership-v1-ScheduleVacationResponse) |  |
| UpdateVacationEndDate | [UpdateVacationEndDateRequest](#service-membership-v1-UpdateVacationEndDateRequest) | [UpdateVacationEndDateResponse](#service-membership-v1-UpdateVacationEndDateResponse) |  |
| ForceEndVacation | [ForceEndVacationRequest](#service-membership-v1-ForceEndVacationRequest) | [ForceEndVacationResponse](#service-membership-v1-ForceEndVacationResponse) |  |
| CancelScheduledVacation | [CancelScheduledVacationRequest](#service-membership-v1-CancelScheduledVacationRequest) | [CancelScheduledVacationResponse](#service-membership-v1-CancelScheduledVacationResponse) |  |
| AssignDepartmentResponsible | [AssignDepartmentResponsibleRequest](#service-membership-v1-AssignDepartmentResponsibleRequest) | [AssignDepartmentResponsibleResponse](#service-membership-v1-AssignDepartmentResponsibleResponse) | ------ DepartmentResponsible ------ |
| RevokeDepartmentResponsible | [RevokeDepartmentResponsibleRequest](#service-membership-v1-RevokeDepartmentResponsibleRequest) | [RevokeDepartmentResponsibleResponse](#service-membership-v1-RevokeDepartmentResponsibleResponse) |  |
| AssignDepartmentResponsibleDeputy | [AssignDepartmentResponsibleDeputyRequest](#service-membership-v1-AssignDepartmentResponsibleDeputyRequest) | [AssignDepartmentResponsibleDeputyResponse](#service-membership-v1-AssignDepartmentResponsibleDeputyResponse) |  |
| RemoveDepartmentResponsibleDeputy | [RemoveDepartmentResponsibleDeputyRequest](#service-membership-v1-RemoveDepartmentResponsibleDeputyRequest) | [RemoveDepartmentResponsibleDeputyResponse](#service-membership-v1-RemoveDepartmentResponsibleDeputyResponse) |  |
| AssignClinicHead | [AssignClinicHeadRequest](#service-membership-v1-AssignClinicHeadRequest) | [AssignClinicHeadResponse](#service-membership-v1-AssignClinicHeadResponse) | ------ ClinicHead ------ |
| RevokeClinicHead | [RevokeClinicHeadRequest](#service-membership-v1-RevokeClinicHeadRequest) | [RevokeClinicHeadResponse](#service-membership-v1-RevokeClinicHeadResponse) |  |
| AssignClinicHeadDeputy | [AssignClinicHeadDeputyRequest](#service-membership-v1-AssignClinicHeadDeputyRequest) | [AssignClinicHeadDeputyResponse](#service-membership-v1-AssignClinicHeadDeputyResponse) |  |
| RemoveClinicHeadDeputy | [RemoveClinicHeadDeputyRequest](#service-membership-v1-RemoveClinicHeadDeputyRequest) | [RemoveClinicHeadDeputyResponse](#service-membership-v1-RemoveClinicHeadDeputyResponse) |  |
| AssignOrganizationAdmin | [AssignOrganizationAdminRequest](#service-membership-v1-AssignOrganizationAdminRequest) | [AssignOrganizationAdminResponse](#service-membership-v1-AssignOrganizationAdminResponse) | ------ OrganizationAdmin ------ |
| RevokeOrganizationAdmin | [RevokeOrganizationAdminRequest](#service-membership-v1-RevokeOrganizationAdminRequest) | [RevokeOrganizationAdminResponse](#service-membership-v1-RevokeOrganizationAdminResponse) |  |
| AssignOrganizationAdminDeputy | [AssignOrganizationAdminDeputyRequest](#service-membership-v1-AssignOrganizationAdminDeputyRequest) | [AssignOrganizationAdminDeputyResponse](#service-membership-v1-AssignOrganizationAdminDeputyResponse) |  |
| RemoveOrganizationAdminDeputy | [RemoveOrganizationAdminDeputyRequest](#service-membership-v1-RemoveOrganizationAdminDeputyRequest) | [RemoveOrganizationAdminDeputyResponse](#service-membership-v1-RemoveOrganizationAdminDeputyResponse) |  |
| AssignOrganizationHead | [AssignOrganizationHeadRequest](#service-membership-v1-AssignOrganizationHeadRequest) | [AssignOrganizationHeadResponse](#service-membership-v1-AssignOrganizationHeadResponse) | ------ OrganizationHead ------ |
| RevokeOrganizationHead | [RevokeOrganizationHeadRequest](#service-membership-v1-RevokeOrganizationHeadRequest) | [RevokeOrganizationHeadResponse](#service-membership-v1-RevokeOrganizationHeadResponse) |  |
| AssignOrganizationHeadDeputy | [AssignOrganizationHeadDeputyRequest](#service-membership-v1-AssignOrganizationHeadDeputyRequest) | [AssignOrganizationHeadDeputyResponse](#service-membership-v1-AssignOrganizationHeadDeputyResponse) |  |
| RemoveOrganizationHeadDeputy | [RemoveOrganizationHeadDeputyRequest](#service-membership-v1-RemoveOrganizationHeadDeputyRequest) | [RemoveOrganizationHeadDeputyResponse](#service-membership-v1-RemoveOrganizationHeadDeputyResponse) |  |
| AssignOrganizationDispatcher | [AssignOrganizationDispatcherRequest](#service-membership-v1-AssignOrganizationDispatcherRequest) | [AssignOrganizationDispatcherResponse](#service-membership-v1-AssignOrganizationDispatcherResponse) | ------ OrganizationDispatcher ------ |
| RevokeOrganizationDispatcher | [RevokeOrganizationDispatcherRequest](#service-membership-v1-RevokeOrganizationDispatcherRequest) | [RevokeOrganizationDispatcherResponse](#service-membership-v1-RevokeOrganizationDispatcherResponse) |  |
| AssignOrganizationDispatcherDeputy | [AssignOrganizationDispatcherDeputyRequest](#service-membership-v1-AssignOrganizationDispatcherDeputyRequest) | [AssignOrganizationDispatcherDeputyResponse](#service-membership-v1-AssignOrganizationDispatcherDeputyResponse) |  |
| RemoveOrganizationDispatcherDeputy | [RemoveOrganizationDispatcherDeputyRequest](#service-membership-v1-RemoveOrganizationDispatcherDeputyRequest) | [RemoveOrganizationDispatcherDeputyResponse](#service-membership-v1-RemoveOrganizationDispatcherDeputyResponse) |  |
| GrantSystemAdmin | [GrantSystemAdminRequest](#service-membership-v1-GrantSystemAdminRequest) | [GrantSystemAdminResponse](#service-membership-v1-GrantSystemAdminResponse) | ------ SystemAdmin ------ |
| RevokeSystemAdmin | [RevokeSystemAdminRequest](#service-membership-v1-RevokeSystemAdminRequest) | [RevokeSystemAdminResponse](#service-membership-v1-RevokeSystemAdminResponse) |  |





<a name="service_orgstructure_v1_orgstructure-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## service/orgstructure/v1/orgstructure.proto



<a name="service-orgstructure-v1-AddressInput"></a>

### AddressInput



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  |  |
| point | [PointInput](#service-orgstructure-v1-PointInput) | optional |  |






<a name="service-orgstructure-v1-CreateClinicRequest"></a>

### CreateClinicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| physical_address | [AddressInput](#service-orgstructure-v1-AddressInput) |  |  |
| description | [string](#string) | optional |  |






<a name="service-orgstructure-v1-CreateClinicResponse"></a>

### CreateClinicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |






<a name="service-orgstructure-v1-CreateDepartmentRequest"></a>

### CreateDepartmentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="service-orgstructure-v1-CreateDepartmentResponse"></a>

### CreateDepartmentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |






<a name="service-orgstructure-v1-CreateOrganizationRequest"></a>

### CreateOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| legal_address | [AddressInput](#service-orgstructure-v1-AddressInput) |  |  |
| description | [string](#string) | optional |  |






<a name="service-orgstructure-v1-CreateOrganizationResponse"></a>

### CreateOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |






<a name="service-orgstructure-v1-PointInput"></a>

### PointInput



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| longitude | [double](#double) |  |  |
| latitude | [double](#double) |  |  |






<a name="service-orgstructure-v1-UpdateClinicDetailsRequest"></a>

### UpdateClinicDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="service-orgstructure-v1-UpdateClinicDetailsResponse"></a>

### UpdateClinicDetailsResponse







<a name="service-orgstructure-v1-UpdateClinicPhysicalAddressRequest"></a>

### UpdateClinicPhysicalAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| physical_address | [AddressInput](#service-orgstructure-v1-AddressInput) |  |  |






<a name="service-orgstructure-v1-UpdateClinicPhysicalAddressResponse"></a>

### UpdateClinicPhysicalAddressResponse







<a name="service-orgstructure-v1-UpdateDepartmentDetailsRequest"></a>

### UpdateDepartmentDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="service-orgstructure-v1-UpdateDepartmentDetailsResponse"></a>

### UpdateDepartmentDetailsResponse







<a name="service-orgstructure-v1-UpdateOrganizationDetailsRequest"></a>

### UpdateOrganizationDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="service-orgstructure-v1-UpdateOrganizationDetailsResponse"></a>

### UpdateOrganizationDetailsResponse







<a name="service-orgstructure-v1-UpdateOrganizationLegalAddressRequest"></a>

### UpdateOrganizationLegalAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| legal_address | [AddressInput](#service-orgstructure-v1-AddressInput) |  |  |






<a name="service-orgstructure-v1-UpdateOrganizationLegalAddressResponse"></a>

### UpdateOrganizationLegalAddressResponse













<a name="service-orgstructure-v1-OrgStructureService"></a>

### OrgStructureService
OrgStructureService is the command-side contract for the three
organisational structure aggregates: Organization, Clinic, and
Department. Every mutation returns either an identifier (Create) or
an empty response (Update). google.api.http annotations drive a
separate REST gateway binary; command-service itself serves pure
gRPC, while this repo generates grpc-gateway stubs under pkg/ for
that gateway binary to consume.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateOrganization | [CreateOrganizationRequest](#service-orgstructure-v1-CreateOrganizationRequest) | [CreateOrganizationResponse](#service-orgstructure-v1-CreateOrganizationResponse) |  |
| UpdateOrganizationDetails | [UpdateOrganizationDetailsRequest](#service-orgstructure-v1-UpdateOrganizationDetailsRequest) | [UpdateOrganizationDetailsResponse](#service-orgstructure-v1-UpdateOrganizationDetailsResponse) |  |
| UpdateOrganizationLegalAddress | [UpdateOrganizationLegalAddressRequest](#service-orgstructure-v1-UpdateOrganizationLegalAddressRequest) | [UpdateOrganizationLegalAddressResponse](#service-orgstructure-v1-UpdateOrganizationLegalAddressResponse) |  |
| CreateClinic | [CreateClinicRequest](#service-orgstructure-v1-CreateClinicRequest) | [CreateClinicResponse](#service-orgstructure-v1-CreateClinicResponse) |  |
| UpdateClinicDetails | [UpdateClinicDetailsRequest](#service-orgstructure-v1-UpdateClinicDetailsRequest) | [UpdateClinicDetailsResponse](#service-orgstructure-v1-UpdateClinicDetailsResponse) |  |
| UpdateClinicPhysicalAddress | [UpdateClinicPhysicalAddressRequest](#service-orgstructure-v1-UpdateClinicPhysicalAddressRequest) | [UpdateClinicPhysicalAddressResponse](#service-orgstructure-v1-UpdateClinicPhysicalAddressResponse) |  |
| CreateDepartment | [CreateDepartmentRequest](#service-orgstructure-v1-CreateDepartmentRequest) | [CreateDepartmentResponse](#service-orgstructure-v1-CreateDepartmentResponse) |  |
| UpdateDepartmentDetails | [UpdateDepartmentDetailsRequest](#service-orgstructure-v1-UpdateDepartmentDetailsRequest) | [UpdateDepartmentDetailsResponse](#service-orgstructure-v1-UpdateDepartmentDetailsResponse) |  |





## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

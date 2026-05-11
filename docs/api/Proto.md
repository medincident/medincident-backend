# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [error/v1/error.proto](#error_v1_error-proto)
    - [ErrorCode](#error-v1-ErrorCode)
    - [ErrorResponse](#error-v1-ErrorResponse)
    - [ValidationFailedDetails](#error-v1-ValidationFailedDetails)
    - [ValidationFailedDetails.FieldViolation](#error-v1-ValidationFailedDetails-FieldViolation)

- [command/announcement/v1/announcement.proto](#command_announcement_v1_announcement-proto)
    - [ArchiveAnnouncementRequest](#command-announcement-v1-ArchiveAnnouncementRequest)
    - [ArchiveAnnouncementResponse](#command-announcement-v1-ArchiveAnnouncementResponse)
    - [CreateAnnouncementRequest](#command-announcement-v1-CreateAnnouncementRequest)
    - [CreateAnnouncementResponse](#command-announcement-v1-CreateAnnouncementResponse)
    - [UnarchiveAnnouncementRequest](#command-announcement-v1-UnarchiveAnnouncementRequest)
    - [UnarchiveAnnouncementResponse](#command-announcement-v1-UnarchiveAnnouncementResponse)
    - [UpdateAnnouncementPriorityRequest](#command-announcement-v1-UpdateAnnouncementPriorityRequest)
    - [UpdateAnnouncementPriorityResponse](#command-announcement-v1-UpdateAnnouncementPriorityResponse)
    - [UpdateAnnouncementRequest](#command-announcement-v1-UpdateAnnouncementRequest)
    - [UpdateAnnouncementResponse](#command-announcement-v1-UpdateAnnouncementResponse)

    - [AnnouncementPriority](#command-announcement-v1-AnnouncementPriority)

    - [AnnouncementCommandService](#command-announcement-v1-AnnouncementCommandService)

- [command/incident/buffer/v1/buffer.proto](#command_incident_buffer_v1_buffer-proto)
    - [CancelPatientIncidentRequest](#command-incident-buffer-v1-CancelPatientIncidentRequest)
    - [CancelPatientIncidentResponse](#command-incident-buffer-v1-CancelPatientIncidentResponse)
    - [PublishPatientIncidentRequest](#command-incident-buffer-v1-PublishPatientIncidentRequest)
    - [PublishPatientIncidentResponse](#command-incident-buffer-v1-PublishPatientIncidentResponse)
    - [RejectPatientIncidentRequest](#command-incident-buffer-v1-RejectPatientIncidentRequest)
    - [RejectPatientIncidentResponse](#command-incident-buffer-v1-RejectPatientIncidentResponse)
    - [SubmitPatientIncidentRequest](#command-incident-buffer-v1-SubmitPatientIncidentRequest)
    - [SubmitPatientIncidentResponse](#command-incident-buffer-v1-SubmitPatientIncidentResponse)
    - [UpdatePatientIncidentRequest](#command-incident-buffer-v1-UpdatePatientIncidentRequest)
    - [UpdatePatientIncidentResponse](#command-incident-buffer-v1-UpdatePatientIncidentResponse)

    - [IncidentBufferCommandService](#command-incident-buffer-v1-IncidentBufferCommandService)

- [command/incident/classifier/v1/incident_classifier.proto](#command_incident_classifier_v1_incident_classifier-proto)
    - [AllowIncidentTypeForPatientsRequest](#command-incident-classifier-v1-AllowIncidentTypeForPatientsRequest)
    - [AllowIncidentTypeForPatientsResponse](#command-incident-classifier-v1-AllowIncidentTypeForPatientsResponse)
    - [CreateIncidentCategoryRequest](#command-incident-classifier-v1-CreateIncidentCategoryRequest)
    - [CreateIncidentCategoryResponse](#command-incident-classifier-v1-CreateIncidentCategoryResponse)
    - [CreateIncidentTypeRequest](#command-incident-classifier-v1-CreateIncidentTypeRequest)
    - [CreateIncidentTypeResponse](#command-incident-classifier-v1-CreateIncidentTypeResponse)
    - [DeactivateIncidentCategoryRequest](#command-incident-classifier-v1-DeactivateIncidentCategoryRequest)
    - [DeactivateIncidentCategoryResponse](#command-incident-classifier-v1-DeactivateIncidentCategoryResponse)
    - [DeactivateIncidentTypeRequest](#command-incident-classifier-v1-DeactivateIncidentTypeRequest)
    - [DeactivateIncidentTypeResponse](#command-incident-classifier-v1-DeactivateIncidentTypeResponse)
    - [DeleteIncidentCategoryRequest](#command-incident-classifier-v1-DeleteIncidentCategoryRequest)
    - [DeleteIncidentCategoryResponse](#command-incident-classifier-v1-DeleteIncidentCategoryResponse)
    - [DeleteIncidentTypeRequest](#command-incident-classifier-v1-DeleteIncidentTypeRequest)
    - [DeleteIncidentTypeResponse](#command-incident-classifier-v1-DeleteIncidentTypeResponse)
    - [DisallowIncidentTypeForPatientsRequest](#command-incident-classifier-v1-DisallowIncidentTypeForPatientsRequest)
    - [DisallowIncidentTypeForPatientsResponse](#command-incident-classifier-v1-DisallowIncidentTypeForPatientsResponse)
    - [MoveIncidentCategoryRequest](#command-incident-classifier-v1-MoveIncidentCategoryRequest)
    - [MoveIncidentCategoryResponse](#command-incident-classifier-v1-MoveIncidentCategoryResponse)
    - [MoveIncidentTypeRequest](#command-incident-classifier-v1-MoveIncidentTypeRequest)
    - [MoveIncidentTypeResponse](#command-incident-classifier-v1-MoveIncidentTypeResponse)
    - [ReactivateIncidentCategoryRequest](#command-incident-classifier-v1-ReactivateIncidentCategoryRequest)
    - [ReactivateIncidentCategoryResponse](#command-incident-classifier-v1-ReactivateIncidentCategoryResponse)
    - [ReactivateIncidentTypeRequest](#command-incident-classifier-v1-ReactivateIncidentTypeRequest)
    - [ReactivateIncidentTypeResponse](#command-incident-classifier-v1-ReactivateIncidentTypeResponse)
    - [UpdateIncidentCategoryDetailsRequest](#command-incident-classifier-v1-UpdateIncidentCategoryDetailsRequest)
    - [UpdateIncidentCategoryDetailsResponse](#command-incident-classifier-v1-UpdateIncidentCategoryDetailsResponse)
    - [UpdateIncidentTypeDetailsRequest](#command-incident-classifier-v1-UpdateIncidentTypeDetailsRequest)
    - [UpdateIncidentTypeDetailsResponse](#command-incident-classifier-v1-UpdateIncidentTypeDetailsResponse)

    - [IncidentClassifierCommandService](#command-incident-classifier-v1-IncidentClassifierCommandService)

- [command/incident/v1/incident.proto](#command_incident_v1_incident-proto)
    - [CancelIncidentRequest](#command-incident-v1-CancelIncidentRequest)
    - [CancelIncidentResponse](#command-incident-v1-CancelIncidentResponse)
    - [CreateIncidentRequest](#command-incident-v1-CreateIncidentRequest)
    - [CreateIncidentResponse](#command-incident-v1-CreateIncidentResponse)
    - [ReopenIncidentRequest](#command-incident-v1-ReopenIncidentRequest)
    - [ReopenIncidentResponse](#command-incident-v1-ReopenIncidentResponse)
    - [UpdateIncidentDescriptionRequest](#command-incident-v1-UpdateIncidentDescriptionRequest)
    - [UpdateIncidentDescriptionResponse](#command-incident-v1-UpdateIncidentDescriptionResponse)
    - [UpdateIncidentPriorityRequest](#command-incident-v1-UpdateIncidentPriorityRequest)
    - [UpdateIncidentPriorityResponse](#command-incident-v1-UpdateIncidentPriorityResponse)
    - [UpdateIncidentStatusRequest](#command-incident-v1-UpdateIncidentStatusRequest)
    - [UpdateIncidentStatusResponse](#command-incident-v1-UpdateIncidentStatusResponse)

    - [IncidentPriority](#command-incident-v1-IncidentPriority)
    - [IncidentStatus](#command-incident-v1-IncidentStatus)

    - [IncidentCommandService](#command-incident-v1-IncidentCommandService)

- [command/membership/v1/membership.proto](#command_membership_v1_membership-proto)
    - [AssignClinicHeadDeputyRequest](#command-membership-v1-AssignClinicHeadDeputyRequest)
    - [AssignClinicHeadDeputyResponse](#command-membership-v1-AssignClinicHeadDeputyResponse)
    - [AssignClinicHeadRequest](#command-membership-v1-AssignClinicHeadRequest)
    - [AssignClinicHeadResponse](#command-membership-v1-AssignClinicHeadResponse)
    - [AssignDepartmentResponsibleDeputyRequest](#command-membership-v1-AssignDepartmentResponsibleDeputyRequest)
    - [AssignDepartmentResponsibleDeputyResponse](#command-membership-v1-AssignDepartmentResponsibleDeputyResponse)
    - [AssignDepartmentResponsibleRequest](#command-membership-v1-AssignDepartmentResponsibleRequest)
    - [AssignDepartmentResponsibleResponse](#command-membership-v1-AssignDepartmentResponsibleResponse)
    - [AssignOrganizationAdminDeputyRequest](#command-membership-v1-AssignOrganizationAdminDeputyRequest)
    - [AssignOrganizationAdminDeputyResponse](#command-membership-v1-AssignOrganizationAdminDeputyResponse)
    - [AssignOrganizationAdminRequest](#command-membership-v1-AssignOrganizationAdminRequest)
    - [AssignOrganizationAdminResponse](#command-membership-v1-AssignOrganizationAdminResponse)
    - [AssignOrganizationDispatcherDeputyRequest](#command-membership-v1-AssignOrganizationDispatcherDeputyRequest)
    - [AssignOrganizationDispatcherDeputyResponse](#command-membership-v1-AssignOrganizationDispatcherDeputyResponse)
    - [AssignOrganizationDispatcherRequest](#command-membership-v1-AssignOrganizationDispatcherRequest)
    - [AssignOrganizationDispatcherResponse](#command-membership-v1-AssignOrganizationDispatcherResponse)
    - [AssignOrganizationHeadDeputyRequest](#command-membership-v1-AssignOrganizationHeadDeputyRequest)
    - [AssignOrganizationHeadDeputyResponse](#command-membership-v1-AssignOrganizationHeadDeputyResponse)
    - [AssignOrganizationHeadRequest](#command-membership-v1-AssignOrganizationHeadRequest)
    - [AssignOrganizationHeadResponse](#command-membership-v1-AssignOrganizationHeadResponse)
    - [CancelScheduledVacationRequest](#command-membership-v1-CancelScheduledVacationRequest)
    - [CancelScheduledVacationResponse](#command-membership-v1-CancelScheduledVacationResponse)
    - [ForceEndVacationRequest](#command-membership-v1-ForceEndVacationRequest)
    - [ForceEndVacationResponse](#command-membership-v1-ForceEndVacationResponse)
    - [GrantSystemAdminRequest](#command-membership-v1-GrantSystemAdminRequest)
    - [GrantSystemAdminResponse](#command-membership-v1-GrantSystemAdminResponse)
    - [HireEmployeeRequest](#command-membership-v1-HireEmployeeRequest)
    - [HireEmployeeResponse](#command-membership-v1-HireEmployeeResponse)
    - [RemoveClinicHeadDeputyRequest](#command-membership-v1-RemoveClinicHeadDeputyRequest)
    - [RemoveClinicHeadDeputyResponse](#command-membership-v1-RemoveClinicHeadDeputyResponse)
    - [RemoveDepartmentResponsibleDeputyRequest](#command-membership-v1-RemoveDepartmentResponsibleDeputyRequest)
    - [RemoveDepartmentResponsibleDeputyResponse](#command-membership-v1-RemoveDepartmentResponsibleDeputyResponse)
    - [RemoveOrganizationAdminDeputyRequest](#command-membership-v1-RemoveOrganizationAdminDeputyRequest)
    - [RemoveOrganizationAdminDeputyResponse](#command-membership-v1-RemoveOrganizationAdminDeputyResponse)
    - [RemoveOrganizationDispatcherDeputyRequest](#command-membership-v1-RemoveOrganizationDispatcherDeputyRequest)
    - [RemoveOrganizationDispatcherDeputyResponse](#command-membership-v1-RemoveOrganizationDispatcherDeputyResponse)
    - [RemoveOrganizationHeadDeputyRequest](#command-membership-v1-RemoveOrganizationHeadDeputyRequest)
    - [RemoveOrganizationHeadDeputyResponse](#command-membership-v1-RemoveOrganizationHeadDeputyResponse)
    - [RevokeClinicHeadRequest](#command-membership-v1-RevokeClinicHeadRequest)
    - [RevokeClinicHeadResponse](#command-membership-v1-RevokeClinicHeadResponse)
    - [RevokeDepartmentResponsibleRequest](#command-membership-v1-RevokeDepartmentResponsibleRequest)
    - [RevokeDepartmentResponsibleResponse](#command-membership-v1-RevokeDepartmentResponsibleResponse)
    - [RevokeOrganizationAdminRequest](#command-membership-v1-RevokeOrganizationAdminRequest)
    - [RevokeOrganizationAdminResponse](#command-membership-v1-RevokeOrganizationAdminResponse)
    - [RevokeOrganizationDispatcherRequest](#command-membership-v1-RevokeOrganizationDispatcherRequest)
    - [RevokeOrganizationDispatcherResponse](#command-membership-v1-RevokeOrganizationDispatcherResponse)
    - [RevokeOrganizationHeadRequest](#command-membership-v1-RevokeOrganizationHeadRequest)
    - [RevokeOrganizationHeadResponse](#command-membership-v1-RevokeOrganizationHeadResponse)
    - [RevokeSystemAdminRequest](#command-membership-v1-RevokeSystemAdminRequest)
    - [RevokeSystemAdminResponse](#command-membership-v1-RevokeSystemAdminResponse)
    - [ScheduleVacationRequest](#command-membership-v1-ScheduleVacationRequest)
    - [ScheduleVacationResponse](#command-membership-v1-ScheduleVacationResponse)
    - [StartVacationNowRequest](#command-membership-v1-StartVacationNowRequest)
    - [StartVacationNowResponse](#command-membership-v1-StartVacationNowResponse)
    - [TerminateEmployeeRequest](#command-membership-v1-TerminateEmployeeRequest)
    - [TerminateEmployeeResponse](#command-membership-v1-TerminateEmployeeResponse)
    - [UpdateEmployeeDepartmentRequest](#command-membership-v1-UpdateEmployeeDepartmentRequest)
    - [UpdateEmployeeDepartmentResponse](#command-membership-v1-UpdateEmployeeDepartmentResponse)
    - [UpdateEmployeePositionRequest](#command-membership-v1-UpdateEmployeePositionRequest)
    - [UpdateEmployeePositionResponse](#command-membership-v1-UpdateEmployeePositionResponse)
    - [UpdateVacationEndDateRequest](#command-membership-v1-UpdateVacationEndDateRequest)
    - [UpdateVacationEndDateResponse](#command-membership-v1-UpdateVacationEndDateResponse)

    - [MembershipCommandService](#command-membership-v1-MembershipCommandService)

- [command/orgstructure/v1/orgstructure.proto](#command_orgstructure_v1_orgstructure-proto)
    - [AddressInput](#command-orgstructure-v1-AddressInput)
    - [CreateClinicRequest](#command-orgstructure-v1-CreateClinicRequest)
    - [CreateClinicResponse](#command-orgstructure-v1-CreateClinicResponse)
    - [CreateDepartmentRequest](#command-orgstructure-v1-CreateDepartmentRequest)
    - [CreateDepartmentResponse](#command-orgstructure-v1-CreateDepartmentResponse)
    - [CreateOrganizationRequest](#command-orgstructure-v1-CreateOrganizationRequest)
    - [CreateOrganizationResponse](#command-orgstructure-v1-CreateOrganizationResponse)
    - [PointInput](#command-orgstructure-v1-PointInput)
    - [UpdateClinicDetailsRequest](#command-orgstructure-v1-UpdateClinicDetailsRequest)
    - [UpdateClinicDetailsResponse](#command-orgstructure-v1-UpdateClinicDetailsResponse)
    - [UpdateClinicPhysicalAddressRequest](#command-orgstructure-v1-UpdateClinicPhysicalAddressRequest)
    - [UpdateClinicPhysicalAddressResponse](#command-orgstructure-v1-UpdateClinicPhysicalAddressResponse)
    - [UpdateDepartmentDetailsRequest](#command-orgstructure-v1-UpdateDepartmentDetailsRequest)
    - [UpdateDepartmentDetailsResponse](#command-orgstructure-v1-UpdateDepartmentDetailsResponse)
    - [UpdateOrganizationDetailsRequest](#command-orgstructure-v1-UpdateOrganizationDetailsRequest)
    - [UpdateOrganizationDetailsResponse](#command-orgstructure-v1-UpdateOrganizationDetailsResponse)
    - [UpdateOrganizationLegalAddressRequest](#command-orgstructure-v1-UpdateOrganizationLegalAddressRequest)
    - [UpdateOrganizationLegalAddressResponse](#command-orgstructure-v1-UpdateOrganizationLegalAddressResponse)

    - [OrgStructureCommandService](#command-orgstructure-v1-OrgStructureCommandService)

- [command/request/classifier/v1/request_classifier.proto](#command_request_classifier_v1_request_classifier-proto)
    - [CreateRequestTypeRequest](#command-request-classifier-v1-CreateRequestTypeRequest)
    - [CreateRequestTypeResponse](#command-request-classifier-v1-CreateRequestTypeResponse)
    - [DeactivateRequestTypeRequest](#command-request-classifier-v1-DeactivateRequestTypeRequest)
    - [DeactivateRequestTypeResponse](#command-request-classifier-v1-DeactivateRequestTypeResponse)
    - [DeleteRequestTypeRequest](#command-request-classifier-v1-DeleteRequestTypeRequest)
    - [DeleteRequestTypeResponse](#command-request-classifier-v1-DeleteRequestTypeResponse)
    - [ReactivateRequestTypeRequest](#command-request-classifier-v1-ReactivateRequestTypeRequest)
    - [ReactivateRequestTypeResponse](#command-request-classifier-v1-ReactivateRequestTypeResponse)
    - [UpdateRequestTypeDetailsRequest](#command-request-classifier-v1-UpdateRequestTypeDetailsRequest)
    - [UpdateRequestTypeDetailsResponse](#command-request-classifier-v1-UpdateRequestTypeDetailsResponse)

    - [RequestClassifierCommandService](#command-request-classifier-v1-RequestClassifierCommandService)

- [command/request/v1/request.proto](#command_request_v1_request-proto)
    - [AssignExecutorsRequest](#command-request-v1-AssignExecutorsRequest)
    - [AssignExecutorsResponse](#command-request-v1-AssignExecutorsResponse)
    - [CreateServiceRequestRequest](#command-request-v1-CreateServiceRequestRequest)
    - [CreateServiceRequestResponse](#command-request-v1-CreateServiceRequestResponse)
    - [UpdateServiceRequestDescriptionRequest](#command-request-v1-UpdateServiceRequestDescriptionRequest)
    - [UpdateServiceRequestDescriptionResponse](#command-request-v1-UpdateServiceRequestDescriptionResponse)
    - [UpdateServiceRequestStatusRequest](#command-request-v1-UpdateServiceRequestStatusRequest)
    - [UpdateServiceRequestStatusResponse](#command-request-v1-UpdateServiceRequestStatusResponse)

    - [ServiceRequestCommandService](#command-request-v1-ServiceRequestCommandService)

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
    - [DeptResponsibleAssigned](#event-department-v1-DeptResponsibleAssigned)
    - [DeptResponsibleDeputyAssigned](#event-department-v1-DeptResponsibleDeputyAssigned)
    - [DeptResponsibleDeputyRemoved](#event-department-v1-DeptResponsibleDeputyRemoved)
    - [DeptResponsibleRevoked](#event-department-v1-DeptResponsibleRevoked)

- [event/employee/v1/events.proto](#event_employee_v1_events-proto)
    - [EmployeeDepartmentChanged](#event-employee-v1-EmployeeDepartmentChanged)
    - [EmployeeHired](#event-employee-v1-EmployeeHired)
    - [EmployeePositionChanged](#event-employee-v1-EmployeePositionChanged)
    - [EmployeeTerminated](#event-employee-v1-EmployeeTerminated)

- [event/incident/buffer/v1/events.proto](#event_incident_buffer_v1_events-proto)
    - [PatientIncidentBufferCreated](#event-incident-buffer-v1-PatientIncidentBufferCreated)
    - [PatientIncidentBufferUpdated](#event-incident-buffer-v1-PatientIncidentBufferUpdated)

- [event/incident/classifier/v1/events.proto](#event_incident_classifier_v1_events-proto)
    - [IncidentCategoryCreated](#event-incident-classifier-v1-IncidentCategoryCreated)
    - [IncidentCategoryDeactivated](#event-incident-classifier-v1-IncidentCategoryDeactivated)
    - [IncidentCategoryDeleted](#event-incident-classifier-v1-IncidentCategoryDeleted)
    - [IncidentCategoryDetailsUpdated](#event-incident-classifier-v1-IncidentCategoryDetailsUpdated)
    - [IncidentCategoryMoved](#event-incident-classifier-v1-IncidentCategoryMoved)
    - [IncidentCategoryReactivated](#event-incident-classifier-v1-IncidentCategoryReactivated)
    - [IncidentTypeAllowedForPatients](#event-incident-classifier-v1-IncidentTypeAllowedForPatients)
    - [IncidentTypeCreated](#event-incident-classifier-v1-IncidentTypeCreated)
    - [IncidentTypeDeactivated](#event-incident-classifier-v1-IncidentTypeDeactivated)
    - [IncidentTypeDeleted](#event-incident-classifier-v1-IncidentTypeDeleted)
    - [IncidentTypeDetailsUpdated](#event-incident-classifier-v1-IncidentTypeDetailsUpdated)
    - [IncidentTypeDisallowedForPatients](#event-incident-classifier-v1-IncidentTypeDisallowedForPatients)
    - [IncidentTypeMoved](#event-incident-classifier-v1-IncidentTypeMoved)
    - [IncidentTypeReactivated](#event-incident-classifier-v1-IncidentTypeReactivated)

- [event/incident/v1/events.proto](#event_incident_v1_events-proto)
    - [IncidentCreated](#event-incident-v1-IncidentCreated)
    - [IncidentDescriptionUpdated](#event-incident-v1-IncidentDescriptionUpdated)
    - [IncidentPriorityChanged](#event-incident-v1-IncidentPriorityChanged)
    - [IncidentStatusChanged](#event-incident-v1-IncidentStatusChanged)

- [event/organization/v1/events.proto](#event_organization_v1_events-proto)
    - [Address](#event-organization-v1-Address)
    - [OrgAdminAssigned](#event-organization-v1-OrgAdminAssigned)
    - [OrgAdminDeputyAssigned](#event-organization-v1-OrgAdminDeputyAssigned)
    - [OrgAdminDeputyRemoved](#event-organization-v1-OrgAdminDeputyRemoved)
    - [OrgAdminRevoked](#event-organization-v1-OrgAdminRevoked)
    - [OrgDispatcherAssigned](#event-organization-v1-OrgDispatcherAssigned)
    - [OrgDispatcherDeputyAssigned](#event-organization-v1-OrgDispatcherDeputyAssigned)
    - [OrgDispatcherDeputyRemoved](#event-organization-v1-OrgDispatcherDeputyRemoved)
    - [OrgDispatcherRevoked](#event-organization-v1-OrgDispatcherRevoked)
    - [OrgHeadAssigned](#event-organization-v1-OrgHeadAssigned)
    - [OrgHeadDeputyAssigned](#event-organization-v1-OrgHeadDeputyAssigned)
    - [OrgHeadDeputyRemoved](#event-organization-v1-OrgHeadDeputyRemoved)
    - [OrgHeadRevoked](#event-organization-v1-OrgHeadRevoked)
    - [OrganizationCreated](#event-organization-v1-OrganizationCreated)
    - [OrganizationDetailsChanged](#event-organization-v1-OrganizationDetailsChanged)
    - [OrganizationLegalAddressChanged](#event-organization-v1-OrganizationLegalAddressChanged)
    - [Point](#event-organization-v1-Point)

- [event/request_type/v1/events.proto](#event_request_type_v1_events-proto)
    - [RequestTypeCreated](#event-request_type-v1-RequestTypeCreated)
    - [RequestTypeDeactivated](#event-request_type-v1-RequestTypeDeactivated)
    - [RequestTypeDeleted](#event-request_type-v1-RequestTypeDeleted)
    - [RequestTypeDetailsUpdated](#event-request_type-v1-RequestTypeDetailsUpdated)
    - [RequestTypeReactivated](#event-request_type-v1-RequestTypeReactivated)

- [event/service_request/v1/events.proto](#event_service_request_v1_events-proto)
    - [ServiceRequestCreated](#event-service_request-v1-ServiceRequestCreated)
    - [ServiceRequestDescriptionUpdated](#event-service_request-v1-ServiceRequestDescriptionUpdated)
    - [ServiceRequestExecutorAssigned](#event-service_request-v1-ServiceRequestExecutorAssigned)
    - [ServiceRequestExecutorRemoved](#event-service_request-v1-ServiceRequestExecutorRemoved)
    - [ServiceRequestStatusChanged](#event-service_request-v1-ServiceRequestStatusChanged)

- [event/system_admin/v1/events.proto](#event_system_admin_v1_events-proto)
    - [SystemAdminGranted](#event-system_admin-v1-SystemAdminGranted)
    - [SystemAdminRevoked](#event-system_admin-v1-SystemAdminRevoked)

- [event/v1/envelope.proto](#event_v1_envelope-proto)
    - [Envelope](#event-v1-Envelope)

- [event/vacation/v1/events.proto](#event_vacation_v1_events-proto)
    - [VacationCancelled](#event-vacation-v1-VacationCancelled)
    - [VacationEndDateChanged](#event-vacation-v1-VacationEndDateChanged)
    - [VacationEnded](#event-vacation-v1-VacationEnded)
    - [VacationScheduled](#event-vacation-v1-VacationScheduled)
    - [VacationStarted](#event-vacation-v1-VacationStarted)

- [query/analytics/v1/analytics.proto](#query_analytics_v1_analytics-proto)
    - [CategoryCount](#query-analytics-v1-CategoryCount)
    - [DepartmentCount](#query-analytics-v1-DepartmentCount)
    - [GetSnapshotRequest](#query-analytics-v1-GetSnapshotRequest)
    - [GetSnapshotResponse](#query-analytics-v1-GetSnapshotResponse)
    - [GetSummaryRequest](#query-analytics-v1-GetSummaryRequest)
    - [GetSummaryResponse](#query-analytics-v1-GetSummaryResponse)
    - [GetTimeSeriesRequest](#query-analytics-v1-GetTimeSeriesRequest)
    - [GetTimeSeriesResponse](#query-analytics-v1-GetTimeSeriesResponse)
    - [IncidentPriorityBreakdown](#query-analytics-v1-IncidentPriorityBreakdown)
    - [IncidentSourceBreakdown](#query-analytics-v1-IncidentSourceBreakdown)
    - [IncidentStatusBreakdown](#query-analytics-v1-IncidentStatusBreakdown)
    - [IncidentSummary](#query-analytics-v1-IncidentSummary)
    - [PatientBufferStatusBreakdown](#query-analytics-v1-PatientBufferStatusBreakdown)
    - [PatientBufferSummary](#query-analytics-v1-PatientBufferSummary)
    - [RequestStatusBreakdown](#query-analytics-v1-RequestStatusBreakdown)
    - [RequestSummary](#query-analytics-v1-RequestSummary)
    - [ResolutionStats](#query-analytics-v1-ResolutionStats)
    - [SnapshotIncident](#query-analytics-v1-SnapshotIncident)
    - [SnapshotPatientBuffer](#query-analytics-v1-SnapshotPatientBuffer)
    - [SnapshotRequest](#query-analytics-v1-SnapshotRequest)
    - [SummaryPeriod](#query-analytics-v1-SummaryPeriod)
    - [TimeSeriesBucket](#query-analytics-v1-TimeSeriesBucket)
    - [TimeSeriesIncidentBucket](#query-analytics-v1-TimeSeriesIncidentBucket)
    - [TimeSeriesRequestBucket](#query-analytics-v1-TimeSeriesRequestBucket)
    - [TypeCount](#query-analytics-v1-TypeCount)

    - [TimeSeriesGranularity](#query-analytics-v1-TimeSeriesGranularity)

    - [AnalyticsQueryService](#query-analytics-v1-AnalyticsQueryService)

- [query/announcement/v1/announcement.proto](#query_announcement_v1_announcement-proto)
    - [AnnouncementView](#query-announcement-v1-AnnouncementView)
    - [GetAnnouncementRequest](#query-announcement-v1-GetAnnouncementRequest)
    - [GetAnnouncementResponse](#query-announcement-v1-GetAnnouncementResponse)
    - [ListAnnouncementsForClinicRequest](#query-announcement-v1-ListAnnouncementsForClinicRequest)
    - [ListAnnouncementsForClinicResponse](#query-announcement-v1-ListAnnouncementsForClinicResponse)
    - [ListAnnouncementsForDepartmentRequest](#query-announcement-v1-ListAnnouncementsForDepartmentRequest)
    - [ListAnnouncementsForDepartmentResponse](#query-announcement-v1-ListAnnouncementsForDepartmentResponse)
    - [ListAnnouncementsForOrganizationRequest](#query-announcement-v1-ListAnnouncementsForOrganizationRequest)
    - [ListAnnouncementsForOrganizationResponse](#query-announcement-v1-ListAnnouncementsForOrganizationResponse)

    - [AnnouncementPriority](#query-announcement-v1-AnnouncementPriority)

    - [AnnouncementQueryService](#query-announcement-v1-AnnouncementQueryService)

- [query/incident/classifier/v1/classifier.proto](#query_incident_classifier_v1_classifier-proto)
    - [Category](#query-incident-classifier-v1-Category)
    - [GetCategoryRequest](#query-incident-classifier-v1-GetCategoryRequest)
    - [GetCategoryResponse](#query-incident-classifier-v1-GetCategoryResponse)
    - [GetTypeRequest](#query-incident-classifier-v1-GetTypeRequest)
    - [GetTypeResponse](#query-incident-classifier-v1-GetTypeResponse)
    - [ListActiveRootCategoriesRequest](#query-incident-classifier-v1-ListActiveRootCategoriesRequest)
    - [ListActiveRootCategoriesResponse](#query-incident-classifier-v1-ListActiveRootCategoriesResponse)
    - [ListActiveTypesByOrganizationRequest](#query-incident-classifier-v1-ListActiveTypesByOrganizationRequest)
    - [ListActiveTypesByOrganizationResponse](#query-incident-classifier-v1-ListActiveTypesByOrganizationResponse)
    - [ListCategoriesByOrganizationRequest](#query-incident-classifier-v1-ListCategoriesByOrganizationRequest)
    - [ListCategoriesByOrganizationResponse](#query-incident-classifier-v1-ListCategoriesByOrganizationResponse)
    - [ListCategorySubtreeRequest](#query-incident-classifier-v1-ListCategorySubtreeRequest)
    - [ListCategorySubtreeResponse](#query-incident-classifier-v1-ListCategorySubtreeResponse)
    - [ListPatientAllowedTypesByOrganizationRequest](#query-incident-classifier-v1-ListPatientAllowedTypesByOrganizationRequest)
    - [ListPatientAllowedTypesByOrganizationResponse](#query-incident-classifier-v1-ListPatientAllowedTypesByOrganizationResponse)
    - [ListPatientVisibleCategoriesByOrganizationRequest](#query-incident-classifier-v1-ListPatientVisibleCategoriesByOrganizationRequest)
    - [ListPatientVisibleCategoriesByOrganizationResponse](#query-incident-classifier-v1-ListPatientVisibleCategoriesByOrganizationResponse)
    - [ListTypesByCategoryRequest](#query-incident-classifier-v1-ListTypesByCategoryRequest)
    - [ListTypesByCategoryResponse](#query-incident-classifier-v1-ListTypesByCategoryResponse)
    - [Type](#query-incident-classifier-v1-Type)

    - [IncidentClassifierQueryService](#query-incident-classifier-v1-IncidentClassifierQueryService)

- [query/incident/v1/incident.proto](#query_incident_v1_incident-proto)
    - [ActorView](#query-incident-v1-ActorView)
    - [BufferEntryView](#query-incident-v1-BufferEntryView)
    - [GetBufferEntryRequest](#query-incident-v1-GetBufferEntryRequest)
    - [GetBufferEntryResponse](#query-incident-v1-GetBufferEntryResponse)
    - [GetIncidentHistoryRequest](#query-incident-v1-GetIncidentHistoryRequest)
    - [GetIncidentHistoryResponse](#query-incident-v1-GetIncidentHistoryResponse)
    - [GetIncidentRequest](#query-incident-v1-GetIncidentRequest)
    - [GetIncidentResponse](#query-incident-v1-GetIncidentResponse)
    - [IncidentView](#query-incident-v1-IncidentView)
    - [ListBufferEntriesRequest](#query-incident-v1-ListBufferEntriesRequest)
    - [ListBufferEntriesResponse](#query-incident-v1-ListBufferEntriesResponse)
    - [ListIncidentsRequest](#query-incident-v1-ListIncidentsRequest)
    - [ListIncidentsResponse](#query-incident-v1-ListIncidentsResponse)
    - [ListMyBufferEntriesRequest](#query-incident-v1-ListMyBufferEntriesRequest)
    - [ListMyBufferEntriesResponse](#query-incident-v1-ListMyBufferEntriesResponse)
    - [ListMyIncidentsRequest](#query-incident-v1-ListMyIncidentsRequest)
    - [ListMyIncidentsResponse](#query-incident-v1-ListMyIncidentsResponse)
    - [PriorityHistoryEntry](#query-incident-v1-PriorityHistoryEntry)
    - [RegistrarView](#query-incident-v1-RegistrarView)
    - [StatusHistoryEntry](#query-incident-v1-StatusHistoryEntry)

    - [BufferStatus](#query-incident-v1-BufferStatus)
    - [IncidentPriority](#query-incident-v1-IncidentPriority)
    - [IncidentStatus](#query-incident-v1-IncidentStatus)
    - [PatientStatus](#query-incident-v1-PatientStatus)

    - [IncidentQueryService](#query-incident-v1-IncidentQueryService)

- [query/membership/v1/membership.proto](#query_membership_v1_membership-proto)
    - [CountEmployeesByClinicRequest](#query-membership-v1-CountEmployeesByClinicRequest)
    - [CountEmployeesByClinicResponse](#query-membership-v1-CountEmployeesByClinicResponse)
    - [CountEmployeesByDepartmentRequest](#query-membership-v1-CountEmployeesByDepartmentRequest)
    - [CountEmployeesByDepartmentResponse](#query-membership-v1-CountEmployeesByDepartmentResponse)
    - [CountEmployeesByOrganizationRequest](#query-membership-v1-CountEmployeesByOrganizationRequest)
    - [CountEmployeesByOrganizationResponse](#query-membership-v1-CountEmployeesByOrganizationResponse)
    - [CountVacationsByEmployeeRequest](#query-membership-v1-CountVacationsByEmployeeRequest)
    - [CountVacationsByEmployeeResponse](#query-membership-v1-CountVacationsByEmployeeResponse)
    - [EmployeeCardView](#query-membership-v1-EmployeeCardView)
    - [GetClinicHeadRequest](#query-membership-v1-GetClinicHeadRequest)
    - [GetClinicHeadResponse](#query-membership-v1-GetClinicHeadResponse)
    - [GetDepartmentResponsibleRequest](#query-membership-v1-GetDepartmentResponsibleRequest)
    - [GetDepartmentResponsibleResponse](#query-membership-v1-GetDepartmentResponsibleResponse)
    - [GetEmployeeRequest](#query-membership-v1-GetEmployeeRequest)
    - [GetEmployeeResponse](#query-membership-v1-GetEmployeeResponse)
    - [ListCandidatesForClinicHeadRequest](#query-membership-v1-ListCandidatesForClinicHeadRequest)
    - [ListCandidatesForClinicHeadResponse](#query-membership-v1-ListCandidatesForClinicHeadResponse)
    - [ListCandidatesForDeptResponsibleRequest](#query-membership-v1-ListCandidatesForDeptResponsibleRequest)
    - [ListCandidatesForDeptResponsibleResponse](#query-membership-v1-ListCandidatesForDeptResponsibleResponse)
    - [ListCandidatesForHireRequest](#query-membership-v1-ListCandidatesForHireRequest)
    - [ListCandidatesForHireResponse](#query-membership-v1-ListCandidatesForHireResponse)
    - [ListCandidatesForOrgAdminRequest](#query-membership-v1-ListCandidatesForOrgAdminRequest)
    - [ListCandidatesForOrgAdminResponse](#query-membership-v1-ListCandidatesForOrgAdminResponse)
    - [ListCandidatesForOrgDispatcherRequest](#query-membership-v1-ListCandidatesForOrgDispatcherRequest)
    - [ListCandidatesForOrgDispatcherResponse](#query-membership-v1-ListCandidatesForOrgDispatcherResponse)
    - [ListCandidatesForOrgHeadRequest](#query-membership-v1-ListCandidatesForOrgHeadRequest)
    - [ListCandidatesForOrgHeadResponse](#query-membership-v1-ListCandidatesForOrgHeadResponse)
    - [ListCandidatesForSystemAdminRequest](#query-membership-v1-ListCandidatesForSystemAdminRequest)
    - [ListCandidatesForSystemAdminResponse](#query-membership-v1-ListCandidatesForSystemAdminResponse)
    - [ListEmployeesByClinicRequest](#query-membership-v1-ListEmployeesByClinicRequest)
    - [ListEmployeesByClinicResponse](#query-membership-v1-ListEmployeesByClinicResponse)
    - [ListEmployeesByDepartmentRequest](#query-membership-v1-ListEmployeesByDepartmentRequest)
    - [ListEmployeesByDepartmentResponse](#query-membership-v1-ListEmployeesByDepartmentResponse)
    - [ListEmployeesByOrganizationRequest](#query-membership-v1-ListEmployeesByOrganizationRequest)
    - [ListEmployeesByOrganizationResponse](#query-membership-v1-ListEmployeesByOrganizationResponse)
    - [ListOrgAdminsRequest](#query-membership-v1-ListOrgAdminsRequest)
    - [ListOrgAdminsResponse](#query-membership-v1-ListOrgAdminsResponse)
    - [ListOrgDispatchersRequest](#query-membership-v1-ListOrgDispatchersRequest)
    - [ListOrgDispatchersResponse](#query-membership-v1-ListOrgDispatchersResponse)
    - [ListOrgHeadsRequest](#query-membership-v1-ListOrgHeadsRequest)
    - [ListOrgHeadsResponse](#query-membership-v1-ListOrgHeadsResponse)
    - [ListSystemAdminsRequest](#query-membership-v1-ListSystemAdminsRequest)
    - [ListSystemAdminsResponse](#query-membership-v1-ListSystemAdminsResponse)
    - [ListVacationsByEmployeeRequest](#query-membership-v1-ListVacationsByEmployeeRequest)
    - [ListVacationsByEmployeeResponse](#query-membership-v1-ListVacationsByEmployeeResponse)
    - [RoleAssignment](#query-membership-v1-RoleAssignment)
    - [SearchEmployeesByOrganizationRequest](#query-membership-v1-SearchEmployeesByOrganizationRequest)
    - [SearchEmployeesByOrganizationResponse](#query-membership-v1-SearchEmployeesByOrganizationResponse)
    - [SystemAdminView](#query-membership-v1-SystemAdminView)
    - [VacationView](#query-membership-v1-VacationView)
    - [ZitadelUserView](#query-membership-v1-ZitadelUserView)

    - [MembershipQueryService](#query-membership-v1-MembershipQueryService)

- [query/orgstructure/v1/orgstructure.proto](#query_orgstructure_v1_orgstructure-proto)
    - [Address](#query-orgstructure-v1-Address)
    - [Clinic](#query-orgstructure-v1-Clinic)
    - [ClinicListItem](#query-orgstructure-v1-ClinicListItem)
    - [CountClinicsByOrganizationRequest](#query-orgstructure-v1-CountClinicsByOrganizationRequest)
    - [CountClinicsByOrganizationResponse](#query-orgstructure-v1-CountClinicsByOrganizationResponse)
    - [CountDepartmentsByClinicRequest](#query-orgstructure-v1-CountDepartmentsByClinicRequest)
    - [CountDepartmentsByClinicResponse](#query-orgstructure-v1-CountDepartmentsByClinicResponse)
    - [CountOrganizationsRequest](#query-orgstructure-v1-CountOrganizationsRequest)
    - [CountOrganizationsResponse](#query-orgstructure-v1-CountOrganizationsResponse)
    - [Department](#query-orgstructure-v1-Department)
    - [DepartmentListItem](#query-orgstructure-v1-DepartmentListItem)
    - [GetClinicRequest](#query-orgstructure-v1-GetClinicRequest)
    - [GetClinicResponse](#query-orgstructure-v1-GetClinicResponse)
    - [GetDepartmentRequest](#query-orgstructure-v1-GetDepartmentRequest)
    - [GetDepartmentResponse](#query-orgstructure-v1-GetDepartmentResponse)
    - [GetOrganizationRequest](#query-orgstructure-v1-GetOrganizationRequest)
    - [GetOrganizationResponse](#query-orgstructure-v1-GetOrganizationResponse)
    - [ListClinicsByOrganizationRequest](#query-orgstructure-v1-ListClinicsByOrganizationRequest)
    - [ListClinicsByOrganizationResponse](#query-orgstructure-v1-ListClinicsByOrganizationResponse)
    - [ListDepartmentsByClinicRequest](#query-orgstructure-v1-ListDepartmentsByClinicRequest)
    - [ListDepartmentsByClinicResponse](#query-orgstructure-v1-ListDepartmentsByClinicResponse)
    - [ListOrganizationsRequest](#query-orgstructure-v1-ListOrganizationsRequest)
    - [ListOrganizationsResponse](#query-orgstructure-v1-ListOrganizationsResponse)
    - [Organization](#query-orgstructure-v1-Organization)
    - [OrganizationListItem](#query-orgstructure-v1-OrganizationListItem)
    - [Point](#query-orgstructure-v1-Point)
    - [SearchOrganizationsRequest](#query-orgstructure-v1-SearchOrganizationsRequest)
    - [SearchOrganizationsResponse](#query-orgstructure-v1-SearchOrganizationsResponse)

    - [OrgStructureQueryService](#query-orgstructure-v1-OrgStructureQueryService)

- [query/request/classifier/v1/classifier.proto](#query_request_classifier_v1_classifier-proto)
    - [GetRequestTypeRequest](#query-request-classifier-v1-GetRequestTypeRequest)
    - [GetRequestTypeResponse](#query-request-classifier-v1-GetRequestTypeResponse)
    - [ListActiveRequestTypesByOrganizationRequest](#query-request-classifier-v1-ListActiveRequestTypesByOrganizationRequest)
    - [ListActiveRequestTypesByOrganizationResponse](#query-request-classifier-v1-ListActiveRequestTypesByOrganizationResponse)
    - [ListRequestTypesByOrganizationRequest](#query-request-classifier-v1-ListRequestTypesByOrganizationRequest)
    - [ListRequestTypesByOrganizationResponse](#query-request-classifier-v1-ListRequestTypesByOrganizationResponse)
    - [RequestType](#query-request-classifier-v1-RequestType)

    - [RequestClassifierQueryService](#query-request-classifier-v1-RequestClassifierQueryService)

- [query/request/v1/request.proto](#query_request_v1_request-proto)
    - [Executor](#query-request-v1-Executor)
    - [ExecutorHistoryEntry](#query-request-v1-ExecutorHistoryEntry)
    - [GetServiceRequestHistoryRequest](#query-request-v1-GetServiceRequestHistoryRequest)
    - [GetServiceRequestHistoryResponse](#query-request-v1-GetServiceRequestHistoryResponse)
    - [GetServiceRequestRequest](#query-request-v1-GetServiceRequestRequest)
    - [GetServiceRequestResponse](#query-request-v1-GetServiceRequestResponse)
    - [ListServiceRequestsByIncidentRequest](#query-request-v1-ListServiceRequestsByIncidentRequest)
    - [ListServiceRequestsByIncidentResponse](#query-request-v1-ListServiceRequestsByIncidentResponse)
    - [ListServiceRequestsRequest](#query-request-v1-ListServiceRequestsRequest)
    - [ListServiceRequestsResponse](#query-request-v1-ListServiceRequestsResponse)
    - [ServiceRequest](#query-request-v1-ServiceRequest)
    - [StatusHistoryEntry](#query-request-v1-StatusHistoryEntry)

    - [ServiceRequestQueryService](#query-request-v1-ServiceRequestQueryService)

- [query/self/v1/self.proto](#query_self_v1_self-proto)
    - [GetMyClinicRoleRequest](#query-self-v1-GetMyClinicRoleRequest)
    - [GetMyClinicRoleResponse](#query-self-v1-GetMyClinicRoleResponse)
    - [GetMyDepartmentRoleRequest](#query-self-v1-GetMyDepartmentRoleRequest)
    - [GetMyDepartmentRoleResponse](#query-self-v1-GetMyDepartmentRoleResponse)
    - [GetMyEmploymentRequest](#query-self-v1-GetMyEmploymentRequest)
    - [GetMyEmploymentResponse](#query-self-v1-GetMyEmploymentResponse)
    - [GetMyIdentityRequest](#query-self-v1-GetMyIdentityRequest)
    - [GetMyIdentityResponse](#query-self-v1-GetMyIdentityResponse)
    - [GetMyOrganizationRoleRequest](#query-self-v1-GetMyOrganizationRoleRequest)
    - [GetMyOrganizationRoleResponse](#query-self-v1-GetMyOrganizationRoleResponse)
    - [ListMyOrganizationsRequest](#query-self-v1-ListMyOrganizationsRequest)
    - [ListMyOrganizationsResponse](#query-self-v1-ListMyOrganizationsResponse)

    - [SelfQueryService](#query-self-v1-SelfQueryService)

- [query/stats/v1/stats.proto](#query_stats_v1_stats-proto)
    - [ClinicStats](#query-stats-v1-ClinicStats)
    - [DepartmentStats](#query-stats-v1-DepartmentStats)
    - [GetClinicStatsRequest](#query-stats-v1-GetClinicStatsRequest)
    - [GetClinicStatsResponse](#query-stats-v1-GetClinicStatsResponse)
    - [GetDepartmentStatsRequest](#query-stats-v1-GetDepartmentStatsRequest)
    - [GetDepartmentStatsResponse](#query-stats-v1-GetDepartmentStatsResponse)
    - [GetOrganizationStatsRequest](#query-stats-v1-GetOrganizationStatsRequest)
    - [GetOrganizationStatsResponse](#query-stats-v1-GetOrganizationStatsResponse)
    - [OrganizationStats](#query-stats-v1-OrganizationStats)

    - [StatsQueryService](#query-stats-v1-StatsQueryService)

- [Scalar Value Types](#scalar-value-types)



<a name="error_v1_error-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## error/v1/error.proto
api/proto/error/v1/error.proto


<a name="error-v1-ErrorCode"></a>

### ErrorCode
ErrorCode is always present in gRPC status details for client-visible
errors. It carries the machine-readable domain error code.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [string](#string) |  |  |






<a name="error-v1-ErrorResponse"></a>

### ErrorResponse
ErrorResponse is the HTTP response body written by the gateway error
handler for all client-visible errors. The details field is present
only when code = &#34;validation_failed&#34;.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| code | [string](#string) |  |  |
| message | [string](#string) |  |  |
| details | [ValidationFailedDetails](#error-v1-ValidationFailedDetails) | optional |  |






<a name="error-v1-ValidationFailedDetails"></a>

### ValidationFailedDetails
ValidationFailedDetails is present only when code = &#34;validation_failed&#34;.
Each violation corresponds to one struct-tag rule failure or one
domain-level leaf error from errors.Join.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| violations | [ValidationFailedDetails.FieldViolation](#error-v1-ValidationFailedDetails-FieldViolation) | repeated |  |






<a name="error-v1-ValidationFailedDetails-FieldViolation"></a>

### ValidationFailedDetails.FieldViolation



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| field | [string](#string) |  |  |
| rule | [string](#string) |  |  |
| message | [string](#string) |  |  |
| param | [string](#string) | optional |  |















<a name="command_announcement_v1_announcement-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/announcement/v1/announcement.proto



<a name="command-announcement-v1-ArchiveAnnouncementRequest"></a>

### ArchiveAnnouncementRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="command-announcement-v1-ArchiveAnnouncementResponse"></a>

### ArchiveAnnouncementResponse







<a name="command-announcement-v1-CreateAnnouncementRequest"></a>

### CreateAnnouncementRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| clinic_id | [string](#string) | optional |  |
| department_id | [string](#string) | optional |  |
| title | [string](#string) |  |  |
| content | [string](#string) |  |  |
| priority | [AnnouncementPriority](#command-announcement-v1-AnnouncementPriority) |  |  |
| starts_at | [string](#string) | optional | RFC3339Nano |
| ends_at | [string](#string) | optional | RFC3339Nano |






<a name="command-announcement-v1-CreateAnnouncementResponse"></a>

### CreateAnnouncementResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="command-announcement-v1-UnarchiveAnnouncementRequest"></a>

### UnarchiveAnnouncementRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="command-announcement-v1-UnarchiveAnnouncementResponse"></a>

### UnarchiveAnnouncementResponse







<a name="command-announcement-v1-UpdateAnnouncementPriorityRequest"></a>

### UpdateAnnouncementPriorityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| priority | [AnnouncementPriority](#command-announcement-v1-AnnouncementPriority) |  |  |






<a name="command-announcement-v1-UpdateAnnouncementPriorityResponse"></a>

### UpdateAnnouncementPriorityResponse







<a name="command-announcement-v1-UpdateAnnouncementRequest"></a>

### UpdateAnnouncementRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| title | [string](#string) |  |  |
| content | [string](#string) |  |  |
| starts_at | [string](#string) | optional | RFC3339Nano; absent = clear |
| ends_at | [string](#string) | optional | RFC3339Nano; absent = clear |






<a name="command-announcement-v1-UpdateAnnouncementResponse"></a>

### UpdateAnnouncementResponse









<a name="command-announcement-v1-AnnouncementPriority"></a>

### AnnouncementPriority


| Name | Number | Description |
| ---- | ------ | ----------- |
| ANNOUNCEMENT_PRIORITY_UNSPECIFIED | 0 |  |
| ANNOUNCEMENT_PRIORITY_NORMAL | 1 |  |
| ANNOUNCEMENT_PRIORITY_HIGH | 2 |  |







<a name="command-announcement-v1-AnnouncementCommandService"></a>

### AnnouncementCommandService
AnnouncementCommandService is the write-side contract for announcements.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateAnnouncement | [CreateAnnouncementRequest](#command-announcement-v1-CreateAnnouncementRequest) | [CreateAnnouncementResponse](#command-announcement-v1-CreateAnnouncementResponse) |  |
| UpdateAnnouncement | [UpdateAnnouncementRequest](#command-announcement-v1-UpdateAnnouncementRequest) | [UpdateAnnouncementResponse](#command-announcement-v1-UpdateAnnouncementResponse) |  |
| UpdateAnnouncementPriority | [UpdateAnnouncementPriorityRequest](#command-announcement-v1-UpdateAnnouncementPriorityRequest) | [UpdateAnnouncementPriorityResponse](#command-announcement-v1-UpdateAnnouncementPriorityResponse) |  |
| ArchiveAnnouncement | [ArchiveAnnouncementRequest](#command-announcement-v1-ArchiveAnnouncementRequest) | [ArchiveAnnouncementResponse](#command-announcement-v1-ArchiveAnnouncementResponse) |  |
| UnarchiveAnnouncement | [UnarchiveAnnouncementRequest](#command-announcement-v1-UnarchiveAnnouncementRequest) | [UnarchiveAnnouncementResponse](#command-announcement-v1-UnarchiveAnnouncementResponse) |  |





<a name="command_incident_buffer_v1_buffer-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/incident/buffer/v1/buffer.proto



<a name="command-incident-buffer-v1-CancelPatientIncidentRequest"></a>

### CancelPatientIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| buffer_id | [string](#string) |  |  |






<a name="command-incident-buffer-v1-CancelPatientIncidentResponse"></a>

### CancelPatientIncidentResponse







<a name="command-incident-buffer-v1-PublishPatientIncidentRequest"></a>

### PublishPatientIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| buffer_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| category_id | [string](#string) |  |  |
| type_id | [string](#string) |  |  |
| description | [string](#string) | optional | dispatcher&#39;s edited description |






<a name="command-incident-buffer-v1-PublishPatientIncidentResponse"></a>

### PublishPatientIncidentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |






<a name="command-incident-buffer-v1-RejectPatientIncidentRequest"></a>

### RejectPatientIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| buffer_id | [string](#string) |  |  |






<a name="command-incident-buffer-v1-RejectPatientIncidentResponse"></a>

### RejectPatientIncidentResponse







<a name="command-incident-buffer-v1-SubmitPatientIncidentRequest"></a>

### SubmitPatientIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| category_id | [string](#string) | optional |  |
| type_id | [string](#string) | optional |  |
| description | [string](#string) | optional |  |
| occurred_at | [string](#string) | optional | RFC3339Nano |






<a name="command-incident-buffer-v1-SubmitPatientIncidentResponse"></a>

### SubmitPatientIncidentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| buffer_id | [string](#string) |  |  |






<a name="command-incident-buffer-v1-UpdatePatientIncidentRequest"></a>

### UpdatePatientIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| buffer_id | [string](#string) |  |  |
| category_id | [string](#string) | optional |  |
| type_id | [string](#string) | optional |  |
| description | [string](#string) | optional |  |
| occurred_at | [string](#string) | optional |  |






<a name="command-incident-buffer-v1-UpdatePatientIncidentResponse"></a>

### UpdatePatientIncidentResponse













<a name="command-incident-buffer-v1-IncidentBufferCommandService"></a>

### IncidentBufferCommandService
IncidentBufferCommandService handles the patient submission flow:
patients SubmitPatientIncident, may UpdatePatientIncident or
CancelPatientIncident while pending; dispatchers PublishPatientIncident
(creates a real incident) or RejectPatientIncident.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SubmitPatientIncident | [SubmitPatientIncidentRequest](#command-incident-buffer-v1-SubmitPatientIncidentRequest) | [SubmitPatientIncidentResponse](#command-incident-buffer-v1-SubmitPatientIncidentResponse) |  |
| UpdatePatientIncident | [UpdatePatientIncidentRequest](#command-incident-buffer-v1-UpdatePatientIncidentRequest) | [UpdatePatientIncidentResponse](#command-incident-buffer-v1-UpdatePatientIncidentResponse) |  |
| CancelPatientIncident | [CancelPatientIncidentRequest](#command-incident-buffer-v1-CancelPatientIncidentRequest) | [CancelPatientIncidentResponse](#command-incident-buffer-v1-CancelPatientIncidentResponse) |  |
| PublishPatientIncident | [PublishPatientIncidentRequest](#command-incident-buffer-v1-PublishPatientIncidentRequest) | [PublishPatientIncidentResponse](#command-incident-buffer-v1-PublishPatientIncidentResponse) |  |
| RejectPatientIncident | [RejectPatientIncidentRequest](#command-incident-buffer-v1-RejectPatientIncidentRequest) | [RejectPatientIncidentResponse](#command-incident-buffer-v1-RejectPatientIncidentResponse) |  |





<a name="command_incident_classifier_v1_incident_classifier-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/incident/classifier/v1/incident_classifier.proto



<a name="command-incident-classifier-v1-AllowIncidentTypeForPatientsRequest"></a>

### AllowIncidentTypeForPatientsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-AllowIncidentTypeForPatientsResponse"></a>

### AllowIncidentTypeForPatientsResponse







<a name="command-incident-classifier-v1-CreateIncidentCategoryRequest"></a>

### CreateIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| parent_category_id | [string](#string) | optional |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-incident-classifier-v1-CreateIncidentCategoryResponse"></a>

### CreateIncidentCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-CreateIncidentTypeRequest"></a>

### CreateIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-incident-classifier-v1-CreateIncidentTypeResponse"></a>

### CreateIncidentTypeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-DeactivateIncidentCategoryRequest"></a>

### DeactivateIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-DeactivateIncidentCategoryResponse"></a>

### DeactivateIncidentCategoryResponse







<a name="command-incident-classifier-v1-DeactivateIncidentTypeRequest"></a>

### DeactivateIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-DeactivateIncidentTypeResponse"></a>

### DeactivateIncidentTypeResponse







<a name="command-incident-classifier-v1-DeleteIncidentCategoryRequest"></a>

### DeleteIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-DeleteIncidentCategoryResponse"></a>

### DeleteIncidentCategoryResponse







<a name="command-incident-classifier-v1-DeleteIncidentTypeRequest"></a>

### DeleteIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-DeleteIncidentTypeResponse"></a>

### DeleteIncidentTypeResponse







<a name="command-incident-classifier-v1-DisallowIncidentTypeForPatientsRequest"></a>

### DisallowIncidentTypeForPatientsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-DisallowIncidentTypeForPatientsResponse"></a>

### DisallowIncidentTypeForPatientsResponse







<a name="command-incident-classifier-v1-MoveIncidentCategoryRequest"></a>

### MoveIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| new_parent_category_id | [string](#string) | optional |  |






<a name="command-incident-classifier-v1-MoveIncidentCategoryResponse"></a>

### MoveIncidentCategoryResponse







<a name="command-incident-classifier-v1-MoveIncidentTypeRequest"></a>

### MoveIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| new_category_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-MoveIncidentTypeResponse"></a>

### MoveIncidentTypeResponse







<a name="command-incident-classifier-v1-ReactivateIncidentCategoryRequest"></a>

### ReactivateIncidentCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-ReactivateIncidentCategoryResponse"></a>

### ReactivateIncidentCategoryResponse







<a name="command-incident-classifier-v1-ReactivateIncidentTypeRequest"></a>

### ReactivateIncidentTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-incident-classifier-v1-ReactivateIncidentTypeResponse"></a>

### ReactivateIncidentTypeResponse







<a name="command-incident-classifier-v1-UpdateIncidentCategoryDetailsRequest"></a>

### UpdateIncidentCategoryDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-incident-classifier-v1-UpdateIncidentCategoryDetailsResponse"></a>

### UpdateIncidentCategoryDetailsResponse







<a name="command-incident-classifier-v1-UpdateIncidentTypeDetailsRequest"></a>

### UpdateIncidentTypeDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-incident-classifier-v1-UpdateIncidentTypeDetailsResponse"></a>

### UpdateIncidentTypeDetailsResponse













<a name="command-incident-classifier-v1-IncidentClassifierCommandService"></a>

### IncidentClassifierCommandService
IncidentClassifierService is the command-side contract for a single
administrative surface that manages an organisation&#39;s incident
classifier: a tree of incident categories (max depth 5) with incident
types living inside any category. Every mutation returns an id (for
Create) or an empty response and publishes one domain event per
touched row.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateIncidentCategory | [CreateIncidentCategoryRequest](#command-incident-classifier-v1-CreateIncidentCategoryRequest) | [CreateIncidentCategoryResponse](#command-incident-classifier-v1-CreateIncidentCategoryResponse) | --- Categories --- |
| UpdateIncidentCategoryDetails | [UpdateIncidentCategoryDetailsRequest](#command-incident-classifier-v1-UpdateIncidentCategoryDetailsRequest) | [UpdateIncidentCategoryDetailsResponse](#command-incident-classifier-v1-UpdateIncidentCategoryDetailsResponse) |  |
| MoveIncidentCategory | [MoveIncidentCategoryRequest](#command-incident-classifier-v1-MoveIncidentCategoryRequest) | [MoveIncidentCategoryResponse](#command-incident-classifier-v1-MoveIncidentCategoryResponse) |  |
| DeactivateIncidentCategory | [DeactivateIncidentCategoryRequest](#command-incident-classifier-v1-DeactivateIncidentCategoryRequest) | [DeactivateIncidentCategoryResponse](#command-incident-classifier-v1-DeactivateIncidentCategoryResponse) |  |
| ReactivateIncidentCategory | [ReactivateIncidentCategoryRequest](#command-incident-classifier-v1-ReactivateIncidentCategoryRequest) | [ReactivateIncidentCategoryResponse](#command-incident-classifier-v1-ReactivateIncidentCategoryResponse) |  |
| DeleteIncidentCategory | [DeleteIncidentCategoryRequest](#command-incident-classifier-v1-DeleteIncidentCategoryRequest) | [DeleteIncidentCategoryResponse](#command-incident-classifier-v1-DeleteIncidentCategoryResponse) |  |
| CreateIncidentType | [CreateIncidentTypeRequest](#command-incident-classifier-v1-CreateIncidentTypeRequest) | [CreateIncidentTypeResponse](#command-incident-classifier-v1-CreateIncidentTypeResponse) | --- Types --- |
| UpdateIncidentTypeDetails | [UpdateIncidentTypeDetailsRequest](#command-incident-classifier-v1-UpdateIncidentTypeDetailsRequest) | [UpdateIncidentTypeDetailsResponse](#command-incident-classifier-v1-UpdateIncidentTypeDetailsResponse) |  |
| MoveIncidentType | [MoveIncidentTypeRequest](#command-incident-classifier-v1-MoveIncidentTypeRequest) | [MoveIncidentTypeResponse](#command-incident-classifier-v1-MoveIncidentTypeResponse) |  |
| DeactivateIncidentType | [DeactivateIncidentTypeRequest](#command-incident-classifier-v1-DeactivateIncidentTypeRequest) | [DeactivateIncidentTypeResponse](#command-incident-classifier-v1-DeactivateIncidentTypeResponse) |  |
| ReactivateIncidentType | [ReactivateIncidentTypeRequest](#command-incident-classifier-v1-ReactivateIncidentTypeRequest) | [ReactivateIncidentTypeResponse](#command-incident-classifier-v1-ReactivateIncidentTypeResponse) |  |
| DeleteIncidentType | [DeleteIncidentTypeRequest](#command-incident-classifier-v1-DeleteIncidentTypeRequest) | [DeleteIncidentTypeResponse](#command-incident-classifier-v1-DeleteIncidentTypeResponse) |  |
| AllowIncidentTypeForPatients | [AllowIncidentTypeForPatientsRequest](#command-incident-classifier-v1-AllowIncidentTypeForPatientsRequest) | [AllowIncidentTypeForPatientsResponse](#command-incident-classifier-v1-AllowIncidentTypeForPatientsResponse) |  |
| DisallowIncidentTypeForPatients | [DisallowIncidentTypeForPatientsRequest](#command-incident-classifier-v1-DisallowIncidentTypeForPatientsRequest) | [DisallowIncidentTypeForPatientsResponse](#command-incident-classifier-v1-DisallowIncidentTypeForPatientsResponse) |  |





<a name="command_incident_v1_incident-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/incident/v1/incident.proto



<a name="command-incident-v1-CancelIncidentRequest"></a>

### CancelIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |






<a name="command-incident-v1-CancelIncidentResponse"></a>

### CancelIncidentResponse







<a name="command-incident-v1-CreateIncidentRequest"></a>

### CreateIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| category_id | [string](#string) |  |  |
| type_id | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| occurred_at | [string](#string) |  | RFC3339Nano |






<a name="command-incident-v1-CreateIncidentResponse"></a>

### CreateIncidentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |






<a name="command-incident-v1-ReopenIncidentRequest"></a>

### ReopenIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |






<a name="command-incident-v1-ReopenIncidentResponse"></a>

### ReopenIncidentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| reopened_incident_id | [string](#string) |  |  |






<a name="command-incident-v1-UpdateIncidentDescriptionRequest"></a>

### UpdateIncidentDescriptionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-incident-v1-UpdateIncidentDescriptionResponse"></a>

### UpdateIncidentDescriptionResponse







<a name="command-incident-v1-UpdateIncidentPriorityRequest"></a>

### UpdateIncidentPriorityRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |
| priority | [IncidentPriority](#command-incident-v1-IncidentPriority) |  |  |






<a name="command-incident-v1-UpdateIncidentPriorityResponse"></a>

### UpdateIncidentPriorityResponse







<a name="command-incident-v1-UpdateIncidentStatusRequest"></a>

### UpdateIncidentStatusRequest
UpdateIncidentStatus only handles forward transitions:
pending -&gt; in_progress, in_progress -&gt; done, in_progress -&gt; rejected.
Cancellation by registrar uses CancelIncident.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |
| new_status | [IncidentStatus](#command-incident-v1-IncidentStatus) |  |  |






<a name="command-incident-v1-UpdateIncidentStatusResponse"></a>

### UpdateIncidentStatusResponse









<a name="command-incident-v1-IncidentPriority"></a>

### IncidentPriority


| Name | Number | Description |
| ---- | ------ | ----------- |
| INCIDENT_PRIORITY_UNSPECIFIED | 0 |  |
| INCIDENT_PRIORITY_LOW | 1 |  |
| INCIDENT_PRIORITY_NORMAL | 2 |  |
| INCIDENT_PRIORITY_HIGH | 3 |  |
| INCIDENT_PRIORITY_CRITICAL | 4 |  |



<a name="command-incident-v1-IncidentStatus"></a>

### IncidentStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| INCIDENT_STATUS_UNSPECIFIED | 0 |  |
| INCIDENT_STATUS_PENDING | 1 |  |
| INCIDENT_STATUS_IN_PROGRESS | 2 |  |
| INCIDENT_STATUS_DONE | 3 |  |
| INCIDENT_STATUS_REJECTED | 4 |  |
| INCIDENT_STATUS_CANCELLED | 5 |  |







<a name="command-incident-v1-IncidentCommandService"></a>

### IncidentCommandService
IncidentCommandService is the write-side contract for incidents
created directly by employees (not the patient buffer).

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateIncident | [CreateIncidentRequest](#command-incident-v1-CreateIncidentRequest) | [CreateIncidentResponse](#command-incident-v1-CreateIncidentResponse) |  |
| CancelIncident | [CancelIncidentRequest](#command-incident-v1-CancelIncidentRequest) | [CancelIncidentResponse](#command-incident-v1-CancelIncidentResponse) |  |
| UpdateIncidentStatus | [UpdateIncidentStatusRequest](#command-incident-v1-UpdateIncidentStatusRequest) | [UpdateIncidentStatusResponse](#command-incident-v1-UpdateIncidentStatusResponse) |  |
| UpdateIncidentPriority | [UpdateIncidentPriorityRequest](#command-incident-v1-UpdateIncidentPriorityRequest) | [UpdateIncidentPriorityResponse](#command-incident-v1-UpdateIncidentPriorityResponse) |  |
| UpdateIncidentDescription | [UpdateIncidentDescriptionRequest](#command-incident-v1-UpdateIncidentDescriptionRequest) | [UpdateIncidentDescriptionResponse](#command-incident-v1-UpdateIncidentDescriptionResponse) |  |
| ReopenIncident | [ReopenIncidentRequest](#command-incident-v1-ReopenIncidentRequest) | [ReopenIncidentResponse](#command-incident-v1-ReopenIncidentResponse) |  |





<a name="command_membership_v1_membership-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/membership/v1/membership.proto



<a name="command-membership-v1-AssignClinicHeadDeputyRequest"></a>

### AssignClinicHeadDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignClinicHeadDeputyResponse"></a>

### AssignClinicHeadDeputyResponse







<a name="command-membership-v1-AssignClinicHeadRequest"></a>

### AssignClinicHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignClinicHeadResponse"></a>

### AssignClinicHeadResponse







<a name="command-membership-v1-AssignDepartmentResponsibleDeputyRequest"></a>

### AssignDepartmentResponsibleDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignDepartmentResponsibleDeputyResponse"></a>

### AssignDepartmentResponsibleDeputyResponse







<a name="command-membership-v1-AssignDepartmentResponsibleRequest"></a>

### AssignDepartmentResponsibleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignDepartmentResponsibleResponse"></a>

### AssignDepartmentResponsibleResponse







<a name="command-membership-v1-AssignOrganizationAdminDeputyRequest"></a>

### AssignOrganizationAdminDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignOrganizationAdminDeputyResponse"></a>

### AssignOrganizationAdminDeputyResponse







<a name="command-membership-v1-AssignOrganizationAdminRequest"></a>

### AssignOrganizationAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignOrganizationAdminResponse"></a>

### AssignOrganizationAdminResponse







<a name="command-membership-v1-AssignOrganizationDispatcherDeputyRequest"></a>

### AssignOrganizationDispatcherDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignOrganizationDispatcherDeputyResponse"></a>

### AssignOrganizationDispatcherDeputyResponse







<a name="command-membership-v1-AssignOrganizationDispatcherRequest"></a>

### AssignOrganizationDispatcherRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignOrganizationDispatcherResponse"></a>

### AssignOrganizationDispatcherResponse







<a name="command-membership-v1-AssignOrganizationHeadDeputyRequest"></a>

### AssignOrganizationHeadDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignOrganizationHeadDeputyResponse"></a>

### AssignOrganizationHeadDeputyResponse







<a name="command-membership-v1-AssignOrganizationHeadRequest"></a>

### AssignOrganizationHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-AssignOrganizationHeadResponse"></a>

### AssignOrganizationHeadResponse







<a name="command-membership-v1-CancelScheduledVacationRequest"></a>

### CancelScheduledVacationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="command-membership-v1-CancelScheduledVacationResponse"></a>

### CancelScheduledVacationResponse







<a name="command-membership-v1-ForceEndVacationRequest"></a>

### ForceEndVacationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="command-membership-v1-ForceEndVacationResponse"></a>

### ForceEndVacationResponse







<a name="command-membership-v1-GrantSystemAdminRequest"></a>

### GrantSystemAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |






<a name="command-membership-v1-GrantSystemAdminResponse"></a>

### GrantSystemAdminResponse







<a name="command-membership-v1-HireEmployeeRequest"></a>

### HireEmployeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| position | [string](#string) | optional |  |






<a name="command-membership-v1-HireEmployeeResponse"></a>

### HireEmployeeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RemoveClinicHeadDeputyRequest"></a>

### RemoveClinicHeadDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RemoveClinicHeadDeputyResponse"></a>

### RemoveClinicHeadDeputyResponse







<a name="command-membership-v1-RemoveDepartmentResponsibleDeputyRequest"></a>

### RemoveDepartmentResponsibleDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RemoveDepartmentResponsibleDeputyResponse"></a>

### RemoveDepartmentResponsibleDeputyResponse







<a name="command-membership-v1-RemoveOrganizationAdminDeputyRequest"></a>

### RemoveOrganizationAdminDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RemoveOrganizationAdminDeputyResponse"></a>

### RemoveOrganizationAdminDeputyResponse







<a name="command-membership-v1-RemoveOrganizationDispatcherDeputyRequest"></a>

### RemoveOrganizationDispatcherDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RemoveOrganizationDispatcherDeputyResponse"></a>

### RemoveOrganizationDispatcherDeputyResponse







<a name="command-membership-v1-RemoveOrganizationHeadDeputyRequest"></a>

### RemoveOrganizationHeadDeputyRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RemoveOrganizationHeadDeputyResponse"></a>

### RemoveOrganizationHeadDeputyResponse







<a name="command-membership-v1-RevokeClinicHeadRequest"></a>

### RevokeClinicHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RevokeClinicHeadResponse"></a>

### RevokeClinicHeadResponse







<a name="command-membership-v1-RevokeDepartmentResponsibleRequest"></a>

### RevokeDepartmentResponsibleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RevokeDepartmentResponsibleResponse"></a>

### RevokeDepartmentResponsibleResponse







<a name="command-membership-v1-RevokeOrganizationAdminRequest"></a>

### RevokeOrganizationAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RevokeOrganizationAdminResponse"></a>

### RevokeOrganizationAdminResponse







<a name="command-membership-v1-RevokeOrganizationDispatcherRequest"></a>

### RevokeOrganizationDispatcherRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RevokeOrganizationDispatcherResponse"></a>

### RevokeOrganizationDispatcherResponse







<a name="command-membership-v1-RevokeOrganizationHeadRequest"></a>

### RevokeOrganizationHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-RevokeOrganizationHeadResponse"></a>

### RevokeOrganizationHeadResponse







<a name="command-membership-v1-RevokeSystemAdminRequest"></a>

### RevokeSystemAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |






<a name="command-membership-v1-RevokeSystemAdminResponse"></a>

### RevokeSystemAdminResponse







<a name="command-membership-v1-ScheduleVacationRequest"></a>

### ScheduleVacationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| starts_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |






<a name="command-membership-v1-ScheduleVacationResponse"></a>

### ScheduleVacationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="command-membership-v1-StartVacationNowRequest"></a>

### StartVacationNowRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) | optional |  |






<a name="command-membership-v1-StartVacationNowResponse"></a>

### StartVacationNowResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |






<a name="command-membership-v1-TerminateEmployeeRequest"></a>

### TerminateEmployeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="command-membership-v1-TerminateEmployeeResponse"></a>

### TerminateEmployeeResponse







<a name="command-membership-v1-UpdateEmployeeDepartmentRequest"></a>

### UpdateEmployeeDepartmentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |






<a name="command-membership-v1-UpdateEmployeeDepartmentResponse"></a>

### UpdateEmployeeDepartmentResponse







<a name="command-membership-v1-UpdateEmployeePositionRequest"></a>

### UpdateEmployeePositionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| position | [string](#string) | optional |  |






<a name="command-membership-v1-UpdateEmployeePositionResponse"></a>

### UpdateEmployeePositionResponse







<a name="command-membership-v1-UpdateVacationEndDateRequest"></a>

### UpdateVacationEndDateRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="command-membership-v1-UpdateVacationEndDateResponse"></a>

### UpdateVacationEndDateResponse













<a name="command-membership-v1-MembershipCommandService"></a>

### MembershipCommandService
MembershipService is the command-side contract for who-works-where
data. It covers the Employee lifecycle (hire / update / terminate),
Vacation lifecycle (start / schedule / end / cancel / change),
role grants with deputy management for DepartmentResponsible,
ClinicHead, OrganizationAdmin, OrganizationHead, and
OrganizationDispatcher, plus the global SystemAdmin grant that
operates directly on Zitadel user identifiers.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| HireEmployee | [HireEmployeeRequest](#command-membership-v1-HireEmployeeRequest) | [HireEmployeeResponse](#command-membership-v1-HireEmployeeResponse) | Employee lifecycle |
| UpdateEmployeePosition | [UpdateEmployeePositionRequest](#command-membership-v1-UpdateEmployeePositionRequest) | [UpdateEmployeePositionResponse](#command-membership-v1-UpdateEmployeePositionResponse) |  |
| UpdateEmployeeDepartment | [UpdateEmployeeDepartmentRequest](#command-membership-v1-UpdateEmployeeDepartmentRequest) | [UpdateEmployeeDepartmentResponse](#command-membership-v1-UpdateEmployeeDepartmentResponse) |  |
| TerminateEmployee | [TerminateEmployeeRequest](#command-membership-v1-TerminateEmployeeRequest) | [TerminateEmployeeResponse](#command-membership-v1-TerminateEmployeeResponse) |  |
| StartVacationNow | [StartVacationNowRequest](#command-membership-v1-StartVacationNowRequest) | [StartVacationNowResponse](#command-membership-v1-StartVacationNowResponse) | Vacation lifecycle |
| ScheduleVacation | [ScheduleVacationRequest](#command-membership-v1-ScheduleVacationRequest) | [ScheduleVacationResponse](#command-membership-v1-ScheduleVacationResponse) |  |
| UpdateVacationEndDate | [UpdateVacationEndDateRequest](#command-membership-v1-UpdateVacationEndDateRequest) | [UpdateVacationEndDateResponse](#command-membership-v1-UpdateVacationEndDateResponse) |  |
| ForceEndVacation | [ForceEndVacationRequest](#command-membership-v1-ForceEndVacationRequest) | [ForceEndVacationResponse](#command-membership-v1-ForceEndVacationResponse) |  |
| CancelScheduledVacation | [CancelScheduledVacationRequest](#command-membership-v1-CancelScheduledVacationRequest) | [CancelScheduledVacationResponse](#command-membership-v1-CancelScheduledVacationResponse) |  |
| AssignDepartmentResponsible | [AssignDepartmentResponsibleRequest](#command-membership-v1-AssignDepartmentResponsibleRequest) | [AssignDepartmentResponsibleResponse](#command-membership-v1-AssignDepartmentResponsibleResponse) | ------ DepartmentResponsible ------ |
| RevokeDepartmentResponsible | [RevokeDepartmentResponsibleRequest](#command-membership-v1-RevokeDepartmentResponsibleRequest) | [RevokeDepartmentResponsibleResponse](#command-membership-v1-RevokeDepartmentResponsibleResponse) |  |
| AssignDepartmentResponsibleDeputy | [AssignDepartmentResponsibleDeputyRequest](#command-membership-v1-AssignDepartmentResponsibleDeputyRequest) | [AssignDepartmentResponsibleDeputyResponse](#command-membership-v1-AssignDepartmentResponsibleDeputyResponse) |  |
| RemoveDepartmentResponsibleDeputy | [RemoveDepartmentResponsibleDeputyRequest](#command-membership-v1-RemoveDepartmentResponsibleDeputyRequest) | [RemoveDepartmentResponsibleDeputyResponse](#command-membership-v1-RemoveDepartmentResponsibleDeputyResponse) |  |
| AssignClinicHead | [AssignClinicHeadRequest](#command-membership-v1-AssignClinicHeadRequest) | [AssignClinicHeadResponse](#command-membership-v1-AssignClinicHeadResponse) | ------ ClinicHead ------ |
| RevokeClinicHead | [RevokeClinicHeadRequest](#command-membership-v1-RevokeClinicHeadRequest) | [RevokeClinicHeadResponse](#command-membership-v1-RevokeClinicHeadResponse) |  |
| AssignClinicHeadDeputy | [AssignClinicHeadDeputyRequest](#command-membership-v1-AssignClinicHeadDeputyRequest) | [AssignClinicHeadDeputyResponse](#command-membership-v1-AssignClinicHeadDeputyResponse) |  |
| RemoveClinicHeadDeputy | [RemoveClinicHeadDeputyRequest](#command-membership-v1-RemoveClinicHeadDeputyRequest) | [RemoveClinicHeadDeputyResponse](#command-membership-v1-RemoveClinicHeadDeputyResponse) |  |
| AssignOrganizationAdmin | [AssignOrganizationAdminRequest](#command-membership-v1-AssignOrganizationAdminRequest) | [AssignOrganizationAdminResponse](#command-membership-v1-AssignOrganizationAdminResponse) | ------ OrganizationAdmin ------ |
| RevokeOrganizationAdmin | [RevokeOrganizationAdminRequest](#command-membership-v1-RevokeOrganizationAdminRequest) | [RevokeOrganizationAdminResponse](#command-membership-v1-RevokeOrganizationAdminResponse) |  |
| AssignOrganizationAdminDeputy | [AssignOrganizationAdminDeputyRequest](#command-membership-v1-AssignOrganizationAdminDeputyRequest) | [AssignOrganizationAdminDeputyResponse](#command-membership-v1-AssignOrganizationAdminDeputyResponse) |  |
| RemoveOrganizationAdminDeputy | [RemoveOrganizationAdminDeputyRequest](#command-membership-v1-RemoveOrganizationAdminDeputyRequest) | [RemoveOrganizationAdminDeputyResponse](#command-membership-v1-RemoveOrganizationAdminDeputyResponse) |  |
| AssignOrganizationHead | [AssignOrganizationHeadRequest](#command-membership-v1-AssignOrganizationHeadRequest) | [AssignOrganizationHeadResponse](#command-membership-v1-AssignOrganizationHeadResponse) | ------ OrganizationHead ------ |
| RevokeOrganizationHead | [RevokeOrganizationHeadRequest](#command-membership-v1-RevokeOrganizationHeadRequest) | [RevokeOrganizationHeadResponse](#command-membership-v1-RevokeOrganizationHeadResponse) |  |
| AssignOrganizationHeadDeputy | [AssignOrganizationHeadDeputyRequest](#command-membership-v1-AssignOrganizationHeadDeputyRequest) | [AssignOrganizationHeadDeputyResponse](#command-membership-v1-AssignOrganizationHeadDeputyResponse) |  |
| RemoveOrganizationHeadDeputy | [RemoveOrganizationHeadDeputyRequest](#command-membership-v1-RemoveOrganizationHeadDeputyRequest) | [RemoveOrganizationHeadDeputyResponse](#command-membership-v1-RemoveOrganizationHeadDeputyResponse) |  |
| AssignOrganizationDispatcher | [AssignOrganizationDispatcherRequest](#command-membership-v1-AssignOrganizationDispatcherRequest) | [AssignOrganizationDispatcherResponse](#command-membership-v1-AssignOrganizationDispatcherResponse) | ------ OrganizationDispatcher ------ |
| RevokeOrganizationDispatcher | [RevokeOrganizationDispatcherRequest](#command-membership-v1-RevokeOrganizationDispatcherRequest) | [RevokeOrganizationDispatcherResponse](#command-membership-v1-RevokeOrganizationDispatcherResponse) |  |
| AssignOrganizationDispatcherDeputy | [AssignOrganizationDispatcherDeputyRequest](#command-membership-v1-AssignOrganizationDispatcherDeputyRequest) | [AssignOrganizationDispatcherDeputyResponse](#command-membership-v1-AssignOrganizationDispatcherDeputyResponse) |  |
| RemoveOrganizationDispatcherDeputy | [RemoveOrganizationDispatcherDeputyRequest](#command-membership-v1-RemoveOrganizationDispatcherDeputyRequest) | [RemoveOrganizationDispatcherDeputyResponse](#command-membership-v1-RemoveOrganizationDispatcherDeputyResponse) |  |
| GrantSystemAdmin | [GrantSystemAdminRequest](#command-membership-v1-GrantSystemAdminRequest) | [GrantSystemAdminResponse](#command-membership-v1-GrantSystemAdminResponse) | ------ SystemAdmin ------ |
| RevokeSystemAdmin | [RevokeSystemAdminRequest](#command-membership-v1-RevokeSystemAdminRequest) | [RevokeSystemAdminResponse](#command-membership-v1-RevokeSystemAdminResponse) |  |





<a name="command_orgstructure_v1_orgstructure-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/orgstructure/v1/orgstructure.proto



<a name="command-orgstructure-v1-AddressInput"></a>

### AddressInput



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  |  |
| point | [PointInput](#command-orgstructure-v1-PointInput) | optional |  |






<a name="command-orgstructure-v1-CreateClinicRequest"></a>

### CreateClinicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| physical_address | [AddressInput](#command-orgstructure-v1-AddressInput) |  |  |
| description | [string](#string) | optional |  |






<a name="command-orgstructure-v1-CreateClinicResponse"></a>

### CreateClinicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |






<a name="command-orgstructure-v1-CreateDepartmentRequest"></a>

### CreateDepartmentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-orgstructure-v1-CreateDepartmentResponse"></a>

### CreateDepartmentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |






<a name="command-orgstructure-v1-CreateOrganizationRequest"></a>

### CreateOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| legal_address | [AddressInput](#command-orgstructure-v1-AddressInput) |  |  |
| description | [string](#string) | optional |  |






<a name="command-orgstructure-v1-CreateOrganizationResponse"></a>

### CreateOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |






<a name="command-orgstructure-v1-PointInput"></a>

### PointInput



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| longitude | [double](#double) |  |  |
| latitude | [double](#double) |  |  |






<a name="command-orgstructure-v1-UpdateClinicDetailsRequest"></a>

### UpdateClinicDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-orgstructure-v1-UpdateClinicDetailsResponse"></a>

### UpdateClinicDetailsResponse







<a name="command-orgstructure-v1-UpdateClinicPhysicalAddressRequest"></a>

### UpdateClinicPhysicalAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| physical_address | [AddressInput](#command-orgstructure-v1-AddressInput) |  |  |






<a name="command-orgstructure-v1-UpdateClinicPhysicalAddressResponse"></a>

### UpdateClinicPhysicalAddressResponse







<a name="command-orgstructure-v1-UpdateDepartmentDetailsRequest"></a>

### UpdateDepartmentDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-orgstructure-v1-UpdateDepartmentDetailsResponse"></a>

### UpdateDepartmentDetailsResponse







<a name="command-orgstructure-v1-UpdateOrganizationDetailsRequest"></a>

### UpdateOrganizationDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-orgstructure-v1-UpdateOrganizationDetailsResponse"></a>

### UpdateOrganizationDetailsResponse







<a name="command-orgstructure-v1-UpdateOrganizationLegalAddressRequest"></a>

### UpdateOrganizationLegalAddressRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| legal_address | [AddressInput](#command-orgstructure-v1-AddressInput) |  |  |






<a name="command-orgstructure-v1-UpdateOrganizationLegalAddressResponse"></a>

### UpdateOrganizationLegalAddressResponse













<a name="command-orgstructure-v1-OrgStructureCommandService"></a>

### OrgStructureCommandService
OrgStructureCommandService is the command-side contract for the three
organisational structure aggregates: Organization, Clinic, and
Department. Every mutation returns either an identifier (Create) or
an empty response (Update). google.api.http annotations drive the
gateway-server binary in this repo (cmd/gateway-server): the
command-server itself serves pure gRPC, and gateway-server dials
into it and exposes a REST facade using the grpc-gateway stubs
generated under pkg/.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateOrganization | [CreateOrganizationRequest](#command-orgstructure-v1-CreateOrganizationRequest) | [CreateOrganizationResponse](#command-orgstructure-v1-CreateOrganizationResponse) |  |
| UpdateOrganizationDetails | [UpdateOrganizationDetailsRequest](#command-orgstructure-v1-UpdateOrganizationDetailsRequest) | [UpdateOrganizationDetailsResponse](#command-orgstructure-v1-UpdateOrganizationDetailsResponse) |  |
| UpdateOrganizationLegalAddress | [UpdateOrganizationLegalAddressRequest](#command-orgstructure-v1-UpdateOrganizationLegalAddressRequest) | [UpdateOrganizationLegalAddressResponse](#command-orgstructure-v1-UpdateOrganizationLegalAddressResponse) |  |
| CreateClinic | [CreateClinicRequest](#command-orgstructure-v1-CreateClinicRequest) | [CreateClinicResponse](#command-orgstructure-v1-CreateClinicResponse) |  |
| UpdateClinicDetails | [UpdateClinicDetailsRequest](#command-orgstructure-v1-UpdateClinicDetailsRequest) | [UpdateClinicDetailsResponse](#command-orgstructure-v1-UpdateClinicDetailsResponse) |  |
| UpdateClinicPhysicalAddress | [UpdateClinicPhysicalAddressRequest](#command-orgstructure-v1-UpdateClinicPhysicalAddressRequest) | [UpdateClinicPhysicalAddressResponse](#command-orgstructure-v1-UpdateClinicPhysicalAddressResponse) |  |
| CreateDepartment | [CreateDepartmentRequest](#command-orgstructure-v1-CreateDepartmentRequest) | [CreateDepartmentResponse](#command-orgstructure-v1-CreateDepartmentResponse) |  |
| UpdateDepartmentDetails | [UpdateDepartmentDetailsRequest](#command-orgstructure-v1-UpdateDepartmentDetailsRequest) | [UpdateDepartmentDetailsResponse](#command-orgstructure-v1-UpdateDepartmentDetailsResponse) |  |





<a name="command_request_classifier_v1_request_classifier-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/request/classifier/v1/request_classifier.proto



<a name="command-request-classifier-v1-CreateRequestTypeRequest"></a>

### CreateRequestTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-request-classifier-v1-CreateRequestTypeResponse"></a>

### CreateRequestTypeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-request-classifier-v1-DeactivateRequestTypeRequest"></a>

### DeactivateRequestTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-request-classifier-v1-DeactivateRequestTypeResponse"></a>

### DeactivateRequestTypeResponse







<a name="command-request-classifier-v1-DeleteRequestTypeRequest"></a>

### DeleteRequestTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-request-classifier-v1-DeleteRequestTypeResponse"></a>

### DeleteRequestTypeResponse







<a name="command-request-classifier-v1-ReactivateRequestTypeRequest"></a>

### ReactivateRequestTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |






<a name="command-request-classifier-v1-ReactivateRequestTypeResponse"></a>

### ReactivateRequestTypeResponse







<a name="command-request-classifier-v1-UpdateRequestTypeDetailsRequest"></a>

### UpdateRequestTypeDetailsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |






<a name="command-request-classifier-v1-UpdateRequestTypeDetailsResponse"></a>

### UpdateRequestTypeDetailsResponse













<a name="command-request-classifier-v1-RequestClassifierCommandService"></a>

### RequestClassifierCommandService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateRequestType | [CreateRequestTypeRequest](#command-request-classifier-v1-CreateRequestTypeRequest) | [CreateRequestTypeResponse](#command-request-classifier-v1-CreateRequestTypeResponse) |  |
| UpdateRequestTypeDetails | [UpdateRequestTypeDetailsRequest](#command-request-classifier-v1-UpdateRequestTypeDetailsRequest) | [UpdateRequestTypeDetailsResponse](#command-request-classifier-v1-UpdateRequestTypeDetailsResponse) |  |
| DeactivateRequestType | [DeactivateRequestTypeRequest](#command-request-classifier-v1-DeactivateRequestTypeRequest) | [DeactivateRequestTypeResponse](#command-request-classifier-v1-DeactivateRequestTypeResponse) |  |
| ReactivateRequestType | [ReactivateRequestTypeRequest](#command-request-classifier-v1-ReactivateRequestTypeRequest) | [ReactivateRequestTypeResponse](#command-request-classifier-v1-ReactivateRequestTypeResponse) |  |
| DeleteRequestType | [DeleteRequestTypeRequest](#command-request-classifier-v1-DeleteRequestTypeRequest) | [DeleteRequestTypeResponse](#command-request-classifier-v1-DeleteRequestTypeResponse) |  |





<a name="command_request_v1_request-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## command/request/v1/request.proto



<a name="command-request-v1-AssignExecutorsRequest"></a>

### AssignExecutorsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_request_id | [string](#string) |  |  |
| executor_employee_ids | [string](#string) | repeated |  |






<a name="command-request-v1-AssignExecutorsResponse"></a>

### AssignExecutorsResponse







<a name="command-request-v1-CreateServiceRequestRequest"></a>

### CreateServiceRequestRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| type_id | [string](#string) |  |  |
| incident_id | [string](#string) | optional |  |
| description | [string](#string) |  |  |
| executor_employee_ids | [string](#string) | repeated |  |






<a name="command-request-v1-CreateServiceRequestResponse"></a>

### CreateServiceRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_request_id | [string](#string) |  |  |






<a name="command-request-v1-UpdateServiceRequestDescriptionRequest"></a>

### UpdateServiceRequestDescriptionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_request_id | [string](#string) |  |  |
| description | [string](#string) |  |  |






<a name="command-request-v1-UpdateServiceRequestDescriptionResponse"></a>

### UpdateServiceRequestDescriptionResponse







<a name="command-request-v1-UpdateServiceRequestStatusRequest"></a>

### UpdateServiceRequestStatusRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_request_id | [string](#string) |  |  |
| new_status | [string](#string) |  |  |






<a name="command-request-v1-UpdateServiceRequestStatusResponse"></a>

### UpdateServiceRequestStatusResponse













<a name="command-request-v1-ServiceRequestCommandService"></a>

### ServiceRequestCommandService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| CreateServiceRequest | [CreateServiceRequestRequest](#command-request-v1-CreateServiceRequestRequest) | [CreateServiceRequestResponse](#command-request-v1-CreateServiceRequestResponse) |  |
| UpdateServiceRequestDescription | [UpdateServiceRequestDescriptionRequest](#command-request-v1-UpdateServiceRequestDescriptionRequest) | [UpdateServiceRequestDescriptionResponse](#command-request-v1-UpdateServiceRequestDescriptionResponse) |  |
| UpdateServiceRequestStatus | [UpdateServiceRequestStatusRequest](#command-request-v1-UpdateServiceRequestStatusRequest) | [UpdateServiceRequestStatusResponse](#command-request-v1-UpdateServiceRequestStatusResponse) |  |
| AssignExecutors | [AssignExecutorsRequest](#command-request-v1-AssignExecutorsRequest) | [AssignExecutorsResponse](#command-request-v1-AssignExecutorsResponse) |  |





<a name="event_clinic_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/clinic/v1/events.proto



<a name="event-clinic-v1-Address"></a>

### Address



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  |  |
| point | [Point](#event-clinic-v1-Point) |  |  |






<a name="event-clinic-v1-ClinicCreated"></a>

### ClinicCreated
ClinicCreated is emitted when a new clinic is persisted.
subject: medincident.event.clinic.v1.created


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| physical_address | [Address](#event-clinic-v1-Address) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-clinic-v1-ClinicDetailsChanged"></a>

### ClinicDetailsChanged
ClinicDetailsChanged is emitted when name or description changes.
subject: medincident.event.clinic.v1.details_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-clinic-v1-ClinicHeadAssigned"></a>

### ClinicHeadAssigned
ClinicHeadAssigned — subject: medincident.event.clinic.v1.clinic_head_assigned


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| assigned_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-clinic-v1-ClinicHeadDeputyAssigned"></a>

### ClinicHeadDeputyAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-clinic-v1-ClinicHeadDeputyRemoved"></a>

### ClinicHeadDeputyRemoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-clinic-v1-ClinicHeadRevoked"></a>

### ClinicHeadRevoked



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-clinic-v1-ClinicPhysicalAddressChanged"></a>

### ClinicPhysicalAddressChanged
ClinicPhysicalAddressChanged is emitted when the physical address changes.
subject: medincident.event.clinic.v1.physical_address_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| physical_address | [Address](#event-clinic-v1-Address) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






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
DepartmentCreated is emitted when a new department is persisted.
subject: medincident.event.department.v1.created


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-department-v1-DepartmentDetailsChanged"></a>

### DepartmentDetailsChanged
DepartmentDetailsChanged is emitted when name or description changes.
subject: medincident.event.department.v1.details_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-department-v1-DeptResponsibleAssigned"></a>

### DeptResponsibleAssigned
DeptResponsibleAssigned — subject: medincident.event.department.v1.dept_responsible_assigned


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| assigned_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-department-v1-DeptResponsibleDeputyAssigned"></a>

### DeptResponsibleDeputyAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-department-v1-DeptResponsibleDeputyRemoved"></a>

### DeptResponsibleDeputyRemoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-department-v1-DeptResponsibleRevoked"></a>

### DeptResponsibleRevoked



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |















<a name="event_employee_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/employee/v1/events.proto



<a name="event-employee-v1-EmployeeDepartmentChanged"></a>

### EmployeeDepartmentChanged
EmployeeDepartmentChanged — subject: medincident.event.employee.v1.department_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| new_department_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-employee-v1-EmployeeHired"></a>

### EmployeeHired
EmployeeHired — subject: medincident.event.employee.v1.hired
aggregate_id = employee UUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| position | [string](#string) |  | empty = not set |
| hired_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-employee-v1-EmployeePositionChanged"></a>

### EmployeePositionChanged
EmployeePositionChanged — subject: medincident.event.employee.v1.position_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| position | [string](#string) |  | empty = cleared |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-employee-v1-EmployeeTerminated"></a>

### EmployeeTerminated
EmployeeTerminated — subject: medincident.event.employee.v1.terminated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  | needed for counter decrement |
| department_id | [string](#string) |  |  |
| terminated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |















<a name="event_incident_buffer_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/incident/buffer/v1/events.proto



<a name="event-incident-buffer-v1-PatientIncidentBufferCreated"></a>

### PatientIncidentBufferCreated
PatientIncidentBufferCreated — subject: medincident.event.patient_incident_buffer.v1.created
aggregate_id = buffer UUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| buffer_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| patient_zitadel_user_id | [string](#string) |  |  |
| category_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| type_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| description | [string](#string) |  |  |
| occurred_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| status | [string](#string) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-buffer-v1-PatientIncidentBufferUpdated"></a>

### PatientIncidentBufferUpdated
PatientIncidentBufferUpdated — subject: medincident.event.patient_incident_buffer.v1.updated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| buffer_id | [string](#string) |  |  |
| category_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| type_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| description | [string](#string) |  |  |
| occurred_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| status | [string](#string) |  |  |
| published_incident_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |















<a name="event_incident_classifier_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/incident/classifier/v1/events.proto



<a name="event-incident-classifier-v1-IncidentCategoryCreated"></a>

### IncidentCategoryCreated
IncidentCategoryCreated — subject: medincident.event.incident_category.v1.created
aggregate_id = category UUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| parent_category_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  | absent = root category |
| name | [string](#string) |  |  |
| description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  | absent = no description |
| is_active | [bool](#bool) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentCategoryDeactivated"></a>

### IncidentCategoryDeactivated
IncidentCategoryDeactivated — subject: medincident.event.incident_category.v1.deactivated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentCategoryDeleted"></a>

### IncidentCategoryDeleted
IncidentCategoryDeleted — subject: medincident.event.incident_category.v1.deleted


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| deleted_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentCategoryDetailsUpdated"></a>

### IncidentCategoryDetailsUpdated
IncidentCategoryDetailsUpdated — subject: medincident.event.incident_category.v1.details_updated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  | absent = description cleared |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentCategoryMoved"></a>

### IncidentCategoryMoved
IncidentCategoryMoved — subject: medincident.event.incident_category.v1.moved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| new_parent_category_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  | absent = moved to root |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentCategoryReactivated"></a>

### IncidentCategoryReactivated
IncidentCategoryReactivated — subject: medincident.event.incident_category.v1.reactivated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentTypeAllowedForPatients"></a>

### IncidentTypeAllowedForPatients
IncidentTypeAllowedForPatients — subject: medincident.event.incident_type.v1.allowed_for_patients


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentTypeCreated"></a>

### IncidentTypeCreated
IncidentTypeCreated — subject: medincident.event.incident_type.v1.created
aggregate_id = type UUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| category_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  | absent = no description |
| is_active | [bool](#bool) |  |  |
| is_allowed_for_patients | [bool](#bool) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentTypeDeactivated"></a>

### IncidentTypeDeactivated
IncidentTypeDeactivated — subject: medincident.event.incident_type.v1.deactivated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentTypeDeleted"></a>

### IncidentTypeDeleted
IncidentTypeDeleted — subject: medincident.event.incident_type.v1.deleted


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| deleted_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentTypeDetailsUpdated"></a>

### IncidentTypeDetailsUpdated
IncidentTypeDetailsUpdated — subject: medincident.event.incident_type.v1.details_updated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  | absent = description cleared |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentTypeDisallowedForPatients"></a>

### IncidentTypeDisallowedForPatients
IncidentTypeDisallowedForPatients — subject: medincident.event.incident_type.v1.disallowed_for_patients


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentTypeMoved"></a>

### IncidentTypeMoved
IncidentTypeMoved — subject: medincident.event.incident_type.v1.moved


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| new_category_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-classifier-v1-IncidentTypeReactivated"></a>

### IncidentTypeReactivated
IncidentTypeReactivated — subject: medincident.event.incident_type.v1.reactivated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |















<a name="event_incident_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/incident/v1/events.proto



<a name="event-incident-v1-IncidentCreated"></a>

### IncidentCreated
IncidentCreated — subject: medincident.event.incident.v1.created
aggregate_id = incident UUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| clinic_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| category_id | [string](#string) |  |  |
| type_id | [string](#string) |  |  |
| status | [string](#string) |  |  |
| priority | [string](#string) |  |  |
| description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| patient_original_description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| occurred_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| registrar_zitadel_user_id | [string](#string) |  |  |
| registrar_employee_id | [string](#string) |  |  |
| source_patient_zitadel_user_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| source_buffer_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| reopened_from_incident_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-v1-IncidentDescriptionUpdated"></a>

### IncidentDescriptionUpdated
IncidentDescriptionUpdated — subject: medincident.event.incident.v1.description_updated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |
| description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-v1-IncidentPriorityChanged"></a>

### IncidentPriorityChanged
IncidentPriorityChanged — subject: medincident.event.incident.v1.priority_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |
| old_priority | [string](#string) |  |  |
| new_priority | [string](#string) |  |  |
| actor_employee_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| actor_zitadel_user_id | [string](#string) |  |  |
| changed_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-incident-v1-IncidentStatusChanged"></a>

### IncidentStatusChanged
IncidentStatusChanged — subject: medincident.event.incident.v1.status_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |
| old_status | [string](#string) |  |  |
| new_status | [string](#string) |  |  |
| actor_employee_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| actor_zitadel_user_id | [string](#string) |  |  |
| changed_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |















<a name="event_organization_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/organization/v1/events.proto



<a name="event-organization-v1-Address"></a>

### Address
Address is embedded in organization events that carry a legal address.
Duplicated per-aggregate so aggregates can evolve independently.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  |  |
| point | [Point](#event-organization-v1-Point) |  |  |






<a name="event-organization-v1-OrgAdminAssigned"></a>

### OrgAdminAssigned
OrgAdminAssigned — subject: medincident.event.organization.v1.org_admin_assigned


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| assigned_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgAdminDeputyAssigned"></a>

### OrgAdminDeputyAssigned
OrgAdminDeputyAssigned — subject: medincident.event.organization.v1.org_admin_deputy_assigned


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgAdminDeputyRemoved"></a>

### OrgAdminDeputyRemoved
OrgAdminDeputyRemoved — subject: medincident.event.organization.v1.org_admin_deputy_removed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgAdminRevoked"></a>

### OrgAdminRevoked
OrgAdminRevoked — subject: medincident.event.organization.v1.org_admin_revoked


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrgDispatcherAssigned"></a>

### OrgDispatcherAssigned
OrgDispatcherAssigned — subject: medincident.event.organization.v1.org_dispatcher_assigned


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| assigned_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgDispatcherDeputyAssigned"></a>

### OrgDispatcherDeputyAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgDispatcherDeputyRemoved"></a>

### OrgDispatcherDeputyRemoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgDispatcherRevoked"></a>

### OrgDispatcherRevoked



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrgHeadAssigned"></a>

### OrgHeadAssigned
OrgHeadAssigned — subject: medincident.event.organization.v1.org_head_assigned


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| assigned_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgHeadDeputyAssigned"></a>

### OrgHeadDeputyAssigned



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgHeadDeputyRemoved"></a>

### OrgHeadDeputyRemoved



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrgHeadRevoked"></a>

### OrgHeadRevoked



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |






<a name="event-organization-v1-OrganizationCreated"></a>

### OrganizationCreated
OrganizationCreated is emitted when a new organization is persisted.
subject: medincident.event.organization.v1.created


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) |  | empty = absent |
| legal_address | [Address](#event-organization-v1-Address) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrganizationDetailsChanged"></a>

### OrganizationDetailsChanged
OrganizationDetailsChanged is emitted when name or description changes.
subject: medincident.event.organization.v1.details_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) |  | empty = absent |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-OrganizationLegalAddressChanged"></a>

### OrganizationLegalAddressChanged
OrganizationLegalAddressChanged is emitted when the legal address changes.
subject: medincident.event.organization.v1.legal_address_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| legal_address | [Address](#event-organization-v1-Address) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-organization-v1-Point"></a>

### Point



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| longitude | [double](#double) |  |  |
| latitude | [double](#double) |  |  |















<a name="event_request_type_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/request_type/v1/events.proto



<a name="event-request_type-v1-RequestTypeCreated"></a>

### RequestTypeCreated
RequestTypeCreated — subject: medincident.event.request_type.v1.created
aggregate_id = request_type UUID


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| is_active | [bool](#bool) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-request_type-v1-RequestTypeDeactivated"></a>

### RequestTypeDeactivated
RequestTypeDeactivated — subject: medincident.event.request_type.v1.deactivated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-request_type-v1-RequestTypeDeleted"></a>

### RequestTypeDeleted
RequestTypeDeleted — subject: medincident.event.request_type.v1.deleted


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| deleted_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-request_type-v1-RequestTypeDetailsUpdated"></a>

### RequestTypeDetailsUpdated
RequestTypeDetailsUpdated — subject: medincident.event.request_type.v1.details_updated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-request_type-v1-RequestTypeReactivated"></a>

### RequestTypeReactivated
RequestTypeReactivated — subject: medincident.event.request_type.v1.reactivated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |















<a name="event_service_request_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/service_request/v1/events.proto



<a name="event-service_request-v1-ServiceRequestCreated"></a>

### ServiceRequestCreated
ServiceRequestCreated — subject: medincident.event.service_request.v1.created
aggregate_id = service_request UUID
author_display_name is NOT included — query-side projector calls lookupUserDisplayName(tx, author_id).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| clinic_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| type_id | [string](#string) |  |  |
| incident_id | [google.protobuf.StringValue](#google-protobuf-StringValue) |  |  |
| description | [string](#string) |  |  |
| status | [string](#string) |  |  |
| author_id | [string](#string) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-service_request-v1-ServiceRequestDescriptionUpdated"></a>

### ServiceRequestDescriptionUpdated
ServiceRequestDescriptionUpdated — subject: medincident.event.service_request.v1.description_updated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  |  |
| description | [string](#string) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-service_request-v1-ServiceRequestExecutorAssigned"></a>

### ServiceRequestExecutorAssigned
ServiceRequestExecutorAssigned — subject: medincident.event.service_request.v1.executor_assigned


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| actor_id | [string](#string) |  |  |
| changed_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-service_request-v1-ServiceRequestExecutorRemoved"></a>

### ServiceRequestExecutorRemoved
ServiceRequestExecutorRemoved — subject: medincident.event.service_request.v1.executor_removed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| actor_id | [string](#string) |  |  |
| changed_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-service_request-v1-ServiceRequestStatusChanged"></a>

### ServiceRequestStatusChanged
ServiceRequestStatusChanged — subject: medincident.event.service_request.v1.status_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_id | [string](#string) |  |  |
| old_status | [string](#string) |  |  |
| new_status | [string](#string) |  |  |
| actor_id | [string](#string) |  |  |
| changed_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |















<a name="event_system_admin_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/system_admin/v1/events.proto



<a name="event-system_admin-v1-SystemAdminGranted"></a>

### SystemAdminGranted
SystemAdminGranted — subject: medincident.event.system_admin.v1.granted


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| granted_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-system_admin-v1-SystemAdminRevoked"></a>

### SystemAdminRevoked
SystemAdminRevoked — subject: medincident.event.system_admin.v1.revoked


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| revoked_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |















<a name="event_v1_envelope-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/v1/envelope.proto



<a name="event-v1-Envelope"></a>

### Envelope
Envelope is the transport wrapper for every domain event published to
NATS JetStream. The subject on the NATS message encodes the routing
key (medincident.event.&lt;aggregate_type&gt;.v1.&lt;action&gt;); fields here
carry the structured metadata and the typed payload.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| occurred_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | occurred_at is the wall-clock time the domain mutation happened. |
| aggregate_type | [string](#string) |  | aggregate_type identifies the aggregate family: &#34;organization&#34;, &#34;clinic&#34;, &#34;department&#34;, &#34;employee&#34;, &#34;vacation&#34;, &#34;system_admin&#34;, &#34;incident_category&#34;, &#34;incident_type&#34;, &#34;incident&#34;, &#34;incident_buffer&#34;, &#34;request_type&#34;, &#34;service_request&#34;. |
| aggregate_id | [string](#string) |  | aggregate_id is the UUID (as a string) of the aggregate instance. For system_admin events it is the zitadel_user_id. |
| payload | [google.protobuf.Any](#google-protobuf-Any) |  | payload holds the concrete event message serialised as Any. |















<a name="event_vacation_v1_events-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## event/vacation/v1/events.proto



<a name="event-vacation-v1-VacationCancelled"></a>

### VacationCancelled
VacationCancelled — subject: medincident.event.vacation.v1.cancelled


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| cancelled_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-vacation-v1-VacationEndDateChanged"></a>

### VacationEndDateChanged
VacationEndDateChanged — subject: medincident.event.vacation.v1.end_date_changed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| updated_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-vacation-v1-VacationEnded"></a>

### VacationEnded
VacationEnded — subject: medincident.event.vacation.v1.ended


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | final ends_at stamp |






<a name="event-vacation-v1-VacationScheduled"></a>

### VacationScheduled
VacationScheduled — subject: medincident.event.vacation.v1.scheduled


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| starts_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  | absent = open-ended |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |






<a name="event-vacation-v1-VacationStarted"></a>

### VacationStarted
VacationStarted — subject: medincident.event.vacation.v1.started


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| vacation_id | [string](#string) |  |  |
| starts_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| ends_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |
| created_at | [google.protobuf.Timestamp](#google-protobuf-Timestamp) |  |  |















<a name="query_analytics_v1_analytics-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/analytics/v1/analytics.proto
api/proto/query/analytics/v1/analytics.proto


<a name="query-analytics-v1-CategoryCount"></a>

### CategoryCount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| category_name | [string](#string) |  |  |
| count | [int64](#int64) |  |  |






<a name="query-analytics-v1-DepartmentCount"></a>

### DepartmentCount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| department_name | [string](#string) |  |  |
| count | [int64](#int64) |  |  |






<a name="query-analytics-v1-GetSnapshotRequest"></a>

### GetSnapshotRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| from | [string](#string) |  |  |
| to | [string](#string) |  |  |
| clinic_id | [string](#string) | optional |  |
| department_id | [string](#string) | optional |  |
| include_patient_buffer | [bool](#bool) |  |  |






<a name="query-analytics-v1-GetSnapshotResponse"></a>

### GetSnapshotResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incidents | [SnapshotIncident](#query-analytics-v1-SnapshotIncident) | repeated |  |
| requests | [SnapshotRequest](#query-analytics-v1-SnapshotRequest) | repeated |  |
| patient_buffer | [SnapshotPatientBuffer](#query-analytics-v1-SnapshotPatientBuffer) | repeated |  |






<a name="query-analytics-v1-GetSummaryRequest"></a>

### GetSummaryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| from | [string](#string) |  |  |
| to | [string](#string) |  |  |
| clinic_id | [string](#string) | optional |  |
| department_id | [string](#string) | optional |  |






<a name="query-analytics-v1-GetSummaryResponse"></a>

### GetSummaryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incidents | [IncidentSummary](#query-analytics-v1-IncidentSummary) |  |  |
| requests | [RequestSummary](#query-analytics-v1-RequestSummary) |  |  |
| patient_buffer | [PatientBufferSummary](#query-analytics-v1-PatientBufferSummary) |  |  |
| period | [SummaryPeriod](#query-analytics-v1-SummaryPeriod) |  |  |






<a name="query-analytics-v1-GetTimeSeriesRequest"></a>

### GetTimeSeriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| from | [string](#string) |  |  |
| to | [string](#string) |  |  |
| clinic_id | [string](#string) | optional |  |
| department_id | [string](#string) | optional |  |
| granularity | [TimeSeriesGranularity](#query-analytics-v1-TimeSeriesGranularity) |  |  |






<a name="query-analytics-v1-GetTimeSeriesResponse"></a>

### GetTimeSeriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| buckets | [TimeSeriesBucket](#query-analytics-v1-TimeSeriesBucket) | repeated |  |






<a name="query-analytics-v1-IncidentPriorityBreakdown"></a>

### IncidentPriorityBreakdown



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| low | [int64](#int64) |  |  |
| normal | [int64](#int64) |  |  |
| high | [int64](#int64) |  |  |
| critical | [int64](#int64) |  |  |






<a name="query-analytics-v1-IncidentSourceBreakdown"></a>

### IncidentSourceBreakdown



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [int64](#int64) |  |  |
| patient | [int64](#int64) |  |  |






<a name="query-analytics-v1-IncidentStatusBreakdown"></a>

### IncidentStatusBreakdown



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pending | [int64](#int64) |  |  |
| in_progress | [int64](#int64) |  |  |
| done | [int64](#int64) |  |  |
| rejected | [int64](#int64) |  |  |
| cancelled | [int64](#int64) |  |  |






<a name="query-analytics-v1-IncidentSummary"></a>

### IncidentSummary



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |
| by_status | [IncidentStatusBreakdown](#query-analytics-v1-IncidentStatusBreakdown) |  |  |
| by_priority | [IncidentPriorityBreakdown](#query-analytics-v1-IncidentPriorityBreakdown) |  |  |
| by_source | [IncidentSourceBreakdown](#query-analytics-v1-IncidentSourceBreakdown) |  |  |
| reopened | [int64](#int64) |  |  |
| with_linked_requests | [int64](#int64) |  |  |
| resolution | [ResolutionStats](#query-analytics-v1-ResolutionStats) | optional |  |
| top_categories | [CategoryCount](#query-analytics-v1-CategoryCount) | repeated |  |
| top_types | [TypeCount](#query-analytics-v1-TypeCount) | repeated |  |
| top_departments | [DepartmentCount](#query-analytics-v1-DepartmentCount) | repeated |  |






<a name="query-analytics-v1-PatientBufferStatusBreakdown"></a>

### PatientBufferStatusBreakdown



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pending | [int64](#int64) |  |  |
| published | [int64](#int64) |  |  |
| rejected | [int64](#int64) |  |  |
| cancelled | [int64](#int64) |  |  |






<a name="query-analytics-v1-PatientBufferSummary"></a>

### PatientBufferSummary



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |
| by_status | [PatientBufferStatusBreakdown](#query-analytics-v1-PatientBufferStatusBreakdown) |  |  |
| acceptance_rate | [double](#double) |  |  |
| rejection_rate | [double](#double) |  |  |






<a name="query-analytics-v1-RequestStatusBreakdown"></a>

### RequestStatusBreakdown



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| created | [int64](#int64) |  |  |
| in_work | [int64](#int64) |  |  |
| on_hold | [int64](#int64) |  |  |
| pending_review | [int64](#int64) |  |  |
| completed | [int64](#int64) |  |  |
| cancelled | [int64](#int64) |  |  |






<a name="query-analytics-v1-RequestSummary"></a>

### RequestSummary



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |
| by_status | [RequestStatusBreakdown](#query-analytics-v1-RequestStatusBreakdown) |  |  |
| linked | [int64](#int64) |  |  |
| unlinked | [int64](#int64) |  |  |
| completion | [ResolutionStats](#query-analytics-v1-ResolutionStats) | optional |  |
| top_types | [TypeCount](#query-analytics-v1-TypeCount) | repeated |  |
| top_departments | [DepartmentCount](#query-analytics-v1-DepartmentCount) | repeated |  |






<a name="query-analytics-v1-ResolutionStats"></a>

### ResolutionStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| avg_minutes | [double](#double) |  |  |
| min_minutes | [double](#double) |  |  |
| max_minutes | [double](#double) |  |  |
| p50_minutes | [double](#double) |  |  |
| p90_minutes | [double](#double) |  |  |
| p95_minutes | [double](#double) |  |  |






<a name="query-analytics-v1-SnapshotIncident"></a>

### SnapshotIncident



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| created_at | [string](#string) |  |  |
| occurred_at | [string](#string) |  |  |
| closed_at | [string](#string) | optional |  |
| status | [string](#string) |  |  |
| priority | [string](#string) |  |  |
| category_id | [string](#string) |  |  |
| category_name | [string](#string) |  |  |
| type_id | [string](#string) |  |  |
| type_name | [string](#string) |  |  |
| clinic_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| is_patient_source | [bool](#bool) |  |  |
| is_reopened | [bool](#bool) |  |  |
| linked_requests_count | [int32](#int32) |  |  |






<a name="query-analytics-v1-SnapshotPatientBuffer"></a>

### SnapshotPatientBuffer



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| created_at | [string](#string) |  |  |
| status | [string](#string) |  |  |
| category_id | [string](#string) | optional |  |






<a name="query-analytics-v1-SnapshotRequest"></a>

### SnapshotRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| created_at | [string](#string) |  |  |
| completed_at | [string](#string) | optional |  |
| status | [string](#string) |  |  |
| type_id | [string](#string) |  |  |
| type_name | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| has_linked_incident | [bool](#bool) |  |  |






<a name="query-analytics-v1-SummaryPeriod"></a>

### SummaryPeriod



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| from | [string](#string) |  |  |
| to | [string](#string) |  |  |
| clinic_id | [string](#string) | optional |  |
| department_id | [string](#string) | optional |  |






<a name="query-analytics-v1-TimeSeriesBucket"></a>

### TimeSeriesBucket



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| bucket_start | [string](#string) |  |  |
| bucket_end | [string](#string) |  |  |
| incidents | [TimeSeriesIncidentBucket](#query-analytics-v1-TimeSeriesIncidentBucket) |  |  |
| requests | [TimeSeriesRequestBucket](#query-analytics-v1-TimeSeriesRequestBucket) |  |  |






<a name="query-analytics-v1-TimeSeriesIncidentBucket"></a>

### TimeSeriesIncidentBucket



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |
| pending | [int64](#int64) |  |  |
| in_progress | [int64](#int64) |  |  |
| done | [int64](#int64) |  |  |
| rejected | [int64](#int64) |  |  |
| cancelled | [int64](#int64) |  |  |
| high_critical | [int64](#int64) |  |  |
| patient_source | [int64](#int64) |  |  |
| reopened | [int64](#int64) |  |  |






<a name="query-analytics-v1-TimeSeriesRequestBucket"></a>

### TimeSeriesRequestBucket



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |
| completed | [int64](#int64) |  |  |
| cancelled | [int64](#int64) |  |  |
| linked | [int64](#int64) |  |  |






<a name="query-analytics-v1-TypeCount"></a>

### TypeCount



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type_id | [string](#string) |  |  |
| type_name | [string](#string) |  |  |
| count | [int64](#int64) |  |  |








<a name="query-analytics-v1-TimeSeriesGranularity"></a>

### TimeSeriesGranularity


| Name | Number | Description |
| ---- | ------ | ----------- |
| TIME_SERIES_GRANULARITY_UNSPECIFIED | 0 |  |
| TIME_SERIES_GRANULARITY_DAY | 1 |  |
| TIME_SERIES_GRANULARITY_WEEK | 2 |  |
| TIME_SERIES_GRANULARITY_MONTH | 3 |  |







<a name="query-analytics-v1-AnalyticsQueryService"></a>

### AnalyticsQueryService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetSnapshot | [GetSnapshotRequest](#query-analytics-v1-GetSnapshotRequest) | [GetSnapshotResponse](#query-analytics-v1-GetSnapshotResponse) |  |
| GetSummary | [GetSummaryRequest](#query-analytics-v1-GetSummaryRequest) | [GetSummaryResponse](#query-analytics-v1-GetSummaryResponse) |  |
| GetTimeSeries | [GetTimeSeriesRequest](#query-analytics-v1-GetTimeSeriesRequest) | [GetTimeSeriesResponse](#query-analytics-v1-GetTimeSeriesResponse) |  |





<a name="query_announcement_v1_announcement-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/announcement/v1/announcement.proto



<a name="query-announcement-v1-AnnouncementView"></a>

### AnnouncementView



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| clinic_id | [string](#string) | optional |  |
| department_id | [string](#string) | optional |  |
| author_id | [string](#string) |  |  |
| title | [string](#string) |  |  |
| content | [string](#string) |  |  |
| priority | [AnnouncementPriority](#query-announcement-v1-AnnouncementPriority) |  |  |
| is_archived | [bool](#bool) |  |  |
| starts_at | [string](#string) | optional | RFC3339Nano |
| ends_at | [string](#string) | optional | RFC3339Nano |
| created_at | [string](#string) |  | RFC3339Nano |
| updated_at | [string](#string) |  | RFC3339Nano |
| view_count | [int64](#int64) |  |  |






<a name="query-announcement-v1-GetAnnouncementRequest"></a>

### GetAnnouncementRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-announcement-v1-GetAnnouncementResponse"></a>

### GetAnnouncementResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| announcement | [AnnouncementView](#query-announcement-v1-AnnouncementView) |  |  |






<a name="query-announcement-v1-ListAnnouncementsForClinicRequest"></a>

### ListAnnouncementsForClinicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| include_archived | [bool](#bool) |  |  |
| priority | [AnnouncementPriority](#query-announcement-v1-AnnouncementPriority) |  |  |
| limit | [int32](#int32) |  |  |
| cursor | [string](#string) | optional |  |






<a name="query-announcement-v1-ListAnnouncementsForClinicResponse"></a>

### ListAnnouncementsForClinicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [AnnouncementView](#query-announcement-v1-AnnouncementView) | repeated |  |
| next_cursor | [string](#string) | optional |  |






<a name="query-announcement-v1-ListAnnouncementsForDepartmentRequest"></a>

### ListAnnouncementsForDepartmentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| include_archived | [bool](#bool) |  |  |
| priority | [AnnouncementPriority](#query-announcement-v1-AnnouncementPriority) |  |  |
| limit | [int32](#int32) |  |  |
| cursor | [string](#string) | optional |  |






<a name="query-announcement-v1-ListAnnouncementsForDepartmentResponse"></a>

### ListAnnouncementsForDepartmentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [AnnouncementView](#query-announcement-v1-AnnouncementView) | repeated |  |
| next_cursor | [string](#string) | optional |  |






<a name="query-announcement-v1-ListAnnouncementsForOrganizationRequest"></a>

### ListAnnouncementsForOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| include_archived | [bool](#bool) |  |  |
| priority | [AnnouncementPriority](#query-announcement-v1-AnnouncementPriority) |  | UNSPECIFIED = all |
| limit | [int32](#int32) |  |  |
| cursor | [string](#string) | optional |  |






<a name="query-announcement-v1-ListAnnouncementsForOrganizationResponse"></a>

### ListAnnouncementsForOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [AnnouncementView](#query-announcement-v1-AnnouncementView) | repeated |  |
| next_cursor | [string](#string) | optional |  |








<a name="query-announcement-v1-AnnouncementPriority"></a>

### AnnouncementPriority


| Name | Number | Description |
| ---- | ------ | ----------- |
| ANNOUNCEMENT_PRIORITY_UNSPECIFIED | 0 |  |
| ANNOUNCEMENT_PRIORITY_NORMAL | 1 |  |
| ANNOUNCEMENT_PRIORITY_HIGH | 2 |  |







<a name="query-announcement-v1-AnnouncementQueryService"></a>

### AnnouncementQueryService
AnnouncementQueryService is the read-side contract for announcements.
GetAnnouncement increments the view counter on each call.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetAnnouncement | [GetAnnouncementRequest](#query-announcement-v1-GetAnnouncementRequest) | [GetAnnouncementResponse](#query-announcement-v1-GetAnnouncementResponse) |  |
| ListAnnouncementsForOrganization | [ListAnnouncementsForOrganizationRequest](#query-announcement-v1-ListAnnouncementsForOrganizationRequest) | [ListAnnouncementsForOrganizationResponse](#query-announcement-v1-ListAnnouncementsForOrganizationResponse) |  |
| ListAnnouncementsForClinic | [ListAnnouncementsForClinicRequest](#query-announcement-v1-ListAnnouncementsForClinicRequest) | [ListAnnouncementsForClinicResponse](#query-announcement-v1-ListAnnouncementsForClinicResponse) |  |
| ListAnnouncementsForDepartment | [ListAnnouncementsForDepartmentRequest](#query-announcement-v1-ListAnnouncementsForDepartmentRequest) | [ListAnnouncementsForDepartmentResponse](#query-announcement-v1-ListAnnouncementsForDepartmentResponse) |  |





<a name="query_incident_classifier_v1_classifier-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/incident/classifier/v1/classifier.proto



<a name="query-incident-classifier-v1-Category"></a>

### Category
Category mirrors projections.incident_categories row.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| parent_category_id | [string](#string) | optional |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| is_active | [bool](#bool) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="query-incident-classifier-v1-GetCategoryRequest"></a>

### GetCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-incident-classifier-v1-GetCategoryResponse"></a>

### GetCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category | [Category](#query-incident-classifier-v1-Category) |  |  |






<a name="query-incident-classifier-v1-GetTypeRequest"></a>

### GetTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-incident-classifier-v1-GetTypeResponse"></a>

### GetTypeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| type | [Type](#query-incident-classifier-v1-Type) |  |  |






<a name="query-incident-classifier-v1-ListActiveRootCategoriesRequest"></a>

### ListActiveRootCategoriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-classifier-v1-ListActiveRootCategoriesResponse"></a>

### ListActiveRootCategoriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Category](#query-incident-classifier-v1-Category) | repeated |  |






<a name="query-incident-classifier-v1-ListActiveTypesByOrganizationRequest"></a>

### ListActiveTypesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-classifier-v1-ListActiveTypesByOrganizationResponse"></a>

### ListActiveTypesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Type](#query-incident-classifier-v1-Type) | repeated |  |






<a name="query-incident-classifier-v1-ListCategoriesByOrganizationRequest"></a>

### ListCategoriesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-classifier-v1-ListCategoriesByOrganizationResponse"></a>

### ListCategoriesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Category](#query-incident-classifier-v1-Category) | repeated |  |






<a name="query-incident-classifier-v1-ListCategorySubtreeRequest"></a>

### ListCategorySubtreeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| root_category_id | [string](#string) |  |  |






<a name="query-incident-classifier-v1-ListCategorySubtreeResponse"></a>

### ListCategorySubtreeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Category](#query-incident-classifier-v1-Category) | repeated |  |






<a name="query-incident-classifier-v1-ListPatientAllowedTypesByOrganizationRequest"></a>

### ListPatientAllowedTypesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-classifier-v1-ListPatientAllowedTypesByOrganizationResponse"></a>

### ListPatientAllowedTypesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Type](#query-incident-classifier-v1-Type) | repeated |  |






<a name="query-incident-classifier-v1-ListPatientVisibleCategoriesByOrganizationRequest"></a>

### ListPatientVisibleCategoriesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-classifier-v1-ListPatientVisibleCategoriesByOrganizationResponse"></a>

### ListPatientVisibleCategoriesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Category](#query-incident-classifier-v1-Category) | repeated |  |






<a name="query-incident-classifier-v1-ListTypesByCategoryRequest"></a>

### ListTypesByCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-classifier-v1-ListTypesByCategoryResponse"></a>

### ListTypesByCategoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [Type](#query-incident-classifier-v1-Type) | repeated |  |






<a name="query-incident-classifier-v1-Type"></a>

### Type
Type mirrors projections.incident_types row.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| category_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| is_active | [bool](#bool) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| is_allowed_for_patients | [bool](#bool) |  |  |












<a name="query-incident-classifier-v1-IncidentClassifierQueryService"></a>

### IncidentClassifierQueryService
IncidentClassifierQueryService exposes read methods over the
projections.incident_categories / projections.incident_types
tables. Subtree walks (descendants of a category) use recursive SQL
CTEs in the reader; the RPC surface stays flat.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetCategory | [GetCategoryRequest](#query-incident-classifier-v1-GetCategoryRequest) | [GetCategoryResponse](#query-incident-classifier-v1-GetCategoryResponse) |  |
| ListCategoriesByOrganization | [ListCategoriesByOrganizationRequest](#query-incident-classifier-v1-ListCategoriesByOrganizationRequest) | [ListCategoriesByOrganizationResponse](#query-incident-classifier-v1-ListCategoriesByOrganizationResponse) |  |
| ListActiveRootCategories | [ListActiveRootCategoriesRequest](#query-incident-classifier-v1-ListActiveRootCategoriesRequest) | [ListActiveRootCategoriesResponse](#query-incident-classifier-v1-ListActiveRootCategoriesResponse) |  |
| ListCategorySubtree | [ListCategorySubtreeRequest](#query-incident-classifier-v1-ListCategorySubtreeRequest) | [ListCategorySubtreeResponse](#query-incident-classifier-v1-ListCategorySubtreeResponse) |  |
| GetType | [GetTypeRequest](#query-incident-classifier-v1-GetTypeRequest) | [GetTypeResponse](#query-incident-classifier-v1-GetTypeResponse) |  |
| ListTypesByCategory | [ListTypesByCategoryRequest](#query-incident-classifier-v1-ListTypesByCategoryRequest) | [ListTypesByCategoryResponse](#query-incident-classifier-v1-ListTypesByCategoryResponse) |  |
| ListActiveTypesByOrganization | [ListActiveTypesByOrganizationRequest](#query-incident-classifier-v1-ListActiveTypesByOrganizationRequest) | [ListActiveTypesByOrganizationResponse](#query-incident-classifier-v1-ListActiveTypesByOrganizationResponse) |  |
| ListPatientAllowedTypesByOrganization | [ListPatientAllowedTypesByOrganizationRequest](#query-incident-classifier-v1-ListPatientAllowedTypesByOrganizationRequest) | [ListPatientAllowedTypesByOrganizationResponse](#query-incident-classifier-v1-ListPatientAllowedTypesByOrganizationResponse) | Patient-facing reads. These are scoped to one organisation and return only the slice of the classifier that a patient may see when filing an incident. |
| ListPatientVisibleCategoriesByOrganization | [ListPatientVisibleCategoriesByOrganizationRequest](#query-incident-classifier-v1-ListPatientVisibleCategoriesByOrganizationRequest) | [ListPatientVisibleCategoriesByOrganizationResponse](#query-incident-classifier-v1-ListPatientVisibleCategoriesByOrganizationResponse) |  |





<a name="query_incident_v1_incident-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/incident/v1/incident.proto



<a name="query-incident-v1-ActorView"></a>

### ActorView



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) | optional |  |
| display_name | [string](#string) | optional |  |






<a name="query-incident-v1-BufferEntryView"></a>

### BufferEntryView



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| patient_zitadel_user_id | [string](#string) |  |  |
| category_id | [string](#string) | optional |  |
| type_id | [string](#string) | optional |  |
| description | [string](#string) | optional |  |
| occurred_at | [string](#string) | optional |  |
| status | [BufferStatus](#query-incident-v1-BufferStatus) |  |  |
| published_incident_id | [string](#string) | optional |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| patient_status | [PatientStatus](#query-incident-v1-PatientStatus) | optional | Populated only for patient callers. |






<a name="query-incident-v1-GetBufferEntryRequest"></a>

### GetBufferEntryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-incident-v1-GetBufferEntryResponse"></a>

### GetBufferEntryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| entry | [BufferEntryView](#query-incident-v1-BufferEntryView) |  |  |






<a name="query-incident-v1-GetIncidentHistoryRequest"></a>

### GetIncidentHistoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |






<a name="query-incident-v1-GetIncidentHistoryResponse"></a>

### GetIncidentHistoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status_history | [StatusHistoryEntry](#query-incident-v1-StatusHistoryEntry) | repeated |  |
| priority_history | [PriorityHistoryEntry](#query-incident-v1-PriorityHistoryEntry) | repeated |  |






<a name="query-incident-v1-GetIncidentRequest"></a>

### GetIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-incident-v1-GetIncidentResponse"></a>

### GetIncidentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident | [IncidentView](#query-incident-v1-IncidentView) |  |  |






<a name="query-incident-v1-IncidentView"></a>

### IncidentView
IncidentView is the unified payload. For patients only id, status,
patient_status, description and timestamps are populated.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) | optional |  |
| clinic_id | [string](#string) | optional |  |
| department_id | [string](#string) | optional |  |
| category_id | [string](#string) | optional |  |
| type_id | [string](#string) | optional |  |
| status | [IncidentStatus](#query-incident-v1-IncidentStatus) |  |  |
| priority | [IncidentPriority](#query-incident-v1-IncidentPriority) |  |  |
| description | [string](#string) | optional |  |
| patient_original_description | [string](#string) | optional |  |
| occurred_at | [string](#string) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |
| registrar | [RegistrarView](#query-incident-v1-RegistrarView) | optional |  |
| source_patient_zitadel_user_id | [string](#string) | optional |  |
| source_buffer_id | [string](#string) | optional |  |
| reopened_from_incident_id | [string](#string) | optional |  |
| patient_status | [PatientStatus](#query-incident-v1-PatientStatus) | optional | Populated only for patient callers. |






<a name="query-incident-v1-ListBufferEntriesRequest"></a>

### ListBufferEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| statuses | [BufferStatus](#query-incident-v1-BufferStatus) | repeated |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-v1-ListBufferEntriesResponse"></a>

### ListBufferEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [BufferEntryView](#query-incident-v1-BufferEntryView) | repeated |  |






<a name="query-incident-v1-ListIncidentsRequest"></a>

### ListIncidentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| statuses | [IncidentStatus](#query-incident-v1-IncidentStatus) | repeated |  |
| priorities | [IncidentPriority](#query-incident-v1-IncidentPriority) | repeated |  |
| clinic_id | [string](#string) | optional |  |
| department_id | [string](#string) | optional |  |
| category_id | [string](#string) | optional |  |
| type_id | [string](#string) | optional |  |
| occurred_from | [string](#string) | optional | RFC3339Nano |
| occurred_to | [string](#string) | optional |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-v1-ListIncidentsResponse"></a>

### ListIncidentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [IncidentView](#query-incident-v1-IncidentView) | repeated |  |






<a name="query-incident-v1-ListMyBufferEntriesRequest"></a>

### ListMyBufferEntriesRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-v1-ListMyBufferEntriesResponse"></a>

### ListMyBufferEntriesResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [BufferEntryView](#query-incident-v1-BufferEntryView) | repeated |  |






<a name="query-incident-v1-ListMyIncidentsRequest"></a>

### ListMyIncidentsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-incident-v1-ListMyIncidentsResponse"></a>

### ListMyIncidentsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [IncidentView](#query-incident-v1-IncidentView) | repeated |  |






<a name="query-incident-v1-PriorityHistoryEntry"></a>

### PriorityHistoryEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| old_priority | [IncidentPriority](#query-incident-v1-IncidentPriority) |  |  |
| new_priority | [IncidentPriority](#query-incident-v1-IncidentPriority) |  |  |
| actor | [ActorView](#query-incident-v1-ActorView) |  |  |
| changed_at | [string](#string) |  |  |






<a name="query-incident-v1-RegistrarView"></a>

### RegistrarView



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| display_name | [string](#string) |  |  |
| position | [string](#string) | optional |  |
| organization_id | [string](#string) |  |  |
| clinic_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |






<a name="query-incident-v1-StatusHistoryEntry"></a>

### StatusHistoryEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| old_status | [IncidentStatus](#query-incident-v1-IncidentStatus) |  | UNSPECIFIED for the initial entry |
| new_status | [IncidentStatus](#query-incident-v1-IncidentStatus) |  |  |
| actor | [ActorView](#query-incident-v1-ActorView) |  |  |
| changed_at | [string](#string) |  |  |








<a name="query-incident-v1-BufferStatus"></a>

### BufferStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| BUFFER_STATUS_UNSPECIFIED | 0 |  |
| BUFFER_STATUS_PENDING | 1 |  |
| BUFFER_STATUS_PUBLISHED | 2 |  |
| BUFFER_STATUS_REJECTED | 3 |  |
| BUFFER_STATUS_CANCELLED | 4 |  |



<a name="query-incident-v1-IncidentPriority"></a>

### IncidentPriority


| Name | Number | Description |
| ---- | ------ | ----------- |
| INCIDENT_PRIORITY_UNSPECIFIED | 0 |  |
| INCIDENT_PRIORITY_LOW | 1 |  |
| INCIDENT_PRIORITY_NORMAL | 2 |  |
| INCIDENT_PRIORITY_HIGH | 3 |  |
| INCIDENT_PRIORITY_CRITICAL | 4 |  |



<a name="query-incident-v1-IncidentStatus"></a>

### IncidentStatus


| Name | Number | Description |
| ---- | ------ | ----------- |
| INCIDENT_STATUS_UNSPECIFIED | 0 |  |
| INCIDENT_STATUS_PENDING | 1 |  |
| INCIDENT_STATUS_IN_PROGRESS | 2 |  |
| INCIDENT_STATUS_DONE | 3 |  |
| INCIDENT_STATUS_REJECTED | 4 |  |
| INCIDENT_STATUS_CANCELLED | 5 |  |



<a name="query-incident-v1-PatientStatus"></a>

### PatientStatus
PatientStatus is the simplified four-value status surfaced to patients.

| Name | Number | Description |
| ---- | ------ | ----------- |
| PATIENT_STATUS_UNSPECIFIED | 0 |  |
| PATIENT_STATUS_PENDING | 1 | buffer pending |
| PATIENT_STATUS_ACCEPTED | 2 | dispatcher accepted; incident pending/in_progress |
| PATIENT_STATUS_CLOSED | 3 | done / rejected / buffer rejected |
| PATIENT_STATUS_CANCELLED | 4 | patient cancelled |







<a name="query-incident-v1-IncidentQueryService"></a>

### IncidentQueryService
IncidentQueryService is the unified read-side. The same RPCs serve
both employees and patients — visibility filtering and field
redaction happen inside the reader based on the caller&#39;s role.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetIncident | [GetIncidentRequest](#query-incident-v1-GetIncidentRequest) | [GetIncidentResponse](#query-incident-v1-GetIncidentResponse) |  |
| ListIncidents | [ListIncidentsRequest](#query-incident-v1-ListIncidentsRequest) | [ListIncidentsResponse](#query-incident-v1-ListIncidentsResponse) |  |
| ListMyIncidents | [ListMyIncidentsRequest](#query-incident-v1-ListMyIncidentsRequest) | [ListMyIncidentsResponse](#query-incident-v1-ListMyIncidentsResponse) |  |
| GetIncidentHistory | [GetIncidentHistoryRequest](#query-incident-v1-GetIncidentHistoryRequest) | [GetIncidentHistoryResponse](#query-incident-v1-GetIncidentHistoryResponse) |  |
| GetBufferEntry | [GetBufferEntryRequest](#query-incident-v1-GetBufferEntryRequest) | [GetBufferEntryResponse](#query-incident-v1-GetBufferEntryResponse) | Buffer reads. |
| ListBufferEntries | [ListBufferEntriesRequest](#query-incident-v1-ListBufferEntriesRequest) | [ListBufferEntriesResponse](#query-incident-v1-ListBufferEntriesResponse) |  |
| ListMyBufferEntries | [ListMyBufferEntriesRequest](#query-incident-v1-ListMyBufferEntriesRequest) | [ListMyBufferEntriesResponse](#query-incident-v1-ListMyBufferEntriesResponse) |  |





<a name="query_membership_v1_membership-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/membership/v1/membership.proto



<a name="query-membership-v1-CountEmployeesByClinicRequest"></a>

### CountEmployeesByClinicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| include_terminated | [bool](#bool) |  | When false (default), rows with terminated_at IS NOT NULL are hidden. Set true to include offboarded employees. |
| on_vacation | [bool](#bool) |  | When true, restrict to employees currently on an active vacation (current_vacation_ends_at IS NOT NULL AND &gt; now()). |
| position | [string](#string) | optional | Optional exact-match filter on employee_cards.position. Trimmed before comparison; all-whitespace is treated as unset. |






<a name="query-membership-v1-CountEmployeesByClinicResponse"></a>

### CountEmployeesByClinicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |






<a name="query-membership-v1-CountEmployeesByDepartmentRequest"></a>

### CountEmployeesByDepartmentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| include_terminated | [bool](#bool) |  | When false (default), rows with terminated_at IS NOT NULL are hidden. Set true to include offboarded employees. |
| on_vacation | [bool](#bool) |  | When true, restrict to employees currently on an active vacation (current_vacation_ends_at IS NOT NULL AND &gt; now()). |
| position | [string](#string) | optional | Optional exact-match filter on employee_cards.position. Trimmed before comparison; all-whitespace is treated as unset. |






<a name="query-membership-v1-CountEmployeesByDepartmentResponse"></a>

### CountEmployeesByDepartmentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |






<a name="query-membership-v1-CountEmployeesByOrganizationRequest"></a>

### CountEmployeesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| include_terminated | [bool](#bool) |  | When false (default), rows with terminated_at IS NOT NULL are hidden. Set true to include offboarded employees. |
| on_vacation | [bool](#bool) |  | When true, restrict to employees currently on an active vacation (current_vacation_ends_at IS NOT NULL AND &gt; now()). |
| position | [string](#string) | optional | Optional exact-match filter on employee_cards.position. Trimmed before comparison; all-whitespace is treated as unset. |






<a name="query-membership-v1-CountEmployeesByOrganizationResponse"></a>

### CountEmployeesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |






<a name="query-membership-v1-CountVacationsByEmployeeRequest"></a>

### CountVacationsByEmployeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| state | [string](#string) |  | Optional state filter. When empty, counts all states. Valid values: scheduled, active, ended, cancelled. |






<a name="query-membership-v1-CountVacationsByEmployeeResponse"></a>

### CountVacationsByEmployeeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |






<a name="query-membership-v1-EmployeeCardView"></a>

### EmployeeCardView
EmployeeCardView is the denormalised card projection returned by
both Get and List endpoints. Fields mirror projections.employee_cards
columns; timestamps are RFC3339 strings. Optional fields stay unset
when the backing column is NULL.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| zitadel_user_id | [string](#string) |  |  |
| first_name | [string](#string) | optional |  |
| last_name | [string](#string) | optional |  |
| display_name | [string](#string) | optional |  |
| email | [string](#string) | optional |  |
| organization_id | [string](#string) |  |  |
| organization_name | [string](#string) | optional |  |
| clinic_id | [string](#string) | optional |  |
| clinic_name | [string](#string) | optional |  |
| department_id | [string](#string) |  |  |
| department_name | [string](#string) | optional |  |
| position | [string](#string) | optional |  |
| terminated_at | [string](#string) | optional |  |
| current_vacation_ends_at | [string](#string) | optional |  |
| next_vacation_starts_at | [string](#string) | optional |  |






<a name="query-membership-v1-GetClinicHeadRequest"></a>

### GetClinicHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |






<a name="query-membership-v1-GetClinicHeadResponse"></a>

### GetClinicHeadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignment | [RoleAssignment](#query-membership-v1-RoleAssignment) | optional |  |






<a name="query-membership-v1-GetDepartmentResponsibleRequest"></a>

### GetDepartmentResponsibleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |






<a name="query-membership-v1-GetDepartmentResponsibleResponse"></a>

### GetDepartmentResponsibleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| assignment | [RoleAssignment](#query-membership-v1-RoleAssignment) | optional |  |






<a name="query-membership-v1-GetEmployeeRequest"></a>

### GetEmployeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-membership-v1-GetEmployeeResponse"></a>

### GetEmployeeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee | [EmployeeCardView](#query-membership-v1-EmployeeCardView) |  |  |






<a name="query-membership-v1-ListCandidatesForClinicHeadRequest"></a>

### ListCandidatesForClinicHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| query | [string](#string) |  | Optional substring search (ILIKE %query%) on first_name, last_name, display_name, email. Trimmed at the handler boundary. Max 256 chars. |
| after | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |






<a name="query-membership-v1-ListCandidatesForClinicHeadResponse"></a>

### ListCandidatesForClinicHeadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |
| next_cursor | [string](#string) |  |  |






<a name="query-membership-v1-ListCandidatesForDeptResponsibleRequest"></a>

### ListCandidatesForDeptResponsibleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| query | [string](#string) |  | Optional substring search (ILIKE %query%) on first_name, last_name, display_name, email. Trimmed at the handler boundary. Max 256 chars. |
| after | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |






<a name="query-membership-v1-ListCandidatesForDeptResponsibleResponse"></a>

### ListCandidatesForDeptResponsibleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |
| next_cursor | [string](#string) |  |  |






<a name="query-membership-v1-ListCandidatesForHireRequest"></a>

### ListCandidatesForHireRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| query | [string](#string) |  | Optional substring search (ILIKE %query%) on first_name, last_name, display_name, email. Trimmed at the handler boundary. Max 256 chars. |
| after | [string](#string) |  | Opaque cursor from next_cursor of a previous response. Empty = first page. |
| limit | [int32](#int32) |  | Page size. 0 → default 50. Range [1, 500]. |






<a name="query-membership-v1-ListCandidatesForHireResponse"></a>

### ListCandidatesForHireResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ZitadelUserView](#query-membership-v1-ZitadelUserView) | repeated |  |
| next_cursor | [string](#string) |  | Opaque cursor for the next page. Empty when this is the last page. |






<a name="query-membership-v1-ListCandidatesForOrgAdminRequest"></a>

### ListCandidatesForOrgAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| query | [string](#string) |  | Optional substring search (ILIKE %query%) on first_name, last_name, display_name, email. Trimmed at the handler boundary. Max 256 chars. |
| after | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |






<a name="query-membership-v1-ListCandidatesForOrgAdminResponse"></a>

### ListCandidatesForOrgAdminResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |
| next_cursor | [string](#string) |  |  |






<a name="query-membership-v1-ListCandidatesForOrgDispatcherRequest"></a>

### ListCandidatesForOrgDispatcherRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| query | [string](#string) |  | Optional substring search (ILIKE %query%) on first_name, last_name, display_name, email. Trimmed at the handler boundary. Max 256 chars. |
| after | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |






<a name="query-membership-v1-ListCandidatesForOrgDispatcherResponse"></a>

### ListCandidatesForOrgDispatcherResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |
| next_cursor | [string](#string) |  |  |






<a name="query-membership-v1-ListCandidatesForOrgHeadRequest"></a>

### ListCandidatesForOrgHeadRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| query | [string](#string) |  | Optional substring search (ILIKE %query%) on first_name, last_name, display_name, email. Trimmed at the handler boundary. Max 256 chars. |
| after | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |






<a name="query-membership-v1-ListCandidatesForOrgHeadResponse"></a>

### ListCandidatesForOrgHeadResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |
| next_cursor | [string](#string) |  |  |






<a name="query-membership-v1-ListCandidatesForSystemAdminRequest"></a>

### ListCandidatesForSystemAdminRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| query | [string](#string) |  | Optional substring search (ILIKE %query%) on first_name, last_name, display_name, email. Trimmed at the handler boundary. Max 256 chars. |
| after | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |






<a name="query-membership-v1-ListCandidatesForSystemAdminResponse"></a>

### ListCandidatesForSystemAdminResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ZitadelUserView](#query-membership-v1-ZitadelUserView) | repeated |  |
| next_cursor | [string](#string) |  |  |






<a name="query-membership-v1-ListEmployeesByClinicRequest"></a>

### ListEmployeesByClinicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |
| include_terminated | [bool](#bool) |  | When false (default), rows with terminated_at IS NOT NULL are hidden. Set true to include offboarded employees. |
| on_vacation | [bool](#bool) |  | When true, restrict to employees currently on an active vacation (current_vacation_ends_at IS NOT NULL AND &gt; now()). |
| position | [string](#string) | optional | Optional exact-match filter on employee_cards.position. Trimmed before comparison; all-whitespace is treated as unset. |






<a name="query-membership-v1-ListEmployeesByClinicResponse"></a>

### ListEmployeesByClinicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |






<a name="query-membership-v1-ListEmployeesByDepartmentRequest"></a>

### ListEmployeesByDepartmentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |
| include_terminated | [bool](#bool) |  | When false (default), rows with terminated_at IS NOT NULL are hidden. Set true to include offboarded employees. |
| on_vacation | [bool](#bool) |  | When true, restrict to employees currently on an active vacation (current_vacation_ends_at IS NOT NULL AND &gt; now()). |
| position | [string](#string) | optional | Optional exact-match filter on employee_cards.position. Trimmed before comparison; all-whitespace is treated as unset. |






<a name="query-membership-v1-ListEmployeesByDepartmentResponse"></a>

### ListEmployeesByDepartmentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |






<a name="query-membership-v1-ListEmployeesByOrganizationRequest"></a>

### ListEmployeesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |
| include_terminated | [bool](#bool) |  | When false (default), rows with terminated_at IS NOT NULL are hidden. Set true to include offboarded employees. |
| on_vacation | [bool](#bool) |  | When true, restrict to employees currently on an active vacation (current_vacation_ends_at IS NOT NULL AND &gt; now()). |
| position | [string](#string) | optional | Optional exact-match filter on employee_cards.position. Trimmed before comparison; all-whitespace is treated as unset. |






<a name="query-membership-v1-ListEmployeesByOrganizationResponse"></a>

### ListEmployeesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |






<a name="query-membership-v1-ListOrgAdminsRequest"></a>

### ListOrgAdminsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-membership-v1-ListOrgAdminsResponse"></a>

### ListOrgAdminsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [RoleAssignment](#query-membership-v1-RoleAssignment) | repeated |  |






<a name="query-membership-v1-ListOrgDispatchersRequest"></a>

### ListOrgDispatchersRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-membership-v1-ListOrgDispatchersResponse"></a>

### ListOrgDispatchersResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [RoleAssignment](#query-membership-v1-RoleAssignment) | repeated |  |






<a name="query-membership-v1-ListOrgHeadsRequest"></a>

### ListOrgHeadsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-membership-v1-ListOrgHeadsResponse"></a>

### ListOrgHeadsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [RoleAssignment](#query-membership-v1-RoleAssignment) | repeated |  |






<a name="query-membership-v1-ListSystemAdminsRequest"></a>

### ListSystemAdminsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-membership-v1-ListSystemAdminsResponse"></a>

### ListSystemAdminsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [SystemAdminView](#query-membership-v1-SystemAdminView) | repeated |  |






<a name="query-membership-v1-ListVacationsByEmployeeRequest"></a>

### ListVacationsByEmployeeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| state | [string](#string) |  | Optional state filter. When empty, all states are returned. Valid values: scheduled, active, ended, cancelled. |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-membership-v1-ListVacationsByEmployeeResponse"></a>

### ListVacationsByEmployeeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [VacationView](#query-membership-v1-VacationView) | repeated |  |






<a name="query-membership-v1-RoleAssignment"></a>

### RoleAssignment
RoleAssignment is a role row enriched with the denormalised card for
the holder and, when present, the deputy. Read-model callers use
this so they do not need to follow role lookups with N&#43;1 GetEmployee
calls to render a name or email.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| holder | [EmployeeCardView](#query-membership-v1-EmployeeCardView) |  |  |
| deputy | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | optional |  |






<a name="query-membership-v1-SearchEmployeesByOrganizationRequest"></a>

### SearchEmployeesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| query | [string](#string) |  | Fuzzy substring matched case-insensitively (ILIKE %query%) against first_name, last_name, display_name, and email. Trimmed at the handler boundary; an empty query degenerates to the same behaviour as ListEmployeesByOrganization with the same filters. Maximum length 256 characters; over-long inputs are rejected before the authz round-trip. |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |
| include_terminated | [bool](#bool) |  | When false (default), rows with terminated_at IS NOT NULL are hidden. Set true to include offboarded employees. |
| on_vacation | [bool](#bool) |  | When true, restrict to employees currently on an active vacation (current_vacation_ends_at IS NOT NULL AND &gt; now()). |
| position | [string](#string) | optional | Optional exact-match filter on employee_cards.position. Trimmed before comparison; all-whitespace is treated as unset. |






<a name="query-membership-v1-SearchEmployeesByOrganizationResponse"></a>

### SearchEmployeesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [EmployeeCardView](#query-membership-v1-EmployeeCardView) | repeated |  |






<a name="query-membership-v1-SystemAdminView"></a>

### SystemAdminView
SystemAdminView is the system-admin role; system admins are rooted in
Zitadel user ids, not employee ids.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |
| created_at | [string](#string) |  |  |






<a name="query-membership-v1-VacationView"></a>

### VacationView
VacationView mirrors projections.employee_vacations. state is one of
{scheduled, active, ended, cancelled}.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| state | [string](#string) |  |  |
| starts_at | [string](#string) |  |  |
| ends_at | [string](#string) | optional |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="query-membership-v1-ZitadelUserView"></a>

### ZitadelUserView
ZitadelUserView is a lightweight projection of projections.users used
by candidate-listing endpoints that operate on raw Zitadel identities
(ForHire, ForSystemAdmin). Fields mirror the NOT NULL columns of the
users table.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| zitadel_user_id | [string](#string) |  |  |
| first_name | [string](#string) |  |  |
| last_name | [string](#string) |  |  |
| display_name | [string](#string) |  |  |
| email | [string](#string) |  |  |












<a name="query-membership-v1-MembershipQueryService"></a>

### MembershipQueryService
MembershipQueryService exposes read methods over the employee
projections and the role assignment projections (clinic heads,
department responsibles, org admins/heads/dispatchers, system admins).

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetEmployee | [GetEmployeeRequest](#query-membership-v1-GetEmployeeRequest) | [GetEmployeeResponse](#query-membership-v1-GetEmployeeResponse) |  |
| ListEmployeesByDepartment | [ListEmployeesByDepartmentRequest](#query-membership-v1-ListEmployeesByDepartmentRequest) | [ListEmployeesByDepartmentResponse](#query-membership-v1-ListEmployeesByDepartmentResponse) |  |
| ListEmployeesByClinic | [ListEmployeesByClinicRequest](#query-membership-v1-ListEmployeesByClinicRequest) | [ListEmployeesByClinicResponse](#query-membership-v1-ListEmployeesByClinicResponse) |  |
| ListEmployeesByOrganization | [ListEmployeesByOrganizationRequest](#query-membership-v1-ListEmployeesByOrganizationRequest) | [ListEmployeesByOrganizationResponse](#query-membership-v1-ListEmployeesByOrganizationResponse) |  |
| CountEmployeesByDepartment | [CountEmployeesByDepartmentRequest](#query-membership-v1-CountEmployeesByDepartmentRequest) | [CountEmployeesByDepartmentResponse](#query-membership-v1-CountEmployeesByDepartmentResponse) |  |
| CountEmployeesByClinic | [CountEmployeesByClinicRequest](#query-membership-v1-CountEmployeesByClinicRequest) | [CountEmployeesByClinicResponse](#query-membership-v1-CountEmployeesByClinicResponse) |  |
| CountEmployeesByOrganization | [CountEmployeesByOrganizationRequest](#query-membership-v1-CountEmployeesByOrganizationRequest) | [CountEmployeesByOrganizationResponse](#query-membership-v1-CountEmployeesByOrganizationResponse) |  |
| SearchEmployeesByOrganization | [SearchEmployeesByOrganizationRequest](#query-membership-v1-SearchEmployeesByOrganizationRequest) | [SearchEmployeesByOrganizationResponse](#query-membership-v1-SearchEmployeesByOrganizationResponse) |  |
| ListVacationsByEmployee | [ListVacationsByEmployeeRequest](#query-membership-v1-ListVacationsByEmployeeRequest) | [ListVacationsByEmployeeResponse](#query-membership-v1-ListVacationsByEmployeeResponse) |  |
| CountVacationsByEmployee | [CountVacationsByEmployeeRequest](#query-membership-v1-CountVacationsByEmployeeRequest) | [CountVacationsByEmployeeResponse](#query-membership-v1-CountVacationsByEmployeeResponse) |  |
| GetClinicHead | [GetClinicHeadRequest](#query-membership-v1-GetClinicHeadRequest) | [GetClinicHeadResponse](#query-membership-v1-GetClinicHeadResponse) |  |
| GetDepartmentResponsible | [GetDepartmentResponsibleRequest](#query-membership-v1-GetDepartmentResponsibleRequest) | [GetDepartmentResponsibleResponse](#query-membership-v1-GetDepartmentResponsibleResponse) |  |
| ListOrgAdmins | [ListOrgAdminsRequest](#query-membership-v1-ListOrgAdminsRequest) | [ListOrgAdminsResponse](#query-membership-v1-ListOrgAdminsResponse) |  |
| ListOrgDispatchers | [ListOrgDispatchersRequest](#query-membership-v1-ListOrgDispatchersRequest) | [ListOrgDispatchersResponse](#query-membership-v1-ListOrgDispatchersResponse) |  |
| ListOrgHeads | [ListOrgHeadsRequest](#query-membership-v1-ListOrgHeadsRequest) | [ListOrgHeadsResponse](#query-membership-v1-ListOrgHeadsResponse) |  |
| ListSystemAdmins | [ListSystemAdminsRequest](#query-membership-v1-ListSystemAdminsRequest) | [ListSystemAdminsResponse](#query-membership-v1-ListSystemAdminsResponse) |  |
| ListCandidatesForHire | [ListCandidatesForHireRequest](#query-membership-v1-ListCandidatesForHireRequest) | [ListCandidatesForHireResponse](#query-membership-v1-ListCandidatesForHireResponse) |  |
| ListCandidatesForSystemAdmin | [ListCandidatesForSystemAdminRequest](#query-membership-v1-ListCandidatesForSystemAdminRequest) | [ListCandidatesForSystemAdminResponse](#query-membership-v1-ListCandidatesForSystemAdminResponse) |  |
| ListCandidatesForOrgAdmin | [ListCandidatesForOrgAdminRequest](#query-membership-v1-ListCandidatesForOrgAdminRequest) | [ListCandidatesForOrgAdminResponse](#query-membership-v1-ListCandidatesForOrgAdminResponse) |  |
| ListCandidatesForOrgHead | [ListCandidatesForOrgHeadRequest](#query-membership-v1-ListCandidatesForOrgHeadRequest) | [ListCandidatesForOrgHeadResponse](#query-membership-v1-ListCandidatesForOrgHeadResponse) |  |
| ListCandidatesForOrgDispatcher | [ListCandidatesForOrgDispatcherRequest](#query-membership-v1-ListCandidatesForOrgDispatcherRequest) | [ListCandidatesForOrgDispatcherResponse](#query-membership-v1-ListCandidatesForOrgDispatcherResponse) |  |
| ListCandidatesForClinicHead | [ListCandidatesForClinicHeadRequest](#query-membership-v1-ListCandidatesForClinicHeadRequest) | [ListCandidatesForClinicHeadResponse](#query-membership-v1-ListCandidatesForClinicHeadResponse) |  |
| ListCandidatesForDeptResponsible | [ListCandidatesForDeptResponsibleRequest](#query-membership-v1-ListCandidatesForDeptResponsibleRequest) | [ListCandidatesForDeptResponsibleResponse](#query-membership-v1-ListCandidatesForDeptResponsibleResponse) |  |





<a name="query_orgstructure_v1_orgstructure-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/orgstructure/v1/orgstructure.proto



<a name="query-orgstructure-v1-Address"></a>

### Address
Address is the read-side view of a stored postal address; Point is
optional because the command-side allows text-only addresses.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| text | [string](#string) |  |  |
| point | [Point](#query-orgstructure-v1-Point) | optional |  |






<a name="query-orgstructure-v1-Clinic"></a>

### Clinic
Clinic mirrors the projections.clinics row returned by GetClinic.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| physical_address | [Address](#query-orgstructure-v1-Address) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="query-orgstructure-v1-ClinicListItem"></a>

### ClinicListItem
ClinicListItem is the minimal shape returned by list endpoints.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |






<a name="query-orgstructure-v1-CountClinicsByOrganizationRequest"></a>

### CountClinicsByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |






<a name="query-orgstructure-v1-CountClinicsByOrganizationResponse"></a>

### CountClinicsByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |






<a name="query-orgstructure-v1-CountDepartmentsByClinicRequest"></a>

### CountDepartmentsByClinicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |






<a name="query-orgstructure-v1-CountDepartmentsByClinicResponse"></a>

### CountDepartmentsByClinicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |






<a name="query-orgstructure-v1-CountOrganizationsRequest"></a>

### CountOrganizationsRequest







<a name="query-orgstructure-v1-CountOrganizationsResponse"></a>

### CountOrganizationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| total | [int64](#int64) |  |  |






<a name="query-orgstructure-v1-Department"></a>

### Department
Department mirrors the projections.departments row returned by
GetDepartment.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| clinic_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="query-orgstructure-v1-DepartmentListItem"></a>

### DepartmentListItem
DepartmentListItem is the minimal shape returned by list endpoints.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| clinic_id | [string](#string) |  |  |
| name | [string](#string) |  |  |






<a name="query-orgstructure-v1-GetClinicRequest"></a>

### GetClinicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-orgstructure-v1-GetClinicResponse"></a>

### GetClinicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic | [Clinic](#query-orgstructure-v1-Clinic) |  |  |






<a name="query-orgstructure-v1-GetDepartmentRequest"></a>

### GetDepartmentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-orgstructure-v1-GetDepartmentResponse"></a>

### GetDepartmentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department | [Department](#query-orgstructure-v1-Department) |  |  |






<a name="query-orgstructure-v1-GetOrganizationRequest"></a>

### GetOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-orgstructure-v1-GetOrganizationResponse"></a>

### GetOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization | [Organization](#query-orgstructure-v1-Organization) |  |  |






<a name="query-orgstructure-v1-ListClinicsByOrganizationRequest"></a>

### ListClinicsByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-orgstructure-v1-ListClinicsByOrganizationResponse"></a>

### ListClinicsByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ClinicListItem](#query-orgstructure-v1-ClinicListItem) | repeated |  |






<a name="query-orgstructure-v1-ListDepartmentsByClinicRequest"></a>

### ListDepartmentsByClinicRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-orgstructure-v1-ListDepartmentsByClinicResponse"></a>

### ListDepartmentsByClinicResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [DepartmentListItem](#query-orgstructure-v1-DepartmentListItem) | repeated |  |






<a name="query-orgstructure-v1-ListOrganizationsRequest"></a>

### ListOrganizationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-orgstructure-v1-ListOrganizationsResponse"></a>

### ListOrganizationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [OrganizationListItem](#query-orgstructure-v1-OrganizationListItem) | repeated |  |






<a name="query-orgstructure-v1-Organization"></a>

### Organization
Organization mirrors the projections.organizations row returned by
GetOrganization. Timestamps are RFC3339 strings.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| legal_address | [Address](#query-orgstructure-v1-Address) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="query-orgstructure-v1-OrganizationListItem"></a>

### OrganizationListItem
OrganizationListItem is the minimal shape returned by list endpoints.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| name | [string](#string) |  |  |






<a name="query-orgstructure-v1-Point"></a>

### Point
Point is the coordinate pair attached to an Address. Absent when the
projection row has neither longitude nor latitude.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| longitude | [double](#double) |  |  |
| latitude | [double](#double) |  |  |






<a name="query-orgstructure-v1-SearchOrganizationsRequest"></a>

### SearchOrganizationsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| query | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-orgstructure-v1-SearchOrganizationsResponse"></a>

### SearchOrganizationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [OrganizationListItem](#query-orgstructure-v1-OrganizationListItem) | repeated |  |












<a name="query-orgstructure-v1-OrgStructureQueryService"></a>

### OrgStructureQueryService
OrgStructureQueryService exposes read methods over the orgstructure
projections. Pagination is offset/limit; backed by COUNT(*) queries
via dedicated CountX RPCs — clients drive pagination controls with
both values.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetOrganization | [GetOrganizationRequest](#query-orgstructure-v1-GetOrganizationRequest) | [GetOrganizationResponse](#query-orgstructure-v1-GetOrganizationResponse) |  |
| ListOrganizations | [ListOrganizationsRequest](#query-orgstructure-v1-ListOrganizationsRequest) | [ListOrganizationsResponse](#query-orgstructure-v1-ListOrganizationsResponse) |  |
| CountOrganizations | [CountOrganizationsRequest](#query-orgstructure-v1-CountOrganizationsRequest) | [CountOrganizationsResponse](#query-orgstructure-v1-CountOrganizationsResponse) |  |
| SearchOrganizations | [SearchOrganizationsRequest](#query-orgstructure-v1-SearchOrganizationsRequest) | [SearchOrganizationsResponse](#query-orgstructure-v1-SearchOrganizationsResponse) |  |
| GetClinic | [GetClinicRequest](#query-orgstructure-v1-GetClinicRequest) | [GetClinicResponse](#query-orgstructure-v1-GetClinicResponse) |  |
| ListClinicsByOrganization | [ListClinicsByOrganizationRequest](#query-orgstructure-v1-ListClinicsByOrganizationRequest) | [ListClinicsByOrganizationResponse](#query-orgstructure-v1-ListClinicsByOrganizationResponse) |  |
| CountClinicsByOrganization | [CountClinicsByOrganizationRequest](#query-orgstructure-v1-CountClinicsByOrganizationRequest) | [CountClinicsByOrganizationResponse](#query-orgstructure-v1-CountClinicsByOrganizationResponse) |  |
| GetDepartment | [GetDepartmentRequest](#query-orgstructure-v1-GetDepartmentRequest) | [GetDepartmentResponse](#query-orgstructure-v1-GetDepartmentResponse) |  |
| ListDepartmentsByClinic | [ListDepartmentsByClinicRequest](#query-orgstructure-v1-ListDepartmentsByClinicRequest) | [ListDepartmentsByClinicResponse](#query-orgstructure-v1-ListDepartmentsByClinicResponse) |  |
| CountDepartmentsByClinic | [CountDepartmentsByClinicRequest](#query-orgstructure-v1-CountDepartmentsByClinicRequest) | [CountDepartmentsByClinicResponse](#query-orgstructure-v1-CountDepartmentsByClinicResponse) |  |





<a name="query_request_classifier_v1_classifier-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/request/classifier/v1/classifier.proto



<a name="query-request-classifier-v1-GetRequestTypeRequest"></a>

### GetRequestTypeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-request-classifier-v1-GetRequestTypeResponse"></a>

### GetRequestTypeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request_type | [RequestType](#query-request-classifier-v1-RequestType) |  |  |






<a name="query-request-classifier-v1-ListActiveRequestTypesByOrganizationRequest"></a>

### ListActiveRequestTypesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-request-classifier-v1-ListActiveRequestTypesByOrganizationResponse"></a>

### ListActiveRequestTypesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [RequestType](#query-request-classifier-v1-RequestType) | repeated |  |






<a name="query-request-classifier-v1-ListRequestTypesByOrganizationRequest"></a>

### ListRequestTypesByOrganizationRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-request-classifier-v1-ListRequestTypesByOrganizationResponse"></a>

### ListRequestTypesByOrganizationResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [RequestType](#query-request-classifier-v1-RequestType) | repeated |  |






<a name="query-request-classifier-v1-RequestType"></a>

### RequestType



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| name | [string](#string) |  |  |
| description | [string](#string) | optional |  |
| is_active | [bool](#bool) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |












<a name="query-request-classifier-v1-RequestClassifierQueryService"></a>

### RequestClassifierQueryService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetRequestType | [GetRequestTypeRequest](#query-request-classifier-v1-GetRequestTypeRequest) | [GetRequestTypeResponse](#query-request-classifier-v1-GetRequestTypeResponse) |  |
| ListRequestTypesByOrganization | [ListRequestTypesByOrganizationRequest](#query-request-classifier-v1-ListRequestTypesByOrganizationRequest) | [ListRequestTypesByOrganizationResponse](#query-request-classifier-v1-ListRequestTypesByOrganizationResponse) |  |
| ListActiveRequestTypesByOrganization | [ListActiveRequestTypesByOrganizationRequest](#query-request-classifier-v1-ListActiveRequestTypesByOrganizationRequest) | [ListActiveRequestTypesByOrganizationResponse](#query-request-classifier-v1-ListActiveRequestTypesByOrganizationResponse) |  |





<a name="query_request_v1_request-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/request/v1/request.proto



<a name="query-request-v1-Executor"></a>

### Executor



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| assigned_at | [string](#string) |  |  |
| assigned_by_id | [string](#string) |  |  |






<a name="query-request-v1-ExecutorHistoryEntry"></a>

### ExecutorHistoryEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| action | [string](#string) |  |  |
| employee_id | [string](#string) |  |  |
| employee_name | [string](#string) |  |  |
| actor_id | [string](#string) |  |  |
| actor_name | [string](#string) |  |  |
| changed_at | [string](#string) |  |  |






<a name="query-request-v1-GetServiceRequestHistoryRequest"></a>

### GetServiceRequestHistoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_request_id | [string](#string) |  |  |






<a name="query-request-v1-GetServiceRequestHistoryResponse"></a>

### GetServiceRequestHistoryResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status_history | [StatusHistoryEntry](#query-request-v1-StatusHistoryEntry) | repeated |  |
| executor_history | [ExecutorHistoryEntry](#query-request-v1-ExecutorHistoryEntry) | repeated |  |






<a name="query-request-v1-GetServiceRequestRequest"></a>

### GetServiceRequestRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-request-v1-GetServiceRequestResponse"></a>

### GetServiceRequestResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_request | [ServiceRequest](#query-request-v1-ServiceRequest) |  |  |






<a name="query-request-v1-ListServiceRequestsByIncidentRequest"></a>

### ListServiceRequestsByIncidentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| incident_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-request-v1-ListServiceRequestsByIncidentResponse"></a>

### ListServiceRequestsByIncidentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ServiceRequest](#query-request-v1-ServiceRequest) | repeated |  |






<a name="query-request-v1-ListServiceRequestsRequest"></a>

### ListServiceRequestsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| limit | [int32](#int32) |  |  |
| offset | [int32](#int32) |  |  |






<a name="query-request-v1-ListServiceRequestsResponse"></a>

### ListServiceRequestsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ServiceRequest](#query-request-v1-ServiceRequest) | repeated |  |






<a name="query-request-v1-ServiceRequest"></a>

### ServiceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| clinic_id | [string](#string) |  |  |
| department_id | [string](#string) |  |  |
| type_id | [string](#string) |  |  |
| incident_id | [string](#string) | optional |  |
| description | [string](#string) |  |  |
| status | [string](#string) |  |  |
| author_id | [string](#string) |  |  |
| author_display_name | [string](#string) |  |  |
| executors | [Executor](#query-request-v1-Executor) | repeated |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="query-request-v1-StatusHistoryEntry"></a>

### StatusHistoryEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| old_status | [string](#string) | optional |  |
| new_status | [string](#string) |  |  |
| actor_id | [string](#string) |  |  |
| actor_name | [string](#string) |  |  |
| changed_at | [string](#string) |  |  |












<a name="query-request-v1-ServiceRequestQueryService"></a>

### ServiceRequestQueryService


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetServiceRequest | [GetServiceRequestRequest](#query-request-v1-GetServiceRequestRequest) | [GetServiceRequestResponse](#query-request-v1-GetServiceRequestResponse) |  |
| ListServiceRequests | [ListServiceRequestsRequest](#query-request-v1-ListServiceRequestsRequest) | [ListServiceRequestsResponse](#query-request-v1-ListServiceRequestsResponse) |  |
| ListServiceRequestsByIncident | [ListServiceRequestsByIncidentRequest](#query-request-v1-ListServiceRequestsByIncidentRequest) | [ListServiceRequestsByIncidentResponse](#query-request-v1-ListServiceRequestsByIncidentResponse) |  |
| GetServiceRequestHistory | [GetServiceRequestHistoryRequest](#query-request-v1-GetServiceRequestHistoryRequest) | [GetServiceRequestHistoryResponse](#query-request-v1-GetServiceRequestHistoryResponse) |  |





<a name="query_self_v1_self-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/self/v1/self.proto



<a name="query-self-v1-GetMyClinicRoleRequest"></a>

### GetMyClinicRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |






<a name="query-self-v1-GetMyClinicRoleResponse"></a>

### GetMyClinicRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| is_clinic_head | [bool](#bool) |  |  |






<a name="query-self-v1-GetMyDepartmentRoleRequest"></a>

### GetMyDepartmentRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |






<a name="query-self-v1-GetMyDepartmentRoleResponse"></a>

### GetMyDepartmentRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| is_department_responsible | [bool](#bool) |  |  |






<a name="query-self-v1-GetMyEmploymentRequest"></a>

### GetMyEmploymentRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |






<a name="query-self-v1-GetMyEmploymentResponse"></a>

### GetMyEmploymentResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee | [query.membership.v1.EmployeeCardView](#query-membership-v1-EmployeeCardView) |  |  |






<a name="query-self-v1-GetMyIdentityRequest"></a>

### GetMyIdentityRequest







<a name="query-self-v1-GetMyIdentityResponse"></a>

### GetMyIdentityResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| is_system_admin | [bool](#bool) |  |  |






<a name="query-self-v1-GetMyOrganizationRoleRequest"></a>

### GetMyOrganizationRoleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |






<a name="query-self-v1-GetMyOrganizationRoleResponse"></a>

### GetMyOrganizationRoleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| is_org_admin | [bool](#bool) |  |  |
| is_org_head | [bool](#bool) |  |  |
| is_org_dispatcher | [bool](#bool) |  |  |






<a name="query-self-v1-ListMyOrganizationsRequest"></a>

### ListMyOrganizationsRequest







<a name="query-self-v1-ListMyOrganizationsResponse"></a>

### ListMyOrganizationsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [query.orgstructure.v1.OrganizationListItem](#query-orgstructure-v1-OrganizationListItem) | repeated |  |












<a name="query-self-v1-SelfQueryService"></a>

### SelfQueryService
SelfQueryService exposes read methods over the caller&#39;s own
memberships and roles. Every RPC is gated only by authentication
(valid Zitadel JWT) — no additional role check is performed inside
these handlers because the caller is querying data about themselves.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetMyIdentity | [GetMyIdentityRequest](#query-self-v1-GetMyIdentityRequest) | [GetMyIdentityResponse](#query-self-v1-GetMyIdentityResponse) | GetMyIdentity returns whether the caller is a system administrator. |
| ListMyOrganizations | [ListMyOrganizationsRequest](#query-self-v1-ListMyOrganizationsRequest) | [ListMyOrganizationsResponse](#query-self-v1-ListMyOrganizationsResponse) | ListMyOrganizations returns every organization where the caller has an active (non-terminated) employee record. |
| GetMyEmployment | [GetMyEmploymentRequest](#query-self-v1-GetMyEmploymentRequest) | [GetMyEmploymentResponse](#query-self-v1-GetMyEmploymentResponse) | GetMyEmployment returns the caller&#39;s employee card in the given organization. Returns NOT_FOUND when the caller is not an active employee of that organization. |
| GetMyOrganizationRole | [GetMyOrganizationRoleRequest](#query-self-v1-GetMyOrganizationRoleRequest) | [GetMyOrganizationRoleResponse](#query-self-v1-GetMyOrganizationRoleResponse) | GetMyOrganizationRole returns the caller&#39;s named roles in the given organization. Returns NOT_FOUND when the caller is not an active employee of that organization. |
| GetMyClinicRole | [GetMyClinicRoleRequest](#query-self-v1-GetMyClinicRoleRequest) | [GetMyClinicRoleResponse](#query-self-v1-GetMyClinicRoleResponse) | GetMyClinicRole returns whether the caller is the clinic head of the given clinic. Returns NOT_FOUND when the caller is not an active employee of that clinic. |
| GetMyDepartmentRole | [GetMyDepartmentRoleRequest](#query-self-v1-GetMyDepartmentRoleRequest) | [GetMyDepartmentRoleResponse](#query-self-v1-GetMyDepartmentRoleResponse) | GetMyDepartmentRole returns whether the caller holds the department-responsible role for the given department. Returns NOT_FOUND when the caller is not an active employee of that department. |





<a name="query_stats_v1_stats-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/stats/v1/stats.proto



<a name="query-stats-v1-ClinicStats"></a>

### ClinicStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |
| organization_id | [string](#string) |  |  |
| employees_total | [int64](#int64) |  |  |
| departments_total | [int64](#int64) |  |  |
| employees_on_vacation | [int64](#int64) |  |  |






<a name="query-stats-v1-DepartmentStats"></a>

### DepartmentStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |
| clinic_id | [string](#string) | optional |  |
| organization_id | [string](#string) |  |  |
| employees_total | [int64](#int64) |  |  |
| employees_on_vacation | [int64](#int64) |  |  |






<a name="query-stats-v1-GetClinicStatsRequest"></a>

### GetClinicStatsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| clinic_id | [string](#string) |  |  |






<a name="query-stats-v1-GetClinicStatsResponse"></a>

### GetClinicStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stats | [ClinicStats](#query-stats-v1-ClinicStats) |  |  |






<a name="query-stats-v1-GetDepartmentStatsRequest"></a>

### GetDepartmentStatsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |






<a name="query-stats-v1-GetDepartmentStatsResponse"></a>

### GetDepartmentStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stats | [DepartmentStats](#query-stats-v1-DepartmentStats) |  |  |






<a name="query-stats-v1-GetOrganizationStatsRequest"></a>

### GetOrganizationStatsRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |






<a name="query-stats-v1-GetOrganizationStatsResponse"></a>

### GetOrganizationStatsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| stats | [OrganizationStats](#query-stats-v1-OrganizationStats) |  |  |






<a name="query-stats-v1-OrganizationStats"></a>

### OrganizationStats



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| organization_id | [string](#string) |  |  |
| employees_total | [int64](#int64) |  |  |
| clinics_total | [int64](#int64) |  |  |
| departments_total | [int64](#int64) |  |  |
| employees_on_vacation | [int64](#int64) |  |  |
| vacations_scheduled | [int64](#int64) |  |  |












<a name="query-stats-v1-StatsQueryService"></a>

### StatsQueryService
StatsQueryService exposes aggregate counters for dashboards. Counts
come from the per-aggregate counters tables; &#34;employees on vacation&#34;
and &#34;scheduled vacations&#34; are computed at read time from
projections.employee_vacations so they never drift.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetOrganizationStats | [GetOrganizationStatsRequest](#query-stats-v1-GetOrganizationStatsRequest) | [GetOrganizationStatsResponse](#query-stats-v1-GetOrganizationStatsResponse) |  |
| GetClinicStats | [GetClinicStatsRequest](#query-stats-v1-GetClinicStatsRequest) | [GetClinicStatsResponse](#query-stats-v1-GetClinicStatsResponse) |  |
| GetDepartmentStats | [GetDepartmentStatsRequest](#query-stats-v1-GetDepartmentStatsRequest) | [GetDepartmentStatsResponse](#query-stats-v1-GetDepartmentStatsResponse) |  |





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

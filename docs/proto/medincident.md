# Protocol Documentation
<a name="top"></a>

## Table of Contents

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

- [query/identity/v1/identity.proto](#query_identity_v1_identity-proto)
    - [GetSessionRequest](#query-identity-v1-GetSessionRequest)
    - [GetSessionResponse](#query-identity-v1-GetSessionResponse)
    - [GetUserRequest](#query-identity-v1-GetUserRequest)
    - [GetUserResponse](#query-identity-v1-GetUserResponse)
    - [Session](#query-identity-v1-Session)
    - [User](#query-identity-v1-User)
    - [UserAgent](#query-identity-v1-UserAgent)
    - [UserAgent.HeadersEntry](#query-identity-v1-UserAgent-HeadersEntry)
    - [UserAgentHeaderValues](#query-identity-v1-UserAgentHeaderValues)

    - [IdentityQueryService](#query-identity-v1-IdentityQueryService)

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
    - [RoleHolder](#query-membership-v1-RoleHolder)
    - [SearchEmployeesByOrganizationRequest](#query-membership-v1-SearchEmployeesByOrganizationRequest)
    - [SearchEmployeesByOrganizationResponse](#query-membership-v1-SearchEmployeesByOrganizationResponse)
    - [SystemAdminView](#query-membership-v1-SystemAdminView)
    - [VacationView](#query-membership-v1-VacationView)

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





<a name="query_identity_v1_identity-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/identity/v1/identity.proto



<a name="query-identity-v1-GetSessionRequest"></a>

### GetSessionRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-identity-v1-GetSessionResponse"></a>

### GetSessionResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| session | [Session](#query-identity-v1-Session) |  |  |






<a name="query-identity-v1-GetUserRequest"></a>

### GetUserRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="query-identity-v1-GetUserResponse"></a>

### GetUserResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#query-identity-v1-User) |  |  |






<a name="query-identity-v1-Session"></a>

### Session
Session mirrors projections.sessions. user_id/user_resource_owner are
populated once a SessionUserChecked event is consumed.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| user_id | [string](#string) | optional |  |
| user_resource_owner | [string](#string) | optional |  |
| preferred_language | [string](#string) | optional |  |
| checked_at | [string](#string) | optional |  |
| user_agent | [UserAgent](#query-identity-v1-UserAgent) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="query-identity-v1-User"></a>

### User
User mirrors projections.users. gender uses the zitadel proto enum
integer: 0 unspecified, 1 female, 2 male, 3 diverse.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| user_name | [string](#string) |  |  |
| first_name | [string](#string) |  |  |
| last_name | [string](#string) |  |  |
| display_name | [string](#string) |  |  |
| nick_name | [string](#string) | optional |  |
| email | [string](#string) |  |  |
| email_verified | [bool](#bool) |  |  |
| preferred_language | [string](#string) |  |  |
| gender | [int32](#int32) |  |  |
| created_at | [string](#string) |  |  |
| updated_at | [string](#string) |  |  |






<a name="query-identity-v1-UserAgent"></a>

### UserAgent
UserAgent is the JSONB user_agent column unpacked to proto.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| ip | [string](#string) |  |  |
| headers | [UserAgent.HeadersEntry](#query-identity-v1-UserAgent-HeadersEntry) | repeated |  |
| fingerprint_id | [string](#string) | optional |  |
| description | [string](#string) | optional |  |






<a name="query-identity-v1-UserAgent-HeadersEntry"></a>

### UserAgent.HeadersEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [UserAgentHeaderValues](#query-identity-v1-UserAgentHeaderValues) |  |  |






<a name="query-identity-v1-UserAgentHeaderValues"></a>

### UserAgentHeaderValues
UserAgentHeaderValues holds the repeated-string values for a single
header entry.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| values | [string](#string) | repeated |  |












<a name="query-identity-v1-IdentityQueryService"></a>

### IdentityQueryService
IdentityQueryService exposes read methods over the Zitadel user and
session projections. Rows are populated asynchronously by the
Zitadel NATS consumer in query-server.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetUser | [GetUserRequest](#query-identity-v1-GetUserRequest) | [GetUserResponse](#query-identity-v1-GetUserResponse) |  |
| GetSession | [GetSessionRequest](#query-identity-v1-GetSessionRequest) | [GetSessionResponse](#query-identity-v1-GetSessionResponse) |  |





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
| holder | [RoleHolder](#query-membership-v1-RoleHolder) | optional |  |






<a name="query-membership-v1-GetDepartmentResponsibleRequest"></a>

### GetDepartmentResponsibleRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| department_id | [string](#string) |  |  |






<a name="query-membership-v1-GetDepartmentResponsibleResponse"></a>

### GetDepartmentResponsibleResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| holder | [RoleHolder](#query-membership-v1-RoleHolder) | optional |  |






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
| items | [RoleHolder](#query-membership-v1-RoleHolder) | repeated |  |






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
| items | [RoleHolder](#query-membership-v1-RoleHolder) | repeated |  |






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
| items | [RoleHolder](#query-membership-v1-RoleHolder) | repeated |  |






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






<a name="query-membership-v1-RoleHolder"></a>

### RoleHolder
RoleHolder represents a single role assignment. deputy_employee_id
stays empty when the holder has no deputy (or the row is itself a
deputy-less holder).


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| employee_id | [string](#string) |  |  |
| deputy_employee_id | [string](#string) | optional |  |






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

# Protocol Documentation
<a name="top"></a>

## Table of Contents

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
    - [ListTypesByCategoryRequest](#query-incident-classifier-v1-ListTypesByCategoryRequest)
    - [ListTypesByCategoryResponse](#query-incident-classifier-v1-ListTypesByCategoryResponse)
    - [Type](#query-incident-classifier-v1-Type)

    - [IncidentClassifierQueryService](#query-incident-classifier-v1-IncidentClassifierQueryService)

- [query/membership/v1/membership.proto](#query_membership_v1_membership-proto)
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
    - [SystemAdminView](#query-membership-v1-SystemAdminView)
    - [VacationView](#query-membership-v1-VacationView)

    - [MembershipQueryService](#query-membership-v1-MembershipQueryService)

- [query/orgstructure/v1/orgstructure.proto](#query_orgstructure_v1_orgstructure-proto)
    - [Address](#query-orgstructure-v1-Address)
    - [Clinic](#query-orgstructure-v1-Clinic)
    - [ClinicListItem](#query-orgstructure-v1-ClinicListItem)
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






<a name="query-incident-classifier-v1-ListTypesByCategoryRequest"></a>

### ListTypesByCategoryRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| category_id | [string](#string) |  |  |






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





<a name="query_membership_v1_membership-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## query/membership/v1/membership.proto



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






<a name="query-membership-v1-ListOrgHeadsResponse"></a>

### ListOrgHeadsResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [RoleHolder](#query-membership-v1-RoleHolder) | repeated |  |






<a name="query-membership-v1-ListSystemAdminsRequest"></a>

### ListSystemAdminsRequest







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
| ListVacationsByEmployee | [ListVacationsByEmployeeRequest](#query-membership-v1-ListVacationsByEmployeeRequest) | [ListVacationsByEmployeeResponse](#query-membership-v1-ListVacationsByEmployeeResponse) |  |
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
| GetClinic | [GetClinicRequest](#query-orgstructure-v1-GetClinicRequest) | [GetClinicResponse](#query-orgstructure-v1-GetClinicResponse) |  |
| ListClinicsByOrganization | [ListClinicsByOrganizationRequest](#query-orgstructure-v1-ListClinicsByOrganizationRequest) | [ListClinicsByOrganizationResponse](#query-orgstructure-v1-ListClinicsByOrganizationResponse) |  |
| GetDepartment | [GetDepartmentRequest](#query-orgstructure-v1-GetDepartmentRequest) | [GetDepartmentResponse](#query-orgstructure-v1-GetDepartmentResponse) |  |
| ListDepartmentsByClinic | [ListDepartmentsByClinicRequest](#query-orgstructure-v1-ListDepartmentsByClinicRequest) | [ListDepartmentsByClinicResponse](#query-orgstructure-v1-ListDepartmentsByClinicResponse) |  |





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

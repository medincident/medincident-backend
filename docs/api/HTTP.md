<!-- Generator: Widdershins v4.0.1 -->

<h1 id="command-announcement-v1-announcement-proto">command/announcement/v1/announcement.proto version not set</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

<h1 id="command-announcement-v1-announcement-proto-announcementcommandservice">AnnouncementCommandService</h1>

## AnnouncementCommandService_CreateAnnouncement

<a id="opIdAnnouncementCommandService_CreateAnnouncement"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/announcements \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/announcements`

> Body parameter

```json
{
  "organizationId": "string",
  "clinicId": "string",
  "departmentId": "string",
  "title": "string",
  "content": "string",
  "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
  "startsAt": "string",
  "endsAt": "string"
}
```

<h3 id="announcementcommandservice_createannouncement-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[v1CreateAnnouncementRequest](#schemav1createannouncementrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "id": "string"
}
```

> Not found. Error codes:
- `announcement_organization_not_found` — organization with the given ID does not exist.
- `announcement_clinic_not_found` — clinic with the given ID does not exist.
- `announcement_department_not_found` — department with the given ID does not exist.

```json
{
  "code": "announcement_organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="announcementcommandservice_createannouncement-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateAnnouncementResponse](#schemav1createannouncementresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `announcement_organization_not_found` — organization with the given ID does not exist.
- `announcement_clinic_not_found` — clinic with the given ID does not exist.
- `announcement_department_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnnouncementCommandService_UpdateAnnouncement

<a id="opIdAnnouncementCommandService_UpdateAnnouncement"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/announcements/{id} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/announcements/{id}`

> Body parameter

```json
{
  "title": "string",
  "content": "string",
  "startsAt": "string",
  "endsAt": "string"
}
```

<h3 id="announcementcommandservice_updateannouncement-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|
|body|body|[AnnouncementCommandServiceUpdateAnnouncementBody](#schemaannouncementcommandserviceupdateannouncementbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Failed precondition. Error codes:
- `announcement_archived` — announcement is archived and cannot be modified.

```json
{
  "code": "announcement_archived",
  "message": "Announcement is archived."
}
```

> Not found. Error codes:
- `announcement_not_found` — announcement with the given ID does not exist.

```json
{
  "code": "announcement_not_found",
  "message": "Announcement not found."
}
```

<h3 id="announcementcommandservice_updateannouncement-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateAnnouncementResponse](#schemav1updateannouncementresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Failed precondition. Error codes:
- `announcement_archived` — announcement is archived and cannot be modified.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `announcement_not_found` — announcement with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnnouncementCommandService_UpdateAnnouncementPriority

<a id="opIdAnnouncementCommandService_UpdateAnnouncementPriority"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/announcements/{id}/priority \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/announcements/{id}/priority`

> Body parameter

```json
{
  "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED"
}
```

<h3 id="announcementcommandservice_updateannouncementpriority-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|
|body|body|[AnnouncementCommandServiceUpdateAnnouncementPriorityBody](#schemaannouncementcommandserviceupdateannouncementprioritybody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Failed precondition. Error codes:
- `announcement_archived` — announcement is archived and cannot be modified.

```json
{
  "code": "announcement_archived",
  "message": "Announcement is archived."
}
```

> Not found. Error codes:
- `announcement_not_found` — announcement with the given ID does not exist.

```json
{
  "code": "announcement_not_found",
  "message": "Announcement not found."
}
```

<h3 id="announcementcommandservice_updateannouncementpriority-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateAnnouncementPriorityResponse](#schemav1updateannouncementpriorityresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Failed precondition. Error codes:
- `announcement_archived` — announcement is archived and cannot be modified.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `announcement_not_found` — announcement with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnnouncementCommandService_ArchiveAnnouncement

<a id="opIdAnnouncementCommandService_ArchiveAnnouncement"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/announcements/{id}:archive \
  -H 'Accept: application/json'

```

`POST /v1/announcements/{id}:archive`

<h3 id="announcementcommandservice_archiveannouncement-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `announcement_not_found` — announcement with the given ID does not exist.

```json
{
  "code": "announcement_not_found",
  "message": "Announcement not found."
}
```

<h3 id="announcementcommandservice_archiveannouncement-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ArchiveAnnouncementResponse](#schemav1archiveannouncementresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `announcement_not_found` — announcement with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnnouncementCommandService_UnarchiveAnnouncement

<a id="opIdAnnouncementCommandService_UnarchiveAnnouncement"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/announcements/{id}:unarchive \
  -H 'Accept: application/json'

```

`POST /v1/announcements/{id}:unarchive`

<h3 id="announcementcommandservice_unarchiveannouncement-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `announcement_not_found` — announcement with the given ID does not exist.

```json
{
  "code": "announcement_not_found",
  "message": "Announcement not found."
}
```

<h3 id="announcementcommandservice_unarchiveannouncement-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UnarchiveAnnouncementResponse](#schemav1unarchiveannouncementresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `announcement_not_found` — announcement with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-incidentbuffercommandservice">IncidentBufferCommandService</h1>

## IncidentBufferCommandService_SubmitPatientIncident

<a id="opIdIncidentBufferCommandService_SubmitPatientIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/patient-incidents \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/patient-incidents`

> Body parameter

```json
{
  "organizationId": "string",
  "categoryId": "string",
  "typeId": "string",
  "description": "string",
  "occurredAt": "string"
}
```

<h3 id="incidentbuffercommandservice_submitpatientincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[v1SubmitPatientIncidentRequest](#schemav1submitpatientincidentrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "bufferId": "string"
}
```

> Validation failed or precondition not met. Error codes:
- `buffer_type_not_allowed_for_patients` — the selected incident type is not available for patient submissions.

```json
{
  "code": "buffer_type_not_allowed_for_patients",
  "message": "Incident type is not allowed for patient submissions."
}
```

> Not found. Error codes:
- `buffer_organization_not_found` — organization with the given ID does not exist.

```json
{
  "code": "buffer_organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="incidentbuffercommandservice_submitpatientincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1SubmitPatientIncidentResponse](#schemav1submitpatientincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `buffer_type_not_allowed_for_patients` — the selected incident type is not available for patient submissions.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `buffer_organization_not_found` — organization with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentBufferCommandService_UpdatePatientIncident

<a id="opIdIncidentBufferCommandService_UpdatePatientIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/patient-incidents/{bufferId} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/patient-incidents/{bufferId}`

> Body parameter

```json
{
  "categoryId": "string",
  "typeId": "string",
  "description": "string",
  "occurredAt": "string"
}
```

<h3 id="incidentbuffercommandservice_updatepatientincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|bufferId|path|string|true|none|
|body|body|[IncidentBufferCommandServiceUpdatePatientIncidentBody](#schemaincidentbuffercommandserviceupdatepatientincidentbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `buffer_not_pending` — patient incident is not in pending status and cannot be updated.

```json
{
  "code": "buffer_not_pending",
  "message": "Patient incident is not in pending status."
}
```

> Permission denied. Error codes:
- `buffer_not_patient_owner` — caller is not the patient who submitted this incident.

```json
{
  "code": "buffer_not_patient_owner",
  "message": "Not the owner of this patient incident."
}
```

> Not found. Error codes:
- `buffer_not_found` — patient incident with the given ID does not exist.

```json
{
  "code": "buffer_not_found",
  "message": "Patient incident not found."
}
```

<h3 id="incidentbuffercommandservice_updatepatientincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdatePatientIncidentResponse](#schemav1updatepatientincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `buffer_not_pending` — patient incident is not in pending status and cannot be updated.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied. Error codes:
- `buffer_not_patient_owner` — caller is not the patient who submitted this incident.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `buffer_not_found` — patient incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentBufferCommandService_CancelPatientIncident

<a id="opIdIncidentBufferCommandService_CancelPatientIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/patient-incidents/{bufferId}:cancel \
  -H 'Accept: application/json'

```

`POST /v1/patient-incidents/{bufferId}:cancel`

<h3 id="incidentbuffercommandservice_cancelpatientincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|bufferId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `buffer_not_pending` — patient incident is not in pending status and cannot be cancelled.

```json
{
  "code": "buffer_not_pending",
  "message": "Patient incident is not in pending status."
}
```

> Permission denied. Error codes:
- `buffer_not_patient_owner` — caller is not the patient who submitted this incident.

```json
{
  "code": "buffer_not_patient_owner",
  "message": "Not the owner of this patient incident."
}
```

> Not found. Error codes:
- `buffer_not_found` — patient incident with the given ID does not exist.

```json
{
  "code": "buffer_not_found",
  "message": "Patient incident not found."
}
```

<h3 id="incidentbuffercommandservice_cancelpatientincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CancelPatientIncidentResponse](#schemav1cancelpatientincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `buffer_not_pending` — patient incident is not in pending status and cannot be cancelled.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied. Error codes:
- `buffer_not_patient_owner` — caller is not the patient who submitted this incident.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `buffer_not_found` — patient incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentBufferCommandService_PublishPatientIncident

<a id="opIdIncidentBufferCommandService_PublishPatientIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/patient-incidents/{bufferId}:publish \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/patient-incidents/{bufferId}:publish`

> Body parameter

```json
{
  "departmentId": "string",
  "categoryId": "string",
  "typeId": "string",
  "description": "string"
}
```

<h3 id="incidentbuffercommandservice_publishpatientincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|bufferId|path|string|true|none|
|body|body|[IncidentBufferCommandServicePublishPatientIncidentBody](#schemaincidentbuffercommandservicepublishpatientincidentbody)|true|none|

> Example responses

> 200 Response

```json
{
  "incidentId": "string"
}
```

> Validation failed or precondition not met. Error codes:
- `buffer_not_pending` — patient incident is not in pending status and cannot be published.

```json
{
  "code": "buffer_not_pending",
  "message": "Patient incident is not in pending status."
}
```

> Not found. Error codes:
- `buffer_not_found` — patient incident with the given ID does not exist.
- `buffer_department_not_found` — department with the given ID does not exist.
- `buffer_category_not_found` — category with the given ID does not exist.
- `buffer_dispatcher_not_found` — dispatcher employee not found.

```json
{
  "code": "buffer_not_found",
  "message": "Patient incident not found."
}
```

<h3 id="incidentbuffercommandservice_publishpatientincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1PublishPatientIncidentResponse](#schemav1publishpatientincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `buffer_not_pending` — patient incident is not in pending status and cannot be published.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `buffer_not_found` — patient incident with the given ID does not exist.
- `buffer_department_not_found` — department with the given ID does not exist.
- `buffer_category_not_found` — category with the given ID does not exist.
- `buffer_dispatcher_not_found` — dispatcher employee not found.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentBufferCommandService_RejectPatientIncident

<a id="opIdIncidentBufferCommandService_RejectPatientIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/patient-incidents/{bufferId}:reject \
  -H 'Accept: application/json'

```

`POST /v1/patient-incidents/{bufferId}:reject`

<h3 id="incidentbuffercommandservice_rejectpatientincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|bufferId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `buffer_not_pending` — patient incident is not in pending status and cannot be rejected.

```json
{
  "code": "buffer_not_pending",
  "message": "Patient incident is not in pending status."
}
```

> Not found. Error codes:
- `buffer_not_found` — patient incident with the given ID does not exist.

```json
{
  "code": "buffer_not_found",
  "message": "Patient incident not found."
}
```

<h3 id="incidentbuffercommandservice_rejectpatientincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RejectPatientIncidentResponse](#schemav1rejectpatientincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `buffer_not_pending` — patient incident is not in pending status and cannot be rejected.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `buffer_not_found` — patient incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-incidentclassifiercommandservice">IncidentClassifierCommandService</h1>

## IncidentClassifierCommandService_DeleteIncidentCategory

<a id="opIdIncidentClassifierCommandService_DeleteIncidentCategory"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/incident-categories/{categoryId} \
  -H 'Accept: application/json'

```

`DELETE /v1/incident-categories/{categoryId}`

<h3 id="incidentclassifiercommandservice_deleteincidentcategory-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|categoryId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.

```json
{
  "code": "incident_category_not_found",
  "message": "Category not found."
}
```

<h3 id="incidentclassifiercommandservice_deleteincidentcategory-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1DeleteIncidentCategoryResponse](#schemav1deleteincidentcategoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_DeactivateIncidentCategory

<a id="opIdIncidentClassifierCommandService_DeactivateIncidentCategory"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incident-categories/{categoryId}/deactivations \
  -H 'Accept: application/json'

```

`POST /v1/incident-categories/{categoryId}/deactivations`

<h3 id="incidentclassifiercommandservice_deactivateincidentcategory-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|categoryId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.

```json
{
  "code": "incident_category_not_found",
  "message": "Category not found."
}
```

<h3 id="incidentclassifiercommandservice_deactivateincidentcategory-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1DeactivateIncidentCategoryResponse](#schemav1deactivateincidentcategoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_UpdateIncidentCategoryDetails

<a id="opIdIncidentClassifierCommandService_UpdateIncidentCategoryDetails"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/incident-categories/{categoryId}/details \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/incident-categories/{categoryId}/details`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="incidentclassifiercommandservice_updateincidentcategorydetails-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|categoryId|path|string|true|none|
|body|body|[IncidentClassifierCommandServiceUpdateIncidentCategoryDetailsBody](#schemaincidentclassifiercommandserviceupdateincidentcategorydetailsbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.

```json
{
  "code": "incident_category_not_found",
  "message": "Category not found."
}
```

> Conflict. Error codes:
- `incident_category_name_conflict` — a category with this name already exists in the same scope.

```json
{
  "code": "incident_category_name_conflict",
  "message": "Category name already exists."
}
```

<h3 id="incidentclassifiercommandservice_updateincidentcategorydetails-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateIncidentCategoryDetailsResponse](#schemav1updateincidentcategorydetailsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `incident_category_name_conflict` — a category with this name already exists in the same scope.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_ReactivateIncidentCategory

<a id="opIdIncidentClassifierCommandService_ReactivateIncidentCategory"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incident-categories/{categoryId}/reactivations \
  -H 'Accept: application/json'

```

`POST /v1/incident-categories/{categoryId}/reactivations`

<h3 id="incidentclassifiercommandservice_reactivateincidentcategory-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|categoryId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `incident_category_reactivate_inactive_ancestor` — an ancestor category is inactive and must be reactivated first.

```json
{
  "code": "incident_category_reactivate_inactive_ancestor",
  "message": "Ancestor category is inactive."
}
```

> Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.

```json
{
  "code": "incident_category_not_found",
  "message": "Category not found."
}
```

> Conflict. Error codes:
- `incident_category_reactivate_name_conflict` — reactivation would create a name conflict with an active category.

```json
{
  "code": "incident_category_reactivate_name_conflict",
  "message": "Name conflict on reactivation."
}
```

<h3 id="incidentclassifiercommandservice_reactivateincidentcategory-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ReactivateIncidentCategoryResponse](#schemav1reactivateincidentcategoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_category_reactivate_inactive_ancestor` — an ancestor category is inactive and must be reactivated first.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `incident_category_reactivate_name_conflict` — reactivation would create a name conflict with an active category.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## --- Types ---

<a id="opIdIncidentClassifierCommandService_CreateIncidentType"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incident-categories/{categoryId}/types \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/incident-categories/{categoryId}/types`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="----types-----parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|categoryId|path|string|true|none|
|body|body|[IncidentClassifierCommandServiceCreateIncidentTypeBody](#schemaincidentclassifiercommandservicecreateincidenttypebody)|true|none|

> Example responses

> 200 Response

```json
{
  "typeId": "string"
}
```

> Validation failed or precondition not met. Error codes:
- `incident_category_inactive` — the category is inactive; new types cannot be created in an inactive category.

```json
{
  "code": "incident_category_inactive",
  "message": "Category is inactive."
}
```

> Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.

```json
{
  "code": "incident_category_not_found",
  "message": "Category not found."
}
```

> Conflict. Error codes:
- `incident_type_name_conflict` — a type with this name already exists in the same category.

```json
{
  "code": "incident_type_name_conflict",
  "message": "Type name already exists."
}
```

<h3 id="----types-----responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateIncidentTypeResponse](#schemav1createincidenttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_category_inactive` — the category is inactive; new types cannot be created in an inactive category.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `incident_type_name_conflict` — a type with this name already exists in the same category.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_MoveIncidentCategory

<a id="opIdIncidentClassifierCommandService_MoveIncidentCategory"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incident-categories/{categoryId}:move \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/incident-categories/{categoryId}:move`

> Body parameter

```json
{
  "newParentCategoryId": "string"
}
```

<h3 id="incidentclassifiercommandservice_moveincidentcategory-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|categoryId|path|string|true|none|
|body|body|[IncidentClassifierCommandServiceMoveIncidentCategoryBody](#schemaincidentclassifiercommandservicemoveincidentcategorybody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `incident_category_move_organization_mismatch` — target parent belongs to a different organization.
- `incident_category_move_would_create_cycle` — move would create a cycle in the category tree.
- `incident_category_move_would_exceed_depth` — move would exceed the maximum category depth.
- `incident_category_parent_inactive` — target parent category is inactive.

```json
{
  "code": "incident_category_move_would_create_cycle",
  "message": "Move would create a cycle."
}
```

> Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.
- `incident_category_parent_not_found` — target parent category with the given ID does not exist.

```json
{
  "code": "incident_category_not_found",
  "message": "Category not found."
}
```

<h3 id="incidentclassifiercommandservice_moveincidentcategory-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1MoveIncidentCategoryResponse](#schemav1moveincidentcategoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_category_move_organization_mismatch` — target parent belongs to a different organization.
- `incident_category_move_would_create_cycle` — move would create a cycle in the category tree.
- `incident_category_move_would_exceed_depth` — move would exceed the maximum category depth.
- `incident_category_parent_inactive` — target parent category is inactive.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.
- `incident_category_parent_not_found` — target parent category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_DeleteIncidentType

<a id="opIdIncidentClassifierCommandService_DeleteIncidentType"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/incident-types/{typeId} \
  -H 'Accept: application/json'

```

`DELETE /v1/incident-types/{typeId}`

<h3 id="incidentclassifiercommandservice_deleteincidenttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.

```json
{
  "code": "incident_type_not_found",
  "message": "Type not found."
}
```

<h3 id="incidentclassifiercommandservice_deleteincidenttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1DeleteIncidentTypeResponse](#schemav1deleteincidenttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_DeactivateIncidentType

<a id="opIdIncidentClassifierCommandService_DeactivateIncidentType"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incident-types/{typeId}/deactivations \
  -H 'Accept: application/json'

```

`POST /v1/incident-types/{typeId}/deactivations`

<h3 id="incidentclassifiercommandservice_deactivateincidenttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.

```json
{
  "code": "incident_type_not_found",
  "message": "Type not found."
}
```

<h3 id="incidentclassifiercommandservice_deactivateincidenttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1DeactivateIncidentTypeResponse](#schemav1deactivateincidenttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_UpdateIncidentTypeDetails

<a id="opIdIncidentClassifierCommandService_UpdateIncidentTypeDetails"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/incident-types/{typeId}/details \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/incident-types/{typeId}/details`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="incidentclassifiercommandservice_updateincidenttypedetails-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|
|body|body|[IncidentClassifierCommandServiceUpdateIncidentTypeDetailsBody](#schemaincidentclassifiercommandserviceupdateincidenttypedetailsbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.

```json
{
  "code": "incident_type_not_found",
  "message": "Type not found."
}
```

> Conflict. Error codes:
- `incident_type_name_conflict` — a type with this name already exists in the same category.

```json
{
  "code": "incident_type_name_conflict",
  "message": "Type name already exists."
}
```

<h3 id="incidentclassifiercommandservice_updateincidenttypedetails-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateIncidentTypeDetailsResponse](#schemav1updateincidenttypedetailsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `incident_type_name_conflict` — a type with this name already exists in the same category.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_DisallowIncidentTypeForPatients

<a id="opIdIncidentClassifierCommandService_DisallowIncidentTypeForPatients"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/incident-types/{typeId}/patient-allowances \
  -H 'Accept: application/json'

```

`DELETE /v1/incident-types/{typeId}/patient-allowances`

<h3 id="incidentclassifiercommandservice_disallowincidenttypeforpatients-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.

```json
{
  "code": "incident_type_not_found",
  "message": "Type not found."
}
```

<h3 id="incidentclassifiercommandservice_disallowincidenttypeforpatients-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1DisallowIncidentTypeForPatientsResponse](#schemav1disallowincidenttypeforpatientsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_AllowIncidentTypeForPatients

<a id="opIdIncidentClassifierCommandService_AllowIncidentTypeForPatients"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incident-types/{typeId}/patient-allowances \
  -H 'Accept: application/json'

```

`POST /v1/incident-types/{typeId}/patient-allowances`

<h3 id="incidentclassifiercommandservice_allowincidenttypeforpatients-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.

```json
{
  "code": "incident_type_not_found",
  "message": "Type not found."
}
```

<h3 id="incidentclassifiercommandservice_allowincidenttypeforpatients-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AllowIncidentTypeForPatientsResponse](#schemav1allowincidenttypeforpatientsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_ReactivateIncidentType

<a id="opIdIncidentClassifierCommandService_ReactivateIncidentType"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incident-types/{typeId}/reactivations \
  -H 'Accept: application/json'

```

`POST /v1/incident-types/{typeId}/reactivations`

<h3 id="incidentclassifiercommandservice_reactivateincidenttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `incident_type_reactivate_inactive_ancestor` — the parent category is inactive and must be reactivated first.

```json
{
  "code": "incident_type_reactivate_inactive_ancestor",
  "message": "Parent category is inactive."
}
```

> Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.

```json
{
  "code": "incident_type_not_found",
  "message": "Type not found."
}
```

> Conflict. Error codes:
- `incident_type_reactivate_name_conflict` — reactivation would create a name conflict with an active type.

```json
{
  "code": "incident_type_reactivate_name_conflict",
  "message": "Name conflict on reactivation."
}
```

<h3 id="incidentclassifiercommandservice_reactivateincidenttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ReactivateIncidentTypeResponse](#schemav1reactivateincidenttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_type_reactivate_inactive_ancestor` — the parent category is inactive and must be reactivated first.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `incident_type_reactivate_name_conflict` — reactivation would create a name conflict with an active type.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierCommandService_MoveIncidentType

<a id="opIdIncidentClassifierCommandService_MoveIncidentType"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incident-types/{typeId}:move \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/incident-types/{typeId}:move`

> Body parameter

```json
{
  "newCategoryId": "string"
}
```

<h3 id="incidentclassifiercommandservice_moveincidenttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|
|body|body|[IncidentClassifierCommandServiceMoveIncidentTypeBody](#schemaincidentclassifiercommandservicemoveincidenttypebody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `incident_type_move_organization_mismatch` — target category belongs to a different organization.
- `incident_type_category_inactive` — target category is inactive.
- `incident_category_inactive` — target category is inactive.

```json
{
  "code": "incident_type_move_organization_mismatch",
  "message": "Target category belongs to a different organization."
}
```

> Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.
- `incident_type_category_not_found` — target category with the given ID does not exist.

```json
{
  "code": "incident_type_not_found",
  "message": "Type not found."
}
```

<h3 id="incidentclassifiercommandservice_moveincidenttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1MoveIncidentTypeResponse](#schemav1moveincidenttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_type_move_organization_mismatch` — target category belongs to a different organization.
- `incident_type_category_inactive` — target category is inactive.
- `incident_category_inactive` — target category is inactive.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.
- `incident_type_category_not_found` — target category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## --- Categories ---

<a id="opIdIncidentClassifierCommandService_CreateIncidentCategory"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/incident-categories \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/incident-categories`

> Body parameter

```json
{
  "parentCategoryId": "string",
  "name": "string",
  "description": "string"
}
```

<h3 id="----categories-----parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|body|body|[IncidentClassifierCommandServiceCreateIncidentCategoryBody](#schemaincidentclassifiercommandservicecreateincidentcategorybody)|true|none|

> Example responses

> 200 Response

```json
{
  "categoryId": "string"
}
```

> Validation failed or precondition not met. Error codes:
- `incident_category_parent_inactive` — parent category is inactive.
- `incident_category_max_depth_exceeded` — category tree depth limit (5) would be exceeded.

```json
{
  "code": "incident_category_parent_inactive",
  "message": "Parent category is inactive."
}
```

> Not found. Error codes:
- `incident_category_not_found` — parent category with the given ID does not exist.

```json
{
  "code": "incident_category_not_found",
  "message": "Parent category not found."
}
```

> Conflict. Error codes:
- `incident_category_name_conflict` — a category with this name already exists in the same scope.

```json
{
  "code": "incident_category_name_conflict",
  "message": "Category name already exists."
}
```

<h3 id="----categories-----responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateIncidentCategoryResponse](#schemav1createincidentcategoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_category_parent_inactive` — parent category is inactive.
- `incident_category_max_depth_exceeded` — category tree depth limit (5) would be exceeded.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_category_not_found` — parent category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `incident_category_name_conflict` — a category with this name already exists in the same scope.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-incidentcommandservice">IncidentCommandService</h1>

## IncidentCommandService_CreateIncident

<a id="opIdIncidentCommandService_CreateIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incidents \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/incidents`

> Body parameter

```json
{
  "departmentId": "string",
  "categoryId": "string",
  "typeId": "string",
  "description": "string",
  "occurredAt": "string"
}
```

<h3 id="incidentcommandservice_createincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[v1CreateIncidentRequest](#schemav1createincidentrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "incidentId": "string"
}
```

> Validation failed or precondition not met. Error codes:
- `incident_category_inactive` — the selected category is inactive.
- `incident_type_inactive` — the selected type is inactive.

```json
{
  "code": "incident_category_inactive",
  "message": "Incident category is inactive."
}
```

> Not found. Error codes:
- `incident_department_not_found` — department with the given ID does not exist.
- `incident_category_not_found` — incident category with the given ID does not exist.
- `incident_type_not_found` — incident type with the given ID does not exist.
- `incident_employee_not_found` — registrar employee not found.

```json
{
  "code": "incident_department_not_found",
  "message": "Department not found."
}
```

<h3 id="incidentcommandservice_createincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateIncidentResponse](#schemav1createincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_category_inactive` — the selected category is inactive.
- `incident_type_inactive` — the selected type is inactive.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_department_not_found` — department with the given ID does not exist.
- `incident_category_not_found` — incident category with the given ID does not exist.
- `incident_type_not_found` — incident type with the given ID does not exist.
- `incident_employee_not_found` — registrar employee not found.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentCommandService_UpdateIncidentDescription

<a id="opIdIncidentCommandService_UpdateIncidentDescription"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/incidents/{incidentId}/description \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/incidents/{incidentId}/description`

> Body parameter

```json
{
  "description": "string"
}
```

<h3 id="incidentcommandservice_updateincidentdescription-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|incidentId|path|string|true|none|
|body|body|[IncidentCommandServiceUpdateIncidentDescriptionBody](#schemaincidentcommandserviceupdateincidentdescriptionbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `incident_frozen` — incident is frozen and cannot be modified.

```json
{
  "code": "incident_frozen",
  "message": "Incident is frozen."
}
```

> Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.

```json
{
  "code": "incident_not_found",
  "message": "Incident not found."
}
```

<h3 id="incidentcommandservice_updateincidentdescription-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateIncidentDescriptionResponse](#schemav1updateincidentdescriptionresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_frozen` — incident is frozen and cannot be modified.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentCommandService_UpdateIncidentPriority

<a id="opIdIncidentCommandService_UpdateIncidentPriority"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/incidents/{incidentId}/priority \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/incidents/{incidentId}/priority`

> Body parameter

```json
{
  "priority": "INCIDENT_PRIORITY_UNSPECIFIED"
}
```

<h3 id="incidentcommandservice_updateincidentpriority-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|incidentId|path|string|true|none|
|body|body|[IncidentCommandServiceUpdateIncidentPriorityBody](#schemaincidentcommandserviceupdateincidentprioritybody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `incident_frozen` — incident is frozen and cannot be modified.

```json
{
  "code": "incident_frozen",
  "message": "Incident is frozen."
}
```

> Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.

```json
{
  "code": "incident_not_found",
  "message": "Incident not found."
}
```

<h3 id="incidentcommandservice_updateincidentpriority-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateIncidentPriorityResponse](#schemav1updateincidentpriorityresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_frozen` — incident is frozen and cannot be modified.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentCommandService_UpdateIncidentStatus

<a id="opIdIncidentCommandService_UpdateIncidentStatus"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/incidents/{incidentId}/status \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/incidents/{incidentId}/status`

> Body parameter

```json
{
  "newStatus": "INCIDENT_STATUS_UNSPECIFIED"
}
```

<h3 id="incidentcommandservice_updateincidentstatus-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|incidentId|path|string|true|none|
|body|body|[IncidentCommandServiceUpdateIncidentStatusBody](#schemaincidentcommandserviceupdateincidentstatusbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `incident_invalid_status_transition` — the requested status transition is not allowed.
- `incident_frozen` — incident is frozen and cannot be modified.

```json
{
  "code": "incident_invalid_status_transition",
  "message": "Invalid status transition."
}
```

> Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.

```json
{
  "code": "incident_not_found",
  "message": "Incident not found."
}
```

<h3 id="incidentcommandservice_updateincidentstatus-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateIncidentStatusResponse](#schemav1updateincidentstatusresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_invalid_status_transition` — the requested status transition is not allowed.
- `incident_frozen` — incident is frozen and cannot be modified.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentCommandService_CancelIncident

<a id="opIdIncidentCommandService_CancelIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incidents/{incidentId}:cancel \
  -H 'Accept: application/json'

```

`POST /v1/incidents/{incidentId}:cancel`

<h3 id="incidentcommandservice_cancelincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|incidentId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Validation failed or precondition not met. Error codes:
- `incident_not_cancellable` — incident is not in a cancellable status.

```json
{
  "code": "incident_not_cancellable",
  "message": "Incident cannot be cancelled in its current status."
}
```

> Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.

```json
{
  "code": "incident_not_found",
  "message": "Incident not found."
}
```

<h3 id="incidentcommandservice_cancelincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CancelIncidentResponse](#schemav1cancelincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_not_cancellable` — incident is not in a cancellable status.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentCommandService_ReopenIncident

<a id="opIdIncidentCommandService_ReopenIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/incidents/{incidentId}:reopen \
  -H 'Accept: application/json'

```

`POST /v1/incidents/{incidentId}:reopen`

<h3 id="incidentcommandservice_reopenincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|incidentId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "reopenedIncidentId": "string"
}
```

> Validation failed or precondition not met. Error codes:
- `incident_not_reopenable` — incident cannot be reopened from its current status.

```json
{
  "code": "incident_not_reopenable",
  "message": "Incident cannot be reopened."
}
```

> Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.

```json
{
  "code": "incident_not_found",
  "message": "Incident not found."
}
```

<h3 id="incidentcommandservice_reopenincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ReopenIncidentResponse](#schemav1reopenincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or precondition not met. Error codes:
- `incident_not_reopenable` — incident cannot be reopened from its current status.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_not_found` — incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-membershipcommandservice">MembershipCommandService</h1>

## ------ ClinicHead ------

<a id="opIdMembershipCommandService_AssignClinicHead"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/clinics/{clinicId}/heads \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/clinics/{clinicId}/heads`

> Body parameter

```json
{
  "employeeId": "string"
}
```

<h3 id="-------clinichead--------parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignClinicHeadBody](#schemamembershipcommandserviceassignclinicheadbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "clinic_not_found",
  "message": "Clinic not found."
}
```

> Conflict. Error codes:
- `clinic_head_already_assigned` — this clinic already has a head assigned.

```json
{
  "code": "clinic_head_already_assigned",
  "message": "Clinic head already assigned."
}
```

<h3 id="-------clinichead--------responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignClinicHeadResponse](#schemav1assignclinicheadresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `clinic_head_already_assigned` — this clinic already has a head assigned.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RevokeClinicHead

<a id="opIdMembershipCommandService_RevokeClinicHead"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/clinics/{clinicId}/heads/{employeeId} \
  -H 'Accept: application/json'

```

`DELETE /v1/clinics/{clinicId}/heads/{employeeId}`

<h3 id="membershipcommandservice_revokeclinichead-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `clinic_head_not_found` — this employee is not the clinic head.

```json
{
  "code": "clinic_head_not_found",
  "message": "Clinic head not found."
}
```

<h3 id="membershipcommandservice_revokeclinichead-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RevokeClinicHeadResponse](#schemav1revokeclinicheadresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `clinic_head_not_found` — this employee is not the clinic head.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RemoveClinicHeadDeputy

<a id="opIdMembershipCommandService_RemoveClinicHeadDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/clinics/{clinicId}/heads/{employeeId}/deputy \
  -H 'Accept: application/json'

```

`DELETE /v1/clinics/{clinicId}/heads/{employeeId}/deputy`

<h3 id="membershipcommandservice_removeclinicheaddeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.

```json
{
  "code": "clinic_not_found",
  "message": "Clinic not found."
}
```

<h3 id="membershipcommandservice_removeclinicheaddeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RemoveClinicHeadDeputyResponse](#schemav1removeclinicheaddeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_AssignClinicHeadDeputy

<a id="opIdMembershipCommandService_AssignClinicHeadDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/clinics/{clinicId}/heads/{employeeId}/deputy \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/clinics/{clinicId}/heads/{employeeId}/deputy`

> Body parameter

```json
{
  "deputyEmployeeId": "string"
}
```

<h3 id="membershipcommandservice_assignclinicheaddeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignClinicHeadDeputyBody](#schemamembershipcommandserviceassignclinicheaddeputybody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.

```json
{
  "code": "clinic_not_found",
  "message": "Clinic not found."
}
```

> Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.

```json
{
  "code": "deputy_already_assigned",
  "message": "Deputy already assigned."
}
```

<h3 id="membershipcommandservice_assignclinicheaddeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignClinicHeadDeputyResponse](#schemav1assignclinicheaddeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ------ DepartmentResponsible ------

<a id="opIdMembershipCommandService_AssignDepartmentResponsible"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/departments/{departmentId}/responsibles \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/departments/{departmentId}/responsibles`

> Body parameter

```json
{
  "employeeId": "string"
}
```

<h3 id="-------departmentresponsible--------parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignDepartmentResponsibleBody](#schemamembershipcommandserviceassigndepartmentresponsiblebody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "department_not_found",
  "message": "Department not found."
}
```

> Conflict. Error codes:
- `department_responsible_already_assigned` — this employee is already the department responsible.

```json
{
  "code": "department_responsible_already_assigned",
  "message": "Department responsible already assigned."
}
```

<h3 id="-------departmentresponsible--------responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignDepartmentResponsibleResponse](#schemav1assigndepartmentresponsibleresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `department_responsible_already_assigned` — this employee is already the department responsible.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RevokeDepartmentResponsible

<a id="opIdMembershipCommandService_RevokeDepartmentResponsible"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/departments/{departmentId}/responsibles/{employeeId} \
  -H 'Accept: application/json'

```

`DELETE /v1/departments/{departmentId}/responsibles/{employeeId}`

<h3 id="membershipcommandservice_revokedepartmentresponsible-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `department_responsible_not_found` — this employee is not the department responsible.
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "department_responsible_not_found",
  "message": "Department responsible not found."
}
```

<h3 id="membershipcommandservice_revokedepartmentresponsible-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RevokeDepartmentResponsibleResponse](#schemav1revokedepartmentresponsibleresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `department_responsible_not_found` — this employee is not the department responsible.
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RemoveDepartmentResponsibleDeputy

<a id="opIdMembershipCommandService_RemoveDepartmentResponsibleDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/departments/{departmentId}/responsibles/{employeeId}/deputy \
  -H 'Accept: application/json'

```

`DELETE /v1/departments/{departmentId}/responsibles/{employeeId}/deputy`

<h3 id="membershipcommandservice_removedepartmentresponsibledeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.

```json
{
  "code": "department_not_found",
  "message": "Department not found."
}
```

<h3 id="membershipcommandservice_removedepartmentresponsibledeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RemoveDepartmentResponsibleDeputyResponse](#schemav1removedepartmentresponsibledeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_AssignDepartmentResponsibleDeputy

<a id="opIdMembershipCommandService_AssignDepartmentResponsibleDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/departments/{departmentId}/responsibles/{employeeId}/deputy \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/departments/{departmentId}/responsibles/{employeeId}/deputy`

> Body parameter

```json
{
  "deputyEmployeeId": "string"
}
```

<h3 id="membershipcommandservice_assigndepartmentresponsibledeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignDepartmentResponsibleDeputyBody](#schemamembershipcommandserviceassigndepartmentresponsibledeputybody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.

```json
{
  "code": "department_not_found",
  "message": "Department not found."
}
```

> Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.

```json
{
  "code": "deputy_already_assigned",
  "message": "Deputy already assigned."
}
```

<h3 id="membershipcommandservice_assigndepartmentresponsibledeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignDepartmentResponsibleDeputyResponse](#schemav1assigndepartmentresponsibledeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Employee lifecycle

<a id="opIdMembershipCommandService_HireEmployee"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/employees \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/employees`

> Body parameter

```json
{
  "zitadelUserId": "string",
  "departmentId": "string",
  "position": "string"
}
```

<h3 id="employee-lifecycle-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[v1HireEmployeeRequest](#schemav1hireemployeerequest)|true|none|

> Example responses

> 200 Response

```json
{
  "employeeId": "string"
}
```

> Not found. Error codes:
- `zitadel_user_not_found` — Zitadel user with the given ID does not exist.
- `department_not_found` — department with the given ID does not exist.

```json
{
  "code": "department_not_found",
  "message": "Department not found."
}
```

> Conflict. Error codes:
- `employee_already_hired` — this Zitadel user is already an active employee of this organization.

```json
{
  "code": "employee_already_hired",
  "message": "Employee already hired."
}
```

> Service unavailable. Error codes:
- `zitadel_verify_failed` — Zitadel identity service is temporarily unavailable.

```json
{
  "code": "zitadel_verify_failed",
  "message": "Identity service unavailable."
}
```

<h3 id="employee-lifecycle-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1HireEmployeeResponse](#schemav1hireemployeeresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `zitadel_user_not_found` — Zitadel user with the given ID does not exist.
- `department_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `employee_already_hired` — this Zitadel user is already an active employee of this organization.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|503|[Service Unavailable](https://tools.ietf.org/html/rfc7231#section-6.6.4)|Service unavailable. Error codes:
- `zitadel_verify_failed` — Zitadel identity service is temporarily unavailable.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_TerminateEmployee

<a id="opIdMembershipCommandService_TerminateEmployee"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/employees/{employeeId} \
  -H 'Accept: application/json'

```

`DELETE /v1/employees/{employeeId}`

<h3 id="membershipcommandservice_terminateemployee-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "employee_not_found",
  "message": "Employee not found."
}
```

<h3 id="membershipcommandservice_terminateemployee-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1TerminateEmployeeResponse](#schemav1terminateemployeeresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_UpdateEmployeeDepartment

<a id="opIdMembershipCommandService_UpdateEmployeeDepartment"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/employees/{employeeId}/department \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/employees/{employeeId}/department`

> Body parameter

```json
{
  "departmentId": "string"
}
```

<h3 id="membershipcommandservice_updateemployeedepartment-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceUpdateEmployeeDepartmentBody](#schemamembershipcommandserviceupdateemployeedepartmentbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.
- `department_not_found` — department with the given ID does not exist.

```json
{
  "code": "employee_not_found",
  "message": "Employee not found."
}
```

<h3 id="membershipcommandservice_updateemployeedepartment-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateEmployeeDepartmentResponse](#schemav1updateemployeedepartmentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.
- `department_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_UpdateEmployeePosition

<a id="opIdMembershipCommandService_UpdateEmployeePosition"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/employees/{employeeId}/position \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/employees/{employeeId}/position`

> Body parameter

```json
{
  "position": "string"
}
```

<h3 id="membershipcommandservice_updateemployeeposition-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceUpdateEmployeePositionBody](#schemamembershipcommandserviceupdateemployeepositionbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "employee_not_found",
  "message": "Employee not found."
}
```

<h3 id="membershipcommandservice_updateemployeeposition-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateEmployeePositionResponse](#schemav1updateemployeepositionresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_ScheduleVacation

<a id="opIdMembershipCommandService_ScheduleVacation"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/employees/{employeeId}/vacations \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/employees/{employeeId}/vacations`

> Body parameter

```json
{
  "startsAt": "2019-08-24T14:15:22Z",
  "endsAt": "2019-08-24T14:15:22Z"
}
```

<h3 id="membershipcommandservice_schedulevacation-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceScheduleVacationBody](#schemamembershipcommandserviceschedulevacationbody)|true|none|

> Example responses

> 200 Response

```json
{
  "vacationId": "string"
}
```

> Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "employee_not_found",
  "message": "Employee not found."
}
```

<h3 id="membershipcommandservice_schedulevacation-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ScheduleVacationResponse](#schemav1schedulevacationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Vacation lifecycle

<a id="opIdMembershipCommandService_StartVacationNow"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/employees/{employeeId}/vacations:start-now \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/employees/{employeeId}/vacations:start-now`

> Body parameter

```json
{
  "endsAt": "2019-08-24T14:15:22Z"
}
```

<h3 id="vacation-lifecycle-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceStartVacationNowBody](#schemamembershipcommandservicestartvacationnowbody)|true|none|

> Example responses

> 200 Response

```json
{
  "vacationId": "string"
}
```

> Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "employee_not_found",
  "message": "Employee not found."
}
```

<h3 id="vacation-lifecycle-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1StartVacationNowResponse](#schemav1startvacationnowresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ------ OrganizationAdmin ------

<a id="opIdMembershipCommandService_AssignOrganizationAdmin"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/admins \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/admins`

> Body parameter

```json
{
  "employeeId": "string"
}
```

<h3 id="-------organizationadmin--------parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignOrganizationAdminBody](#schemamembershipcommandserviceassignorganizationadminbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

> Conflict. Error codes:
- `organization_admin_already_assigned` — this employee is already an organization admin.

```json
{
  "code": "organization_admin_already_assigned",
  "message": "Organization admin already assigned."
}
```

<h3 id="-------organizationadmin--------responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignOrganizationAdminResponse](#schemav1assignorganizationadminresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `organization_admin_already_assigned` — this employee is already an organization admin.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RevokeOrganizationAdmin

<a id="opIdMembershipCommandService_RevokeOrganizationAdmin"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/organizations/{organizationId}/admins/{employeeId} \
  -H 'Accept: application/json'

```

`DELETE /v1/organizations/{organizationId}/admins/{employeeId}`

<h3 id="membershipcommandservice_revokeorganizationadmin-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_admin_not_found` — this employee is not an organization admin.

```json
{
  "code": "organization_admin_not_found",
  "message": "Organization admin not found."
}
```

<h3 id="membershipcommandservice_revokeorganizationadmin-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RevokeOrganizationAdminResponse](#schemav1revokeorganizationadminresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_admin_not_found` — this employee is not an organization admin.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RemoveOrganizationAdminDeputy

<a id="opIdMembershipCommandService_RemoveOrganizationAdminDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/organizations/{organizationId}/admins/{employeeId}/deputy \
  -H 'Accept: application/json'

```

`DELETE /v1/organizations/{organizationId}/admins/{employeeId}/deputy`

<h3 id="membershipcommandservice_removeorganizationadmindeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="membershipcommandservice_removeorganizationadmindeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RemoveOrganizationAdminDeputyResponse](#schemav1removeorganizationadmindeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_AssignOrganizationAdminDeputy

<a id="opIdMembershipCommandService_AssignOrganizationAdminDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/admins/{employeeId}/deputy \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/admins/{employeeId}/deputy`

> Body parameter

```json
{
  "deputyEmployeeId": "string"
}
```

<h3 id="membershipcommandservice_assignorganizationadmindeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignOrganizationAdminDeputyBody](#schemamembershipcommandserviceassignorganizationadmindeputybody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

> Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.

```json
{
  "code": "deputy_already_assigned",
  "message": "Deputy already assigned."
}
```

<h3 id="membershipcommandservice_assignorganizationadmindeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignOrganizationAdminDeputyResponse](#schemav1assignorganizationadmindeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ------ OrganizationDispatcher ------

<a id="opIdMembershipCommandService_AssignOrganizationDispatcher"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/dispatchers \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/dispatchers`

> Body parameter

```json
{
  "employeeId": "string"
}
```

<h3 id="-------organizationdispatcher--------parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignOrganizationDispatcherBody](#schemamembershipcommandserviceassignorganizationdispatcherbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

> Conflict. Error codes:
- `organization_dispatcher_already_assigned` — this employee is already an organization dispatcher.

```json
{
  "code": "organization_dispatcher_already_assigned",
  "message": "Organization dispatcher already assigned."
}
```

<h3 id="-------organizationdispatcher--------responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignOrganizationDispatcherResponse](#schemav1assignorganizationdispatcherresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `organization_dispatcher_already_assigned` — this employee is already an organization dispatcher.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RevokeOrganizationDispatcher

<a id="opIdMembershipCommandService_RevokeOrganizationDispatcher"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/organizations/{organizationId}/dispatchers/{employeeId} \
  -H 'Accept: application/json'

```

`DELETE /v1/organizations/{organizationId}/dispatchers/{employeeId}`

<h3 id="membershipcommandservice_revokeorganizationdispatcher-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_dispatcher_not_found` — this employee is not an organization dispatcher.

```json
{
  "code": "organization_dispatcher_not_found",
  "message": "Organization dispatcher not found."
}
```

<h3 id="membershipcommandservice_revokeorganizationdispatcher-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RevokeOrganizationDispatcherResponse](#schemav1revokeorganizationdispatcherresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_dispatcher_not_found` — this employee is not an organization dispatcher.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RemoveOrganizationDispatcherDeputy

<a id="opIdMembershipCommandService_RemoveOrganizationDispatcherDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/organizations/{organizationId}/dispatchers/{employeeId}/deputy \
  -H 'Accept: application/json'

```

`DELETE /v1/organizations/{organizationId}/dispatchers/{employeeId}/deputy`

<h3 id="membershipcommandservice_removeorganizationdispatcherdeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="membershipcommandservice_removeorganizationdispatcherdeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RemoveOrganizationDispatcherDeputyResponse](#schemav1removeorganizationdispatcherdeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_AssignOrganizationDispatcherDeputy

<a id="opIdMembershipCommandService_AssignOrganizationDispatcherDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/dispatchers/{employeeId}/deputy \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/dispatchers/{employeeId}/deputy`

> Body parameter

```json
{
  "deputyEmployeeId": "string"
}
```

<h3 id="membershipcommandservice_assignorganizationdispatcherdeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignOrganizationDispatcherDeputyBody](#schemamembershipcommandserviceassignorganizationdispatcherdeputybody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

> Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.

```json
{
  "code": "deputy_already_assigned",
  "message": "Deputy already assigned."
}
```

<h3 id="membershipcommandservice_assignorganizationdispatcherdeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignOrganizationDispatcherDeputyResponse](#schemav1assignorganizationdispatcherdeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ------ OrganizationHead ------

<a id="opIdMembershipCommandService_AssignOrganizationHead"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/heads \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/heads`

> Body parameter

```json
{
  "employeeId": "string"
}
```

<h3 id="-------organizationhead--------parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignOrganizationHeadBody](#schemamembershipcommandserviceassignorganizationheadbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

> Conflict. Error codes:
- `organization_head_already_assigned` — this organization already has a head assigned.

```json
{
  "code": "organization_head_already_assigned",
  "message": "Organization head already assigned."
}
```

<h3 id="-------organizationhead--------responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignOrganizationHeadResponse](#schemav1assignorganizationheadresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `organization_head_already_assigned` — this organization already has a head assigned.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RevokeOrganizationHead

<a id="opIdMembershipCommandService_RevokeOrganizationHead"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/organizations/{organizationId}/heads/{employeeId} \
  -H 'Accept: application/json'

```

`DELETE /v1/organizations/{organizationId}/heads/{employeeId}`

<h3 id="membershipcommandservice_revokeorganizationhead-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_head_not_found` — this employee is not the organization head.

```json
{
  "code": "organization_head_not_found",
  "message": "Organization head not found."
}
```

<h3 id="membershipcommandservice_revokeorganizationhead-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RevokeOrganizationHeadResponse](#schemav1revokeorganizationheadresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_head_not_found` — this employee is not the organization head.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RemoveOrganizationHeadDeputy

<a id="opIdMembershipCommandService_RemoveOrganizationHeadDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/organizations/{organizationId}/heads/{employeeId}/deputy \
  -H 'Accept: application/json'

```

`DELETE /v1/organizations/{organizationId}/heads/{employeeId}/deputy`

<h3 id="membershipcommandservice_removeorganizationheaddeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="membershipcommandservice_removeorganizationheaddeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RemoveOrganizationHeadDeputyResponse](#schemav1removeorganizationheaddeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_AssignOrganizationHeadDeputy

<a id="opIdMembershipCommandService_AssignOrganizationHeadDeputy"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/heads/{employeeId}/deputy \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/heads/{employeeId}/deputy`

> Body parameter

```json
{
  "deputyEmployeeId": "string"
}
```

<h3 id="membershipcommandservice_assignorganizationheaddeputy-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|employeeId|path|string|true|none|
|body|body|[MembershipCommandServiceAssignOrganizationHeadDeputyBody](#schemamembershipcommandserviceassignorganizationheaddeputybody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

> Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.

```json
{
  "code": "deputy_already_assigned",
  "message": "Deputy already assigned."
}
```

<h3 id="membershipcommandservice_assignorganizationheaddeputy-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignOrganizationHeadDeputyResponse](#schemav1assignorganizationheaddeputyresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.
- `employee_not_found` — employee with the given ID does not exist.
- `deputy_not_found` — deputy employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `deputy_already_assigned` — this employee is already a deputy for this role.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ------ SystemAdmin ------

<a id="opIdMembershipCommandService_GrantSystemAdmin"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/system-admins \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/system-admins`

> Body parameter

```json
{
  "zitadelUserId": "string"
}
```

<h3 id="-------systemadmin--------parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[v1GrantSystemAdminRequest](#schemav1grantsystemadminrequest)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `zitadel_user_not_found` — Zitadel user with the given ID does not exist.

```json
{
  "code": "zitadel_user_not_found",
  "message": "Zitadel user not found."
}
```

> Conflict. Error codes:
- `system_admin_already_granted` — this user is already a system admin.

```json
{
  "code": "system_admin_already_granted",
  "message": "System admin already granted."
}
```

> Service unavailable. Error codes:
- `zitadel_verify_failed` — Zitadel identity service is temporarily unavailable.

```json
{
  "code": "zitadel_verify_failed",
  "message": "Identity service unavailable."
}
```

<h3 id="-------systemadmin--------responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GrantSystemAdminResponse](#schemav1grantsystemadminresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `zitadel_user_not_found` — Zitadel user with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `system_admin_already_granted` — this user is already a system admin.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|503|[Service Unavailable](https://tools.ietf.org/html/rfc7231#section-6.6.4)|Service unavailable. Error codes:
- `zitadel_verify_failed` — Zitadel identity service is temporarily unavailable.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_RevokeSystemAdmin

<a id="opIdMembershipCommandService_RevokeSystemAdmin"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/system-admins/{zitadelUserId} \
  -H 'Accept: application/json'

```

`DELETE /v1/system-admins/{zitadelUserId}`

<h3 id="membershipcommandservice_revokesystemadmin-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|zitadelUserId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `system_admin_not_found` — this user is not a system admin.

```json
{
  "code": "system_admin_not_found",
  "message": "System admin not found."
}
```

<h3 id="membershipcommandservice_revokesystemadmin-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1RevokeSystemAdminResponse](#schemav1revokesystemadminresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `system_admin_not_found` — this user is not a system admin.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_CancelScheduledVacation

<a id="opIdMembershipCommandService_CancelScheduledVacation"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/vacations/{vacationId}/cancellations \
  -H 'Accept: application/json'

```

`POST /v1/vacations/{vacationId}/cancellations`

<h3 id="membershipcommandservice_cancelscheduledvacation-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|vacationId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `vacation_not_found` — vacation with the given ID does not exist.

```json
{
  "code": "vacation_not_found",
  "message": "Vacation not found."
}
```

<h3 id="membershipcommandservice_cancelscheduledvacation-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CancelScheduledVacationResponse](#schemav1cancelscheduledvacationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `vacation_not_found` — vacation with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_UpdateVacationEndDate

<a id="opIdMembershipCommandService_UpdateVacationEndDate"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/vacations/{vacationId}/end-date \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/vacations/{vacationId}/end-date`

> Body parameter

```json
{
  "endsAt": "2019-08-24T14:15:22Z"
}
```

<h3 id="membershipcommandservice_updatevacationenddate-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|vacationId|path|string|true|none|
|body|body|[MembershipCommandServiceUpdateVacationEndDateBody](#schemamembershipcommandserviceupdatevacationenddatebody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `vacation_not_found` — vacation with the given ID does not exist.

```json
{
  "code": "vacation_not_found",
  "message": "Vacation not found."
}
```

<h3 id="membershipcommandservice_updatevacationenddate-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateVacationEndDateResponse](#schemav1updatevacationenddateresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `vacation_not_found` — vacation with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipCommandService_ForceEndVacation

<a id="opIdMembershipCommandService_ForceEndVacation"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/vacations/{vacationId}/terminations \
  -H 'Accept: application/json'

```

`POST /v1/vacations/{vacationId}/terminations`

<h3 id="membershipcommandservice_forceendvacation-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|vacationId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `vacation_not_found` — vacation with the given ID does not exist.

```json
{
  "code": "vacation_not_found",
  "message": "Vacation not found."
}
```

<h3 id="membershipcommandservice_forceendvacation-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ForceEndVacationResponse](#schemav1forceendvacationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `vacation_not_found` — vacation with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-orgstructurecommandservice">OrgStructureCommandService</h1>

## OrgStructureCommandService_CreateDepartment

<a id="opIdOrgStructureCommandService_CreateDepartment"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/clinics/{clinicId}/departments \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/clinics/{clinicId}/departments`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="orgstructurecommandservice_createdepartment-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|body|body|[OrgStructureCommandServiceCreateDepartmentBody](#schemaorgstructurecommandservicecreatedepartmentbody)|true|none|

> Example responses

> 200 Response

```json
{
  "departmentId": "string"
}
```

> Not found. Error codes:
- `department_clinic_not_found` — clinic with the given ID does not exist.

```json
{
  "code": "department_clinic_not_found",
  "message": "Clinic not found."
}
```

<h3 id="orgstructurecommandservice_createdepartment-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateDepartmentResponse](#schemav1createdepartmentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `department_clinic_not_found` — clinic with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureCommandService_UpdateClinicDetails

<a id="opIdOrgStructureCommandService_UpdateClinicDetails"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/clinics/{clinicId}/details \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/clinics/{clinicId}/details`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="orgstructurecommandservice_updateclinicdetails-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|body|body|[OrgStructureCommandServiceUpdateClinicDetailsBody](#schemaorgstructurecommandserviceupdateclinicdetailsbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.

```json
{
  "code": "clinic_not_found",
  "message": "Clinic not found."
}
```

<h3 id="orgstructurecommandservice_updateclinicdetails-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateClinicDetailsResponse](#schemav1updateclinicdetailsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureCommandService_UpdateClinicPhysicalAddress

<a id="opIdOrgStructureCommandService_UpdateClinicPhysicalAddress"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/clinics/{clinicId}/physical-address \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/clinics/{clinicId}/physical-address`

> Body parameter

```json
{
  "physicalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  }
}
```

<h3 id="orgstructurecommandservice_updateclinicphysicaladdress-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|body|body|[OrgStructureCommandServiceUpdateClinicPhysicalAddressBody](#schemaorgstructurecommandserviceupdateclinicphysicaladdressbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.

```json
{
  "code": "clinic_not_found",
  "message": "Clinic not found."
}
```

<h3 id="orgstructurecommandservice_updateclinicphysicaladdress-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateClinicPhysicalAddressResponse](#schemav1updateclinicphysicaladdressresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureCommandService_UpdateDepartmentDetails

<a id="opIdOrgStructureCommandService_UpdateDepartmentDetails"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/departments/{departmentId}/details \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/departments/{departmentId}/details`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="orgstructurecommandservice_updatedepartmentdetails-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|
|body|body|[OrgStructureCommandServiceUpdateDepartmentDetailsBody](#schemaorgstructurecommandserviceupdatedepartmentdetailsbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.

```json
{
  "code": "department_not_found",
  "message": "Department not found."
}
```

<h3 id="orgstructurecommandservice_updatedepartmentdetails-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateDepartmentDetailsResponse](#schemav1updatedepartmentdetailsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureCommandService_CreateOrganization

<a id="opIdOrgStructureCommandService_CreateOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations`

> Body parameter

```json
{
  "name": "string",
  "legalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  },
  "description": "string"
}
```

<h3 id="orgstructurecommandservice_createorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[v1CreateOrganizationRequest](#schemav1createorganizationrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "organizationId": "string"
}
```

<h3 id="orgstructurecommandservice_createorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateOrganizationResponse](#schemav1createorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureCommandService_CreateClinic

<a id="opIdOrgStructureCommandService_CreateClinic"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/clinics \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/clinics`

> Body parameter

```json
{
  "name": "string",
  "physicalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  },
  "description": "string"
}
```

<h3 id="orgstructurecommandservice_createclinic-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|body|body|[OrgStructureCommandServiceCreateClinicBody](#schemaorgstructurecommandservicecreateclinicbody)|true|none|

> Example responses

> 200 Response

```json
{
  "clinicId": "string"
}
```

> Not found. Error codes:
- `clinic_organization_not_found` — organization with the given ID does not exist.

```json
{
  "code": "clinic_organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="orgstructurecommandservice_createclinic-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateClinicResponse](#schemav1createclinicresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `clinic_organization_not_found` — organization with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureCommandService_UpdateOrganizationDetails

<a id="opIdOrgStructureCommandService_UpdateOrganizationDetails"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/organizations/{organizationId}/details \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/organizations/{organizationId}/details`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="orgstructurecommandservice_updateorganizationdetails-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|body|body|[OrgStructureCommandServiceUpdateOrganizationDetailsBody](#schemaorgstructurecommandserviceupdateorganizationdetailsbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="orgstructurecommandservice_updateorganizationdetails-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateOrganizationDetailsResponse](#schemav1updateorganizationdetailsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureCommandService_UpdateOrganizationLegalAddress

<a id="opIdOrgStructureCommandService_UpdateOrganizationLegalAddress"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/organizations/{organizationId}/legal-address \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/organizations/{organizationId}/legal-address`

> Body parameter

```json
{
  "legalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  }
}
```

<h3 id="orgstructurecommandservice_updateorganizationlegaladdress-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|body|body|[OrgStructureCommandServiceUpdateOrganizationLegalAddressBody](#schemaorgstructurecommandserviceupdateorganizationlegaladdressbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="orgstructurecommandservice_updateorganizationlegaladdress-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateOrganizationLegalAddressResponse](#schemav1updateorganizationlegaladdressresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-requestclassifiercommandservice">RequestClassifierCommandService</h1>

## RequestClassifierCommandService_CreateRequestType

<a id="opIdRequestClassifierCommandService_CreateRequestType"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/organizations/{organizationId}/request-types \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/organizations/{organizationId}/request-types`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="requestclassifiercommandservice_createrequesttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|body|body|[RequestClassifierCommandServiceCreateRequestTypeBody](#schemarequestclassifiercommandservicecreaterequesttypebody)|true|none|

> Example responses

> 200 Response

```json
{
  "typeId": "string"
}
```

> Conflict. Error codes:
- `request_type_name_conflict` — a request type with this name already exists in the organization.

```json
{
  "code": "request_type_name_conflict",
  "message": "Request type name already exists."
}
```

<h3 id="requestclassifiercommandservice_createrequesttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateRequestTypeResponse](#schemav1createrequesttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `request_type_name_conflict` — a request type with this name already exists in the organization.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## RequestClassifierCommandService_DeleteRequestType

<a id="opIdRequestClassifierCommandService_DeleteRequestType"></a>

> Code samples

```shell
# You can also use wget
curl -X DELETE /v1/request-types/{typeId} \
  -H 'Accept: application/json'

```

`DELETE /v1/request-types/{typeId}`

<h3 id="requestclassifiercommandservice_deleterequesttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.

```json
{
  "code": "request_type_not_found",
  "message": "Request type not found."
}
```

<h3 id="requestclassifiercommandservice_deleterequesttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1DeleteRequestTypeResponse](#schemav1deleterequesttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## RequestClassifierCommandService_DeactivateRequestType

<a id="opIdRequestClassifierCommandService_DeactivateRequestType"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/request-types/{typeId}/deactivations \
  -H 'Accept: application/json'

```

`POST /v1/request-types/{typeId}/deactivations`

<h3 id="requestclassifiercommandservice_deactivaterequesttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.

```json
{
  "code": "request_type_not_found",
  "message": "Request type not found."
}
```

<h3 id="requestclassifiercommandservice_deactivaterequesttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1DeactivateRequestTypeResponse](#schemav1deactivaterequesttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## RequestClassifierCommandService_UpdateRequestTypeDetails

<a id="opIdRequestClassifierCommandService_UpdateRequestTypeDetails"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/request-types/{typeId}/details \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/request-types/{typeId}/details`

> Body parameter

```json
{
  "name": "string",
  "description": "string"
}
```

<h3 id="requestclassifiercommandservice_updaterequesttypedetails-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|
|body|body|[RequestClassifierCommandServiceUpdateRequestTypeDetailsBody](#schemarequestclassifiercommandserviceupdaterequesttypedetailsbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.

```json
{
  "code": "request_type_not_found",
  "message": "Request type not found."
}
```

> Conflict. Error codes:
- `request_type_name_conflict` — a request type with this name already exists in the organization.

```json
{
  "code": "request_type_name_conflict",
  "message": "Request type name already exists."
}
```

<h3 id="requestclassifiercommandservice_updaterequesttypedetails-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateRequestTypeDetailsResponse](#schemav1updaterequesttypedetailsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `request_type_name_conflict` — a request type with this name already exists in the organization.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## RequestClassifierCommandService_ReactivateRequestType

<a id="opIdRequestClassifierCommandService_ReactivateRequestType"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/request-types/{typeId}/reactivations \
  -H 'Accept: application/json'

```

`POST /v1/request-types/{typeId}/reactivations`

<h3 id="requestclassifiercommandservice_reactivaterequesttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|typeId|path|string|true|none|

> Example responses

> 200 Response

```json
{}
```

> Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.

```json
{
  "code": "request_type_not_found",
  "message": "Request type not found."
}
```

> Conflict. Error codes:
- `request_type_name_conflict` — reactivation would create a name conflict with an active type.

```json
{
  "code": "request_type_name_conflict",
  "message": "Request type name already exists."
}
```

<h3 id="requestclassifiercommandservice_reactivaterequesttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ReactivateRequestTypeResponse](#schemav1reactivaterequesttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|409|[Conflict](https://tools.ietf.org/html/rfc7231#section-6.5.8)|Conflict. Error codes:
- `request_type_name_conflict` — reactivation would create a name conflict with an active type.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-servicerequestcommandservice">ServiceRequestCommandService</h1>

## ServiceRequestCommandService_CreateServiceRequest

<a id="opIdServiceRequestCommandService_CreateServiceRequest"></a>

> Code samples

```shell
# You can also use wget
curl -X POST /v1/service-requests \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`POST /v1/service-requests`

> Body parameter

```json
{
  "departmentId": "string",
  "typeId": "string",
  "incidentId": "string",
  "description": "string",
  "executorEmployeeIds": [
    "string"
  ]
}
```

<h3 id="servicerequestcommandservice_createservicerequest-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[v1CreateServiceRequestRequest](#schemav1createservicerequestrequest)|true|none|

> Example responses

> 200 Response

```json
{
  "serviceRequestId": "string"
}
```

> Failed precondition. Error codes:
- `service_request_type_inactive` — the request type is inactive.
- `service_request_type_org_mismatch` — the request type belongs to a different organization.
- `service_request_incident_org_mismatch` — the linked incident belongs to a different organization.
- `service_request_employee_dept_mismatch` — an executor employee does not belong to the specified department.

```json
{
  "code": "service_request_type_inactive",
  "message": "Request type is inactive."
}
```

> Not found. Error codes:
- `service_request_department_not_found` — department with the given ID does not exist.
- `service_request_type_not_found` — request type with the given ID does not exist.
- `service_request_incident_not_found` — linked incident with the given ID does not exist.
- `service_request_employee_not_found` — one of the executor employees was not found.

```json
{
  "code": "service_request_department_not_found",
  "message": "Department not found."
}
```

<h3 id="servicerequestcommandservice_createservicerequest-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CreateServiceRequestResponse](#schemav1createservicerequestresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Failed precondition. Error codes:
- `service_request_type_inactive` — the request type is inactive.
- `service_request_type_org_mismatch` — the request type belongs to a different organization.
- `service_request_incident_org_mismatch` — the linked incident belongs to a different organization.
- `service_request_employee_dept_mismatch` — an executor employee does not belong to the specified department.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `service_request_department_not_found` — department with the given ID does not exist.
- `service_request_type_not_found` — request type with the given ID does not exist.
- `service_request_incident_not_found` — linked incident with the given ID does not exist.
- `service_request_employee_not_found` — one of the executor employees was not found.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ServiceRequestCommandService_UpdateServiceRequestDescription

<a id="opIdServiceRequestCommandService_UpdateServiceRequestDescription"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/service-requests/{serviceRequestId}/description \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/service-requests/{serviceRequestId}/description`

> Body parameter

```json
{
  "description": "string"
}
```

<h3 id="servicerequestcommandservice_updateservicerequestdescription-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|serviceRequestId|path|string|true|none|
|body|body|[ServiceRequestCommandServiceUpdateServiceRequestDescriptionBody](#schemaservicerequestcommandserviceupdateservicerequestdescriptionbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Failed precondition. Error codes:
- `service_request_frozen` — service request is frozen and cannot be modified.

```json
{
  "code": "service_request_frozen",
  "message": "Service request is frozen."
}
```

> Not found. Error codes:
- `service_request_not_found` — service request with the given ID does not exist.

```json
{
  "code": "service_request_not_found",
  "message": "Service request not found."
}
```

<h3 id="servicerequestcommandservice_updateservicerequestdescription-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateServiceRequestDescriptionResponse](#schemav1updateservicerequestdescriptionresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Failed precondition. Error codes:
- `service_request_frozen` — service request is frozen and cannot be modified.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `service_request_not_found` — service request with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ServiceRequestCommandService_AssignExecutors

<a id="opIdServiceRequestCommandService_AssignExecutors"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/service-requests/{serviceRequestId}/executors \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/service-requests/{serviceRequestId}/executors`

> Body parameter

```json
{
  "executorEmployeeIds": [
    "string"
  ]
}
```

<h3 id="servicerequestcommandservice_assignexecutors-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|serviceRequestId|path|string|true|none|
|body|body|[ServiceRequestCommandServiceAssignExecutorsBody](#schemaservicerequestcommandserviceassignexecutorsbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Failed precondition. Error codes:
- `service_request_frozen` — service request is frozen and cannot be modified.
- `service_request_employee_dept_mismatch` — an executor employee does not belong to the service request's department.

```json
{
  "code": "service_request_frozen",
  "message": "Service request is frozen."
}
```

> Not found. Error codes:
- `service_request_not_found` — service request with the given ID does not exist.
- `service_request_employee_not_found` — one of the executor employees was not found.

```json
{
  "code": "service_request_not_found",
  "message": "Service request not found."
}
```

<h3 id="servicerequestcommandservice_assignexecutors-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1AssignExecutorsResponse](#schemav1assignexecutorsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Failed precondition. Error codes:
- `service_request_frozen` — service request is frozen and cannot be modified.
- `service_request_employee_dept_mismatch` — an executor employee does not belong to the service request's department.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `service_request_not_found` — service request with the given ID does not exist.
- `service_request_employee_not_found` — one of the executor employees was not found.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ServiceRequestCommandService_UpdateServiceRequestStatus

<a id="opIdServiceRequestCommandService_UpdateServiceRequestStatus"></a>

> Code samples

```shell
# You can also use wget
curl -X PUT /v1/service-requests/{serviceRequestId}/status \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'

```

`PUT /v1/service-requests/{serviceRequestId}/status`

> Body parameter

```json
{
  "newStatus": "string"
}
```

<h3 id="servicerequestcommandservice_updateservicerequeststatus-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|serviceRequestId|path|string|true|none|
|body|body|[ServiceRequestCommandServiceUpdateServiceRequestStatusBody](#schemaservicerequestcommandserviceupdateservicerequeststatusbody)|true|none|

> Example responses

> 200 Response

```json
{}
```

> Failed precondition. Error codes:
- `service_request_invalid_status_transition` — the requested status transition is not allowed.
- `service_request_frozen` — service request is frozen and cannot be modified.

```json
{
  "code": "service_request_invalid_status_transition",
  "message": "Invalid status transition."
}
```

> Not found. Error codes:
- `service_request_not_found` — service request with the given ID does not exist.

```json
{
  "code": "service_request_not_found",
  "message": "Service request not found."
}
```

<h3 id="servicerequestcommandservice_updateservicerequeststatus-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1UpdateServiceRequestStatusResponse](#schemav1updateservicerequeststatusresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Failed precondition. Error codes:
- `service_request_invalid_status_transition` — the requested status transition is not allowed.
- `service_request_frozen` — service request is frozen and cannot be modified.|[v1ErrorResponse](#schemav1errorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `service_request_not_found` — service request with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-analyticsqueryservice">AnalyticsQueryService</h1>

## AnalyticsQueryService_GetSnapshot

<a id="opIdAnalyticsQueryService_GetSnapshot"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/analytics/snapshot \
  -H 'Accept: application/json'

```

`GET /v1/analytics/snapshot`

<h3 id="analyticsqueryservice_getsnapshot-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|query|string|false|none|
|from|query|string|false|none|
|to|query|string|false|none|
|clinicId|query|string|false|none|
|departmentId|query|string|false|none|
|includePatientBuffer|query|boolean|false|none|

> Example responses

> 200 Response

```json
{
  "incidents": [
    {
      "createdAt": "string",
      "occurredAt": "string",
      "closedAt": "string",
      "status": "string",
      "priority": "string",
      "categoryId": "string",
      "categoryName": "string",
      "typeId": "string",
      "typeName": "string",
      "clinicId": "string",
      "departmentId": "string",
      "isPatientSource": true,
      "isReopened": true,
      "linkedRequestsCount": 0
    }
  ],
  "requests": [
    {
      "createdAt": "string",
      "completedAt": "string",
      "status": "string",
      "typeId": "string",
      "typeName": "string",
      "departmentId": "string",
      "hasLinkedIncident": true
    }
  ],
  "patientBuffer": [
    {
      "createdAt": "string",
      "status": "string",
      "categoryId": "string"
    }
  ]
}
```

> Not found. Error codes:
- `analytics_org_not_found` — organization with the given ID does not exist.
- `analytics_clinic_not_found` — clinic with the given ID does not exist.
- `analytics_dept_not_found` — department with the given ID does not exist.

```json
{
  "code": "analytics_org_not_found",
  "message": "Organization not found."
}
```

<h3 id="analyticsqueryservice_getsnapshot-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetSnapshotResponse](#schemav1getsnapshotresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `analytics_org_not_found` — organization with the given ID does not exist.
- `analytics_clinic_not_found` — clinic with the given ID does not exist.
- `analytics_dept_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnalyticsQueryService_GetSummary

<a id="opIdAnalyticsQueryService_GetSummary"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/analytics/summary \
  -H 'Accept: application/json'

```

`GET /v1/analytics/summary`

<h3 id="analyticsqueryservice_getsummary-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|query|string|false|none|
|from|query|string|false|none|
|to|query|string|false|none|
|clinicId|query|string|false|none|
|departmentId|query|string|false|none|

> Example responses

> 200 Response

```json
{
  "incidents": {
    "total": "string",
    "byStatus": {
      "pending": "string",
      "inProgress": "string",
      "done": "string",
      "rejected": "string",
      "cancelled": "string"
    },
    "byPriority": {
      "low": "string",
      "normal": "string",
      "high": "string",
      "critical": "string"
    },
    "bySource": {
      "staff": "string",
      "patient": "string"
    },
    "reopened": "string",
    "withLinkedRequests": "string",
    "resolution": {
      "avgMinutes": 0.1,
      "minMinutes": 0.1,
      "maxMinutes": 0.1,
      "p50Minutes": 0.1,
      "p90Minutes": 0.1,
      "p95Minutes": 0.1
    },
    "topCategories": [
      {
        "categoryId": "string",
        "categoryName": "string",
        "count": "string"
      }
    ],
    "topTypes": [
      {
        "typeId": "string",
        "typeName": "string",
        "count": "string"
      }
    ],
    "topDepartments": [
      {
        "departmentId": "string",
        "departmentName": "string",
        "count": "string"
      }
    ]
  },
  "requests": {
    "total": "string",
    "byStatus": {
      "created": "string",
      "inWork": "string",
      "onHold": "string",
      "pendingReview": "string",
      "completed": "string",
      "cancelled": "string"
    },
    "linked": "string",
    "unlinked": "string",
    "completion": {
      "avgMinutes": 0.1,
      "minMinutes": 0.1,
      "maxMinutes": 0.1,
      "p50Minutes": 0.1,
      "p90Minutes": 0.1,
      "p95Minutes": 0.1
    },
    "topTypes": [
      {
        "typeId": "string",
        "typeName": "string",
        "count": "string"
      }
    ],
    "topDepartments": [
      {
        "departmentId": "string",
        "departmentName": "string",
        "count": "string"
      }
    ]
  },
  "patientBuffer": {
    "total": "string",
    "byStatus": {
      "pending": "string",
      "published": "string",
      "rejected": "string",
      "cancelled": "string"
    },
    "acceptanceRate": 0.1,
    "rejectionRate": 0.1
  },
  "period": {
    "organizationId": "string",
    "from": "string",
    "to": "string",
    "clinicId": "string",
    "departmentId": "string"
  }
}
```

> Not found. Error codes:
- `analytics_org_not_found` — organization with the given ID does not exist.
- `analytics_clinic_not_found` — clinic with the given ID does not exist.
- `analytics_dept_not_found` — department with the given ID does not exist.

```json
{
  "code": "analytics_org_not_found",
  "message": "Organization not found."
}
```

<h3 id="analyticsqueryservice_getsummary-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetSummaryResponse](#schemav1getsummaryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `analytics_org_not_found` — organization with the given ID does not exist.
- `analytics_clinic_not_found` — clinic with the given ID does not exist.
- `analytics_dept_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnalyticsQueryService_GetTimeSeries

<a id="opIdAnalyticsQueryService_GetTimeSeries"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/analytics/timeseries \
  -H 'Accept: application/json'

```

`GET /v1/analytics/timeseries`

<h3 id="analyticsqueryservice_gettimeseries-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|query|string|false|none|
|from|query|string|false|none|
|to|query|string|false|none|
|clinicId|query|string|false|none|
|departmentId|query|string|false|none|
|granularity|query|string|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|granularity|TIME_SERIES_GRANULARITY_UNSPECIFIED|
|granularity|TIME_SERIES_GRANULARITY_DAY|
|granularity|TIME_SERIES_GRANULARITY_WEEK|
|granularity|TIME_SERIES_GRANULARITY_MONTH|

> Example responses

> 200 Response

```json
{
  "buckets": [
    {
      "bucketStart": "string",
      "bucketEnd": "string",
      "incidents": {
        "total": "string",
        "pending": "string",
        "inProgress": "string",
        "done": "string",
        "rejected": "string",
        "cancelled": "string",
        "highCritical": "string",
        "patientSource": "string",
        "reopened": "string"
      },
      "requests": {
        "total": "string",
        "completed": "string",
        "cancelled": "string",
        "linked": "string"
      }
    }
  ]
}
```

> Not found. Error codes:
- `analytics_org_not_found` — organization with the given ID does not exist.
- `analytics_clinic_not_found` — clinic with the given ID does not exist.
- `analytics_dept_not_found` — department with the given ID does not exist.

```json
{
  "code": "analytics_org_not_found",
  "message": "Organization not found."
}
```

<h3 id="analyticsqueryservice_gettimeseries-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetTimeSeriesResponse](#schemav1gettimeseriesresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `analytics_org_not_found` — organization with the given ID does not exist.
- `analytics_clinic_not_found` — clinic with the given ID does not exist.
- `analytics_dept_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-announcementqueryservice">AnnouncementQueryService</h1>

## AnnouncementQueryService_GetAnnouncement

<a id="opIdAnnouncementQueryService_GetAnnouncement"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/announcements/{id} \
  -H 'Accept: application/json'

```

`GET /v1/query/announcements/{id}`

<h3 id="announcementqueryservice_getannouncement-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "announcement": {
    "id": "string",
    "organizationId": "string",
    "clinicId": "string",
    "departmentId": "string",
    "authorId": "string",
    "title": "string",
    "content": "string",
    "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
    "isArchived": true,
    "startsAt": "string",
    "endsAt": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "viewCount": "string"
  }
}
```

> Not found. Error codes:
- `announcement_query_not_found` — announcement with the given ID does not exist.

```json
{
  "code": "announcement_query_not_found",
  "message": "Announcement not found."
}
```

<h3 id="announcementqueryservice_getannouncement-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetAnnouncementResponse](#schemav1getannouncementresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `announcement_query_not_found` — announcement with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnnouncementQueryService_ListAnnouncementsForClinic

<a id="opIdAnnouncementQueryService_ListAnnouncementsForClinic"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/clinics/{clinicId}/announcements \
  -H 'Accept: application/json'

```

`GET /v1/query/clinics/{clinicId}/announcements`

<h3 id="announcementqueryservice_listannouncementsforclinic-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|includeArchived|query|boolean|false|none|
|priority|query|string|false|none|
|limit|query|integer(int32)|false|none|
|cursor|query|string|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|priority|ANNOUNCEMENT_PRIORITY_UNSPECIFIED|
|priority|ANNOUNCEMENT_PRIORITY_NORMAL|
|priority|ANNOUNCEMENT_PRIORITY_HIGH|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "authorId": "string",
      "title": "string",
      "content": "string",
      "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
      "isArchived": true,
      "startsAt": "string",
      "endsAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "viewCount": "string"
    }
  ],
  "nextCursor": "string"
}
```

<h3 id="announcementqueryservice_listannouncementsforclinic-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListAnnouncementsForClinicResponse](#schemav1listannouncementsforclinicresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnnouncementQueryService_ListAnnouncementsForDepartment

<a id="opIdAnnouncementQueryService_ListAnnouncementsForDepartment"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/departments/{departmentId}/announcements \
  -H 'Accept: application/json'

```

`GET /v1/query/departments/{departmentId}/announcements`

<h3 id="announcementqueryservice_listannouncementsfordepartment-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|
|includeArchived|query|boolean|false|none|
|priority|query|string|false|none|
|limit|query|integer(int32)|false|none|
|cursor|query|string|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|priority|ANNOUNCEMENT_PRIORITY_UNSPECIFIED|
|priority|ANNOUNCEMENT_PRIORITY_NORMAL|
|priority|ANNOUNCEMENT_PRIORITY_HIGH|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "authorId": "string",
      "title": "string",
      "content": "string",
      "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
      "isArchived": true,
      "startsAt": "string",
      "endsAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "viewCount": "string"
    }
  ],
  "nextCursor": "string"
}
```

<h3 id="announcementqueryservice_listannouncementsfordepartment-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListAnnouncementsForDepartmentResponse](#schemav1listannouncementsfordepartmentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## AnnouncementQueryService_ListAnnouncementsForOrganization

<a id="opIdAnnouncementQueryService_ListAnnouncementsForOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/organizations/{organizationId}/announcements \
  -H 'Accept: application/json'

```

`GET /v1/query/organizations/{organizationId}/announcements`

<h3 id="announcementqueryservice_listannouncementsfororganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|includeArchived|query|boolean|false|none|
|priority|query|string|false|UNSPECIFIED = all|
|limit|query|integer(int32)|false|none|
|cursor|query|string|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|priority|ANNOUNCEMENT_PRIORITY_UNSPECIFIED|
|priority|ANNOUNCEMENT_PRIORITY_NORMAL|
|priority|ANNOUNCEMENT_PRIORITY_HIGH|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "authorId": "string",
      "title": "string",
      "content": "string",
      "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
      "isArchived": true,
      "startsAt": "string",
      "endsAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "viewCount": "string"
    }
  ],
  "nextCursor": "string"
}
```

<h3 id="announcementqueryservice_listannouncementsfororganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListAnnouncementsForOrganizationResponse](#schemav1listannouncementsfororganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-incidentclassifierqueryservice">IncidentClassifierQueryService</h1>

## IncidentClassifierQueryService_ListTypesByCategory

<a id="opIdIncidentClassifierQueryService_ListTypesByCategory"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/incident-categories/{categoryId}/types \
  -H 'Accept: application/json'

```

`GET /v1/incident-categories/{categoryId}/types`

<h3 id="incidentclassifierqueryservice_listtypesbycategory-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|categoryId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "categoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string",
      "isAllowedForPatients": true
    }
  ]
}
```

<h3 id="incidentclassifierqueryservice_listtypesbycategory-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListTypesByCategoryResponse](#schemav1listtypesbycategoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierQueryService_GetCategory

<a id="opIdIncidentClassifierQueryService_GetCategory"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/incident-categories/{id} \
  -H 'Accept: application/json'

```

`GET /v1/incident-categories/{id}`

<h3 id="incidentclassifierqueryservice_getcategory-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "category": {
    "id": "string",
    "organizationId": "string",
    "parentCategoryId": "string",
    "name": "string",
    "description": "string",
    "isActive": true,
    "createdAt": "string",
    "updatedAt": "string"
  }
}
```

> Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.

```json
{
  "code": "incident_category_not_found",
  "message": "Category not found."
}
```

<h3 id="incidentclassifierqueryservice_getcategory-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetCategoryResponse](#schemav1getcategoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_category_not_found` — category with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierQueryService_ListCategorySubtree

<a id="opIdIncidentClassifierQueryService_ListCategorySubtree"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/incident-categories/{rootCategoryId}:subtree \
  -H 'Accept: application/json'

```

`GET /v1/incident-categories/{rootCategoryId}:subtree`

<h3 id="incidentclassifierqueryservice_listcategorysubtree-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|rootCategoryId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "parentCategoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

<h3 id="incidentclassifierqueryservice_listcategorysubtree-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListCategorySubtreeResponse](#schemav1listcategorysubtreeresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierQueryService_GetType

<a id="opIdIncidentClassifierQueryService_GetType"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/incident-types/{id} \
  -H 'Accept: application/json'

```

`GET /v1/incident-types/{id}`

<h3 id="incidentclassifierqueryservice_gettype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "type": {
    "id": "string",
    "organizationId": "string",
    "categoryId": "string",
    "name": "string",
    "description": "string",
    "isActive": true,
    "createdAt": "string",
    "updatedAt": "string",
    "isAllowedForPatients": true
  }
}
```

> Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.

```json
{
  "code": "incident_type_not_found",
  "message": "Type not found."
}
```

<h3 id="incidentclassifierqueryservice_gettype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetTypeResponse](#schemav1gettyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_type_not_found` — type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierQueryService_ListCategoriesByOrganization

<a id="opIdIncidentClassifierQueryService_ListCategoriesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/incident-categories \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/incident-categories`

<h3 id="incidentclassifierqueryservice_listcategoriesbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "parentCategoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

<h3 id="incidentclassifierqueryservice_listcategoriesbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListCategoriesByOrganizationResponse](#schemav1listcategoriesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierQueryService_ListPatientVisibleCategoriesByOrganization

<a id="opIdIncidentClassifierQueryService_ListPatientVisibleCategoriesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/incident-categories:patient-visible \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/incident-categories:patient-visible`

<h3 id="incidentclassifierqueryservice_listpatientvisiblecategoriesbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "parentCategoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

<h3 id="incidentclassifierqueryservice_listpatientvisiblecategoriesbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListPatientVisibleCategoriesByOrganizationResponse](#schemav1listpatientvisiblecategoriesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierQueryService_ListActiveRootCategories

<a id="opIdIncidentClassifierQueryService_ListActiveRootCategories"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/incident-categories:roots \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/incident-categories:roots`

<h3 id="incidentclassifierqueryservice_listactiverootcategories-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "parentCategoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

<h3 id="incidentclassifierqueryservice_listactiverootcategories-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListActiveRootCategoriesResponse](#schemav1listactiverootcategoriesresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentClassifierQueryService_ListActiveTypesByOrganization

<a id="opIdIncidentClassifierQueryService_ListActiveTypesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/incident-types:active \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/incident-types:active`

<h3 id="incidentclassifierqueryservice_listactivetypesbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "categoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string",
      "isAllowedForPatients": true
    }
  ]
}
```

<h3 id="incidentclassifierqueryservice_listactivetypesbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListActiveTypesByOrganizationResponse](#schemav1listactivetypesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Patient-facing reads. These are scoped to one organisation and return
only the slice of the classifier that a patient may see when filing
an incident.

<a id="opIdIncidentClassifierQueryService_ListPatientAllowedTypesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/incident-types:patient-allowed \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/incident-types:patient-allowed`

<h3 id="patient-facing-reads.-these-are-scoped-to-one-organisation-and-return
only-the-slice-of-the-classifier-that-a-patient-may-see-when-filing
an-incident.-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "categoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string",
      "isAllowedForPatients": true
    }
  ]
}
```

<h3 id="patient-facing-reads.-these-are-scoped-to-one-organisation-and-return
only-the-slice-of-the-classifier-that-a-patient-may-see-when-filing
an-incident.-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListPatientAllowedTypesByOrganizationResponse](#schemav1listpatientallowedtypesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-incidentqueryservice">IncidentQueryService</h1>

## IncidentQueryService_GetIncident

<a id="opIdIncidentQueryService_GetIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/incidents/{id} \
  -H 'Accept: application/json'

```

`GET /v1/query/incidents/{id}`

<h3 id="incidentqueryservice_getincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "incident": {
    "id": "string",
    "organizationId": "string",
    "clinicId": "string",
    "departmentId": "string",
    "categoryId": "string",
    "typeId": "string",
    "status": "INCIDENT_STATUS_UNSPECIFIED",
    "priority": "INCIDENT_PRIORITY_UNSPECIFIED",
    "description": "string",
    "patientOriginalDescription": "string",
    "occurredAt": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "registrar": {
      "employeeId": "string",
      "displayName": "string",
      "position": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string"
    },
    "sourcePatientZitadelUserId": "string",
    "sourceBufferId": "string",
    "reopenedFromIncidentId": "string",
    "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
  }
}
```

> Not found. Error codes:
- `incident_query_not_found` — incident with the given ID does not exist.

```json
{
  "code": "incident_query_not_found",
  "message": "Incident not found."
}
```

<h3 id="incidentqueryservice_getincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetIncidentResponse](#schemav1getincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_query_not_found` — incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentQueryService_GetIncidentHistory

<a id="opIdIncidentQueryService_GetIncidentHistory"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/incidents/{incidentId}/history \
  -H 'Accept: application/json'

```

`GET /v1/query/incidents/{incidentId}/history`

<h3 id="incidentqueryservice_getincidenthistory-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|incidentId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "statusHistory": [
    {
      "id": "string",
      "oldStatus": "INCIDENT_STATUS_UNSPECIFIED",
      "newStatus": "INCIDENT_STATUS_UNSPECIFIED",
      "actor": {
        "employeeId": "string",
        "displayName": "string"
      },
      "changedAt": "string"
    }
  ],
  "priorityHistory": [
    {
      "id": "string",
      "oldPriority": "INCIDENT_PRIORITY_UNSPECIFIED",
      "newPriority": "INCIDENT_PRIORITY_UNSPECIFIED",
      "actor": {
        "employeeId": "string",
        "displayName": "string"
      },
      "changedAt": "string"
    }
  ]
}
```

> Not found. Error codes:
- `incident_query_not_found` — incident with the given ID does not exist.

```json
{
  "code": "incident_query_not_found",
  "message": "Incident not found."
}
```

<h3 id="incidentqueryservice_getincidenthistory-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetIncidentHistoryResponse](#schemav1getincidenthistoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `incident_query_not_found` — incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentQueryService_ListMyIncidents

<a id="opIdIncidentQueryService_ListMyIncidents"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/incidents:mine \
  -H 'Accept: application/json'

```

`GET /v1/query/incidents:mine`

<h3 id="incidentqueryservice_listmyincidents-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "categoryId": "string",
      "typeId": "string",
      "status": "INCIDENT_STATUS_UNSPECIFIED",
      "priority": "INCIDENT_PRIORITY_UNSPECIFIED",
      "description": "string",
      "patientOriginalDescription": "string",
      "occurredAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "registrar": {
        "employeeId": "string",
        "displayName": "string",
        "position": "string",
        "organizationId": "string",
        "clinicId": "string",
        "departmentId": "string"
      },
      "sourcePatientZitadelUserId": "string",
      "sourceBufferId": "string",
      "reopenedFromIncidentId": "string",
      "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
    }
  ]
}
```

<h3 id="incidentqueryservice_listmyincidents-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListMyIncidentsResponse](#schemav1listmyincidentsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentQueryService_ListIncidents

<a id="opIdIncidentQueryService_ListIncidents"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/organizations/{organizationId}/incidents \
  -H 'Accept: application/json'

```

`GET /v1/query/organizations/{organizationId}/incidents`

<h3 id="incidentqueryservice_listincidents-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|statuses|query|array[string]|false|none|
|priorities|query|array[string]|false|none|
|clinicId|query|string|false|none|
|departmentId|query|string|false|none|
|categoryId|query|string|false|none|
|typeId|query|string|false|none|
|occurredFrom|query|string|false|RFC3339Nano|
|occurredTo|query|string|false|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|statuses|INCIDENT_STATUS_UNSPECIFIED|
|statuses|INCIDENT_STATUS_PENDING|
|statuses|INCIDENT_STATUS_IN_PROGRESS|
|statuses|INCIDENT_STATUS_DONE|
|statuses|INCIDENT_STATUS_REJECTED|
|statuses|INCIDENT_STATUS_CANCELLED|
|priorities|INCIDENT_PRIORITY_UNSPECIFIED|
|priorities|INCIDENT_PRIORITY_LOW|
|priorities|INCIDENT_PRIORITY_NORMAL|
|priorities|INCIDENT_PRIORITY_HIGH|
|priorities|INCIDENT_PRIORITY_CRITICAL|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "categoryId": "string",
      "typeId": "string",
      "status": "INCIDENT_STATUS_UNSPECIFIED",
      "priority": "INCIDENT_PRIORITY_UNSPECIFIED",
      "description": "string",
      "patientOriginalDescription": "string",
      "occurredAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "registrar": {
        "employeeId": "string",
        "displayName": "string",
        "position": "string",
        "organizationId": "string",
        "clinicId": "string",
        "departmentId": "string"
      },
      "sourcePatientZitadelUserId": "string",
      "sourceBufferId": "string",
      "reopenedFromIncidentId": "string",
      "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
    }
  ]
}
```

<h3 id="incidentqueryservice_listincidents-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListIncidentsResponse](#schemav1listincidentsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentQueryService_ListBufferEntries

<a id="opIdIncidentQueryService_ListBufferEntries"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/organizations/{organizationId}/patient-incidents \
  -H 'Accept: application/json'

```

`GET /v1/query/organizations/{organizationId}/patient-incidents`

<h3 id="incidentqueryservice_listbufferentries-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|statuses|query|array[string]|false|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

#### Enumerated Values

|Parameter|Value|
|---|---|
|statuses|BUFFER_STATUS_UNSPECIFIED|
|statuses|BUFFER_STATUS_PENDING|
|statuses|BUFFER_STATUS_PUBLISHED|
|statuses|BUFFER_STATUS_REJECTED|
|statuses|BUFFER_STATUS_CANCELLED|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "patientZitadelUserId": "string",
      "categoryId": "string",
      "typeId": "string",
      "description": "string",
      "occurredAt": "string",
      "status": "BUFFER_STATUS_UNSPECIFIED",
      "publishedIncidentId": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
    }
  ]
}
```

<h3 id="incidentqueryservice_listbufferentries-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListBufferEntriesResponse](#schemav1listbufferentriesresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## Buffer reads.

<a id="opIdIncidentQueryService_GetBufferEntry"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/patient-incidents/{id} \
  -H 'Accept: application/json'

```

`GET /v1/query/patient-incidents/{id}`

<h3 id="buffer-reads.-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "entry": {
    "id": "string",
    "organizationId": "string",
    "patientZitadelUserId": "string",
    "categoryId": "string",
    "typeId": "string",
    "description": "string",
    "occurredAt": "string",
    "status": "BUFFER_STATUS_UNSPECIFIED",
    "publishedIncidentId": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
  }
}
```

> Not found. Error codes:
- `buffer_query_not_found` — patient incident with the given ID does not exist.

```json
{
  "code": "buffer_query_not_found",
  "message": "Patient incident not found."
}
```

<h3 id="buffer-reads.-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetBufferEntryResponse](#schemav1getbufferentryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `buffer_query_not_found` — patient incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## IncidentQueryService_ListMyBufferEntries

<a id="opIdIncidentQueryService_ListMyBufferEntries"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/query/patient-incidents:mine \
  -H 'Accept: application/json'

```

`GET /v1/query/patient-incidents:mine`

<h3 id="incidentqueryservice_listmybufferentries-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "patientZitadelUserId": "string",
      "categoryId": "string",
      "typeId": "string",
      "description": "string",
      "occurredAt": "string",
      "status": "BUFFER_STATUS_UNSPECIFIED",
      "publishedIncidentId": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
    }
  ]
}
```

<h3 id="incidentqueryservice_listmybufferentries-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListMyBufferEntriesResponse](#schemav1listmybufferentriesresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-membershipqueryservice">MembershipQueryService</h1>

## MembershipQueryService_ListEmployeesByClinic

<a id="opIdMembershipQueryService_ListEmployeesByClinic"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/clinics/{clinicId}/employees \
  -H 'Accept: application/json'

```

`GET /v1/clinics/{clinicId}/employees`

<h3 id="membershipqueryservice_listemployeesbyclinic-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|
|includeTerminated|query|boolean|false|When false (default), rows with terminated_at IS NOT NULL are|
|onVacation|query|boolean|false|When true, restrict to employees currently on an active vacation|
|position|query|string|false|Optional exact-match filter on employee_cards.position. Trimmed|

#### Detailed descriptions

**includeTerminated**: When false (default), rows with terminated_at IS NOT NULL are
hidden. Set true to include offboarded employees.

**onVacation**: When true, restrict to employees currently on an active vacation
(current_vacation_ends_at IS NOT NULL AND > now()).

**position**: Optional exact-match filter on employee_cards.position. Trimmed
before comparison; all-whitespace is treated as unset.

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  ]
}
```

<h3 id="membershipqueryservice_listemployeesbyclinic-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListEmployeesByClinicResponse](#schemav1listemployeesbyclinicresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_CountEmployeesByClinic

<a id="opIdMembershipQueryService_CountEmployeesByClinic"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/clinics/{clinicId}/employees:count \
  -H 'Accept: application/json'

```

`GET /v1/clinics/{clinicId}/employees:count`

<h3 id="membershipqueryservice_countemployeesbyclinic-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|includeTerminated|query|boolean|false|When false (default), rows with terminated_at IS NOT NULL are|
|onVacation|query|boolean|false|When true, restrict to employees currently on an active vacation|
|position|query|string|false|Optional exact-match filter on employee_cards.position. Trimmed|

#### Detailed descriptions

**includeTerminated**: When false (default), rows with terminated_at IS NOT NULL are
hidden. Set true to include offboarded employees.

**onVacation**: When true, restrict to employees currently on an active vacation
(current_vacation_ends_at IS NOT NULL AND > now()).

**position**: Optional exact-match filter on employee_cards.position. Trimmed
before comparison; all-whitespace is treated as unset.

> Example responses

> 200 Response

```json
{
  "total": "string"
}
```

<h3 id="membershipqueryservice_countemployeesbyclinic-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CountEmployeesByClinicResponse](#schemav1countemployeesbyclinicresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_GetClinicHead

<a id="opIdMembershipQueryService_GetClinicHead"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/clinics/{clinicId}/head \
  -H 'Accept: application/json'

```

`GET /v1/clinics/{clinicId}/head`

<h3 id="membershipqueryservice_getclinichead-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "assignment": {
    "holder": {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    },
    "deputy": {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  }
}
```

<h3 id="membershipqueryservice_getclinichead-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetClinicHeadResponse](#schemav1getclinicheadresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_ListEmployeesByDepartment

<a id="opIdMembershipQueryService_ListEmployeesByDepartment"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/departments/{departmentId}/employees \
  -H 'Accept: application/json'

```

`GET /v1/departments/{departmentId}/employees`

<h3 id="membershipqueryservice_listemployeesbydepartment-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|
|includeTerminated|query|boolean|false|When false (default), rows with terminated_at IS NOT NULL are|
|onVacation|query|boolean|false|When true, restrict to employees currently on an active vacation|
|position|query|string|false|Optional exact-match filter on employee_cards.position. Trimmed|

#### Detailed descriptions

**includeTerminated**: When false (default), rows with terminated_at IS NOT NULL are
hidden. Set true to include offboarded employees.

**onVacation**: When true, restrict to employees currently on an active vacation
(current_vacation_ends_at IS NOT NULL AND > now()).

**position**: Optional exact-match filter on employee_cards.position. Trimmed
before comparison; all-whitespace is treated as unset.

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  ]
}
```

<h3 id="membershipqueryservice_listemployeesbydepartment-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListEmployeesByDepartmentResponse](#schemav1listemployeesbydepartmentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_CountEmployeesByDepartment

<a id="opIdMembershipQueryService_CountEmployeesByDepartment"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/departments/{departmentId}/employees:count \
  -H 'Accept: application/json'

```

`GET /v1/departments/{departmentId}/employees:count`

<h3 id="membershipqueryservice_countemployeesbydepartment-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|
|includeTerminated|query|boolean|false|When false (default), rows with terminated_at IS NOT NULL are|
|onVacation|query|boolean|false|When true, restrict to employees currently on an active vacation|
|position|query|string|false|Optional exact-match filter on employee_cards.position. Trimmed|

#### Detailed descriptions

**includeTerminated**: When false (default), rows with terminated_at IS NOT NULL are
hidden. Set true to include offboarded employees.

**onVacation**: When true, restrict to employees currently on an active vacation
(current_vacation_ends_at IS NOT NULL AND > now()).

**position**: Optional exact-match filter on employee_cards.position. Trimmed
before comparison; all-whitespace is treated as unset.

> Example responses

> 200 Response

```json
{
  "total": "string"
}
```

<h3 id="membershipqueryservice_countemployeesbydepartment-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CountEmployeesByDepartmentResponse](#schemav1countemployeesbydepartmentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_GetDepartmentResponsible

<a id="opIdMembershipQueryService_GetDepartmentResponsible"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/departments/{departmentId}/responsible \
  -H 'Accept: application/json'

```

`GET /v1/departments/{departmentId}/responsible`

<h3 id="membershipqueryservice_getdepartmentresponsible-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "assignment": {
    "holder": {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    },
    "deputy": {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  }
}
```

<h3 id="membershipqueryservice_getdepartmentresponsible-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetDepartmentResponsibleResponse](#schemav1getdepartmentresponsibleresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_ListVacationsByEmployee

<a id="opIdMembershipQueryService_ListVacationsByEmployee"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/employees/{employeeId}/vacations \
  -H 'Accept: application/json'

```

`GET /v1/employees/{employeeId}/vacations`

<h3 id="membershipqueryservice_listvacationsbyemployee-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|employeeId|path|string|true|none|
|state|query|string|false|Optional state filter. When empty, all states are returned.|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

#### Detailed descriptions

**state**: Optional state filter. When empty, all states are returned.
Valid values: scheduled, active, ended, cancelled.

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "employeeId": "string",
      "state": "string",
      "startsAt": "string",
      "endsAt": "string",
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

<h3 id="membershipqueryservice_listvacationsbyemployee-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListVacationsByEmployeeResponse](#schemav1listvacationsbyemployeeresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_CountVacationsByEmployee

<a id="opIdMembershipQueryService_CountVacationsByEmployee"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/employees/{employeeId}/vacations:count \
  -H 'Accept: application/json'

```

`GET /v1/employees/{employeeId}/vacations:count`

<h3 id="membershipqueryservice_countvacationsbyemployee-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|employeeId|path|string|true|none|
|state|query|string|false|Optional state filter. When empty, counts all states.|

#### Detailed descriptions

**state**: Optional state filter. When empty, counts all states.
Valid values: scheduled, active, ended, cancelled.

> Example responses

> 200 Response

```json
{
  "total": "string"
}
```

<h3 id="membershipqueryservice_countvacationsbyemployee-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CountVacationsByEmployeeResponse](#schemav1countvacationsbyemployeeresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_GetEmployee

<a id="opIdMembershipQueryService_GetEmployee"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/employees/{id} \
  -H 'Accept: application/json'

```

`GET /v1/employees/{id}`

<h3 id="membershipqueryservice_getemployee-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "employee": {
    "employeeId": "string",
    "zitadelUserId": "string",
    "firstName": "string",
    "lastName": "string",
    "displayName": "string",
    "email": "string",
    "organizationId": "string",
    "organizationName": "string",
    "clinicId": "string",
    "clinicName": "string",
    "departmentId": "string",
    "departmentName": "string",
    "position": "string",
    "terminatedAt": "string",
    "currentVacationEndsAt": "string",
    "nextVacationStartsAt": "string"
  }
}
```

> Not found. Error codes:
- `employee_card_not_found` — employee with the given ID does not exist.

```json
{
  "code": "employee_card_not_found",
  "message": "Employee not found."
}
```

<h3 id="membershipqueryservice_getemployee-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetEmployeeResponse](#schemav1getemployeeresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `employee_card_not_found` — employee with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_ListOrgAdmins

<a id="opIdMembershipQueryService_ListOrgAdmins"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/admins \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/admins`

<h3 id="membershipqueryservice_listorgadmins-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "holder": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      },
      "deputy": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      }
    }
  ]
}
```

<h3 id="membershipqueryservice_listorgadmins-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListOrgAdminsResponse](#schemav1listorgadminsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_ListOrgDispatchers

<a id="opIdMembershipQueryService_ListOrgDispatchers"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/dispatchers \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/dispatchers`

<h3 id="membershipqueryservice_listorgdispatchers-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "holder": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      },
      "deputy": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      }
    }
  ]
}
```

<h3 id="membershipqueryservice_listorgdispatchers-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListOrgDispatchersResponse](#schemav1listorgdispatchersresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_ListEmployeesByOrganization

<a id="opIdMembershipQueryService_ListEmployeesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/employees \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/employees`

<h3 id="membershipqueryservice_listemployeesbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|
|includeTerminated|query|boolean|false|When false (default), rows with terminated_at IS NOT NULL are|
|onVacation|query|boolean|false|When true, restrict to employees currently on an active vacation|
|position|query|string|false|Optional exact-match filter on employee_cards.position. Trimmed|

#### Detailed descriptions

**includeTerminated**: When false (default), rows with terminated_at IS NOT NULL are
hidden. Set true to include offboarded employees.

**onVacation**: When true, restrict to employees currently on an active vacation
(current_vacation_ends_at IS NOT NULL AND > now()).

**position**: Optional exact-match filter on employee_cards.position. Trimmed
before comparison; all-whitespace is treated as unset.

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  ]
}
```

<h3 id="membershipqueryservice_listemployeesbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListEmployeesByOrganizationResponse](#schemav1listemployeesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_CountEmployeesByOrganization

<a id="opIdMembershipQueryService_CountEmployeesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/employees:count \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/employees:count`

<h3 id="membershipqueryservice_countemployeesbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|includeTerminated|query|boolean|false|When false (default), rows with terminated_at IS NOT NULL are|
|onVacation|query|boolean|false|When true, restrict to employees currently on an active vacation|
|position|query|string|false|Optional exact-match filter on employee_cards.position. Trimmed|

#### Detailed descriptions

**includeTerminated**: When false (default), rows with terminated_at IS NOT NULL are
hidden. Set true to include offboarded employees.

**onVacation**: When true, restrict to employees currently on an active vacation
(current_vacation_ends_at IS NOT NULL AND > now()).

**position**: Optional exact-match filter on employee_cards.position. Trimmed
before comparison; all-whitespace is treated as unset.

> Example responses

> 200 Response

```json
{
  "total": "string"
}
```

<h3 id="membershipqueryservice_countemployeesbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CountEmployeesByOrganizationResponse](#schemav1countemployeesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_SearchEmployeesByOrganization

<a id="opIdMembershipQueryService_SearchEmployeesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/employees:search \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/employees:search`

<h3 id="membershipqueryservice_searchemployeesbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|query|query|string|false|Fuzzy substring matched case-insensitively (ILIKE %query%) against|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|
|includeTerminated|query|boolean|false|When false (default), rows with terminated_at IS NOT NULL are|
|onVacation|query|boolean|false|When true, restrict to employees currently on an active vacation|
|position|query|string|false|Optional exact-match filter on employee_cards.position. Trimmed|

#### Detailed descriptions

**query**: Fuzzy substring matched case-insensitively (ILIKE %query%) against
first_name, last_name, display_name, and email. Trimmed at the
handler boundary; an empty query degenerates to the same behaviour
as ListEmployeesByOrganization with the same filters. Maximum
length 256 characters; over-long inputs are rejected before the
authz round-trip.

**includeTerminated**: When false (default), rows with terminated_at IS NOT NULL are
hidden. Set true to include offboarded employees.

**onVacation**: When true, restrict to employees currently on an active vacation
(current_vacation_ends_at IS NOT NULL AND > now()).

**position**: Optional exact-match filter on employee_cards.position. Trimmed
before comparison; all-whitespace is treated as unset.

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  ]
}
```

<h3 id="membershipqueryservice_searchemployeesbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1SearchEmployeesByOrganizationResponse](#schemav1searchemployeesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_ListOrgHeads

<a id="opIdMembershipQueryService_ListOrgHeads"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/heads \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/heads`

<h3 id="membershipqueryservice_listorgheads-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "holder": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      },
      "deputy": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      }
    }
  ]
}
```

<h3 id="membershipqueryservice_listorgheads-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListOrgHeadsResponse](#schemav1listorgheadsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## MembershipQueryService_ListSystemAdmins

<a id="opIdMembershipQueryService_ListSystemAdmins"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/system-admins \
  -H 'Accept: application/json'

```

`GET /v1/system-admins`

<h3 id="membershipqueryservice_listsystemadmins-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "zitadelUserId": "string",
      "createdAt": "string"
    }
  ]
}
```

<h3 id="membershipqueryservice_listsystemadmins-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListSystemAdminsResponse](#schemav1listsystemadminsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-orgstructurequeryservice">OrgStructureQueryService</h1>

## OrgStructureQueryService_ListDepartmentsByClinic

<a id="opIdOrgStructureQueryService_ListDepartmentsByClinic"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/clinics/{clinicId}/departments \
  -H 'Accept: application/json'

```

`GET /v1/clinics/{clinicId}/departments`

<h3 id="orgstructurequeryservice_listdepartmentsbyclinic-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "clinicId": "string",
      "name": "string"
    }
  ]
}
```

<h3 id="orgstructurequeryservice_listdepartmentsbyclinic-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListDepartmentsByClinicResponse](#schemav1listdepartmentsbyclinicresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_CountDepartmentsByClinic

<a id="opIdOrgStructureQueryService_CountDepartmentsByClinic"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/clinics/{clinicId}/departments:count \
  -H 'Accept: application/json'

```

`GET /v1/clinics/{clinicId}/departments:count`

<h3 id="orgstructurequeryservice_countdepartmentsbyclinic-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "total": "string"
}
```

<h3 id="orgstructurequeryservice_countdepartmentsbyclinic-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CountDepartmentsByClinicResponse](#schemav1countdepartmentsbyclinicresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_GetClinic

<a id="opIdOrgStructureQueryService_GetClinic"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/clinics/{id} \
  -H 'Accept: application/json'

```

`GET /v1/clinics/{id}`

<h3 id="orgstructurequeryservice_getclinic-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "clinic": {
    "id": "string",
    "organizationId": "string",
    "name": "string",
    "description": "string",
    "physicalAddress": {
      "text": "string",
      "point": {
        "longitude": 0.1,
        "latitude": 0.1
      }
    },
    "createdAt": "string",
    "updatedAt": "string"
  }
}
```

> Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.

```json
{
  "code": "clinic_not_found",
  "message": "Clinic not found."
}
```

<h3 id="orgstructurequeryservice_getclinic-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetClinicResponse](#schemav1getclinicresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `clinic_not_found` — clinic with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_GetDepartment

<a id="opIdOrgStructureQueryService_GetDepartment"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/departments/{id} \
  -H 'Accept: application/json'

```

`GET /v1/departments/{id}`

<h3 id="orgstructurequeryservice_getdepartment-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "department": {
    "id": "string",
    "clinicId": "string",
    "name": "string",
    "description": "string",
    "createdAt": "string",
    "updatedAt": "string"
  }
}
```

> Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.

```json
{
  "code": "department_not_found",
  "message": "Department not found."
}
```

<h3 id="orgstructurequeryservice_getdepartment-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetDepartmentResponse](#schemav1getdepartmentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `department_not_found` — department with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_ListOrganizations

<a id="opIdOrgStructureQueryService_ListOrganizations"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations \
  -H 'Accept: application/json'

```

`GET /v1/organizations`

<h3 id="orgstructurequeryservice_listorganizations-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "name": "string"
    }
  ]
}
```

<h3 id="orgstructurequeryservice_listorganizations-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListOrganizationsResponse](#schemav1listorganizationsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_GetOrganization

<a id="opIdOrgStructureQueryService_GetOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{id} \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{id}`

<h3 id="orgstructurequeryservice_getorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "organization": {
    "id": "string",
    "name": "string",
    "description": "string",
    "legalAddress": {
      "text": "string",
      "point": {
        "longitude": 0.1,
        "latitude": 0.1
      }
    },
    "createdAt": "string",
    "updatedAt": "string"
  }
}
```

> Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.

```json
{
  "code": "organization_not_found",
  "message": "Organization not found."
}
```

<h3 id="orgstructurequeryservice_getorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetOrganizationResponse](#schemav1getorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `organization_not_found` — organization with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_ListClinicsByOrganization

<a id="opIdOrgStructureQueryService_ListClinicsByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/clinics \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/clinics`

<h3 id="orgstructurequeryservice_listclinicsbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "name": "string"
    }
  ]
}
```

<h3 id="orgstructurequeryservice_listclinicsbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListClinicsByOrganizationResponse](#schemav1listclinicsbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_CountClinicsByOrganization

<a id="opIdOrgStructureQueryService_CountClinicsByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/clinics:count \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/clinics:count`

<h3 id="orgstructurequeryservice_countclinicsbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "total": "string"
}
```

<h3 id="orgstructurequeryservice_countclinicsbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CountClinicsByOrganizationResponse](#schemav1countclinicsbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_CountOrganizations

<a id="opIdOrgStructureQueryService_CountOrganizations"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations:count \
  -H 'Accept: application/json'

```

`GET /v1/organizations:count`

> Example responses

> 200 Response

```json
{
  "total": "string"
}
```

<h3 id="orgstructurequeryservice_countorganizations-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1CountOrganizationsResponse](#schemav1countorganizationsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## OrgStructureQueryService_SearchOrganizations

<a id="opIdOrgStructureQueryService_SearchOrganizations"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations:search \
  -H 'Accept: application/json'

```

`GET /v1/organizations:search`

<h3 id="orgstructurequeryservice_searchorganizations-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|query|query|string|false|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "name": "string"
    }
  ]
}
```

<h3 id="orgstructurequeryservice_searchorganizations-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1SearchOrganizationsResponse](#schemav1searchorganizationsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-requestclassifierqueryservice">RequestClassifierQueryService</h1>

## RequestClassifierQueryService_ListRequestTypesByOrganization

<a id="opIdRequestClassifierQueryService_ListRequestTypesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/request-types \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/request-types`

<h3 id="requestclassifierqueryservice_listrequesttypesbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

<h3 id="requestclassifierqueryservice_listrequesttypesbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListRequestTypesByOrganizationResponse](#schemav1listrequesttypesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## RequestClassifierQueryService_ListActiveRequestTypesByOrganization

<a id="opIdRequestClassifierQueryService_ListActiveRequestTypesByOrganization"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/request-types:active \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/request-types:active`

<h3 id="requestclassifierqueryservice_listactiverequesttypesbyorganization-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

<h3 id="requestclassifierqueryservice_listactiverequesttypesbyorganization-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListActiveRequestTypesByOrganizationResponse](#schemav1listactiverequesttypesbyorganizationresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## RequestClassifierQueryService_GetRequestType

<a id="opIdRequestClassifierQueryService_GetRequestType"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/request-types/{id} \
  -H 'Accept: application/json'

```

`GET /v1/request-types/{id}`

<h3 id="requestclassifierqueryservice_getrequesttype-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "requestType": {
    "id": "string",
    "organizationId": "string",
    "name": "string",
    "description": "string",
    "isActive": true,
    "createdAt": "string",
    "updatedAt": "string"
  }
}
```

> Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.

```json
{
  "code": "request_type_not_found",
  "message": "Request type not found."
}
```

<h3 id="requestclassifierqueryservice_getrequesttype-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetRequestTypeResponse](#schemav1getrequesttyperesponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `request_type_not_found` — request type with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-servicerequestqueryservice">ServiceRequestQueryService</h1>

## ServiceRequestQueryService_ListServiceRequestsByIncident

<a id="opIdServiceRequestQueryService_ListServiceRequestsByIncident"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/incidents/{incidentId}/service-requests \
  -H 'Accept: application/json'

```

`GET /v1/incidents/{incidentId}/service-requests`

<h3 id="servicerequestqueryservice_listservicerequestsbyincident-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|incidentId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "typeId": "string",
      "incidentId": "string",
      "description": "string",
      "status": "string",
      "authorId": "string",
      "authorDisplayName": "string",
      "executors": [
        {
          "employeeId": "string",
          "assignedAt": "string",
          "assignedById": "string"
        }
      ],
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

> Not found. Error codes:
- `service_request_query_incident_not_found` — incident with the given ID does not exist.

```json
{
  "code": "service_request_query_incident_not_found",
  "message": "Incident not found."
}
```

<h3 id="servicerequestqueryservice_listservicerequestsbyincident-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListServiceRequestsByIncidentResponse](#schemav1listservicerequestsbyincidentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `service_request_query_incident_not_found` — incident with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ServiceRequestQueryService_ListServiceRequests

<a id="opIdServiceRequestQueryService_ListServiceRequests"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/service-requests \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/service-requests`

<h3 id="servicerequestqueryservice_listservicerequests-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|
|limit|query|integer(int32)|false|none|
|offset|query|integer(int32)|false|none|

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "typeId": "string",
      "incidentId": "string",
      "description": "string",
      "status": "string",
      "authorId": "string",
      "authorDisplayName": "string",
      "executors": [
        {
          "employeeId": "string",
          "assignedAt": "string",
          "assignedById": "string"
        }
      ],
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}
```

<h3 id="servicerequestqueryservice_listservicerequests-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListServiceRequestsResponse](#schemav1listservicerequestsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ServiceRequestQueryService_GetServiceRequest

<a id="opIdServiceRequestQueryService_GetServiceRequest"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/service-requests/{id} \
  -H 'Accept: application/json'

```

`GET /v1/service-requests/{id}`

<h3 id="servicerequestqueryservice_getservicerequest-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|id|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "serviceRequest": {
    "id": "string",
    "organizationId": "string",
    "clinicId": "string",
    "departmentId": "string",
    "typeId": "string",
    "incidentId": "string",
    "description": "string",
    "status": "string",
    "authorId": "string",
    "authorDisplayName": "string",
    "executors": [
      {
        "employeeId": "string",
        "assignedAt": "string",
        "assignedById": "string"
      }
    ],
    "createdAt": "string",
    "updatedAt": "string"
  }
}
```

> Not found. Error codes:
- `service_request_query_not_found` — service request with the given ID does not exist.

```json
{
  "code": "service_request_query_not_found",
  "message": "Service request not found."
}
```

<h3 id="servicerequestqueryservice_getservicerequest-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetServiceRequestResponse](#schemav1getservicerequestresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `service_request_query_not_found` — service request with the given ID does not exist.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ServiceRequestQueryService_GetServiceRequestHistory

<a id="opIdServiceRequestQueryService_GetServiceRequestHistory"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/service-requests/{serviceRequestId}/history \
  -H 'Accept: application/json'

```

`GET /v1/service-requests/{serviceRequestId}/history`

<h3 id="servicerequestqueryservice_getservicerequesthistory-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|serviceRequestId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "statusHistory": [
    {
      "id": "string",
      "oldStatus": "string",
      "newStatus": "string",
      "actorId": "string",
      "actorName": "string",
      "changedAt": "string"
    }
  ],
  "executorHistory": [
    {
      "id": "string",
      "action": "string",
      "employeeId": "string",
      "employeeName": "string",
      "actorId": "string",
      "actorName": "string",
      "changedAt": "string"
    }
  ]
}
```

<h3 id="servicerequestqueryservice_getservicerequesthistory-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetServiceRequestHistoryResponse](#schemav1getservicerequesthistoryresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-selfqueryservice">SelfQueryService</h1>

## GetMyIdentity returns whether the caller is a system administrator.

<a id="opIdSelfQueryService_GetMyIdentity"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/me \
  -H 'Accept: application/json'

```

`GET /v1/me`

> Example responses

> 200 Response

```json
{
  "isSystemAdmin": true
}
```

<h3 id="getmyidentity-returns-whether-the-caller-is-a-system-administrator.-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetMyIdentityResponse](#schemav1getmyidentityresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetMyClinicRole returns whether the caller is the clinic head of
the given clinic. Returns NOT_FOUND when the caller is not an
active employee of that clinic.

<a id="opIdSelfQueryService_GetMyClinicRole"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/me/clinics/{clinicId}/role \
  -H 'Accept: application/json'

```

`GET /v1/me/clinics/{clinicId}/role`

<h3 id="getmyclinicrole-returns-whether-the-caller-is-the-clinic-head-of
the-given-clinic.-returns-not_found-when-the-caller-is-not-an
active-employee-of-that-clinic.-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "isClinicHead": true
}
```

> Not found. Error codes:
- `self_clinic_role_not_found` — caller is not an active employee of the given clinic.

```json
{
  "code": "self_clinic_role_not_found",
  "message": "Clinic role not found."
}
```

<h3 id="getmyclinicrole-returns-whether-the-caller-is-the-clinic-head-of
the-given-clinic.-returns-not_found-when-the-caller-is-not-an
active-employee-of-that-clinic.-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetMyClinicRoleResponse](#schemav1getmyclinicroleresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `self_clinic_role_not_found` — caller is not an active employee of the given clinic.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetMyDepartmentRole returns whether the caller holds the
department-responsible role for the given department. Returns
NOT_FOUND when the caller is not an active employee of that
department.

<a id="opIdSelfQueryService_GetMyDepartmentRole"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/me/departments/{departmentId}/role \
  -H 'Accept: application/json'

```

`GET /v1/me/departments/{departmentId}/role`

<h3 id="getmydepartmentrole-returns-whether-the-caller-holds-the
department-responsible-role-for-the-given-department.-returns
not_found-when-the-caller-is-not-an-active-employee-of-that
department.-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "isDepartmentResponsible": true
}
```

> Not found. Error codes:
- `self_dept_role_not_found` — caller is not an active employee of the given department.

```json
{
  "code": "self_dept_role_not_found",
  "message": "Department role not found."
}
```

<h3 id="getmydepartmentrole-returns-whether-the-caller-holds-the
department-responsible-role-for-the-given-department.-returns
not_found-when-the-caller-is-not-an-active-employee-of-that
department.-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetMyDepartmentRoleResponse](#schemav1getmydepartmentroleresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `self_dept_role_not_found` — caller is not an active employee of the given department.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## ListMyOrganizations returns every organization where the caller
has an active (non-terminated) employee record.

<a id="opIdSelfQueryService_ListMyOrganizations"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/me/organizations \
  -H 'Accept: application/json'

```

`GET /v1/me/organizations`

> Example responses

> 200 Response

```json
{
  "items": [
    {
      "id": "string",
      "name": "string"
    }
  ]
}
```

<h3 id="listmyorganizations-returns-every-organization-where-the-caller
has-an-active-(non-terminated)-employee-record.-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1ListMyOrganizationsResponse](#schemav1listmyorganizationsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetMyEmployment returns the caller's employee card in the given
organization. Returns NOT_FOUND when the caller is not an active
employee of that organization.

<a id="opIdSelfQueryService_GetMyEmployment"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/me/organizations/{organizationId}/employment \
  -H 'Accept: application/json'

```

`GET /v1/me/organizations/{organizationId}/employment`

<h3 id="getmyemployment-returns-the-caller's-employee-card-in-the-given
organization.-returns-not_found-when-the-caller-is-not-an-active
employee-of-that-organization.-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "employee": {
    "employeeId": "string",
    "zitadelUserId": "string",
    "firstName": "string",
    "lastName": "string",
    "displayName": "string",
    "email": "string",
    "organizationId": "string",
    "organizationName": "string",
    "clinicId": "string",
    "clinicName": "string",
    "departmentId": "string",
    "departmentName": "string",
    "position": "string",
    "terminatedAt": "string",
    "currentVacationEndsAt": "string",
    "nextVacationStartsAt": "string"
  }
}
```

> Not found. Error codes:
- `self_employment_not_found` — caller is not an active employee of the given organization.

```json
{
  "code": "self_employment_not_found",
  "message": "Employment not found."
}
```

<h3 id="getmyemployment-returns-the-caller's-employee-card-in-the-given
organization.-returns-not_found-when-the-caller-is-not-an-active
employee-of-that-organization.-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetMyEmploymentResponse](#schemav1getmyemploymentresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `self_employment_not_found` — caller is not an active employee of the given organization.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## GetMyOrganizationRole returns the caller's named roles in the
given organization. Returns NOT_FOUND when the caller is not an
active employee of that organization.

<a id="opIdSelfQueryService_GetMyOrganizationRole"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/me/organizations/{organizationId}/role \
  -H 'Accept: application/json'

```

`GET /v1/me/organizations/{organizationId}/role`

<h3 id="getmyorganizationrole-returns-the-caller's-named-roles-in-the
given-organization.-returns-not_found-when-the-caller-is-not-an
active-employee-of-that-organization.-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "isOrgAdmin": true,
  "isOrgHead": true,
  "isOrgDispatcher": true
}
```

> Not found. Error codes:
- `self_org_role_not_found` — caller is not an active employee of the given organization.

```json
{
  "code": "self_org_role_not_found",
  "message": "Organization role not found."
}
```

<h3 id="getmyorganizationrole-returns-the-caller's-named-roles-in-the
given-organization.-returns-not_found-when-the-caller-is-not-an
active-employee-of-that-organization.-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetMyOrganizationRoleResponse](#schemav1getmyorganizationroleresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not found. Error codes:
- `self_org_role_not_found` — caller is not an active employee of the given organization.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

<h1 id="command-announcement-v1-announcement-proto-statsqueryservice">StatsQueryService</h1>

## StatsQueryService_GetClinicStats

<a id="opIdStatsQueryService_GetClinicStats"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/clinics/{clinicId}/stats \
  -H 'Accept: application/json'

```

`GET /v1/clinics/{clinicId}/stats`

<h3 id="statsqueryservice_getclinicstats-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|clinicId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "stats": {
    "clinicId": "string",
    "organizationId": "string",
    "employeesTotal": "string",
    "departmentsTotal": "string",
    "employeesOnVacation": "string"
  }
}
```

<h3 id="statsqueryservice_getclinicstats-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetClinicStatsResponse](#schemav1getclinicstatsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## StatsQueryService_GetDepartmentStats

<a id="opIdStatsQueryService_GetDepartmentStats"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/departments/{departmentId}/stats \
  -H 'Accept: application/json'

```

`GET /v1/departments/{departmentId}/stats`

<h3 id="statsqueryservice_getdepartmentstats-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|departmentId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "stats": {
    "departmentId": "string",
    "clinicId": "string",
    "organizationId": "string",
    "employeesTotal": "string",
    "employeesOnVacation": "string"
  }
}
```

<h3 id="statsqueryservice_getdepartmentstats-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetDepartmentStatsResponse](#schemav1getdepartmentstatsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

## StatsQueryService_GetOrganizationStats

<a id="opIdStatsQueryService_GetOrganizationStats"></a>

> Code samples

```shell
# You can also use wget
curl -X GET /v1/organizations/{organizationId}/stats \
  -H 'Accept: application/json'

```

`GET /v1/organizations/{organizationId}/stats`

<h3 id="statsqueryservice_getorganizationstats-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|organizationId|path|string|true|none|

> Example responses

> 200 Response

```json
{
  "stats": {
    "organizationId": "string",
    "employeesTotal": "string",
    "clinicsTotal": "string",
    "departmentsTotal": "string",
    "employeesOnVacation": "string",
    "vacationsScheduled": "string"
  }
}
```

<h3 id="statsqueryservice_getorganizationstats-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|A successful response.|[v1GetOrganizationStatsResponse](#schemav1getorganizationstatsresponse)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Validation failed or invalid input.|[v1ValidationErrorResponse](#schemav1validationerrorresponse)|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthenticated — missing or invalid token.|[v1ErrorResponse](#schemav1errorresponse)|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Permission denied.|[v1ErrorResponse](#schemav1errorresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Unexpected server error.|[v1ErrorResponse](#schemav1errorresponse)|
|default|Default|An unexpected error response.|[v1ErrorResponse](#schemav1errorresponse)|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_AnnouncementCommandServiceUpdateAnnouncementBody">AnnouncementCommandServiceUpdateAnnouncementBody</h2>
<!-- backwards compatibility -->
<a id="schemaannouncementcommandserviceupdateannouncementbody"></a>
<a id="schema_AnnouncementCommandServiceUpdateAnnouncementBody"></a>
<a id="tocSannouncementcommandserviceupdateannouncementbody"></a>
<a id="tocsannouncementcommandserviceupdateannouncementbody"></a>

```json
{
  "title": "string",
  "content": "string",
  "startsAt": "string",
  "endsAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|title|string|true|none|none|
|content|string|true|none|none|
|startsAt|string|false|none|none|
|endsAt|string|false|none|none|

<h2 id="tocS_AnnouncementCommandServiceUpdateAnnouncementPriorityBody">AnnouncementCommandServiceUpdateAnnouncementPriorityBody</h2>
<!-- backwards compatibility -->
<a id="schemaannouncementcommandserviceupdateannouncementprioritybody"></a>
<a id="schema_AnnouncementCommandServiceUpdateAnnouncementPriorityBody"></a>
<a id="tocSannouncementcommandserviceupdateannouncementprioritybody"></a>
<a id="tocsannouncementcommandserviceupdateannouncementprioritybody"></a>

```json
{
  "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|priority|[commandAnnouncementV1AnnouncementPriority](#schemacommandannouncementv1announcementpriority)|true|none|none|

<h2 id="tocS_IncidentBufferCommandServicePublishPatientIncidentBody">IncidentBufferCommandServicePublishPatientIncidentBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentbuffercommandservicepublishpatientincidentbody"></a>
<a id="schema_IncidentBufferCommandServicePublishPatientIncidentBody"></a>
<a id="tocSincidentbuffercommandservicepublishpatientincidentbody"></a>
<a id="tocsincidentbuffercommandservicepublishpatientincidentbody"></a>

```json
{
  "departmentId": "string",
  "categoryId": "string",
  "typeId": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|departmentId|string|true|none|none|
|categoryId|string|true|none|none|
|typeId|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_IncidentBufferCommandServiceUpdatePatientIncidentBody">IncidentBufferCommandServiceUpdatePatientIncidentBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentbuffercommandserviceupdatepatientincidentbody"></a>
<a id="schema_IncidentBufferCommandServiceUpdatePatientIncidentBody"></a>
<a id="tocSincidentbuffercommandserviceupdatepatientincidentbody"></a>
<a id="tocsincidentbuffercommandserviceupdatepatientincidentbody"></a>

```json
{
  "categoryId": "string",
  "typeId": "string",
  "description": "string",
  "occurredAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|categoryId|string|false|none|none|
|typeId|string|false|none|none|
|description|string|false|none|none|
|occurredAt|string|false|none|none|

<h2 id="tocS_IncidentClassifierCommandServiceCreateIncidentCategoryBody">IncidentClassifierCommandServiceCreateIncidentCategoryBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentclassifiercommandservicecreateincidentcategorybody"></a>
<a id="schema_IncidentClassifierCommandServiceCreateIncidentCategoryBody"></a>
<a id="tocSincidentclassifiercommandservicecreateincidentcategorybody"></a>
<a id="tocsincidentclassifiercommandservicecreateincidentcategorybody"></a>

```json
{
  "parentCategoryId": "string",
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|parentCategoryId|string|false|none|none|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_IncidentClassifierCommandServiceCreateIncidentTypeBody">IncidentClassifierCommandServiceCreateIncidentTypeBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentclassifiercommandservicecreateincidenttypebody"></a>
<a id="schema_IncidentClassifierCommandServiceCreateIncidentTypeBody"></a>
<a id="tocSincidentclassifiercommandservicecreateincidenttypebody"></a>
<a id="tocsincidentclassifiercommandservicecreateincidenttypebody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_IncidentClassifierCommandServiceMoveIncidentCategoryBody">IncidentClassifierCommandServiceMoveIncidentCategoryBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentclassifiercommandservicemoveincidentcategorybody"></a>
<a id="schema_IncidentClassifierCommandServiceMoveIncidentCategoryBody"></a>
<a id="tocSincidentclassifiercommandservicemoveincidentcategorybody"></a>
<a id="tocsincidentclassifiercommandservicemoveincidentcategorybody"></a>

```json
{
  "newParentCategoryId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|newParentCategoryId|string|false|none|none|

<h2 id="tocS_IncidentClassifierCommandServiceMoveIncidentTypeBody">IncidentClassifierCommandServiceMoveIncidentTypeBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentclassifiercommandservicemoveincidenttypebody"></a>
<a id="schema_IncidentClassifierCommandServiceMoveIncidentTypeBody"></a>
<a id="tocSincidentclassifiercommandservicemoveincidenttypebody"></a>
<a id="tocsincidentclassifiercommandservicemoveincidenttypebody"></a>

```json
{
  "newCategoryId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|newCategoryId|string|true|none|none|

<h2 id="tocS_IncidentClassifierCommandServiceUpdateIncidentCategoryDetailsBody">IncidentClassifierCommandServiceUpdateIncidentCategoryDetailsBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentclassifiercommandserviceupdateincidentcategorydetailsbody"></a>
<a id="schema_IncidentClassifierCommandServiceUpdateIncidentCategoryDetailsBody"></a>
<a id="tocSincidentclassifiercommandserviceupdateincidentcategorydetailsbody"></a>
<a id="tocsincidentclassifiercommandserviceupdateincidentcategorydetailsbody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_IncidentClassifierCommandServiceUpdateIncidentTypeDetailsBody">IncidentClassifierCommandServiceUpdateIncidentTypeDetailsBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentclassifiercommandserviceupdateincidenttypedetailsbody"></a>
<a id="schema_IncidentClassifierCommandServiceUpdateIncidentTypeDetailsBody"></a>
<a id="tocSincidentclassifiercommandserviceupdateincidenttypedetailsbody"></a>
<a id="tocsincidentclassifiercommandserviceupdateincidenttypedetailsbody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_IncidentCommandServiceUpdateIncidentDescriptionBody">IncidentCommandServiceUpdateIncidentDescriptionBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentcommandserviceupdateincidentdescriptionbody"></a>
<a id="schema_IncidentCommandServiceUpdateIncidentDescriptionBody"></a>
<a id="tocSincidentcommandserviceupdateincidentdescriptionbody"></a>
<a id="tocsincidentcommandserviceupdateincidentdescriptionbody"></a>

```json
{
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|description|string|false|none|none|

<h2 id="tocS_IncidentCommandServiceUpdateIncidentPriorityBody">IncidentCommandServiceUpdateIncidentPriorityBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentcommandserviceupdateincidentprioritybody"></a>
<a id="schema_IncidentCommandServiceUpdateIncidentPriorityBody"></a>
<a id="tocSincidentcommandserviceupdateincidentprioritybody"></a>
<a id="tocsincidentcommandserviceupdateincidentprioritybody"></a>

```json
{
  "priority": "INCIDENT_PRIORITY_UNSPECIFIED"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|priority|[commandIncidentV1IncidentPriority](#schemacommandincidentv1incidentpriority)|true|none|none|

<h2 id="tocS_IncidentCommandServiceUpdateIncidentStatusBody">IncidentCommandServiceUpdateIncidentStatusBody</h2>
<!-- backwards compatibility -->
<a id="schemaincidentcommandserviceupdateincidentstatusbody"></a>
<a id="schema_IncidentCommandServiceUpdateIncidentStatusBody"></a>
<a id="tocSincidentcommandserviceupdateincidentstatusbody"></a>
<a id="tocsincidentcommandserviceupdateincidentstatusbody"></a>

```json
{
  "newStatus": "INCIDENT_STATUS_UNSPECIFIED"
}

```

UpdateIncidentStatus only handles forward transitions:
pending -> in_progress, in_progress -> done, in_progress -> rejected.
Cancellation by registrar uses CancelIncident.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|newStatus|[commandIncidentV1IncidentStatus](#schemacommandincidentv1incidentstatus)|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignClinicHeadBody">MembershipCommandServiceAssignClinicHeadBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassignclinicheadbody"></a>
<a id="schema_MembershipCommandServiceAssignClinicHeadBody"></a>
<a id="tocSmembershipcommandserviceassignclinicheadbody"></a>
<a id="tocsmembershipcommandserviceassignclinicheadbody"></a>

```json
{
  "employeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignClinicHeadDeputyBody">MembershipCommandServiceAssignClinicHeadDeputyBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassignclinicheaddeputybody"></a>
<a id="schema_MembershipCommandServiceAssignClinicHeadDeputyBody"></a>
<a id="tocSmembershipcommandserviceassignclinicheaddeputybody"></a>
<a id="tocsmembershipcommandserviceassignclinicheaddeputybody"></a>

```json
{
  "deputyEmployeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|deputyEmployeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignDepartmentResponsibleBody">MembershipCommandServiceAssignDepartmentResponsibleBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassigndepartmentresponsiblebody"></a>
<a id="schema_MembershipCommandServiceAssignDepartmentResponsibleBody"></a>
<a id="tocSmembershipcommandserviceassigndepartmentresponsiblebody"></a>
<a id="tocsmembershipcommandserviceassigndepartmentresponsiblebody"></a>

```json
{
  "employeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignDepartmentResponsibleDeputyBody">MembershipCommandServiceAssignDepartmentResponsibleDeputyBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassigndepartmentresponsibledeputybody"></a>
<a id="schema_MembershipCommandServiceAssignDepartmentResponsibleDeputyBody"></a>
<a id="tocSmembershipcommandserviceassigndepartmentresponsibledeputybody"></a>
<a id="tocsmembershipcommandserviceassigndepartmentresponsibledeputybody"></a>

```json
{
  "deputyEmployeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|deputyEmployeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignOrganizationAdminBody">MembershipCommandServiceAssignOrganizationAdminBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassignorganizationadminbody"></a>
<a id="schema_MembershipCommandServiceAssignOrganizationAdminBody"></a>
<a id="tocSmembershipcommandserviceassignorganizationadminbody"></a>
<a id="tocsmembershipcommandserviceassignorganizationadminbody"></a>

```json
{
  "employeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignOrganizationAdminDeputyBody">MembershipCommandServiceAssignOrganizationAdminDeputyBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassignorganizationadmindeputybody"></a>
<a id="schema_MembershipCommandServiceAssignOrganizationAdminDeputyBody"></a>
<a id="tocSmembershipcommandserviceassignorganizationadmindeputybody"></a>
<a id="tocsmembershipcommandserviceassignorganizationadmindeputybody"></a>

```json
{
  "deputyEmployeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|deputyEmployeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignOrganizationDispatcherBody">MembershipCommandServiceAssignOrganizationDispatcherBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassignorganizationdispatcherbody"></a>
<a id="schema_MembershipCommandServiceAssignOrganizationDispatcherBody"></a>
<a id="tocSmembershipcommandserviceassignorganizationdispatcherbody"></a>
<a id="tocsmembershipcommandserviceassignorganizationdispatcherbody"></a>

```json
{
  "employeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignOrganizationDispatcherDeputyBody">MembershipCommandServiceAssignOrganizationDispatcherDeputyBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassignorganizationdispatcherdeputybody"></a>
<a id="schema_MembershipCommandServiceAssignOrganizationDispatcherDeputyBody"></a>
<a id="tocSmembershipcommandserviceassignorganizationdispatcherdeputybody"></a>
<a id="tocsmembershipcommandserviceassignorganizationdispatcherdeputybody"></a>

```json
{
  "deputyEmployeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|deputyEmployeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignOrganizationHeadBody">MembershipCommandServiceAssignOrganizationHeadBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassignorganizationheadbody"></a>
<a id="schema_MembershipCommandServiceAssignOrganizationHeadBody"></a>
<a id="tocSmembershipcommandserviceassignorganizationheadbody"></a>
<a id="tocsmembershipcommandserviceassignorganizationheadbody"></a>

```json
{
  "employeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceAssignOrganizationHeadDeputyBody">MembershipCommandServiceAssignOrganizationHeadDeputyBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceassignorganizationheaddeputybody"></a>
<a id="schema_MembershipCommandServiceAssignOrganizationHeadDeputyBody"></a>
<a id="tocSmembershipcommandserviceassignorganizationheaddeputybody"></a>
<a id="tocsmembershipcommandserviceassignorganizationheaddeputybody"></a>

```json
{
  "deputyEmployeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|deputyEmployeeId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceScheduleVacationBody">MembershipCommandServiceScheduleVacationBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceschedulevacationbody"></a>
<a id="schema_MembershipCommandServiceScheduleVacationBody"></a>
<a id="tocSmembershipcommandserviceschedulevacationbody"></a>
<a id="tocsmembershipcommandserviceschedulevacationbody"></a>

```json
{
  "startsAt": "2019-08-24T14:15:22Z",
  "endsAt": "2019-08-24T14:15:22Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|startsAt|string(date-time)|true|none|none|
|endsAt|string(date-time)|false|none|none|

<h2 id="tocS_MembershipCommandServiceStartVacationNowBody">MembershipCommandServiceStartVacationNowBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandservicestartvacationnowbody"></a>
<a id="schema_MembershipCommandServiceStartVacationNowBody"></a>
<a id="tocSmembershipcommandservicestartvacationnowbody"></a>
<a id="tocsmembershipcommandservicestartvacationnowbody"></a>

```json
{
  "endsAt": "2019-08-24T14:15:22Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|endsAt|string(date-time)|false|none|none|

<h2 id="tocS_MembershipCommandServiceUpdateEmployeeDepartmentBody">MembershipCommandServiceUpdateEmployeeDepartmentBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceupdateemployeedepartmentbody"></a>
<a id="schema_MembershipCommandServiceUpdateEmployeeDepartmentBody"></a>
<a id="tocSmembershipcommandserviceupdateemployeedepartmentbody"></a>
<a id="tocsmembershipcommandserviceupdateemployeedepartmentbody"></a>

```json
{
  "departmentId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|departmentId|string|true|none|none|

<h2 id="tocS_MembershipCommandServiceUpdateEmployeePositionBody">MembershipCommandServiceUpdateEmployeePositionBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceupdateemployeepositionbody"></a>
<a id="schema_MembershipCommandServiceUpdateEmployeePositionBody"></a>
<a id="tocSmembershipcommandserviceupdateemployeepositionbody"></a>
<a id="tocsmembershipcommandserviceupdateemployeepositionbody"></a>

```json
{
  "position": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|position|string|false|none|none|

<h2 id="tocS_MembershipCommandServiceUpdateVacationEndDateBody">MembershipCommandServiceUpdateVacationEndDateBody</h2>
<!-- backwards compatibility -->
<a id="schemamembershipcommandserviceupdatevacationenddatebody"></a>
<a id="schema_MembershipCommandServiceUpdateVacationEndDateBody"></a>
<a id="tocSmembershipcommandserviceupdatevacationenddatebody"></a>
<a id="tocsmembershipcommandserviceupdatevacationenddatebody"></a>

```json
{
  "endsAt": "2019-08-24T14:15:22Z"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|endsAt|string(date-time)|true|none|none|

<h2 id="tocS_OrgStructureCommandServiceCreateClinicBody">OrgStructureCommandServiceCreateClinicBody</h2>
<!-- backwards compatibility -->
<a id="schemaorgstructurecommandservicecreateclinicbody"></a>
<a id="schema_OrgStructureCommandServiceCreateClinicBody"></a>
<a id="tocSorgstructurecommandservicecreateclinicbody"></a>
<a id="tocsorgstructurecommandservicecreateclinicbody"></a>

```json
{
  "name": "string",
  "physicalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  },
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|physicalAddress|[v1AddressInput](#schemav1addressinput)|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_OrgStructureCommandServiceCreateDepartmentBody">OrgStructureCommandServiceCreateDepartmentBody</h2>
<!-- backwards compatibility -->
<a id="schemaorgstructurecommandservicecreatedepartmentbody"></a>
<a id="schema_OrgStructureCommandServiceCreateDepartmentBody"></a>
<a id="tocSorgstructurecommandservicecreatedepartmentbody"></a>
<a id="tocsorgstructurecommandservicecreatedepartmentbody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_OrgStructureCommandServiceUpdateClinicDetailsBody">OrgStructureCommandServiceUpdateClinicDetailsBody</h2>
<!-- backwards compatibility -->
<a id="schemaorgstructurecommandserviceupdateclinicdetailsbody"></a>
<a id="schema_OrgStructureCommandServiceUpdateClinicDetailsBody"></a>
<a id="tocSorgstructurecommandserviceupdateclinicdetailsbody"></a>
<a id="tocsorgstructurecommandserviceupdateclinicdetailsbody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_OrgStructureCommandServiceUpdateClinicPhysicalAddressBody">OrgStructureCommandServiceUpdateClinicPhysicalAddressBody</h2>
<!-- backwards compatibility -->
<a id="schemaorgstructurecommandserviceupdateclinicphysicaladdressbody"></a>
<a id="schema_OrgStructureCommandServiceUpdateClinicPhysicalAddressBody"></a>
<a id="tocSorgstructurecommandserviceupdateclinicphysicaladdressbody"></a>
<a id="tocsorgstructurecommandserviceupdateclinicphysicaladdressbody"></a>

```json
{
  "physicalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|physicalAddress|[v1AddressInput](#schemav1addressinput)|true|none|none|

<h2 id="tocS_OrgStructureCommandServiceUpdateDepartmentDetailsBody">OrgStructureCommandServiceUpdateDepartmentDetailsBody</h2>
<!-- backwards compatibility -->
<a id="schemaorgstructurecommandserviceupdatedepartmentdetailsbody"></a>
<a id="schema_OrgStructureCommandServiceUpdateDepartmentDetailsBody"></a>
<a id="tocSorgstructurecommandserviceupdatedepartmentdetailsbody"></a>
<a id="tocsorgstructurecommandserviceupdatedepartmentdetailsbody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_OrgStructureCommandServiceUpdateOrganizationDetailsBody">OrgStructureCommandServiceUpdateOrganizationDetailsBody</h2>
<!-- backwards compatibility -->
<a id="schemaorgstructurecommandserviceupdateorganizationdetailsbody"></a>
<a id="schema_OrgStructureCommandServiceUpdateOrganizationDetailsBody"></a>
<a id="tocSorgstructurecommandserviceupdateorganizationdetailsbody"></a>
<a id="tocsorgstructurecommandserviceupdateorganizationdetailsbody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_OrgStructureCommandServiceUpdateOrganizationLegalAddressBody">OrgStructureCommandServiceUpdateOrganizationLegalAddressBody</h2>
<!-- backwards compatibility -->
<a id="schemaorgstructurecommandserviceupdateorganizationlegaladdressbody"></a>
<a id="schema_OrgStructureCommandServiceUpdateOrganizationLegalAddressBody"></a>
<a id="tocSorgstructurecommandserviceupdateorganizationlegaladdressbody"></a>
<a id="tocsorgstructurecommandserviceupdateorganizationlegaladdressbody"></a>

```json
{
  "legalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|legalAddress|[v1AddressInput](#schemav1addressinput)|true|none|none|

<h2 id="tocS_RequestClassifierCommandServiceCreateRequestTypeBody">RequestClassifierCommandServiceCreateRequestTypeBody</h2>
<!-- backwards compatibility -->
<a id="schemarequestclassifiercommandservicecreaterequesttypebody"></a>
<a id="schema_RequestClassifierCommandServiceCreateRequestTypeBody"></a>
<a id="tocSrequestclassifiercommandservicecreaterequesttypebody"></a>
<a id="tocsrequestclassifiercommandservicecreaterequesttypebody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_RequestClassifierCommandServiceUpdateRequestTypeDetailsBody">RequestClassifierCommandServiceUpdateRequestTypeDetailsBody</h2>
<!-- backwards compatibility -->
<a id="schemarequestclassifiercommandserviceupdaterequesttypedetailsbody"></a>
<a id="schema_RequestClassifierCommandServiceUpdateRequestTypeDetailsBody"></a>
<a id="tocSrequestclassifiercommandserviceupdaterequesttypedetailsbody"></a>
<a id="tocsrequestclassifiercommandserviceupdaterequesttypedetailsbody"></a>

```json
{
  "name": "string",
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_ServiceRequestCommandServiceAssignExecutorsBody">ServiceRequestCommandServiceAssignExecutorsBody</h2>
<!-- backwards compatibility -->
<a id="schemaservicerequestcommandserviceassignexecutorsbody"></a>
<a id="schema_ServiceRequestCommandServiceAssignExecutorsBody"></a>
<a id="tocSservicerequestcommandserviceassignexecutorsbody"></a>
<a id="tocsservicerequestcommandserviceassignexecutorsbody"></a>

```json
{
  "executorEmployeeIds": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|executorEmployeeIds|[string]|true|none|none|

<h2 id="tocS_ServiceRequestCommandServiceUpdateServiceRequestDescriptionBody">ServiceRequestCommandServiceUpdateServiceRequestDescriptionBody</h2>
<!-- backwards compatibility -->
<a id="schemaservicerequestcommandserviceupdateservicerequestdescriptionbody"></a>
<a id="schema_ServiceRequestCommandServiceUpdateServiceRequestDescriptionBody"></a>
<a id="tocSservicerequestcommandserviceupdateservicerequestdescriptionbody"></a>
<a id="tocsservicerequestcommandserviceupdateservicerequestdescriptionbody"></a>

```json
{
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|description|string|true|none|none|

<h2 id="tocS_ServiceRequestCommandServiceUpdateServiceRequestStatusBody">ServiceRequestCommandServiceUpdateServiceRequestStatusBody</h2>
<!-- backwards compatibility -->
<a id="schemaservicerequestcommandserviceupdateservicerequeststatusbody"></a>
<a id="schema_ServiceRequestCommandServiceUpdateServiceRequestStatusBody"></a>
<a id="tocSservicerequestcommandserviceupdateservicerequeststatusbody"></a>
<a id="tocsservicerequestcommandserviceupdateservicerequeststatusbody"></a>

```json
{
  "newStatus": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|newStatus|string|true|none|none|

<h2 id="tocS_ValidationFailedDetailsFieldViolation">ValidationFailedDetailsFieldViolation</h2>
<!-- backwards compatibility -->
<a id="schemavalidationfaileddetailsfieldviolation"></a>
<a id="schema_ValidationFailedDetailsFieldViolation"></a>
<a id="tocSvalidationfaileddetailsfieldviolation"></a>
<a id="tocsvalidationfaileddetailsfieldviolation"></a>

```json
{
  "field": "string",
  "rule": "string",
  "message": "string",
  "param": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|field|string|false|none|none|
|rule|string|false|none|none|
|message|string|false|none|none|
|param|string|false|none|none|

<h2 id="tocS_classifierV1Type">classifierV1Type</h2>
<!-- backwards compatibility -->
<a id="schemaclassifierv1type"></a>
<a id="schema_classifierV1Type"></a>
<a id="tocSclassifierv1type"></a>
<a id="tocsclassifierv1type"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "categoryId": "string",
  "name": "string",
  "description": "string",
  "isActive": true,
  "createdAt": "string",
  "updatedAt": "string",
  "isAllowedForPatients": true
}

```

Type mirrors projections.incident_types row.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|categoryId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|isActive|boolean|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|
|isAllowedForPatients|boolean|false|none|none|

<h2 id="tocS_commandAnnouncementV1AnnouncementPriority">commandAnnouncementV1AnnouncementPriority</h2>
<!-- backwards compatibility -->
<a id="schemacommandannouncementv1announcementpriority"></a>
<a id="schema_commandAnnouncementV1AnnouncementPriority"></a>
<a id="tocScommandannouncementv1announcementpriority"></a>
<a id="tocscommandannouncementv1announcementpriority"></a>

```json
"ANNOUNCEMENT_PRIORITY_UNSPECIFIED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|ANNOUNCEMENT_PRIORITY_UNSPECIFIED|
|*anonymous*|ANNOUNCEMENT_PRIORITY_NORMAL|
|*anonymous*|ANNOUNCEMENT_PRIORITY_HIGH|

<h2 id="tocS_commandIncidentV1IncidentPriority">commandIncidentV1IncidentPriority</h2>
<!-- backwards compatibility -->
<a id="schemacommandincidentv1incidentpriority"></a>
<a id="schema_commandIncidentV1IncidentPriority"></a>
<a id="tocScommandincidentv1incidentpriority"></a>
<a id="tocscommandincidentv1incidentpriority"></a>

```json
"INCIDENT_PRIORITY_UNSPECIFIED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|INCIDENT_PRIORITY_UNSPECIFIED|
|*anonymous*|INCIDENT_PRIORITY_LOW|
|*anonymous*|INCIDENT_PRIORITY_NORMAL|
|*anonymous*|INCIDENT_PRIORITY_HIGH|
|*anonymous*|INCIDENT_PRIORITY_CRITICAL|

<h2 id="tocS_commandIncidentV1IncidentStatus">commandIncidentV1IncidentStatus</h2>
<!-- backwards compatibility -->
<a id="schemacommandincidentv1incidentstatus"></a>
<a id="schema_commandIncidentV1IncidentStatus"></a>
<a id="tocScommandincidentv1incidentstatus"></a>
<a id="tocscommandincidentv1incidentstatus"></a>

```json
"INCIDENT_STATUS_UNSPECIFIED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|INCIDENT_STATUS_UNSPECIFIED|
|*anonymous*|INCIDENT_STATUS_PENDING|
|*anonymous*|INCIDENT_STATUS_IN_PROGRESS|
|*anonymous*|INCIDENT_STATUS_DONE|
|*anonymous*|INCIDENT_STATUS_REJECTED|
|*anonymous*|INCIDENT_STATUS_CANCELLED|

<h2 id="tocS_protobufAny">protobufAny</h2>
<!-- backwards compatibility -->
<a id="schemaprotobufany"></a>
<a id="schema_protobufAny"></a>
<a id="tocSprotobufany"></a>
<a id="tocsprotobufany"></a>

```json
{
  "@type": "string",
  "property1": null,
  "property2": null
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|**additionalProperties**|any|false|none|none|
|@type|string|false|none|none|

<h2 id="tocS_queryAnnouncementV1AnnouncementPriority">queryAnnouncementV1AnnouncementPriority</h2>
<!-- backwards compatibility -->
<a id="schemaqueryannouncementv1announcementpriority"></a>
<a id="schema_queryAnnouncementV1AnnouncementPriority"></a>
<a id="tocSqueryannouncementv1announcementpriority"></a>
<a id="tocsqueryannouncementv1announcementpriority"></a>

```json
"ANNOUNCEMENT_PRIORITY_UNSPECIFIED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|ANNOUNCEMENT_PRIORITY_UNSPECIFIED|
|*anonymous*|ANNOUNCEMENT_PRIORITY_NORMAL|
|*anonymous*|ANNOUNCEMENT_PRIORITY_HIGH|

<h2 id="tocS_queryIncidentV1IncidentPriority">queryIncidentV1IncidentPriority</h2>
<!-- backwards compatibility -->
<a id="schemaqueryincidentv1incidentpriority"></a>
<a id="schema_queryIncidentV1IncidentPriority"></a>
<a id="tocSqueryincidentv1incidentpriority"></a>
<a id="tocsqueryincidentv1incidentpriority"></a>

```json
"INCIDENT_PRIORITY_UNSPECIFIED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|INCIDENT_PRIORITY_UNSPECIFIED|
|*anonymous*|INCIDENT_PRIORITY_LOW|
|*anonymous*|INCIDENT_PRIORITY_NORMAL|
|*anonymous*|INCIDENT_PRIORITY_HIGH|
|*anonymous*|INCIDENT_PRIORITY_CRITICAL|

<h2 id="tocS_queryIncidentV1IncidentStatus">queryIncidentV1IncidentStatus</h2>
<!-- backwards compatibility -->
<a id="schemaqueryincidentv1incidentstatus"></a>
<a id="schema_queryIncidentV1IncidentStatus"></a>
<a id="tocSqueryincidentv1incidentstatus"></a>
<a id="tocsqueryincidentv1incidentstatus"></a>

```json
"INCIDENT_STATUS_UNSPECIFIED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|INCIDENT_STATUS_UNSPECIFIED|
|*anonymous*|INCIDENT_STATUS_PENDING|
|*anonymous*|INCIDENT_STATUS_IN_PROGRESS|
|*anonymous*|INCIDENT_STATUS_DONE|
|*anonymous*|INCIDENT_STATUS_REJECTED|
|*anonymous*|INCIDENT_STATUS_CANCELLED|

<h2 id="tocS_queryIncidentV1StatusHistoryEntry">queryIncidentV1StatusHistoryEntry</h2>
<!-- backwards compatibility -->
<a id="schemaqueryincidentv1statushistoryentry"></a>
<a id="schema_queryIncidentV1StatusHistoryEntry"></a>
<a id="tocSqueryincidentv1statushistoryentry"></a>
<a id="tocsqueryincidentv1statushistoryentry"></a>

```json
{
  "id": "string",
  "oldStatus": "INCIDENT_STATUS_UNSPECIFIED",
  "newStatus": "INCIDENT_STATUS_UNSPECIFIED",
  "actor": {
    "employeeId": "string",
    "displayName": "string"
  },
  "changedAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|oldStatus|[queryIncidentV1IncidentStatus](#schemaqueryincidentv1incidentstatus)|false|none|none|
|newStatus|[queryIncidentV1IncidentStatus](#schemaqueryincidentv1incidentstatus)|false|none|none|
|actor|[v1ActorView](#schemav1actorview)|false|none|none|
|changedAt|string|false|none|none|

<h2 id="tocS_queryRequestV1StatusHistoryEntry">queryRequestV1StatusHistoryEntry</h2>
<!-- backwards compatibility -->
<a id="schemaqueryrequestv1statushistoryentry"></a>
<a id="schema_queryRequestV1StatusHistoryEntry"></a>
<a id="tocSqueryrequestv1statushistoryentry"></a>
<a id="tocsqueryrequestv1statushistoryentry"></a>

```json
{
  "id": "string",
  "oldStatus": "string",
  "newStatus": "string",
  "actorId": "string",
  "actorName": "string",
  "changedAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|oldStatus|string|false|none|none|
|newStatus|string|false|none|none|
|actorId|string|false|none|none|
|actorName|string|false|none|none|
|changedAt|string|false|none|none|

<h2 id="tocS_rpcStatus">rpcStatus</h2>
<!-- backwards compatibility -->
<a id="schemarpcstatus"></a>
<a id="schema_rpcStatus"></a>
<a id="tocSrpcstatus"></a>
<a id="tocsrpcstatus"></a>

```json
{
  "code": 0,
  "message": "string",
  "details": [
    {
      "@type": "string",
      "property1": null,
      "property2": null
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|integer(int32)|false|none|none|
|message|string|false|none|none|
|details|[[protobufAny](#schemaprotobufany)]|false|none|none|

<h2 id="tocS_v1ActorView">v1ActorView</h2>
<!-- backwards compatibility -->
<a id="schemav1actorview"></a>
<a id="schema_v1ActorView"></a>
<a id="tocSv1actorview"></a>
<a id="tocsv1actorview"></a>

```json
{
  "employeeId": "string",
  "displayName": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|false|none|none|
|displayName|string|false|none|none|

<h2 id="tocS_v1Address">v1Address</h2>
<!-- backwards compatibility -->
<a id="schemav1address"></a>
<a id="schema_v1Address"></a>
<a id="tocSv1address"></a>
<a id="tocsv1address"></a>

```json
{
  "text": "string",
  "point": {
    "longitude": 0.1,
    "latitude": 0.1
  }
}

```

Address is the read-side view of a stored postal address; Point is
optional because the command-side allows text-only addresses.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|text|string|false|none|none|
|point|[v1Point](#schemav1point)|false|none|Point is the coordinate pair attached to an Address. Absent when the<br>projection row has neither longitude nor latitude.|

<h2 id="tocS_v1AddressInput">v1AddressInput</h2>
<!-- backwards compatibility -->
<a id="schemav1addressinput"></a>
<a id="schema_v1AddressInput"></a>
<a id="tocSv1addressinput"></a>
<a id="tocsv1addressinput"></a>

```json
{
  "text": "string",
  "point": {
    "longitude": 0.1,
    "latitude": 0.1
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|text|string|true|none|none|
|point|[v1PointInput](#schemav1pointinput)|false|none|none|

<h2 id="tocS_v1AllowIncidentTypeForPatientsResponse">v1AllowIncidentTypeForPatientsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1allowincidenttypeforpatientsresponse"></a>
<a id="schema_v1AllowIncidentTypeForPatientsResponse"></a>
<a id="tocSv1allowincidenttypeforpatientsresponse"></a>
<a id="tocsv1allowincidenttypeforpatientsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AnnouncementView">v1AnnouncementView</h2>
<!-- backwards compatibility -->
<a id="schemav1announcementview"></a>
<a id="schema_v1AnnouncementView"></a>
<a id="tocSv1announcementview"></a>
<a id="tocsv1announcementview"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "clinicId": "string",
  "departmentId": "string",
  "authorId": "string",
  "title": "string",
  "content": "string",
  "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
  "isArchived": true,
  "startsAt": "string",
  "endsAt": "string",
  "createdAt": "string",
  "updatedAt": "string",
  "viewCount": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|clinicId|string|false|none|none|
|departmentId|string|false|none|none|
|authorId|string|false|none|none|
|title|string|false|none|none|
|content|string|false|none|none|
|priority|[queryAnnouncementV1AnnouncementPriority](#schemaqueryannouncementv1announcementpriority)|false|none|none|
|isArchived|boolean|false|none|none|
|startsAt|string|false|none|none|
|endsAt|string|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|
|viewCount|string(int64)|false|none|none|

<h2 id="tocS_v1ArchiveAnnouncementResponse">v1ArchiveAnnouncementResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1archiveannouncementresponse"></a>
<a id="schema_v1ArchiveAnnouncementResponse"></a>
<a id="tocSv1archiveannouncementresponse"></a>
<a id="tocsv1archiveannouncementresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignClinicHeadDeputyResponse">v1AssignClinicHeadDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignclinicheaddeputyresponse"></a>
<a id="schema_v1AssignClinicHeadDeputyResponse"></a>
<a id="tocSv1assignclinicheaddeputyresponse"></a>
<a id="tocsv1assignclinicheaddeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignClinicHeadResponse">v1AssignClinicHeadResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignclinicheadresponse"></a>
<a id="schema_v1AssignClinicHeadResponse"></a>
<a id="tocSv1assignclinicheadresponse"></a>
<a id="tocsv1assignclinicheadresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignDepartmentResponsibleDeputyResponse">v1AssignDepartmentResponsibleDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assigndepartmentresponsibledeputyresponse"></a>
<a id="schema_v1AssignDepartmentResponsibleDeputyResponse"></a>
<a id="tocSv1assigndepartmentresponsibledeputyresponse"></a>
<a id="tocsv1assigndepartmentresponsibledeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignDepartmentResponsibleResponse">v1AssignDepartmentResponsibleResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assigndepartmentresponsibleresponse"></a>
<a id="schema_v1AssignDepartmentResponsibleResponse"></a>
<a id="tocSv1assigndepartmentresponsibleresponse"></a>
<a id="tocsv1assigndepartmentresponsibleresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignExecutorsResponse">v1AssignExecutorsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignexecutorsresponse"></a>
<a id="schema_v1AssignExecutorsResponse"></a>
<a id="tocSv1assignexecutorsresponse"></a>
<a id="tocsv1assignexecutorsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignOrganizationAdminDeputyResponse">v1AssignOrganizationAdminDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignorganizationadmindeputyresponse"></a>
<a id="schema_v1AssignOrganizationAdminDeputyResponse"></a>
<a id="tocSv1assignorganizationadmindeputyresponse"></a>
<a id="tocsv1assignorganizationadmindeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignOrganizationAdminResponse">v1AssignOrganizationAdminResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignorganizationadminresponse"></a>
<a id="schema_v1AssignOrganizationAdminResponse"></a>
<a id="tocSv1assignorganizationadminresponse"></a>
<a id="tocsv1assignorganizationadminresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignOrganizationDispatcherDeputyResponse">v1AssignOrganizationDispatcherDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignorganizationdispatcherdeputyresponse"></a>
<a id="schema_v1AssignOrganizationDispatcherDeputyResponse"></a>
<a id="tocSv1assignorganizationdispatcherdeputyresponse"></a>
<a id="tocsv1assignorganizationdispatcherdeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignOrganizationDispatcherResponse">v1AssignOrganizationDispatcherResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignorganizationdispatcherresponse"></a>
<a id="schema_v1AssignOrganizationDispatcherResponse"></a>
<a id="tocSv1assignorganizationdispatcherresponse"></a>
<a id="tocsv1assignorganizationdispatcherresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignOrganizationHeadDeputyResponse">v1AssignOrganizationHeadDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignorganizationheaddeputyresponse"></a>
<a id="schema_v1AssignOrganizationHeadDeputyResponse"></a>
<a id="tocSv1assignorganizationheaddeputyresponse"></a>
<a id="tocsv1assignorganizationheaddeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1AssignOrganizationHeadResponse">v1AssignOrganizationHeadResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1assignorganizationheadresponse"></a>
<a id="schema_v1AssignOrganizationHeadResponse"></a>
<a id="tocSv1assignorganizationheadresponse"></a>
<a id="tocsv1assignorganizationheadresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1BufferEntryView">v1BufferEntryView</h2>
<!-- backwards compatibility -->
<a id="schemav1bufferentryview"></a>
<a id="schema_v1BufferEntryView"></a>
<a id="tocSv1bufferentryview"></a>
<a id="tocsv1bufferentryview"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "patientZitadelUserId": "string",
  "categoryId": "string",
  "typeId": "string",
  "description": "string",
  "occurredAt": "string",
  "status": "BUFFER_STATUS_UNSPECIFIED",
  "publishedIncidentId": "string",
  "createdAt": "string",
  "updatedAt": "string",
  "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|patientZitadelUserId|string|false|none|none|
|categoryId|string|false|none|none|
|typeId|string|false|none|none|
|description|string|false|none|none|
|occurredAt|string|false|none|none|
|status|[v1BufferStatus](#schemav1bufferstatus)|false|none|none|
|publishedIncidentId|string|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|
|patientStatus|[v1PatientStatus](#schemav1patientstatus)|false|none|PatientStatus is the simplified four-value status surfaced to patients.<br><br> - PATIENT_STATUS_PENDING: buffer pending<br> - PATIENT_STATUS_ACCEPTED: dispatcher accepted; incident pending/in_progress<br> - PATIENT_STATUS_CLOSED: done / rejected / buffer rejected<br> - PATIENT_STATUS_CANCELLED: patient cancelled|

<h2 id="tocS_v1BufferStatus">v1BufferStatus</h2>
<!-- backwards compatibility -->
<a id="schemav1bufferstatus"></a>
<a id="schema_v1BufferStatus"></a>
<a id="tocSv1bufferstatus"></a>
<a id="tocsv1bufferstatus"></a>

```json
"BUFFER_STATUS_UNSPECIFIED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|BUFFER_STATUS_UNSPECIFIED|
|*anonymous*|BUFFER_STATUS_PENDING|
|*anonymous*|BUFFER_STATUS_PUBLISHED|
|*anonymous*|BUFFER_STATUS_REJECTED|
|*anonymous*|BUFFER_STATUS_CANCELLED|

<h2 id="tocS_v1CancelIncidentResponse">v1CancelIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1cancelincidentresponse"></a>
<a id="schema_v1CancelIncidentResponse"></a>
<a id="tocSv1cancelincidentresponse"></a>
<a id="tocsv1cancelincidentresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1CancelPatientIncidentResponse">v1CancelPatientIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1cancelpatientincidentresponse"></a>
<a id="schema_v1CancelPatientIncidentResponse"></a>
<a id="tocSv1cancelpatientincidentresponse"></a>
<a id="tocsv1cancelpatientincidentresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1CancelScheduledVacationResponse">v1CancelScheduledVacationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1cancelscheduledvacationresponse"></a>
<a id="schema_v1CancelScheduledVacationResponse"></a>
<a id="tocSv1cancelscheduledvacationresponse"></a>
<a id="tocsv1cancelscheduledvacationresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1Category">v1Category</h2>
<!-- backwards compatibility -->
<a id="schemav1category"></a>
<a id="schema_v1Category"></a>
<a id="tocSv1category"></a>
<a id="tocsv1category"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "parentCategoryId": "string",
  "name": "string",
  "description": "string",
  "isActive": true,
  "createdAt": "string",
  "updatedAt": "string"
}

```

Category mirrors projections.incident_categories row.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|parentCategoryId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|isActive|boolean|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|

<h2 id="tocS_v1CategoryCount">v1CategoryCount</h2>
<!-- backwards compatibility -->
<a id="schemav1categorycount"></a>
<a id="schema_v1CategoryCount"></a>
<a id="tocSv1categorycount"></a>
<a id="tocsv1categorycount"></a>

```json
{
  "categoryId": "string",
  "categoryName": "string",
  "count": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|categoryId|string|false|none|none|
|categoryName|string|false|none|none|
|count|string(int64)|false|none|none|

<h2 id="tocS_v1Clinic">v1Clinic</h2>
<!-- backwards compatibility -->
<a id="schemav1clinic"></a>
<a id="schema_v1Clinic"></a>
<a id="tocSv1clinic"></a>
<a id="tocsv1clinic"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "name": "string",
  "description": "string",
  "physicalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  },
  "createdAt": "string",
  "updatedAt": "string"
}

```

Clinic mirrors the projections.clinics row returned by GetClinic.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|physicalAddress|[v1Address](#schemav1address)|false|none|Address is the read-side view of a stored postal address; Point is<br>optional because the command-side allows text-only addresses.|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|

<h2 id="tocS_v1ClinicListItem">v1ClinicListItem</h2>
<!-- backwards compatibility -->
<a id="schemav1cliniclistitem"></a>
<a id="schema_v1ClinicListItem"></a>
<a id="tocSv1cliniclistitem"></a>
<a id="tocsv1cliniclistitem"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "name": "string"
}

```

ClinicListItem is the minimal shape returned by list endpoints.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|name|string|false|none|none|

<h2 id="tocS_v1ClinicStats">v1ClinicStats</h2>
<!-- backwards compatibility -->
<a id="schemav1clinicstats"></a>
<a id="schema_v1ClinicStats"></a>
<a id="tocSv1clinicstats"></a>
<a id="tocsv1clinicstats"></a>

```json
{
  "clinicId": "string",
  "organizationId": "string",
  "employeesTotal": "string",
  "departmentsTotal": "string",
  "employeesOnVacation": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|clinicId|string|false|none|none|
|organizationId|string|false|none|none|
|employeesTotal|string(int64)|false|none|none|
|departmentsTotal|string(int64)|false|none|none|
|employeesOnVacation|string(int64)|false|none|none|

<h2 id="tocS_v1CountClinicsByOrganizationResponse">v1CountClinicsByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1countclinicsbyorganizationresponse"></a>
<a id="schema_v1CountClinicsByOrganizationResponse"></a>
<a id="tocSv1countclinicsbyorganizationresponse"></a>
<a id="tocsv1countclinicsbyorganizationresponse"></a>

```json
{
  "total": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|

<h2 id="tocS_v1CountDepartmentsByClinicResponse">v1CountDepartmentsByClinicResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1countdepartmentsbyclinicresponse"></a>
<a id="schema_v1CountDepartmentsByClinicResponse"></a>
<a id="tocSv1countdepartmentsbyclinicresponse"></a>
<a id="tocsv1countdepartmentsbyclinicresponse"></a>

```json
{
  "total": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|

<h2 id="tocS_v1CountEmployeesByClinicResponse">v1CountEmployeesByClinicResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1countemployeesbyclinicresponse"></a>
<a id="schema_v1CountEmployeesByClinicResponse"></a>
<a id="tocSv1countemployeesbyclinicresponse"></a>
<a id="tocsv1countemployeesbyclinicresponse"></a>

```json
{
  "total": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|

<h2 id="tocS_v1CountEmployeesByDepartmentResponse">v1CountEmployeesByDepartmentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1countemployeesbydepartmentresponse"></a>
<a id="schema_v1CountEmployeesByDepartmentResponse"></a>
<a id="tocSv1countemployeesbydepartmentresponse"></a>
<a id="tocsv1countemployeesbydepartmentresponse"></a>

```json
{
  "total": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|

<h2 id="tocS_v1CountEmployeesByOrganizationResponse">v1CountEmployeesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1countemployeesbyorganizationresponse"></a>
<a id="schema_v1CountEmployeesByOrganizationResponse"></a>
<a id="tocSv1countemployeesbyorganizationresponse"></a>
<a id="tocsv1countemployeesbyorganizationresponse"></a>

```json
{
  "total": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|

<h2 id="tocS_v1CountOrganizationsResponse">v1CountOrganizationsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1countorganizationsresponse"></a>
<a id="schema_v1CountOrganizationsResponse"></a>
<a id="tocSv1countorganizationsresponse"></a>
<a id="tocsv1countorganizationsresponse"></a>

```json
{
  "total": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|

<h2 id="tocS_v1CountVacationsByEmployeeResponse">v1CountVacationsByEmployeeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1countvacationsbyemployeeresponse"></a>
<a id="schema_v1CountVacationsByEmployeeResponse"></a>
<a id="tocSv1countvacationsbyemployeeresponse"></a>
<a id="tocsv1countvacationsbyemployeeresponse"></a>

```json
{
  "total": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|

<h2 id="tocS_v1CreateAnnouncementRequest">v1CreateAnnouncementRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1createannouncementrequest"></a>
<a id="schema_v1CreateAnnouncementRequest"></a>
<a id="tocSv1createannouncementrequest"></a>
<a id="tocsv1createannouncementrequest"></a>

```json
{
  "organizationId": "string",
  "clinicId": "string",
  "departmentId": "string",
  "title": "string",
  "content": "string",
  "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
  "startsAt": "string",
  "endsAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|organizationId|string|true|none|none|
|clinicId|string|false|none|none|
|departmentId|string|false|none|none|
|title|string|true|none|none|
|content|string|true|none|none|
|priority|[commandAnnouncementV1AnnouncementPriority](#schemacommandannouncementv1announcementpriority)|true|none|none|
|startsAt|string|false|none|none|
|endsAt|string|false|none|none|

<h2 id="tocS_v1CreateAnnouncementResponse">v1CreateAnnouncementResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createannouncementresponse"></a>
<a id="schema_v1CreateAnnouncementResponse"></a>
<a id="tocSv1createannouncementresponse"></a>
<a id="tocsv1createannouncementresponse"></a>

```json
{
  "id": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|

<h2 id="tocS_v1CreateClinicResponse">v1CreateClinicResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createclinicresponse"></a>
<a id="schema_v1CreateClinicResponse"></a>
<a id="tocSv1createclinicresponse"></a>
<a id="tocsv1createclinicresponse"></a>

```json
{
  "clinicId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|clinicId|string|false|none|none|

<h2 id="tocS_v1CreateDepartmentResponse">v1CreateDepartmentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createdepartmentresponse"></a>
<a id="schema_v1CreateDepartmentResponse"></a>
<a id="tocSv1createdepartmentresponse"></a>
<a id="tocsv1createdepartmentresponse"></a>

```json
{
  "departmentId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|departmentId|string|false|none|none|

<h2 id="tocS_v1CreateIncidentCategoryResponse">v1CreateIncidentCategoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createincidentcategoryresponse"></a>
<a id="schema_v1CreateIncidentCategoryResponse"></a>
<a id="tocSv1createincidentcategoryresponse"></a>
<a id="tocsv1createincidentcategoryresponse"></a>

```json
{
  "categoryId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|categoryId|string|false|none|none|

<h2 id="tocS_v1CreateIncidentRequest">v1CreateIncidentRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1createincidentrequest"></a>
<a id="schema_v1CreateIncidentRequest"></a>
<a id="tocSv1createincidentrequest"></a>
<a id="tocsv1createincidentrequest"></a>

```json
{
  "departmentId": "string",
  "categoryId": "string",
  "typeId": "string",
  "description": "string",
  "occurredAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|departmentId|string|true|none|none|
|categoryId|string|true|none|none|
|typeId|string|true|none|none|
|description|string|false|none|none|
|occurredAt|string|true|none|none|

<h2 id="tocS_v1CreateIncidentResponse">v1CreateIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createincidentresponse"></a>
<a id="schema_v1CreateIncidentResponse"></a>
<a id="tocSv1createincidentresponse"></a>
<a id="tocsv1createincidentresponse"></a>

```json
{
  "incidentId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|incidentId|string|false|none|none|

<h2 id="tocS_v1CreateIncidentTypeResponse">v1CreateIncidentTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createincidenttyperesponse"></a>
<a id="schema_v1CreateIncidentTypeResponse"></a>
<a id="tocSv1createincidenttyperesponse"></a>
<a id="tocsv1createincidenttyperesponse"></a>

```json
{
  "typeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|typeId|string|false|none|none|

<h2 id="tocS_v1CreateOrganizationRequest">v1CreateOrganizationRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1createorganizationrequest"></a>
<a id="schema_v1CreateOrganizationRequest"></a>
<a id="tocSv1createorganizationrequest"></a>
<a id="tocsv1createorganizationrequest"></a>

```json
{
  "name": "string",
  "legalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  },
  "description": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|name|string|true|none|none|
|legalAddress|[v1AddressInput](#schemav1addressinput)|true|none|none|
|description|string|false|none|none|

<h2 id="tocS_v1CreateOrganizationResponse">v1CreateOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createorganizationresponse"></a>
<a id="schema_v1CreateOrganizationResponse"></a>
<a id="tocSv1createorganizationresponse"></a>
<a id="tocsv1createorganizationresponse"></a>

```json
{
  "organizationId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|organizationId|string|false|none|none|

<h2 id="tocS_v1CreateRequestTypeResponse">v1CreateRequestTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createrequesttyperesponse"></a>
<a id="schema_v1CreateRequestTypeResponse"></a>
<a id="tocSv1createrequesttyperesponse"></a>
<a id="tocsv1createrequesttyperesponse"></a>

```json
{
  "typeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|typeId|string|false|none|none|

<h2 id="tocS_v1CreateServiceRequestRequest">v1CreateServiceRequestRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1createservicerequestrequest"></a>
<a id="schema_v1CreateServiceRequestRequest"></a>
<a id="tocSv1createservicerequestrequest"></a>
<a id="tocsv1createservicerequestrequest"></a>

```json
{
  "departmentId": "string",
  "typeId": "string",
  "incidentId": "string",
  "description": "string",
  "executorEmployeeIds": [
    "string"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|departmentId|string|true|none|none|
|typeId|string|true|none|none|
|incidentId|string|false|none|none|
|description|string|true|none|none|
|executorEmployeeIds|[string]|true|none|none|

<h2 id="tocS_v1CreateServiceRequestResponse">v1CreateServiceRequestResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1createservicerequestresponse"></a>
<a id="schema_v1CreateServiceRequestResponse"></a>
<a id="tocSv1createservicerequestresponse"></a>
<a id="tocsv1createservicerequestresponse"></a>

```json
{
  "serviceRequestId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|serviceRequestId|string|false|none|none|

<h2 id="tocS_v1DeactivateIncidentCategoryResponse">v1DeactivateIncidentCategoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1deactivateincidentcategoryresponse"></a>
<a id="schema_v1DeactivateIncidentCategoryResponse"></a>
<a id="tocSv1deactivateincidentcategoryresponse"></a>
<a id="tocsv1deactivateincidentcategoryresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1DeactivateIncidentTypeResponse">v1DeactivateIncidentTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1deactivateincidenttyperesponse"></a>
<a id="schema_v1DeactivateIncidentTypeResponse"></a>
<a id="tocSv1deactivateincidenttyperesponse"></a>
<a id="tocsv1deactivateincidenttyperesponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1DeactivateRequestTypeResponse">v1DeactivateRequestTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1deactivaterequesttyperesponse"></a>
<a id="schema_v1DeactivateRequestTypeResponse"></a>
<a id="tocSv1deactivaterequesttyperesponse"></a>
<a id="tocsv1deactivaterequesttyperesponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1DeleteIncidentCategoryResponse">v1DeleteIncidentCategoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1deleteincidentcategoryresponse"></a>
<a id="schema_v1DeleteIncidentCategoryResponse"></a>
<a id="tocSv1deleteincidentcategoryresponse"></a>
<a id="tocsv1deleteincidentcategoryresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1DeleteIncidentTypeResponse">v1DeleteIncidentTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1deleteincidenttyperesponse"></a>
<a id="schema_v1DeleteIncidentTypeResponse"></a>
<a id="tocSv1deleteincidenttyperesponse"></a>
<a id="tocsv1deleteincidenttyperesponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1DeleteRequestTypeResponse">v1DeleteRequestTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1deleterequesttyperesponse"></a>
<a id="schema_v1DeleteRequestTypeResponse"></a>
<a id="tocSv1deleterequesttyperesponse"></a>
<a id="tocsv1deleterequesttyperesponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1Department">v1Department</h2>
<!-- backwards compatibility -->
<a id="schemav1department"></a>
<a id="schema_v1Department"></a>
<a id="tocSv1department"></a>
<a id="tocsv1department"></a>

```json
{
  "id": "string",
  "clinicId": "string",
  "name": "string",
  "description": "string",
  "createdAt": "string",
  "updatedAt": "string"
}

```

Department mirrors the projections.departments row returned by
GetDepartment.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|clinicId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|

<h2 id="tocS_v1DepartmentCount">v1DepartmentCount</h2>
<!-- backwards compatibility -->
<a id="schemav1departmentcount"></a>
<a id="schema_v1DepartmentCount"></a>
<a id="tocSv1departmentcount"></a>
<a id="tocsv1departmentcount"></a>

```json
{
  "departmentId": "string",
  "departmentName": "string",
  "count": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|departmentId|string|false|none|none|
|departmentName|string|false|none|none|
|count|string(int64)|false|none|none|

<h2 id="tocS_v1DepartmentListItem">v1DepartmentListItem</h2>
<!-- backwards compatibility -->
<a id="schemav1departmentlistitem"></a>
<a id="schema_v1DepartmentListItem"></a>
<a id="tocSv1departmentlistitem"></a>
<a id="tocsv1departmentlistitem"></a>

```json
{
  "id": "string",
  "clinicId": "string",
  "name": "string"
}

```

DepartmentListItem is the minimal shape returned by list endpoints.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|clinicId|string|false|none|none|
|name|string|false|none|none|

<h2 id="tocS_v1DepartmentStats">v1DepartmentStats</h2>
<!-- backwards compatibility -->
<a id="schemav1departmentstats"></a>
<a id="schema_v1DepartmentStats"></a>
<a id="tocSv1departmentstats"></a>
<a id="tocsv1departmentstats"></a>

```json
{
  "departmentId": "string",
  "clinicId": "string",
  "organizationId": "string",
  "employeesTotal": "string",
  "employeesOnVacation": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|departmentId|string|false|none|none|
|clinicId|string|false|none|none|
|organizationId|string|false|none|none|
|employeesTotal|string(int64)|false|none|none|
|employeesOnVacation|string(int64)|false|none|none|

<h2 id="tocS_v1DisallowIncidentTypeForPatientsResponse">v1DisallowIncidentTypeForPatientsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1disallowincidenttypeforpatientsresponse"></a>
<a id="schema_v1DisallowIncidentTypeForPatientsResponse"></a>
<a id="tocSv1disallowincidenttypeforpatientsresponse"></a>
<a id="tocsv1disallowincidenttypeforpatientsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1EmployeeCardView">v1EmployeeCardView</h2>
<!-- backwards compatibility -->
<a id="schemav1employeecardview"></a>
<a id="schema_v1EmployeeCardView"></a>
<a id="tocSv1employeecardview"></a>
<a id="tocsv1employeecardview"></a>

```json
{
  "employeeId": "string",
  "zitadelUserId": "string",
  "firstName": "string",
  "lastName": "string",
  "displayName": "string",
  "email": "string",
  "organizationId": "string",
  "organizationName": "string",
  "clinicId": "string",
  "clinicName": "string",
  "departmentId": "string",
  "departmentName": "string",
  "position": "string",
  "terminatedAt": "string",
  "currentVacationEndsAt": "string",
  "nextVacationStartsAt": "string"
}

```

EmployeeCardView is the denormalised card projection returned by
both Get and List endpoints. Fields mirror projections.employee_cards
columns; timestamps are RFC3339 strings. Optional fields stay unset
when the backing column is NULL.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|false|none|none|
|zitadelUserId|string|false|none|none|
|firstName|string|false|none|none|
|lastName|string|false|none|none|
|displayName|string|false|none|none|
|email|string|false|none|none|
|organizationId|string|false|none|none|
|organizationName|string|false|none|none|
|clinicId|string|false|none|none|
|clinicName|string|false|none|none|
|departmentId|string|false|none|none|
|departmentName|string|false|none|none|
|position|string|false|none|none|
|terminatedAt|string|false|none|none|
|currentVacationEndsAt|string|false|none|none|
|nextVacationStartsAt|string|false|none|none|

<h2 id="tocS_v1ErrorResponse">v1ErrorResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1errorresponse"></a>
<a id="schema_v1ErrorResponse"></a>
<a id="tocSv1errorresponse"></a>
<a id="tocsv1errorresponse"></a>

```json
{
  "code": "string",
  "message": "string"
}

```

ErrorResponse

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|string|false|none|none|
|message|string|false|none|none|

<h2 id="tocS_v1Executor">v1Executor</h2>
<!-- backwards compatibility -->
<a id="schemav1executor"></a>
<a id="schema_v1Executor"></a>
<a id="tocSv1executor"></a>
<a id="tocsv1executor"></a>

```json
{
  "employeeId": "string",
  "assignedAt": "string",
  "assignedById": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|false|none|none|
|assignedAt|string|false|none|none|
|assignedById|string|false|none|none|

<h2 id="tocS_v1ExecutorHistoryEntry">v1ExecutorHistoryEntry</h2>
<!-- backwards compatibility -->
<a id="schemav1executorhistoryentry"></a>
<a id="schema_v1ExecutorHistoryEntry"></a>
<a id="tocSv1executorhistoryentry"></a>
<a id="tocsv1executorhistoryentry"></a>

```json
{
  "id": "string",
  "action": "string",
  "employeeId": "string",
  "employeeName": "string",
  "actorId": "string",
  "actorName": "string",
  "changedAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|action|string|false|none|none|
|employeeId|string|false|none|none|
|employeeName|string|false|none|none|
|actorId|string|false|none|none|
|actorName|string|false|none|none|
|changedAt|string|false|none|none|

<h2 id="tocS_v1ForceEndVacationResponse">v1ForceEndVacationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1forceendvacationresponse"></a>
<a id="schema_v1ForceEndVacationResponse"></a>
<a id="tocSv1forceendvacationresponse"></a>
<a id="tocsv1forceendvacationresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1GetAnnouncementResponse">v1GetAnnouncementResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getannouncementresponse"></a>
<a id="schema_v1GetAnnouncementResponse"></a>
<a id="tocSv1getannouncementresponse"></a>
<a id="tocsv1getannouncementresponse"></a>

```json
{
  "announcement": {
    "id": "string",
    "organizationId": "string",
    "clinicId": "string",
    "departmentId": "string",
    "authorId": "string",
    "title": "string",
    "content": "string",
    "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
    "isArchived": true,
    "startsAt": "string",
    "endsAt": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "viewCount": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|announcement|[v1AnnouncementView](#schemav1announcementview)|false|none|none|

<h2 id="tocS_v1GetBufferEntryResponse">v1GetBufferEntryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getbufferentryresponse"></a>
<a id="schema_v1GetBufferEntryResponse"></a>
<a id="tocSv1getbufferentryresponse"></a>
<a id="tocsv1getbufferentryresponse"></a>

```json
{
  "entry": {
    "id": "string",
    "organizationId": "string",
    "patientZitadelUserId": "string",
    "categoryId": "string",
    "typeId": "string",
    "description": "string",
    "occurredAt": "string",
    "status": "BUFFER_STATUS_UNSPECIFIED",
    "publishedIncidentId": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|entry|[v1BufferEntryView](#schemav1bufferentryview)|false|none|none|

<h2 id="tocS_v1GetCategoryResponse">v1GetCategoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getcategoryresponse"></a>
<a id="schema_v1GetCategoryResponse"></a>
<a id="tocSv1getcategoryresponse"></a>
<a id="tocsv1getcategoryresponse"></a>

```json
{
  "category": {
    "id": "string",
    "organizationId": "string",
    "parentCategoryId": "string",
    "name": "string",
    "description": "string",
    "isActive": true,
    "createdAt": "string",
    "updatedAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|category|[v1Category](#schemav1category)|false|none|Category mirrors projections.incident_categories row.|

<h2 id="tocS_v1GetClinicHeadResponse">v1GetClinicHeadResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getclinicheadresponse"></a>
<a id="schema_v1GetClinicHeadResponse"></a>
<a id="tocSv1getclinicheadresponse"></a>
<a id="tocsv1getclinicheadresponse"></a>

```json
{
  "assignment": {
    "holder": {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    },
    "deputy": {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|assignment|[v1RoleAssignment](#schemav1roleassignment)|false|none|RoleAssignment is a role row enriched with the denormalised card for<br>the holder and, when present, the deputy. Read-model callers use<br>this so they do not need to follow role lookups with N+1 GetEmployee<br>calls to render a name or email.|

<h2 id="tocS_v1GetClinicResponse">v1GetClinicResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getclinicresponse"></a>
<a id="schema_v1GetClinicResponse"></a>
<a id="tocSv1getclinicresponse"></a>
<a id="tocsv1getclinicresponse"></a>

```json
{
  "clinic": {
    "id": "string",
    "organizationId": "string",
    "name": "string",
    "description": "string",
    "physicalAddress": {
      "text": "string",
      "point": {
        "longitude": 0.1,
        "latitude": 0.1
      }
    },
    "createdAt": "string",
    "updatedAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|clinic|[v1Clinic](#schemav1clinic)|false|none|Clinic mirrors the projections.clinics row returned by GetClinic.|

<h2 id="tocS_v1GetClinicStatsResponse">v1GetClinicStatsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getclinicstatsresponse"></a>
<a id="schema_v1GetClinicStatsResponse"></a>
<a id="tocSv1getclinicstatsresponse"></a>
<a id="tocsv1getclinicstatsresponse"></a>

```json
{
  "stats": {
    "clinicId": "string",
    "organizationId": "string",
    "employeesTotal": "string",
    "departmentsTotal": "string",
    "employeesOnVacation": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|stats|[v1ClinicStats](#schemav1clinicstats)|false|none|none|

<h2 id="tocS_v1GetDepartmentResponse">v1GetDepartmentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getdepartmentresponse"></a>
<a id="schema_v1GetDepartmentResponse"></a>
<a id="tocSv1getdepartmentresponse"></a>
<a id="tocsv1getdepartmentresponse"></a>

```json
{
  "department": {
    "id": "string",
    "clinicId": "string",
    "name": "string",
    "description": "string",
    "createdAt": "string",
    "updatedAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|department|[v1Department](#schemav1department)|false|none|Department mirrors the projections.departments row returned by<br>GetDepartment.|

<h2 id="tocS_v1GetDepartmentResponsibleResponse">v1GetDepartmentResponsibleResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getdepartmentresponsibleresponse"></a>
<a id="schema_v1GetDepartmentResponsibleResponse"></a>
<a id="tocSv1getdepartmentresponsibleresponse"></a>
<a id="tocsv1getdepartmentresponsibleresponse"></a>

```json
{
  "assignment": {
    "holder": {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    },
    "deputy": {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|assignment|[v1RoleAssignment](#schemav1roleassignment)|false|none|RoleAssignment is a role row enriched with the denormalised card for<br>the holder and, when present, the deputy. Read-model callers use<br>this so they do not need to follow role lookups with N+1 GetEmployee<br>calls to render a name or email.|

<h2 id="tocS_v1GetDepartmentStatsResponse">v1GetDepartmentStatsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getdepartmentstatsresponse"></a>
<a id="schema_v1GetDepartmentStatsResponse"></a>
<a id="tocSv1getdepartmentstatsresponse"></a>
<a id="tocsv1getdepartmentstatsresponse"></a>

```json
{
  "stats": {
    "departmentId": "string",
    "clinicId": "string",
    "organizationId": "string",
    "employeesTotal": "string",
    "employeesOnVacation": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|stats|[v1DepartmentStats](#schemav1departmentstats)|false|none|none|

<h2 id="tocS_v1GetEmployeeResponse">v1GetEmployeeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getemployeeresponse"></a>
<a id="schema_v1GetEmployeeResponse"></a>
<a id="tocSv1getemployeeresponse"></a>
<a id="tocsv1getemployeeresponse"></a>

```json
{
  "employee": {
    "employeeId": "string",
    "zitadelUserId": "string",
    "firstName": "string",
    "lastName": "string",
    "displayName": "string",
    "email": "string",
    "organizationId": "string",
    "organizationName": "string",
    "clinicId": "string",
    "clinicName": "string",
    "departmentId": "string",
    "departmentName": "string",
    "position": "string",
    "terminatedAt": "string",
    "currentVacationEndsAt": "string",
    "nextVacationStartsAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employee|[v1EmployeeCardView](#schemav1employeecardview)|false|none|EmployeeCardView is the denormalised card projection returned by<br>both Get and List endpoints. Fields mirror projections.employee_cards<br>columns; timestamps are RFC3339 strings. Optional fields stay unset<br>when the backing column is NULL.|

<h2 id="tocS_v1GetIncidentHistoryResponse">v1GetIncidentHistoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getincidenthistoryresponse"></a>
<a id="schema_v1GetIncidentHistoryResponse"></a>
<a id="tocSv1getincidenthistoryresponse"></a>
<a id="tocsv1getincidenthistoryresponse"></a>

```json
{
  "statusHistory": [
    {
      "id": "string",
      "oldStatus": "INCIDENT_STATUS_UNSPECIFIED",
      "newStatus": "INCIDENT_STATUS_UNSPECIFIED",
      "actor": {
        "employeeId": "string",
        "displayName": "string"
      },
      "changedAt": "string"
    }
  ],
  "priorityHistory": [
    {
      "id": "string",
      "oldPriority": "INCIDENT_PRIORITY_UNSPECIFIED",
      "newPriority": "INCIDENT_PRIORITY_UNSPECIFIED",
      "actor": {
        "employeeId": "string",
        "displayName": "string"
      },
      "changedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|statusHistory|[[queryIncidentV1StatusHistoryEntry](#schemaqueryincidentv1statushistoryentry)]|false|none|none|
|priorityHistory|[[v1PriorityHistoryEntry](#schemav1priorityhistoryentry)]|false|none|none|

<h2 id="tocS_v1GetIncidentResponse">v1GetIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getincidentresponse"></a>
<a id="schema_v1GetIncidentResponse"></a>
<a id="tocSv1getincidentresponse"></a>
<a id="tocsv1getincidentresponse"></a>

```json
{
  "incident": {
    "id": "string",
    "organizationId": "string",
    "clinicId": "string",
    "departmentId": "string",
    "categoryId": "string",
    "typeId": "string",
    "status": "INCIDENT_STATUS_UNSPECIFIED",
    "priority": "INCIDENT_PRIORITY_UNSPECIFIED",
    "description": "string",
    "patientOriginalDescription": "string",
    "occurredAt": "string",
    "createdAt": "string",
    "updatedAt": "string",
    "registrar": {
      "employeeId": "string",
      "displayName": "string",
      "position": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string"
    },
    "sourcePatientZitadelUserId": "string",
    "sourceBufferId": "string",
    "reopenedFromIncidentId": "string",
    "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|incident|[v1IncidentView](#schemav1incidentview)|false|none|IncidentView is the unified payload. For patients only id, status,<br>patient_status, description and timestamps are populated.|

<h2 id="tocS_v1GetMyClinicRoleResponse">v1GetMyClinicRoleResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getmyclinicroleresponse"></a>
<a id="schema_v1GetMyClinicRoleResponse"></a>
<a id="tocSv1getmyclinicroleresponse"></a>
<a id="tocsv1getmyclinicroleresponse"></a>

```json
{
  "isClinicHead": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|isClinicHead|boolean|false|none|none|

<h2 id="tocS_v1GetMyDepartmentRoleResponse">v1GetMyDepartmentRoleResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getmydepartmentroleresponse"></a>
<a id="schema_v1GetMyDepartmentRoleResponse"></a>
<a id="tocSv1getmydepartmentroleresponse"></a>
<a id="tocsv1getmydepartmentroleresponse"></a>

```json
{
  "isDepartmentResponsible": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|isDepartmentResponsible|boolean|false|none|none|

<h2 id="tocS_v1GetMyEmploymentResponse">v1GetMyEmploymentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getmyemploymentresponse"></a>
<a id="schema_v1GetMyEmploymentResponse"></a>
<a id="tocSv1getmyemploymentresponse"></a>
<a id="tocsv1getmyemploymentresponse"></a>

```json
{
  "employee": {
    "employeeId": "string",
    "zitadelUserId": "string",
    "firstName": "string",
    "lastName": "string",
    "displayName": "string",
    "email": "string",
    "organizationId": "string",
    "organizationName": "string",
    "clinicId": "string",
    "clinicName": "string",
    "departmentId": "string",
    "departmentName": "string",
    "position": "string",
    "terminatedAt": "string",
    "currentVacationEndsAt": "string",
    "nextVacationStartsAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employee|[v1EmployeeCardView](#schemav1employeecardview)|false|none|EmployeeCardView is the denormalised card projection returned by<br>both Get and List endpoints. Fields mirror projections.employee_cards<br>columns; timestamps are RFC3339 strings. Optional fields stay unset<br>when the backing column is NULL.|

<h2 id="tocS_v1GetMyIdentityResponse">v1GetMyIdentityResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getmyidentityresponse"></a>
<a id="schema_v1GetMyIdentityResponse"></a>
<a id="tocSv1getmyidentityresponse"></a>
<a id="tocsv1getmyidentityresponse"></a>

```json
{
  "isSystemAdmin": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|isSystemAdmin|boolean|false|none|none|

<h2 id="tocS_v1GetMyOrganizationRoleResponse">v1GetMyOrganizationRoleResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getmyorganizationroleresponse"></a>
<a id="schema_v1GetMyOrganizationRoleResponse"></a>
<a id="tocSv1getmyorganizationroleresponse"></a>
<a id="tocsv1getmyorganizationroleresponse"></a>

```json
{
  "isOrgAdmin": true,
  "isOrgHead": true,
  "isOrgDispatcher": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|isOrgAdmin|boolean|false|none|none|
|isOrgHead|boolean|false|none|none|
|isOrgDispatcher|boolean|false|none|none|

<h2 id="tocS_v1GetOrganizationResponse">v1GetOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getorganizationresponse"></a>
<a id="schema_v1GetOrganizationResponse"></a>
<a id="tocSv1getorganizationresponse"></a>
<a id="tocsv1getorganizationresponse"></a>

```json
{
  "organization": {
    "id": "string",
    "name": "string",
    "description": "string",
    "legalAddress": {
      "text": "string",
      "point": {
        "longitude": 0.1,
        "latitude": 0.1
      }
    },
    "createdAt": "string",
    "updatedAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|organization|[v1Organization](#schemav1organization)|false|none|Organization mirrors the projections.organizations row returned by<br>GetOrganization. Timestamps are RFC3339 strings.|

<h2 id="tocS_v1GetOrganizationStatsResponse">v1GetOrganizationStatsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getorganizationstatsresponse"></a>
<a id="schema_v1GetOrganizationStatsResponse"></a>
<a id="tocSv1getorganizationstatsresponse"></a>
<a id="tocsv1getorganizationstatsresponse"></a>

```json
{
  "stats": {
    "organizationId": "string",
    "employeesTotal": "string",
    "clinicsTotal": "string",
    "departmentsTotal": "string",
    "employeesOnVacation": "string",
    "vacationsScheduled": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|stats|[v1OrganizationStats](#schemav1organizationstats)|false|none|none|

<h2 id="tocS_v1GetRequestTypeResponse">v1GetRequestTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getrequesttyperesponse"></a>
<a id="schema_v1GetRequestTypeResponse"></a>
<a id="tocSv1getrequesttyperesponse"></a>
<a id="tocsv1getrequesttyperesponse"></a>

```json
{
  "requestType": {
    "id": "string",
    "organizationId": "string",
    "name": "string",
    "description": "string",
    "isActive": true,
    "createdAt": "string",
    "updatedAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|requestType|[v1RequestType](#schemav1requesttype)|false|none|none|

<h2 id="tocS_v1GetServiceRequestHistoryResponse">v1GetServiceRequestHistoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getservicerequesthistoryresponse"></a>
<a id="schema_v1GetServiceRequestHistoryResponse"></a>
<a id="tocSv1getservicerequesthistoryresponse"></a>
<a id="tocsv1getservicerequesthistoryresponse"></a>

```json
{
  "statusHistory": [
    {
      "id": "string",
      "oldStatus": "string",
      "newStatus": "string",
      "actorId": "string",
      "actorName": "string",
      "changedAt": "string"
    }
  ],
  "executorHistory": [
    {
      "id": "string",
      "action": "string",
      "employeeId": "string",
      "employeeName": "string",
      "actorId": "string",
      "actorName": "string",
      "changedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|statusHistory|[[queryRequestV1StatusHistoryEntry](#schemaqueryrequestv1statushistoryentry)]|false|none|none|
|executorHistory|[[v1ExecutorHistoryEntry](#schemav1executorhistoryentry)]|false|none|none|

<h2 id="tocS_v1GetServiceRequestResponse">v1GetServiceRequestResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getservicerequestresponse"></a>
<a id="schema_v1GetServiceRequestResponse"></a>
<a id="tocSv1getservicerequestresponse"></a>
<a id="tocsv1getservicerequestresponse"></a>

```json
{
  "serviceRequest": {
    "id": "string",
    "organizationId": "string",
    "clinicId": "string",
    "departmentId": "string",
    "typeId": "string",
    "incidentId": "string",
    "description": "string",
    "status": "string",
    "authorId": "string",
    "authorDisplayName": "string",
    "executors": [
      {
        "employeeId": "string",
        "assignedAt": "string",
        "assignedById": "string"
      }
    ],
    "createdAt": "string",
    "updatedAt": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|serviceRequest|[v1ServiceRequest](#schemav1servicerequest)|false|none|none|

<h2 id="tocS_v1GetSnapshotResponse">v1GetSnapshotResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getsnapshotresponse"></a>
<a id="schema_v1GetSnapshotResponse"></a>
<a id="tocSv1getsnapshotresponse"></a>
<a id="tocsv1getsnapshotresponse"></a>

```json
{
  "incidents": [
    {
      "createdAt": "string",
      "occurredAt": "string",
      "closedAt": "string",
      "status": "string",
      "priority": "string",
      "categoryId": "string",
      "categoryName": "string",
      "typeId": "string",
      "typeName": "string",
      "clinicId": "string",
      "departmentId": "string",
      "isPatientSource": true,
      "isReopened": true,
      "linkedRequestsCount": 0
    }
  ],
  "requests": [
    {
      "createdAt": "string",
      "completedAt": "string",
      "status": "string",
      "typeId": "string",
      "typeName": "string",
      "departmentId": "string",
      "hasLinkedIncident": true
    }
  ],
  "patientBuffer": [
    {
      "createdAt": "string",
      "status": "string",
      "categoryId": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|incidents|[[v1SnapshotIncident](#schemav1snapshotincident)]|false|none|none|
|requests|[[v1SnapshotRequest](#schemav1snapshotrequest)]|false|none|none|
|patientBuffer|[[v1SnapshotPatientBuffer](#schemav1snapshotpatientbuffer)]|false|none|none|

<h2 id="tocS_v1GetSummaryResponse">v1GetSummaryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1getsummaryresponse"></a>
<a id="schema_v1GetSummaryResponse"></a>
<a id="tocSv1getsummaryresponse"></a>
<a id="tocsv1getsummaryresponse"></a>

```json
{
  "incidents": {
    "total": "string",
    "byStatus": {
      "pending": "string",
      "inProgress": "string",
      "done": "string",
      "rejected": "string",
      "cancelled": "string"
    },
    "byPriority": {
      "low": "string",
      "normal": "string",
      "high": "string",
      "critical": "string"
    },
    "bySource": {
      "staff": "string",
      "patient": "string"
    },
    "reopened": "string",
    "withLinkedRequests": "string",
    "resolution": {
      "avgMinutes": 0.1,
      "minMinutes": 0.1,
      "maxMinutes": 0.1,
      "p50Minutes": 0.1,
      "p90Minutes": 0.1,
      "p95Minutes": 0.1
    },
    "topCategories": [
      {
        "categoryId": "string",
        "categoryName": "string",
        "count": "string"
      }
    ],
    "topTypes": [
      {
        "typeId": "string",
        "typeName": "string",
        "count": "string"
      }
    ],
    "topDepartments": [
      {
        "departmentId": "string",
        "departmentName": "string",
        "count": "string"
      }
    ]
  },
  "requests": {
    "total": "string",
    "byStatus": {
      "created": "string",
      "inWork": "string",
      "onHold": "string",
      "pendingReview": "string",
      "completed": "string",
      "cancelled": "string"
    },
    "linked": "string",
    "unlinked": "string",
    "completion": {
      "avgMinutes": 0.1,
      "minMinutes": 0.1,
      "maxMinutes": 0.1,
      "p50Minutes": 0.1,
      "p90Minutes": 0.1,
      "p95Minutes": 0.1
    },
    "topTypes": [
      {
        "typeId": "string",
        "typeName": "string",
        "count": "string"
      }
    ],
    "topDepartments": [
      {
        "departmentId": "string",
        "departmentName": "string",
        "count": "string"
      }
    ]
  },
  "patientBuffer": {
    "total": "string",
    "byStatus": {
      "pending": "string",
      "published": "string",
      "rejected": "string",
      "cancelled": "string"
    },
    "acceptanceRate": 0.1,
    "rejectionRate": 0.1
  },
  "period": {
    "organizationId": "string",
    "from": "string",
    "to": "string",
    "clinicId": "string",
    "departmentId": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|incidents|[v1IncidentSummary](#schemav1incidentsummary)|false|none|none|
|requests|[v1RequestSummary](#schemav1requestsummary)|false|none|none|
|patientBuffer|[v1PatientBufferSummary](#schemav1patientbuffersummary)|false|none|none|
|period|[v1SummaryPeriod](#schemav1summaryperiod)|false|none|none|

<h2 id="tocS_v1GetTimeSeriesResponse">v1GetTimeSeriesResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1gettimeseriesresponse"></a>
<a id="schema_v1GetTimeSeriesResponse"></a>
<a id="tocSv1gettimeseriesresponse"></a>
<a id="tocsv1gettimeseriesresponse"></a>

```json
{
  "buckets": [
    {
      "bucketStart": "string",
      "bucketEnd": "string",
      "incidents": {
        "total": "string",
        "pending": "string",
        "inProgress": "string",
        "done": "string",
        "rejected": "string",
        "cancelled": "string",
        "highCritical": "string",
        "patientSource": "string",
        "reopened": "string"
      },
      "requests": {
        "total": "string",
        "completed": "string",
        "cancelled": "string",
        "linked": "string"
      }
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|buckets|[[v1TimeSeriesBucket](#schemav1timeseriesbucket)]|false|none|none|

<h2 id="tocS_v1GetTypeResponse">v1GetTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1gettyperesponse"></a>
<a id="schema_v1GetTypeResponse"></a>
<a id="tocSv1gettyperesponse"></a>
<a id="tocsv1gettyperesponse"></a>

```json
{
  "type": {
    "id": "string",
    "organizationId": "string",
    "categoryId": "string",
    "name": "string",
    "description": "string",
    "isActive": true,
    "createdAt": "string",
    "updatedAt": "string",
    "isAllowedForPatients": true
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|type|[classifierV1Type](#schemaclassifierv1type)|false|none|Type mirrors projections.incident_types row.|

<h2 id="tocS_v1GrantSystemAdminRequest">v1GrantSystemAdminRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1grantsystemadminrequest"></a>
<a id="schema_v1GrantSystemAdminRequest"></a>
<a id="tocSv1grantsystemadminrequest"></a>
<a id="tocsv1grantsystemadminrequest"></a>

```json
{
  "zitadelUserId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|zitadelUserId|string|true|none|none|

<h2 id="tocS_v1GrantSystemAdminResponse">v1GrantSystemAdminResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1grantsystemadminresponse"></a>
<a id="schema_v1GrantSystemAdminResponse"></a>
<a id="tocSv1grantsystemadminresponse"></a>
<a id="tocsv1grantsystemadminresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1HireEmployeeRequest">v1HireEmployeeRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1hireemployeerequest"></a>
<a id="schema_v1HireEmployeeRequest"></a>
<a id="tocSv1hireemployeerequest"></a>
<a id="tocsv1hireemployeerequest"></a>

```json
{
  "zitadelUserId": "string",
  "departmentId": "string",
  "position": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|zitadelUserId|string|true|none|none|
|departmentId|string|true|none|none|
|position|string|false|none|none|

<h2 id="tocS_v1HireEmployeeResponse">v1HireEmployeeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1hireemployeeresponse"></a>
<a id="schema_v1HireEmployeeResponse"></a>
<a id="tocSv1hireemployeeresponse"></a>
<a id="tocsv1hireemployeeresponse"></a>

```json
{
  "employeeId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|false|none|none|

<h2 id="tocS_v1IncidentPriorityBreakdown">v1IncidentPriorityBreakdown</h2>
<!-- backwards compatibility -->
<a id="schemav1incidentprioritybreakdown"></a>
<a id="schema_v1IncidentPriorityBreakdown"></a>
<a id="tocSv1incidentprioritybreakdown"></a>
<a id="tocsv1incidentprioritybreakdown"></a>

```json
{
  "low": "string",
  "normal": "string",
  "high": "string",
  "critical": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|low|string(int64)|false|none|none|
|normal|string(int64)|false|none|none|
|high|string(int64)|false|none|none|
|critical|string(int64)|false|none|none|

<h2 id="tocS_v1IncidentSourceBreakdown">v1IncidentSourceBreakdown</h2>
<!-- backwards compatibility -->
<a id="schemav1incidentsourcebreakdown"></a>
<a id="schema_v1IncidentSourceBreakdown"></a>
<a id="tocSv1incidentsourcebreakdown"></a>
<a id="tocsv1incidentsourcebreakdown"></a>

```json
{
  "staff": "string",
  "patient": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|staff|string(int64)|false|none|none|
|patient|string(int64)|false|none|none|

<h2 id="tocS_v1IncidentStatusBreakdown">v1IncidentStatusBreakdown</h2>
<!-- backwards compatibility -->
<a id="schemav1incidentstatusbreakdown"></a>
<a id="schema_v1IncidentStatusBreakdown"></a>
<a id="tocSv1incidentstatusbreakdown"></a>
<a id="tocsv1incidentstatusbreakdown"></a>

```json
{
  "pending": "string",
  "inProgress": "string",
  "done": "string",
  "rejected": "string",
  "cancelled": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pending|string(int64)|false|none|none|
|inProgress|string(int64)|false|none|none|
|done|string(int64)|false|none|none|
|rejected|string(int64)|false|none|none|
|cancelled|string(int64)|false|none|none|

<h2 id="tocS_v1IncidentSummary">v1IncidentSummary</h2>
<!-- backwards compatibility -->
<a id="schemav1incidentsummary"></a>
<a id="schema_v1IncidentSummary"></a>
<a id="tocSv1incidentsummary"></a>
<a id="tocsv1incidentsummary"></a>

```json
{
  "total": "string",
  "byStatus": {
    "pending": "string",
    "inProgress": "string",
    "done": "string",
    "rejected": "string",
    "cancelled": "string"
  },
  "byPriority": {
    "low": "string",
    "normal": "string",
    "high": "string",
    "critical": "string"
  },
  "bySource": {
    "staff": "string",
    "patient": "string"
  },
  "reopened": "string",
  "withLinkedRequests": "string",
  "resolution": {
    "avgMinutes": 0.1,
    "minMinutes": 0.1,
    "maxMinutes": 0.1,
    "p50Minutes": 0.1,
    "p90Minutes": 0.1,
    "p95Minutes": 0.1
  },
  "topCategories": [
    {
      "categoryId": "string",
      "categoryName": "string",
      "count": "string"
    }
  ],
  "topTypes": [
    {
      "typeId": "string",
      "typeName": "string",
      "count": "string"
    }
  ],
  "topDepartments": [
    {
      "departmentId": "string",
      "departmentName": "string",
      "count": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|
|byStatus|[v1IncidentStatusBreakdown](#schemav1incidentstatusbreakdown)|false|none|none|
|byPriority|[v1IncidentPriorityBreakdown](#schemav1incidentprioritybreakdown)|false|none|none|
|bySource|[v1IncidentSourceBreakdown](#schemav1incidentsourcebreakdown)|false|none|none|
|reopened|string(int64)|false|none|none|
|withLinkedRequests|string(int64)|false|none|none|
|resolution|[v1ResolutionStats](#schemav1resolutionstats)|false|none|none|
|topCategories|[[v1CategoryCount](#schemav1categorycount)]|false|none|none|
|topTypes|[[v1TypeCount](#schemav1typecount)]|false|none|none|
|topDepartments|[[v1DepartmentCount](#schemav1departmentcount)]|false|none|none|

<h2 id="tocS_v1IncidentView">v1IncidentView</h2>
<!-- backwards compatibility -->
<a id="schemav1incidentview"></a>
<a id="schema_v1IncidentView"></a>
<a id="tocSv1incidentview"></a>
<a id="tocsv1incidentview"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "clinicId": "string",
  "departmentId": "string",
  "categoryId": "string",
  "typeId": "string",
  "status": "INCIDENT_STATUS_UNSPECIFIED",
  "priority": "INCIDENT_PRIORITY_UNSPECIFIED",
  "description": "string",
  "patientOriginalDescription": "string",
  "occurredAt": "string",
  "createdAt": "string",
  "updatedAt": "string",
  "registrar": {
    "employeeId": "string",
    "displayName": "string",
    "position": "string",
    "organizationId": "string",
    "clinicId": "string",
    "departmentId": "string"
  },
  "sourcePatientZitadelUserId": "string",
  "sourceBufferId": "string",
  "reopenedFromIncidentId": "string",
  "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
}

```

IncidentView is the unified payload. For patients only id, status,
patient_status, description and timestamps are populated.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|clinicId|string|false|none|none|
|departmentId|string|false|none|none|
|categoryId|string|false|none|none|
|typeId|string|false|none|none|
|status|[queryIncidentV1IncidentStatus](#schemaqueryincidentv1incidentstatus)|false|none|none|
|priority|[queryIncidentV1IncidentPriority](#schemaqueryincidentv1incidentpriority)|false|none|none|
|description|string|false|none|none|
|patientOriginalDescription|string|false|none|none|
|occurredAt|string|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|
|registrar|[v1RegistrarView](#schemav1registrarview)|false|none|none|
|sourcePatientZitadelUserId|string|false|none|none|
|sourceBufferId|string|false|none|none|
|reopenedFromIncidentId|string|false|none|none|
|patientStatus|[v1PatientStatus](#schemav1patientstatus)|false|none|PatientStatus is the simplified four-value status surfaced to patients.<br><br> - PATIENT_STATUS_PENDING: buffer pending<br> - PATIENT_STATUS_ACCEPTED: dispatcher accepted; incident pending/in_progress<br> - PATIENT_STATUS_CLOSED: done / rejected / buffer rejected<br> - PATIENT_STATUS_CANCELLED: patient cancelled|

<h2 id="tocS_v1ListActiveRequestTypesByOrganizationResponse">v1ListActiveRequestTypesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listactiverequesttypesbyorganizationresponse"></a>
<a id="schema_v1ListActiveRequestTypesByOrganizationResponse"></a>
<a id="tocSv1listactiverequesttypesbyorganizationresponse"></a>
<a id="tocsv1listactiverequesttypesbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1RequestType](#schemav1requesttype)]|false|none|none|

<h2 id="tocS_v1ListActiveRootCategoriesResponse">v1ListActiveRootCategoriesResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listactiverootcategoriesresponse"></a>
<a id="schema_v1ListActiveRootCategoriesResponse"></a>
<a id="tocSv1listactiverootcategoriesresponse"></a>
<a id="tocsv1listactiverootcategoriesresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "parentCategoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1Category](#schemav1category)]|false|none|[Category mirrors projections.incident_categories row.]|

<h2 id="tocS_v1ListActiveTypesByOrganizationResponse">v1ListActiveTypesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listactivetypesbyorganizationresponse"></a>
<a id="schema_v1ListActiveTypesByOrganizationResponse"></a>
<a id="tocSv1listactivetypesbyorganizationresponse"></a>
<a id="tocsv1listactivetypesbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "categoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string",
      "isAllowedForPatients": true
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[classifierV1Type](#schemaclassifierv1type)]|false|none|[Type mirrors projections.incident_types row.]|

<h2 id="tocS_v1ListAnnouncementsForClinicResponse">v1ListAnnouncementsForClinicResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listannouncementsforclinicresponse"></a>
<a id="schema_v1ListAnnouncementsForClinicResponse"></a>
<a id="tocSv1listannouncementsforclinicresponse"></a>
<a id="tocsv1listannouncementsforclinicresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "authorId": "string",
      "title": "string",
      "content": "string",
      "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
      "isArchived": true,
      "startsAt": "string",
      "endsAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "viewCount": "string"
    }
  ],
  "nextCursor": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1AnnouncementView](#schemav1announcementview)]|false|none|none|
|nextCursor|string|false|none|none|

<h2 id="tocS_v1ListAnnouncementsForDepartmentResponse">v1ListAnnouncementsForDepartmentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listannouncementsfordepartmentresponse"></a>
<a id="schema_v1ListAnnouncementsForDepartmentResponse"></a>
<a id="tocSv1listannouncementsfordepartmentresponse"></a>
<a id="tocsv1listannouncementsfordepartmentresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "authorId": "string",
      "title": "string",
      "content": "string",
      "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
      "isArchived": true,
      "startsAt": "string",
      "endsAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "viewCount": "string"
    }
  ],
  "nextCursor": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1AnnouncementView](#schemav1announcementview)]|false|none|none|
|nextCursor|string|false|none|none|

<h2 id="tocS_v1ListAnnouncementsForOrganizationResponse">v1ListAnnouncementsForOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listannouncementsfororganizationresponse"></a>
<a id="schema_v1ListAnnouncementsForOrganizationResponse"></a>
<a id="tocSv1listannouncementsfororganizationresponse"></a>
<a id="tocsv1listannouncementsfororganizationresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "authorId": "string",
      "title": "string",
      "content": "string",
      "priority": "ANNOUNCEMENT_PRIORITY_UNSPECIFIED",
      "isArchived": true,
      "startsAt": "string",
      "endsAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "viewCount": "string"
    }
  ],
  "nextCursor": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1AnnouncementView](#schemav1announcementview)]|false|none|none|
|nextCursor|string|false|none|none|

<h2 id="tocS_v1ListBufferEntriesResponse">v1ListBufferEntriesResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listbufferentriesresponse"></a>
<a id="schema_v1ListBufferEntriesResponse"></a>
<a id="tocSv1listbufferentriesresponse"></a>
<a id="tocsv1listbufferentriesresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "patientZitadelUserId": "string",
      "categoryId": "string",
      "typeId": "string",
      "description": "string",
      "occurredAt": "string",
      "status": "BUFFER_STATUS_UNSPECIFIED",
      "publishedIncidentId": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1BufferEntryView](#schemav1bufferentryview)]|false|none|none|

<h2 id="tocS_v1ListCategoriesByOrganizationResponse">v1ListCategoriesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listcategoriesbyorganizationresponse"></a>
<a id="schema_v1ListCategoriesByOrganizationResponse"></a>
<a id="tocSv1listcategoriesbyorganizationresponse"></a>
<a id="tocsv1listcategoriesbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "parentCategoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1Category](#schemav1category)]|false|none|[Category mirrors projections.incident_categories row.]|

<h2 id="tocS_v1ListCategorySubtreeResponse">v1ListCategorySubtreeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listcategorysubtreeresponse"></a>
<a id="schema_v1ListCategorySubtreeResponse"></a>
<a id="tocSv1listcategorysubtreeresponse"></a>
<a id="tocsv1listcategorysubtreeresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "parentCategoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1Category](#schemav1category)]|false|none|[Category mirrors projections.incident_categories row.]|

<h2 id="tocS_v1ListClinicsByOrganizationResponse">v1ListClinicsByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listclinicsbyorganizationresponse"></a>
<a id="schema_v1ListClinicsByOrganizationResponse"></a>
<a id="tocSv1listclinicsbyorganizationresponse"></a>
<a id="tocsv1listclinicsbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "name": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1ClinicListItem](#schemav1cliniclistitem)]|false|none|[ClinicListItem is the minimal shape returned by list endpoints.]|

<h2 id="tocS_v1ListDepartmentsByClinicResponse">v1ListDepartmentsByClinicResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listdepartmentsbyclinicresponse"></a>
<a id="schema_v1ListDepartmentsByClinicResponse"></a>
<a id="tocSv1listdepartmentsbyclinicresponse"></a>
<a id="tocsv1listdepartmentsbyclinicresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "clinicId": "string",
      "name": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1DepartmentListItem](#schemav1departmentlistitem)]|false|none|[DepartmentListItem is the minimal shape returned by list endpoints.]|

<h2 id="tocS_v1ListEmployeesByClinicResponse">v1ListEmployeesByClinicResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listemployeesbyclinicresponse"></a>
<a id="schema_v1ListEmployeesByClinicResponse"></a>
<a id="tocSv1listemployeesbyclinicresponse"></a>
<a id="tocsv1listemployeesbyclinicresponse"></a>

```json
{
  "items": [
    {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1EmployeeCardView](#schemav1employeecardview)]|false|none|[EmployeeCardView is the denormalised card projection returned by<br>both Get and List endpoints. Fields mirror projections.employee_cards<br>columns; timestamps are RFC3339 strings. Optional fields stay unset<br>when the backing column is NULL.]|

<h2 id="tocS_v1ListEmployeesByDepartmentResponse">v1ListEmployeesByDepartmentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listemployeesbydepartmentresponse"></a>
<a id="schema_v1ListEmployeesByDepartmentResponse"></a>
<a id="tocSv1listemployeesbydepartmentresponse"></a>
<a id="tocsv1listemployeesbydepartmentresponse"></a>

```json
{
  "items": [
    {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1EmployeeCardView](#schemav1employeecardview)]|false|none|[EmployeeCardView is the denormalised card projection returned by<br>both Get and List endpoints. Fields mirror projections.employee_cards<br>columns; timestamps are RFC3339 strings. Optional fields stay unset<br>when the backing column is NULL.]|

<h2 id="tocS_v1ListEmployeesByOrganizationResponse">v1ListEmployeesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listemployeesbyorganizationresponse"></a>
<a id="schema_v1ListEmployeesByOrganizationResponse"></a>
<a id="tocSv1listemployeesbyorganizationresponse"></a>
<a id="tocsv1listemployeesbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1EmployeeCardView](#schemav1employeecardview)]|false|none|[EmployeeCardView is the denormalised card projection returned by<br>both Get and List endpoints. Fields mirror projections.employee_cards<br>columns; timestamps are RFC3339 strings. Optional fields stay unset<br>when the backing column is NULL.]|

<h2 id="tocS_v1ListIncidentsResponse">v1ListIncidentsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listincidentsresponse"></a>
<a id="schema_v1ListIncidentsResponse"></a>
<a id="tocSv1listincidentsresponse"></a>
<a id="tocsv1listincidentsresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "categoryId": "string",
      "typeId": "string",
      "status": "INCIDENT_STATUS_UNSPECIFIED",
      "priority": "INCIDENT_PRIORITY_UNSPECIFIED",
      "description": "string",
      "patientOriginalDescription": "string",
      "occurredAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "registrar": {
        "employeeId": "string",
        "displayName": "string",
        "position": "string",
        "organizationId": "string",
        "clinicId": "string",
        "departmentId": "string"
      },
      "sourcePatientZitadelUserId": "string",
      "sourceBufferId": "string",
      "reopenedFromIncidentId": "string",
      "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1IncidentView](#schemav1incidentview)]|false|none|[IncidentView is the unified payload. For patients only id, status,<br>patient_status, description and timestamps are populated.]|

<h2 id="tocS_v1ListMyBufferEntriesResponse">v1ListMyBufferEntriesResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listmybufferentriesresponse"></a>
<a id="schema_v1ListMyBufferEntriesResponse"></a>
<a id="tocSv1listmybufferentriesresponse"></a>
<a id="tocsv1listmybufferentriesresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "patientZitadelUserId": "string",
      "categoryId": "string",
      "typeId": "string",
      "description": "string",
      "occurredAt": "string",
      "status": "BUFFER_STATUS_UNSPECIFIED",
      "publishedIncidentId": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1BufferEntryView](#schemav1bufferentryview)]|false|none|none|

<h2 id="tocS_v1ListMyIncidentsResponse">v1ListMyIncidentsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listmyincidentsresponse"></a>
<a id="schema_v1ListMyIncidentsResponse"></a>
<a id="tocSv1listmyincidentsresponse"></a>
<a id="tocsv1listmyincidentsresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "categoryId": "string",
      "typeId": "string",
      "status": "INCIDENT_STATUS_UNSPECIFIED",
      "priority": "INCIDENT_PRIORITY_UNSPECIFIED",
      "description": "string",
      "patientOriginalDescription": "string",
      "occurredAt": "string",
      "createdAt": "string",
      "updatedAt": "string",
      "registrar": {
        "employeeId": "string",
        "displayName": "string",
        "position": "string",
        "organizationId": "string",
        "clinicId": "string",
        "departmentId": "string"
      },
      "sourcePatientZitadelUserId": "string",
      "sourceBufferId": "string",
      "reopenedFromIncidentId": "string",
      "patientStatus": "PATIENT_STATUS_UNSPECIFIED"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1IncidentView](#schemav1incidentview)]|false|none|[IncidentView is the unified payload. For patients only id, status,<br>patient_status, description and timestamps are populated.]|

<h2 id="tocS_v1ListMyOrganizationsResponse">v1ListMyOrganizationsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listmyorganizationsresponse"></a>
<a id="schema_v1ListMyOrganizationsResponse"></a>
<a id="tocSv1listmyorganizationsresponse"></a>
<a id="tocsv1listmyorganizationsresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "name": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1OrganizationListItem](#schemav1organizationlistitem)]|false|none|[OrganizationListItem is the minimal shape returned by list endpoints.]|

<h2 id="tocS_v1ListOrgAdminsResponse">v1ListOrgAdminsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listorgadminsresponse"></a>
<a id="schema_v1ListOrgAdminsResponse"></a>
<a id="tocSv1listorgadminsresponse"></a>
<a id="tocsv1listorgadminsresponse"></a>

```json
{
  "items": [
    {
      "holder": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      },
      "deputy": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      }
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1RoleAssignment](#schemav1roleassignment)]|false|none|[RoleAssignment is a role row enriched with the denormalised card for<br>the holder and, when present, the deputy. Read-model callers use<br>this so they do not need to follow role lookups with N+1 GetEmployee<br>calls to render a name or email.]|

<h2 id="tocS_v1ListOrgDispatchersResponse">v1ListOrgDispatchersResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listorgdispatchersresponse"></a>
<a id="schema_v1ListOrgDispatchersResponse"></a>
<a id="tocSv1listorgdispatchersresponse"></a>
<a id="tocsv1listorgdispatchersresponse"></a>

```json
{
  "items": [
    {
      "holder": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      },
      "deputy": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      }
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1RoleAssignment](#schemav1roleassignment)]|false|none|[RoleAssignment is a role row enriched with the denormalised card for<br>the holder and, when present, the deputy. Read-model callers use<br>this so they do not need to follow role lookups with N+1 GetEmployee<br>calls to render a name or email.]|

<h2 id="tocS_v1ListOrgHeadsResponse">v1ListOrgHeadsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listorgheadsresponse"></a>
<a id="schema_v1ListOrgHeadsResponse"></a>
<a id="tocSv1listorgheadsresponse"></a>
<a id="tocsv1listorgheadsresponse"></a>

```json
{
  "items": [
    {
      "holder": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      },
      "deputy": {
        "employeeId": "string",
        "zitadelUserId": "string",
        "firstName": "string",
        "lastName": "string",
        "displayName": "string",
        "email": "string",
        "organizationId": "string",
        "organizationName": "string",
        "clinicId": "string",
        "clinicName": "string",
        "departmentId": "string",
        "departmentName": "string",
        "position": "string",
        "terminatedAt": "string",
        "currentVacationEndsAt": "string",
        "nextVacationStartsAt": "string"
      }
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1RoleAssignment](#schemav1roleassignment)]|false|none|[RoleAssignment is a role row enriched with the denormalised card for<br>the holder and, when present, the deputy. Read-model callers use<br>this so they do not need to follow role lookups with N+1 GetEmployee<br>calls to render a name or email.]|

<h2 id="tocS_v1ListOrganizationsResponse">v1ListOrganizationsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listorganizationsresponse"></a>
<a id="schema_v1ListOrganizationsResponse"></a>
<a id="tocSv1listorganizationsresponse"></a>
<a id="tocsv1listorganizationsresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "name": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1OrganizationListItem](#schemav1organizationlistitem)]|false|none|[OrganizationListItem is the minimal shape returned by list endpoints.]|

<h2 id="tocS_v1ListPatientAllowedTypesByOrganizationResponse">v1ListPatientAllowedTypesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listpatientallowedtypesbyorganizationresponse"></a>
<a id="schema_v1ListPatientAllowedTypesByOrganizationResponse"></a>
<a id="tocSv1listpatientallowedtypesbyorganizationresponse"></a>
<a id="tocsv1listpatientallowedtypesbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "categoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string",
      "isAllowedForPatients": true
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[classifierV1Type](#schemaclassifierv1type)]|false|none|[Type mirrors projections.incident_types row.]|

<h2 id="tocS_v1ListPatientVisibleCategoriesByOrganizationResponse">v1ListPatientVisibleCategoriesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listpatientvisiblecategoriesbyorganizationresponse"></a>
<a id="schema_v1ListPatientVisibleCategoriesByOrganizationResponse"></a>
<a id="tocSv1listpatientvisiblecategoriesbyorganizationresponse"></a>
<a id="tocsv1listpatientvisiblecategoriesbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "parentCategoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1Category](#schemav1category)]|false|none|[Category mirrors projections.incident_categories row.]|

<h2 id="tocS_v1ListRequestTypesByOrganizationResponse">v1ListRequestTypesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listrequesttypesbyorganizationresponse"></a>
<a id="schema_v1ListRequestTypesByOrganizationResponse"></a>
<a id="tocSv1listrequesttypesbyorganizationresponse"></a>
<a id="tocsv1listrequesttypesbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1RequestType](#schemav1requesttype)]|false|none|none|

<h2 id="tocS_v1ListServiceRequestsByIncidentResponse">v1ListServiceRequestsByIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listservicerequestsbyincidentresponse"></a>
<a id="schema_v1ListServiceRequestsByIncidentResponse"></a>
<a id="tocSv1listservicerequestsbyincidentresponse"></a>
<a id="tocsv1listservicerequestsbyincidentresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "typeId": "string",
      "incidentId": "string",
      "description": "string",
      "status": "string",
      "authorId": "string",
      "authorDisplayName": "string",
      "executors": [
        {
          "employeeId": "string",
          "assignedAt": "string",
          "assignedById": "string"
        }
      ],
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1ServiceRequest](#schemav1servicerequest)]|false|none|none|

<h2 id="tocS_v1ListServiceRequestsResponse">v1ListServiceRequestsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listservicerequestsresponse"></a>
<a id="schema_v1ListServiceRequestsResponse"></a>
<a id="tocSv1listservicerequestsresponse"></a>
<a id="tocsv1listservicerequestsresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "clinicId": "string",
      "departmentId": "string",
      "typeId": "string",
      "incidentId": "string",
      "description": "string",
      "status": "string",
      "authorId": "string",
      "authorDisplayName": "string",
      "executors": [
        {
          "employeeId": "string",
          "assignedAt": "string",
          "assignedById": "string"
        }
      ],
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1ServiceRequest](#schemav1servicerequest)]|false|none|none|

<h2 id="tocS_v1ListSystemAdminsResponse">v1ListSystemAdminsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listsystemadminsresponse"></a>
<a id="schema_v1ListSystemAdminsResponse"></a>
<a id="tocSv1listsystemadminsresponse"></a>
<a id="tocsv1listsystemadminsresponse"></a>

```json
{
  "items": [
    {
      "zitadelUserId": "string",
      "createdAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1SystemAdminView](#schemav1systemadminview)]|false|none|[SystemAdminView is the system-admin role; system admins are rooted in<br>Zitadel user ids, not employee ids.]|

<h2 id="tocS_v1ListTypesByCategoryResponse">v1ListTypesByCategoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listtypesbycategoryresponse"></a>
<a id="schema_v1ListTypesByCategoryResponse"></a>
<a id="tocSv1listtypesbycategoryresponse"></a>
<a id="tocsv1listtypesbycategoryresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "organizationId": "string",
      "categoryId": "string",
      "name": "string",
      "description": "string",
      "isActive": true,
      "createdAt": "string",
      "updatedAt": "string",
      "isAllowedForPatients": true
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[classifierV1Type](#schemaclassifierv1type)]|false|none|[Type mirrors projections.incident_types row.]|

<h2 id="tocS_v1ListVacationsByEmployeeResponse">v1ListVacationsByEmployeeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1listvacationsbyemployeeresponse"></a>
<a id="schema_v1ListVacationsByEmployeeResponse"></a>
<a id="tocSv1listvacationsbyemployeeresponse"></a>
<a id="tocsv1listvacationsbyemployeeresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "employeeId": "string",
      "state": "string",
      "startsAt": "string",
      "endsAt": "string",
      "createdAt": "string",
      "updatedAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1VacationView](#schemav1vacationview)]|false|none|[VacationView mirrors projections.employee_vacations. state is one of<br>{scheduled, active, ended, cancelled}.]|

<h2 id="tocS_v1MoveIncidentCategoryResponse">v1MoveIncidentCategoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1moveincidentcategoryresponse"></a>
<a id="schema_v1MoveIncidentCategoryResponse"></a>
<a id="tocSv1moveincidentcategoryresponse"></a>
<a id="tocsv1moveincidentcategoryresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1MoveIncidentTypeResponse">v1MoveIncidentTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1moveincidenttyperesponse"></a>
<a id="schema_v1MoveIncidentTypeResponse"></a>
<a id="tocSv1moveincidenttyperesponse"></a>
<a id="tocsv1moveincidenttyperesponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1Organization">v1Organization</h2>
<!-- backwards compatibility -->
<a id="schemav1organization"></a>
<a id="schema_v1Organization"></a>
<a id="tocSv1organization"></a>
<a id="tocsv1organization"></a>

```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "legalAddress": {
    "text": "string",
    "point": {
      "longitude": 0.1,
      "latitude": 0.1
    }
  },
  "createdAt": "string",
  "updatedAt": "string"
}

```

Organization mirrors the projections.organizations row returned by
GetOrganization. Timestamps are RFC3339 strings.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|legalAddress|[v1Address](#schemav1address)|false|none|Address is the read-side view of a stored postal address; Point is<br>optional because the command-side allows text-only addresses.|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|

<h2 id="tocS_v1OrganizationListItem">v1OrganizationListItem</h2>
<!-- backwards compatibility -->
<a id="schemav1organizationlistitem"></a>
<a id="schema_v1OrganizationListItem"></a>
<a id="tocSv1organizationlistitem"></a>
<a id="tocsv1organizationlistitem"></a>

```json
{
  "id": "string",
  "name": "string"
}

```

OrganizationListItem is the minimal shape returned by list endpoints.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|name|string|false|none|none|

<h2 id="tocS_v1OrganizationStats">v1OrganizationStats</h2>
<!-- backwards compatibility -->
<a id="schemav1organizationstats"></a>
<a id="schema_v1OrganizationStats"></a>
<a id="tocSv1organizationstats"></a>
<a id="tocsv1organizationstats"></a>

```json
{
  "organizationId": "string",
  "employeesTotal": "string",
  "clinicsTotal": "string",
  "departmentsTotal": "string",
  "employeesOnVacation": "string",
  "vacationsScheduled": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|organizationId|string|false|none|none|
|employeesTotal|string(int64)|false|none|none|
|clinicsTotal|string(int64)|false|none|none|
|departmentsTotal|string(int64)|false|none|none|
|employeesOnVacation|string(int64)|false|none|none|
|vacationsScheduled|string(int64)|false|none|none|

<h2 id="tocS_v1PatientBufferStatusBreakdown">v1PatientBufferStatusBreakdown</h2>
<!-- backwards compatibility -->
<a id="schemav1patientbufferstatusbreakdown"></a>
<a id="schema_v1PatientBufferStatusBreakdown"></a>
<a id="tocSv1patientbufferstatusbreakdown"></a>
<a id="tocsv1patientbufferstatusbreakdown"></a>

```json
{
  "pending": "string",
  "published": "string",
  "rejected": "string",
  "cancelled": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|pending|string(int64)|false|none|none|
|published|string(int64)|false|none|none|
|rejected|string(int64)|false|none|none|
|cancelled|string(int64)|false|none|none|

<h2 id="tocS_v1PatientBufferSummary">v1PatientBufferSummary</h2>
<!-- backwards compatibility -->
<a id="schemav1patientbuffersummary"></a>
<a id="schema_v1PatientBufferSummary"></a>
<a id="tocSv1patientbuffersummary"></a>
<a id="tocsv1patientbuffersummary"></a>

```json
{
  "total": "string",
  "byStatus": {
    "pending": "string",
    "published": "string",
    "rejected": "string",
    "cancelled": "string"
  },
  "acceptanceRate": 0.1,
  "rejectionRate": 0.1
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|
|byStatus|[v1PatientBufferStatusBreakdown](#schemav1patientbufferstatusbreakdown)|false|none|none|
|acceptanceRate|number(double)|false|none|none|
|rejectionRate|number(double)|false|none|none|

<h2 id="tocS_v1PatientStatus">v1PatientStatus</h2>
<!-- backwards compatibility -->
<a id="schemav1patientstatus"></a>
<a id="schema_v1PatientStatus"></a>
<a id="tocSv1patientstatus"></a>
<a id="tocsv1patientstatus"></a>

```json
"PATIENT_STATUS_UNSPECIFIED"

```

PatientStatus is the simplified four-value status surfaced to patients.

 - PATIENT_STATUS_PENDING: buffer pending
 - PATIENT_STATUS_ACCEPTED: dispatcher accepted; incident pending/in_progress
 - PATIENT_STATUS_CLOSED: done / rejected / buffer rejected
 - PATIENT_STATUS_CANCELLED: patient cancelled

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|PatientStatus is the simplified four-value status surfaced to patients.<br><br> - PATIENT_STATUS_PENDING: buffer pending<br> - PATIENT_STATUS_ACCEPTED: dispatcher accepted; incident pending/in_progress<br> - PATIENT_STATUS_CLOSED: done / rejected / buffer rejected<br> - PATIENT_STATUS_CANCELLED: patient cancelled|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|PATIENT_STATUS_UNSPECIFIED|
|*anonymous*|PATIENT_STATUS_PENDING|
|*anonymous*|PATIENT_STATUS_ACCEPTED|
|*anonymous*|PATIENT_STATUS_CLOSED|
|*anonymous*|PATIENT_STATUS_CANCELLED|

<h2 id="tocS_v1Point">v1Point</h2>
<!-- backwards compatibility -->
<a id="schemav1point"></a>
<a id="schema_v1Point"></a>
<a id="tocSv1point"></a>
<a id="tocsv1point"></a>

```json
{
  "longitude": 0.1,
  "latitude": 0.1
}

```

Point is the coordinate pair attached to an Address. Absent when the
projection row has neither longitude nor latitude.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|longitude|number(double)|false|none|none|
|latitude|number(double)|false|none|none|

<h2 id="tocS_v1PointInput">v1PointInput</h2>
<!-- backwards compatibility -->
<a id="schemav1pointinput"></a>
<a id="schema_v1PointInput"></a>
<a id="tocSv1pointinput"></a>
<a id="tocsv1pointinput"></a>

```json
{
  "longitude": 0.1,
  "latitude": 0.1
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|longitude|number(double)|false|none|none|
|latitude|number(double)|false|none|none|

<h2 id="tocS_v1PriorityHistoryEntry">v1PriorityHistoryEntry</h2>
<!-- backwards compatibility -->
<a id="schemav1priorityhistoryentry"></a>
<a id="schema_v1PriorityHistoryEntry"></a>
<a id="tocSv1priorityhistoryentry"></a>
<a id="tocsv1priorityhistoryentry"></a>

```json
{
  "id": "string",
  "oldPriority": "INCIDENT_PRIORITY_UNSPECIFIED",
  "newPriority": "INCIDENT_PRIORITY_UNSPECIFIED",
  "actor": {
    "employeeId": "string",
    "displayName": "string"
  },
  "changedAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|oldPriority|[queryIncidentV1IncidentPriority](#schemaqueryincidentv1incidentpriority)|false|none|none|
|newPriority|[queryIncidentV1IncidentPriority](#schemaqueryincidentv1incidentpriority)|false|none|none|
|actor|[v1ActorView](#schemav1actorview)|false|none|none|
|changedAt|string|false|none|none|

<h2 id="tocS_v1PublishPatientIncidentResponse">v1PublishPatientIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1publishpatientincidentresponse"></a>
<a id="schema_v1PublishPatientIncidentResponse"></a>
<a id="tocSv1publishpatientincidentresponse"></a>
<a id="tocsv1publishpatientincidentresponse"></a>

```json
{
  "incidentId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|incidentId|string|false|none|none|

<h2 id="tocS_v1ReactivateIncidentCategoryResponse">v1ReactivateIncidentCategoryResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1reactivateincidentcategoryresponse"></a>
<a id="schema_v1ReactivateIncidentCategoryResponse"></a>
<a id="tocSv1reactivateincidentcategoryresponse"></a>
<a id="tocsv1reactivateincidentcategoryresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1ReactivateIncidentTypeResponse">v1ReactivateIncidentTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1reactivateincidenttyperesponse"></a>
<a id="schema_v1ReactivateIncidentTypeResponse"></a>
<a id="tocSv1reactivateincidenttyperesponse"></a>
<a id="tocsv1reactivateincidenttyperesponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1ReactivateRequestTypeResponse">v1ReactivateRequestTypeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1reactivaterequesttyperesponse"></a>
<a id="schema_v1ReactivateRequestTypeResponse"></a>
<a id="tocSv1reactivaterequesttyperesponse"></a>
<a id="tocsv1reactivaterequesttyperesponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RegistrarView">v1RegistrarView</h2>
<!-- backwards compatibility -->
<a id="schemav1registrarview"></a>
<a id="schema_v1RegistrarView"></a>
<a id="tocSv1registrarview"></a>
<a id="tocsv1registrarview"></a>

```json
{
  "employeeId": "string",
  "displayName": "string",
  "position": "string",
  "organizationId": "string",
  "clinicId": "string",
  "departmentId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|employeeId|string|false|none|none|
|displayName|string|false|none|none|
|position|string|false|none|none|
|organizationId|string|false|none|none|
|clinicId|string|false|none|none|
|departmentId|string|false|none|none|

<h2 id="tocS_v1RejectPatientIncidentResponse">v1RejectPatientIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1rejectpatientincidentresponse"></a>
<a id="schema_v1RejectPatientIncidentResponse"></a>
<a id="tocSv1rejectpatientincidentresponse"></a>
<a id="tocsv1rejectpatientincidentresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RemoveClinicHeadDeputyResponse">v1RemoveClinicHeadDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1removeclinicheaddeputyresponse"></a>
<a id="schema_v1RemoveClinicHeadDeputyResponse"></a>
<a id="tocSv1removeclinicheaddeputyresponse"></a>
<a id="tocsv1removeclinicheaddeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RemoveDepartmentResponsibleDeputyResponse">v1RemoveDepartmentResponsibleDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1removedepartmentresponsibledeputyresponse"></a>
<a id="schema_v1RemoveDepartmentResponsibleDeputyResponse"></a>
<a id="tocSv1removedepartmentresponsibledeputyresponse"></a>
<a id="tocsv1removedepartmentresponsibledeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RemoveOrganizationAdminDeputyResponse">v1RemoveOrganizationAdminDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1removeorganizationadmindeputyresponse"></a>
<a id="schema_v1RemoveOrganizationAdminDeputyResponse"></a>
<a id="tocSv1removeorganizationadmindeputyresponse"></a>
<a id="tocsv1removeorganizationadmindeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RemoveOrganizationDispatcherDeputyResponse">v1RemoveOrganizationDispatcherDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1removeorganizationdispatcherdeputyresponse"></a>
<a id="schema_v1RemoveOrganizationDispatcherDeputyResponse"></a>
<a id="tocSv1removeorganizationdispatcherdeputyresponse"></a>
<a id="tocsv1removeorganizationdispatcherdeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RemoveOrganizationHeadDeputyResponse">v1RemoveOrganizationHeadDeputyResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1removeorganizationheaddeputyresponse"></a>
<a id="schema_v1RemoveOrganizationHeadDeputyResponse"></a>
<a id="tocSv1removeorganizationheaddeputyresponse"></a>
<a id="tocsv1removeorganizationheaddeputyresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1ReopenIncidentResponse">v1ReopenIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1reopenincidentresponse"></a>
<a id="schema_v1ReopenIncidentResponse"></a>
<a id="tocSv1reopenincidentresponse"></a>
<a id="tocsv1reopenincidentresponse"></a>

```json
{
  "reopenedIncidentId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|reopenedIncidentId|string|false|none|none|

<h2 id="tocS_v1RequestStatusBreakdown">v1RequestStatusBreakdown</h2>
<!-- backwards compatibility -->
<a id="schemav1requeststatusbreakdown"></a>
<a id="schema_v1RequestStatusBreakdown"></a>
<a id="tocSv1requeststatusbreakdown"></a>
<a id="tocsv1requeststatusbreakdown"></a>

```json
{
  "created": "string",
  "inWork": "string",
  "onHold": "string",
  "pendingReview": "string",
  "completed": "string",
  "cancelled": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|created|string(int64)|false|none|none|
|inWork|string(int64)|false|none|none|
|onHold|string(int64)|false|none|none|
|pendingReview|string(int64)|false|none|none|
|completed|string(int64)|false|none|none|
|cancelled|string(int64)|false|none|none|

<h2 id="tocS_v1RequestSummary">v1RequestSummary</h2>
<!-- backwards compatibility -->
<a id="schemav1requestsummary"></a>
<a id="schema_v1RequestSummary"></a>
<a id="tocSv1requestsummary"></a>
<a id="tocsv1requestsummary"></a>

```json
{
  "total": "string",
  "byStatus": {
    "created": "string",
    "inWork": "string",
    "onHold": "string",
    "pendingReview": "string",
    "completed": "string",
    "cancelled": "string"
  },
  "linked": "string",
  "unlinked": "string",
  "completion": {
    "avgMinutes": 0.1,
    "minMinutes": 0.1,
    "maxMinutes": 0.1,
    "p50Minutes": 0.1,
    "p90Minutes": 0.1,
    "p95Minutes": 0.1
  },
  "topTypes": [
    {
      "typeId": "string",
      "typeName": "string",
      "count": "string"
    }
  ],
  "topDepartments": [
    {
      "departmentId": "string",
      "departmentName": "string",
      "count": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|
|byStatus|[v1RequestStatusBreakdown](#schemav1requeststatusbreakdown)|false|none|none|
|linked|string(int64)|false|none|none|
|unlinked|string(int64)|false|none|none|
|completion|[v1ResolutionStats](#schemav1resolutionstats)|false|none|none|
|topTypes|[[v1TypeCount](#schemav1typecount)]|false|none|none|
|topDepartments|[[v1DepartmentCount](#schemav1departmentcount)]|false|none|none|

<h2 id="tocS_v1RequestType">v1RequestType</h2>
<!-- backwards compatibility -->
<a id="schemav1requesttype"></a>
<a id="schema_v1RequestType"></a>
<a id="tocSv1requesttype"></a>
<a id="tocsv1requesttype"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "name": "string",
  "description": "string",
  "isActive": true,
  "createdAt": "string",
  "updatedAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|name|string|false|none|none|
|description|string|false|none|none|
|isActive|boolean|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|

<h2 id="tocS_v1ResolutionStats">v1ResolutionStats</h2>
<!-- backwards compatibility -->
<a id="schemav1resolutionstats"></a>
<a id="schema_v1ResolutionStats"></a>
<a id="tocSv1resolutionstats"></a>
<a id="tocsv1resolutionstats"></a>

```json
{
  "avgMinutes": 0.1,
  "minMinutes": 0.1,
  "maxMinutes": 0.1,
  "p50Minutes": 0.1,
  "p90Minutes": 0.1,
  "p95Minutes": 0.1
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|avgMinutes|number(double)|false|none|none|
|minMinutes|number(double)|false|none|none|
|maxMinutes|number(double)|false|none|none|
|p50Minutes|number(double)|false|none|none|
|p90Minutes|number(double)|false|none|none|
|p95Minutes|number(double)|false|none|none|

<h2 id="tocS_v1RevokeClinicHeadResponse">v1RevokeClinicHeadResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1revokeclinicheadresponse"></a>
<a id="schema_v1RevokeClinicHeadResponse"></a>
<a id="tocSv1revokeclinicheadresponse"></a>
<a id="tocsv1revokeclinicheadresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RevokeDepartmentResponsibleResponse">v1RevokeDepartmentResponsibleResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1revokedepartmentresponsibleresponse"></a>
<a id="schema_v1RevokeDepartmentResponsibleResponse"></a>
<a id="tocSv1revokedepartmentresponsibleresponse"></a>
<a id="tocsv1revokedepartmentresponsibleresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RevokeOrganizationAdminResponse">v1RevokeOrganizationAdminResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1revokeorganizationadminresponse"></a>
<a id="schema_v1RevokeOrganizationAdminResponse"></a>
<a id="tocSv1revokeorganizationadminresponse"></a>
<a id="tocsv1revokeorganizationadminresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RevokeOrganizationDispatcherResponse">v1RevokeOrganizationDispatcherResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1revokeorganizationdispatcherresponse"></a>
<a id="schema_v1RevokeOrganizationDispatcherResponse"></a>
<a id="tocSv1revokeorganizationdispatcherresponse"></a>
<a id="tocsv1revokeorganizationdispatcherresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RevokeOrganizationHeadResponse">v1RevokeOrganizationHeadResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1revokeorganizationheadresponse"></a>
<a id="schema_v1RevokeOrganizationHeadResponse"></a>
<a id="tocSv1revokeorganizationheadresponse"></a>
<a id="tocsv1revokeorganizationheadresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RevokeSystemAdminResponse">v1RevokeSystemAdminResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1revokesystemadminresponse"></a>
<a id="schema_v1RevokeSystemAdminResponse"></a>
<a id="tocSv1revokesystemadminresponse"></a>
<a id="tocsv1revokesystemadminresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1RoleAssignment">v1RoleAssignment</h2>
<!-- backwards compatibility -->
<a id="schemav1roleassignment"></a>
<a id="schema_v1RoleAssignment"></a>
<a id="tocSv1roleassignment"></a>
<a id="tocsv1roleassignment"></a>

```json
{
  "holder": {
    "employeeId": "string",
    "zitadelUserId": "string",
    "firstName": "string",
    "lastName": "string",
    "displayName": "string",
    "email": "string",
    "organizationId": "string",
    "organizationName": "string",
    "clinicId": "string",
    "clinicName": "string",
    "departmentId": "string",
    "departmentName": "string",
    "position": "string",
    "terminatedAt": "string",
    "currentVacationEndsAt": "string",
    "nextVacationStartsAt": "string"
  },
  "deputy": {
    "employeeId": "string",
    "zitadelUserId": "string",
    "firstName": "string",
    "lastName": "string",
    "displayName": "string",
    "email": "string",
    "organizationId": "string",
    "organizationName": "string",
    "clinicId": "string",
    "clinicName": "string",
    "departmentId": "string",
    "departmentName": "string",
    "position": "string",
    "terminatedAt": "string",
    "currentVacationEndsAt": "string",
    "nextVacationStartsAt": "string"
  }
}

```

RoleAssignment is a role row enriched with the denormalised card for
the holder and, when present, the deputy. Read-model callers use
this so they do not need to follow role lookups with N+1 GetEmployee
calls to render a name or email.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|holder|[v1EmployeeCardView](#schemav1employeecardview)|false|none|EmployeeCardView is the denormalised card projection returned by<br>both Get and List endpoints. Fields mirror projections.employee_cards<br>columns; timestamps are RFC3339 strings. Optional fields stay unset<br>when the backing column is NULL.|
|deputy|[v1EmployeeCardView](#schemav1employeecardview)|false|none|EmployeeCardView is the denormalised card projection returned by<br>both Get and List endpoints. Fields mirror projections.employee_cards<br>columns; timestamps are RFC3339 strings. Optional fields stay unset<br>when the backing column is NULL.|

<h2 id="tocS_v1ScheduleVacationResponse">v1ScheduleVacationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1schedulevacationresponse"></a>
<a id="schema_v1ScheduleVacationResponse"></a>
<a id="tocSv1schedulevacationresponse"></a>
<a id="tocsv1schedulevacationresponse"></a>

```json
{
  "vacationId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|vacationId|string|false|none|none|

<h2 id="tocS_v1SearchEmployeesByOrganizationResponse">v1SearchEmployeesByOrganizationResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1searchemployeesbyorganizationresponse"></a>
<a id="schema_v1SearchEmployeesByOrganizationResponse"></a>
<a id="tocSv1searchemployeesbyorganizationresponse"></a>
<a id="tocsv1searchemployeesbyorganizationresponse"></a>

```json
{
  "items": [
    {
      "employeeId": "string",
      "zitadelUserId": "string",
      "firstName": "string",
      "lastName": "string",
      "displayName": "string",
      "email": "string",
      "organizationId": "string",
      "organizationName": "string",
      "clinicId": "string",
      "clinicName": "string",
      "departmentId": "string",
      "departmentName": "string",
      "position": "string",
      "terminatedAt": "string",
      "currentVacationEndsAt": "string",
      "nextVacationStartsAt": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1EmployeeCardView](#schemav1employeecardview)]|false|none|[EmployeeCardView is the denormalised card projection returned by<br>both Get and List endpoints. Fields mirror projections.employee_cards<br>columns; timestamps are RFC3339 strings. Optional fields stay unset<br>when the backing column is NULL.]|

<h2 id="tocS_v1SearchOrganizationsResponse">v1SearchOrganizationsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1searchorganizationsresponse"></a>
<a id="schema_v1SearchOrganizationsResponse"></a>
<a id="tocSv1searchorganizationsresponse"></a>
<a id="tocsv1searchorganizationsresponse"></a>

```json
{
  "items": [
    {
      "id": "string",
      "name": "string"
    }
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|items|[[v1OrganizationListItem](#schemav1organizationlistitem)]|false|none|[OrganizationListItem is the minimal shape returned by list endpoints.]|

<h2 id="tocS_v1ServiceRequest">v1ServiceRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1servicerequest"></a>
<a id="schema_v1ServiceRequest"></a>
<a id="tocSv1servicerequest"></a>
<a id="tocsv1servicerequest"></a>

```json
{
  "id": "string",
  "organizationId": "string",
  "clinicId": "string",
  "departmentId": "string",
  "typeId": "string",
  "incidentId": "string",
  "description": "string",
  "status": "string",
  "authorId": "string",
  "authorDisplayName": "string",
  "executors": [
    {
      "employeeId": "string",
      "assignedAt": "string",
      "assignedById": "string"
    }
  ],
  "createdAt": "string",
  "updatedAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|organizationId|string|false|none|none|
|clinicId|string|false|none|none|
|departmentId|string|false|none|none|
|typeId|string|false|none|none|
|incidentId|string|false|none|none|
|description|string|false|none|none|
|status|string|false|none|none|
|authorId|string|false|none|none|
|authorDisplayName|string|false|none|none|
|executors|[[v1Executor](#schemav1executor)]|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|

<h2 id="tocS_v1SnapshotIncident">v1SnapshotIncident</h2>
<!-- backwards compatibility -->
<a id="schemav1snapshotincident"></a>
<a id="schema_v1SnapshotIncident"></a>
<a id="tocSv1snapshotincident"></a>
<a id="tocsv1snapshotincident"></a>

```json
{
  "createdAt": "string",
  "occurredAt": "string",
  "closedAt": "string",
  "status": "string",
  "priority": "string",
  "categoryId": "string",
  "categoryName": "string",
  "typeId": "string",
  "typeName": "string",
  "clinicId": "string",
  "departmentId": "string",
  "isPatientSource": true,
  "isReopened": true,
  "linkedRequestsCount": 0
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|createdAt|string|false|none|none|
|occurredAt|string|false|none|none|
|closedAt|string|false|none|none|
|status|string|false|none|none|
|priority|string|false|none|none|
|categoryId|string|false|none|none|
|categoryName|string|false|none|none|
|typeId|string|false|none|none|
|typeName|string|false|none|none|
|clinicId|string|false|none|none|
|departmentId|string|false|none|none|
|isPatientSource|boolean|false|none|none|
|isReopened|boolean|false|none|none|
|linkedRequestsCount|integer(int32)|false|none|none|

<h2 id="tocS_v1SnapshotPatientBuffer">v1SnapshotPatientBuffer</h2>
<!-- backwards compatibility -->
<a id="schemav1snapshotpatientbuffer"></a>
<a id="schema_v1SnapshotPatientBuffer"></a>
<a id="tocSv1snapshotpatientbuffer"></a>
<a id="tocsv1snapshotpatientbuffer"></a>

```json
{
  "createdAt": "string",
  "status": "string",
  "categoryId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|createdAt|string|false|none|none|
|status|string|false|none|none|
|categoryId|string|false|none|none|

<h2 id="tocS_v1SnapshotRequest">v1SnapshotRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1snapshotrequest"></a>
<a id="schema_v1SnapshotRequest"></a>
<a id="tocSv1snapshotrequest"></a>
<a id="tocsv1snapshotrequest"></a>

```json
{
  "createdAt": "string",
  "completedAt": "string",
  "status": "string",
  "typeId": "string",
  "typeName": "string",
  "departmentId": "string",
  "hasLinkedIncident": true
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|createdAt|string|false|none|none|
|completedAt|string|false|none|none|
|status|string|false|none|none|
|typeId|string|false|none|none|
|typeName|string|false|none|none|
|departmentId|string|false|none|none|
|hasLinkedIncident|boolean|false|none|none|

<h2 id="tocS_v1StartVacationNowResponse">v1StartVacationNowResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1startvacationnowresponse"></a>
<a id="schema_v1StartVacationNowResponse"></a>
<a id="tocSv1startvacationnowresponse"></a>
<a id="tocsv1startvacationnowresponse"></a>

```json
{
  "vacationId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|vacationId|string|false|none|none|

<h2 id="tocS_v1SubmitPatientIncidentRequest">v1SubmitPatientIncidentRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1submitpatientincidentrequest"></a>
<a id="schema_v1SubmitPatientIncidentRequest"></a>
<a id="tocSv1submitpatientincidentrequest"></a>
<a id="tocsv1submitpatientincidentrequest"></a>

```json
{
  "organizationId": "string",
  "categoryId": "string",
  "typeId": "string",
  "description": "string",
  "occurredAt": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|organizationId|string|true|none|none|
|categoryId|string|false|none|none|
|typeId|string|false|none|none|
|description|string|false|none|none|
|occurredAt|string|false|none|none|

<h2 id="tocS_v1SubmitPatientIncidentResponse">v1SubmitPatientIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1submitpatientincidentresponse"></a>
<a id="schema_v1SubmitPatientIncidentResponse"></a>
<a id="tocSv1submitpatientincidentresponse"></a>
<a id="tocsv1submitpatientincidentresponse"></a>

```json
{
  "bufferId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|bufferId|string|false|none|none|

<h2 id="tocS_v1SummaryPeriod">v1SummaryPeriod</h2>
<!-- backwards compatibility -->
<a id="schemav1summaryperiod"></a>
<a id="schema_v1SummaryPeriod"></a>
<a id="tocSv1summaryperiod"></a>
<a id="tocsv1summaryperiod"></a>

```json
{
  "organizationId": "string",
  "from": "string",
  "to": "string",
  "clinicId": "string",
  "departmentId": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|organizationId|string|false|none|none|
|from|string|false|none|none|
|to|string|false|none|none|
|clinicId|string|false|none|none|
|departmentId|string|false|none|none|

<h2 id="tocS_v1SystemAdminView">v1SystemAdminView</h2>
<!-- backwards compatibility -->
<a id="schemav1systemadminview"></a>
<a id="schema_v1SystemAdminView"></a>
<a id="tocSv1systemadminview"></a>
<a id="tocsv1systemadminview"></a>

```json
{
  "zitadelUserId": "string",
  "createdAt": "string"
}

```

SystemAdminView is the system-admin role; system admins are rooted in
Zitadel user ids, not employee ids.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|zitadelUserId|string|false|none|none|
|createdAt|string|false|none|none|

<h2 id="tocS_v1TerminateEmployeeResponse">v1TerminateEmployeeResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1terminateemployeeresponse"></a>
<a id="schema_v1TerminateEmployeeResponse"></a>
<a id="tocSv1terminateemployeeresponse"></a>
<a id="tocsv1terminateemployeeresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1TimeSeriesBucket">v1TimeSeriesBucket</h2>
<!-- backwards compatibility -->
<a id="schemav1timeseriesbucket"></a>
<a id="schema_v1TimeSeriesBucket"></a>
<a id="tocSv1timeseriesbucket"></a>
<a id="tocsv1timeseriesbucket"></a>

```json
{
  "bucketStart": "string",
  "bucketEnd": "string",
  "incidents": {
    "total": "string",
    "pending": "string",
    "inProgress": "string",
    "done": "string",
    "rejected": "string",
    "cancelled": "string",
    "highCritical": "string",
    "patientSource": "string",
    "reopened": "string"
  },
  "requests": {
    "total": "string",
    "completed": "string",
    "cancelled": "string",
    "linked": "string"
  }
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|bucketStart|string|false|none|none|
|bucketEnd|string|false|none|none|
|incidents|[v1TimeSeriesIncidentBucket](#schemav1timeseriesincidentbucket)|false|none|none|
|requests|[v1TimeSeriesRequestBucket](#schemav1timeseriesrequestbucket)|false|none|none|

<h2 id="tocS_v1TimeSeriesGranularity">v1TimeSeriesGranularity</h2>
<!-- backwards compatibility -->
<a id="schemav1timeseriesgranularity"></a>
<a id="schema_v1TimeSeriesGranularity"></a>
<a id="tocSv1timeseriesgranularity"></a>
<a id="tocsv1timeseriesgranularity"></a>

```json
"TIME_SERIES_GRANULARITY_UNSPECIFIED"

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|*anonymous*|string|false|none|none|

#### Enumerated Values

|Property|Value|
|---|---|
|*anonymous*|TIME_SERIES_GRANULARITY_UNSPECIFIED|
|*anonymous*|TIME_SERIES_GRANULARITY_DAY|
|*anonymous*|TIME_SERIES_GRANULARITY_WEEK|
|*anonymous*|TIME_SERIES_GRANULARITY_MONTH|

<h2 id="tocS_v1TimeSeriesIncidentBucket">v1TimeSeriesIncidentBucket</h2>
<!-- backwards compatibility -->
<a id="schemav1timeseriesincidentbucket"></a>
<a id="schema_v1TimeSeriesIncidentBucket"></a>
<a id="tocSv1timeseriesincidentbucket"></a>
<a id="tocsv1timeseriesincidentbucket"></a>

```json
{
  "total": "string",
  "pending": "string",
  "inProgress": "string",
  "done": "string",
  "rejected": "string",
  "cancelled": "string",
  "highCritical": "string",
  "patientSource": "string",
  "reopened": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|
|pending|string(int64)|false|none|none|
|inProgress|string(int64)|false|none|none|
|done|string(int64)|false|none|none|
|rejected|string(int64)|false|none|none|
|cancelled|string(int64)|false|none|none|
|highCritical|string(int64)|false|none|none|
|patientSource|string(int64)|false|none|none|
|reopened|string(int64)|false|none|none|

<h2 id="tocS_v1TimeSeriesRequestBucket">v1TimeSeriesRequestBucket</h2>
<!-- backwards compatibility -->
<a id="schemav1timeseriesrequestbucket"></a>
<a id="schema_v1TimeSeriesRequestBucket"></a>
<a id="tocSv1timeseriesrequestbucket"></a>
<a id="tocsv1timeseriesrequestbucket"></a>

```json
{
  "total": "string",
  "completed": "string",
  "cancelled": "string",
  "linked": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|total|string(int64)|false|none|none|
|completed|string(int64)|false|none|none|
|cancelled|string(int64)|false|none|none|
|linked|string(int64)|false|none|none|

<h2 id="tocS_v1TypeCount">v1TypeCount</h2>
<!-- backwards compatibility -->
<a id="schemav1typecount"></a>
<a id="schema_v1TypeCount"></a>
<a id="tocSv1typecount"></a>
<a id="tocsv1typecount"></a>

```json
{
  "typeId": "string",
  "typeName": "string",
  "count": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|typeId|string|false|none|none|
|typeName|string|false|none|none|
|count|string(int64)|false|none|none|

<h2 id="tocS_v1UnarchiveAnnouncementResponse">v1UnarchiveAnnouncementResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1unarchiveannouncementresponse"></a>
<a id="schema_v1UnarchiveAnnouncementResponse"></a>
<a id="tocSv1unarchiveannouncementresponse"></a>
<a id="tocsv1unarchiveannouncementresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateAnnouncementPriorityResponse">v1UpdateAnnouncementPriorityResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateannouncementpriorityresponse"></a>
<a id="schema_v1UpdateAnnouncementPriorityResponse"></a>
<a id="tocSv1updateannouncementpriorityresponse"></a>
<a id="tocsv1updateannouncementpriorityresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateAnnouncementResponse">v1UpdateAnnouncementResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateannouncementresponse"></a>
<a id="schema_v1UpdateAnnouncementResponse"></a>
<a id="tocSv1updateannouncementresponse"></a>
<a id="tocsv1updateannouncementresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateClinicDetailsResponse">v1UpdateClinicDetailsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateclinicdetailsresponse"></a>
<a id="schema_v1UpdateClinicDetailsResponse"></a>
<a id="tocSv1updateclinicdetailsresponse"></a>
<a id="tocsv1updateclinicdetailsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateClinicPhysicalAddressResponse">v1UpdateClinicPhysicalAddressResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateclinicphysicaladdressresponse"></a>
<a id="schema_v1UpdateClinicPhysicalAddressResponse"></a>
<a id="tocSv1updateclinicphysicaladdressresponse"></a>
<a id="tocsv1updateclinicphysicaladdressresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateDepartmentDetailsResponse">v1UpdateDepartmentDetailsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updatedepartmentdetailsresponse"></a>
<a id="schema_v1UpdateDepartmentDetailsResponse"></a>
<a id="tocSv1updatedepartmentdetailsresponse"></a>
<a id="tocsv1updatedepartmentdetailsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateEmployeeDepartmentResponse">v1UpdateEmployeeDepartmentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateemployeedepartmentresponse"></a>
<a id="schema_v1UpdateEmployeeDepartmentResponse"></a>
<a id="tocSv1updateemployeedepartmentresponse"></a>
<a id="tocsv1updateemployeedepartmentresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateEmployeePositionResponse">v1UpdateEmployeePositionResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateemployeepositionresponse"></a>
<a id="schema_v1UpdateEmployeePositionResponse"></a>
<a id="tocSv1updateemployeepositionresponse"></a>
<a id="tocsv1updateemployeepositionresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateIncidentCategoryDetailsResponse">v1UpdateIncidentCategoryDetailsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateincidentcategorydetailsresponse"></a>
<a id="schema_v1UpdateIncidentCategoryDetailsResponse"></a>
<a id="tocSv1updateincidentcategorydetailsresponse"></a>
<a id="tocsv1updateincidentcategorydetailsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateIncidentDescriptionResponse">v1UpdateIncidentDescriptionResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateincidentdescriptionresponse"></a>
<a id="schema_v1UpdateIncidentDescriptionResponse"></a>
<a id="tocSv1updateincidentdescriptionresponse"></a>
<a id="tocsv1updateincidentdescriptionresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateIncidentPriorityResponse">v1UpdateIncidentPriorityResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateincidentpriorityresponse"></a>
<a id="schema_v1UpdateIncidentPriorityResponse"></a>
<a id="tocSv1updateincidentpriorityresponse"></a>
<a id="tocsv1updateincidentpriorityresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateIncidentStatusResponse">v1UpdateIncidentStatusResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateincidentstatusresponse"></a>
<a id="schema_v1UpdateIncidentStatusResponse"></a>
<a id="tocSv1updateincidentstatusresponse"></a>
<a id="tocsv1updateincidentstatusresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateIncidentTypeDetailsResponse">v1UpdateIncidentTypeDetailsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateincidenttypedetailsresponse"></a>
<a id="schema_v1UpdateIncidentTypeDetailsResponse"></a>
<a id="tocSv1updateincidenttypedetailsresponse"></a>
<a id="tocsv1updateincidenttypedetailsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateOrganizationDetailsResponse">v1UpdateOrganizationDetailsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateorganizationdetailsresponse"></a>
<a id="schema_v1UpdateOrganizationDetailsResponse"></a>
<a id="tocSv1updateorganizationdetailsresponse"></a>
<a id="tocsv1updateorganizationdetailsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateOrganizationLegalAddressResponse">v1UpdateOrganizationLegalAddressResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateorganizationlegaladdressresponse"></a>
<a id="schema_v1UpdateOrganizationLegalAddressResponse"></a>
<a id="tocSv1updateorganizationlegaladdressresponse"></a>
<a id="tocsv1updateorganizationlegaladdressresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdatePatientIncidentResponse">v1UpdatePatientIncidentResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updatepatientincidentresponse"></a>
<a id="schema_v1UpdatePatientIncidentResponse"></a>
<a id="tocSv1updatepatientincidentresponse"></a>
<a id="tocsv1updatepatientincidentresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateRequestTypeDetailsResponse">v1UpdateRequestTypeDetailsResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updaterequesttypedetailsresponse"></a>
<a id="schema_v1UpdateRequestTypeDetailsResponse"></a>
<a id="tocSv1updaterequesttypedetailsresponse"></a>
<a id="tocsv1updaterequesttypedetailsresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateServiceRequestDescriptionResponse">v1UpdateServiceRequestDescriptionResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateservicerequestdescriptionresponse"></a>
<a id="schema_v1UpdateServiceRequestDescriptionResponse"></a>
<a id="tocSv1updateservicerequestdescriptionresponse"></a>
<a id="tocsv1updateservicerequestdescriptionresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateServiceRequestStatusResponse">v1UpdateServiceRequestStatusResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updateservicerequeststatusresponse"></a>
<a id="schema_v1UpdateServiceRequestStatusResponse"></a>
<a id="tocSv1updateservicerequeststatusresponse"></a>
<a id="tocsv1updateservicerequeststatusresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1UpdateVacationEndDateResponse">v1UpdateVacationEndDateResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1updatevacationenddateresponse"></a>
<a id="schema_v1UpdateVacationEndDateResponse"></a>
<a id="tocSv1updatevacationenddateresponse"></a>
<a id="tocsv1updatevacationenddateresponse"></a>

```json
{}

```

### Properties

*None*

<h2 id="tocS_v1VacationView">v1VacationView</h2>
<!-- backwards compatibility -->
<a id="schemav1vacationview"></a>
<a id="schema_v1VacationView"></a>
<a id="tocSv1vacationview"></a>
<a id="tocsv1vacationview"></a>

```json
{
  "id": "string",
  "employeeId": "string",
  "state": "string",
  "startsAt": "string",
  "endsAt": "string",
  "createdAt": "string",
  "updatedAt": "string"
}

```

VacationView mirrors projections.employee_vacations. state is one of
{scheduled, active, ended, cancelled}.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string|false|none|none|
|employeeId|string|false|none|none|
|state|string|false|none|none|
|startsAt|string|false|none|none|
|endsAt|string|false|none|none|
|createdAt|string|false|none|none|
|updatedAt|string|false|none|none|

<h2 id="tocS_v1ValidationFailedDetails">v1ValidationFailedDetails</h2>
<!-- backwards compatibility -->
<a id="schemav1validationfaileddetails"></a>
<a id="schema_v1ValidationFailedDetails"></a>
<a id="tocSv1validationfaileddetails"></a>
<a id="tocsv1validationfaileddetails"></a>

```json
{
  "violations": [
    {
      "field": "string",
      "rule": "string",
      "message": "string",
      "param": "string"
    }
  ]
}

```

ValidationFailedDetails is present only when code = "validation_failed".
Each violation corresponds to one struct-tag rule failure or one
domain-level leaf error from errors.Join.

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|violations|[[ValidationFailedDetailsFieldViolation](#schemavalidationfaileddetailsfieldviolation)]|false|none|none|

<h2 id="tocS_v1ValidationErrorResponse">v1ValidationErrorResponse</h2>
<!-- backwards compatibility -->
<a id="schemav1validationerrorresponse"></a>
<a id="schema_v1ValidationErrorResponse"></a>
<a id="tocSv1validationerrorresponse"></a>
<a id="tocsv1validationerrorresponse"></a>

```json
{
  "code": "string",
  "message": "string",
  "details": {
    "violations": [
      {
        "field": "string",
        "rule": "string",
        "message": "string",
        "param": "string"
      }
    ]
  }
}

```

ValidationErrorResponse

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|string|false|none|none|
|message|string|false|none|none|
|details|[v1ValidationFailedDetails](#schemav1validationfaileddetails)|false|none|ValidationFailedDetails is present only when code = "validation_failed".<br>Each violation corresponds to one struct-tag rule failure or one<br>domain-level leaf error from errors.Join.|

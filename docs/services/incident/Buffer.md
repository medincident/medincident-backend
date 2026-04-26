← [Документация](../../README.md)

# Буфер пациента

Промежуточный этап перед созданием инцидента. Пациент подаёт заявку; диспетчер (`OrgDispatcher`) обрабатывает её — публикует (создаёт инцидент), отклоняет или пациент отменяет сам.

**gRPC:** `IncidentBufferCommandService` (command) / `IncidentQueryService` (query, буфер)

Статусная машина: [Статусные машины](../../architecture/Status-Machines.md)

---

## SubmitPatientIncident

**HTTP:** `POST /v1/organizations/{organization_id}/patient-incidents`
**gRPC:** `IncidentBufferCommandService.SubmitPatientIncident`

### Права доступа

`Authenticated` — любой аутентифицированный пользователь может подать заявку.

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `incident_type_id` | string (UUID) | required, uuid |
| `description` | string | required, min=1, max=4096 |

### Инварианты

- Тип инцидента должен быть активным и иметь флаг `is_allowed_for_patients=true`.
- Начальный статус заявки — `pending`.

### Ошибки

| Код | Описание |
|---|---|
| `incident_type_not_found` | Тип не найден, не активен или не доступен для пациентов |

---

## UpdatePatientIncident

**HTTP:** `PATCH /v1/patient-incidents/{patient_incident_id}`
**gRPC:** `IncidentBufferCommandService.UpdatePatientIncident`

### Права доступа

Пациент, подавший заявку (по `zitadel_user_id` из токена), или `OrgDispatcherOf.Organization`.

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `patient_incident_id` | string (UUID) | required, uuid |
| `description` | string | omitempty, min=1, max=4096 |
| `incident_type_id` | string (UUID) | omitempty, uuid |

### Инварианты

- Доступно только для заявок в статусе `pending`.

---

## CancelPatientIncident

**HTTP:** `POST /v1/patient-incidents/{patient_incident_id}/cancel`
**gRPC:** `IncidentBufferCommandService.CancelPatientIncident`

### Права доступа

Пациент, подавший заявку.

### Инварианты

- Доступно только для заявок в статусе `pending`.

---

## PublishPatientIncident

**HTTP:** `POST /v1/patient-incidents/{patient_incident_id}/publish`
**gRPC:** `IncidentBufferCommandService.PublishPatientIncident`

### Права доступа

`OrgDispatcherOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `patient_incident_id` | string (UUID) | required, uuid |
| `priority` | enum | required |

### Инварианты

- Публикация атомарно переводит заявку в `published` и создаёт инцидент (`CreateIncident` с `source_buffer_id`).
- `incident.source_buffer_id` ссылается на опубликованную заявку.

---

## RejectPatientIncident

**HTTP:** `POST /v1/patient-incidents/{patient_incident_id}/reject`
**gRPC:** `IncidentBufferCommandService.RejectPatientIncident`

### Права доступа

`OrgDispatcherOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `patient_incident_id` | string (UUID) | required, uuid |
| `reason` | string | required, min=1, max=1024 |

### Инварианты

- Доступно только для заявок в статусе `pending`.

---

## Query-методы

| Метод | Права |
|---|---|
| `ListPatientIncidents` | OrgDispatcherOf или пациент-владелец |
| `GetPatientIncident` | OrgDispatcherOf или пациент-владелец |

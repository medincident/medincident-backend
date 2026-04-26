← [Документация](../../README.md)

# Инциденты

Управление инцидентами — основная рабочая сущность системы.

**gRPC:** `IncidentCommandService` (command) / `IncidentQueryService` (query)

Статусная машина: [Статусные машины](../../architecture/Status-Machines.md)

---

## CreateIncident

**HTTP:** `POST /v1/organizations/{organization_id}/incidents`
**gRPC:** `IncidentCommandService.CreateIncident`

### Права доступа

`AdminOf.Organization(organizationID)` или `OrgDispatcherOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `incident_type_id` | string (UUID) | required, uuid |
| `description` | string | required, min=1, max=4096 |
| `priority` | enum | required |
| `source_buffer_id` | string (UUID) | omitempty, uuid |

### Инварианты

- Тип инцидента должен принадлежать той же организации и быть активным.
- При создании из буфера (`source_buffer_id`) запись буфера переходит в статус `published`.
- Начальный статус инцидента — `open`.

### Ошибки

| Код | Описание |
|---|---|
| `incident_type_not_found` | Тип инцидента не найден или не активен |
| `patient_incident_not_found` | Запись буфера не найдена |
| `patient_incident_invalid_status` | Запись буфера не в статусе `pending` |

---

## UpdateIncidentDescription

**HTTP:** `PATCH /v1/incidents/{incident_id}/description`
**gRPC:** `IncidentCommandService.UpdateIncidentDescription`

### Права доступа

`AdminOf.Organization(organizationID)` (через incident → organization)

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `incident_id` | string (UUID) | required, uuid |
| `description` | string | required, min=1, max=4096 |

### Инварианты

- Доступно только для инцидентов в статусе `open`.

---

## UpdateIncidentStatus

**HTTP:** `PUT /v1/incidents/{incident_id}/status`
**gRPC:** `IncidentCommandService.UpdateIncidentStatus`

### Права доступа

`SystemAdmin` или `OrgAdminOf.Organization` или `OrgHeadOf.Organization` или `OrgDispatcherOf.Organization` или `ClinicHeadOf.Clinic` или `DeptResponsibleOf.Department`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `incident_id` | string (UUID) | required, uuid |
| `new_status` | enum | required, oneof=in_progress done rejected |

### Инварианты

- Допустимые переходы: `pending → in_progress`, `in_progress → done`, `in_progress → rejected`.
- Отмена (`cancelled`) выполняется через отдельный метод `CancelIncident`.

### Ошибки

| Код | Описание |
|---|---|
| `incident_not_found` | Инцидент не найден |
| `incident_invalid_status_flow` | Недопустимый переход статуса |

---

## UpdateIncidentPriority

**HTTP:** `PUT /v1/incidents/{incident_id}/priority`
**gRPC:** `IncidentCommandService.UpdateIncidentPriority`

### Права доступа

`SystemAdmin` или `OrgAdminOf.Organization` или `OrgHeadOf.Organization` или `OrgDispatcherOf.Organization` или `ClinicHeadOf.Clinic` или `DeptResponsibleOf.Department`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `incident_id` | string (UUID) | required, uuid |
| `priority` | enum | required, oneof=low normal high critical |

### Инварианты

- Недоступно для инцидентов в терминальном статусе (`done`, `rejected`, `cancelled`).
- Если приоритет не изменился, строка истории не записывается.

### Ошибки

| Код | Описание |
|---|---|
| `incident_not_found` | Инцидент не найден |
| `incident_frozen` | Инцидент находится в терминальном статусе |

---

## CancelIncident

**HTTP:** `POST /v1/incidents/{incident_id}/cancel`
**gRPC:** `IncidentCommandService.CancelIncident`

### Права доступа

`AdminOf.Organization(organizationID)`

### Инварианты

- Доступно только для инцидентов в статусе `open`.

---

## ReopenIncident

**HTTP:** `POST /v1/incidents/{incident_id}/reopen`
**gRPC:** `IncidentCommandService.ReopenIncident`

### Права доступа

`AdminOf.Organization(organizationID)`

### Инварианты

- Доступно только для инцидентов в статусе `closed`.

---

## Query-методы

| Метод | Права |
|---|---|
| `ListIncidents` | ReaderOf.Organization |
| `GetIncident` | ReaderOf.Organization |
| `ListIncidentStatusHistory` | ReaderOf.Organization |
| `ListIncidentPriorityHistory` | ReaderOf.Organization |

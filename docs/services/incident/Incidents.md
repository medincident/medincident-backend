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

**HTTP:** `PATCH /v1/incidents/{incident_id}/status`
**gRPC:** `IncidentCommandService.UpdateIncidentStatus`

### Права доступа

`AdminOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `incident_id` | string (UUID) | required, uuid |
| `status` | enum | required |
| `priority` | enum | omitempty |

### Инварианты

- Допустимые переходы статусов: [Статусные машины](../../architecture/Status-Machines.md).
- `cancelled` — терминальный статус, переход из него запрещён.

### Ошибки

| Код | Описание |
|---|---|
| `incident_not_found` | Инцидент не найден |
| `incident_invalid_status_transition` | Недопустимый переход статуса |

---

## UpdateIncidentPriority

**HTTP:** `PUT /v1/incidents/{incident_id}/priority`
**gRPC:** `IncidentCommandService.UpdateIncidentPriority`

### Права доступа

Привилегированные роли: `AdminOf.Organization(organizationID)`,
`OrgHeadOf.Organization(organizationID)`, `OrgDispatcherOf.Organization(organizationID)`,
`ClinicHeadOf.Clinic(clinicID)`, `DeptResponsibleOf.Department(departmentID)`.

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `incident_id` | string (UUID) | required, uuid |
| `priority` | enum (`low`, `normal`, `high`, `critical`) | required |

### Инварианты

- Инцидент не должен быть в терминальном статусе (`done`, `rejected`, `cancelled`).
- Если новый приоритет совпадает с текущим — операция пропускается без ошибки (no-op), запись в историю не создаётся.

### Ошибки

| Код | Описание |
|---|---|
| `incident_not_found` | Инцидент не найден |
| `incident_frozen` | Инцидент в терминальном статусе и не может быть изменён |

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

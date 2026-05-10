← [Документация](../../README.md)

# Инциденты

Управление инцидентами — основная рабочая сущность системы.

**gRPC:** `IncidentCommandService` (command) / `IncidentQueryService` (query)

Статусная машина: [Статусные машины](../../architecture/Status-Machines.md)

---

## CreateIncident

**HTTP:** `POST /v1/incidents`
**gRPC:** `IncidentCommandService.CreateIncident`

### Права доступа

`AdminOf.Organization(organizationID)` или `OrgDispatcherOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `department_id` | string (UUID) | required, uuid |
| `category_id` | string (UUID) | required, uuid |
| `type_id` | string (UUID) | required, uuid |
| `description` | string | omitempty, max=4096 |
| `occurred_at` | string (RFC3339Nano) | required |

### Инварианты

- Тип инцидента должен принадлежать той же организации и быть активным.
- При создании через `PublishPatientIncident` бэкенд автоматически связывает инцидент с записью буфера через внутреннее поле `source_buffer_id` и переводит её в статус `published` — клиент это поле не передаёт.
- Начальный статус инцидента — `pending`.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `incident_category_inactive` | 400 | Категория деактивирована |
| `incident_type_inactive` | 400 | Тип инцидента деактивирован |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_department_not_found` | 404 | Отдел не найден |
| `incident_category_not_found` | 404 | Категория не найдена |
| `incident_type_not_found` | 404 | Тип инцидента не найден |
| `incident_employee_not_found` | 404 | Регистратор (сотрудник) не найден |
| `incident_registrar_user_not_found` | 404 | Проекция пользователя регистратора не найдена |

---

## CancelIncident

**HTTP:** `POST /v1/incidents/{incident_id}:cancel`
**gRPC:** `IncidentCommandService.CancelIncident`

### Права доступа

`AdminOf.Organization(organizationID)`

### Инварианты

- Доступно только для инцидентов в статусе `pending` или `in_progress`.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `incident_not_cancellable` | 400 | Инцидент нельзя отменить в текущем статусе |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_not_found` | 404 | Инцидент не найден |

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

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `incident_invalid_status_transition` | 400 | Недопустимый переход статуса |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_not_found` | 404 | Инцидент не найден |

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

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `incident_frozen` | 400 | Инцидент находится в терминальном статусе |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_not_found` | 404 | Инцидент не найден |

---

## UpdateIncidentDescription

**HTTP:** `PUT /v1/incidents/{incident_id}/description`
**gRPC:** `IncidentCommandService.UpdateIncidentDescription`

### Права доступа

`AdminOf.Organization(organizationID)` (через incident → organization)

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `incident_id` | string (UUID) | required, uuid |
| `description` | string | omitempty, max=4096 |

### Инварианты

- Доступно только для инцидентов в статусе `pending` или `in_progress`.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_not_found` | 404 | Инцидент не найден |

---

## ReopenIncident

**HTTP:** `POST /v1/incidents/{incident_id}:reopen`
**gRPC:** `IncidentCommandService.ReopenIncident`

### Права доступа

`AdminOf.Organization(organizationID)`

### Инварианты

- Доступно только для инцидентов в терминальном статусе (`done`, `rejected`, `cancelled`).

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `incident_not_reopenable` | 400 | Инцидент нельзя переоткрыть в текущем статусе |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_not_found` | 404 | Инцидент не найден |

---

## Query-методы

| Метод | Права |
|---|---|
| `ListIncidents` | ReaderOf.Organization |
| `GetIncident` | ReaderOf.Organization |
| `ListMyIncidents` | Authenticated (только свои) |
| `GetIncidentHistory` | ReaderOf.Organization |

### GetIncident

**HTTP:** `GET /v1/query/incidents/{id}`
**gRPC:** `IncidentQueryService.GetIncident`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `incident_query_not_found` | 404 | Инцидент не найден |

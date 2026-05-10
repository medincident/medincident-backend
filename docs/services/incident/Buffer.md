← [Документация](../../README.md)

# Буфер пациента

Промежуточный этап перед созданием инцидента. Пациент подаёт заявку; диспетчер (`OrgDispatcher`) обрабатывает её — публикует (создаёт инцидент), отклоняет или пациент отменяет сам.

**gRPC:** `IncidentBufferCommandService` (command) / `IncidentQueryService` (query, буфер)

Статусная машина: [Статусные машины](../../architecture/Status-Machines.md)

---

## SubmitPatientIncident

**HTTP:** `POST /v1/patient-incidents`
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

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `buffer_organization_not_found` | 404 | Организация не найдена |

---

## UpdatePatientIncident

**HTTP:** `PUT /v1/patient-incidents/{buffer_id}`
**gRPC:** `IncidentBufferCommandService.UpdatePatientIncident`

### Права доступа

Пациент, подавший заявку (по `zitadel_user_id` из токена), или `OrgDispatcherOf.Organization`.

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `buffer_id` | string (UUID) | required, uuid |
| `description` | string | omitempty, min=1, max=4096 |
| `incident_type_id` | string (UUID) | omitempty, uuid |

### Инварианты

- Доступно только для заявок в статусе `pending`.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `buffer_not_pending` | 400 | Заявка не в статусе pending |
| `permission_denied` / `buffer_not_patient_owner` | 403 | Нет прав доступа |

---

## CancelPatientIncident

**HTTP:** `POST /v1/patient-incidents/{buffer_id}:cancel`
**gRPC:** `IncidentBufferCommandService.CancelPatientIncident`

### Права доступа

Пациент, подавший заявку.

### Инварианты

- Доступно только для заявок в статусе `pending`.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |

---

## PublishPatientIncident

**HTTP:** `POST /v1/patient-incidents/{buffer_id}:publish`
**gRPC:** `IncidentBufferCommandService.PublishPatientIncident`

### Права доступа

`OrgDispatcherOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `buffer_id` | string (UUID) | required, uuid |
| `department_id` | string (UUID) | required, uuid |
| `category_id` | string (UUID) | required, uuid |
| `type_id` | string (UUID) | required, uuid |
| `description` | string | omitempty |

### Инварианты

- Публикация атомарно переводит заявку в `published` и создаёт инцидент (`CreateIncident` с `source_buffer_id`).
- `incident.source_buffer_id` ссылается на опубликованную заявку.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `buffer_department_not_found` | 404 | Отдел не найден |
| `buffer_category_not_found` | 404 | Категория не найдена |
| `buffer_type_not_found` | 404 | Тип инцидента не найден |
| `buffer_dispatcher_not_found` | 404 | Диспетчер не найден |

---

## RejectPatientIncident

**HTTP:** `POST /v1/patient-incidents/{buffer_id}:reject`
**gRPC:** `IncidentBufferCommandService.RejectPatientIncident`

### Права доступа

`OrgDispatcherOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `buffer_id` | string (UUID) | required, uuid |

### Инварианты

- Доступно только для заявок в статусе `pending`.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |

---

## Query-методы

| Метод | Права |
|---|---|
| `ListBufferEntries` | OrgDispatcherOf или пациент-владелец |
| `GetBufferEntry` | OrgDispatcherOf или пациент-владелец |
| `ListMyBufferEntries` | Authenticated (только свои) |

### GetBufferEntry

**HTTP:** `GET /v1/query/patient-incidents/{id}`
**gRPC:** `IncidentQueryService.GetBufferEntry`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `buffer_query_not_found` | 404 | Запись буфера не найдена |

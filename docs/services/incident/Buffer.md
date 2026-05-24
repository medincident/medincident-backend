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
| `category_id` | string (UUID) | omitempty, uuid |
| `type_id` | string (UUID) | omitempty, uuid |
| `description` | string | required, min=1, max=10000 |
| `summary` | string | required, min=1, max=10000 — AI-формализованный текст обращения |
| `priority` | BufferPriority (`BUFFER_PRIORITY_NORMAL` \| `BUFFER_PRIORITY_HIGH`) | required — приоритет, выставляемый AI-сервисом |
| `occurred_at` | string (RFC3339Nano) | omitempty |

### Инварианты

- Тип инцидента должен иметь флаг `is_allowed_for_patients=true`.
- Начальный статус заявки — `pending`.
- `description` — исходный текст обращения пациента (обязателен).
- `summary` — формализованный вариант текста, генерируется AI-сервисом (обязателен).
- `priority` — оценка AI-сервиса; диспетчер видит её, но устанавливает приоритет НС самостоятельно через `UpdatePriority`.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `buffer_occurred_at_invalid` | 400 | Некорректный формат даты occurred_at |
| `buffer_type_not_allowed_for_patients` | 422 | Тип инцидента недоступен для пациентов |
| `buffer_organization_not_found` | 404 | Организация не найдена |
| `buffer_category_not_found` | 404 | Категория не найдена |
| `buffer_type_not_found` | 404 | Тип инцидента не найден |

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
| `category_id` | string (UUID) | omitempty, uuid |
| `type_id` | string (UUID) | omitempty, uuid |
| `description` | string | omitempty, min=1, max=10000 |
| `summary` | string | omitempty, min=1, max=10000 — обновлённый AI-текст |
| `priority` | BufferPriority (`BUFFER_PRIORITY_NORMAL` \| `BUFFER_PRIORITY_HIGH`) | omitempty — обновлённый приоритет AI |
| `occurred_at` | string (RFC3339Nano) | omitempty |

### Инварианты

- Доступно только для заявок в статусе `pending`.
- Тип инцидента должен иметь флаг `is_allowed_for_patients=true`.
- Все поля независимо обновляемы; `nil` оставляет текущее значение без изменений.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `buffer_occurred_at_invalid` | 400 | Некорректный формат даты occurred_at |
| `buffer_not_pending` | 422 | Заявка не в статусе pending |
| `buffer_type_not_allowed_for_patients` | 422 | Тип инцидента недоступен для пациентов |
| `permission_denied` / `buffer_not_patient_owner` | 403 | Нет прав доступа |
| `buffer_not_found` | 404 | Заявка не найдена |
| `buffer_category_not_found` | 404 | Категория не найдена |
| `buffer_type_not_found` | 404 | Тип инцидента не найден |

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
| `buffer_not_pending` | 422 | Заявка не в статусе pending |
| `permission_denied` | 403 | Нет прав доступа |
| `buffer_not_found` | 404 | Заявка не найдена |
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

## BufferEntryView — поля ответа

| Поле | Тип | Описание |
|---|---|---|
| `id` | string (UUID) | Идентификатор буфера |
| `organization_id` | string (UUID) | Организация |
| `patient_zitadel_user_id` | string | Zitadel ID пациента |
| `category_id` | string (UUID)? | Категория |
| `type_id` | string (UUID)? | Тип |
| `description` | string | Исходный текст обращения (всегда непустой) |
| `summary` | string | AI-формализованный текст (всегда непустой) |
| `priority` | string (`normal` \| `high`) | AI-приоритет |
| `occurred_at` | string (RFC3339Nano)? | Время инцидента |
| `status` | BufferStatus | Статус заявки |
| `published_incident_id` | string (UUID)? | ID созданного инцидента (при `published`) |
| `created_at` | string (RFC3339Nano) | Время создания |
| `updated_at` | string (RFC3339Nano) | Время обновления |
| `patient_status` | PatientStatus? | Только для пациента-владельца |

## Данные буфера в ответе инцидента

Если инцидент создан из буфера (`source_buffer_id` заполнен), в ответе `GetIncident` / `ListIncidents` / `ListMyIncidents` присутствует поле `patient_buffer`:

| Поле | Тип | Описание |
|---|---|---|
| `buffer_id` | string (UUID) | ID буфера-источника |
| `description` | string | Исходный текст обращения |
| `summary` | string | AI-формализованный текст |
| `priority` | string (`normal` \| `high`) | AI-приоритет из буфера |

Диспетчер видит `patient_buffer.priority` как AI-оценку, но устанавливает приоритет НС самостоятельно через `UpdatePriority`.

---

## Query-методы

| Метод | Права |
|---|---|
| `ListBufferEntries` | OrgDispatcherOf или пациент-владелец |
| `GetBufferEntry` | OrgDispatcherOf или пациент-владелец |
| `ListMyBufferEntries` | Authenticated (только свои) |

### ListBufferEntries

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `incident_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |
| `permission_denied` | 403 | Нет прав доступа |

---

### ListMyBufferEntries

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `incident_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |

---

### GetBufferEntry

**HTTP:** `GET /v1/query/patient-incidents/{id}`
**gRPC:** `IncidentQueryService.GetBufferEntry`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `buffer_query_not_found` | 404 | Запись буфера не найдена |

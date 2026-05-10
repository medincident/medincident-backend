← [Документация](../../README.md)

# Аналитика инцидентов

Сервис предоставляет агрегированную статистику по инцидентам, сервисным заявкам
и буферу пациентов в рамках заданного временного диапазона. Данные используются
для построения дашбордов и трендовых графиков на стороне клиента.

**gRPC:** `AnalyticsQueryService` (query)

---

## GetSnapshot

**HTTP:** `GET /v1/analytics/snapshot`
**gRPC:** `AnalyticsQueryService.GetSnapshot`

Возвращает необработанные аналитические записи по каждому инциденту, сервисной
заявке и (опционально) записи буфера пациентов за указанный период. Предназначен
для построения графиков на стороне клиента.

Максимальный допустимый период — **366 дней**. Запросы с большим диапазоном
отклоняются с кодом `analytics_period_too_large`.

### Параметры

> HTTP-параметры используют camelCase (grpc-gateway): `organizationId`, `from`, `to`, `clinicId`, `departmentId`, `includePatientBuffer`. gRPC-поля — snake_case.

| Поле | HTTP-параметр | Тип | Правила |
|---|---|---|---|
| `organization_id` | `organizationId` | string (UUID) | обязательный |
| `from` | `from` | string (RFC3339) | обязательный; должен быть меньше `to` |
| `to` | `to` | string (RFC3339) | обязательный |
| `clinic_id` | `clinicId` | string (UUID) | необязательный; фильтр по клинике |
| `department_id` | `departmentId` | string (UUID) | необязательный; фильтр по отделению |
| `include_patient_buffer` | `includePatientBuffer` | bool | необязательный; по умолчанию `false` |

### Ответ

| Поле | Описание |
|---|---|
| `incidents` | Список записей `SnapshotIncident` с полями: статус, приоритет, категория, тип, клиника, отделение, признак источника-пациента, признак повторного открытия, количество связанных заявок |
| `requests` | Список записей `SnapshotRequest` с полями: статус, тип, отделение, признак связи с инцидентом |
| `patient_buffer` | Список записей `SnapshotPatientBuffer` (только при `include_patient_buffer = true`): статус, категория |

---

## GetSummary

**HTTP:** `GET /v1/analytics/summary`
**gRPC:** `AnalyticsQueryService.GetSummary`

Возвращает агрегированные KPI-показатели за заданный период. Ограничения на
длину периода отсутствуют.

### Параметры

> HTTP-параметры используют camelCase (grpc-gateway): `organizationId`, `from`, `to`, `clinicId`, `departmentId`. gRPC-поля — snake_case.

| Поле | HTTP-параметр | Тип | Правила |
|---|---|---|---|
| `organization_id` | `organizationId` | string (UUID) | обязательный |
| `from` | `from` | string (RFC3339) | обязательный; должен быть меньше `to` |
| `to` | `to` | string (RFC3339) | обязательный |
| `clinic_id` | `clinicId` | string (UUID) | необязательный; фильтр по клинике |
| `department_id` | `departmentId` | string (UUID) | необязательный; фильтр по отделению |

### Ответ

**Инциденты:**

| Поле | Описание |
|---|---|
| Счётчики по статусам | `pending`, `in_progress`, `done`, `rejected`, `cancelled` |
| Счётчики по приоритетам | `low`, `normal`, `high`, `critical` |
| Счётчики по источнику | `patient_source` (источник — пациент), `staff_source` (источник — персонал) |
| `reopened` | Количество повторно открытых инцидентов |
| `with_linked_requests` | Количество инцидентов со связанными заявками |
| `resolution_avg/min/max` | Среднее/минимальное/максимальное время закрытия (минуты) |
| Перцентили | P50 / P90 / P95 времени закрытия (минуты) |
| Топ-10 категорий | Категории с наибольшим числом инцидентов |
| Топ-10 типов | Типы инцидентов с наибольшим числом инцидентов |
| Топ-10 отделений | Отделения с наибольшим числом инцидентов |

**Сервисные заявки:**

| Поле | Описание |
|---|---|
| Счётчики по статусам | `created`, `in_work`, `on_hold`, `pending_review`, `completed`, `cancelled` |
| `linked` / `unlinked` | Заявки со связанным инцидентом / без него |
| `completion_avg/min/max` | Среднее/минимальное/максимальное время выполнения (минуты) |
| Перцентили | P50 / P90 / P95 времени выполнения (минуты) |
| Топ-10 типов | Типы заявок с наибольшим числом |
| Топ-10 отделений | Отделения с наибольшим числом заявок |

**Буфер пациентов:**

| Поле | Описание |
|---|---|
| Счётчики по статусам | `pending`, `published`, `rejected`, `cancelled` |

---

## GetTimeSeries

**HTTP:** `GET /v1/analytics/timeseries`
**gRPC:** `AnalyticsQueryService.GetTimeSeries`

Возвращает счётчики по временным корзинам (бакетам) для построения трендовых
графиков. Все бакеты заполняются нулями, что обеспечивает равномерные ряды
одинаковой длины вне зависимости от наличия данных. Ограничений на длину периода
нет.

### Параметры

> HTTP-параметры используют camelCase (grpc-gateway): `organizationId`, `from`, `to`, `clinicId`, `departmentId`, `granularity`. gRPC-поля — snake_case.

Максимальное количество бакетов — **1000**. Запросы, превышающие это значение (например, 2 года с `DAY`), отклоняются с кодом `analytics_period_too_large`.

| Поле | HTTP-параметр | Тип | Правила |
|---|---|---|---|
| `organization_id` | `organizationId` | string (UUID) | обязательный |
| `from` | `from` | string (RFC3339) | обязательный; должен быть меньше `to` |
| `to` | `to` | string (RFC3339) | обязательный |
| `clinic_id` | `clinicId` | string (UUID) | необязательный; фильтр по клинике |
| `department_id` | `departmentId` | string (UUID) | необязательный; фильтр по отделению |
| `granularity` | `granularity` | enum | `DAY` / `WEEK` / `MONTH` |

### Ответ

Список бакетов `TimeSeriesBucket`. Каждый бакет содержит:

| Поле | Описание |
|---|---|
| `bucket_start` / `bucket_end` | Временные границы бакета |
| `incident_total` | Всего инцидентов |
| `i_pending` / `i_in_progress` / `i_done` / `i_rejected` / `i_cancelled` | По статусу |
| `i_high_critical` | Инциденты с приоритетом `high` или `critical` |
| `i_patient_source` | Инциденты от пациентов |
| `i_reopened` | Повторно открытые инциденты |
| `req_total` / `req_completed` / `req_cancelled` / `req_linked` | Счётчики заявок |

---

## Разграничение прав доступа

| Запрашиваемый скоуп | Допустимые роли |
|---|---|
| Орг-уровень (без `clinic_id` и `department_id`) | SystemAdmin · OrgAdmin · OrgHead · OrgDispatcher |
| Клиника (с `clinic_id`) | SystemAdmin · OrgAdmin · OrgHead · OrgDispatcher · ClinicHead этой клиники |
| Отделение (с `department_id`) | SystemAdmin · OrgAdmin · OrgHead · OrgDispatcher · ClinicHead клиники отделения (если передан `clinic_id`) · DeptResponsible этого отделения |

---

## Мультитенантная защита

Применяется двухуровневая защита от межтенантных запросов:

1. **Проверка роли** — `requireAuth` гарантирует, что вызывающая сторона имеет
   необходимую роль в рамках запрашиваемого скоупа (организация / клиника /
   отделение).

2. **Валидация принадлежности** — после успешной проверки роли `validateScope`
   верифицирует, что `clinic_id` принадлежит переданному `organization_id`, а
   `department_id` принадлежит той же организации (через связь
   `departments → clinics → organization_id`). Несоответствие возвращает
   `analytics_clinic_not_found` или `analytics_dept_not_found` соответственно —
   без различия между «не существует» и «не в этой организации».

---

## Источник данных

Все три метода обращаются напрямую к таблицам `domain.*`, а не к
`projections.*`. Причина: проекционные таблицы не содержат поля `priority`,
`source_patient_zitadel_user_id` и `reopened_from_incident_id`, необходимые
для полноценной аналитики. Использование `domain.*` является архитектурно
намеренным решением, а не временным компромиссом.

---

## Коды ошибок

Общие для всех трёх методов (`GetSnapshot`, `GetSummary`, `GetTimeSeries`):

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Недостаточно прав |
| `analytics_org_not_found` | 404 | Некорректный или отсутствующий `organization_id` |
| `analytics_clinic_not_found` | 404 | `clinic_id` не принадлежит организации |
| `analytics_dept_not_found` | 404 | `department_id` не принадлежит организации |
| `analytics_invalid_scope` | 400 | `clinic_id` или `department_id` не является валидным UUID |
| `analytics_period_invalid` | 400 | Некорректный диапазон дат (`from >= to` или неверный формат RFC3339) |
| `analytics_period_too_large` | 400 | Период превышает 366 дней (`GetSnapshot`) или 1000 бакетов (`GetTimeSeries`) |
| `analytics_query_failed` | 500 | Ошибка выполнения SQL-запроса |

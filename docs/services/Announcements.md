# Объявления (Announcements)

Объявления — односторонняя информационная рассылка от администрации сотрудникам. Могут быть адресованы на уровне организации, клиники или отделения.

---

## Модель данных

### `domain.announcements`

| Поле | Тип | Описание |
|---|---|---|
| `id` | UUID PK (uuidv7) | Идентификатор |
| `organization_id` | UUID NOT NULL FK → `domain.organizations` | Организация (всегда заполнено, CASCADE DELETE) |
| `clinic_id` | UUID FK → `domain.clinics` | NULL = объявление уровня организации |
| `department_id` | UUID FK → `domain.departments` | NULL = объявление уровня орг или клиники |
| `author_id` | TEXT NOT NULL | Zitadel user ID автора (берётся из JWT, фронт не передаёт) |
| `title` | TEXT NOT NULL | Заголовок (min=3, max=200) |
| `content` | TEXT NOT NULL | Текст (min=10, max=5000) |
| `priority` | `announcement_priority` NOT NULL DEFAULT `'normal'` | Приоритет: `normal` или `high` |
| `is_archived` | BOOL NOT NULL DEFAULT `false` | Флаг архивирования |
| `starts_at` | TIMESTAMPTZ | NULL = без ограничения начала показа |
| `ends_at` | TIMESTAMPTZ | NULL = без ограничения конца показа |
| `created_at` | TIMESTAMPTZ NOT NULL | |
| `updated_at` | TIMESTAMPTZ NOT NULL | |

**Ограничение:** `department_id IS NULL OR clinic_id IS NOT NULL` — нельзя указать отдел без клиники.

### Scope объявления

| `clinic_id` | `department_id` | Scope |
|---|---|---|
| NULL | NULL | Вся организация |
| заполнен | NULL | Конкретная клиника |
| заполнен | заполнен | Конкретный отдел |

### `projections.announcement_views`

| Поле | Тип | Описание |
|---|---|---|
| `announcement_id` | UUID PK FK → `domain.announcements` ON DELETE CASCADE | |
| `view_count` | BIGINT NOT NULL DEFAULT 0 | Счётчик просмотров |

Инкрементируется при каждом вызове `GetAnnouncement`. Начальная строка создаётся при создании объявления.

---

## Жизненный цикл

### Архивирование

- `is_archived = true` — объявление скрыто от сотрудников, видно только администраторам scope.
- Архивировать и де-архивировать можно **всегда**, без ограничений (операция идемпотентна).
- Де-архивированное объявление с `ends_at` в прошлом не отображается сотрудникам (фильтр по времени).

### Приоритет

- Значения: `normal` | `high`.
- Изменение приоритета **запрещено** когда `is_archived = true`.

### Физическое удаление

Объявления **никогда не удаляются напрямую**. Удаляются только каскадом при удалении организации (`ON DELETE CASCADE`).

---

## Видимость

### Для сотрудников (неадминистраторов)

Сотрудник видит объявление если:
- `is_archived = false`
- `starts_at IS NULL OR starts_at <= now()`
- `ends_at IS NULL OR ends_at >= now()`
- И одно из:
  - `clinic_id IS NULL` → объявление уровня организации сотрудника
  - `clinic_id = clinic_id сотрудника AND department_id IS NULL` → уровень его клиники
  - `department_id = department_id сотрудника` → уровень его отдела

### Для администраторов

Видят **все** объявления в рамках своего scope, включая архивированные, при передаче `include_archived = true` в List-запросах.

---

## Пагинация

Все List-методы используют cursor-based (keyset) пагинацию:
- Параметры: `limit` (int32, default=50, max=500) + `cursor` (optional string).
- Курсор — base64(JSON{`created_at`, `id`}), сортировка `created_at DESC, id DESC`.
- Ответ содержит `next_cursor` (absent, если страниц больше нет).

---

## Права доступа

### Уровень организации

| Операция | Роли |
|---|---|
| Создать, редактировать, управлять приоритетом, архивировать | `SystemAdmin`, `OrganizationAdmin` |
| Читать (все, включая архив) | `SystemAdmin`, `OrganizationAdmin` |
| Читать (активные) | Все сотрудники организации |

### Уровень клиники

| Операция | Роли |
|---|---|
| Создать, редактировать, управлять приоритетом, архивировать | `SystemAdmin`, `OrganizationAdmin`, `ClinicHead` (+ заместитель, своя клиника) |
| Читать (все, включая архив) | Те же |
| Читать (активные) | Все сотрудники клиники |

### Уровень отдела

| Операция | Роли |
|---|---|
| Создать, редактировать, управлять приоритетом, архивировать | `SystemAdmin`, `OrganizationAdmin`, `ClinicHead` (+ заместитель), `DepartmentResponsible` (+ заместитель, свой отдел) |
| Читать (все, включая архив) | Те же |
| Читать (активные) | Все сотрудники отдела |

---

## gRPC API

### Command (`AnnouncementCommandService`)

| Метод | HTTP | Описание |
|---|---|---|
| `CreateAnnouncement` | `POST /v1/announcements` | Создать объявление |
| `UpdateAnnouncement` | `PUT /v1/announcements/{id}` | Обновить title, content, starts_at, ends_at |
| `UpdateAnnouncementPriority` | `PUT /v1/announcements/{id}/priority` | Изменить приоритет (запрещено для архивных) |
| `ArchiveAnnouncement` | `POST /v1/announcements/{id}:archive` | Архивировать (идемпотентно) |
| `UnarchiveAnnouncement` | `POST /v1/announcements/{id}:unarchive` | Де-архивировать (идемпотентно) |

### Query (`AnnouncementQueryService`)

| Метод | HTTP | Описание |
|---|---|---|
| `GetAnnouncement` | `GET /v1/query/announcements/{id}` | Получить объявление, инкрементирует `view_count` |
| `ListAnnouncementsForOrganization` | `GET /v1/query/organizations/{org_id}/announcements` | Список объявлений уровня организации |
| `ListAnnouncementsForClinic` | `GET /v1/query/clinics/{clinic_id}/announcements` | Список объявлений уровня клиники (включая org-level) |
| `ListAnnouncementsForDepartment` | `GET /v1/query/departments/{dept_id}/announcements` | Список объявлений уровня отдела (включая clinic- и org-level) |

#### Фильтры List-методов

| Параметр | Тип | Описание |
|---|---|---|
| `include_archived` | bool | Включить архивные (только для администраторов scope) |
| `priority` | enum | Фильтр по приоритету (UNSPECIFIED = все) |
| `limit` | int32 | Количество на странице (default=50, max=500) |
| `cursor` | optional string | Курсор для следующей страницы |

---

## Замечания

- `author_id` бэкенд проставляет из JWT, фронт не передаёт.
- `organization_id` всегда денормализуется из `clinic_id`/`department_id` при создании.
- При создании автоматически создаётся строка в `projections.announcement_views` с `view_count = 0`.
- Если `ends_at` задан, он должен быть позже `starts_at` (валидация на уровне сервиса).
- `UpdateAnnouncement`: передать `starts_at = nil` = очистить поле (нет ограничения начала).

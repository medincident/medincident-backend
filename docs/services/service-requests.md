# Заявки (ServiceRequest)

Заявки — второй основной агрегат системы после инцидентов. Заявка создается в контексте конкретного отдела и может быть связана с инцидентом. В отличие от инцидентов, заявки всегда привязаны к отделу и имеют назначенных исполнителей.

**gRPC:** `ServiceRequestCommandService` (command) / `ServiceRequestQueryService` (query)

---

## Классификатор типов заявок

Плоский справочник типов заявок, привязанных к организации (без иерархии категорий, в отличие от классификатора инцидентов). Каждая организация ведет свой независимый набор типов.

**gRPC:** `RequestClassifierCommandService` (command) / `RequestClassifierQueryService` (query)

Подробнее: [Классификатор заявок](request/Classifier.md)

### Инварианты классификатора

- Имя типа заявки уникально среди **активных** типов в рамках одной организации.
- Деактивированный тип недоступен для выбора при создании новой заявки.
- Удаление типа невозможно, если он используется в существующих заявках (FK RESTRICT).

---

## Модель данных

### `domain.service_requests`

| Поле | Тип | Описание |
|---|---|---|
| `id` | UUID (PK) | Идентификатор заявки |
| `organization_id` | UUID (FK) | Организация |
| `clinic_id` | UUID (FK) | Клиника |
| `department_id` | UUID (FK) | Отдел (обязателен) |
| `type_id` | UUID (FK) | Тип заявки из классификатора |
| `incident_id` | UUID (FK, nullable) | Связанный инцидент (необязательный) |
| `description` | text | Описание заявки |
| `status` | string | Текущий статус |
| `author_id` | UUID (FK) | Автор заявки |
| `created_at` | timestamp | Время создания |
| `updated_at` | timestamp | Время последнего обновления |

### `domain.service_request_executors`

| Поле | Тип | Описание |
|---|---|---|
| `service_request_id` | UUID (FK) | Ссылка на заявку |
| `employee_id` | UUID (FK) | Ссылка на сотрудника-исполнителя |

---

## Статусная машина

```
created --> in_work <--> on_hold
               |
               v
        pending_review --> completed (терминальный)
               |
               v
            in_work (на доработку)

Любой нетерминальный --> cancelled (терминальный)
```

### Переходы и права

| Переход | Кто может выполнить |
|---|---|
| `created` -> `in_work` | Исполнитель |
| `in_work` -> `on_hold` | Исполнитель |
| `on_hold` -> `in_work` | Исполнитель |
| `in_work` -> `pending_review` | Исполнитель |
| `pending_review` -> `completed` | Ответственные роли |
| `pending_review` -> `in_work` | Ответственные роли |
| `*` -> `cancelled` | Ответственные роли |

**Исполнитель** — сотрудник, назначенный исполнителем заявки (`service_request_executors`).

**Ответственные роли** — `SystemAdmin`, `OrganizationAdmin`, `OrganizationHead`, `ClinicHead`, `DepartmentResponsible` (с учетом scope через цепочку `service_request -> department -> clinic -> organization`).

### Терминальные статусы

- `completed` — заявка выполнена и принята.
- `cancelled` — заявка отменена. Переход возможен из любого нетерминального статуса.

Переход из терминального статуса запрещен.

---

## Права на создание и редактирование

Создание и редактирование заявок доступно ролям:

- `SystemAdmin`
- `AdminOf.Organization(organizationID)`
- `OrgHeadOf.Organization(organizationID)`
- `ClinicHeadOf.Clinic(clinicID)`
- `DepartmentResponsibleOf.Department(departmentID)`

Роли `Employee` и `OrgDispatcher` **не могут** создавать или редактировать заявки.

---

## Инварианты

- `department_id` обязателен — каждая заявка привязана к конкретному отделу.
- Если передан `incident_id`, инцидент должен принадлежать той же организации, что и заявка.
- При создании заявки должен быть указан минимум один исполнитель.
- Каждый исполнитель должен быть сотрудником отдела (`department_id`), указанного в заявке.
- `completed` и `cancelled` — терминальные статусы, переход из них запрещен.

---

## gRPC Command API

### RequestClassifierCommandService

| Метод | HTTP | Описание |
|---|---|---|
| `CreateRequestType` | `POST /v1/organizations/{organization_id}/request-types` | Создание типа заявки |
| `UpdateRequestTypeDetails` | `PUT /v1/request-types/{type_id}/details` | Обновление имени/описания |
| `DeactivateRequestType` | `POST /v1/request-types/{type_id}/deactivations` | Деактивация типа |
| `ReactivateRequestType` | `POST /v1/request-types/{type_id}/reactivations` | Реактивация типа |
| `DeleteRequestType` | `DELETE /v1/request-types/{type_id}` | Удаление типа |

### ServiceRequestCommandService

| Метод | HTTP | Описание |
|---|---|---|
| `CreateServiceRequest` | `POST /v1/service-requests` | Создание заявки |
| `UpdateServiceRequestDescription` | `PUT /v1/service-requests/{service_request_id}/description` | Обновление описания |
| `UpdateServiceRequestStatus` | `PUT /v1/service-requests/{service_request_id}/status` | Смена статуса |
| `AssignExecutors` | `PUT /v1/service-requests/{service_request_id}/executors` | Назначение исполнителей |

---

## gRPC Query API

### RequestClassifierQueryService

| Метод | Права |
|---|---|
| `GetRequestType` | ReaderOf.Organization |
| `ListRequestTypesByOrganization` | ReaderOf.Organization |
| `ListActiveRequestTypesByOrganization` | ReaderOf.Organization |

### ServiceRequestQueryService

| Метод | Права |
|---|---|
| `GetServiceRequest` | ReaderOf.Department |
| `ListServiceRequests` | ReaderOf.Department |
| `ListServiceRequestsByIncident` | ReaderOf.Organization |
| `GetServiceRequestHistory` | ReaderOf.Department |

---

## История

Изменения заявок отслеживаются в проекционных таблицах:

- **Статусы** — каждый переход статуса фиксируется с указанием предыдущего и нового статуса, автора перехода и временной метки.
- **Исполнители** — назначение и снятие исполнителей фиксируется с указанием автора действия и временной метки.

Историю можно получить через `GetServiceRequestHistory` (query API).

← [Документация](../../README.md)

# Заявки (ServiceRequest)

Заявки — второй основной агрегат системы после инцидентов. Создаются в контексте конкретного отдела и могут быть связаны с инцидентом. Имеют назначенных исполнителей и статусную машину с разделением ролей.

**gRPC:** `ServiceRequestCommandService` (command) / `ServiceRequestQueryService` (query)

---

## Статусная машина

```
created → in_work ←→ on_hold
              │
              ▼
       pending_review → completed (терминальный)
              │
              ▼
          in_work (на доработку)

Любой нетерминальный → cancelled (терминальный)
```

### Переходы и права

| Переход | Кто может выполнить |
|---|---|
| `created` → `in_work` | Исполнитель |
| `in_work` → `on_hold` | Исполнитель |
| `on_hold` → `in_work` | Исполнитель |
| `in_work` → `pending_review` | Исполнитель |
| `pending_review` → `completed` | Ответственные роли |
| `pending_review` → `in_work` | Ответственные роли |
| `*` → `cancelled` | Ответственные роли |

**Исполнитель** — сотрудник, назначенный исполнителем заявки (`service_request_executors`).

**Ответственные роли** — `SystemAdmin`, `OrganizationAdmin`, `OrganizationHead`, `ClinicHead`, `DepartmentResponsible` (с учётом scope).

---

## Command-методы

### CreateServiceRequest

**HTTP:** `POST /v1/service-requests`
**gRPC:** `ServiceRequestCommandService.CreateServiceRequest`

#### Права доступа

`privilegedActorPolicy(organizationID, clinicID, departmentID)` — SystemAdmin, OrgAdmin, OrgHead, ClinicHead, DeptResponsible.

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `department_id` | string (UUID) | required, uuid |
| `type_id` | string (UUID) | required, uuid |
| `incident_id` | string (UUID) | omitempty, uuid |
| `description` | string | required, min=1, max=10000 |
| `executor_employee_ids` | []string (UUID) | required, min=1, каждый — uuid |

#### Инварианты

- `department_id` обязателен.
- `type_id` должен принадлежать той же организации и быть активным.
- Если передан `incident_id`, инцидент должен принадлежать той же организации.
- Минимум один исполнитель; каждый должен быть сотрудником указанного отдела.
- `author_id` проставляется автоматически из JWT.

#### Ошибки

| Код | Описание |
|---|---|
| `service_request_department_not_found` | Отдел не найден |
| `service_request_type_not_found` | Тип заявки не найден |
| `service_request_type_inactive` | Тип заявки деактивирован |
| `service_request_type_org_mismatch` | Тип не принадлежит организации |
| `service_request_incident_not_found` | Инцидент не найден |
| `service_request_incident_org_mismatch` | Инцидент из другой организации |
| `service_request_employee_not_found` | Сотрудник-исполнитель не найден |
| `service_request_employee_dept_mismatch` | Исполнитель не из указанного отдела |

---

### UpdateServiceRequestDescription

**HTTP:** `PUT /v1/service-requests/{service_request_id}/description`
**gRPC:** `ServiceRequestCommandService.UpdateServiceRequestDescription`

#### Права доступа

`privilegedActorPolicy` (те же роли, что и при создании).

#### Инварианты

- Нельзя изменять терминальные заявки (`completed`, `cancelled`).

---

### UpdateServiceRequestStatus

**HTTP:** `PUT /v1/service-requests/{service_request_id}/status`
**gRPC:** `ServiceRequestCommandService.UpdateServiceRequestStatus`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `service_request_id` | string (UUID) | required, uuid |
| `new_status` | string | required, oneof: in_work, on_hold, pending_review, completed, cancelled |

#### Инварианты

- Переход проверяется по статусной машине (см. выше).
- Переходы исполнителя не требуют дополнительной авторизации (caller проверяется как executor).
- Переходы ответственных ролей требуют `privilegedActorPolicy`.
- Нельзя переходить из терминальных статусов.

---

### AssignExecutors

**HTTP:** `PUT /v1/service-requests/{service_request_id}/executors`
**gRPC:** `ServiceRequestCommandService.AssignExecutors`

#### Права доступа

`privilegedActorPolicy` (те же роли, что и при создании).

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `service_request_id` | string (UUID) | required, uuid |
| `executor_employee_ids` | []string (UUID) | required, min=1, каждый — uuid |

#### Инварианты

- Замена выполняется diff-based: удаляются отсутствующие, добавляются новые.
- Каждый исполнитель должен быть сотрудником отдела заявки.
- Нельзя изменять терминальные заявки.
- Изменения фиксируются в `projections.service_request_executor_history`.

---

## Query-методы

| Метод | HTTP | Права |
|---|---|---|
| `GetServiceRequest` | `GET /v1/service-requests/{id}` | ReaderOf.Organization |
| `ListServiceRequests` | `GET /v1/organizations/{organization_id}/service-requests` | ReaderOf.Organization |
| `ListServiceRequestsByIncident` | `GET /v1/incidents/{incident_id}/service-requests` | ReaderOf.Organization |
| `GetServiceRequestHistory` | `GET /v1/service-requests/{service_request_id}/history` | ReaderOf.Organization |

### GetServiceRequestHistory

Возвращает две временные линии:
- **Статусы** — каждый переход с `old_status`, `new_status`, `actor_id`, `actor_name`, `changed_at`.
- **Исполнители** — каждое назначение/снятие с `action` (`assigned`/`removed`), `employee_id`, `employee_name`, `actor_id`, `actor_name`, `changed_at`.

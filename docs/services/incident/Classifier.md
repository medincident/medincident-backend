# Классификатор инцидентов

Управление иерархией категорий и типов инцидентов. Каждая организация ведёт свой независимый классификатор.

**gRPC:** `IncidentClassifierCommandService` (command) / `IncidentClassifierQueryService` (query)

---

## Категории

### CreateIncidentCategory

**HTTP:** `POST /v1/organizations/{organization_id}/incident-categories`
**gRPC:** `IncidentClassifierCommandService.CreateIncidentCategory`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `name` | string | required, min=2, max=256 |
| `parent_category_id` | string (UUID) | omitempty, uuid |
| `description` | string | omitempty, max=2048 |

#### Инварианты

- Категории образуют дерево. Глубина ограничена константой `incidentClassifierMaxDepth`.
- Если `parent_category_id` задан — родитель должен принадлежать той же организации и быть активным.
- При создании категория имеет статус `active`.

#### Ошибки

| Код | Описание |
|---|---|
| `incident_category_not_found` | Родительская категория не найдена |
| `incident_category_max_depth_exceeded` | Превышена максимальная глубина дерева |

---

### UpdateIncidentCategoryDetails

**HTTP:** `PATCH /v1/incident-categories/{category_id}`
**gRPC:** `IncidentClassifierCommandService.UpdateIncidentCategoryDetails`

#### Права доступа

`AdminOf.Category(categoryID)`

---

### MoveIncidentCategory

**HTTP:** `PUT /v1/incident-categories/{category_id}/parent`
**gRPC:** `IncidentClassifierCommandService.MoveIncidentCategory`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Инварианты

- Нельзя переместить категорию в саму себя или в своего потомка (цикл).
- Ограничение глубины применяется к новому положению.

---

### DeactivateIncidentCategory / ReactivateIncidentCategory

**HTTP:** `POST /v1/incident-categories/{category_id}/deactivate` / `reactivate`
**gRPC:** `IncidentClassifierCommandService.DeactivateIncidentCategory` / `ReactivateIncidentCategory`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Инварианты

- Деактивация категории рекурсивно деактивирует все дочерние категории и типы.
- Деактивированный тип недоступен для выбора при создании нового инцидента.

---

### DeleteIncidentCategory

**HTTP:** `DELETE /v1/incident-categories/{category_id}`
**gRPC:** `IncidentClassifierCommandService.DeleteIncidentCategory`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Инварианты

- Удаление невозможно, если к категории привязаны инциденты или дочерние типы с инцидентами (FK RESTRICT).

---

## Типы инцидентов

### CreateIncidentType

**HTTP:** `POST /v1/incident-categories/{category_id}/incident-types`
**gRPC:** `IncidentClassifierCommandService.CreateIncidentType`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `category_id` | string (UUID) | required, uuid |
| `name` | string | required, min=2, max=256 |
| `description` | string | omitempty, max=2048 |

---

### UpdateIncidentTypeDetails

**HTTP:** `PATCH /v1/incident-types/{type_id}`
**gRPC:** `IncidentClassifierCommandService.UpdateIncidentTypeDetails`

#### Права доступа

`AdminOf.IncidentType(typeID)`

---

### MoveIncidentType

**HTTP:** `PUT /v1/incident-types/{type_id}/category`
**gRPC:** `IncidentClassifierCommandService.MoveIncidentType`

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Инварианты

- Целевая категория должна принадлежать той же организации.

---

### DeactivateIncidentType / ReactivateIncidentType

#### Права доступа

`AdminOf.IncidentType(typeID)`

---

### DeleteIncidentType

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Инварианты

- Удаление невозможно, если тип используется в существующих инцидентах (FK RESTRICT).

---

### AllowIncidentTypeForPatients / DisallowIncidentTypeForPatients

Управление флагом `is_allowed_for_patients` — доступность типа для выбора пациентом при подаче заявки через буфер.

**gRPC:** `IncidentClassifierCommandService.AllowIncidentTypeForPatients` / `DisallowIncidentTypeForPatients`

#### Права доступа

`AdminOf.IncidentType(typeID)`

---

## Query-методы

| Метод | Права |
|---|---|
| `ListIncidentCategories` | ReaderOf.Organization |
| `GetIncidentCategory` | ReaderOf.Category |
| `ListIncidentTypes` | ReaderOf.Organization |
| `GetIncidentType` | ReaderOf.IncidentType |

Для пациентов (patient-facing): `ListIncidentCategories` и `ListIncidentTypes` — `Authenticated` (только `is_allowed_for_patients=true` типы).

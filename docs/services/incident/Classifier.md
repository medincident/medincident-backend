← [Документация](../../README.md)

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

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `incident_category_parent_inactive` | 422 | Родительская категория деактивирована |
| `incident_category_parent_organization_mismatch` | 422 | Родительская категория принадлежит другой организации |
| `incident_category_max_depth_exceeded` | 422 | Превышена максимальная глубина дерева категорий |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_category_parent_not_found` | 404 | Родительская категория не найдена |
| `incident_category_name_conflict` | 409 | Категория с таким именем уже существует |

---

### UpdateIncidentCategoryDetails

**HTTP:** `PUT /v1/incident-categories/{category_id}/details`
**gRPC:** `IncidentClassifierCommandService.UpdateIncidentCategoryDetails`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_category_not_found` | 404 | Категория не найдена |
| `incident_category_name_conflict` | 409 | Категория с таким именем уже существует |

---

### MoveIncidentCategory

**HTTP:** `POST /v1/incident-categories/{category_id}:move`
**gRPC:** `IncidentClassifierCommandService.MoveIncidentCategory`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Инварианты

- Нельзя переместить категорию в саму себя или в своего потомка (цикл).
- Ограничение глубины применяется к новому положению.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_category_not_found` | 404 | Категория не найдена |
| `incident_category_parent_not_found` | 404 | Родительская категория не найдена |

---

### DeactivateIncidentCategory

**HTTP:** `POST /v1/incident-categories/{category_id}/deactivations`
**gRPC:** `IncidentClassifierCommandService.DeactivateIncidentCategory`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Инварианты

- Деактивация категории рекурсивно деактивирует все дочерние категории и типы.
- Деактивированный тип недоступен для выбора при создании нового инцидента.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_category_not_found` | 404 | Категория не найдена |

---

### ReactivateIncidentCategory

**HTTP:** `POST /v1/incident-categories/{category_id}/reactivations`
**gRPC:** `IncidentClassifierCommandService.ReactivateIncidentCategory`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_category_not_found` | 404 | Категория не найдена |

---

### DeleteIncidentCategory

**HTTP:** `DELETE /v1/incident-categories/{category_id}`
**gRPC:** `IncidentClassifierCommandService.DeleteIncidentCategory`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Инварианты

- Удаление невозможно, если к категории привязаны инциденты или дочерние типы с инцидентами (FK RESTRICT).

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_category_not_found` | 404 | Категория не найдена |

---

## Типы инцидентов

### CreateIncidentType

**HTTP:** `POST /v1/incident-categories/{category_id}/types`
**gRPC:** `IncidentClassifierCommandService.CreateIncidentType`

#### Права доступа

`AdminOf.Category(categoryID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `category_id` | string (UUID) | required, uuid |
| `name` | string | required, min=2, max=256 |
| `description` | string | omitempty, max=2048 |

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_category_not_found` | 404 | Категория не найдена |
| `incident_type_name_conflict` | 409 | Тип с таким именем уже существует в категории |

---

### UpdateIncidentTypeDetails

**HTTP:** `PUT /v1/incident-types/{type_id}/details`
**gRPC:** `IncidentClassifierCommandService.UpdateIncidentTypeDetails`

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_not_found` | 404 | Тип не найден |
| `incident_type_name_conflict` | 409 | Тип с таким именем уже существует в категории |

---

### MoveIncidentType

**HTTP:** `POST /v1/incident-types/{type_id}:move`
**gRPC:** `IncidentClassifierCommandService.MoveIncidentType`

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Инварианты

- Целевая категория должна принадлежать той же организации.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_not_found` | 404 | Тип не найден |
| `incident_type_category_not_found` | 404 | Целевая категория не найдена |

---

### DeactivateIncidentType

**HTTP:** `POST /v1/incident-types/{type_id}/deactivations`
**gRPC:** `IncidentClassifierCommandService.DeactivateIncidentType`

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_not_found` | 404 | Тип не найден |

---

### ReactivateIncidentType

**HTTP:** `POST /v1/incident-types/{type_id}/reactivations`
**gRPC:** `IncidentClassifierCommandService.ReactivateIncidentType`

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_not_found` | 404 | Тип не найден |

---

### DeleteIncidentType

**HTTP:** `DELETE /v1/incident-types/{type_id}`
**gRPC:** `IncidentClassifierCommandService.DeleteIncidentType`

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Инварианты

- Удаление невозможно, если тип используется в существующих инцидентах (FK RESTRICT).

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_not_found` | 404 | Тип не найден |

---

### AllowIncidentTypeForPatients

**HTTP:** `POST /v1/incident-types/{type_id}/patient-allowances`
**gRPC:** `IncidentClassifierCommandService.AllowIncidentTypeForPatients`

Управление флагом `is_allowed_for_patients` — доступность типа для выбора пациентом при подаче заявки через буфер.

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_not_found` | 404 | Тип не найден |

---

### DisallowIncidentTypeForPatients

**HTTP:** `DELETE /v1/incident-types/{type_id}/patient-allowances`
**gRPC:** `IncidentClassifierCommandService.DisallowIncidentTypeForPatients`

#### Права доступа

`AdminOf.IncidentType(typeID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_not_found` | 404 | Тип не найден |

---

## Query-методы

Все list-методы принимают любого аутентифицированного вызывающего (`Authenticated`). Сотрудники организации (`ReaderOf.Organization` / `ReaderOf.Category`) получают полный набор данных. Пациенты (аутентифицированные, но не являющиеся сотрудниками) автоматически получают отфильтрованный результат: только активные категории, в поддереве которых есть хотя бы один активный тип с `is_allowed_for_patients=true`, и только такие типы. Переключение между режимами происходит внутри сервиса через `authz.Satisfies` — `permission_denied` для пациентов не возникает.

| Метод | Права | Примечание |
|---|---|---|
| `GetCategory` | ReaderOf.Category | Только сотрудники |
| `ListCategoriesByOrganization` | Authenticated | Пациенты видят только patient-visible категории |
| `ListRootCategories` | Authenticated | Пациенты видят только активные корни с patient-allowed типами в поддереве |
| `ListCategorySubtree` | Authenticated | Пациенты видят только активные узлы поддерева с patient-allowed типами |
| `GetType` | ReaderOf.IncidentType | Только сотрудники |
| `ListTypesByCategory` | Authenticated | Пациенты видят только активные типы с `is_allowed_for_patients=true` |
| `ListTypesByOrganization` | Authenticated | Пациенты видят только типы с `is_allowed_for_patients=true` |

### Параметр include_deactivated

Методы `ListCategoriesByOrganization`, `ListRootCategories` и `ListTypesByOrganization` поддерживают параметр `include_deactivated bool` (по умолчанию `false`).

- `include_deactivated=false` (по умолчанию): возвращаются только активные записи. Пациентам всегда возвращается только активное подмножество независимо от значения флага.
- `include_deactivated=true`: возвращаются все записи, включая деактивированные. Требует `SystemAdmin` или `OrgAdminOf` соответствующей организации. Без нужных прав возвращается `permission_denied`.

**Поведение Get-методов для деактивированных записей:**

`GetCategory` и `GetType` не принимают параметр `include_deactivated`. Видимость деактивированных записей определяется автоматически по роли вызывающего:
- Не-администратор: деактивированная запись возвращается как `*_not_found` (неотличимо от отсутствующей).
- Администратор: возвращается фактическая запись.

### GetCategory

**HTTP:** `GET /v1/incident-categories/{id}`
**gRPC:** `IncidentClassifierQueryService.GetCategory`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |
| `incident_category_not_found` | 404 | Категория не найдена или деактивирована (для не-администратора) |

### ListCategoriesByOrganization

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `incident_classifier_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |
| `permission_denied` | 403 | Нет прав доступа (include_deactivated=true без прав администратора) |

### ListRootCategories

Заменяет удалённый метод `ListActiveRootCategories`. Возвращает корневые категории (без родителя) организации. По умолчанию (`include_deactivated=false`) — только активные.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `incident_classifier_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |
| `permission_denied` | 403 | Нет прав доступа (include_deactivated=true без прав администратора) |

### ListCategorySubtree

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа (не аутентифицирован) |

### GetType

**HTTP:** `GET /v1/incident-types/{id}`
**gRPC:** `IncidentClassifierQueryService.GetType`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |
| `incident_type_not_found` | 404 | Тип не найден или деактивирован (для не-администратора) |

### ListTypesByCategory

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `incident_classifier_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |
| `permission_denied` | 403 | Нет прав доступа (не аутентифицирован) |

### ListTypesByOrganization

Заменяет удалённый метод `ListActiveTypesByOrganization`. Возвращает все типы инцидентов организации. По умолчанию (`include_deactivated=false`) — только активные.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `incident_classifier_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |
| `permission_denied` | 403 | Нет прав доступа (include_deactivated=true без прав администратора) |

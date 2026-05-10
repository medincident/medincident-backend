← [Документация](../../README.md)

# Классификатор заявок

Плоский справочник типов заявок, привязанных к организации. В ��тличие от классификатора инцидентов, здесь нет иерархии категорий — только типы.

**gRPC:** `RequestClassifierCommandService` (command) / `RequestClassifierQueryService` (query)

---

## Типы заявок

### CreateRequestType

**HTTP:** `POST /v1/organizations/{organization_id}/request-types`
**gRPC:** `RequestClassifierCommandService.CreateRequestType`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `name` | string | required, min=2, max=256 |
| `description` | string | omitempty, min=8, max=2048 |

#### Инварианты

- Имя типа уникально среди активных типов в рамках одной организации.
- При создании тип имеет статус `active`.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Недостаточно прав |
| `request_type_name_conflict` | 409 | Активный тип с таким именем уже существует |

---

### UpdateRequestTypeDetails

**HTTP:** `PUT /v1/request-types/{type_id}/details`
**gRPC:** `RequestClassifierCommandService.UpdateRequestTypeDetails`

#### Права доступа

`AdminOf.RequestType(typeID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `type_id` | string (UUID) | required, uuid |
| `name` | string | required, min=2, max=256 |
| `description` | string | omitempty, min=8, max=2048 |

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Недостаточно прав |
| `request_type_not_found` | 404 | Тип не найден |
| `request_type_name_conflict` | 409 | Активный тип с таким именем уже существует |

---

### DeactivateRequestType

**HTTP:** `POST /v1/request-types/{type_id}/deactivations`
**gRPC:** `RequestClassifierCommandService.DeactivateRequestType`

#### Права доступа

`AdminOf.RequestType(typeID)`

#### Инварианты

- Идемпотентная операция: если тип уже неактивен, операция завершается без ошибки.
- Деактивированный тип недоступен для выбора при создании новой заявки.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Недостаточно прав |
| `request_type_not_found` | 404 | Тип не найден |

---

### ReactivateRequestType

**HTTP:** `POST /v1/request-types/{type_id}/reactivations`
**gRPC:** `RequestClassifierCommandService.ReactivateRequestType`

#### Права доступа

`AdminOf.RequestType(typeID)`

#### Инварианты

- Идемпотентная операция: если тип уже активен, операция завершается без ошибки.
- Может вызвать `request_type_name_conflict`, если активный тип с таким именем уже существует.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Недостаточно прав |
| `request_type_not_found` | 404 | Тип не найден |
| `request_type_name_conflict` | 409 | Активный тип с таким именем уже существует |

---

### DeleteRequestType

**HTTP:** `DELETE /v1/request-types/{type_id}`
**gRPC:** `RequestClassifierCommandService.DeleteRequestType`

#### Права доступа

`AdminOf.RequestType(typeID)`

#### Инварианты

- Удаление невозможно, если тип используется в существующих заявках (FK RESTRICT).

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Недостаточно прав |
| `request_type_not_found` | 404 | Тип не найден |

---

## Query-методы

| Метод | HTTP | Права |
|---|---|---|
| `GetRequestType` | `GET /v1/request-types/{id}` | ReaderOf.RequestType |
| `ListRequestTypesByOrganization` | `GET /v1/organizations/{organization_id}/request-types` | ReaderOf.Organization |
| `ListActiveRequestTypesByOrganization` | `GET /v1/organizations/{organization_id}/request-types:active` | ReaderOf.Organization |

### GetRequestType

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибка валидации |
| `permission_denied` | 403 | Недостаточно прав |
| `request_type_not_found` | 404 | Тип не найден |

← [Документация](../README.md)

# Орг. структура

Управление организациями, клиниками и отделами.

**gRPC:** `OrgStructureCommandService` (command) / `OrgStructureQueryService` (query)

---

## CreateOrganization

**HTTP:** `POST /v1/organizations`
**gRPC:** `OrgStructureCommandService.CreateOrganization`

### Права доступа

`SystemAdmin`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `name` | string | required, min=4, max=256 |
| `legal_address` | Address | required |

### Инварианты

- Имя обрезается (trim) перед валидацией.
- Адрес валидируется через `address.ValidateAddress`.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |

---

## UpdateOrganizationDetails

**HTTP:** `PUT /v1/organizations/{organization_id}/details`
**gRPC:** `OrgStructureCommandService.UpdateOrganizationDetails`

### Права доступа

`AdminOf.Organization(organizationID)` — SystemAdmin или OrgAdmin организации.

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `name` | string | required, min=4, max=256 |

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_not_found` | 404 | Организация не найдена |

---

## UpdateOrganizationLegalAddress

**HTTP:** `PUT /v1/organizations/{organization_id}/legal-address`
**gRPC:** `OrgStructureCommandService.UpdateOrganizationLegalAddress`

### Права доступа

`AdminOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `legal_address` | Address | required |

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_not_found` | 404 | Организация не найдена |

---

## CreateClinic

**HTTP:** `POST /v1/organizations/{organization_id}/clinics`
**gRPC:** `OrgStructureCommandService.CreateClinic`

### Права доступа

`AdminOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `name` | string | required, min=4, max=256 |
| `physical_address` | Address | required |

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_organization_not_found` | 404 | Организация не найдена |

---

## UpdateClinicDetails

**HTTP:** `PUT /v1/clinics/{clinic_id}/details`
**gRPC:** `OrgStructureCommandService.UpdateClinicDetails`

### Права доступа

`AdminOf.Clinic(clinicID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_not_found` | 404 | Клиника не найдена |

---

## UpdateClinicPhysicalAddress

**HTTP:** `PUT /v1/clinics/{clinic_id}/physical-address`
**gRPC:** `OrgStructureCommandService.UpdateClinicPhysicalAddress`

### Права доступа

`AdminOf.Clinic(clinicID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_not_found` | 404 | Клиника не найдена |

---

## CreateDepartment

**HTTP:** `POST /v1/clinics/{clinic_id}/departments`
**gRPC:** `OrgStructureCommandService.CreateDepartment`

### Права доступа

`AdminOf.Clinic(clinicID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `clinic_id` | string (UUID) | required, uuid |
| `name` | string | required, min=4, max=256 |

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_clinic_not_found` | 404 | Клиника не найдена |

---

## UpdateDepartmentDetails

**HTTP:** `PUT /v1/departments/{department_id}/details`
**gRPC:** `OrgStructureCommandService.UpdateDepartmentDetails`

### Права доступа

`AdminOf.Department(departmentID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_not_found` | 404 | Отдел не найден |

---

## DeactivateOrganization

**HTTP:** `POST /v1/organizations/{organization_id}:deactivate`
**gRPC:** `OrgStructureCommandService.DeactivateOrganization`

Каскадно деактивирует организацию, все её клиники, все отделы клиник и всех сотрудников отделов. Каждый сотрудник теряет все свои роли. На каждую фактически изменённую строку публикуется соответствующее доменное событие. Операция идемпотентна: повторный вызов для уже неактивной организации завершается успехом без дополнительных событий (для самой организации), но каскад по активным дочерним сущностям выполняется.

### Права доступа

`AdminOf.Organization(organizationID)` — SystemAdmin или OrgAdmin организации.

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_not_found` | 404 | Организация не найдена |

---

## ActivateOrganization

**HTTP:** `POST /v1/organizations/{organization_id}:activate`
**gRPC:** `OrgStructureCommandService.ActivateOrganization`

Активирует организацию. Дочерние сущности (клиники, отделы, сотрудники) не затрагиваются — каждая остаётся в своём текущем состоянии. Идемпотентен: если организация уже активна, событие не публикуется.

### Права доступа

`AdminOf.Organization(organizationID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_not_found` | 404 | Организация не найдена |

---

## DeleteOrganization

**HTTP:** `DELETE /v1/organizations/{organization_id}`
**gRPC:** `OrgStructureCommandService.DeleteOrganization`

Жёстко удаляет организацию. Перед удалением организации каскадно удаляются все роли сотрудников, сами сотрудники, отделы и клиники в порядке, диктуемом FK-ограничениями. Если на организацию ссылаются внешние сущности (инциденты и т.п.), возвращается `organization_delete_has_dependents`.

### Права доступа

`AdminOf.Organization(organizationID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_not_found` | 404 | Организация не найдена |
| `organization_delete_has_dependents` | 409 | Организация имеет зависимые объекты и не может быть удалена |

---

## DeactivateClinic

**HTTP:** `POST /v1/clinics/{clinic_id}:deactivate`
**gRPC:** `OrgStructureCommandService.DeactivateClinic`

Каскадно деактивирует клинику, все её отделы и всех сотрудников отделов. Каждый сотрудник теряет все свои роли.

### Права доступа

`AdminOf.Clinic(clinicID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_not_found` | 404 | Клиника не найдена |

---

## ActivateClinic

**HTTP:** `POST /v1/clinics/{clinic_id}:activate`
**gRPC:** `OrgStructureCommandService.ActivateClinic`

Активирует клинику. Требует, чтобы родительская организация была активна. Дочерние отделы и сотрудники не затрагиваются.

### Права доступа

`AdminOf.Clinic(clinicID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_not_found` | 404 | Клиника не найдена |
| `clinic_activate_parent_inactive` | 422 | Родительская организация неактивна |

---

## DeleteClinic

**HTTP:** `DELETE /v1/clinics/{clinic_id}`
**gRPC:** `OrgStructureCommandService.DeleteClinic`

Жёстко удаляет клинику. Перед удалением каскадно удаляются все роли сотрудников, сами сотрудники и отделы клиники.

### Права доступа

`AdminOf.Clinic(clinicID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_not_found` | 404 | Клиника не найдена |
| `clinic_delete_has_dependents` | 409 | Клиника имеет зависимые объекты и не может быть удалена |

---

## DeactivateDepartment

**HTTP:** `POST /v1/departments/{department_id}:deactivate`
**gRPC:** `OrgStructureCommandService.DeactivateDepartment`

Каскадно деактивирует отдел и всех его сотрудников. Каждый сотрудник теряет все роли.

### Права доступа

`AdminOf.Department(departmentID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_not_found` | 404 | Отдел не найден |

---

## ActivateDepartment

**HTTP:** `POST /v1/departments/{department_id}:activate`
**gRPC:** `OrgStructureCommandService.ActivateDepartment`

Активирует отдел. Требует, чтобы родительская клиника была активна. Сотрудники не затрагиваются.

### Права доступа

`AdminOf.Department(departmentID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_not_found` | 404 | Отдел не найден |
| `department_activate_parent_inactive` | 422 | Родительская клиника неактивна |

---

## DeleteDepartment

**HTTP:** `DELETE /v1/departments/{department_id}`
**gRPC:** `OrgStructureCommandService.DeleteDepartment`

Жёстко удаляет отдел. Перед удалением каскадно удаляются все роли сотрудников и сами сотрудники отдела.

### Права доступа

`AdminOf.Department(departmentID)`

### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_not_found` | 404 | Отдел не найден |
| `department_delete_has_dependents` | 409 | Отдел имеет зависимые объекты и не может быть удалён |

---

## Query-методы

Доступны через `OrgStructureQueryService`. Организации — публичный каталог (только аутентификация). Клиники и отделы требуют `ReaderOf.*`.

| Метод | Права |
|---|---|
| `ListOrganizations` | Authenticated |
| `GetOrganization` | Authenticated |
| `ListClinicsByOrganization` | ReaderOf.Organization |
| `CountClinicsByOrganization` | ReaderOf.Organization |
| `GetClinic` | ReaderOf.Clinic |
| `ListDepartmentsByClinic` | ReaderOf.Clinic |
| `CountDepartmentsByClinic` | ReaderOf.Clinic |
| `GetDepartment` | ReaderOf.Department |

### Ошибки

| Метод | Код | HTTP | Описание |
|---|---|---|---|
| Все | `permission_denied` | 403 | Нет прав доступа |
| `GetOrganization` | `organization_not_found` | 404 | Организация не найдена |
| `GetClinic` | `clinic_not_found` | 404 | Клиника не найдена |
| `GetDepartment` | `department_not_found` | 404 | Отдел не найден |

### ListOrganizations

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `orgstructure_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |
| `permission_denied` | 403 | Нет прав доступа |

### ListClinicsByOrganization

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `orgstructure_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |
| `permission_denied` | 403 | Нет прав доступа |

### ListDepartmentsByClinic

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `orgstructure_bad_cursor` | 400 | Недопустимый или некорректный курсор пагинации |
| `permission_denied` | 403 | Нет прав доступа |

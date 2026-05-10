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

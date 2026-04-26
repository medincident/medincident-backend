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

| Код | Описание |
|---|---|
| `validation_failed` | Нарушение правил валидации |
| `permission_denied` | Нет прав SystemAdmin |

---

## UpdateOrganizationDetails

**HTTP:** `PATCH /v1/organizations/{organization_id}`
**gRPC:** `OrgStructureCommandService.UpdateOrganizationDetails`

### Права доступа

`AdminOf.Organization(organizationID)` — SystemAdmin или OrgAdmin организации.

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `name` | string | required, min=4, max=256 |

### Ошибки

| Код | Описание |
|---|---|
| `validation_failed` | Нарушение правил валидации |
| `permission_denied` | Нет прав |
| `organization_not_found` | Организация не найдена |

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

---

## UpdateClinicDetails

**HTTP:** `PATCH /v1/clinics/{clinic_id}`
**gRPC:** `OrgStructureCommandService.UpdateClinicDetails`

### Права доступа

`AdminOf.Clinic(clinicID)`

---

## UpdateClinicPhysicalAddress

**HTTP:** `PUT /v1/clinics/{clinic_id}/physical-address`
**gRPC:** `OrgStructureCommandService.UpdateClinicPhysicalAddress`

### Права доступа

`AdminOf.Clinic(clinicID)`

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

---

## UpdateDepartmentDetails

**HTTP:** `PATCH /v1/departments/{department_id}`
**gRPC:** `OrgStructureCommandService.UpdateDepartmentDetails`

### Права доступа

`AdminOf.Department(departmentID)`

---

## Query-методы

Доступны через `OrgStructureQueryService`. Организации — публичный каталог (только аутентификация). Клиники и отделы требуют `ReaderOf.*`.

| Метод | Права |
|---|---|
| `ListOrganizations` | Authenticated |
| `GetOrganization` | Authenticated |
| `ListClinics` | ReaderOf.Organization |
| `GetClinic` | ReaderOf.Clinic |
| `ListDepartments` | ReaderOf.Clinic |
| `GetDepartment` | ReaderOf.Department |

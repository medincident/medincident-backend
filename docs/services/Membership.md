← [Документация](../README.md)

# Сотрудники и роли

Управление сотрудниками, назначение организационных ролей, учёт отпусков.

**gRPC:** `MembershipCommandService` (command) / `MembershipQueryService` (query)

---

## Сотрудники

### HireEmployee

**HTTP:** `POST /v1/organizations/{organization_id}/employees`
**gRPC:** `MembershipCommandService.HireEmployee`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |
| `zitadel_user_id` | string | required, min=1 |
| `department_id` | string (UUID) | omitempty, uuid |

#### Инварианты

- Сотрудник привязан к организации; отдел опциональный.
- Один Zitadel-пользователь не может быть нанят дважды в одну организацию.

#### Ошибки

| Код | Описание |
|---|---|
| `employee_already_hired` | Пользователь уже является сотрудником |
| `department_not_found` | Отдел не найден или не принадлежит организации |

---

### TerminateEmployee

**HTTP:** `DELETE /v1/employees/{employee_id}`
**gRPC:** `MembershipCommandService.TerminateEmployee`

#### Права доступа

`AdminOf.Employee(employeeID)`

#### Инварианты

- При увольнении автоматически очищаются все роли сотрудника и слоты заместителей в которых он числится.
- Активные отпуска закрываются.

---

### UpdateEmployeeDepartment

**HTTP:** `PATCH /v1/employees/{employee_id}/department`
**gRPC:** `MembershipCommandService.UpdateEmployeeDepartment`

#### Права доступа

`AdminOf.Employee(employeeID)`

---

## Отпуска

### ScheduleVacation

**HTTP:** `POST /v1/employees/{employee_id}/vacations`
**gRPC:** `MembershipCommandService.ScheduleVacation`

#### Права доступа

`AdminOf.Employee(employeeID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `employee_id` | string (UUID) | required, uuid |
| `starts_at` | timestamp | required |
| `ends_at` | timestamp | omitempty |

#### Инварианты

- `ends_at` — опционально (открытый отпуск).
- Если задано, `ends_at > starts_at`.
- Перекрывающиеся отпуска не допускаются.

---

### CancelVacation

**HTTP:** `DELETE /v1/vacations/{vacation_id}`
**gRPC:** `MembershipCommandService.CancelVacation`

#### Права доступа

`AdminOf.Vacation(vacationID)`

---

## Роль OrgAdmin

### AssignOrgAdmin / RevokeOrgAdmin

Назначение/снятие OrgAdmin на уровне организации.

#### Права доступа

`SystemAdmin`

### AssignOrgAdminDeputy / RemoveOrgAdminDeputy

Назначение/снятие заместителя OrgAdmin.

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Инварианты

- Заместитель должен быть сотрудником той же организации.
- Заместитель не может совпадать с основным держателем роли.

---

## Роль OrgHead

### AssignOrgHead / RevokeOrgHead

Назначение/снятие OrgHead.

#### Права доступа

`AdminOf.Organization(organizationID)`

### AssignOrgHeadDeputy / RemoveOrgHeadDeputy

#### Права доступа

`AdminOf.Organization(organizationID)`

---

## Роль OrgDispatcher

### AssignOrgDispatcher / RevokeOrgDispatcher

Назначение/снятие OrgDispatcher.

#### Права доступа

`AdminOf.Organization(organizationID)`

### AssignOrgDispatcherDeputy / RemoveOrgDispatcherDeputy

#### Права доступа

`AdminOf.Organization(organizationID)`

---

## Роль ClinicHead

### AssignClinicHead / RevokeClinicHead

#### Права доступа

`AdminOf.Clinic(clinicID)`

### AssignClinicHeadDeputy / RemoveClinicHeadDeputy

#### Права доступа

`AdminOf.Clinic(clinicID)`

---

## Роль DeptResponsible

### AssignDepartmentResponsible / RevokeDepartmentResponsible

#### Права доступа

`AdminOf.Department(departmentID)`

### AssignDepartmentResponsibleDeputy / RemoveDepartmentResponsibleDeputy

#### Права доступа

`AdminOf.Department(departmentID)`

---

## Query-методы

| Метод | Права |
|---|---|
| `GetEmployee` | ReaderOf.Employee |
| `ListEmployees` | ReaderOf.Organization |
| `GetVacation` | SystemAdmin + OrgAdminOf.Employee + SelfEmployee |
| `ListVacations` | SystemAdmin + OrgAdminOf.Employee + SelfEmployee |
| `ListRoles` | ReaderOf.Organization |

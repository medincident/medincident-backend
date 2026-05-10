← [Документация](../README.md)

# Сотрудники и роли

Управление сотрудниками, назначение организационных ролей, учёт отпусков.

**gRPC:** `MembershipCommandService` (command) / `MembershipQueryService` (query)

---

## Сотрудники

### HireEmployee

**HTTP:** `POST /v1/employees`
**gRPC:** `MembershipCommandService.HireEmployee`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `zitadel_user_id` | string | required, min=1 |
| `department_id` | string (UUID) | required, uuid |
| `position` | string | omitempty |

#### Инварианты

- Один Zitadel-пользователь не может быть нанят дважды в одну организацию.
- Отдел должен принадлежать организации.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_not_found` | 404 | Отдел не найден |
| `zitadel_user_not_found` | 404 | Пользователь Zitadel не найден |
| `employee_already_hired` | 409 | Пользователь уже является сотрудником |
| `zitadel_verify_failed` | 503 | Внешний сервис идентификации недоступен |

---

### UpdateEmployeePosition

**HTTP:** `PUT /v1/employees/{employee_id}/position`
**gRPC:** `MembershipCommandService.UpdateEmployeePosition`

#### Права доступа

`AdminOf.Employee(employeeID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `employee_id` | string (UUID) | required, uuid |
| `position` | string | omitempty |

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |

---

### UpdateEmployeeDepartment

**HTTP:** `PUT /v1/employees/{employee_id}/department`
**gRPC:** `MembershipCommandService.UpdateEmployeeDepartment`

#### Права доступа

`AdminOf.Employee(employeeID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `employee_id` | string (UUID) | required, uuid |
| `department_id` | string (UUID) | required, uuid |

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |
| `department_not_found` | 404 | Отдел не найден |

---

### TerminateEmployee

**HTTP:** `DELETE /v1/employees/{employee_id}`
**gRPC:** `MembershipCommandService.TerminateEmployee`

#### Права доступа

`AdminOf.Employee(employeeID)`

#### Инварианты

- При увольнении автоматически очищаются все роли сотрудника и слоты заместителей в которых он числится.
- Активные отпуска закрываются.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |

---

## Отпуска

### StartVacationNow

**HTTP:** `POST /v1/employees/{employee_id}/vacations:start-now`
**gRPC:** `MembershipCommandService.StartVacationNow`

#### Права доступа

`AdminOf.Employee(employeeID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `employee_id` | string (UUID) | required, uuid |
| `ends_at` | timestamp | omitempty |

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |

---

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

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |

---

### UpdateVacationEndDate

**HTTP:** `PUT /v1/vacations/{vacation_id}/end-date`
**gRPC:** `MembershipCommandService.UpdateVacationEndDate`

#### Права доступа

`AdminOf.Vacation(vacationID)`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `vacation_id` | string (UUID) | required, uuid |
| `ends_at` | timestamp | required |

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `vacation_not_found` | 404 | Отпуск не найден |

---

### ForceEndVacation

**HTTP:** `POST /v1/vacations/{vacation_id}/terminations`
**gRPC:** `MembershipCommandService.ForceEndVacation`

#### Права доступа

`AdminOf.Vacation(vacationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `vacation_not_found` | 404 | Отпуск не найден |

---

### CancelScheduledVacation

**HTTP:** `POST /v1/vacations/{vacation_id}/cancellations`
**gRPC:** `MembershipCommandService.CancelScheduledVacation`

#### Права доступа

`AdminOf.Vacation(vacationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `vacation_not_found` | 404 | Отпуск не найден |

---

## Роль DeptResponsible

### AssignDepartmentResponsible

**HTTP:** `POST /v1/departments/{department_id}/responsibles`
**gRPC:** `MembershipCommandService.AssignDepartmentResponsible`

#### Права доступа

`AdminOf.Department(departmentID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |
| `department_responsible_already_assigned` | 409 | Ответственный за отдел уже назначен |

---

### RevokeDepartmentResponsible

**HTTP:** `DELETE /v1/departments/{department_id}/responsibles/{employee_id}`
**gRPC:** `MembershipCommandService.RevokeDepartmentResponsible`

#### Права доступа

`AdminOf.Department(departmentID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_responsible_not_found` | 404 | Ответственный за отдел не найден |

---

### AssignDepartmentResponsibleDeputy

**HTTP:** `POST /v1/departments/{department_id}/responsibles/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.AssignDepartmentResponsibleDeputy`

#### Права доступа

`AdminOf.Department(departmentID)`

#### Инварианты

- Заместитель должен быть сотрудником той же организации.
- Заместитель не может совпадать с основным держателем роли.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_responsible_not_found` | 404 | Ответственный за отдел не найден |
| `deputy_not_found` | 404 | Заместитель не найден |
| `deputy_already_assigned` | 409 | Заместитель уже назначен |

---

### RemoveDepartmentResponsibleDeputy

**HTTP:** `DELETE /v1/departments/{department_id}/responsibles/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.RemoveDepartmentResponsibleDeputy`

#### Права доступа

`AdminOf.Department(departmentID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `department_responsible_not_found` | 404 | Ответственный за отдел не найден |

---

## Роль ClinicHead

### AssignClinicHead

**HTTP:** `POST /v1/clinics/{clinic_id}/heads`
**gRPC:** `MembershipCommandService.AssignClinicHead`

#### Права доступа

`AdminOf.Clinic(clinicID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |
| `clinic_head_already_assigned` | 409 | Руководитель клиники уже назначен |

---

### RevokeClinicHead

**HTTP:** `DELETE /v1/clinics/{clinic_id}/heads/{employee_id}`
**gRPC:** `MembershipCommandService.RevokeClinicHead`

#### Права доступа

`AdminOf.Clinic(clinicID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_head_not_found` | 404 | Руководитель клиники не найден |

---

### AssignClinicHeadDeputy

**HTTP:** `POST /v1/clinics/{clinic_id}/heads/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.AssignClinicHeadDeputy`

#### Права доступа

`AdminOf.Clinic(clinicID)`

#### Инварианты

- Заместитель должен быть сотрудником той же организации.
- Заместитель не может совпадать с основным держателем роли.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_head_not_found` | 404 | Руководитель клиники не найден |
| `deputy_not_found` | 404 | Заместитель не найден |
| `deputy_already_assigned` | 409 | Заместитель уже назначен |

---

### RemoveClinicHeadDeputy

**HTTP:** `DELETE /v1/clinics/{clinic_id}/heads/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.RemoveClinicHeadDeputy`

#### Права доступа

`AdminOf.Clinic(clinicID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `clinic_head_not_found` | 404 | Руководитель клиники не найден |

---

## Роль OrgAdmin

### AssignOrganizationAdmin

**HTTP:** `POST /v1/organizations/{organization_id}/admins`
**gRPC:** `MembershipCommandService.AssignOrganizationAdmin`

#### Права доступа

`SystemAdmin`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |
| `organization_admin_already_assigned` | 409 | Администратор организации уже назначен |

---

### RevokeOrganizationAdmin

**HTTP:** `DELETE /v1/organizations/{organization_id}/admins/{employee_id}`
**gRPC:** `MembershipCommandService.RevokeOrganizationAdmin`

#### Права доступа

`SystemAdmin`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_admin_not_found` | 404 | Администратор организации не найден |

---

### AssignOrganizationAdminDeputy

**HTTP:** `POST /v1/organizations/{organization_id}/admins/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.AssignOrganizationAdminDeputy`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Инварианты

- Заместитель должен быть сотрудником той же организации.
- Заместитель не может совпадать с основным держателем роли.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_admin_not_found` | 404 | Администратор организации не найден |
| `deputy_not_found` | 404 | Заместитель не найден |
| `deputy_already_assigned` | 409 | Заместитель уже назначен |

---

### RemoveOrganizationAdminDeputy

**HTTP:** `DELETE /v1/organizations/{organization_id}/admins/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.RemoveOrganizationAdminDeputy`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_admin_not_found` | 404 | Администратор организации не найден |

---

## Роль OrgHead

### AssignOrganizationHead

**HTTP:** `POST /v1/organizations/{organization_id}/heads`
**gRPC:** `MembershipCommandService.AssignOrganizationHead`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |
| `organization_head_already_assigned` | 409 | Руководитель организации уже назначен |

---

### RevokeOrganizationHead

**HTTP:** `DELETE /v1/organizations/{organization_id}/heads/{employee_id}`
**gRPC:** `MembershipCommandService.RevokeOrganizationHead`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_head_not_found` | 404 | Руководитель организации не найден |

---

### AssignOrganizationHeadDeputy

**HTTP:** `POST /v1/organizations/{organization_id}/heads/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.AssignOrganizationHeadDeputy`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Инварианты

- Заместитель должен быть сотрудником той же организации.
- Заместитель не может совпадать с основным держателем роли.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_head_not_found` | 404 | Руководитель организации не найден |
| `deputy_not_found` | 404 | Заместитель не найден |
| `deputy_already_assigned` | 409 | Заместитель уже назначен |

---

### RemoveOrganizationHeadDeputy

**HTTP:** `DELETE /v1/organizations/{organization_id}/heads/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.RemoveOrganizationHeadDeputy`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_head_not_found` | 404 | Руководитель организации не найден |

---

## Роль OrgDispatcher

### AssignOrganizationDispatcher

**HTTP:** `POST /v1/organizations/{organization_id}/dispatchers`
**gRPC:** `MembershipCommandService.AssignOrganizationDispatcher`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `employee_not_found` | 404 | Сотрудник не найден |
| `organization_dispatcher_already_assigned` | 409 | Диспетчер организации уже назначен |

---

### RevokeOrganizationDispatcher

**HTTP:** `DELETE /v1/organizations/{organization_id}/dispatchers/{employee_id}`
**gRPC:** `MembershipCommandService.RevokeOrganizationDispatcher`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_dispatcher_not_found` | 404 | Диспетчер организации не найден |

---

### AssignOrganizationDispatcherDeputy

**HTTP:** `POST /v1/organizations/{organization_id}/dispatchers/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.AssignOrganizationDispatcherDeputy`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Инварианты

- Заместитель должен быть сотрудником той же организации.
- Заместитель не может совпадать с основным держателем роли.

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_dispatcher_not_found` | 404 | Диспетчер организации не найден |
| `deputy_not_found` | 404 | Заместитель не найден |
| `deputy_already_assigned` | 409 | Заместитель уже назначен |

---

### RemoveOrganizationDispatcherDeputy

**HTTP:** `DELETE /v1/organizations/{organization_id}/dispatchers/{employee_id}/deputy`
**gRPC:** `MembershipCommandService.RemoveOrganizationDispatcherDeputy`

#### Права доступа

`AdminOf.Organization(organizationID)`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `organization_dispatcher_not_found` | 404 | Диспетчер организации не найден |

---

## Роль SystemAdmin

### GrantSystemAdmin

**HTTP:** `POST /v1/system-admins`
**gRPC:** `MembershipCommandService.GrantSystemAdmin`

#### Права доступа

`SystemAdmin`

#### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `zitadel_user_id` | string | required, min=1 |

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `zitadel_user_not_found` | 404 | Пользователь Zitadel не найден |
| `system_admin_already_granted` | 409 | Системный администратор уже назначен |
| `zitadel_verify_failed` | 503 | Внешний сервис идентификации недоступен |

---

### RevokeSystemAdmin

**HTTP:** `DELETE /v1/system-admins/{zitadel_user_id}`
**gRPC:** `MembershipCommandService.RevokeSystemAdmin`

#### Права доступа

`SystemAdmin`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `validation_failed` | 400 | Ошибки валидации полей |
| `permission_denied` | 403 | Нет прав доступа |
| `system_admin_not_found` | 404 | Системный администратор не найден |

---

## Query-методы

### GetEmployee

**HTTP:** `GET /v1/employees/{id}`
**gRPC:** `MembershipQueryService.GetEmployee`

#### Права доступа

`ReaderOf.Employee`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |
| `employee_card_not_found` | 404 | Карточка сотрудника не найдена |

---

### ListEmployeesByDepartment

**HTTP:** `GET /v1/departments/{department_id}/employees`
**gRPC:** `MembershipQueryService.ListEmployeesByDepartment`

#### Права доступа

`ReaderOf.Department`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### ListEmployeesByClinic

**HTTP:** `GET /v1/clinics/{clinic_id}/employees`
**gRPC:** `MembershipQueryService.ListEmployeesByClinic`

#### Права доступа

`ReaderOf.Clinic`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### ListEmployeesByOrganization

**HTTP:** `GET /v1/organizations/{organization_id}/employees`
**gRPC:** `MembershipQueryService.ListEmployeesByOrganization`

#### Права доступа

`ReaderOf.Organization`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### CountEmployeesByDepartment

**HTTP:** `GET /v1/departments/{department_id}/employees:count`
**gRPC:** `MembershipQueryService.CountEmployeesByDepartment`

#### Права доступа

`ReaderOf.Department`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### CountEmployeesByClinic

**HTTP:** `GET /v1/clinics/{clinic_id}/employees:count`
**gRPC:** `MembershipQueryService.CountEmployeesByClinic`

#### Права доступа

`ReaderOf.Clinic`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### CountEmployeesByOrganization

**HTTP:** `GET /v1/organizations/{organization_id}/employees:count`
**gRPC:** `MembershipQueryService.CountEmployeesByOrganization`

#### Права доступа

`ReaderOf.Organization`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### SearchEmployeesByOrganization

**HTTP:** `GET /v1/organizations/{organization_id}/employees:search`
**gRPC:** `MembershipQueryService.SearchEmployeesByOrganization`

#### Права доступа

`ReaderOf.Organization`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### ListVacationsByEmployee

**HTTP:** `GET /v1/employees/{employee_id}/vacations`
**gRPC:** `MembershipQueryService.ListVacationsByEmployee`

#### Права доступа

`SystemAdmin` + `OrgAdminOf.Employee` + `SelfEmployee`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### CountVacationsByEmployee

**HTTP:** `GET /v1/employees/{employee_id}/vacations:count`
**gRPC:** `MembershipQueryService.CountVacationsByEmployee`

#### Права доступа

`SystemAdmin` + `OrgAdminOf.Employee` + `SelfEmployee`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### GetClinicHead

**HTTP:** `GET /v1/clinics/{clinic_id}/head`
**gRPC:** `MembershipQueryService.GetClinicHead`

#### Права доступа

`ReaderOf.Clinic`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### GetDepartmentResponsible

**HTTP:** `GET /v1/departments/{department_id}/responsible`
**gRPC:** `MembershipQueryService.GetDepartmentResponsible`

#### Права доступа

`ReaderOf.Department`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### ListOrgAdmins

**HTTP:** `GET /v1/organizations/{organization_id}/admins`
**gRPC:** `MembershipQueryService.ListOrgAdmins`

#### Права доступа

`ReaderOf.Organization`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### ListOrgDispatchers

**HTTP:** `GET /v1/organizations/{organization_id}/dispatchers`
**gRPC:** `MembershipQueryService.ListOrgDispatchers`

#### Права доступа

`ReaderOf.Organization`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### ListOrgHeads

**HTTP:** `GET /v1/organizations/{organization_id}/heads`
**gRPC:** `MembershipQueryService.ListOrgHeads`

#### Права доступа

`ReaderOf.Organization`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

---

### ListSystemAdmins

**HTTP:** `GET /v1/system-admins`
**gRPC:** `MembershipQueryService.ListSystemAdmins`

#### Права доступа

`SystemAdmin`

#### Ошибки

| Код | HTTP | Описание |
|---|---|---|
| `permission_denied` | 403 | Нет прав доступа |

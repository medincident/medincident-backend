# Роли и права

## Ролевая модель

Все роли — организационные. Нет глобальных ролей кроме `SystemAdmin`.

| Роль | Область | Описание |
|---|---|---|
| `SystemAdmin` | Глобально | Обходит все scope-проверки |
| `OrgAdmin` | Организация | Полное управление орг. структурой и сотрудниками |
| `OrgHead` | Организация | Руководитель организации (с заместителем) |
| `OrgDispatcher` | Организация | Диспетчер — принимает заявки пациентов |
| `ClinicHead` | Клиника | Руководитель клиники (с заместителем) |
| `DeptResponsible` | Отдел | Ответственный за отдел (с заместителем) |
| `MemberOf` | Организация | Любой сотрудник организации (read-only доступ) |

## Иерархия scope-проверок

```
SystemAdmin
    └── OrgAdmin (Organization / Clinic / Department / Employee / Vacation / Category / IncidentType)
            └── MemberOf (Organization / Clinic / Department / Employee / Category / IncidentType)
```

`AdminOf.X(id)` — сокращение для `AnyOf(SystemAdmin, OrgAdminOf.X(id))`.
`ReaderOf.X(id)` — сокращение для `AnyOf(SystemAdmin, OrgAdminOf.X(id), MemberOf.X(id))`.

## Заместители

Роли `OrgAdmin`, `OrgHead`, `OrgDispatcher`, `ClinicHead`, `DeptResponsible` поддерживают заместителя.

**Правила:**

- Заместитель активируется только пока основной держатель роли находится в активном отпуске (`employee_vacations.starts_at <= now() AND (ends_at IS NULL OR ends_at > now())`).
- Заместитель не активируется, если он сам находится в отпуске — делегирование не каскадируется.
- Слот заместителя существует один на роль; старый заместитель заменяется при переназначении.
- При увольнении основного держателя слот заместителя очищается.

**Пример:** `OrgAdmin` назначил заместителя. Пока `OrgAdmin` в отпуске, заместитель может выполнять все операции OrgAdmin уровня. Если заместитель тоже уходит в отпуск — доступ прекращается.

## Политика видимости организаций

Организации — публичный справочник: любой аутентифицированный пользователь (валидный Bearer-токен) может их перечислять и читать без дополнительной проверки роли. Клиники, отделы и сотрудники требуют `ReaderOf.*`.

> **Предупреждение:** если когда-либо появятся приватные или черновые организации, `OrganizationReader` потребует `*authz.Authz`. Это ломающее изменение по всему коду.

## Особые политики

| Политика | Когда применяется |
|---|---|
| `Authenticated` | Публичные endpoint'ы — только проверка JWT |
| `SelfEmployee(id)` | Сотрудник читает свои собственные данные |
| `SelfSession(id)` | Владелец Zitadel-сессии читает свою сессию |

## Где живёт код

- Политики: [`internal/service/authz/policy.go`](../internal/service/authz/policy.go)
- Роли: `org_admins.go`, `org_head.go`, `clinic_head.go`, `dept_responsible.go`, `org_dispatcher.go`
- Применение: `s.authz.Require(ctx, cmd.Caller.ZitadelUserID, authz.AdminOf.X(scopeID))`

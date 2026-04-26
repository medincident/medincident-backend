← [Документация](../README.md)

# Аутентификация и авторизация

## Аутентификация

Используется [Zitadel](https://zitadel.com/) как IdP. Каждый gRPC-запрос должен содержать Bearer-токен в метаданных `authorization`.

**Цепочка проверки:**

1. gRPC-интерсептор `authn.go` извлекает Bearer-токен из метаданных запроса.
2. Токен верифицируется через Zitadel introspection endpoint.
3. При успехе в контекст записывается `callerID` (Zitadel user ID).
4. При неуспехе — запрос отклоняется с `codes.Unauthenticated`.

Хендлеры извлекают `callerID` через `grpcmw.CallerID(ctx)` и формируют `authz.Caller{ZitadelUserID: callerID}`.

## Авторизация

Авторизация — RBAC с scope-изоляцией. Реализована в [`internal/service/authz/`](../../internal/service/authz/).

### Принцип работы

Каждый `authz.Require(ctx, callerID, policy)`:
1. Рендерит policy в SQL-ветки (`branches`).
2. Объединяет ветки в `SELECT EXISTS(branch1 UNION ALL branch2 ...)`.
3. Выполняет один DB-запрос.
4. Возвращает `nil` при `true`, `permission_denied` при `false`.

Неавторизованный вызывающий **не может отличить** "scope не существует" от "нет прав": оба случая возвращают `permission_denied` с одним и тем же сообщением.

### Иерархия политик

```
Policy (interface)
├── authenticatedPolicy     — любой аутентифицированный (SELECT 1)
├── sysAdminPolicy          — domain.system_admins
├── anyOfPolicy             — UNION ALL нескольких политик
├── orgAdminPolicy          — domain.org_admins (прямой + заместитель)
├── orgHeadOfOrgPolicy      — domain.org_heads (прямой + заместитель)
├── orgDispatcherOfOrgPolicy — domain.org_dispatchers (прямой + заместитель)
├── clinicHeadOfClinicPolicy — domain.clinic_heads (прямой + заместитель)
├── deptResponsibleOfDeptPolicy — domain.department_responsibles (прямой + заместитель)
├── memberOfPolicy          — domain.employees (только прямой)
├── selfEmployeePolicy      — конкретный employee по ID
└── selfSessionPolicy       — конкретная сессия по ID
```

**Батареи** (сокращения для типичных комбинаций):
- `AdminOf.X(id)` → `AnyOf(SystemAdmin, OrgAdminOf.X(id))`
- `ReaderOf.X(id)` → `AnyOf(SystemAdmin, OrgAdminOf.X(id), MemberOf.X(id))`

### Порядок выполнения в сервисном методе

```
validation.Struct(cmd.Payload)   // 1. валидация
uuid.MustParse(...)               // 2. парсинг UUID
s.authz.Require(ctx, ...)         // 3. проверка прав
// бизнес-логика                  // 4. выполнение
```

### Коды ошибок

| Код | gRPC-статус | Когда |
|---|---|---|
| `permission_denied` | `PermissionDenied` | Политика вернула false |
| `authz_check_failed` | `Internal` | Ошибка DB при проверке |

### Защита от zero UUID

Если в политику передан `uuid.Nil` — Require возвращает `authz_check_failed` без DB-запроса. Это защищает от ситуации, когда whitespace-баг на транспортном уровне превращает невалидный UUID в запрос, который "молча всё запрещает".

## Межсервисная аутентификация

`command-server` и `query-server` доверяют `gateway-server` без дополнительной верификации (все в одной сети). Каждый входящий запрос к command/query всё равно проходит через authn-интерсептор: Bearer-токен клиента проксируется gateway без изменений.

## Безопасность событий в NATS

События в NATS JetStream считаются доверенными: они пишутся только через `outbox.AppendEvent` в той же транзакции, что и доменная мутация. Consumer'ы на query-стороне принимают их без дополнительной верификации.

## Scope-изоляция

- Организации изолированы друг от друга: `OrgAdmin` одной организации не видит данные другой.
- Все проверки реализованы через EXISTS-запросы с JOIN на целевую сущность через цепочку FK — утечка данных через timing-атаки невозможна.
- Повышение прав через границы организаций невозможно на уровне SQL.

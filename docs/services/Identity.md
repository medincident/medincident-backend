← [Документация](../README.md)

# Идентификация

Сервис для работы с Zitadel-сессиями.

**gRPC:** `IdentityQueryService` (query)

---

## GetSession

**gRPC:** `IdentityQueryService.GetSession`

### Права доступа

`AnyOf(SystemAdmin, SelfSession(sessionID))` — системный администратор или владелец сессии.

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `session_id` | string | required |

### Описание

Возвращает данные Zitadel-сессии по ID. Непрозрачная обёртка над Zitadel API — возвращает то, что возвращает Zitadel, без дополнительной бизнес-логики.

### Ошибки

| Код | Описание |
|---|---|
| `permission_denied` | Нет прав (не владелец и не SystemAdmin) |
| `session_not_found` | Сессия не найдена |

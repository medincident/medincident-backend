← [Документация](../README.md)

# Статистика

Агрегированная статистика по организации.

**gRPC:** `StatsQueryService` (query)

---

## GetOrganizationStats

**gRPC:** `StatsQueryService.GetOrganizationStats`

### Права доступа

`ReaderOf.Organization(organizationID)`

### Параметры

| Поле | Тип | Правила |
|---|---|---|
| `organization_id` | string (UUID) | required, uuid |

### Ответ

Возвращает агрегированные счётчики по организации:

| Поле | Описание |
|---|---|
| `total_incidents` | Всего инцидентов |
| `open_incidents` | Инцидентов в статусе `open` |
| `closed_incidents` | Инцидентов в статусе `closed` |
| `cancelled_incidents` | Инцидентов в статусе `cancelled` |
| `total_employees` | Всего сотрудников |
| `total_clinics` | Всего клиник |
| `total_departments` | Всего отделов |

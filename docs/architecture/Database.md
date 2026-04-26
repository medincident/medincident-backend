← [Документация](../README.md)

# База данных

PostgreSQL 17+. Три схемы с чёткими ролями.

## Схемы

### `domain.*` — источник правды

Транзакционные данные command-стороны. Только NOT NULL + PK + FK, без CHECK-констрейнтов (валидация на сервисном уровне).

| Таблица | Назначение |
|---|---|
| `domain.organizations` | Организации |
| `domain.clinics` | Клиники (FK → organizations) |
| `domain.departments` | Отделы (FK → clinics) |
| `domain.employees` | Сотрудники (FK → organizations) |
| `domain.employee_vacations` | Отпуска сотрудников |
| `domain.org_admins` | Роль OrgAdmin + слот заместителя |
| `domain.org_heads` | Роль OrgHead + слот заместителя |
| `domain.org_dispatchers` | Роль OrgDispatcher + слот заместителя |
| `domain.clinic_heads` | Роль ClinicHead + слот заместителя |
| `domain.department_responsibles` | Роль DeptResponsible + слот заместителя |
| `domain.system_admins` | Системные администраторы |
| `domain.incident_categories` | Категории инцидентов (FK → organizations) |
| `domain.incident_types` | Типы инцидентов (FK → categories) |
| `domain.incidents` | Инциденты |
| `domain.patient_incident_buffer` | Буфер заявок пациентов |

### `outbox.*` — транзакционный outbox

| Таблица | Назначение |
|---|---|
| `outbox.events` | Очередь исходящих событий |

Колонки: `id, subject, payload, headers, created_at, published_at`.
- `subject` — NATS subject вида `medincident.event.<aggregate>.v1.<action>`
- `payload` — сериализованный `event.v1.Envelope` (protobuf binary)
- `headers` — заполняется publisher'ом (dedup-ключи)
- `published_at` — выставляется publisher'ом после успешной публикации

Command-сторона пишет только `subject` и `payload`.

### `projections.*` — read-модели

Строятся NATS consumer'ами на query-стороне. Не используются command-стороной.

| Таблица | Назначение |
|---|---|
| `projections.organizations` | Проекция организаций |
| `projections.clinics` | Проекция клиник |
| `projections.departments` | Проекция отделов |
| `projections.employees` | Проекция сотрудников |
| `projections.org_admins` | Проекция ролей OrgAdmin |
| `projections.org_heads` | Проекция ролей OrgHead |
| `projections.org_dispatchers` | Проекция ролей OrgDispatcher |
| `projections.clinic_heads` | Проекция ролей ClinicHead |
| `projections.department_responsibles` | Проекция ролей DeptResponsible |
| `projections.system_admins` | Проекция системных администраторов |
| `projections.incident_categories` | Проекция категорий инцидентов |
| `projections.incident_types` | Проекция типов инцидентов |
| `projections.incidents` | Проекция инцидентов |
| `projections.patient_incident_buffer` | Проекция буфера пациента |
| `projections.incident_status_history` | История смен статуса |
| `projections.incident_priority_history` | История смен приоритета |

## Миграции

Управляются через `dbmate`. Файлы в `db/migrations/` **только** создаются командой:

```bash
task migrate:new -- <snake_case_name>
```

Никогда не создавать файлы миграций вручную.

```bash
task migrate   # dbmate up — применить все pending-миграции
```

## Инварианты

- **FK — immutable.** Клиника всегда принадлежит одной организации, отдел — одной клинике. После создания FK не меняется.
- **Нет CHECK-констрейнтов.** Все бизнес-правила проверяются на сервисном уровне через `validation.Struct`.
- **gorm field-level permissions.** PK и родительские FK помечены `gorm:"<-:create"`. Мутабельные поля — `gorm:"<-"`.

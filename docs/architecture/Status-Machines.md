← [Документация](../README.md)

# Статусные машины

## Инцидент (`domain.incidents`)

### Статусы

| Статус | Описание |
|---|---|
| `open` | Открыт, в работе |
| `closed` | Закрыт |
| `cancelled` | Отменён |

### Переходы

```
         ┌──────────────┐
         ▼              │
[open] ──► [closed]    reopen
  │         │
  │         └──► [open]  (reopen)
  │
  └──► [cancelled]
```

| Из | В | Метод |
|---|---|---|
| `open` | `closed` | `UpdateIncidentStatus` (status=closed) |
| `open` | `cancelled` | `CancelIncident` |
| `closed` | `open` | `ReopenIncident` |

`cancelled` — терминальный статус, переходов из него нет.

### Приоритеты

Приоритет инцидента изменяется независимо от статуса через `UpdateIncidentStatus` (поле priority).

| Приоритет | Описание |
|---|---|
| `low` | Низкий |
| `medium` | Средний |
| `high` | Высокий |
| `critical` | Критический |

## Буфер пациента (`domain.patient_incident_buffer`)

Буфер — промежуточный этап перед созданием инцидента. Пациент подаёт заявку, диспетчер обрабатывает её.

### Статусы

| Статус | Описание |
|---|---|
| `pending` | Ожидает обработки диспетчером |
| `published` | Преобразована в инцидент |
| `rejected` | Отклонена диспетчером |
| `cancelled` | Отменена пациентом |

### Переходы

```
[pending] ──► [published]   (PublishPatientIncident)
    │
    ├──► [rejected]         (RejectPatientIncident)
    │
    └──► [cancelled]        (CancelPatientIncident)
```

| Из | В | Метод | Инициатор |
|---|---|---|---|
| `pending` | `published` | `PublishPatientIncident` | Диспетчер |
| `pending` | `rejected` | `RejectPatientIncident` | Диспетчер |
| `pending` | `cancelled` | `CancelPatientIncident` | Пациент |

Все три финальных статуса (`published`, `rejected`, `cancelled`) — терминальные.

При `published` создаётся запись в `domain.incidents` и `source_buffer_id` инцидента ссылается на запись буфера.

## Классификатор инцидентов

### Статусы категорий и типов

| Статус | Описание |
|---|---|
| `active` | Доступна для выбора |
| `inactive` | Деактивирована, недоступна для новых инцидентов |

Деактивированная категория/тип не удаляется — исторические инциденты сохраняют ссылку. Реактивация возможна.

Удаление (`DeleteIncidentCategory`, `DeleteIncidentType`) возможно только если к сущности нет прикреплённых инцидентов (FK RESTRICT).

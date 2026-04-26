# medincident-backend — документация

Техническая документация backend-сервиса medincident. Язык — русский.

## Архитектура

| | |
|---|---|
| [Обзор](architecture/Overview.md) | CQRS-разбивка, три сервиса, транспорт |
| [Роли и права](architecture/Roles.md) | Ролевая модель, иерархия, заместители |
| [Аутентификация и авторизация](architecture/Auth.md) | Zitadel, JWT, политики authz |
| [База данных](architecture/Database.md) | Схемы domain / outbox / projections |
| [Статусные машины](architecture/Status-Machines.md) | Инциденты, буфер пациента |

## Сервисы

| | |
|---|---|
| [Орг. структура](services/OrgStructure.md) | Organization / Clinic / Department |
| [Сотрудники и роли](services/Membership.md) | Employee, назначение ролей, отпуска |
| [Классификатор инцидентов](services/incident/Classifier.md) | Категории и типы |
| [Инциденты](services/incident/Incidents.md) | Создание, статусы, приоритеты |
| [Буфер пациента](services/incident/Buffer.md) | Заявки пациентов до публикации |
| [Классификатор заявок](services/request/Classifier.md) | Категории и типы заявок |
| [Заявки](services/request/Requests.md) | ServiceRequest (в разработке) |
| [Идентификация](services/Identity.md) | Сессии Zitadel |
| [Статистика](services/Stats.md) | Агрегированная статистика |

## API Reference

| | |
|---|---|
| [gRPC (proto)](api/Proto.md) | Все RPC и сообщения — автогенерация |
| [HTTP (OpenAPI)](api/HTTP.md) | REST endpoints — автогенерация |

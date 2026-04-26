# medincident-backend

Техническая документация backend-сервиса medincident.

## Архитектура

| Страница | Описание |
|---|---|
| [[Обзор\|Architecture-Overview]] | CQRS-разбивка, три сервиса, транспорт |
| [[Роли и права\|Architecture-Roles]] | Ролевая модель, иерархия, заместители |
| [[Аутентификация и авторизация\|Architecture-Auth]] | Zitadel, JWT, политики authz |
| [[База данных\|Architecture-Database]] | Схемы domain / outbox / projections |
| [[Статусные машины\|Architecture-Status-Machines]] | Инциденты, буфер пациента |

## Сервисы

| Страница | Описание |
|---|---|
| [[Орг. структура\|Service-OrgStructure]] | Organization / Clinic / Department |
| [[Сотрудники и роли\|Service-Membership]] | Employee, назначение ролей, отпуска |
| [[Классификатор инцидентов\|Service-Incident-Classifier]] | Категории и типы |
| [[Инциденты\|Service-Incident-Incidents]] | Создание, статусы, приоритеты |
| [[Буфер пациента\|Service-Incident-Buffer]] | Заявки пациентов до публикации |
| [[Классификатор заявок\|Service-Request-Classifier]] | Категории и типы заявок |
| [[Заявки\|Service-Request-Requests]] | ServiceRequest (в разработке) |
| [[Идентификация\|Service-Identity]] | Сессии Zitadel |
| [[Статистика\|Service-Stats]] | Агрегированная статистика |

## API Reference

| Страница | Описание |
|---|---|
| [[gRPC (proto)\|API-Proto]] | Все RPC и сообщения — автогенерация |
| [[HTTP (OpenAPI)\|API-HTTP]] | REST endpoints — автогенерация |

← [Документация](../README.md)

# Обзор архитектуры

## CQRS-разбивка

Система разделена на три самостоятельных бинаря, живущих в одном репозитории:

| Сервис | Роль | Порт по умолчанию |
|---|---|---|
| `command-server` | Запись: gRPC, бизнес-логика, outbox | `:8080` |
| `query-server` | Чтение: gRPC, проекции | `:8080` |
| `gateway-server` | HTTP-фасад: grpc-gateway mux | `:8080` |

Все три сервиса дефолтно слушают `:8080` внутри своего контейнера. Внешний маппинг портов — зона ответственности операций.

## Поток данных

```
Клиент
  │
  ├─ HTTP ──► gateway-server (grpc-gateway)
  │               │
  │    ┌──────────┤
  │    ▼          ▼
  │ command-server   query-server
  │    │               ▲
  │    ▼               │
  │ PostgreSQL ─────── │  (domain.* / projections.*)
  │    │               │
  │    ▼           NATS JetStream
  │ outbox.events ──► publisher ──► NATS ──► query-server
  │
  └─ gRPC ──► command-server / query-server (напрямую)
```

## Command-side

- Принимает gRPC-команды через `command-server`.
- Каждый метод: валидация → парсинг UUID → проверка прав → бизнес-логика → транзакция (domain + outbox).
- Никогда не читает проекции и не знает о query-стороне.
- Исходящие события пишутся в `outbox.events` в той же транзакции, что и мутация агрегата. Отдельный publisher-сервис драйнит outbox и публикует envelope'ы в NATS JetStream.

## Query-side

- Принимает gRPC-запросы через `query-server`.
- Читает из `projections.*` таблиц, которые строятся consumer'ами NATS.
- Авторизация: большинство endpoint'ов требуют роли (`ReaderOf.*`), организации — доступны любому аутентифицированному пользователю.

## Транспорт

- **gRPC** — основной транспорт между сервисами и для прямых клиентов.
- **HTTP/JSON** — через `gateway-server` (grpc-gateway). Не все RPC экспонируются в HTTP; список — в [`docs/api/HTTP.md`](API-HTTP).
- **NATS JetStream** — асинхронная доставка доменных событий от command к query.

## Объектное хранилище (Garage S3)

`command-server` подключается к [Garage](https://garagehq.deuxfleurs.fr/) — S3-совместимому объектному хранилищу. Клиент создаётся через `bootstrap.OpenGarage` на основе `config.GarageConfig` и использует AWS SDK v2 со статическими credentials и path-style адресацией (требование Garage).

Подробнее: [`docs/architecture/Garage.md`](Garage.md)

## Технологический стек

- Go 1.26.3, gRPC, grpc-gateway
- gorm v2 + PostgreSQL 17
- Garage (S3-совместимое объектное хранилище) через AWS SDK v2
- buf + protoc-gen-go + protoc-gen-grpc-gateway + protoc-gen-openapiv2 + protoc-gen-doc
- dbmate для миграций
- zerolog для логирования
- samber/oops для доменных ошибок
- go-playground/validator для валидации команд
- Zitadel для аутентификации (JWT introspection)

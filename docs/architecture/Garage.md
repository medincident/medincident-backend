<- [Документация](../README.md) | [Обзор архитектуры](Overview.md)

# Garage S3 — объектное хранилище

## Назначение

[Garage](https://garagehq.deuxfleurs.fr/) — self-hosted S3-совместимое объектное хранилище. Используется для хранения файловых вложений, медиа и других бинарных объектов.

## Подключение

Клиент создаётся в `internal/bootstrap/garage.go` функцией `OpenGarage`. Используется AWS SDK v2 со следующими особенностями:

- **Статические credentials** — access key и secret key передаются через конфигурацию.
- **Path-style адресация** — обязательно для Garage (`UsePathStyle: true`).
- **Без health check при старте** — клиент возвращается немедленно (правило 13: bootstrap-хелперы не делают Ping/warm-up).

## Конфигурация

Блок `garage` опциональный — отсутствие блока отключает S3-wiring и (для `gateway-server`) пропускает проверку Garage в `/readyz`. Если блок присутствует, все его поля валидируются.

Блок `garage` в YAML-конфигурации `command-server` и `gateway-server`:

```yaml
garage:
  endpoint: "http://garage.internal:3900"
  access_key_id: "GKxxxxxxxxxxxx"
  secret_access_key: "xxxxxxxxxxxxxxxxxxxx"
  region: "garage"
  bucket: "medincident"
```

| Поле | Описание |
|---|---|
| `endpoint` | URL Garage S3 API (обязательно, валидируется как URL) |
| `access_key_id` | Ключ доступа Garage (обязательно) |
| `secret_access_key` | Секретный ключ Garage (обязательно) |
| `region` | Регион (Garage игнорирует, но AWS SDK требует) |
| `bucket` | Имя бакета для хранения объектов (обязательно) |

Секреты рекомендуется передавать через переменные окружения (`${GARAGE_ACCESS_KEY_ID}`, `${GARAGE_SECRET_ACCESS_KEY}`).

## Readiness probe

`gateway-server` выполняет `HeadBucket` против настроенного бакета на `/readyz`. При недоступности Garage эндпоинт возвращает `503 Service Unavailable` с ключом `"garage": "UNAVAILABLE"` в JSON-теле. Если блок `garage` отсутствует в конфиге gateway-server, ключ `garage` в ответе не появляется и проверка не выполняется. Таймаут одной проверки — 2 секунды (см. `garageProbeTimeout` в `internal/handler/gateway/health.go`).

## Файлы

- `internal/config/garage.go` — структура конфигурации `GarageConfig`
- `internal/bootstrap/garage.go` — создание S3-клиента `OpenGarage`
- `internal/handler/gateway/health.go` — readiness-проба Garage в `/readyz`
- `cmd/command-server/config.go` — подключение к конфигурации command-server
- `cmd/command-server/main.go` — wiring в точке входа command-server
- `cmd/gateway-server/config.go` — подключение к конфигурации gateway-server
- `cmd/gateway-server/main.go` — wiring + readiness в точке входа gateway-server

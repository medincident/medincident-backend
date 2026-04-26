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

Блок `garage` в YAML-конфигурации `command-server`:

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

## Файлы

- `internal/config/garage.go` — структура конфигурации `GarageConfig`
- `internal/bootstrap/garage.go` — создание S3-клиента `OpenGarage`
- `cmd/command-server/config.go` — подключение к конфигурации command-server
- `cmd/command-server/main.go` — wiring в точке входа

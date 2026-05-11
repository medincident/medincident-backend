# Быстрый старт

## Настройка NATS JetStream стрима

Перед первым запуском необходимо создать стрим в NATS. Убедитесь, что
`nats` CLI установлен и NATS-сервер доступен.

```bash
nats stream add medincident.events \
  --subjects "medincident.event.>" \
  --storage file \
  --retention limits \
  --replicas 1 \
  --max-msgs -1 \
  --max-bytes -1 \
  --max-age 0 \
  --dupe-window 2m \
  --no-allow-rollup \
  --deny-delete \
  --deny-purge
```

Для локальной разработки (одна реплика, без дедупликации):

```bash
nats stream add medincident.events \
  --subjects "medincident.event.>" \
  --storage memory \
  --replicas 1
```

После создания стрима запустите сервисы в следующем порядке:
1. PostgreSQL (command DB + query DB)
2. NATS JetStream
3. `command-server`
4. `publisher-server`
5. `query-server`
6. `gateway-server`

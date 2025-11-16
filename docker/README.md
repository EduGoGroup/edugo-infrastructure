# docker - Módulo de Docker Compose

Docker Compose con perfiles para diferentes necesidades.

## Perfiles Disponibles

| Perfil | Servicios | Uso |
|--------|-----------|-----|
| **(default)** | PostgreSQL, MongoDB | Desarrollo básico |
| `messaging` | + RabbitMQ | APIs con eventos |
| `cache` | + Redis | Si necesitas caché |
| `tools` | + PgAdmin, Mongo Express | Debugging visual |

## Ejemplos de Uso

```bash
# Solo core (PostgreSQL + MongoDB)
docker-compose up -d

# Core + RabbitMQ
docker-compose --profile messaging up -d

# Core + RabbitMQ + Redis
docker-compose --profile messaging --profile cache up -d

# Todo + herramientas
docker-compose --profile messaging --profile cache --profile tools up -d
```

## Variables de Entorno

Ver `.env.example` en la raíz del proyecto.

## Desde Otros Proyectos

```makefile
# api-mobile/Makefile
INFRA_PATH := ../edugo-infrastructure

.PHONY: dev-setup
dev-setup:
	@cd $(INFRA_PATH)/docker && docker-compose --profile messaging up -d
```

# docker processes

## Procesos propios de la superficie

### 1. Levantar servicios core

Sin perfiles adicionales, la topologia arranca:

- PostgreSQL
- MongoDB

### 2. Activar mensajeria

Con profile `messaging` se agrega RabbitMQ.

### 3. Activar cache

Con profile `cache` se agrega Redis.

### 4. Activar herramientas visuales

Con profile `tools` se agregan:

- PgAdmin
- Mongo Express

### 5. Aislar servicios en red local

Todos los servicios comparten la red `edugo-network` y persisten datos en volumes dedicados.

## Realidades que importan documentar

- La configuracion local se apoya en `.env.example` desde la raiz del repo.
- `scripts/dev-setup.sh` usa esta superficie, pero su paso de migracion automatica no esta alineado con la estructura actual del repo.

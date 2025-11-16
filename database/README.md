# database - Módulo de Migraciones

Centraliza todas las migraciones de PostgreSQL y MongoDB.

## PostgreSQL - Migraciones

### Uso

```bash
cd database
go run migrate.go up      # Ejecutar migraciones
go run migrate.go down    # Revertir última
go run migrate.go status  # Ver estado
```

### Crear Nueva Migración

```bash
go run migrate.go create "add_avatar_to_users"
```

Genera:
- `migrations/postgres/00X_add_avatar_to_users.up.sql`
- `migrations/postgres/00X_add_avatar_to_users.down.sql`

## MongoDB - Migraciones

### Uso

```bash
cd database
go run mongodb_migrate.go up      # Ejecutar migraciones
go run mongodb_migrate.go down    # Revertir última
go run mongodb_migrate.go status  # Ver estado
```

### Crear Nueva Migración

```bash
go run mongodb_migrate.go create "add_new_collection"
```

Genera:
- `migrations/mongodb/00X_add_new_collection.up.js`
- `migrations/mongodb/00X_add_new_collection.down.js`

### Variables de Entorno

```bash
MONGO_HOST=localhost       # default: localhost
MONGO_PORT=27017          # default: 27017
MONGO_DB_NAME=edugo       # default: edugo
MONGO_USER=               # opcional
MONGO_PASSWORD=           # opcional
```

## Documentación

- **PostgreSQL Schema:** Ver `TABLE_OWNERSHIP.md`
- **MongoDB Schema:** Ver `MONGODB_SCHEMA.md`

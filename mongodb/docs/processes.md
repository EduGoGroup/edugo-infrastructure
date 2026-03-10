# mongodb processes

## Procesos propios del modulo

### 1. Crear estructura documental

El paquete `mongodb/migrations` define tres funciones de estructura internas que crean collections con validator JSON Schema:

- `material_summary`
- `material_assessment_worker`
- `material_event`

La API publica expone:

- `ApplyAll(ctx, db)`
- `ApplyStructure(ctx, db)`
- `ApplyConstraints(ctx, db)`
- `ApplySeeds(ctx, db)`
- `ApplyMockData(ctx, db)`
- `ListFunctions()`

### 2. Crear indices y TTL

El modulo tambien define indices para cada collection. En `material_event` existe un TTL de 90 dias sobre `created_at`.

### 3. Sembrar documentos canonicos

`migrations/seeds.go` carga datos de referencia en:

- `material_assessment_worker`
- `material_summary`

Los comentarios del propio archivo alinean esos documentos con materiales sembrados en Postgres, pero la descripcion detallada de integracion queda fuera de fase 1.

### 4. Sembrar mock data

`migrations/mock_data.go` carga documentos de prueba para desarrollo y testing sobre las mismas collections principales.

### 5. Mantener estado de migraciones

`cmd/migrate/migrate.go` no ejecuta toda la estructura. Su responsabilidad actual es:

- asegurar la collection `schema_migrations`
- mostrar estado
- forzar una version

El propio help del comando deja claro que las migraciones reales se aplican desde Go usando el paquete `migrations`.

### 6. Exponer entities Go

El modulo publica tres entities:

- `MaterialAssessment`
- `MaterialSummary`
- `MaterialEvent`

Estas structs reflejan el shape actual de las collections activas.

### 7. Tests del modulo

El paquete `mongodb/migrations` tiene tests de integracion y los tests cortos pasan en esta fase.

## Realidades que importan documentar

- El centro del modulo es codigo Go embebido, no scripts JS externos.
- Aun existen artefactos `_deprecated/` y un script JS heredado en `migrations/001_setup_collections.js`, pero no son la superficie recomendada.
- `mongodb/Makefile` tiene targets `migrate-*` que apuntan a `go run migrate.go`, ruta que hoy no existe en la raiz del modulo.

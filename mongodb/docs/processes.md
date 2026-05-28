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

### 5. Operacion por CLI

El modulo tiene dos entrypoints operativos:

- `cmd/runner/main.go`: ejecuta estructura + constraints embebidos; comandos: `all`, `structure`, `constraints`
- `cmd/seed/main.go`: ejecuta solo seeds embebidos; comandos: `all`, `canonical`, `mock`

### 6. Exponer entities Go

El modulo publica tres entities:

- `MaterialAssessment`
- `MaterialSummary`
- `MaterialEvent`

Estas structs reflejan el shape actual de las collections activas.

### 7. Tests del modulo

El paquete `mongodb/migrations` tiene tests de integracion y los tests cortos pasan en esta fase.


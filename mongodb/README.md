# MongoDB Migrations - EduGo Infrastructure

## Arquitectura de Migraciones

Este proyecto utiliza una arquitectura de migraciones separada por responsabilidades para MongoDB:

### Estructura de Directorios

```
mongodb/
├── structure/          # Definiciones de esquemas y colecciones
│   ├── 001_material_assessment.go
│   ├── 002_material_content.go
│   ├── ...
│   └── 009_material_event.go
│
├── constraints/        # Índices y validaciones
│   ├── 001_material_assessment_indexes.go
│   ├── 002_material_content_indexes.go
│   ├── ...
│   └── 009_material_event_indexes.go
│
├── seeds/             # Datos iniciales (producción)
│   └── (vacío por ahora)
│
├── testing/           # Datos de prueba (desarrollo)
│   └── (vacío por ahora)
│
├── runner.go          # CLI para ejecutar migraciones
├── migrate.go         # (Legacy - será removido)
└── go.mod
```

### Filosofía de Separación

**¿Por qué separamos las migraciones?**

1. **Structure** - Define la estructura base de las colecciones
   - Schemas con validación JSON
   - Campos requeridos y opcionales
   - Tipos de datos y validaciones básicas
   
2. **Constraints** - Define optimizaciones y reglas adicionales
   - Índices para consultas rápidas
   - Índices únicos para integridad
   - Índices TTL para limpieza automática
   - Índices de texto para búsquedas

3. **Seeds** - Datos iniciales necesarios en producción
   - Configuraciones del sistema
   - Datos maestros
   - Valores por defecto

4. **Testing** - Datos para desarrollo y pruebas
   - Datos de ejemplo
   - Casos de prueba
   - Escenarios de desarrollo

### Ventajas de esta Arquitectura

✅ **Granularidad**: Puedes ejecutar solo lo que necesitas
✅ **Claridad**: Cada archivo tiene una responsabilidad única
✅ **Mantenimiento**: Fácil de entender y modificar
✅ **Testing**: Puedes probar estructura sin datos
✅ **Performance**: Índices separados permiten análisis de rendimiento

## Uso

### Ejecutar todas las migraciones

```bash
go run runner.go all
```

### Ejecutar solo estructura (schemas)

```bash
go run runner.go structure
```

### Ejecutar solo constraints (índices)

```bash
go run runner.go constraints
```

## Variables de Entorno

El runner utiliza las siguientes variables de entorno:

```bash
MONGO_HOST=localhost        # Host de MongoDB
MONGO_PORT=27017           # Puerto de MongoDB
MONGO_USER=edugo           # Usuario de MongoDB
MONGO_PASSWORD=edugo123    # Contraseña de MongoDB
MONGO_DB_NAME=edugo        # Nombre de la base de datos
```

## Colecciones Actuales

1. **material_assessment** - Evaluaciones generadas por IA
2. **material_content** - Contenido extraído de materiales
3. **assessment_attempt_result** - Resultados de intentos de evaluación
4. **audit_logs** - Registro de auditoría del sistema (TTL: 90 días)
5. **notifications** - Notificaciones de usuarios (TTL: configurable)
6. **analytics_events** - Eventos de analítica (TTL: 365 días)
7. **material_summary** - Resúmenes generados por IA
8. **material_assessment_worker** - Evaluaciones procesadas por worker
9. **material_event** - Cola de eventos para procesamiento (TTL: 90 días)

## Notas Importantes

### Índices TTL (Time To Live)

Algunas colecciones tienen índices TTL configurados para limpieza automática:

- **audit_logs**: 90 días
- **analytics_events**: 365 días
- **material_event**: 90 días
- **notifications**: 
  - `expires_at`: expiración inmediata cuando se setea
  - `archived_at`: 30 días después de archivar

### Validación de Schemas

Todas las colecciones tienen validación JSON Schema estricta. Los documentos que no cumplan con el schema serán rechazados por MongoDB.

### Idempotencia

Todas las migraciones son idempotentes - pueden ejecutarse múltiples veces sin efectos secundarios. Si una colección ya existe, se omite su creación.

## Integración con Migrator

El servicio `migrator` en docker-compose ejecuta automáticamente:

1. `git pull` del repositorio edugo-infrastructure
2. Migraciones de PostgreSQL
3. Migraciones de MongoDB (usando `runner.go all`)

No es necesario ejecutar las migraciones manualmente en desarrollo con Docker.

## Desarrollo

### Agregar una Nueva Migración

1. **Crear archivo de estructura** en `structure/XXX_nombre_coleccion.go`:
   ```go
   func CreateNombreColeccion(ctx context.Context, db *mongo.Database) error {
       // Define el schema con validación
   }
   ```

2. **Crear archivo de constraints** en `constraints/XXX_nombre_coleccion_indexes.go`:
   ```go
   func CreateNombreColeccionIndexes(ctx context.Context, db *mongo.Database) error {
       // Define los índices
   }
   ```

3. **Actualizar runner.go**:
   - Agregar llamada en `runStructure()`
   - Agregar llamada en `runConstraints()`

4. **Compilar y probar**:
   ```bash
   go build -o mongodb-runner runner.go
   ./mongodb-runner all
   ```

## Referencias

- [MongoDB Schema Validation](https://www.mongodb.com/docs/manual/core/schema-validation/)
- [MongoDB Indexes](https://www.mongodb.com/docs/manual/indexes/)
- [MongoDB TTL Indexes](https://www.mongodb.com/docs/manual/core/index-ttl/)

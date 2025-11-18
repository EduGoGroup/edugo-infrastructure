# MongoDB Migration Architecture

## 📁 Nueva Estructura

```
mongodb/
├── structure/          # Schemas y colecciones (paso 1)
│   ├── 001_material_assessment.go
│   ├── 002_material_content.go
│   └── ...
├── constraints/        # Índices y validaciones (paso 2)
│   ├── 001_material_assessment_indexes.go
│   ├── 002_material_content_indexes.go
│   └── ...
├── seeds/             # Datos iniciales para producción (paso 3)
│   ├── 001_default_settings.go
│   └── ...
├── testing/           # Datos de prueba para development (paso 4)
│   ├── 001_test_materials.go
│   └── ...
└── migrations/        # DEPRECATED - Scripts JS antiguos
    └── ...
```

## 🎯 Filosofía de Separación

### 1. Structure (Estructura)
**Qué contiene:** Definición de colecciones y schemas de validación
**Cuándo se ejecuta:** Siempre (producción y desarrollo)
**Propósito:** Crear la estructura base de datos

```go
// structure/001_material_assessment.go
func CreateMaterialAssessment(ctx context.Context, db *mongo.Database) error {
    // Define validator schema
    // Create collection with validator
}
```

### 2. Constraints (Restricciones)
**Qué contiene:** Índices, validaciones adicionales
**Cuándo se ejecuta:** Siempre (producción y desarrollo)
**Propósito:** Optimizar queries y garantizar integridad

```go
// constraints/001_material_assessment_indexes.go
func CreateMaterialAssessmentIndexes(ctx context.Context, db *mongo.Database) error {
    // Create indexes for efficient queries
}
```

### 3. Seeds (Datos Iniciales)
**Qué contiene:** Datos mínimos necesarios para funcionamiento
**Cuándo se ejecuta:** Siempre (producción y desarrollo)
**Propósito:** Datos base del sistema

```go
// seeds/001_default_settings.go
func SeedDefaultSettings(ctx context.Context, db *mongo.Database) error {
    // Insert default system settings
}
```

### 4. Testing (Datos de Prueba)
**Qué contiene:** Datos de ejemplo para desarrollo/testing
**Cuándo se ejecuta:** Solo en development
**Propósito:** Facilitar desarrollo y pruebas

```go
// testing/001_test_materials.go
func SeedTestMaterials(ctx context.Context, db *mongo.Database) error {
    // Insert test data
}
```

## 🔄 Orden de Ejecución

1. **Structure** → Crear colecciones y schemas
2. **Constraints** → Aplicar índices y validaciones
3. **Seeds** → Insertar datos iniciales
4. **Testing** → Insertar datos de prueba (solo dev)

## 🎁 Ventajas

### vs Migraciones Tradicionales
- ❌ Antes: Todo mezclado en un archivo JS
- ✅ Ahora: Separación clara de responsabilidades

### vs Scripts JavaScript
- ❌ Antes: Requiere mongosh (incompatible con Alpine/ARM)
- ✅ Ahora: Código Go nativo (funciona en cualquier plataforma)

### vs Migraciones Secuenciales
- ❌ Antes: Un solo flujo lineal
- ✅ Ahora: Control granular de qué ejecutar

## 🚀 Uso en Migrator

```go
// En edugo-dev-environment/migrator
func runMongoMigrations(env string) error {
    // 1. Structure (siempre)
    structure.CreateMaterialAssessment(ctx, db)
    structure.CreateMaterialContent(ctx, db)
    
    // 2. Constraints (siempre)
    constraints.CreateMaterialAssessmentIndexes(ctx, db)
    constraints.CreateMaterialContentIndexes(ctx, db)
    
    // 3. Seeds (siempre)
    seeds.SeedDefaultSettings(ctx, db)
    
    // 4. Testing (solo development)
    if env == "development" {
        testing.SeedTestMaterials(ctx, db)
    }
    
    return nil
}
```

## 📋 Migraciones a Convertir

| # | Script Antiguo | Structure | Constraints | Seeds | Testing |
|---|----------------|-----------|-------------|-------|---------|
| 001 | material_assessment | ✅ | ✅ | - | - |
| 002 | material_content | ⏳ | ⏳ | - | - |
| 003 | assessment_attempt_result | ⏳ | ⏳ | - | - |
| 004 | audit_logs | ⏳ | ⏳ | - | - |
| 005 | notifications | ⏳ | ⏳ | - | - |
| 006 | analytics_events | ⏳ | ⏳ | - | - |
| 007 | material_summary | ⏳ | ⏳ | - | - |
| 008 | material_assessment_worker | ⏳ | ⏳ | - | - |
| 009 | material_event | ⏳ | ⏳ | - | - |

## 🎯 Estado del Proyecto

- ✅ Arquitectura diseñada
- ✅ Directorios creados
- ✅ Primer ejemplo implementado (001_material_assessment)
- ⏳ Pendiente: Convertir las 8 migraciones restantes
- ⏳ Pendiente: Actualizar migrator para usar nueva estructura

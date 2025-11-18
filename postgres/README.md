# PostgreSQL - Estructura de 4 Capas

## 📁 Estructura

```
postgres/
├── structure/          # Capa 1: Definición de tablas (sin FK, sin CHECK)
│   ├── 001_users.sql
│   ├── 002_schools.sql
│   ├── 003_academic_units.sql
│   ├── 004_memberships.sql
│   ├── 005_materials.sql
│   ├── 006_assessment.sql
│   ├── 007_assessment_attempt.sql
│   ├── 008_assessment_attempt_answer.sql
│   ├── 009_extend_assessment.sql
│   ├── 010_extend_attempt.sql
│   └── 011_extend_answer.sql
│
├── constraints/        # Capa 2: Foreign Keys, UNIQUE, CHECK, Triggers, Views
│   ├── 001_users.sql
│   ├── 002_schools.sql
│   ├── 003_academic_units.sql
│   ├── 004_memberships.sql
│   ├── 005_materials.sql
│   ├── 006_assessment.sql
│   ├── 007_assessment_attempt.sql
│   ├── 008_assessment_attempt_answer.sql
│   ├── 009_extend_assessment.sql
│   ├── 010_extend_attempt.sql
│   └── 011_extend_answer.sql
│
├── seeds/              # Capa 3: Datos iniciales/demo
│   └── (vacío por ahora)
│
├── testing/            # Capa 4: Tests SQL
│   └── (vacío por ahora)
│
├── runner.go           # Runner Go para ejecutar las 4 capas
├── go.mod              # Módulo Go
└── migrations/         # (legacy) Migraciones originales
```

## 🚀 Uso

### Ejecutar con runner.go

```bash
# Configurar variables de entorno (opcional)
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=edugo
export POSTGRES_PASSWORD=edugo_dev_2024
export POSTGRES_DB=edugo_db

# Ejecutar runner
go run runner.go

# O compilar y ejecutar
go build -o runner runner.go
./runner
```

### Salida esperada

```
✓ Conectado a PostgreSQL: edugo@localhost:5432/edugo_db

═══════════════════════════════════════════════════════════════
  CAPA: STRUCTURE
═══════════════════════════════════════════════════════════════

▸ Ejecutando: 001_users.sql
  ✓ Éxito
▸ Ejecutando: 002_schools.sql
  ✓ Éxito
...

═══════════════════════════════════════════════════════════════
  CAPA: CONSTRAINTS
═══════════════════════════════════════════════════════════════

▸ Ejecutando: 001_users.sql
  ✓ Éxito
...

═══════════════════════════════════════════════════════════════
  RESUMEN FINAL
═══════════════════════════════════════════════════════════════
✓ Archivos ejecutados: 22
⊘ Archivos omitidos: 0
✓ Todas las capas procesadas exitosamente
```

## 📋 Orden de ejecución

1. **STRUCTURE** (azul): Crea tablas, índices, comentarios
2. **CONSTRAINTS** (púrpura): Agrega FK, UNIQUE, CHECK, triggers, views
3. **SEEDS** (verde): Inserta datos iniciales
4. **TESTING** (cyan): Ejecuta tests SQL

## 🔧 Características

- ✅ Separa estructura de constraints para evitar dependencias circulares
- ✅ Preserva TODO: CHECK constraints, COMMENTS, UNIQUE, triggers, views
- ✅ Nombres cortos: `001_users.sql` en lugar de `001_create_users.sql`
- ✅ Mantiene orden 001-011 de migraciones originales
- ✅ Runner Go con colores y resumen detallado
- ✅ Idempotente: se puede ejecutar múltiples veces
- ✅ Archivos 009-011 vacíos (extensiones ya incluidas en 006-008)

## 🎯 Ventajas sobre migraciones

- **Atómico**: Se ejecuta todo o nada
- **Rápido**: No hay control de versiones, ideal para dev
- **Simple**: Un comando ejecuta todo
- **Flexible**: Fácil agregar seeds y tests
- **Visual**: Colores por capa para seguimiento claro

## 📝 Notas

- Los archivos 009, 010, 011 existen solo para compatibilidad
- Las extensiones ya están incluidas en 006, 007, 008
- Los directorios `seeds/` y `testing/` están preparados para uso futuro

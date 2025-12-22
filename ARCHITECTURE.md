# 🏗️ Arquitectura del Proyecto edugo-infrastructure

Documentación técnica completa de la arquitectura, decisiones de diseño y patrones utilizados en el proyecto `edugo-infrastructure`.

---

## 📋 Tabla de Contenidos

- [Visión General](#-visión-general)
- [Estructura de Módulos](#-estructura-de-módulos)
- [Arquitectura de Migraciones](#-arquitectura-de-migraciones)
- [Decisiones de Diseño](#-decisiones-de-diseño)
- [Patrones y Convenciones](#-patrones-y-convenciones)
- [Casos de Uso Comunes](#-casos-de-uso-comunes)
- [Integración con Otros Proyectos](#-integración-con-otros-proyectos)
- [Referencias](#-referencias)

---

## 🎯 Visión General

### Propósito del Proyecto

`edugo-infrastructure` es un monorepo que centraliza la infraestructura de datos para el ecosistema EduGo:

- **Migraciones de bases de datos** (PostgreSQL y MongoDB)
- **Schemas compartidos** entre microservicios
- **CLIs de administración** para operaciones de BD
- **Configuración de entornos** de desarrollo y testing

### Audiencia de Este Documento

| Rol | Uso de Este Documento |
|-----|----------------------|
| 👨‍💻 **Desarrolladores** | Entender estructura, agregar migraciones, consumir módulos |
| 🛠️ **DevOps** | Setup de ambientes, CI/CD, troubleshooting |
| 🏗️ **Arquitectos** | Decisiones de diseño, trade-offs, evolución del sistema |
| 🧪 **QA** | Estrategia de testing, datos de prueba |

### Principios de Diseño

```
┌─────────────────────────────────────────────────────┐
│  1. Modularidad        → Módulos Go independientes  │
│  2. Reproducibilidad   → Migraciones versionadas    │
│  3. Simplicidad        → Sin frameworks complejos   │
│  4. Observabilidad     → Logs estructurados         │
│  5. Testabilidad       → Testcontainers + fixtures  │
│  6. Versionado Claro   → Tags por módulo           │
└─────────────────────────────────────────────────────┘
```

### Contexto del Negocio

EduGo es una plataforma educativa que gestiona:

- 👥 **Usuarios**: Estudiantes, profesores, administradores
- 📚 **Cursos**: Contenidos educativos estructurados
- 💳 **Transacciones**: Pagos, suscripciones
- 📊 **Analytics**: Métricas de uso y progreso

Este proyecto asegura que todos los microservicios trabajen con la misma estructura de datos.

---

## 📁 Estructura de Módulos

### Árbol del Proyecto

```
edugo-infrastructure/
├── postgres/               # Módulo PostgreSQL
│   ├── cmd/
│   │   ├── migrate/       # CLI de migraciones (up, down, status)
│   │   └── runner/        # Runner de 4 capas (structure, migrations, seeds, mock)
│   ├── migrations/
│   │   ├── structure/     # Capa 1: Esquema base (DDL)
│   │   ├── migrations/    # Capa 2: Migraciones versionadas
│   │   ├── seeds/         # Capa 3: Datos iniciales (producción)
│   │   ├── mock/          # Capa 4: Datos de prueba (desarrollo)
│   │   └── embed.go       # Go embeds para archivos SQL
│   ├── go.mod
│   └── README.md
│
├── mongodb/                # Módulo MongoDB
│   ├── cmd/
│   │   └── migrate/       # CLI de migraciones (up, down, status, force)
│   ├── migrations/
│   │   ├── migrations.go  # Definición de migraciones
│   │   ├── seeds.go       # ApplySeeds() - Datos iniciales
│   │   ├── mock.go        # ApplyMockData() - Datos de prueba
│   │   └── embed.go       # Go embeds para archivos JSON
│   ├── go.mod
│   └── README.md
│
├── schemas/                # Módulo de Schemas
│   ├── user_schema.go
│   ├── course_schema.go
│   ├── transaction_schema.go
│   ├── validator.go
│   ├── go.mod
│   └── README.md
│
├── documents/              # Documentación
│   ├── README.md
│   ├── RELEASING.md       # Guía de versionado y releases
│   └── ARCHITECTURE.md    # Este archivo (simbólico link a raíz)
│
├── improvements/           # Mejoras y deuda técnica
│   ├── README.md
│   ├── TECHNICAL_DEBT.md
│   ├── DUPLICATED_CODE.md
│   ├── DEPRECATED_PATTERNS.md
│   └── MISSING_FEATURES.md
│
├── .github/
│   └── workflows/
│       └── ci.yml         # CI/CD con GitHub Actions
│
├── go.work                # Go workspace (desarrollo local)
└── README.md
```

### Desglose por Módulo

#### 🐘 `postgres/`

**Propósito**: Gestión completa de migraciones PostgreSQL

**Componentes clave**:

1. **CLI `migrate`** (`cmd/migrate/migrate.go`):
   ```bash
   postgres-migrate up        # Aplicar migraciones pendientes
   postgres-migrate down      # Revertir última migración
   postgres-migrate status    # Ver estado de migraciones
   postgres-migrate create    # Crear nueva migración
   ```

2. **Runner de 4 capas** (`cmd/runner/runner.go`):
   ```bash
   postgres-runner structure  # Ejecuta solo capa 1 (DDL)
   postgres-runner migrations # Ejecuta capas 1+2
   postgres-runner seeds      # Ejecuta capas 1+2+3
   postgres-runner mock       # Ejecuta capas 1+2+3+4 (full)
   ```

**Uso en otros proyectos**:
```go
import "github.com/edugo/edugo-infrastructure/postgres/migrations"

// En tests de integración
migrations.ApplyStructure(db)
migrations.ApplyMigrations(db)
migrations.ApplySeeds(db)
```

#### 🍃 `mongodb/`

**Propósito**: Gestión completa de migraciones MongoDB

**Componentes clave**:

1. **CLI `migrate`** (`cmd/migrate/migrate.go`):
   ```bash
   mongodb-migrate up         # Aplicar migraciones pendientes
   mongodb-migrate down       # Revertir última migración
   mongodb-migrate status     # Ver estado de migraciones
   mongodb-migrate force 5    # Forzar versión (cuidado!)
   ```

2. **Funciones Go públicas** (`migrations/migrations.go`):
   ```go
   ApplyAll()      // Aplica todas las migraciones versionadas
   ApplySeeds()    // Aplica datos iniciales (22 docs, 6 colecciones)
   ApplyMockData() // Aplica datos de prueba (35 docs, 6 colecciones)
   ```

**Uso en otros proyectos**:
```go
import "github.com/edugo/edugo-infrastructure/mongodb/migrations"

// En tests de integración
migrations.ApplyAll(ctx, db)
migrations.ApplySeeds(ctx, db)
migrations.ApplyMockData(ctx, db)
```

#### 📋 `schemas/`

**Propósito**: Schemas compartidos entre microservicios

**Componentes clave**:

```go
// user_schema.go
type User struct {
    ID        string    `json:"id" bson:"_id"`
    Email     string    `json:"email" bson:"email"`
    Name      string    `json:"name" bson:"name"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// validator.go
func ValidateUserEmail(email string) error
func ValidateCourseName(name string) error
```

**Uso en otros proyectos**:
```go
import "github.com/edugo/edugo-infrastructure/schemas"

user := schemas.User{
    Email: "alumno@example.com",
    Name:  "Juan Pérez",
}

if err := schemas.ValidateUserEmail(user.Email); err != nil {
    return err
}
```

---

## 🔄 Arquitectura de Migraciones

### Sistema de 4 Capas (PostgreSQL)

El runner de PostgreSQL utiliza una arquitectura en capas que permite setup incremental:

```
┌─────────────────────────────────────────────────────────┐
│  CAPA 1: STRUCTURE (Base Schema - DDL)                  │
│  ├─ 001_initial_schema.sql                              │
│  ├─ 002_core_tables.sql                                 │
│  └─ 003_indexes.sql                                     │
│                                                          │
│  Propósito: Crear estructura base de tablas             │
│  Cuándo ejecutar: Primera vez, ambientes nuevos         │
│  Idempotente: Sí (usa CREATE IF NOT EXISTS)             │
└─────────────────────────────────────────────────────────┘
          ↓
┌─────────────────────────────────────────────────────────┐
│  CAPA 2: MIGRATIONS (Cambios Versionados)               │
│  ├─ 001_add_user_preferences.sql                        │
│  ├─ 002_add_course_categories.sql                       │
│  └─ 003_add_audit_logs.sql                              │
│                                                          │
│  Propósito: Evolución del schema en el tiempo           │
│  Cuándo ejecutar: Después de structure, al actualizar   │
│  Idempotente: Depende (usar IF NOT EXISTS)              │
└─────────────────────────────────────────────────────────┘
          ↓
┌─────────────────────────────────────────────────────────┐
│  CAPA 3: SEEDS (Datos Iniciales - Producción)           │
│  ├─ 001_seed_roles.sql                                  │
│  ├─ 002_seed_permissions.sql                            │
│  └─ 003_seed_system_config.sql                          │
│                                                          │
│  Propósito: Datos requeridos para funcionamiento        │
│  Cuándo ejecutar: Setup inicial, ambientes productivos  │
│  Idempotente: Sí (usa INSERT ... ON CONFLICT DO NOTHING)│
└─────────────────────────────────────────────────────────┘
          ↓
┌─────────────────────────────────────────────────────────┐
│  CAPA 4: MOCK (Datos de Prueba - Desarrollo)            │
│  ├─ 001_mock_users.sql                                  │
│  ├─ 002_mock_courses.sql                                │
│  └─ 003_mock_transactions.sql                           │
│                                                          │
│  Propósito: Datos realistas para desarrollo y testing   │
│  Cuándo ejecutar: SOLO en dev/test, NUNCA en producción │
│  Idempotente: Sí (usa TRUNCATE + INSERT)                │
└─────────────────────────────────────────────────────────┘
```

### Diferencia entre ApplyAll(), ApplySeeds(), ApplyMockData() (MongoDB)

MongoDB no usa archivos SQL, sino funciones Go:

| Función | Descripción | Cuándo Usar | Colecciones | Documentos |
|---------|-------------|-------------|-------------|------------|
| `ApplyAll()` | Aplica migraciones versionadas | Siempre (prod, dev, test) | Todas | - |
| `ApplySeeds()` | Inserta datos iniciales | Producción y desarrollo | 6 | 22 |
| `ApplyMockData()` | Inserta datos de prueba | Solo desarrollo y tests | 6 | 35 |

**Ejemplo de flujo**:

```go
// Ambiente de producción
migrations.ApplyAll(ctx, db)     // ✅ Migraciones
migrations.ApplySeeds(ctx, db)   // ✅ Datos iniciales
// migrations.ApplyMockData()    // ❌ NUNCA en producción

// Ambiente de desarrollo/testing
migrations.ApplyAll(ctx, db)     // ✅ Migraciones
migrations.ApplySeeds(ctx, db)   // ✅ Datos iniciales
migrations.ApplyMockData(ctx, db)// ✅ Datos de prueba
```

### Flujo de Ejecución

#### PostgreSQL (Runner)

```
Usuario ejecuta: postgres-runner mock

    ↓
[1. Conectar a BD]
    ↓
[2. Ejecutar Structure (Capa 1)]
    ├─ Leer archivos .sql de migrations/structure/
    ├─ Ejecutar en orden alfabético
    └─ Log: "Aplicado structure: 001_initial_schema.sql"
    ↓
[3. Ejecutar Migrations (Capa 2)]
    ├─ Crear tabla schema_migrations si no existe
    ├─ Leer versión actual
    ├─ Aplicar migraciones pendientes
    └─ Log: "Migración 003 aplicada exitosamente"
    ↓
[4. Ejecutar Seeds (Capa 3)]
    ├─ Leer archivos .sql de migrations/seeds/
    ├─ Ejecutar con ON CONFLICT DO NOTHING
    └─ Log: "Aplicado seeds: 001_seed_roles.sql"
    ↓
[5. Ejecutar Mock (Capa 4)]
    ├─ Leer archivos .sql de migrations/mock/
    ├─ Ejecutar (TRUNCATE + INSERT)
    └─ Log: "Aplicado mock: 001_mock_users.sql"
    ↓
✅ Base de datos lista para desarrollo
```

#### MongoDB (Funciones Go)

```
Usuario ejecuta: migrations.ApplyMockData(ctx, db)

    ↓
[1. Verificar conexión]
    ↓
[2. ApplyAll() - Migraciones]
    ├─ Crear colección schema_migrations si no existe
    ├─ Leer versión actual
    ├─ Aplicar migraciones pendientes (funciones Go)
    └─ Log: "Aplicadas 5 migraciones"
    ↓
[3. ApplySeeds() - Datos iniciales]
    ├─ Insertar docs en 6 colecciones (22 docs)
    ├─ Usar insertMany con ordered: false
    └─ Log: "Seeds aplicados: 22 documentos"
    ↓
[4. ApplyMockData() - Datos de prueba]
    ├─ Eliminar datos previos de testing
    ├─ Insertar docs en 6 colecciones (35 docs)
    └─ Log: "Mock data aplicado: 35 documentos"
    ↓
✅ MongoDB lista para desarrollo/testing
```

### Tabla Comparativa de Funciones

| Aspecto | PostgreSQL Runner | MongoDB Functions |
|---------|------------------|------------------|
| **Ejecutable** | CLI binario | Importación Go |
| **Formato** | Archivos .sql | Código Go |
| **Capas** | 4 (structure, migrations, seeds, mock) | 3 (migrations, seeds, mock) |
| **Versionado** | Tabla `schema_migrations` | Colección `schema_migrations` |
| **Rollback** | `migrate down` | `migrate down` |
| **Testing** | Testcontainers + runner | Testcontainers + funciones |

---

## 💡 Decisiones de Diseño

### ¿Por qué Go Embeds en lugar de Archivos Externos?

**Problema antes de embeds**:
```
❌ Distribuir archivos .sql por separado
❌ Paths relativos rompen en diferentes ambientes
❌ Riesgo de archivos faltantes en producción
❌ Complejidad en CI/CD
```

**Solución con Go embeds**:
```go
//go:embed migrations/*.sql
var migrationsFS embed.FS

// ✅ Archivos compilados en el binario
// ✅ Paths siempre correctos
// ✅ Deployment simple (un solo binario)
// ✅ Funciona en cualquier ambiente
```

**Diagrama de comparación**:
```
SIN EMBEDS                        CON EMBEDS
┌──────────────┐                  ┌──────────────┐
│  app         │                  │  app         │
│  ├─ bin/     │                  │  (binario)   │
│  ├─ sql/     │ ← Debe copiar    │              │ ← Todo incluido
│  └─ config/  │                  │              │
└──────────────┘                  └──────────────┘
```

### ¿Por qué 4 Capas en PostgreSQL?

**Problema**: Setup de BD requiere diferentes niveles según el ambiente

**Antes (sin capas)**:
```sql
-- Un solo script gigante: all_in_one.sql
CREATE TABLE users (...);
INSERT INTO roles VALUES ('admin'), ('user');
INSERT INTO users VALUES ('test@example.com'); -- ¡Datos de prueba en producción!
```

**Después (con capas)**:
```bash
# Producción
postgres-runner seeds  # Solo structure + migrations + seeds

# Desarrollo
postgres-runner mock   # Incluye también datos de prueba
```

**Beneficios**:
- ✅ Control granular de qué datos cargar
- ✅ Previene datos de prueba en producción
- ✅ Setup más rápido (solo capas necesarias)
- ✅ Tests de integración más fáciles

### ¿Por qué NO Usar ORMs?

**Problemas con ORMs**:

1. **Performance**:
   ```go
   // ORM genera N+1 queries
   users := orm.Find("users")
   for user := range users {
       user.Courses() // ← 1 query por usuario!
   }
   
   // SQL manual: 1 query con JOIN
   SELECT u.*, c.* FROM users u
   LEFT JOIN courses c ON c.user_id = u.id
   ```

2. **Migraciones complejas**:
   ```go
   // ORM limita a cambios simples
   orm.AddColumn("users", "age", "int")
   
   // SQL permite lógica compleja
   ALTER TABLE users ADD COLUMN age INT;
   UPDATE users SET age = EXTRACT(YEAR FROM AGE(NOW(), birth_date));
   ```

3. **Control total**:
   - PostgreSQL tiene features avanzadas (CTEs, window functions, JSONB)
   - ORMs abstraen demasiado y limitan expresividad
   - SQL raw es más explícito y debuggeable

4. **Simplicidad**:
   - Sin DSL propietario que aprender
   - Sin "magia" oculta
   - Stack más pequeño (menos dependencias)

### ¿Por qué Versionado con Tags por Módulo?

**Problema**: Monorepo con múltiples módulos evolucionando independientemente

**Solución**: Tags con patrón `<módulo>/v<SemVer>`

```bash
postgres/v0.11.1    # PostgreSQL evoluciona a ritmo diferente
mongodb/v0.10.1     # MongoDB tiene su propia versión
schemas/v0.1.2      # Schemas raramente cambia
```

**Beneficios**:
- ✅ Consumidores pueden fijar versiones específicas por módulo
- ✅ go get funciona correctamente: `go get .../postgres@v0.11.1`
- ✅ Changelog separado por módulo
- ✅ Rollback granular

**Ver más en**: [documents/RELEASING.md](./documents/RELEASING.md)

### Trade-offs Documentados

| Decisión | ✅ Ventajas | ❌ Desventajas |
|----------|------------|---------------|
| **Go Embeds** | Binarios autosuficientes, deployment simple | Binario más grande |
| **4 Capas** | Control granular | Más archivos que mantener |
| **Sin ORM** | Performance, control total | Más código SQL manual |
| **Tags por módulo** | Versionado independiente | Más tags que gestionar |
| **Monorepo** | Código centralizado | Requiere Go workspace |

---

## 📐 Patrones y Convenciones

### Naming Conventions

#### SQL (PostgreSQL)

```sql
-- Tablas: plural, snake_case
users
course_enrollments
payment_transactions

-- Columnas: singular, snake_case
user_id
created_at
is_active

-- Índices: {table}_{columns}_idx
users_email_idx
courses_category_id_idx

-- Foreign keys: fk_{table}_{referenced_table}
fk_enrollments_users
fk_transactions_courses

-- Migraciones: {version}_{description}.sql
001_initial_schema.sql
002_add_user_preferences.sql
```

#### Go (MongoDB)

```go
// Structs: PascalCase
type User struct {}
type CourseEnrollment struct {}

// Campos: PascalCase (exportados)
type User struct {
    ID        string
    Email     string
    CreatedAt time.Time
}

// Tags JSON/BSON: snake_case
type User struct {
    ID string `json:"id" bson:"_id"`
}

// Funciones públicas: PascalCase
func ApplyMigrations() {}
func ValidateEmail() {}

// Migraciones: Migration{version}{Description}
func Migration001InitialCollections() {}
func Migration002AddUserPreferences() {}
```

### Estructura de Archivos

#### PostgreSQL Migration

```sql
-- migrations/migrations/003_add_audit_logs.sql

-- Up migration
CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id INTEGER,
    changes JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- Índices
    CONSTRAINT chk_action CHECK (action IN ('CREATE', 'UPDATE', 'DELETE'))
);

CREATE INDEX IF NOT EXISTS audit_logs_user_id_idx ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS audit_logs_created_at_idx ON audit_logs(created_at DESC);

-- Down migration (comentado, se usa con CLI)
-- DROP TABLE IF EXISTS audit_logs CASCADE;
```

#### MongoDB Migration

```go
// migrations/migration_003_add_user_preferences.go

func Migration003AddUserPreferences() Migration {
    return Migration{
        Version: 3,
        Name:    "add_user_preferences",
        Up: func(ctx context.Context, db *mongo.Database) error {
            // Crear colección con validación
            validator := bson.M{
                "$jsonSchema": bson.M{
                    "bsonType": "object",
                    "required": []string{"user_id", "preferences"},
                    "properties": bson.M{
                        "user_id": bson.M{"bsonType": "string"},
                        "preferences": bson.M{
                            "bsonType": "object",
                            "properties": bson.M{
                                "email_notifications": bson.M{"bsonType": "bool"},
                                "theme": bson.M{"enum": []string{"light", "dark"}},
                            },
                        },
                    },
                },
            }
            
            opts := options.CreateCollection().SetValidator(validator)
            return db.CreateCollection(ctx, "user_preferences", opts)
        },
        Down: func(ctx context.Context, db *mongo.Database) error {
            return db.Collection("user_preferences").Drop(ctx)
        },
    }
}
```

### Testing Strategy

#### Tests de Integración (recomendado)

```go
// postgres/migrations_test.go

func TestMigrations(t *testing.T) {
    ctx := context.Background()
    
    // Testcontainer PostgreSQL
    container, err := postgres.RunContainer(ctx,
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    require.NoError(t, err)
    defer container.Terminate(ctx)
    
    connStr, err := container.ConnectionString(ctx)
    require.NoError(t, err)
    
    db, err := sql.Open("postgres", connStr)
    require.NoError(t, err)
    defer db.Close()
    
    // Aplicar migraciones
    err = ApplyAll(db)
    require.NoError(t, err)
    
    // Verificar estructura
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
    require.NoError(t, err)
}
```

#### Tests Unitarios (validadores)

```go
// schemas/validator_test.go

func TestValidateUserEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateUserEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### CI/CD Workflow

```yaml
# .github/workflows/ci.yml

name: CI

on: [push, pull_request]

jobs:
  test-postgres:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      
      - name: Run Postgres Tests
        run: |
          cd postgres
          go test -v ./...
  
  test-mongodb:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      
      - name: Run MongoDB Tests
        run: |
          cd mongodb
          go test -v ./...
```

---

## 🎬 Casos de Uso Comunes

### 1. Setup de Desarrollo desde Cero

```bash
# 1. Clonar repositorio
git clone https://github.com/edugo/edugo-infrastructure.git
cd edugo-infrastructure

# 2. Levantar bases de datos con Docker
docker run -d --name postgres-dev \
  -e POSTGRES_PASSWORD=dev \
  -e POSTGRES_DB=edugo \
  -p 5432:5432 \
  postgres:15

docker run -d --name mongo-dev \
  -p 27017:27017 \
  mongo:7.0

# 3. Compilar CLIs
cd postgres/cmd/runner && go build -o postgres-runner
cd ../../../mongodb/cmd/migrate && go build -o mongodb-migrate

# 4. Inicializar PostgreSQL con datos de prueba
export POSTGRES_DSN="postgres://postgres:dev@localhost:5432/edugo?sslmode=disable"
./postgres-runner mock  # Aplica las 4 capas

# 5. Inicializar MongoDB con datos de prueba
export MONGODB_URI="mongodb://localhost:27017/edugo"
./mongodb-migrate up
# En código Go:
# migrations.ApplySeeds(ctx, db)
# migrations.ApplyMockData(ctx, db)

# 6. Verificar
psql $POSTGRES_DSN -c "SELECT COUNT(*) FROM users;"
mongosh mongodb://localhost:27017/edugo --eval "db.users.countDocuments()"
```

### 2. Inicializar BD en Producción

```bash
# PostgreSQL (solo structure + migrations + seeds, SIN mock)
export POSTGRES_DSN="postgres://prod_user:prod_pass@prod-host:5432/edugo"
./postgres-runner seeds

# MongoDB (solo migrations + seeds, SIN mock data)
export MONGODB_URI="mongodb://prod-host:27017/edugo?authSource=admin"
./mongodb-migrate up
# En código Go:
# migrations.ApplyAll(ctx, db)
# migrations.ApplySeeds(ctx, db)
```

### 3. Agregar Nueva Migración PostgreSQL

```bash
# 1. Crear archivo de migración
cd postgres/migrations/migrations/
touch 012_add_course_reviews.sql

# 2. Escribir SQL
cat > 012_add_course_reviews.sql <<'EOF'
CREATE TABLE IF NOT EXISTS course_reviews (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(course_id, user_id)
);

CREATE INDEX course_reviews_course_id_idx ON course_reviews(course_id);
EOF

# 3. Probar localmente
./postgres-migrate up

# 4. Verificar
psql $POSTGRES_DSN -c "\d course_reviews"

# 5. Commit y PR
git add postgres/migrations/migrations/012_add_course_reviews.sql
git commit -m "feat(postgres): add course_reviews table"
```

### 4. Agregar Nueva Migración MongoDB

```bash
# 1. Editar migrations/migrations.go
cd mongodb/migrations/

# Agregar al final de AllMigrations():
func AllMigrations() []Migration {
    return []Migration{
        // ... migraciones existentes
        Migration012AddCourseReviews(),
    }
}

func Migration012AddCourseReviews() Migration {
    return Migration{
        Version: 12,
        Name:    "add_course_reviews",
        Up: func(ctx context.Context, db *mongo.Database) error {
            validator := bson.M{
                "$jsonSchema": bson.M{
                    "bsonType": "object",
                    "required": []string{"course_id", "user_id", "rating"},
                    "properties": bson.M{
                        "course_id": bson.M{"bsonType": "string"},
                        "user_id":   bson.M{"bsonType": "string"},
                        "rating":    bson.M{"bsonType": "int", "minimum": 1, "maximum": 5},
                        "comment":   bson.M{"bsonType": "string"},
                    },
                },
            }
            
            opts := options.CreateCollection().SetValidator(validator)
            if err := db.CreateCollection(ctx, "course_reviews", opts); err != nil {
                return err
            }
            
            // Crear índice único
            _, err := db.Collection("course_reviews").Indexes().CreateOne(ctx,
                mongo.IndexModel{
                    Keys:    bson.D{{Key: "course_id", Value: 1}, {Key: "user_id", Value: 1}},
                    Options: options.Index().SetUnique(true),
                },
            )
            return err
        },
        Down: func(ctx context.Context, db *mongo.Database) error {
            return db.Collection("course_reviews").Drop(ctx)
        },
    }
}

# 2. Probar localmente
./mongodb-migrate up

# 3. Verificar
mongosh --eval "db.course_reviews.getIndexes()"

# 4. Commit
git commit -am "feat(mongodb): add course_reviews collection"
```

### 5. Ejecutar Tests de Integración

```bash
# Asegurar que Docker esté corriendo

# PostgreSQL
cd postgres
ENABLE_INTEGRATION_TESTS=true go test ./... -v

# MongoDB
cd ../mongodb
ENABLE_INTEGRATION_TESTS=true go test ./... -v

# Schemas
cd ../schemas
go test ./... -v
```

### 6. Consumir Módulos en Otro Proyecto

```go
// go.mod
module github.com/edugo/edugo-api-mobile

require (
    github.com/edugo/edugo-infrastructure/postgres v0.11.1
    github.com/edugo/edugo-infrastructure/mongodb v0.10.1
    github.com/edugo/edugo-infrastructure/schemas v0.1.2
)

// main_test.go
import (
    "github.com/edugo/edugo-infrastructure/postgres/migrations"
    "github.com/edugo/edugo-infrastructure/schemas"
)

func TestUserCreation(t *testing.T) {
    // Setup BD con testcontainer
    container, _ := postgres.RunContainer(ctx)
    db, _ := sql.Open("postgres", connStr)
    
    // Aplicar migraciones
    migrations.ApplyAll(db)
    migrations.ApplySeeds(db)
    
    // Crear usuario usando schema compartido
    user := schemas.User{
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    // ... test logic
}
```

---

## 🔗 Integración con Otros Proyectos

### Proyectos Consumidores

```
edugo-infrastructure (este repo)
        ↓
        ├─→ edugo-api-mobile (usa postgres + schemas)
        ├─→ edugo-api-administracion (usa postgres + mongodb + schemas)
        ├─→ edugo-worker (usa mongodb + schemas)
        └─→ edugo-shared (usa schemas)
```

### Diagrama de Dependencias

```
┌─────────────────────────────────────────────────────┐
│           edugo-infrastructure                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐         │
│  │ postgres │  │ mongodb  │  │ schemas  │         │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘         │
└───────┼─────────────┼─────────────┼────────────────┘
        │             │             │
        ↓             ↓             ↓
┌─────────────────────────────────────────────────────┐
│         Microservicios Consumidores                 │
│                                                     │
│  ┌────────────────┐  ┌────────────────┐           │
│  │ api-mobile     │  │ api-admin      │           │
│  │ - postgres ✓   │  │ - postgres ✓   │           │
│  │ - schemas ✓    │  │ - mongodb ✓    │           │
│  └────────────────┘  │ - schemas ✓    │           │
│                      └────────────────┘           │
│                                                     │
│  ┌────────────────┐  ┌────────────────┐           │
│  │ worker         │  │ shared         │           │
│  │ - mongodb ✓    │  │ - schemas ✓    │           │
│  │ - schemas ✓    │  └────────────────┘           │
│  └────────────────┘                               │
└─────────────────────────────────────────────────────┘
```

### Flujo de Eventos (Ejemplo con RabbitMQ)

```
[edugo-api-mobile]
      ↓ Publish evento
  (RabbitMQ)
      ↓ Consume evento
[edugo-worker]
      ↓ Usa schemas compartidos
  (Valida con schemas.ValidateUser)
      ↓ Persiste en MongoDB
  (Usa mongodb/migrations)
```

### Código de Ejemplo: Publisher

```go
// edugo-api-mobile/internal/events/publisher.go

import "github.com/edugo/edugo-infrastructure/schemas"

func PublishUserCreated(user schemas.User) error {
    event := UserCreatedEvent{
        UserID:    user.ID,
        Email:     user.Email,
        Timestamp: time.Now(),
    }
    
    payload, _ := json.Marshal(event)
    return rabbitMQ.Publish("user.created", payload)
}
```

### Código de Ejemplo: Consumer

```go
// edugo-worker/internal/handlers/user_handler.go

import (
    "github.com/edugo/edugo-infrastructure/schemas"
    "github.com/edugo/edugo-infrastructure/mongodb/migrations"
)

func HandleUserCreated(msg []byte) error {
    var event UserCreatedEvent
    json.Unmarshal(msg, &event)
    
    // Validar usando schema compartido
    if err := schemas.ValidateUserEmail(event.Email); err != nil {
        return err
    }
    
    // Persistir en MongoDB
    user := schemas.User{
        ID:    event.UserID,
        Email: event.Email,
    }
    
    _, err := mongoDB.Collection("users").InsertOne(ctx, user)
    return err
}
```

---

## 📚 Referencias

### Documentación del Proyecto

| Documento | Descripción | Ubicación |
|-----------|-------------|-----------|
| README.md principal | Visión general del proyecto | `/README.md` |
| ARCHITECTURE.md | Este documento | `/ARCHITECTURE.md` |
| RELEASING.md | Guía de versionado y releases | `/documents/RELEASING.md` |
| TECHNICAL_DEBT.md | Deuda técnica identificada | `/improvements/TECHNICAL_DEBT.md` |
| PostgreSQL README | Docs específicas de PostgreSQL | `/postgres/README.md` |
| MongoDB README | Docs específicas de MongoDB | `/mongodb/README.md` |
| Schemas README | Docs de schemas compartidos | `/schemas/README.md` |

### Enlaces Externos

- [Go Modules](https://go.dev/ref/mod)
- [Go Embed Directive](https://pkg.go.dev/embed)
- [Testcontainers Go](https://golang.testcontainers.org/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [MongoDB Documentation](https://www.mongodb.com/docs/)
- [Semantic Versioning](https://semver.org/)

### Comandos Útiles

```bash
# PostgreSQL
psql $POSTGRES_DSN -c "\dt"              # Listar tablas
psql $POSTGRES_DSN -c "\d users"         # Describir tabla
psql $POSTGRES_DSN -f script.sql         # Ejecutar script

# MongoDB
mongosh $MONGODB_URI                     # Conectar
db.getCollectionNames()                  # Listar colecciones
db.users.find().pretty()                 # Query con formato

# Go
go work sync                             # Sincronizar workspace
go list -m all                           # Listar dependencias
go mod tidy                              # Limpiar go.mod
go test ./... -v                         # Ejecutar todos los tests
```

### Estado del Proyecto

**Última actualización:** Diciembre 2025

**Métricas**:
- 🐘 Migraciones PostgreSQL: 11
- 🍃 Migraciones MongoDB: 11
- 📋 Schemas definidos: 3
- ✅ Tests de integración: Funcionando
- 📦 Módulos Go: 3 (postgres, mongodb, schemas)
- 🏷️ Tags de versión: postgres/v0.11.1, mongodb/v0.10.1, schemas/v0.1.2

**Estado de Mejoras**: Ver [improvements/README.md](./improvements/README.md)

---

## 🤝 Contribución

### Cómo Contribuir

1. **Crear branch**: `git checkout -b feature/nueva-migracion`
2. **Hacer cambios**: Agregar migración, tests, docs
3. **Ejecutar tests**: `go test ./...`
4. **Commit**: `git commit -m "feat(postgres): add new migration"`
5. **Push y PR**: `git push origin feature/nueva-migracion`

### Guías de Estilo

- SQL: snake_case, siempre usar `IF NOT EXISTS/IF EXISTS`
- Go: gofmt, golangci-lint
- Commits: [Conventional Commits](https://www.conventionalcommits.org/)

---

**Mantenido por:** Equipo de Infraestructura EduGo  
**Preguntas:** Abrir issue en el repositorio

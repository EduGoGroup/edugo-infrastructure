# Plan de Migración: Sistema de Roles Simple a RBAC (Role-Based Access Control)

**Fecha**: 2026-02-13
**Autor**: Análisis técnico
**Estado**: Ambiente preparado - Listo para ejecución
**Versión del Plan**: 3.0 (Ambiente sincronizado, feature branches creados, go.work verificado)

---

## PREPARACION DEL AMBIENTE (Completado)

Antes de iniciar la implementacion RBAC, se realizo una sincronizacion completa del ecosistema:

### Trabajo completado

| Accion | Detalle | Estado |
|--------|---------|--------|
| **Infrastructure cleanup** | Eliminados constraint stubs huerfanos 011-015 (archivos vacios sin structure correspondiente) | Completado |
| **Release postgres/v0.14.0** | Nuevo release con cleanup de stubs | Completado |
| **Release auth/v0.11.0** | Release de auth con cambios pendientes en main | Completado |
| **API-admin deps** | Actualizado postgres v0.13.0→v0.14.0, auth v0.9.0→v0.11.0. Eliminado campo `EmailVerified` de entities, repos, services y mocks | Completado |
| **API-mobile deps** | Actualizado postgres v0.14.0, common v0.9.0, middleware/gin v0.9.0, lifecycle v0.9.0, logger v0.9.0, testing v0.9.0. Eliminado `EmailVerified` de mocks. Swagger regenerado | Completado |
| **Dev-environment** | Actualizado postgres v0.13.0→v0.14.0 en migrator | Completado |
| **go.work** | Creado en `/Users/jhoanmedina/source/EduGo/repos-separados/go.work` para desarrollo local | Completado |
| **Feature branches** | Creados desde dev sincronizado en los 4 repos | Completado |

### Estado actual de los repos

| Repo | Branch actual | Tag base | main=dev |
|------|--------------|----------|----------|
| edugo-infrastructure | `feature/rbac-postgres-tables` | postgres/v0.14.0 | SI |
| edugo-shared | `feature/rbac-auth-middleware` | auth/v0.11.0, common/v0.9.0, middleware/gin/v0.9.0 | SI |
| edugo-api-administracion | `feature/rbac-api-admin` | v0.12.0 | SI |
| edugo-api-mobile | `feature/rbac-api-mobile` | v0.18.0 | SI |

### go.work (Desarrollo local)

```go
// Archivo: /Users/jhoanmedina/source/EduGo/repos-separados/go.work
go 1.25

use (
    ./edugo-infrastructure/postgres
    ./edugo-shared/auth
    ./edugo-shared/common
    ./edugo-shared/middleware/gin
    ./edugo-api-administracion
    ./edugo-api-mobile
)
```

**Verificado**: Ambas APIs compilan correctamente con y sin go.work (`GOWORK=off go build ./...`).

### Numeracion de migraciones disponible

Tras la limpieza de stubs:
- **Structure**: 000-010 ocupados, siguiente disponible: **011**
- **Constraints**: 001-010 ocupados (011-015 eliminados), siguiente disponible: **011**
- **Seeds**: solo 001 existe, siguiente: **002**
- **Testing**: 001-005 ocupados, siguiente: **006**

### Campo EmailVerified

Eliminado de todo el ecosistema. El campo `email_verified` **no existe** en el schema de BD actual ni en ninguna entity/repo/service. No es necesario considerarlo en RBAC.

---

## IMPORTANTE: Proceso de Trabajo Git

**ESTE PROCESO SE APLICA A TODOS LOS PROYECTOS QUE SE DEBEN MODIFICAR:**

### Antes de iniciar cualquier modificación:
0. **Revisar si hay modificaciones sin comitear**
Si tienes cambios locales, detener indicar al usuario un resumen pequeños de los cambios y preguntar que hacer con ellos (stash, commit, descartar)
1. **Verificar estado de main y dev**
   ```bash
   git fetch origin
   git checkout main
   git pull origin main
   git checkout dev
   git pull origin dev
   
   # Verificar que main y dev tengan EL MISMO CONTENIDO
   # (no necesariamente la misma cantidad de commits, sino el mismo código)
   git diff main dev
   ```
   
2. **Crear rama de trabajo desde dev**
   ```bash
   git checkout dev
   git checkout -b feature/rbac-migration-[modulo]
   ```

3. **Hacer modificaciones en la rama de feature**

4. **Al finalizar, crear PR de feature → dev**
   ```bash
   git push origin feature/rbac-migration-[modulo]
   # Crear PR en GitHub de: feature/rbac-migration-[modulo] → dev
   ```

5. **NUNCA hacer merge directo a main**

6. **Esperar aprobación del PR antes de continuar con siguientes pasos**

---

## Índice

1. [Contexto y Alcance](#1-contexto-y-alcance)
2. [Análisis del Sistema Actual](#2-análisis-del-sistema-actual)
3. [Arquitectura Propuesta](#3-arquitectura-propuesta)
4. [Módulo: edugo-infrastructure (postgres)](#4-módulo-edugo-infrastructure-postgres)
5. [Módulo: edugo-shared (auth, common, middleware)](#5-módulo-edugo-shared-auth-common-middleware)
6. [Proyecto: edugo-api-administracion](#6-proyecto-edugo-api-administracion)
7. [Proyecto: edugo-api-mobile](#7-proyecto-edugo-api-mobile)
8. [Plan de Implementación por Fases](#8-plan-de-implementación-por-fases)
9. [Flujo de Releases](#9-flujo-de-releases)
10. [Casos de Uso y Ejemplos](#10-casos-de-uso-y-ejemplos)
11. [Testing](#11-testing)
12. [Riesgos y Mitigación](#12-riesgos-y-mitigación)

---

## 1. Contexto y Alcance

### 1.1 Objetivo

Migrar de un sistema de roles simple basado en un string (`users.role`) a un sistema **RBAC completo** que permita:

- **Múltiples roles por usuario** según el contexto (escuela/unidad académica)
- **Permisos granulares** asociados a cada rol
- **Gestión flexible** de permisos a nivel de sistema, escuela y unidad

### 1.2 Caso de Uso Motivador

```
Usuario: Juan Pérez
- Administrador en Escuela A (school_id: 1)
- Profesor en Escuela B, Unidad "Matemáticas 3° Básico" (school_id: 2, unit_id: 5)
- Alumno en Escuela B, Unidad "Física Avanzada" (school_id: 2, unit_id: 8)
```

**Problema actual**: No es posible modelar esto con `users.role` que solo soporta UN rol global.

### 1.3 Requisitos Funcionales

1. Un usuario puede tener **múltiples roles** en diferentes contextos
2. Cada rol tiene un **conjunto de permisos predefinidos**
3. Los permisos son **verificables en tiempo de ejecución** (middleware/lógica)
4. Debe soportar **cambio de contexto** entre escuelas/unidades
5. Debe mantener **compatibilidad con JWT existente** (actualizado con contexto)
6. Debe permitir **extensión futura** de permisos sin cambios de schema

### 1.4 Estrategia de Migración

**IMPORTANTE**: Como estamos en etapa de desarrollo sin datos de producción, podemos hacer **borrón y cuenta nueva** en la base de datos.

**Esto significa**:
- ❌ **NO usaremos ALTER TABLE** para remover el campo `users.role`
- ✅ **SÍ modificaremos directamente** el script `001_create_users.up.sql` eliminando el campo `role`
- ✅ **Podemos eliminar completamente** la base de datos y recrearla sin problemas
- ✅ **Simplifica** el proceso de migración (menos scripts)

---

## 2. Análisis del Sistema Actual

### 2.1 Estado Actual - edugo-infrastructure/postgres

#### Script actual de users

**Archivo**: `postgres/migrations/structure/001_create_users.sql` (línea 8-9)

```sql
role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'teacher', 'student', 'guardian')),
```

**Constraint actual**: `postgres/migrations/constraints/001_create_users.sql`

```sql
-- Índice en campo role
CREATE INDEX idx_users_role ON users(role);
```

**Problema**: Campo `role` simple que solo permite 1 rol por usuario.

#### Estructura actual de memberships

**Archivo**: `postgres/migrations/structure/004_create_memberships.sql`

```sql
role VARCHAR(50) NOT NULL CHECK (role IN ('teacher', 'student', 'guardian', 'coordinator', 'admin', 'assistant')),
```

Ya es contextual (escuela/unidad) pero **sin tabla de permisos**.

### 2.2 Estado Actual - edugo-shared/auth

**Módulo**: `edugo-shared/auth` (versión actual: `auth/v0.11.0`)

**Archivo**: `auth/jwt.go` (líneas 20-28)

```go
type Claims struct {
    UserID    string          `json:"user_id"`
    Email     string          `json:"email"`
    Role      enum.SystemRole `json:"role"`
    SchoolID  string          `json:"school_id,omitempty"`
    ExpiresAt time.Time
    IssuedAt  time.Time
    Issuer    string
    TokenID   string
}
```

**Problema**: Solo 1 rol (`Role`) en el token.

### 2.3 Estado Actual - edugo-api-administracion

#### Sistema de autenticación

**Archivo**: `internal/auth/service/auth_service.go`

**Método Login()** genera tokens con:
- JWT con Claims actuales (1 rol)
- Refresh token almacenado en tabla `refresh_tokens`

#### Swagger

**Ubicación**: `/docs/swagger.json` (3,227 líneas)

**Generación**:
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion
make swagger  # Ejecuta: swag init -g cmd/main.go -o docs --parseInternal --parseDependency
```

**Importante**: Después de modificar endpoints o DTOs, **SIEMPRE** regenerar swagger y commitear el archivo generado.

### 2.4 Estado Actual - edugo-api-mobile

#### Sistema de autenticación

**Archivo**: `internal/client/auth_client.go`

- Valida tokens localmente con `JWTManager` de shared
- Fallback a api-admin `/v1/auth/verify` si validación local falla
- Circuit breaker para proteger llamadas remotas

#### Middleware de autorización

**Archivo**: `internal/infrastructure/http/middleware/remote_auth.go`

```go
func RequireTeacher() gin.HandlerFunc { ... }  // Basado en roles
func RequireAdmin() gin.HandlerFunc { ... }    // Basado en roles
```

**Problema**: Valida roles, no permisos.

#### Swagger

**Ubicación**: `/docs/swagger.json` (1,540 líneas)

**Generación**:
```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
make swagger  # Ejecuta: swag init -g cmd/main.go -o docs --parseInternal
```

**Característica especial**: Detección dinámica de host en tiempo de ejecución (ver `internal/infrastructure/http/router/swagger.go`).

---

## 3. Arquitectura Propuesta

### 3.1 Modelo de Datos RBAC

#### Nuevas Tablas

```sql
-- Catálogo de Roles
roles
├── id (UUID PK)
├── name (VARCHAR UNIQUE)
├── display_name (VARCHAR)
├── description (TEXT)
├── scope (ENUM: 'system', 'school', 'unit')
├── is_active (BOOLEAN)
├── created_at, updated_at

-- Catálogo de Permisos
permissions
├── id (UUID PK)
├── name (VARCHAR UNIQUE)  -- Formato: 'resource:action'
├── display_name (VARCHAR)
├── description (TEXT)
├── resource (VARCHAR)
├── action (VARCHAR)
├── scope (ENUM: 'system', 'school', 'unit')
├── is_active (BOOLEAN)
├── created_at, updated_at

-- Relación Roles → Permisos
role_permissions
├── id (UUID PK)
├── role_id (UUID FK → roles)
├── permission_id (UUID FK → permissions)
├── created_at
└── UNIQUE(role_id, permission_id)

-- Roles Asignados a Usuarios
user_roles
├── id (UUID PK)
├── user_id (UUID FK → users)
├── role_id (UUID FK → roles)
├── school_id (UUID FK → schools, nullable)
├── academic_unit_id (UUID FK → academic_units, nullable)
├── is_active (BOOLEAN)
├── granted_by (UUID FK → users)
├── granted_at (TIMESTAMP)
├── expires_at (TIMESTAMP, nullable)
├── created_at, updated_at
└── UNIQUE(user_id, role_id, school_id, academic_unit_id)
```

### 3.2 Nomenclatura de Permisos

**Patrón**: `{resource}:{action}` o `{resource}:{action}:{scope}`

Ejemplos:
- `users:create` - Crear usuarios
- `users:read:own` - Ver propio perfil
- `materials:create` - Crear materiales
- `materials:publish` - Publicar materiales
- `assessments:grade` - Calificar evaluaciones

---

## 4. Módulo: edugo-infrastructure (postgres)

### 4.1 Información del Módulo

**Ubicación**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure/postgres`

**Tipo**: Módulo Go independiente con versionamiento por tag

**Versión actual**: `postgres/v0.14.0`

**go.mod**:
```go
module github.com/EduGoGroup/edugo-infrastructure/postgres
go 1.25
```

**Estructura de migración**: Arquitectura de 4 capas
- `structure/` - Tablas base (DDL)
- `constraints/` - FK, UNIQUE, CHECK, Triggers
- `seeds/` - Datos iniciales (producción)
- `testing/` - Datos de prueba

### 4.2 Proceso de Trabajo Git para este Módulo

```bash
# 1. Verificar estado
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure
git fetch origin
git checkout main && git pull origin main
git checkout dev && git pull origin dev
git diff main dev  # Verificar que contenido sea igual

# 2. Feature branch ya creado desde dev sincronizado
git checkout feature/rbac-postgres-tables  # YA EXISTE

# 3. Hacer modificaciones (ver sección 4.3)

# 4. Commit y push
git add postgres/
git commit -m "feat(postgres): agregar tablas RBAC (roles, permissions, user_roles)"
git push origin feature/rbac-postgres-tables

# 5. Crear PR en GitHub: feature/rbac-postgres-tables → dev
# 6. Esperar aprobación y merge

# 7. Después del merge a dev, crear tag desde dev
git checkout dev
git pull origin dev
git tag postgres/v0.15.0 -m "feat: Sistema RBAC con roles y permisos"
git push origin postgres/v0.15.0
```

### 4.3 Archivos a Crear/Modificar

#### 4.3.1 MODIFICAR (no ALTER): Tabla users

**Archivo**: `postgres/migrations/structure/001_create_users.sql`

**Acción**: ELIMINAR la línea del campo `role`

**Línea a ELIMINAR** (línea ~8-9):
```sql
role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'teacher', 'student', 'guardian')),
```

**También ELIMINAR** el comentario asociado (si existe).

**Archivo**: `postgres/migrations/constraints/001_create_users.sql`

**Acción**: ELIMINAR el índice de role

**Línea a ELIMINAR**:
```sql
CREATE INDEX idx_users_role ON users(role);
```

**Rationale**: Como podemos hacer borrón y cuenta nueva, simplemente eliminamos el campo del script original. No necesitamos migración 022_remove_role_from_users.

---

#### 4.3.2 CREAR: Tabla roles

**Archivo NUEVO**: `postgres/migrations/structure/012_create_roles.sql`

```sql
-- ====================================================================
-- TABLA: roles
-- DESCRIPCIÓN: Catálogo maestro de roles del sistema RBAC
-- VERSIÓN: postgres/v0.15.0
-- FECHA: 2026-02-13
-- ====================================================================

-- Tipo ENUM para scope de roles
CREATE TYPE role_scope AS ENUM ('system', 'school', 'unit');

-- Tabla de roles
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    display_name VARCHAR(100) NOT NULL,
    description TEXT,
    scope role_scope NOT NULL DEFAULT 'school',
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Comentarios
COMMENT ON TABLE roles IS 'Catálogo maestro de roles del sistema RBAC';
COMMENT ON COLUMN roles.name IS 'Nombre único del rol (snake_case)';
COMMENT ON COLUMN roles.display_name IS 'Nombre para mostrar en UI';
COMMENT ON COLUMN roles.scope IS 'Alcance del rol: system (global), school (institución), unit (clase/sección)';
```

**Archivo NUEVO**: `postgres/migrations/constraints/012_create_roles.sql`

```sql
-- ====================================================================
-- CONSTRAINTS: roles
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Índices
CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_roles_scope ON roles(scope);
CREATE INDEX idx_roles_active ON roles(is_active);

-- Trigger para updated_at (usa función existente update_updated_at_column)
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

---

#### 4.3.3 CREAR: Tabla permissions

**Archivo NUEVO**: `postgres/migrations/structure/013_create_permissions.sql`

```sql
-- ====================================================================
-- TABLA: permissions
-- DESCRIPCIÓN: Catálogo maestro de permisos del sistema RBAC
-- VERSIÓN: postgres/v0.15.0
-- FECHA: 2026-02-13
-- ====================================================================

-- Tipo ENUM para scope de permisos
CREATE TYPE permission_scope AS ENUM ('system', 'school', 'unit');

-- Tabla de permisos
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    display_name VARCHAR(150) NOT NULL,
    description TEXT,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    scope permission_scope NOT NULL DEFAULT 'school',
    is_active BOOLEAN DEFAULT true NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Comentarios
COMMENT ON TABLE permissions IS 'Catálogo maestro de permisos del sistema RBAC';
COMMENT ON COLUMN permissions.name IS 'Nombre único del permiso en formato resource:action (ej: users:create)';
COMMENT ON COLUMN permissions.resource IS 'Recurso sobre el que aplica el permiso (users, materials, schools, etc.)';
COMMENT ON COLUMN permissions.action IS 'Acción que se puede realizar sobre el recurso (create, read, update, delete, etc.)';
```

**Archivo NUEVO**: `postgres/migrations/constraints/013_create_permissions.sql`

```sql
-- ====================================================================
-- CONSTRAINTS: permissions
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Constraint: name debe seguir patrón resource:action
ALTER TABLE permissions ADD CONSTRAINT chk_permission_name_format 
    CHECK (name ~* '^[a-z_]+:[a-z_]+(:[a-z_]+)?$');

-- Índices
CREATE INDEX idx_permissions_name ON permissions(name);
CREATE INDEX idx_permissions_resource ON permissions(resource);
CREATE INDEX idx_permissions_scope ON permissions(scope);
CREATE INDEX idx_permissions_active ON permissions(is_active);

-- Trigger para updated_at
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

---

#### 4.3.4 CREAR: Tabla role_permissions

**Archivo NUEVO**: `postgres/migrations/structure/014_create_role_permissions.sql`

```sql
-- ====================================================================
-- TABLA: role_permissions
-- DESCRIPCIÓN: Relación N:N entre roles y permisos
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

CREATE TABLE role_permissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON TABLE role_permissions IS 'Relación N:N entre roles y permisos (RBAC)';
```

**Archivo NUEVO**: `postgres/migrations/constraints/014_create_role_permissions.sql`

```sql
-- ====================================================================
-- CONSTRAINTS: role_permissions
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Foreign Keys
ALTER TABLE role_permissions
    ADD CONSTRAINT fk_role_permissions_role
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

ALTER TABLE role_permissions
    ADD CONSTRAINT fk_role_permissions_permission
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;

-- Unique constraint: Un rol no puede tener el mismo permiso duplicado
ALTER TABLE role_permissions
    ADD CONSTRAINT uq_role_permission UNIQUE (role_id, permission_id);

-- Índices
CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission ON role_permissions(permission_id);
```

---

#### 4.3.5 CREAR: Tabla user_roles

**Archivo NUEVO**: `postgres/migrations/structure/015_create_user_roles.sql`

```sql
-- ====================================================================
-- TABLA: user_roles
-- DESCRIPCIÓN: Asignación de roles a usuarios en contextos específicos
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    school_id UUID,
    academic_unit_id UUID,
    is_active BOOLEAN DEFAULT true NOT NULL,
    granted_by UUID,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

COMMENT ON TABLE user_roles IS 'Asignación de roles a usuarios en contextos específicos (RBAC)';
COMMENT ON COLUMN user_roles.school_id IS 'Escuela en la que aplica el rol. NULL = rol a nivel sistema';
COMMENT ON COLUMN user_roles.academic_unit_id IS 'Unidad académica en la que aplica el rol. NULL = rol a nivel escuela';
COMMENT ON COLUMN user_roles.granted_by IS 'Usuario que otorgó el rol (auditoría)';
COMMENT ON COLUMN user_roles.expires_at IS 'Fecha de expiración del rol. NULL = no expira';
```

**Archivo NUEVO**: `postgres/migrations/constraints/015_create_user_roles.sql`

```sql
-- ====================================================================
-- CONSTRAINTS: user_roles
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Foreign Keys
ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_role
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_school
    FOREIGN KEY (school_id) REFERENCES schools(id) ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_unit
    FOREIGN KEY (academic_unit_id) REFERENCES academic_units(id) ON DELETE CASCADE;

ALTER TABLE user_roles
    ADD CONSTRAINT fk_user_roles_granted_by
    FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL;

-- Unique constraint
ALTER TABLE user_roles
    ADD CONSTRAINT uq_user_role_context 
    UNIQUE (user_id, role_id, school_id, academic_unit_id);

-- Check constraint: Si academic_unit_id está presente, school_id debe estar presente
ALTER TABLE user_roles
    ADD CONSTRAINT chk_user_roles_unit_requires_school 
    CHECK (academic_unit_id IS NULL OR school_id IS NOT NULL);

-- Índices
CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role_id);
CREATE INDEX idx_user_roles_school ON user_roles(school_id);
CREATE INDEX idx_user_roles_unit ON user_roles(academic_unit_id);
CREATE INDEX idx_user_roles_active ON user_roles(is_active);
CREATE INDEX idx_user_roles_expires ON user_roles(expires_at) WHERE expires_at IS NOT NULL;
CREATE INDEX idx_user_roles_user_active ON user_roles(user_id, is_active);
CREATE INDEX idx_user_roles_context ON user_roles(user_id, school_id, academic_unit_id);

-- Trigger para updated_at
CREATE TRIGGER set_updated_at
    BEFORE UPDATE ON user_roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
```

---

#### 4.3.6 CREAR: Funciones de utilidad RBAC

**Archivo**: `postgres/migrations/structure/000_create_functions.sql`

**Acción**: AGREGAR al final del archivo (después de las funciones existentes)

```sql
-- ====================================================================
-- FUNCIÓN: get_user_permissions
-- DESCRIPCIÓN: Obtiene permisos de un usuario en un contexto específico
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

CREATE OR REPLACE FUNCTION get_user_permissions(
    p_user_id UUID,
    p_school_id UUID DEFAULT NULL,
    p_unit_id UUID DEFAULT NULL
) RETURNS TABLE(permission_name VARCHAR, permission_scope permission_scope) AS $$
BEGIN
    RETURN QUERY
    SELECT DISTINCT p.name::VARCHAR, p.scope
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
    JOIN role_permissions rp ON r.id = rp.role_id
    JOIN permissions p ON rp.permission_id = p.id
    WHERE ur.user_id = p_user_id
      AND ur.is_active = true
      AND r.is_active = true
      AND p.is_active = true
      AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
      AND (
          -- Permisos a nivel sistema (sin contexto)
          (ur.school_id IS NULL AND p_school_id IS NULL)
          OR
          -- Permisos a nivel escuela (coincide school_id)
          (ur.school_id = p_school_id AND ur.academic_unit_id IS NULL AND p_unit_id IS NULL)
          OR
          -- Permisos a nivel unidad (coincide school_id y unit_id)
          (ur.school_id = p_school_id AND ur.academic_unit_id = p_unit_id)
          OR
          -- Permisos globales siempre aplican (super_admin)
          (ur.school_id IS NULL)
      );
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_user_permissions IS 'Obtiene lista de permisos de un usuario en un contexto específico';

-- ====================================================================
-- FUNCIÓN: user_has_permission
-- DESCRIPCIÓN: Verifica si un usuario tiene un permiso específico
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

CREATE OR REPLACE FUNCTION user_has_permission(
    p_user_id UUID,
    p_permission_name VARCHAR,
    p_school_id UUID DEFAULT NULL,
    p_unit_id UUID DEFAULT NULL
) RETURNS BOOLEAN AS $$
DECLARE
    has_perm BOOLEAN;
BEGIN
    SELECT EXISTS(
        SELECT 1
        FROM user_roles ur
        JOIN roles r ON ur.role_id = r.id
        JOIN role_permissions rp ON r.id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE ur.user_id = p_user_id
          AND p.name = p_permission_name
          AND ur.is_active = true
          AND r.is_active = true
          AND p.is_active = true
          AND (ur.expires_at IS NULL OR ur.expires_at > NOW())
          AND (
              (ur.school_id IS NULL)
              OR (ur.school_id = p_school_id AND ur.academic_unit_id IS NULL AND p_unit_id IS NULL)
              OR (ur.school_id = p_school_id AND ur.academic_unit_id = p_unit_id)
          )
    ) INTO has_perm;
    
    RETURN has_perm;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION user_has_permission IS 'Verifica si un usuario tiene un permiso específico en un contexto dado';
```

---

#### 4.3.7 CREAR: Seeds de roles

**Archivo NUEVO**: `postgres/migrations/seeds/002_seed_roles.sql`

```sql
-- ====================================================================
-- SEEDS: Roles predefinidos del sistema
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

INSERT INTO roles (id, name, display_name, description, scope, is_active) VALUES

-- Roles a nivel sistema
('10000000-0000-0000-0000-000000000001', 'super_admin', 'Super Administrador', 
 'Administrador con acceso total al sistema', 'system', true),

('10000000-0000-0000-0000-000000000002', 'platform_admin', 'Administrador de Plataforma', 
 'Administrador de plataforma con permisos de gestión global', 'system', true),

-- Roles a nivel escuela
('10000000-0000-0000-0000-000000000003', 'school_admin', 'Administrador de Escuela', 
 'Administrador con control total de la institución', 'school', true),

('10000000-0000-0000-0000-000000000004', 'school_director', 'Director', 
 'Director de la institución educativa', 'school', true),

('10000000-0000-0000-0000-000000000005', 'school_coordinator', 'Coordinador', 
 'Coordinador académico de la institución', 'school', true),

('10000000-0000-0000-0000-000000000006', 'school_assistant', 'Asistente Administrativo', 
 'Personal de soporte administrativo', 'school', true),

-- Roles a nivel unidad académica
('10000000-0000-0000-0000-000000000007', 'teacher', 'Profesor', 
 'Docente con permisos de gestión de clase', 'unit', true),

('10000000-0000-0000-0000-000000000008', 'assistant_teacher', 'Profesor Asistente', 
 'Asistente de docente', 'unit', true),

('10000000-0000-0000-0000-000000000009', 'student', 'Estudiante', 
 'Alumno inscrito en la unidad', 'unit', true),

('10000000-0000-0000-0000-000000000010', 'guardian', 'Apoderado', 
 'Tutor legal o apoderado de estudiante', 'unit', true),

('10000000-0000-0000-0000-000000000011', 'observer', 'Observador', 
 'Rol de solo lectura para auditoría', 'unit', true);
```

---

#### 4.3.8 CREAR: Seeds de permisos

**Archivo NUEVO**: `postgres/migrations/seeds/003_seed_permissions.sql`

```sql
-- ====================================================================
-- SEEDS: Permisos predefinidos del sistema
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- Permisos sobre USUARIOS
INSERT INTO permissions (name, display_name, description, resource, action, scope) VALUES
('users:create', 'Crear Usuarios', 'Crear nuevos usuarios en el sistema', 'users', 'create', 'system'),
('users:read', 'Ver Usuarios', 'Ver información de usuarios', 'users', 'read', 'school'),
('users:update', 'Editar Usuarios', 'Modificar datos de usuarios', 'users', 'update', 'school'),
('users:delete', 'Eliminar Usuarios', 'Eliminar usuarios del sistema', 'users', 'delete', 'system'),
('users:read:own', 'Ver Perfil Propio', 'Ver propio perfil de usuario', 'users', 'read:own', 'system'),
('users:update:own', 'Editar Perfil Propio', 'Modificar propio perfil', 'users', 'update:own', 'system'),

-- Permisos sobre ESCUELAS
('schools:create', 'Crear Escuelas', 'Crear nuevas instituciones educativas', 'schools', 'create', 'system'),
('schools:read', 'Ver Escuelas', 'Ver información de escuelas', 'schools', 'read', 'system'),
('schools:update', 'Editar Escuelas', 'Modificar datos de escuelas', 'schools', 'update', 'school'),
('schools:delete', 'Eliminar Escuelas', 'Eliminar escuelas del sistema', 'schools', 'delete', 'system'),
('schools:manage', 'Gestionar Escuela', 'Control total de la escuela', 'schools', 'manage', 'school'),

-- Permisos sobre UNIDADES ACADÉMICAS
('units:create', 'Crear Unidades', 'Crear unidades académicas (clases, grados)', 'units', 'create', 'school'),
('units:read', 'Ver Unidades', 'Ver unidades académicas', 'units', 'read', 'school'),
('units:update', 'Editar Unidades', 'Modificar unidades académicas', 'units', 'update', 'school'),
('units:delete', 'Eliminar Unidades', 'Eliminar unidades académicas', 'units', 'delete', 'school'),

-- Permisos sobre MATERIALES
('materials:create', 'Crear Materiales', 'Crear materiales educativos', 'materials', 'create', 'unit'),
('materials:read', 'Ver Materiales', 'Ver materiales educativos', 'materials', 'read', 'unit'),
('materials:update', 'Editar Materiales', 'Modificar materiales', 'materials', 'update', 'unit'),
('materials:delete', 'Eliminar Materiales', 'Eliminar materiales', 'materials', 'delete', 'unit'),
('materials:publish', 'Publicar Materiales', 'Publicar materiales para estudiantes', 'materials', 'publish', 'unit'),
('materials:download', 'Descargar Materiales', 'Descargar materiales educativos', 'materials', 'download', 'unit'),

-- Permisos sobre EVALUACIONES
('assessments:create', 'Crear Evaluaciones', 'Crear evaluaciones y exámenes', 'assessments', 'create', 'unit'),
('assessments:read', 'Ver Evaluaciones', 'Ver evaluaciones', 'assessments', 'read', 'unit'),
('assessments:update', 'Editar Evaluaciones', 'Modificar evaluaciones', 'assessments', 'update', 'unit'),
('assessments:delete', 'Eliminar Evaluaciones', 'Eliminar evaluaciones', 'assessments', 'delete', 'unit'),
('assessments:publish', 'Publicar Evaluaciones', 'Publicar evaluaciones para estudiantes', 'assessments', 'publish', 'unit'),
('assessments:grade', 'Calificar Evaluaciones', 'Calificar respuestas de estudiantes', 'assessments', 'grade', 'unit'),
('assessments:attempt', 'Rendir Evaluaciones', 'Intentar evaluaciones como estudiante', 'assessments', 'attempt', 'unit'),
('assessments:view_results', 'Ver Resultados', 'Ver resultados propios', 'assessments', 'view_results', 'unit'),

-- Permisos sobre PROGRESO
('progress:read', 'Ver Progreso', 'Ver progreso académico', 'progress', 'read', 'unit'),
('progress:update', 'Actualizar Progreso', 'Actualizar progreso de estudiantes', 'progress', 'update', 'unit'),
('progress:read:own', 'Ver Progreso Propio', 'Ver propio progreso', 'progress', 'read:own', 'unit'),

-- Permisos sobre ESTADÍSTICAS
('stats:global', 'Estadísticas Globales', 'Ver estadísticas de toda la plataforma', 'stats', 'global', 'system'),
('stats:school', 'Estadísticas de Escuela', 'Ver estadísticas de la institución', 'stats', 'school', 'school'),
('stats:unit', 'Estadísticas de Unidad', 'Ver estadísticas de la clase', 'stats', 'unit', 'unit');
```

---

#### 4.3.9 CREAR: Seeds de role_permissions

**Archivo NUEVO**: `postgres/migrations/seeds/004_seed_role_permissions.sql`

```sql
-- ====================================================================
-- SEEDS: Asignación de permisos a roles
-- VERSIÓN: postgres/v0.15.0
-- ====================================================================

-- SUPER_ADMIN: Todos los permisos
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'super_admin'),
    id
FROM permissions;

-- PLATFORM_ADMIN: Gestión de escuelas y usuarios
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'platform_admin'),
    id
FROM permissions
WHERE name IN (
    'users:create', 'users:read', 'users:update',
    'schools:create', 'schools:read', 'schools:update', 'schools:delete',
    'stats:global'
);

-- SCHOOL_ADMIN: Control total de escuela
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'school_admin'),
    id
FROM permissions
WHERE name IN (
    'users:read', 'users:update',
    'schools:read', 'schools:update', 'schools:manage',
    'units:create', 'units:read', 'units:update', 'units:delete',
    'materials:read', 'materials:update', 'materials:delete',
    'assessments:read', 'assessments:update', 'assessments:delete',
    'progress:read', 'progress:update',
    'stats:school'
);

-- TEACHER: Gestión de clase
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'teacher'),
    id
FROM permissions
WHERE name IN (
    'users:read:own', 'users:update:own',
    'units:read',
    'materials:create', 'materials:read', 'materials:update', 'materials:publish', 'materials:download',
    'assessments:create', 'assessments:read', 'assessments:update', 'assessments:publish', 'assessments:grade',
    'progress:read', 'progress:update',
    'stats:unit'
);

-- STUDENT: Consumo de contenido
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'student'),
    id
FROM permissions
WHERE name IN (
    'users:read:own', 'users:update:own',
    'materials:read', 'materials:download',
    'assessments:read', 'assessments:attempt', 'assessments:view_results',
    'progress:read:own'
);

-- GUARDIAN: Ver progreso de estudiantes
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'guardian'),
    id
FROM permissions
WHERE name IN (
    'users:read:own', 'users:update:own',
    'materials:read',
    'assessments:view_results',
    'progress:read'
);
```

---

#### 4.3.10 CREAR: Entities Go

**Archivo NUEVO**: `postgres/entities/role.go`

```go
package entities

import (
    "github.com/google/uuid"
    "time"
)

// Role representa un rol del sistema RBAC
type Role struct {
    ID          uuid.UUID `db:"id"`
    Name        string    `db:"name"`
    DisplayName string    `db:"display_name"`
    Description *string   `db:"description"`
    Scope       string    `db:"scope"` // 'system', 'school', 'unit'
    IsActive    bool      `db:"is_active"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}
```

**Archivo NUEVO**: `postgres/entities/permission.go`

```go
package entities

import (
    "github.com/google/uuid"
    "time"
)

// Permission representa un permiso del sistema RBAC
type Permission struct {
    ID          uuid.UUID `db:"id"`
    Name        string    `db:"name"`
    DisplayName string    `db:"display_name"`
    Description *string   `db:"description"`
    Resource    string    `db:"resource"`
    Action      string    `db:"action"`
    Scope       string    `db:"scope"` // 'system', 'school', 'unit'
    IsActive    bool      `db:"is_active"`
    CreatedAt   time.Time `db:"created_at"`
    UpdatedAt   time.Time `db:"updated_at"`
}
```

**Archivo NUEVO**: `postgres/entities/user_role.go`

```go
package entities

import (
    "github.com/google/uuid"
    "time"
)

// UserRole representa la asignación de un rol a un usuario en un contexto específico
type UserRole struct {
    ID             uuid.UUID  `db:"id"`
    UserID         uuid.UUID  `db:"user_id"`
    RoleID         uuid.UUID  `db:"role_id"`
    SchoolID       *uuid.UUID `db:"school_id"`        // NULL = rol a nivel sistema
    AcademicUnitID *uuid.UUID `db:"academic_unit_id"` // NULL = rol a nivel escuela
    IsActive       bool       `db:"is_active"`
    GrantedBy      *uuid.UUID `db:"granted_by"`
    GrantedAt      time.Time  `db:"granted_at"`
    ExpiresAt      *time.Time `db:"expires_at"` // NULL = no expira
    CreatedAt      time.Time  `db:"created_at"`
    UpdatedAt      time.Time  `db:"updated_at"`
}
```

---

### 4.4 Resumen de Archivos en edugo-infrastructure

| Tipo | Archivo | Acción |
|------|---------|--------|
| **Structure (Modificar)** | `postgres/migrations/structure/001_create_users.sql` | ELIMINAR línea del campo `role` |
| **Constraints (Modificar)** | `postgres/migrations/constraints/001_create_users.sql` | ELIMINAR índice `idx_users_role` |
| **Structure (Crear)** | `postgres/migrations/structure/012_create_roles.sql` | Crear tabla roles |
| **Constraints (Crear)** | `postgres/migrations/constraints/012_create_roles.sql` | Crear índices y trigger |
| **Structure (Crear)** | `postgres/migrations/structure/013_create_permissions.sql` | Crear tabla permissions |
| **Constraints (Crear)** | `postgres/migrations/constraints/013_create_permissions.sql` | Crear índices y trigger |
| **Structure (Crear)** | `postgres/migrations/structure/014_create_role_permissions.sql` | Crear tabla role_permissions |
| **Constraints (Crear)** | `postgres/migrations/constraints/014_create_role_permissions.sql` | Crear FK e índices |
| **Structure (Crear)** | `postgres/migrations/structure/015_create_user_roles.sql` | Crear tabla user_roles |
| **Constraints (Crear)** | `postgres/migrations/constraints/015_create_user_roles.sql` | Crear FK, índices y trigger |
| **Functions (Modificar)** | `postgres/migrations/structure/000_create_functions.sql` | AGREGAR 2 funciones RBAC |
| **Seeds (Crear)** | `postgres/migrations/seeds/002_seed_roles.sql` | Seed de roles |
| **Seeds (Crear)** | `postgres/migrations/seeds/003_seed_permissions.sql` | Seed de permisos |
| **Seeds (Crear)** | `postgres/migrations/seeds/004_seed_role_permissions.sql` | Seed de relación roles-permisos |
| **Entity (Crear)** | `postgres/entities/role.go` | Entity Role |
| **Entity (Crear)** | `postgres/entities/permission.go` | Entity Permission |
| **Entity (Crear)** | `postgres/entities/user_role.go` | Entity UserRole |

**Total**: 2 modificaciones + 14 nuevos archivos

---

## 5. Módulo: edugo-shared (auth, common, middleware)

### 5.1 Información del Módulo

**Ubicación**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared`

**Tipo**: Monorepo con 12 módulos Go independientes

**Módulos afectados por RBAC**:
- `auth/` (versión actual: `auth/v0.11.0`)
- `common/` (versión actual: `common/v0.9.0`)
- `middleware/gin/` (versión actual: `middleware/gin/v0.9.0`)

### 5.2 Proceso de Trabajo Git para este Módulo

```bash
# 1. Verificar estado
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared
git fetch origin
git checkout main && git pull origin main
git checkout dev && git pull origin dev
git diff main dev  # Verificar que contenido sea igual

# 2. Feature branch ya creado desde dev sincronizado
git checkout feature/rbac-auth-middleware  # YA EXISTE

# 3. Hacer modificaciones (ver secciones 5.3, 5.4, 5.5)

# 4. Commit y push
git add auth/ common/ middleware/
git commit -m "feat(auth,common,middleware): soporte RBAC con permisos en JWT"
git push origin feature/rbac-auth-middleware

# 5. Crear PR en GitHub: feature/rbac-auth-middleware → dev
# 6. Esperar aprobación y merge

# 7. Después del merge a dev, crear tags desde dev
git checkout dev
git pull origin dev
git tag auth/v0.12.0 -m "feat: JWT con contextos y permisos RBAC"
git tag common/v0.10.0 -m "feat: enum de permisos RBAC"
git tag middleware/gin/v0.10.0 -m "feat: middleware RequirePermission"
git push origin auth/v0.12.0
git push origin common/v0.10.0
git push origin middleware/gin/v0.10.0
```

### 5.3 Cambios en Módulo auth

#### 5.3.1 MODIFICAR: Claims JWT

**Archivo**: `auth/jwt.go`

**Cambio en estructura Claims** (líneas ~20-28):

```go
// ANTES
type Claims struct {
    UserID    string          `json:"user_id"`
    Email     string          `json:"email"`
    Role      enum.SystemRole `json:"role"`
    SchoolID  string          `json:"school_id,omitempty"`
    jwt.RegisteredClaims
}

// DESPUÉS
type Claims struct {
    UserID   string   `json:"user_id"`
    Email    string   `json:"email"`
    
    // Contexto activo
    ActiveContext *UserContext `json:"active_context,omitempty"`
    
    jwt.RegisteredClaims
}

// NUEVO tipo
type UserContext struct {
    RoleID           string   `json:"role_id"`
    RoleName         string   `json:"role_name"`
    SchoolID         string   `json:"school_id,omitempty"`
    SchoolName       string   `json:"school_name,omitempty"`
    AcademicUnitID   string   `json:"academic_unit_id,omitempty"`
    AcademicUnitName string   `json:"academic_unit_name,omitempty"`
    Permissions      []string `json:"permissions"`
}
```

**Rationale**: Solo incluimos `ActiveContext` en el token para evitar que crezca demasiado. Los contextos disponibles se obtienen via API.

---

#### 5.3.2 MODIFICAR: JWTManager

**Archivo**: `auth/jwt.go`

**Agregar nuevo método**:

```go
// GenerateTokenWithContext genera un JWT con contexto RBAC
func (m *JWTManager) GenerateTokenWithContext(
    userID, email string,
    activeContext *UserContext,
    expiresIn time.Duration,
) (string, error) {
    now := time.Now()
    expiresAt := now.Add(expiresIn)
    
    claims := Claims{
        UserID:        userID,
        Email:         email,
        ActiveContext: activeContext,
        RegisteredClaims: jwt.RegisteredClaims{
            ID:        uuid.New().String(),
            Issuer:    m.issuer,
            Subject:   userID,
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(expiresAt),
            NotBefore: jwt.NewNumericDate(now),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(m.secretKey)
}
```

**Mantener métodos legacy** para compatibilidad temporal:
- `GenerateToken()` - Versión antigua (deprecar)
- `GenerateTokenWithSchool()` - Versión antigua (deprecar)

**Marcar como deprecated**:
```go
// Deprecated: Usar GenerateTokenWithContext en su lugar
func (m *JWTManager) GenerateToken(...) {...}
```

---

### 5.4 Cambios en Módulo common

#### 5.4.1 CREAR: Enum de permisos

**Archivo NUEVO**: `common/types/enum/permission.go`

```go
package enum

type Permission string

// Permisos de usuarios
const (
    PermissionUsersCreate    Permission = "users:create"
    PermissionUsersRead      Permission = "users:read"
    PermissionUsersUpdate    Permission = "users:update"
    PermissionUsersDelete    Permission = "users:delete"
    PermissionUsersReadOwn   Permission = "users:read:own"
    PermissionUsersUpdateOwn Permission = "users:update:own"
)

// Permisos de escuelas
const (
    PermissionSchoolsCreate Permission = "schools:create"
    PermissionSchoolsRead   Permission = "schools:read"
    PermissionSchoolsUpdate Permission = "schools:update"
    PermissionSchoolsDelete Permission = "schools:delete"
    PermissionSchoolsManage Permission = "schools:manage"
)

// Permisos de unidades académicas
const (
    PermissionUnitsCreate Permission = "units:create"
    PermissionUnitsRead   Permission = "units:read"
    PermissionUnitsUpdate Permission = "units:update"
    PermissionUnitsDelete Permission = "units:delete"
)

// Permisos de materiales
const (
    PermissionMaterialsCreate   Permission = "materials:create"
    PermissionMaterialsRead     Permission = "materials:read"
    PermissionMaterialsUpdate   Permission = "materials:update"
    PermissionMaterialsDelete   Permission = "materials:delete"
    PermissionMaterialsPublish  Permission = "materials:publish"
    PermissionMaterialsDownload Permission = "materials:download"
)

// Permisos de evaluaciones
const (
    PermissionAssessmentsCreate      Permission = "assessments:create"
    PermissionAssessmentsRead        Permission = "assessments:read"
    PermissionAssessmentsUpdate      Permission = "assessments:update"
    PermissionAssessmentsDelete      Permission = "assessments:delete"
    PermissionAssessmentsPublish     Permission = "assessments:publish"
    PermissionAssessmentsGrade       Permission = "assessments:grade"
    PermissionAssessmentsAttempt     Permission = "assessments:attempt"
    PermissionAssessmentsViewResults Permission = "assessments:view_results"
)

// Permisos de progreso
const (
    PermissionProgressRead    Permission = "progress:read"
    PermissionProgressUpdate  Permission = "progress:update"
    PermissionProgressReadOwn Permission = "progress:read:own"
)

// Permisos de estadísticas
const (
    PermissionStatsGlobal Permission = "stats:global"
    PermissionStatsSchool Permission = "stats:school"
    PermissionStatsUnit   Permission = "stats:unit"
)

// String retorna la representación en string del permiso
func (p Permission) String() string {
    return string(p)
}

// IsValid verifica si el permiso es válido
func (p Permission) IsValid() bool {
    return AllPermissions[p]
}

// AllPermissions es un mapa de todos los permisos válidos
var AllPermissions = map[Permission]bool{
    // Usuarios
    PermissionUsersCreate:    true,
    PermissionUsersRead:      true,
    PermissionUsersUpdate:    true,
    PermissionUsersDelete:    true,
    PermissionUsersReadOwn:   true,
    PermissionUsersUpdateOwn: true,
    // Escuelas
    PermissionSchoolsCreate: true,
    PermissionSchoolsRead:   true,
    PermissionSchoolsUpdate: true,
    PermissionSchoolsDelete: true,
    PermissionSchoolsManage: true,
    // Unidades
    PermissionUnitsCreate: true,
    PermissionUnitsRead:   true,
    PermissionUnitsUpdate: true,
    PermissionUnitsDelete: true,
    // Materiales
    PermissionMaterialsCreate:   true,
    PermissionMaterialsRead:     true,
    PermissionMaterialsUpdate:   true,
    PermissionMaterialsDelete:   true,
    PermissionMaterialsPublish:  true,
    PermissionMaterialsDownload: true,
    // Evaluaciones
    PermissionAssessmentsCreate:      true,
    PermissionAssessmentsRead:        true,
    PermissionAssessmentsUpdate:      true,
    PermissionAssessmentsDelete:      true,
    PermissionAssessmentsPublish:     true,
    PermissionAssessmentsGrade:       true,
    PermissionAssessmentsAttempt:     true,
    PermissionAssessmentsViewResults: true,
    // Progreso
    PermissionProgressRead:    true,
    PermissionProgressUpdate:  true,
    PermissionProgressReadOwn: true,
    // Estadísticas
    PermissionStatsGlobal: true,
    PermissionStatsSchool: true,
    PermissionStatsUnit:   true,
}

// AllPermissionsSlice retorna todos los permisos como slice
func AllPermissionsSlice() []Permission {
    perms := make([]Permission, 0, len(AllPermissions))
    for perm := range AllPermissions {
        perms = append(perms, perm)
    }
    return perms
}
```

---

### 5.5 Cambios en Módulo middleware/gin

#### 5.5.1 CREAR: Middleware de permisos

**Archivo NUEVO**: `middleware/gin/permission_auth.go`

```go
package gin

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/EduGoGroup/edugo-shared/auth"
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// RequirePermission valida que el usuario tenga un permiso específico
func RequirePermission(permission enum.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Obtener claims del contexto (inyectados por JWTAuthMiddleware)
        claims, exists := c.Get("jwt_claims")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "unauthorized",
                "code":  "NO_CLAIMS",
            })
            c.Abort()
            return
        }
        
        userClaims, ok := claims.(*auth.Claims)
        if !ok {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "invalid claims type",
                "code":  "INVALID_CLAIMS_TYPE",
            })
            c.Abort()
            return
        }
        
        // Verificar si el contexto activo tiene el permiso
        if userClaims.ActiveContext == nil {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "no active context",
                "code":  "NO_ACTIVE_CONTEXT",
            })
            c.Abort()
            return
        }
        
        hasPermission := false
        for _, perm := range userClaims.ActiveContext.Permissions {
            if perm == permission.String() {
                hasPermission = true
                break
            }
        }
        
        if !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{
                "error":    "forbidden",
                "code":     "INSUFFICIENT_PERMISSIONS",
                "required": permission.String(),
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// RequireAnyPermission valida que el usuario tenga AL MENOS uno de los permisos
func RequireAnyPermission(permissions ...enum.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, exists := c.Get("jwt_claims")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        
        userClaims, ok := claims.(*auth.Claims)
        if !ok || userClaims.ActiveContext == nil {
            c.JSON(http.StatusForbidden, gin.H{"error": "no active context"})
            c.Abort()
            return
        }
        
        userPerms := make(map[string]bool)
        for _, perm := range userClaims.ActiveContext.Permissions {
            userPerms[perm] = true
        }
        
        for _, requiredPerm := range permissions {
            if userPerms[requiredPerm.String()] {
                c.Next()
                return
            }
        }
        
        c.JSON(http.StatusForbidden, gin.H{
            "error": "insufficient permissions",
            "code":  "INSUFFICIENT_PERMISSIONS",
        })
        c.Abort()
    }
}

// RequireAllPermissions valida que el usuario tenga TODOS los permisos
func RequireAllPermissions(permissions ...enum.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims, exists := c.Get("jwt_claims")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        
        userClaims, ok := claims.(*auth.Claims)
        if !ok || userClaims.ActiveContext == nil {
            c.JSON(http.StatusForbidden, gin.H{"error": "no active context"})
            c.Abort()
            return
        }
        
        userPerms := make(map[string]bool)
        for _, perm := range userClaims.ActiveContext.Permissions {
            userPerms[perm] = true
        }
        
        missing := []string{}
        for _, requiredPerm := range permissions {
            if !userPerms[requiredPerm.String()] {
                missing = append(missing, requiredPerm.String())
            }
        }
        
        if len(missing) > 0 {
            c.JSON(http.StatusForbidden, gin.H{
                "error":   "insufficient permissions",
                "code":    "INSUFFICIENT_PERMISSIONS",
                "missing": missing,
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

### 5.6 Resumen de Archivos en edugo-shared

| Módulo | Archivo | Acción |
|--------|---------|--------|
| **auth** | `auth/jwt.go` | MODIFICAR Claims + AGREGAR GenerateTokenWithContext |
| **common** | `common/types/enum/permission.go` | CREAR enum de permisos |
| **middleware/gin** | `middleware/gin/permission_auth.go` | CREAR middleware RequirePermission |

**Total**: 1 modificación + 2 nuevos archivos

---

## 6. Proyecto: edugo-api-administracion

### 6.1 Información del Proyecto

**Ubicación**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion`

**Tipo**: Proyecto Go (API REST)

**Swagger**: `/docs/swagger.json` (3,227 líneas)

**Comando de regeneración**:
```bash
make swagger  # Ejecuta: swag init -g cmd/main.go -o docs --parseInternal --parseDependency
```

### 6.2 Proceso de Trabajo Git

```bash
# 1. Verificar estado
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion
git fetch origin
git checkout main && git pull origin main
git checkout dev && git pull origin dev
git diff main dev

# 2. Feature branch ya creado desde dev sincronizado
git checkout feature/rbac-api-admin  # YA EXISTE

# 3. Hacer modificaciones (ver sección 6.3)

# 4. Regenerar swagger SIEMPRE
make swagger

# 5. Commit y push (incluir swagger.json)
git add .
git commit -m "feat: endpoints RBAC y autenticación con permisos"
git push origin feature/rbac-api-admin

# 6. Crear PR: feature/rbac-api-admin → dev
```

### 6.3 Archivos a Crear/Modificar

#### 6.3.1 CREAR: Repositorios RBAC

**Archivo NUEVO**: `internal/domain/repository/role_repository.go`

```go
package repository

import (
    "context"
    "github.com/google/uuid"
)

type Role struct {
    ID          uuid.UUID
    Name        string
    DisplayName string
    Description string
    Scope       string // 'system', 'school', 'unit'
    IsActive    bool
}

type RoleRepository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*Role, error)
    FindByName(ctx context.Context, name string) (*Role, error)
    FindAll(ctx context.Context) ([]*Role, error)
    FindByScope(ctx context.Context, scope string) ([]*Role, error)
    Create(ctx context.Context, role *Role) error
    Update(ctx context.Context, role *Role) error
}
```

**Archivo NUEVO**: `internal/domain/repository/permission_repository.go`

```go
package repository

import (
    "context"
    "github.com/google/uuid"
)

type Permission struct {
    ID          uuid.UUID
    Name        string
    DisplayName string
    Description string
    Resource    string
    Action      string
    Scope       string
    IsActive    bool
}

type PermissionRepository interface {
    FindByID(ctx context.Context, id uuid.UUID) (*Permission, error)
    FindByName(ctx context.Context, name string) (*Permission, error)
    FindAll(ctx context.Context) ([]*Permission, error)
    FindByRole(ctx context.Context, roleID uuid.UUID) ([]*Permission, error)
    FindByResource(ctx context.Context, resource string) ([]*Permission, error)
}
```

**Archivo NUEVO**: `internal/domain/repository/user_role_repository.go`

```go
package repository

import (
    "context"
    "github.com/google/uuid"
    "time"
)

type UserRole struct {
    ID             uuid.UUID
    UserID         uuid.UUID
    RoleID         uuid.UUID
    SchoolID       *uuid.UUID
    AcademicUnitID *uuid.UUID
    IsActive       bool
    GrantedBy      *uuid.UUID
    GrantedAt      time.Time
    ExpiresAt      *time.Time
}

type UserRoleRepository interface {
    // Buscar roles de un usuario
    FindByUser(ctx context.Context, userID uuid.UUID) ([]*UserRole, error)
    FindByUserInContext(ctx context.Context, userID uuid.UUID, schoolID *uuid.UUID, unitID *uuid.UUID) ([]*UserRole, error)
    
    // Gestión de roles
    Grant(ctx context.Context, userRole *UserRole) error
    Revoke(ctx context.Context, id uuid.UUID) error
    RevokeByUserAndRole(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) error
    
    // Verificación
    UserHasRole(ctx context.Context, userID, roleID uuid.UUID, schoolID, unitID *uuid.UUID) (bool, error)
    
    // Permisos
    GetUserPermissions(ctx context.Context, userID uuid.UUID, schoolID, unitID *uuid.UUID) ([]*Permission, error)
    UserHasPermission(ctx context.Context, userID uuid.UUID, permissionName string, schoolID, unitID *uuid.UUID) (bool, error)
}
```

---

#### 6.3.2 CREAR: Implementaciones PostgreSQL

**Archivo NUEVO**: `internal/infrastructure/persistence/postgres/role_repository.go`

(Implementación completa con queries SQL)

**Archivo NUEVO**: `internal/infrastructure/persistence/postgres/permission_repository.go`

(Implementación completa con queries SQL)

**Archivo NUEVO**: `internal/infrastructure/persistence/postgres/user_role_repository.go`

```go
package postgres

import (
    "context"
    "database/sql"
    "github.com/google/uuid"
    "edugo-api-administracion/internal/domain/repository"
)

type UserRoleRepository struct {
    db *sql.DB
}

func NewUserRoleRepository(db *sql.DB) repository.UserRoleRepository {
    return &UserRoleRepository{db: db}
}

func (r *UserRoleRepository) GetUserPermissions(
    ctx context.Context, 
    userID uuid.UUID, 
    schoolID, unitID *uuid.UUID,
) ([]*repository.Permission, error) {
    // Usar función PostgreSQL get_user_permissions()
    query := `SELECT permission_name, permission_scope FROM get_user_permissions($1, $2, $3)`
    
    rows, err := r.db.QueryContext(ctx, query, userID, schoolID, unitID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    permissions := []*repository.Permission{}
    for rows.Next() {
        perm := &repository.Permission{}
        if err := rows.Scan(&perm.Name, &perm.Scope); err != nil {
            return nil, err
        }
        permissions = append(permissions, perm)
    }
    
    return permissions, nil
}

// ... implementar resto de métodos
```

---

#### 6.3.3 MODIFICAR: AuthService

**Archivo**: `internal/auth/service/auth_service.go`

**MODIFICAR método Login()** (usar nuevo GenerateTokenWithContext):

```go
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
    // ... validación de usuario y password (sin cambios)
    
    // NUEVO: Obtener roles del usuario
    userRoles, err := s.userRoleRepo.FindByUser(ctx, user.ID)
    if err != nil {
        return nil, err
    }
    
    if len(userRoles) == 0 {
        return nil, errors.NewBusinessRuleError("user has no assigned roles")
    }
    
    // NUEVO: Determinar contexto activo por defecto (primer rol activo)
    activeContext, err := s.buildUserContext(ctx, userRoles[0])
    if err != nil {
        return nil, err
    }
    
    // Generar tokens con nuevo formato
    accessToken, expiresAt, err := s.jwtManager.GenerateTokenWithContext(
        user.ID.String(),
        user.Email,
        activeContext,
        s.jwtConfig.AccessTokenDuration,
    )
    if err != nil {
        return nil, err
    }
    
    // ... resto del código (refresh token, etc.)
}

// NUEVA función auxiliar
func (s *AuthService) buildUserContext(ctx context.Context, userRole *repository.UserRole) (*auth.UserContext, error) {
    // Obtener información del rol
    role, err := s.roleRepo.FindByID(ctx, userRole.RoleID)
    if err != nil {
        return nil, err
    }
    
    // Obtener permisos del usuario en este contexto
    permissions, err := s.userRoleRepo.GetUserPermissions(
        ctx, 
        userRole.UserID, 
        userRole.SchoolID, 
        userRole.AcademicUnitID,
    )
    if err != nil {
        return nil, err
    }
    
    permissionNames := []string{}
    for _, perm := range permissions {
        permissionNames = append(permissionNames, perm.Name)
    }
    
    // Obtener nombres de escuela y unidad (si aplica)
    var schoolName, unitName string
    if userRole.SchoolID != nil {
        school, _ := s.schoolRepo.FindByID(ctx, *userRole.SchoolID)
        if school != nil {
            schoolName = school.Name
        }
    }
    if userRole.AcademicUnitID != nil {
        unit, _ := s.unitRepo.FindByID(ctx, *userRole.AcademicUnitID)
        if unit != nil {
            unitName = unit.Name
        }
    }
    
    return &auth.UserContext{
        RoleID:           userRole.RoleID.String(),
        RoleName:         role.Name,
        SchoolID:         uuidPtrToString(userRole.SchoolID),
        SchoolName:       schoolName,
        AcademicUnitID:   uuidPtrToString(userRole.AcademicUnitID),
        AcademicUnitName: unitName,
        Permissions:      permissionNames,
    }, nil
}

func uuidPtrToString(id *uuid.UUID) string {
    if id == nil {
        return ""
    }
    return id.String()
}
```

**AGREGAR método SwitchContext()**:

```go
func (s *AuthService) SwitchContext(
    ctx context.Context,
    userID uuid.UUID,
    roleID uuid.UUID,
    schoolID *uuid.UUID,
    unitID *uuid.UUID,
) (*LoginResponse, error) {
    // Verificar que el usuario tiene ese rol en ese contexto
    userRoles, err := s.userRoleRepo.FindByUserInContext(ctx, userID, schoolID, unitID)
    if err != nil || len(userRoles) == 0 {
        return nil, errors.NewForbiddenError("user does not have access to this context")
    }
    
    // Filtrar por roleID específico
    var targetRole *repository.UserRole
    for _, ur := range userRoles {
        if ur.RoleID == roleID {
            targetRole = ur
            break
        }
    }
    
    if targetRole == nil {
        return nil, errors.NewForbiddenError("role not found in this context")
    }
    
    // Construir nuevo contexto activo
    activeContext, err := s.buildUserContext(ctx, targetRole)
    if err != nil {
        return nil, err
    }
    
    // Obtener email del usuario
    user, err := s.userRepo.FindByID(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    // Generar nuevos tokens
    accessToken, _, err := s.jwtManager.GenerateTokenWithContext(
        userID.String(),
        user.Email,
        activeContext,
        s.jwtConfig.AccessTokenDuration,
    )
    if err != nil {
        return nil, err
    }
    
    refreshToken, err := s.tokenService.CreateRefreshToken(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    return &LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int(s.jwtConfig.AccessTokenDuration.Seconds()),
        TokenType:    "Bearer",
        User: &UserInfo{
            ID:            user.ID.String(),
            Email:         user.Email,
            ActiveContext: activeContext,
        },
    }, nil
}
```

---

#### 6.3.4 CREAR: RoleHandler (nuevos endpoints)

**Archivo NUEVO**: `internal/application/handler/role_handler.go`

```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "edugo-api-administracion/internal/application/service"
    "edugo-api-administracion/internal/application/dto"
)

type RoleHandler struct {
    roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
    return &RoleHandler{roleService: roleService}
}

// RegisterRoutes registra las rutas del handler
func (h *RoleHandler) RegisterRoutes(r *gin.RouterGroup) {
    roles := r.Group("/roles")
    {
        roles.GET("", h.ListRoles)                         // GET /roles
        roles.GET("/:id", h.GetRole)                       // GET /roles/:id
        roles.GET("/:id/permissions", h.GetRolePermissions) // GET /roles/:id/permissions
    }
    
    userRoles := r.Group("/users/:user_id/roles")
    {
        userRoles.GET("", h.GetUserRoles)          // GET /users/:user_id/roles
        userRoles.POST("", h.GrantRole)            // POST /users/:user_id/roles
        userRoles.DELETE("/:role_id", h.RevokeRole) // DELETE /users/:user_id/roles/:role_id
    }
}

// ListRoles godoc
// @Summary Listar roles
// @Description Obtiene lista de todos los roles del sistema
// @Tags roles
// @Accept json
// @Produce json
// @Param scope query string false "Filtrar por scope (system, school, unit)"
// @Success 200 {object} dto.RolesResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
    scope := c.Query("scope")
    
    roles, err := h.roleService.GetRoles(c.Request.Context(), scope)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, dto.RolesResponse{Roles: roles})
}

// GetUserRoles godoc
// @Summary Obtener roles de un usuario
// @Description Obtiene todos los roles asignados a un usuario
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} dto.UserRolesResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{user_id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
    userIDStr := c.Param("user_id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
        return
    }
    
    userRoles, err := h.roleService.GetUserRoles(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, dto.UserRolesResponse{UserRoles: userRoles})
}

// GrantRole godoc
// @Summary Asignar rol a usuario
// @Description Otorga un rol a un usuario en un contexto específico
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body dto.GrantRoleRequest true "Datos del rol a otorgar"
// @Success 201 {object} dto.GrantRoleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{user_id}/roles [post]
func (h *RoleHandler) GrantRole(c *gin.Context) {
    userIDStr := c.Param("user_id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
        return
    }
    
    var req dto.GrantRoleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    // Obtener ID del usuario que otorga el rol (del token JWT)
    grantedBy := c.GetString("user_id")
    
    userRole, err := h.roleService.GrantRoleToUser(c.Request.Context(), userID, &req, grantedBy)
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, dto.GrantRoleResponse{UserRole: userRole})
}

// RevokeRole godoc
// @Summary Revocar rol de usuario
// @Description Revoca un rol asignado a un usuario
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param role_id path string true "Role ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/users/{user_id}/roles/{role_id} [delete]
func (h *RoleHandler) RevokeRole(c *gin.Context) {
    userIDStr := c.Param("user_id")
    roleIDStr := c.Param("role_id")
    
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid user_id"})
        return
    }
    
    roleID, err := uuid.Parse(roleIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid role_id"})
        return
    }
    
    if err := h.roleService.RevokeRoleFromUser(c.Request.Context(), userID, roleID); err != nil {
        c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    c.Status(http.StatusNoContent)
}

// ... resto de handlers
```

---

#### 6.3.5 CREAR: AuthHandler - endpoint SwitchContext

**Archivo**: `internal/auth/handler/auth_handler.go`

**AGREGAR método**:

```go
// SwitchContext godoc
// @Summary Cambiar contexto activo
// @Description Cambia el contexto activo del usuario (escuela/unidad/rol)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.SwitchContextRequest true "Nuevo contexto"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse "Usuario no tiene acceso a ese contexto"
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /v1/auth/switch-context [post]
func (h *AuthHandler) SwitchContext(c *gin.Context) {
    var req dto.SwitchContextRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    // Obtener user_id del token JWT
    userIDStr := c.GetString("user_id")
    userID, err := uuid.Parse(userIDStr)
    if err != nil {
        c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "invalid token"})
        return
    }
    
    roleID, err := uuid.Parse(req.RoleID)
    if err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid role_id"})
        return
    }
    
    var schoolID, unitID *uuid.UUID
    if req.SchoolID != "" {
        sid, err := uuid.Parse(req.SchoolID)
        if err != nil {
            c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid school_id"})
            return
        }
        schoolID = &sid
    }
    if req.AcademicUnitID != "" {
        uid, err := uuid.Parse(req.AcademicUnitID)
        if err != nil {
            c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid academic_unit_id"})
            return
        }
        unitID = &uid
    }
    
    response, err := h.authService.SwitchContext(c.Request.Context(), userID, roleID, schoolID, unitID)
    if err != nil {
        c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, response)
}
```

---

#### 6.3.6 MODIFICAR: Router (cmd/main.go)

**Archivo**: `cmd/main.go`

**MODIFICAR** para agregar middleware de permisos y nuevas rutas:

```go
import (
    ginmiddleware "github.com/EduGoGroup/edugo-shared/middleware/gin"
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

func main() {
    // ... inicialización
    
    // Rutas públicas (sin JWT)
    v1Public := r.Group("/v1")
    {
        c.AuthHandler.RegisterRoutes(v1Public)  // login, verify
    }
    
    // Rutas protegidas (requieren JWT)
    v1 := r.Group("/v1")
    v1.Use(ginmiddleware.JWTAuthMiddleware(c.JWTManager))
    {
        // Auth protected
        auth := v1.Group("/auth")
        {
            auth.POST("/logout", c.AuthHandler.Logout)
            auth.POST("/switch-context", c.AuthHandler.SwitchContext) // NUEVO
        }
        
        // Schools - CON permisos
        schools := v1.Group("/schools")
        {
            schools.POST("", 
                ginmiddleware.RequirePermission(enum.PermissionSchoolsCreate),
                c.SchoolHandler.CreateSchool,
            )
            schools.GET("", 
                ginmiddleware.RequirePermission(enum.PermissionSchoolsRead),
                c.SchoolHandler.ListSchools,
            )
            schools.PUT("/:id", 
                ginmiddleware.RequirePermission(enum.PermissionSchoolsUpdate),
                c.SchoolHandler.UpdateSchool,
            )
            schools.DELETE("/:id", 
                ginmiddleware.RequirePermission(enum.PermissionSchoolsDelete),
                c.SchoolHandler.DeleteSchool,
            )
        }
        
        // Users - CON permisos
        users := v1.Group("/users")
        {
            users.POST("", 
                ginmiddleware.RequirePermission(enum.PermissionUsersCreate),
                c.UserHandler.CreateUser,
            )
            users.GET("", 
                ginmiddleware.RequirePermission(enum.PermissionUsersRead),
                c.UserHandler.ListUsers,
            )
            users.PUT("/:id", 
                ginmiddleware.RequirePermission(enum.PermissionUsersUpdate),
                c.UserHandler.UpdateUser,
            )
        }
        
        // Roles - NUEVO
        c.RoleHandler.RegisterRoutes(v1)
        
        // Materials - CON permisos
        materials := v1.Group("/materials")
        {
            materials.POST("", 
                ginmiddleware.RequirePermission(enum.PermissionMaterialsCreate),
                c.MaterialHandler.CreateMaterial,
            )
            materials.PUT("/:id/publish", 
                ginmiddleware.RequirePermission(enum.PermissionMaterialsPublish),
                c.MaterialHandler.PublishMaterial,
            )
        }
    }
    
    // ... resto del código
}
```

---

#### 6.3.7 CREAR: DTOs

**Archivo NUEVO**: `internal/application/dto/role_dto.go`

```go
package dto

type RoleDTO struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    DisplayName string `json:"display_name"`
    Description string `json:"description,omitempty"`
    Scope       string `json:"scope"`
    IsActive    bool   `json:"is_active"`
}

type RolesResponse struct {
    Roles []*RoleDTO `json:"roles"`
}

type UserRoleDTO struct {
    ID             string  `json:"id"`
    UserID         string  `json:"user_id"`
    RoleID         string  `json:"role_id"`
    RoleName       string  `json:"role_name"`
    SchoolID       *string `json:"school_id,omitempty"`
    SchoolName     *string `json:"school_name,omitempty"`
    AcademicUnitID *string `json:"academic_unit_id,omitempty"`
    AcademicUnitName *string `json:"academic_unit_name,omitempty"`
    IsActive       bool    `json:"is_active"`
    GrantedAt      string  `json:"granted_at"`
}

type UserRolesResponse struct {
    UserRoles []*UserRoleDTO `json:"user_roles"`
}

type GrantRoleRequest struct {
    RoleID         string  `json:"role_id" binding:"required"`
    SchoolID       *string `json:"school_id,omitempty"`
    AcademicUnitID *string `json:"academic_unit_id,omitempty"`
    ExpiresAt      *string `json:"expires_at,omitempty"`
}

type GrantRoleResponse struct {
    UserRole *UserRoleDTO `json:"user_role"`
}
```

**Archivo NUEVO**: `internal/application/dto/auth_dto.go`

**MODIFICAR LoginResponse** para incluir ActiveContext:

```go
type LoginResponse struct {
    AccessToken  string      `json:"access_token"`
    RefreshToken string      `json:"refresh_token"`
    ExpiresIn    int         `json:"expires_in"`
    TokenType    string      `json:"token_type"`
    User         *UserInfo   `json:"user"`
}

type UserInfo struct {
    ID            string                `json:"id"`
    Email         string                `json:"email"`
    ActiveContext *UserContextDTO       `json:"active_context,omitempty"` // NUEVO
}

type UserContextDTO struct {
    RoleID           string   `json:"role_id"`
    RoleName         string   `json:"role_name"`
    SchoolID         string   `json:"school_id,omitempty"`
    SchoolName       string   `json:"school_name,omitempty"`
    AcademicUnitID   string   `json:"academic_unit_id,omitempty"`
    AcademicUnitName string   `json:"academic_unit_name,omitempty"`
    Permissions      []string `json:"permissions"`
}

type SwitchContextRequest struct {
    RoleID         string `json:"role_id" binding:"required"`
    SchoolID       string `json:"school_id,omitempty"`
    AcademicUnitID string `json:"academic_unit_id,omitempty"`
}
```

---

### 6.4 Resumen de Archivos en edugo-api-administracion

| Tipo | Archivo | Acción |
|------|---------|--------|
| **Repository (Crear)** | `internal/domain/repository/role_repository.go` | Interface |
| **Repository (Crear)** | `internal/domain/repository/permission_repository.go` | Interface |
| **Repository (Crear)** | `internal/domain/repository/user_role_repository.go` | Interface |
| **Infrastructure (Crear)** | `internal/infrastructure/persistence/postgres/role_repository.go` | Implementación |
| **Infrastructure (Crear)** | `internal/infrastructure/persistence/postgres/permission_repository.go` | Implementación |
| **Infrastructure (Crear)** | `internal/infrastructure/persistence/postgres/user_role_repository.go` | Implementación |
| **Service (Crear)** | `internal/application/service/role_service.go` | Lógica de negocio |
| **Service (Modificar)** | `internal/auth/service/auth_service.go` | Actualizar Login + agregar SwitchContext |
| **Handler (Crear)** | `internal/application/handler/role_handler.go` | Endpoints de roles |
| **Handler (Modificar)** | `internal/auth/handler/auth_handler.go` | Agregar SwitchContext |
| **DTO (Crear)** | `internal/application/dto/role_dto.go` | DTOs de roles |
| **DTO (Modificar)** | `internal/application/dto/auth_dto.go` | Actualizar LoginResponse |
| **Router (Modificar)** | `cmd/main.go` | Agregar rutas y middleware de permisos |
| **Swagger (Regenerar)** | `docs/swagger.json` | make swagger |

**Total**: 8 nuevos archivos + 5 modificaciones + regenerar swagger

---

## 7. Proyecto: edugo-api-mobile

### 7.1 Información del Proyecto

**Ubicación**: `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile`

**Tipo**: Proyecto Go (API REST)

**Swagger**: `/docs/swagger.json` (1,540 líneas)

**Comando de regeneración**:
```bash
make swagger  # Ejecuta: swag init -g cmd/main.go -o docs --parseInternal
```

### 7.2 Proceso de Trabajo Git

```bash
# 1. Verificar estado
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile
git fetch origin
git checkout main && git pull origin main
git checkout dev && git pull origin dev
git diff main dev

# 2. Feature branch ya creado desde dev sincronizado
git checkout feature/rbac-api-mobile  # YA EXISTE

# 3. Hacer modificaciones (ver sección 7.3)

# 4. Regenerar swagger SIEMPRE
make swagger

# 5. Commit y push (incluir swagger.json)
git add .
git commit -m "feat: autorización con permisos RBAC"
git push origin feature/rbac-api-mobile

# 6. Crear PR: feature/rbac-api-mobile → dev
```

### 7.3 Archivos a Modificar

#### 7.3.1 MODIFICAR: AuthClient

**Archivo**: `internal/client/auth_client.go`

**MODIFICAR** para soportar nuevo formato de Claims:

```go
func (c *AuthClient) VerifyToken(ctx context.Context, token string) (*TokenInfo, error) {
    // Intentar validación local con JWTManager (ya actualizado en shared)
    claims, err := c.jwtManager.ValidateToken(token)
    if err == nil {
        // Extraer permisos del active context
        permissions := []string{}
        if claims.ActiveContext != nil {
            permissions = claims.ActiveContext.Permissions
        }
        
        return &TokenInfo{
            UserID:      claims.UserID,
            Email:       claims.Email,
            Permissions: permissions,  // NUEVO campo
            ExpiresAt:   claims.ExpiresAt.Time,
        }, nil
    }
    
    // Fallback: Llamar a api-admin /v1/auth/verify
    // ... (sin cambios en lógica de fallback)
}
```

---

#### 7.3.2 MODIFICAR: Middleware de autorización

**Archivo**: `internal/infrastructure/http/middleware/remote_auth.go`

**AGREGAR** nuevos middleware basados en permisos:

```go
import (
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// RequirePermission valida que el usuario tenga un permiso específico
func RequirePermission(permission enum.Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Obtener claims del contexto
        claims, exists := c.Get("jwt_claims")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        
        userClaims, ok := claims.(*auth.Claims)
        if !ok || userClaims.ActiveContext == nil {
            c.JSON(http.StatusForbidden, gin.H{"error": "no active context"})
            c.Abort()
            return
        }
        
        // Verificar permiso
        hasPermission := false
        for _, perm := range userClaims.ActiveContext.Permissions {
            if perm == permission.String() {
                hasPermission = true
                break
            }
        }
        
        if !hasPermission {
            c.JSON(http.StatusForbidden, gin.H{
                "error":    "forbidden",
                "required": permission.String(),
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

#### 7.3.3 DEPRECAR: Middleware basado en roles

**Archivo**: `internal/infrastructure/http/middleware/auth.go`

**MARCAR como deprecated**:

```go
// Deprecated: Usar RequirePermission() en su lugar
// Este middleware será removido en v1.0.0
func RequireTeacher() gin.HandlerFunc {
    // ... código existente
}

// Deprecated: Usar RequirePermission() en su lugar
func RequireAdmin() gin.HandlerFunc {
    // ... código existente
}
```

---

#### 7.3.4 MODIFICAR: Router

**Archivo**: `internal/infrastructure/http/router/router.go`

**REEMPLAZAR** middleware de roles por permisos:

```go
import (
    "github.com/EduGoGroup/edugo-shared/common/types/enum"
)

// ANTES (basado en roles)
materials.POST("", 
    middleware.RequireTeacher(),  // ❌ Basado en roles
    h.materialHandler.CreateMaterial,
)

// DESPUÉS (basado en permisos)
materials.POST("", 
    middleware.RequirePermission(enum.PermissionMaterialsCreate),  // ✅ Basado en permisos
    h.materialHandler.CreateMaterial,
)

materials.POST("/:id/upload-complete", 
    middleware.RequirePermission(enum.PermissionMaterialsPublish),
    h.materialHandler.NotifyUploadComplete,
)

stats.GET("/global", 
    middleware.RequirePermission(enum.PermissionStatsGlobal),
    h.statsHandler.GetGlobalStats,
)
```

---

### 7.4 Resumen de Archivos en edugo-api-mobile

| Tipo | Archivo | Acción |
|------|---------|--------|
| **Client (Modificar)** | `internal/client/auth_client.go` | Soportar ActiveContext |
| **Middleware (Modificar)** | `internal/infrastructure/http/middleware/remote_auth.go` | Agregar RequirePermission |
| **Middleware (Modificar)** | `internal/infrastructure/http/middleware/auth.go` | Deprecar RequireTeacher/Admin |
| **Router (Modificar)** | `internal/infrastructure/http/router/router.go` | Usar RequirePermission |
| **Swagger (Regenerar)** | `docs/swagger.json` | make swagger |

**Total**: 4 modificaciones + regenerar swagger

---

## 8. Plan de Implementación por Fases

### FASE 1: Infraestructura de Base de Datos (Semana 1)

**Proyecto**: `edugo-infrastructure/postgres`

**Rama**: `feature/rbac-postgres-tables`

**Tareas**:
1. ✅ Modificar `structure/001_create_users.sql` (eliminar campo role)
2. ✅ Modificar `constraints/001_create_users.sql` (eliminar índice)
3. ✅ Crear tablas RBAC (012-015) en structure/
4. ✅ Crear constraints RBAC (012-015) en constraints/
5. ✅ Crear funciones PostgreSQL (en 000_create_functions.sql)
6. ✅ Crear seeds (002-004)
7. ✅ Crear entities Go (role.go, permission.go, user_role.go)
8. ✅ Validar con runner de 4 capas
9. ✅ Commit, PR a dev, merge
10. ✅ Tag: `postgres/v0.15.0`

**Entregable**: Base de datos RBAC funcional + tag liberado

---

### FASE 2: Actualización de Shared (Semana 2)

**Proyecto**: `edugo-shared`

**Rama**: `feature/rbac-auth-middleware`

**Tareas**:
1. ✅ Modificar `auth/jwt.go` (Claims con ActiveContext)
2. ✅ Agregar `GenerateTokenWithContext()` en JWTManager
3. ✅ Crear `common/types/enum/permission.go`
4. ✅ Crear `middleware/gin/permission_auth.go`
5. ✅ Tests unitarios
6. ✅ Commit, PR a dev, merge
7. ✅ Tags: `auth/v0.12.0`, `common/v0.10.0`, `middleware/gin/v0.10.0`

**Entregable**: Librerías shared actualizadas + tags liberados

---

### FASE 3: API Administración - Backend (Semana 3-4)

**Proyecto**: `edugo-api-administracion`

**Rama**: `feature/rbac-api-admin`

**Tareas**:
1. ✅ Actualizar `go.mod` para usar nuevas versiones de shared
2. ✅ Crear repositorios RBAC (interfaces + implementaciones PostgreSQL)
3. ✅ Crear RoleService
4. ✅ Modificar AuthService (Login + SwitchContext)
5. ✅ Crear RoleHandler
6. ✅ Modificar AuthHandler (SwitchContext endpoint)
7. ✅ Crear DTOs
8. ✅ Modificar router en `cmd/main.go`
9. ✅ **Regenerar swagger**: `make swagger`
10. ✅ Tests unitarios e integración
11. ✅ Commit (incluir swagger.json), PR a dev, merge

**Entregable**: API Admin con RBAC completo

---

### FASE 4: API Mobile (Semana 5)

**Proyecto**: `edugo-api-mobile`

**Rama**: `feature/rbac-api-mobile`

**Tareas**:
1. ✅ Actualizar `go.mod` para usar nuevas versiones de shared
2. ✅ Modificar AuthClient
3. ✅ Modificar middleware remote_auth.go (agregar RequirePermission)
4. ✅ Deprecar middleware auth.go (RequireTeacher/Admin)
5. ✅ Modificar router.go (usar RequirePermission)
6. ✅ **Regenerar swagger**: `make swagger`
7. ✅ Tests de integración
8. ✅ Commit (incluir swagger.json), PR a dev, merge

**Entregable**: API Mobile con autorización RBAC

---

### FASE 5: Testing End-to-End (Semana 6)

**Todos los proyectos**

**Tareas**:
1. ✅ Testing de flujo completo:
   - Login con múltiples roles
   - Switch context entre escuelas/unidades
   - Verificación de permisos en endpoints protegidos
2. ✅ Testing de casos edge:
   - Usuario sin roles
   - Rol expirado
   - Permisos insuficientes
3. ✅ Performance testing (queries de permisos)
4. ✅ Actualizar documentación de API

**Entregable**: Sistema RBAC validado end-to-end

---

## 9. Flujo de Releases

### 9.1 Release de Módulo PostgreSQL

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Desde dev (después de merge del PR)
git checkout dev
git pull origin dev

# Crear tag
git tag postgres/v0.15.0 -m "feat: Sistema RBAC con roles, permisos y user_roles"
git push origin postgres/v0.15.0

# GitHub Actions ejecutará:
# - Validación de todos los módulos
# - Tests
# - Creación de GitHub Release
```

**Consumo en otros proyectos**:
```bash
go get github.com/EduGoGroup/edugo-infrastructure/postgres@postgres/v0.15.0
```

---

### 9.2 Release de Módulos Shared

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-shared

git checkout dev
git pull origin dev

# Crear tags simultáneos (porque están relacionados)
git tag auth/v0.12.0 -m "feat: JWT con contextos y permisos RBAC"
git tag common/v0.10.0 -m "feat: enum de permisos RBAC"
git tag middleware/gin/v0.10.0 -m "feat: middleware RequirePermission"

git push origin auth/v0.12.0
git push origin common/v0.10.0
git push origin middleware/gin/v0.10.0
```

**Consumo en otros proyectos**:
```bash
go get github.com/EduGoGroup/edugo-shared/auth@auth/v0.12.0
go get github.com/EduGoGroup/edugo-shared/common@common/v0.10.0
go get github.com/EduGoGroup/edugo-shared/middleware/gin@middleware/gin/v0.10.0
```

---

### 9.3 Actualizar go.mod en API Admin

**Archivo**: `edugo-api-administracion/go.mod`

**NOTA**: Durante desarrollo local con go.work no es necesario actualizar go.mod. Solo actualizar al momento de crear el PR (con `GOWORK=off`).

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-administracion

# Actualizar dependencias (ejecutar con GOWORK=off para que tome los tags de GitHub)
GOWORK=off go get github.com/EduGoGroup/edugo-infrastructure/postgres@postgres/v0.15.0
GOWORK=off go get github.com/EduGoGroup/edugo-shared/auth@auth/v0.12.0
GOWORK=off go get github.com/EduGoGroup/edugo-shared/common@common/v0.10.0
GOWORK=off go get github.com/EduGoGroup/edugo-shared/middleware/gin@middleware/gin/v0.10.0

GOWORK=off go mod tidy
```

---

### 9.4 Actualizar go.mod en API Mobile

**NOTA**: Durante desarrollo local con go.work no es necesario actualizar go.mod. Solo actualizar al momento de crear el PR (con `GOWORK=off`).

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-api-mobile

# Actualizar dependencias (ejecutar con GOWORK=off para que tome los tags de GitHub)
GOWORK=off go get github.com/EduGoGroup/edugo-shared/auth@auth/v0.12.0
GOWORK=off go get github.com/EduGoGroup/edugo-shared/common@common/v0.10.0
GOWORK=off go get github.com/EduGoGroup/edugo-shared/middleware/gin@middleware/gin/v0.10.0

GOWORK=off go mod tidy
```

---

## 10. Casos de Uso y Ejemplos

### 10.1 Login con Múltiples Roles

**Request**:
```http
POST /v1/auth/login
Content-Type: application/json

{
  "email": "juan.perez@edugo.com",
  "password": "password123"
}
```

**Response**:
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "refresh_token_here",
  "expires_in": 900,
  "token_type": "Bearer",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "juan.perez@edugo.com",
    "active_context": {
      "role_id": "10000000-0000-0000-0000-000000000007",
      "role_name": "teacher",
      "school_id": "school-uuid-1",
      "school_name": "Colegio Los Alamos",
      "academic_unit_id": "unit-uuid-5",
      "academic_unit_name": "Matemáticas 3° Básico",
      "permissions": [
        "users:read:own",
        "users:update:own",
        "units:read",
        "materials:create",
        "materials:read",
        "materials:update",
        "materials:publish",
        "materials:download",
        "assessments:create",
        "assessments:read",
        "assessments:update",
        "assessments:publish",
        "assessments:grade",
        "progress:read",
        "progress:update",
        "stats:unit"
      ]
    }
  }
}
```

---

### 10.2 Cambiar Contexto

**Request**:
```http
POST /v1/auth/switch-context
Authorization: Bearer {access_token}
Content-Type: application/json

{
  "role_id": "10000000-0000-0000-0000-000000000009",
  "school_id": "school-uuid-2",
  "academic_unit_id": "unit-uuid-8"
}
```

**Response**: Nuevo token con contexto de estudiante

---

### 10.3 Endpoint Protegido con Permisos

**Request**:
```http
POST /v1/materials
Authorization: Bearer {teacher_token}
Content-Type: application/json

{
  "title": "Ecuaciones Cuadráticas",
  "description": "Material de apoyo",
  "academic_unit_id": "unit-uuid-5"
}
```

**Flujo**:
1. `JWTAuthMiddleware` valida token
2. `RequirePermission("materials:create")` verifica permisos
3. Si tiene permiso → continúa
4. Si NO tiene permiso → 403 Forbidden

---

## 11. Testing

### 11.1 Tests Unitarios

**edugo-shared**:
- `auth/jwt_test.go`: Generación y validación de tokens con ActiveContext
- `middleware/gin/permission_auth_test.go`: Middleware de permisos

**edugo-api-administracion**:
- `repository/*_test.go`: CRUD de roles, permissions, user_roles
- `service/auth_service_test.go`: Login con múltiples roles, SwitchContext
- `service/role_service_test.go`: Asignación/revocación de roles

---

### 11.2 Tests de Integración

```go
func TestE2E_MultiRoleLogin(t *testing.T) {
    // 1. Login como usuario con múltiples roles
    loginResp := loginUser("juan.perez@edugo.com", "password")
    assert.NotNil(t, loginResp.User.ActiveContext)
    assert.Contains(t, loginResp.User.ActiveContext.Permissions, "materials:create")
    
    // 2. Crear material (requiere materials:create)
    token := loginResp.AccessToken
    materials := createMaterial(token, materialData)
    assert.Equal(t, http.StatusCreated, materials.StatusCode)
    
    // 3. Cambiar a contexto sin permisos de creación (student)
    switchResp := switchContext(token, studentRoleID, schoolID, unitID)
    newToken := switchResp.AccessToken
    
    // 4. Intentar crear material (debería fallar)
    materialsResp := createMaterial(newToken, materialData)
    assert.Equal(t, http.StatusForbidden, materialsResp.StatusCode)
    assert.Equal(t, "INSUFFICIENT_PERMISSIONS", materialsResp.Code)
}
```

---

## 12. Riesgos y Mitigación

| Riesgo | Probabilidad | Impacto | Mitigación |
|--------|--------------|---------|------------|
| Tokens JWT demasiado grandes (> 8KB) | Media | Alto | Solo incluir `ActiveContext` en token. Obtener contextos disponibles via API `/users/me/contexts`. |
| Performance de queries de permisos | Media | Medio | Índices compuestos en `user_roles`. Funciones PostgreSQL optimizadas. |
| Breaking changes en JWT para clientes existentes | Alta | Alto | **NO HAY CLIENTES EN PRODUCCIÓN**. Podemos hacer cambio limpio. |
| Olvidar regenerar swagger | Alta | Medio | Hacer parte obligatoria del checklist de PR. Validar en CI que swagger está actualizado. |
| Desincronización de versiones entre módulos | Media | Alto | Documentar claramente qué versiones de shared requiere cada API. Actualizar `go.mod` en mismo PR. |

---

## CHECKLIST FINAL ANTES DE CREAR PR

### Para edugo-infrastructure
- [ ] Modificar structure/001_create_users.sql (eliminar campo role)
- [ ] Modificar constraints/001_create_users.sql (eliminar índice)
- [ ] Crear 4 nuevos structure/ (012-015)
- [ ] Crear 4 nuevos constraints/ (012-015)
- [ ] Modificar 000_create_functions.sql (agregar 2 funciones)
- [ ] Crear 3 seeds/ (002-004)
- [ ] Crear 3 entities/ (.go)
- [ ] Probar con runner de 4 capas localmente
- [ ] Verificar que main y dev tienen mismo contenido
- [ ] PR a dev (NO a main)

### Para edugo-shared
- [ ] Modificar auth/jwt.go
- [ ] Crear common/types/enum/permission.go
- [ ] Crear middleware/gin/permission_auth.go
- [ ] Tests unitarios pasando
- [ ] Verificar que main y dev tienen mismo contenido
- [ ] PR a dev (NO a main)

### Para edugo-api-administracion
- [ ] Actualizar go.mod con nuevas versiones de shared e infrastructure
- [ ] Crear 3 repositorios (interfaces)
- [ ] Crear 3 implementaciones PostgreSQL
- [ ] Crear RoleService
- [ ] Modificar AuthService
- [ ] Crear RoleHandler
- [ ] Modificar AuthHandler
- [ ] Crear DTOs
- [ ] Modificar cmd/main.go
- [ ] **REGENERAR SWAGGER**: `make swagger`
- [ ] **VERIFICAR** que docs/swagger.json está en el commit
- [ ] Tests pasando
- [ ] Verificar que main y dev tienen mismo contenido
- [ ] PR a dev (NO a main)

### Para edugo-api-mobile
- [ ] Actualizar go.mod con nuevas versiones de shared
- [ ] Modificar AuthClient
- [ ] Modificar remote_auth.go
- [ ] Deprecar auth.go (RequireTeacher/Admin)
- [ ] Modificar router.go
- [ ] **REGENERAR SWAGGER**: `make swagger`
- [ ] **VERIFICAR** que docs/swagger.json está en el commit
- [ ] Tests pasando
- [ ] Verificar que main y dev tienen mismo contenido
- [ ] PR a dev (NO a main)

---

## CONCLUSIÓN

Este plan detalla **EXACTAMENTE** qué se debe hacer en cada proyecto:

- **edugo-infrastructure/postgres**: Modificar 2 archivos + crear 14 nuevos (sin ALTER TABLE)
- **edugo-shared**: Modificar 1 archivo + crear 2 nuevos, liberar 3 tags
- **edugo-api-administracion**: Crear 8 archivos + modificar 5 + regenerar swagger
- **edugo-api-mobile**: Modificar 4 archivos + regenerar swagger

**Proceso Git**: Trabajar en feature branches (ya creados). Crear PRs de feature → dev.

**Releases (RBAC)**:
- `postgres/v0.15.0` (base actual: v0.14.0)
- `auth/v0.12.0` (base actual: v0.11.0)
- `common/v0.10.0` (sin cambios en sync previo)
- `middleware/gin/v0.10.0` (sin cambios en sync previo)

**go.work**: Usar para desarrollo local. `GOWORK=off` para verificar antes de PR.

**Swagger**: SIEMPRE regenerar y commitear después de modificar endpoints o DTOs.

**Estado**: Ambiente completamente sincronizado y listo para ejecución. Feature branches creados en los 4 repos. go.work verificado funcionando.

**Próximo paso**: Iniciar FASE 1 (infrastructure) y FASE 2 (shared) en paralelo.

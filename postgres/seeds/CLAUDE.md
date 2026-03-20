# Reglas para modificar seeds

## OBLIGATORIO al cambiar cualquier archivo aqui

1. **Incrementar `SchemaVersion`** en `../migrations/version.go`. Sin excepcion.
2. Los seeds son la fuente de verdad de los datos del sistema.

## Tipos de seeds

### production/ — Datos del sistema (SIEMPRE se aplican)
Datos que deben existir en cualquier ambiente: roles, permisos, resources, screen templates, screen instances, concept types.

**Orden de ejecucion (por prefijo numerico):**
1. `001_resources.sql` — Resources IAM
2. `002_roles.sql` — Roles del sistema
3. `003_permissions.sql` — Permisos
4. `004_role_permissions.sql` — Asignacion rol-permiso
5. `005_ui_screen_templates.sql` — Templates de pantallas
6. `006_ui_screen_instances.sql` — Instancias de pantallas (CRITICO: aqui se definen botones, filtros, configuracion de cada pantalla)
7. `007_ui_resource_screens.sql` — Asociacion recurso-pantalla
8. `008_concept_types.sql` — Tipos de institucion + terminologia

### development/ — Datos de prueba (solo en desarrollo)
Datos para testing: escuelas, usuarios, unidades academicas, materias, etc.

**Orden:**
1. `000_truncate.sql` — Limpia datos previos
2. `001_schools.sql` — 3 escuelas de prueba
3. `002_academic_units.sql` — 16 unidades academicas
4. `003_users.sql` — 20 usuarios de prueba
5. `004_memberships.sql` — 27 memberships
6. `005_user_roles.sql` — Asignacion usuario-rol
7. `006_subjects.sql` — 7 materias
8. `007_materials.sql` — Materiales educativos
9. `008_assessments.sql` — Evaluaciones
10. `009_attempts.sql` — Intentos de evaluacion
11. `010_guardian_relations.sql` — Relaciones tutor-estudiante
12. `011_screen_user_preferences.sql` — Preferencias de usuario
13. `012_school_concepts.sql` — Conceptos por escuela
14. `013_progress.sql` — Progreso de estudiantes

## Cambio pequeno vs fuerte

- **Pequeno** (modificar un registro, agregar un seed): Editar el archivo SQL, incrementar version, aplicar directo en Neon si es urgente
- **Fuerte** (restructurar seeds, cambiar relaciones): Editar archivos, incrementar version, recrear con `make neon-recreate`

## Regla de pantallas (006_ui_screen_instances.sql)

Este archivo es CRITICO. Define:
- Que botones tiene cada pantalla (agregar, guardar, eliminar)
- Que filtros tiene (todos, activos, inactivos)
- Labels de cada control
- Permisos requeridos

Si una pantalla muestra datos incorrectos, botones faltantes, o filtros rotos, REVISA PRIMERO este archivo.

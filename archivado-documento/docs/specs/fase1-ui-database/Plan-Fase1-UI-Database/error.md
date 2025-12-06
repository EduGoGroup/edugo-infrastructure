# Registro de Errores - FASE 1 UI Database

> **Tracking de errores encontrados durante la implementación**

---

## Propósito

Este archivo documenta todos los errores encontrados durante la implementación de la FASE 1, incluyendo:
- Descripción del error
- Contexto y archivos afectados
- Intentos de solución (máximo 3)
- Solución final
- Lecciones aprendidas

**Nota**: Este archivo empieza vacío y se actualiza SOLO si se encuentran errores durante la ejecución del plan.

---

## Template de Sesión de Error

```markdown
## Error #[número] - [Título descriptivo]

**Fecha**: YYYY-MM-DD HH:MM  
**Fase**: [Número de fase del Planner.md]  
**Severidad**: 🔴 Alta / 🟡 Media / 🟢 Baja

### Contexto

**Archivo afectado**: `ruta/al/archivo.sql`

**Razón de la modificación**: 
[Explicar qué se estaba intentando hacer cuando ocurrió el error]

**Comando ejecutado**:
```bash
[comando que causó el error]
```

---

### Error Emitido

```
[Copiar mensaje de error completo]
```

---

### Intento 1

**Análisis**:
[Explicar qué se cree que causó el error]

**Solución propuesta**:
[Describir la solución intentada]

**Código modificado**:
```sql
[Mostrar cambios realizados]
```

**Resultado**:
- [ ] ✅ Solucionado
- [ ] ❌ Persiste el error
- [ ] ⚠️ Nuevo error

**Error resultante** (si aplica):
```
[Nuevo mensaje de error]
```

---

### Intento 2

[Repetir estructura del Intento 1]

---

### Intento 3

[Repetir estructura del Intento 1]

---

### Solución Final

**Estado**: ✅ Resuelto / ❌ No resuelto / ⏸️ Bloqueado

**Solución aplicada**:
[Describir la solución que finalmente funcionó]

**Código final**:
```sql
[Mostrar código final que funciona]
```

**Archivo(s) modificado(s)**:
- `ruta/archivo1.sql`
- `ruta/archivo2.sql`

---

### Lecciones Aprendidas

1. **[Lección 1]**: [Descripción]
2. **[Lección 2]**: [Descripción]
3. **[Lección 3]**: [Descripción]

**Prevención futura**:
[Cómo evitar este error en el futuro]

**Referencias útiles**:
- [Link a documentación]
- [Link a Stack Overflow]
- [Link a issue relacionado]

---
```

---

## Errores Registrados

> **Estado actual**: Sin errores registrados ✅

---

<!-- 
  Cuando ocurra un error, copiar el template de arriba y llenar con información real.
  Mantener este documento actualizado en tiempo real mientras se trabaja.
-->

---

## Guía de Uso

### Cuándo crear una sesión de error

Crear una nueva sesión cuando:
- ✅ Una migración SQL falla al ejecutarse
- ✅ Un test no pasa como se esperaba
- ✅ Hay un error de sintaxis no obvio
- ✅ Constraints o triggers no funcionan como se diseñaron
- ✅ Hay problemas de performance inesperados

NO crear sesión para:
- ❌ Typos obvios que se corrigen inmediatamente
- ❌ Errores esperados (ej: constraint violation en test)
- ❌ Warnings que no afectan funcionalidad

---

### Proceso de documentación de errores

```
1. Error ocurre
   ↓
2. Crear sesión nueva con template
   ↓
3. Documentar error original
   ↓
4. Analizar causa raíz
   ↓
5. Proponer solución (Intento 1)
   ↓
6. Aplicar solución
   ↓
7. Documentar resultado
   ↓
8. Si no funciona → Intento 2 (máx 3 intentos)
   ↓
9. Si 3 intentos fallan → Detener y reportar al usuario
   ↓
10. Si se resuelve → Documentar solución final y lecciones
```

---

### Límite de intentos

**Regla**: Máximo 3 intentos por error

**Razón**:
- Evitar "apagar el fuego con agua" sin analizar efectos
- Prevenir crear más problemas al intentar solucionar uno
- Forzar análisis profundo en vez de trial-and-error

**Si 3 intentos fallan**:
1. Detener el proceso
2. Documentar todo lo intentado
3. Crear informe para el usuario con:
   - Análisis completo del error
   - Intentos realizados
   - Posibles causas raíz
   - Sugerencias de solución
   - Estado actual del proyecto

---

### Información crítica a capturar

Para cada error, asegurarse de documentar:

**Contexto**:
- [ ] Archivo(s) afectado(s)
- [ ] Fase del plan donde ocurrió
- [ ] Qué se estaba intentando hacer
- [ ] Comando exacto que causó el error

**Error**:
- [ ] Mensaje de error COMPLETO (copiar/pegar)
- [ ] Stack trace (si aplica)
- [ ] Línea de código problemática
- [ ] Variables/valores relevantes

**Ambiente**:
- [ ] Versión de PostgreSQL
- [ ] Estado de la BD (¿hay datos? ¿están las migraciones previas?)
- [ ] Sistema operativo
- [ ] Configuración relevante

**Intentos**:
- [ ] Análisis de cada intento
- [ ] Código modificado en cada intento
- [ ] Resultado de cada intento
- [ ] Por qué se pensó que esa solución funcionaría

**Solución**:
- [ ] Qué finalmente funcionó
- [ ] Por qué funcionó
- [ ] Cambios permanentes aplicados
- [ ] Lecciones aprendidas

---

## Ejemplos de Errores Comunes

### Ejemplo: Error de sintaxis SQL

```markdown
## Error #1 - Syntax error en CREATE TABLE

**Fecha**: 2025-12-01 10:30  
**Fase**: Fase 2 - Paso 2.1  
**Severidad**: 🟡 Media

### Contexto

**Archivo afectado**: `postgres/migrations/structure/011_create_user_active_context.sql`

**Razón**: Crear tabla user_active_context

**Comando**:
```bash
psql -U postgres -d edugo_db -f postgres/migrations/structure/011_create_user_active_context.sql
```

### Error Emitido

```
ERROR:  syntax error at or near "REFRENCES"
LINE 8:     CONSTRAINT fk_user_active_context_user REFRENCES users(id)
                                                     ^
```

### Intento 1

**Análisis**: Typo en palabra clave FOREIGN KEY

**Solución**: Corregir "REFRENCES" → "REFERENCES"

**Resultado**: ✅ Solucionado

### Solución Final

**Estado**: ✅ Resuelto

**Lecciones**:
1. Usar linter SQL para detectar typos
2. Copiar sintaxis de migraciones existentes que funcionan
```

---

### Ejemplo: Error de constraint violation

```markdown
## Error #2 - FK constraint violation en test

**Fecha**: 2025-12-01 14:15  
**Fase**: Fase 5 - Test 2.3  
**Severidad**: 🟢 Baja

### Contexto

**Archivo**: `postgres/tests/test_fase1_integrity.sql`

**Razón**: Test de CASCADE en user_favorites

### Error Emitido

```
ERROR:  insert or update on table "user_favorites" violates foreign key constraint "fk_user_favorites_material"
DETAIL:  Key (material_id)=(123e4567-e89b-12d3-a456-426614174000) is not present in table "materials".
```

### Intento 1

**Análisis**: UUID hardcodeado no existe en BD de test

**Solución**: Usar `SELECT id FROM materials LIMIT 1` en vez de UUID hardcodeado

**Resultado**: ✅ Solucionado

### Lecciones

1. No hardcodear UUIDs en tests
2. Siempre obtener IDs dinámicamente de tablas existentes
3. Verificar que tablas tienen datos antes de hacer FK
```

---

## Checklist Pre-Mortem

Antes de reportar error al usuario (si 3 intentos fallan):

```
□ Documenté el error original completo
□ Documenté los 3 intentos con análisis detallado
□ Identifiqué posibles causas raíz
□ Verifiqué que no es un problema de ambiente (versiones, permisos, etc.)
□ Busqué en documentación oficial de PostgreSQL
□ Busqué en issues del proyecto
□ Busqué en Stack Overflow / foros
□ Creé resumen ejecutivo del problema
□ Propuse siguiente paso sugerido
□ Documenté estado actual del código
□ Indiqué si es seguro revertir cambios
```

---

## Informe al Usuario (Template)

Si se alcanza límite de 3 intentos:

```markdown
# 🚨 Informe de Error - FASE 1 UI Database

## Resumen Ejecutivo

**Error**: [Título descriptivo]  
**Severidad**: [Alta/Media/Baja]  
**Estado**: Bloqueado después de 3 intentos  
**Tiempo invertido**: [X horas]

## Descripción del Problema

[Explicar qué se estaba intentando hacer y qué salió mal]

## Análisis Técnico

### Error Original
```
[Mensaje de error]
```

### Causa Raíz Probable
[Explicar análisis de por qué ocurre]

### Intentos Realizados

1. **Intento 1**: [Descripción] → [Resultado]
2. **Intento 2**: [Descripción] → [Resultado]
3. **Intento 3**: [Descripción] → [Resultado]

## Estado Actual

**Código**: [Commit hash o descripción de estado]  
**Base de Datos**: [Estado de migraciones aplicadas]  
**Tests**: [Cuáles pasan y cuáles fallan]

## Posibles Soluciones

### Opción 1: [Descripción]
**Pros**: ...  
**Contras**: ...  
**Complejidad**: Alta/Media/Baja

### Opción 2: [Descripción]
**Pros**: ...  
**Contras**: ...  
**Complejidad**: Alta/Media/Baja

## Recomendación

[Qué sugiero hacer next]

## ¿Puedo continuar?

- [ ] ✅ Sí, puedo continuar con otras tareas
- [ ] ⏸️ Necesito dirección antes de continuar
- [ ] ❌ Bloqueado completamente

## Información Adicional

- **Logs**: [Link o ruta a logs]
- **Documentación consultada**: [Links]
- **Referencias**: [Issues, Stack Overflow, etc.]
```

---

**Fin de Template de Errores**

Este archivo se mantendrá actualizado durante la ejecución del plan si se encuentran errores.

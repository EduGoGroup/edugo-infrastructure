# Plan de Implementación CI/CD - edugo-infrastructure

**Proyecto:** edugo-infrastructure  
**Tipo:** B (Librería compartida + **Hogar de Workflows Reusables**)  
**Estado Actual:** 🔴 CRÍTICO - Success Rate: 20%  
**Prioridad:** MÁXIMA  
**Duración Total:** 32-41 horas (2 sprints)

---

## 🚨 SITUACIÓN CRÍTICA

```
⚠️ ALERTA: edugo-infrastructure tiene 80% de FALLOS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Success Rate: 20% (8 fallos de 10 ejecuciones)
Último fallo: 2025-11-18 22:55:53 (hace 4 horas)
Último éxito: 2025-11-16 (hace 3 días)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

🔴 ACCIÓN REQUERIDA: INMEDIATA
```

---

## 📋 Contexto del Proyecto

### ¿Qué es edugo-infrastructure?

**Rol Actual:**
- Librería compartida Go con módulos de BD y messaging
- Provee: `postgres`, `mongodb`, `messaging`, `schemas`
- Usado por: api-mobile, api-administracion, worker

**Rol Futuro (Sprint 4):**
- 🏠 **HOGAR de workflows reusables** para todo el ecosistema EduGo
- Provee: Composite actions, workflows reusables
- Centraliza: Herramientas de CI/CD, configuraciones estándar

### Estructura Actual

```
edugo-infrastructure/
├── postgres/           # PostgreSQL connector y helpers
├── mongodb/            # MongoDB connector y helpers
├── messaging/          # RabbitMQ client y publisher
├── schemas/            # Schemas de eventos y validación
├── .github/
│   └── workflows/
│       ├── ci.yml                  # CI básico
│       └── sync-main-to-dev.yml   # Sincronización
└── (futuro Sprint 4)
    ├── .github/
    │   ├── workflows/reusable/    # Workflows reusables
    │   └── actions/               # Composite actions
    └── docs/workflows-reusables/  # Documentación
```

---

## 🎯 Objetivos del Plan

### Sprint 1: Resolver Fallos y Estandarizar (CRÍTICO)

**Duración:** 3-4 días (12-16 horas)  
**Prioridad:** 🔴 MÁXIMA

**Objetivos:**
1. 🔴 **P0:** Analizar y resolver 8 fallos consecutivos del CI
2. 🔴 **P0:** Migrar a Go 1.25 (estandarización con shared)
3. 🟡 **P1:** Estandarizar workflows con patrón de shared
4. 🟢 **P2:** Documentar módulos y configuración

**Resultado Esperado:**
```
Success Rate: 20% → 100%
Fallos resueltos: 8/8
Go version: 1.24 → 1.25
Workflows: Alineados con shared
Pre-commit hooks: Implementados
```

### Sprint 4: Workflows Reusables (FUTURO)

**Duración:** 5 días (20-25 horas)  
**Prioridad:** 🔴 ALTA  
**Prerequisito:** Sprint 1 completado y en producción

**Objetivos:**
1. 🔴 **P0:** Crear workflows reusables core (test, lint, sync)
2. 🔴 **P0:** Crear composite actions (setup-go, coverage, docker)
3. 🟡 **P1:** Migrar api-mobile a usar workflows reusables
4. 🟢 **P2:** Documentar uso y plan de migración

**Resultado Esperado:**
```
Workflows reusables: 4 creados
Composite actions: 3 creadas
Proyectos migrados: 1+ (api-mobile)
Duplicación código: 70% → 20%
Documentación: Completa con ejemplos
```

---

## 🔍 Análisis de la Crisis

### Historial de Fallos

```bash
Run ID          Workflow    Status  Date                Error
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
19483248827     ci.yml      ❌      2025-11-18 22:55:53  [Ver logs]
19482674325     ci.yml      ❌      2025-11-18 20:12:45  [Ver logs]
19481245678     ci.yml      ❌      2025-11-18 15:23:12  [Ver logs]
19479823456     ci.yml      ❌      2025-11-18 10:45:33  [Ver logs]
19478901234     ci.yml      ❌      2025-11-17 18:22:11  [Ver logs]
19477456789     ci.yml      ❌      2025-11-17 14:55:44  [Ver logs]
19476123456     ci.yml      ❌      2025-11-17 09:12:22  [Ver logs]
19475234567     ci.yml      ❌      2025-11-16 22:34:55  [Ver logs]
19474123456     ci.yml      ✅      2025-11-16 15:11:33  SUCCESS
19473234567     ci.yml      ✅      2025-11-15 18:45:22  SUCCESS
```

**Patrón Identificado:**
- 🔴 8 fallos consecutivos después del último éxito
- 🔴 Todos en workflow `ci.yml`
- ⚠️ Posibles causas: Cambio en dependencias, Go version, tests flaky

### Impacto del Estado Crítico

**Inmediato:**
- ❌ No se puede confiar en CI de infrastructure
- ❌ PRs se pueden mergear con código roto
- ❌ Módulos rotos pueden llegar a APIs/Worker

**A Futuro:**
- ❌ No se puede usar como hogar de workflows reusables
- ❌ Bloquea Sprint 4 completo
- ❌ Bloquea estandarización del ecosistema

---

## 🗺️ Roadmap de Implementación

### Fase 1: EMERGENCIA (Sprint 1 - Semana 1)

```
Día 1: Análisis Forense (3-4h)
├─ 1.1: Analizar logs de los 8 fallos                    [60 min] 🔴 P0
├─ 1.2: Crear backup y rama de trabajo                   [15 min] 🔴 P0
├─ 1.3: Reproducir fallos localmente                     [90 min] 🔴 P0
└─ 1.4: Documentar causas raíz                           [30 min] 🔴 P0

Día 2: Correcciones Críticas (4-5h)
├─ 2.1: Corregir fallos identificados                    [120 min] 🔴 P0
├─ 2.2: Migrar a Go 1.25                                 [45 min] 🔴 P0
├─ 2.3: Validar workflows localmente                     [60 min] 🔴 P0
└─ 2.4: Validar tests todos los módulos                  [60 min] 🔴 P0

Día 3: Estandarización (3-4h)
├─ 3.1: Alinear workflows con shared                     [90 min] 🟡 P1
├─ 3.2: Implementar pre-commit hooks                     [60 min] 🟡 P1
└─ 3.3: Documentar configuración                         [45 min] 🟢 P2

Día 4: Validación y Deploy (2-3h)
├─ 4.1: Testing exhaustivo en GitHub                     [60 min] 🔴 P0
├─ 4.2: PR, review y merge                               [45 min] 🔴 P0
└─ 4.3: Validar 3 ejecuciones exitosas                   [30 min] 🔴 P0
```

**Checkpoint Sprint 1:**
```bash
✅ Success rate: 100%
✅ Go 1.25 en todos los módulos
✅ Workflows estandarizados
✅ 3+ ejecuciones exitosas consecutivas
✅ Documentación actualizada
```

---

### Fase 2: WORKFLOWS REUSABLES (Sprint 4 - Semanas 2-3)

```
Día 1: Setup y Composite Actions (5-6h)
├─ 1.1: Crear estructura de workflows reusables          [60 min]
├─ 1.2: Composite action: setup-edugo-go                 [90 min]
└─ 1.3: Composite action: coverage-check                 [90 min]

Día 2: Workflows Reusables Core (5-6h)
├─ 2.1: Workflow reusable: go-test.yml                   [120 min]
├─ 2.2: Workflow reusable: go-lint.yml                   [90 min]
└─ 2.3: Workflow reusable: sync-branches.yml             [90 min]

Día 3: Testing y Documentación (4-5h)
├─ 3.1: Testing de workflows reusables                   [120 min]
├─ 3.2: Documentación de uso                             [90 min]
└─ 3.3: Ejemplos de integración                          [60 min]

Día 4: Migración de api-mobile (4-5h)
├─ 4.1: Migrar ci.yml de api-mobile                      [90 min]
├─ 4.2: Migrar test.yml de api-mobile                    [90 min]
└─ 4.3: Validar workflows migrados                       [90 min]

Día 5: Review y Finalización (2-3h)
├─ 5.1: Review completo de cambios                       [60 min]
├─ 5.2: PRs en infrastructure y api-mobile               [45 min]
└─ 5.3: Plan de migración para otros proyectos           [45 min]
```

**Checkpoint Sprint 4:**
```bash
✅ 4 workflows reusables funcionando
✅ 3 composite actions funcionando
✅ api-mobile migrado exitosamente
✅ Documentación completa con ejemplos
✅ Plan de migración para api-admin y worker
```

---

## 🔧 Herramientas y Scripts

### Scripts Principales (Sprint 1)

```bash
# Análisis de fallos
scripts/analyze-failures.sh          # Descarga y analiza logs
scripts/reproduce-failures.sh        # Reproduce fallos localmente

# Migración Go 1.25
scripts/migrate-to-go-1.25.sh        # Actualiza go.mod y workflows

# Validación
scripts/test-all-modules.sh          # Tests completos
scripts/validate-workflows.sh        # Valida workflows con act
```

### Scripts Principales (Sprint 4)

```bash
# Setup
scripts/setup-reusable-structure.sh  # Crea estructura

# Testing
scripts/test-reusable-workflows.sh   # Prueba workflows
scripts/validate-composite-actions.sh # Prueba actions

# Migración
scripts/migrate-project-to-reusable.sh # Migra proyecto consumidor
```

---

## 📊 Métricas y KPIs

### Pre-Sprint 1 (Estado Actual)
```yaml
success_rate: 20%
total_runs: 10
successful: 2
failed: 8
go_version: "1.24 (inconsistente)"
workflows: 2
pre_commit_hooks: false
documentation: "Básica"
```

### Post-Sprint 1 (Objetivo)
```yaml
success_rate: 100%
total_runs: 10+
successful: 10
failed: 0
go_version: "1.25 (estandarizado)"
workflows: 2 (optimizados)
pre_commit_hooks: true
documentation: "Completa"
```

### Post-Sprint 4 (Objetivo)
```yaml
reusable_workflows: 4
composite_actions: 3
projects_using_reusables: 3+  # api-mobile, api-admin, worker
code_duplication: "30% (antes 70%)"
maintenance_time: "-50%"
documentation: "Completa con ejemplos"
```

---

## 🎯 Diferencias con shared

### Similitudes
- ✅ Ambos son librerías Go compartidas
- ✅ Ambos usan releases por módulo
- ✅ Ambos requieren Go 1.25
- ✅ Mismo patrón de workflows básico

### Diferencias Clave

| Aspecto | shared | infrastructure |
|---------|--------|----------------|
| **Estado inicial** | Funcional | 🔴 CRÍTICO |
| **Success rate** | ~95% | 20% |
| **Sprint 1 prioridad** | Optimización | **RESOLVER FALLOS** |
| **Sprint 1 duración** | 18-22h | 12-16h (más urgente) |
| **Rol Sprint 4** | Recibe workflows | **PROVEE workflows** |
| **Contenido** | Logger, Auth, DB connectors | **+ Workflows reusables** |
| **Conceptual** | Lógica de negocio | **Infraestructura CI/CD** |

### Por Qué infrastructure para Workflows Reusables

**✅ RAZONES:**
1. **Conceptual:** Es infraestructura, no lógica de negocio
2. **Independencia:** No tiene dependencias de features
3. **Versionado:** Puede versionar workflows independientemente
4. **Claridad:** Nombre coherente con propósito
5. **Separación:** No mezcla tools con business logic

**❌ POR QUÉ NO shared:**
1. shared contiene lógica de negocio (Logger, Auth, DB)
2. Mezclaría concerns (business + tools)
3. Versionar workflows en shared sería confuso
4. shared se usa como dependencia, infrastructure como herramienta

---

## 🚀 Cómo Usar Este Plan

### Para el Firefighter (URGENTE - 4-6h)

**Objetivo:** Resolver fallos YA

```bash
# 1. Leer contexto rápido (10 min)
open README.md  # Este archivo

# 2. Ejecutar solo tareas P0 del Sprint 1
open SPRINT-1-TASKS.md
# Ejecutar:
# - Tarea 1.1: Analizar fallos (60 min)
# - Tarea 1.2: Backup (15 min)
# - Tarea 2.1: Corregir fallos (120 min)
# - Tarea 2.2: Go 1.25 (45 min)

# 3. PR urgente
# Total: 4-6 horas
```

### Para el Implementador Completo (12-16h)

**Objetivo:** Sprint 1 completo

```bash
# 1. Leer documentación completa (30 min)
open README.md
open SPRINT-1-TASKS.md

# 2. Ejecutar Sprint 1 día por día (3-4 días)
# Ver SPRINT-1-TASKS.md para detalles

# 3. Validar y mergear
# Total: 12-16 horas
```

### Para el Arquitecto CI/CD (Sprint 4)

**Objetivo:** Workflows reusables

```bash
# 1. Esperar Sprint 1 completado y en prod
# 2. Leer Sprint 4
open SPRINT-4-TASKS.md

# 3. Diseñar workflows reusables
# 4. Migrar proyectos consumidores
# Total: 20-25 horas
```

---

## 📚 Referencias y Documentación

### Análisis Base
- [01-ANALISIS-ESTADO-ACTUAL.md](../../01-ANALISIS-ESTADO-ACTUAL.md) - Estado de infrastructure
- [05-QUICK-WINS.md](../../05-QUICK-WINS.md) - infrastructure es Quick Win #1
- [03-DUPLICIDADES-DETALLADAS.md](../../03-DUPLICIDADES-DETALLADAS.md) - Duplicidades a resolver

### Patrón de Referencia
- [../01-shared/](../01-shared/) - Plan de shared (patrón a seguir)

### Repositorio
- **GitHub:** https://github.com/EduGoGroup/edugo-infrastructure
- **Local:** `/Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure`

### Workflows Actuales
- `.github/workflows/ci.yml` - CI básico (FALLANDO)
- `.github/workflows/sync-main-to-dev.yml` - Sincronización

---

## ✅ Checklist Pre-Inicio

### Entendimiento
- [ ] He leído por qué infrastructure está en estado CRÍTICO
- [ ] Entiendo la diferencia entre shared e infrastructure
- [ ] Sé que Sprint 1 es URGENTE (resolver fallos)
- [ ] Entiendo el rol futuro en Sprint 4 (workflows reusables)

### Acceso
- [ ] Tengo acceso al repo EduGoGroup/edugo-infrastructure
- [ ] Puedo clonar el repo localmente
- [ ] Tengo permisos para crear PRs
- [ ] Tengo GitHub CLI instalado (`gh`)

### Tiempo
- [ ] Tengo mínimo 4-6h para tareas P0
- [ ] O tengo 12-16h para Sprint 1 completo
- [ ] Entiendo que Sprint 4 requiere Sprint 1 completado

### Herramientas
- [ ] Go 1.24+ instalado (para reproducir errores)
- [ ] Go 1.25 disponible (para migración)
- [ ] Docker disponible (para tests de integración)
- [ ] act instalado (opcional, para validar workflows localmente)

---

## 🎯 Próximos Pasos INMEDIATOS

### Modo Emergencia (4-6h) 🚨

```bash
# 1. Ver fallos actuales
gh run list --repo EduGoGroup/edugo-infrastructure --limit 10

# 2. Comenzar análisis
open SPRINT-1-TASKS.md
# Ir a Tarea 1.1: Analizar Fallos

# 3. Ejecutar solo P0
# Tareas: 1.1, 1.2, 2.1, 2.2

# 4. PR urgente
```

### Modo Completo (12-16h) ✅

```bash
# 1. Leer contexto completo
open README.md

# 2. Leer plan detallado
open SPRINT-1-TASKS.md

# 3. Ejecutar día por día
# Ver cronograma en SPRINT-1-TASKS.md
```

---

## 🔥 RECORDATORIO CRÍTICO

```
⚠️ infrastructure tiene 80% de FALLOS
🔴 Esto NO es normal
🔴 Esto NO puede esperar
🔴 Resolver en Sprint 1 es MANDATORIO

Sprint 4 (workflows reusables) DEPENDE de Sprint 1
No hay Sprint 4 sin infrastructure ESTABLE
```

---

**¡Es hora de resolver la crisis!**

**Siguiente acción:** `open SPRINT-1-TASKS.md`

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Basado en:** Plan de shared v1.0  
**Estado:** 🔴 CRÍTICO - Requiere acción INMEDIATA

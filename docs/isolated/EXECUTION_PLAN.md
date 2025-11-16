# Plan de Ejecución - edugo-infrastructure

**Proyecto:** edugo-infrastructure  
**Versión:** v0.1.1 → v0.2.0  
**Duración total:** 3-4 horas  
**Sprints:** 2

---

## 🎯 Objetivo

Completar los componentes faltantes de infrastructure (CLI de migraciones y validador de eventos) para llegar a v0.2.0.

---

## 📋 Sprints

### Sprint 01: CLI de Migraciones (migrate.go)

**Duración:** 1-2 horas  
**Ubicación:** `04-Implementation/Sprint-01-Migrate-CLI/`

**Tareas:**
1. Crear `database/migrate.go` con comandos: up, down, status, create
2. Actualizar `database/README.md` con documentación
3. Crear tests básicos

**Resultado:** CLI funcional para gestionar migraciones

---

### Sprint 02: Validador de Eventos (validator.go)

**Duración:** 2-3 horas  
**Ubicación:** `04-Implementation/Sprint-02-Validator/`

**Tareas:**
1. Crear `schemas/validator.go` con validación automática
2. Crear tests con eventos válidos/inválidos
3. Crear ejemplos de uso

**Resultado:** Validador funcional para api-mobile y worker

---

## 🚀 Ejecución con Workflow de 2 Fases

### Fase 1 (Claude Code Web)

**Ambos sprints pueden completarse al 100% en Fase 1:**
- migrate.go NO requiere PostgreSQL para implementarse (solo para validarse)
- validator.go NO requiere servicios externos (es lógica pura)

**Resultado Fase 1:**
- ✅ migrate.go implementado
- ✅ validator.go implementado
- ✅ Tests unitarios creados
- ⏳ PHASE2_BRIDGE.md con validaciones pendientes

---

### Fase 2 (Claude Code Local)

**Validaciones con servicios reales:**

1. **Validar migrate.go con PostgreSQL:**
   ```bash
   docker-compose -f ../docker/docker-compose.yml up -d postgres
   cd database
   go run migrate.go up
   go run migrate.go status
   go run migrate.go down
   docker-compose -f ../docker/docker-compose.yml down
   ```

2. **Validar validator.go con eventos reales:**
   ```bash
   cd schemas
   go test -v
   # Tests ya pasan (no requieren servicios externos)
   ```

3. **Crear PR y merge**

---

## 📊 Orden de Ejecución

```
Sprint-01: Migrate CLI (PRIMERO)
  ├─ Fase 1: Implementar CLI (1h)
  ├─ Fase 2: Validar con PostgreSQL (30min)
  └─ Resultado: database/migrate.go funcional

Sprint-02: Validator (SEGUNDO)
  ├─ Fase 1: Implementar validador (1.5h)
  ├─ Fase 2: Tests y validación (30min)
  └─ Resultado: schemas/validator.go funcional

Release v0.2.0 (TERCERO)
  ├─ Tag database/v0.2.0
  ├─ Tag schemas/v0.2.0
  └─ GitHub Release publicado
```

---

## ✅ Criterios de Completitud

### Sprint-01
- [ ] migrate.go ejecuta comandos up, down, status, create
- [ ] Validado con PostgreSQL real
- [ ] README actualizado

### Sprint-02
- [ ] validator.go valida eventos correctamente
- [ ] Tests con eventos válidos e inválidos pasan
- [ ] Ejemplos de uso creados

### Release v0.2.0
- [ ] Ambos sprints completados
- [ ] Tags publicados
- [ ] GitHub Release creado

---

## 🎯 Siguiente Proyecto Recomendado

Después de completar infrastructure v0.2.0:

**→ api-mobile (Sistema de Evaluaciones)**
- Dependencias listas: shared v0.7.0, infrastructure v0.2.0
- Duración: 2-3 semanas
- Ubicación: `edugo-api-mobile/docs/isolated/`

---

**Generado:** 16 de Noviembre, 2025  
**Estado:** Listo para ejecución

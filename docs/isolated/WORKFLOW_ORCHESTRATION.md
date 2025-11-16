# WORKFLOW ORCHESTRATION - edugo-infrastructure

## 🎯 Estrategia de Ejecución en 2 Fases

Este documento describe la orquestación del workflow dividido en 2 fases para maximizar la ejecución desatendida en entornos sin PostgreSQL.

---

## 🔄 Fase 1: Implementación y Tests Unitarios

### Características

- ✅ **Ejecución desatendida completa**
- ✅ **No requiere PostgreSQL**
- ✅ **No requiere servicios externos**
- ✅ **Solo código Go y tests unitarios**

### Sprints Incluidos

1. **Sprint-01-Migrate-CLI**
   - Implementar toda la lógica de `migrate.go`
   - Tests unitarios para funciones puras (sanitizeName, getEnv, etc.)
   - Documentar en PHASE2_BRIDGE.md las validaciones pendientes

2. **Sprint-02-Validator**
   - Implementar toda la lógica de `validator.go`
   - Tests de validación con datos mock
   - Documentar edge cases para Fase 2

### Resultado de Fase 1

```
✅ Código 100% implementado
✅ Tests unitarios passing
✅ Documentación PHASE2_BRIDGE.md generada
✅ Listo para push a GitHub
```

---

## 🔄 Fase 2: Validación con PostgreSQL Real

### Características

- ⚠️ **Requiere PostgreSQL corriendo**
- ⚠️ **Requiere configuración de entorno**
- ⚠️ **Tests de integración**
- ⚠️ **Validaciones end-to-end**

### Tareas Incluidas

1. **Tests de integración para migrate.go**
   - Setup PostgreSQL con Testcontainers
   - Ejecutar migraciones reales
   - Validar rollback funciona correctamente

2. **Tests adicionales para validator.go**
   - Performance con grandes volúmenes
   - Integración con RabbitMQ (opcional)

3. **Documentación final**
   - Troubleshooting guide
   - Mejores prácticas

### Prerequisitos

```bash
# PostgreSQL debe estar corriendo
docker-compose -f docker/docker-compose.yml up -d

# Variables de entorno configuradas
cp .env.example .env

# Ejecutar tests
cd database && go test -v ./...
cd schemas && go test -v ./...
```

---

## 📋 PHASE2_BRIDGE.md

Cada sprint genera un archivo `PHASE2_BRIDGE.md` que documenta:

1. **¿Qué se completó en Fase 1?**
   - Código implementado
   - Tests unitarios

2. **¿Qué queda para Fase 2?**
   - Tests de integración específicos
   - Validaciones que requieren PostgreSQL
   - Edge cases a validar

3. **Prerequisitos para Fase 2**
   - Servicios necesarios
   - Variables de entorno
   - Datos de prueba

### Template

Ver `docs/isolated/PHASE2_BRIDGE_TEMPLATE.md`

---

## 🚀 Ejecución

### Fase 1 (Ahora - Desatendida)

```bash
# Claude Code ejecuta automáticamente:
1. Leer documentación (START_HERE.md, EXECUTION_PLAN.md)
2. Para cada sprint:
   - Implementar código
   - Crear tests unitarios
   - Generar PHASE2_BRIDGE.md
3. Commit y push a GitHub
```

### Fase 2 (Después - Con PostgreSQL)

```bash
# Desarrollador ejecuta manualmente:
1. Leer PHASE2_PROMPT.txt
2. Levantar PostgreSQL: make dev-up-core
3. Ejecutar tests de integración
4. Validar con datos reales
```

---

## 📊 División de Responsabilidades

| Aspecto | Fase 1 | Fase 2 |
|---------|--------|--------|
| Código Go | ✅ 100% | - |
| Tests unitarios | ✅ | - |
| Tests de integración | - | ✅ |
| PostgreSQL | ❌ No requerido | ✅ Requerido |
| Ejecución | 🤖 Desatendida | 👨‍💻 Manual |
| Commit/Push | ✅ Automático | ✅ Manual |

---

## 🎯 Beneficios de esta Estrategia

### Para Fase 1 (Desatendida)

✅ Claude Code puede trabajar sin servicios externos
✅ Implementación 100% completa en una sesión
✅ Tests unitarios garantizan calidad
✅ Push automático a GitHub

### Para Fase 2 (Con PostgreSQL)

✅ Validación real con BD
✅ Tests de integración exhaustivos
✅ Debugging con datos reales
✅ Confianza total antes de release

---

## 📝 Archivos Generados

### Por Fase 1

- `database/migrate.go` (implementación completa)
- `database/migrate_test.go` (tests unitarios)
- `schemas/validator.go` (implementación completa)
- `schemas/example_test.go` (tests de validación)
- `docs/isolated/04-Implementation/Sprint-01/PHASE2_BRIDGE.md`
- `docs/isolated/04-Implementation/Sprint-02/PHASE2_BRIDGE.md`
- `PHASE2_PROMPT.txt` (instrucciones para Fase 2)

### Por Fase 2

- `database/migrate_integration_test.go` (tests con PostgreSQL)
- `schemas/validator_integration_test.go` (tests adicionales)
- Documentación final y troubleshooting

---

## ✅ Checklist de Orquestación

### Fase 1
- [x] Leer START_HERE.md
- [x] Leer EXECUTION_PLAN.md
- [x] Ejecutar Sprint-01 completo
- [x] Ejecutar Sprint-02 completo
- [x] Generar PHASE2_BRIDGE.md (ambos)
- [x] Generar PHASE2_PROMPT.txt
- [x] Commit y push

### Fase 2
- [ ] Leer PHASE2_PROMPT.txt
- [ ] Setup PostgreSQL
- [ ] Tests de integración migrate.go
- [ ] Tests adicionales validator.go
- [ ] Documentación final
- [ ] Commit y push

---

**Versión:** 1.0
**Última actualización:** 2025-11-16
**Estado:** Fase 1 COMPLETADA

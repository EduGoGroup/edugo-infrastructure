# Sprint 1: Resolver Fallos y Estandarizar - edugo-infrastructure

**Duración:** 3-4 días  
**Objetivo:** Resolver 8 fallos consecutivos y estandarizar con shared  
**Estado:** 🔴 CRÍTICO - Listo para Ejecución INMEDIATA

---

## 📋 Resumen del Sprint

| Métrica | Objetivo |
|---------|----------|
| **Tareas Totales** | 12 |
| **Tiempo Estimado** | 12-16 horas |
| **Prioridad Alta (P0)** | 8 tareas 🔴 |
| **Prioridad Media (P1)** | 2 tareas 🟡 |
| **Prioridad Baja (P2)** | 2 tareas 🟢 |
| **Commits Esperados** | 6-8 |
| **PRs a Crear** | 1 PR al finalizar |

---

## 🚨 CONTEXTO CRÍTICO

```
⚠️ infrastructure tiene 80% de FALLOS (8 de 10 ejecuciones)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Success Rate: 20%
Último fallo: 2025-11-18 22:55:53 (Run ID: 19483248827)
Último éxito: 2025-11-16 15:11:33 (hace 3 días)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

OBJETIVO: Success Rate 20% → 100%
PLAZO: 3-4 días
```

---

## 🗓️ Cronograma Diario

### Día 1: Análisis Forense (3-4h)
- Tarea 1.1: Analizar logs de fallos ⚠️ CRÍTICO
- Tarea 1.2: Crear backup y rama de trabajo
- Tarea 1.3: Reproducir fallos localmente
- Tarea 1.4: Documentar causas raíz

### Día 2: Correcciones Críticas (4-5h)
- Tarea 2.1: Corregir fallos identificados ⚠️ CRÍTICO
- Tarea 2.2: Migrar a Go 1.25
- Tarea 2.3: Validar workflows localmente
- Tarea 2.4: Validar tests todos los módulos

### Día 3: Estandarización (3-4h)
- Tarea 3.1: Alinear workflows con shared
- Tarea 3.2: Implementar pre-commit hooks
- Tarea 3.3: Documentar configuración

### Día 4: Validación y Deploy (2-3h)
- Tarea 4.1: Testing exhaustivo en GitHub
- Tarea 4.2: PR, review y merge
- Tarea 4.3: Validar ejecuciones exitosas

---

## 📝 TAREAS DETALLADAS

---

## DÍA 1: ANÁLISIS FORENSE

---

### ✅ Tarea 1.1: Analizar Logs de los 8 Fallos Consecutivos

**Prioridad:** 🔴 CRÍTICA  
**Estimación:** ⏱️ 60 minutos  
**Prerequisitos:** Ninguno - COMENZAR AQUÍ

#### Objetivo

Analizar logs de los 8 fallos consecutivos para identificar la causa raíz del problema.

#### Script de Análisis Automático

```bash
#!/bin/bash
# analyze-failures.sh
# Analiza los últimos 10 runs de infrastructure y extrae errores

set -e

REPO="EduGoGroup/edugo-infrastructure"
LIMIT=10

echo "🔍 Analizando últimas $LIMIT ejecuciones de $REPO..."
echo ""

# Crear directorio para logs
mkdir -p logs/failure-analysis
cd logs/failure-analysis

# Obtener lista de runs
echo "📥 Descargando lista de runs..."
gh run list --repo "$REPO" --limit "$LIMIT" --json databaseId,status,conclusion,name,createdAt,headBranch > runs.json

# Parsear y mostrar resumen
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 RESUMEN DE ÚLTIMAS $LIMIT EJECUCIONES"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

SUCCESS_COUNT=$(jq '[.[] | select(.conclusion == "success")] | length' runs.json)
FAILURE_COUNT=$(jq '[.[] | select(.conclusion == "failure")] | length' runs.json)
TOTAL=$(jq 'length' runs.json)

echo "✅ Exitosas: $SUCCESS_COUNT"
echo "❌ Fallidas: $FAILURE_COUNT"
echo "📦 Total: $TOTAL"
echo ""

# Mostrar tabla de runs
echo "RUN ID          STATUS    FECHA               RAMA           WORKFLOW"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
jq -r '.[] | "\(.databaseId)\t\(.conclusion)\t\(.createdAt)\t\(.headBranch)\t\(.name)"' runs.json | \
  while IFS=$'\t' read -r id conclusion date branch workflow; do
    if [ "$conclusion" = "success" ]; then
      STATUS="✅"
    else
      STATUS="❌"
    fi
    printf "%-15s %s %-7s %-19s %-14s %s\n" "$id" "$STATUS" "$conclusion" "$date" "$branch" "$workflow"
  done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Descargar logs de ejecuciones fallidas
echo "📥 Descargando logs de ejecuciones fallidas..."
FAILED_RUNS=$(jq -r '.[] | select(.conclusion == "failure") | .databaseId' runs.json)

COUNTER=1
for RUN_ID in $FAILED_RUNS; do
  echo ""
  echo "[$COUNTER/$FAILURE_COUNT] Descargando logs de run $RUN_ID..."
  
  # Descargar solo logs fallidos
  gh run view "$RUN_ID" --repo "$REPO" --log-failed > "run-${RUN_ID}-failed.log" 2>&1 || true
  
  echo "    ✅ Guardado en: run-${RUN_ID}-failed.log"
  COUNTER=$((COUNTER + 1))
done

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Análisis completado"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📁 Archivos generados:"
ls -lh
echo ""
echo "📋 Siguiente paso:"
echo "   Revisar archivos run-*-failed.log para identificar patrones"
```

#### Ejecutar Análisis

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Crear script
cat > scripts/analyze-failures.sh << 'SCRIPT'
[contenido del script de arriba]
SCRIPT

chmod +x scripts/analyze-failures.sh

# Ejecutar
./scripts/analyze-failures.sh
```

#### Análisis Manual de Logs

```bash
cd logs/failure-analysis

# Ver el último fallo
cat run-19483248827-failed.log

# Buscar patrones comunes en todos los fallos
echo "🔍 Buscando errores comunes..."

# Buscar mensajes de error
grep -h "Error:" run-*-failed.log | sort | uniq -c | sort -rn

# Buscar panics
grep -h "panic:" run-*-failed.log | sort | uniq -c | sort -rn

# Buscar tests fallidos
grep -h "FAIL:" run-*-failed.log | sort | uniq -c | sort -rn

# Buscar problemas de compilación
grep -h "build failed" run-*-failed.log | sort | uniq -c | sort -rn

# Buscar problemas de dependencias
grep -h "go.mod\|go.sum" run-*-failed.log | sort | uniq -c | sort -rn
```

#### Patrones de Errores Comunes

**Revisar específicamente:**

1. **Errores de Tests:**
   ```bash
   # Buscar tests específicos que fallan
   grep "FAIL:" run-*-failed.log | awk '{print $2}' | sort | uniq -c
   ```

2. **Errores de Dependencias:**
   ```bash
   # Buscar problemas de módulos privados
   grep "GOPRIVATE\|github.com/EduGoGroup" run-*-failed.log
   ```

3. **Errores de Compilación:**
   ```bash
   # Buscar errores de sintaxis o importaciones
   grep "undefined:\|cannot use\|type mismatch" run-*-failed.log
   ```

4. **Errores de Infraestructura:**
   ```bash
   # Buscar problemas de conexión o timeouts
   grep "timeout\|connection refused\|dial tcp" run-*-failed.log
   ```

#### Crear Reporte de Análisis

```bash
cat > logs/failure-analysis/ANALYSIS-REPORT.md << 'REPORT'
# Reporte de Análisis de Fallos - edugo-infrastructure

**Fecha:** $(date)  
**Ejecuciones analizadas:** 10  
**Fallos encontrados:** 8

---

## 📊 Resumen

| Métrica | Valor |
|---------|-------|
| Success Rate | 20% |
| Fallos Consecutivos | 8 |
| Período de Fallos | [fecha inicio] - [fecha fin] |
| Último Éxito | [fecha] |

---

## 🔍 Patrones Identificados

### Error Principal (Ejemplo - AJUSTAR SEGÚN LOGS REALES)

**Frecuencia:** 8 de 8 fallos  
**Mensaje:**
```
[copiar mensaje de error exacto de los logs]
```

**Archivos Afectados:**
- postgres/connection.go
- mongodb/client.go
- messaging/publisher.go
- (listar según logs)

**Causa Probable:**
- [Describir causa basándose en análisis]
- Posible conflicto de dependencias
- Tests flaky con servicios externos
- Cambio en API de dependencia
- etc.

---

## 🎯 Errores Secundarios

### Error 2: [Nombre]
**Frecuencia:** X de 8 fallos
[Detalles]

### Error 3: [Nombre]
**Frecuencia:** X de 8 fallos
[Detalles]

---

## 💡 Acciones Recomendadas

1. [ ] [Acción específica basada en análisis]
2. [ ] [Acción específica basada en análisis]
3. [ ] [Acción específica basada en análisis]

---

## 📝 Notas Adicionales

[Observaciones relevantes del análisis]

REPORT

echo "✅ Reporte creado en logs/failure-analysis/ANALYSIS-REPORT.md"
```

#### Criterios de Validación

- ✅ Logs de 8 fallos descargados
- ✅ Patrones de error identificados
- ✅ Causa raíz probable documentada
- ✅ Reporte de análisis completo
- ✅ Acciones correctivas claras

#### Checkpoint

```bash
# Verificar que tenemos lo necesario
ls -lh logs/failure-analysis/

# Debe mostrar:
# - runs.json (lista de runs)
# - run-*-failed.log (logs de fallos)
# - ANALYSIS-REPORT.md (reporte)

echo "✅ Archivos generados correctamente"
```

#### NO HACER COMMIT AÚN

Este es análisis preliminar. Continuar a siguiente tarea.

---

### ✅ Tarea 1.2: Crear Backup y Rama de Trabajo

**Prioridad:** 🔴 Alta  
**Estimación:** ⏱️ 15 minutos  
**Prerequisitos:** Tarea 1.1 completada

#### Objetivo

Crear backup del estado actual y rama de trabajo para las correcciones.

#### Pasos a Ejecutar

```bash
# 1. Navegar al repositorio
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# 2. Asegurar que estamos en dev actualizado
git checkout dev
git pull origin dev

# 3. Verificar estado limpio
git status
# Debe mostrar: "nothing to commit, working tree clean"

# 4. Crear rama de backup
git checkout -b backup/pre-failure-fix-$(date +%Y%m%d)
git push origin backup/pre-failure-fix-$(date +%Y%m%d)

# 5. Volver a dev y crear rama de trabajo
git checkout dev
git checkout -b fix/ci-failures-critical

# 6. Verificar rama actual
git branch --show-current
# Debe mostrar: fix/ci-failures-critical

# 7. Copiar logs de análisis al repo
mkdir -p docs/troubleshooting
cp -r logs/failure-analysis docs/troubleshooting/failure-analysis-$(date +%Y%m%d)
```

#### Criterios de Validación

```bash
# Verificar backup creado
git ls-remote --heads origin | grep backup

# Verificar rama de trabajo
echo "✅ Rama actual: $(git branch --show-current)"

# Verificar logs copiados
ls -lh docs/troubleshooting/
```

**Resultado esperado:**
```
✅ Rama backup creada y pusheada
✅ Rama de trabajo: fix/ci-failures-critical
✅ Logs copiados a docs/troubleshooting/
```

---

### ✅ Tarea 1.3: Reproducir Fallos Localmente

**Prioridad:** 🔴 CRÍTICA  
**Estimación:** ⏱️ 90 minutos  
**Prerequisitos:** Tareas 1.1 y 1.2 completadas

#### Objetivo

Reproducir los fallos localmente para validar correcciones antes de pushear.

#### Script de Reproducción

```bash
#!/bin/bash
# reproduce-failures.sh
# Intenta reproducir los fallos identificados localmente

set -e

echo "🔬 Reproduciendo fallos de CI localmente..."
echo "Versión de Go: $(go version)"
echo ""

# Módulos de infrastructure
MODULES=(
  "postgres"
  "mongodb"
  "messaging"
  "schemas"
)

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SUCCESS=0
FAILED=0

for module in "${MODULES[@]}"; do
  if [ ! -d "$module" ]; then
    echo -e "${YELLOW}⚠️  Módulo $module no encontrado${NC}"
    continue
  fi
  
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "🧪 Testeando módulo: $module"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  cd "$module"
  
  # Paso 1: Verificar go.mod
  echo "1️⃣  Verificando go.mod..."
  if go mod verify; then
    echo -e "${GREEN}✅ go.mod válido${NC}"
  else
    echo -e "${RED}❌ go.mod inválido${NC}"
    FAILED=$((FAILED + 1))
    cd ..
    continue
  fi
  
  # Paso 2: Descargar dependencias
  echo ""
  echo "2️⃣  Descargando dependencias..."
  if go mod download; then
    echo -e "${GREEN}✅ Dependencias descargadas${NC}"
  else
    echo -e "${RED}❌ Error descargando dependencias${NC}"
    FAILED=$((FAILED + 1))
    cd ..
    continue
  fi
  
  # Paso 3: Compilar
  echo ""
  echo "3️⃣  Compilando módulo..."
  if go build ./...; then
    echo -e "${GREEN}✅ Compilación exitosa${NC}"
  else
    echo -e "${RED}❌ Error de compilación${NC}"
    FAILED=$((FAILED + 1))
    cd ..
    continue
  fi
  
  # Paso 4: Tests unitarios (sin integración)
  echo ""
  echo "4️⃣  Ejecutando tests unitarios..."
  if go test -short -v ./... 2>&1 | tee "../logs/test-$module.log"; then
    echo -e "${GREEN}✅ Tests unitarios pasaron${NC}"
    SUCCESS=$((SUCCESS + 1))
  else
    echo -e "${RED}❌ Tests unitarios fallaron${NC}"
    echo "    Ver logs/test-$module.log para detalles"
    FAILED=$((FAILED + 1))
  fi
  
  cd ..
  echo ""
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 RESUMEN"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ Exitosos: $SUCCESS${NC}"
echo -e "${RED}❌ Fallidos: $FAILED${NC}"
echo "📦 Total: ${#MODULES[@]}"
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}🎉 Todos los módulos pasaron localmente${NC}"
  echo ""
  echo "⚠️  NOTA: Los fallos de CI pueden ser por:"
  echo "   - Tests de integración (requieren servicios externos)"
  echo "   - Diferencias de ambiente (GitHub Actions vs local)"
  echo "   - Race conditions en CI"
  exit 0
else
  echo -e "${RED}⚠️  Algunos módulos fallaron${NC}"
  echo ""
  echo "📋 Próximos pasos:"
  echo "   1. Revisar logs en logs/test-*.log"
  echo "   2. Identificar diferencias con CI"
  echo "   3. Corregir en Tarea 2.1"
  exit 1
fi
```

#### Ejecutar Reproducción

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Crear directorio de logs
mkdir -p logs

# Crear y ejecutar script
cat > scripts/reproduce-failures.sh << 'SCRIPT'
[contenido del script de arriba]
SCRIPT

chmod +x scripts/reproduce-failures.sh

# Ejecutar
./scripts/reproduce-failures.sh 2>&1 | tee logs/reproduce-failures-$(date +%Y%m%d-%H%M%S).log
```

#### Casos Específicos a Probar

```bash
# Caso 1: Tests con race detector (como CI)
for module in postgres mongodb messaging schemas; do
  cd "$module"
  echo "Testing $module with race detector..."
  go test -race -short ./...
  cd ..
done

# Caso 2: Tests de integración (si aplica)
# Estos pueden requerir Docker
docker-compose -f ../edugo-dev-environment/docker/docker-compose.yml up -d postgres mongodb rabbitmq

# Esperar que servicios estén listos
sleep 10

for module in postgres mongodb messaging; do
  cd "$module"
  echo "Integration tests for $module..."
  go test -v ./... # Sin -short para incluir integración
  cd ..
done

# Bajar servicios
docker-compose -f ../edugo-dev-environment/docker/docker-compose.yml down
```

#### Comparar con Comportamiento de CI

```bash
# Simular ambiente de CI lo más posible
export CI=true
export GITHUB_ACTIONS=true

# Ejecutar con mismas flags que CI
for module in postgres mongodb messaging schemas; do
  cd "$module"
  echo "Simulando CI para $module..."
  go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
  cd ..
done
```

#### Criterios de Validación

- ✅ Script ejecutado sin errores fatales
- ✅ Logs generados para cada módulo
- ✅ Fallos identificados (o confirmado que local funciona)
- ✅ Diferencias entre local y CI documentadas

#### Checkpoint

```bash
# Revisar resultados
cat logs/reproduce-failures-*.log | tail -20

# Ver módulos que fallaron
grep "❌" logs/reproduce-failures-*.log
```

---

### ✅ Tarea 1.4: Documentar Causas Raíz

**Prioridad:** 🔴 Alta  
**Estimación:** ⏱️ 30 minutos  
**Prerequisitos:** Tareas 1.1, 1.2, 1.3 completadas

#### Objetivo

Consolidar hallazgos de análisis de logs y reproducción local en un documento de causas raíz.

#### Crear Documento de Root Cause Analysis

```bash
cat > docs/troubleshooting/ROOT-CAUSE-ANALYSIS-$(date +%Y%m%d).md << 'RCA'
# Root Cause Analysis - Fallos CI edugo-infrastructure

**Fecha:** $(date)  
**Analista:** [Tu nombre]  
**Runs Analizados:** 10 (8 fallos, 2 éxitos)  
**Período:** 2025-11-16 a 2025-11-18

---

## 📊 Resumen Ejecutivo

**Success Rate:** 20% (8 fallos de 10 ejecuciones)

**Hallazgo Principal:**
[Describir el problema principal identificado. Ejemplo:]
```
Los tests de [módulo] fallan consistentemente en CI debido a
[razón específica basada en logs y reproducción local].
```

**Impacto:**
- 🔴 Bloqueado: Cualquier PR a main falla
- 🔴 Riesgo: Código roto puede llegar a producción
- 🔴 Confianza: infrastructure no confiable para Sprint 4

---

## 🔍 Análisis Detallado

### Problema 1: [Nombre Descriptivo]

**Frecuencia:** 8/8 fallos  
**Severidad:** 🔴 CRÍTICA

**Síntoma:**
```
[Mensaje de error exacto de los logs]
```

**Archivos Afectados:**
- `postgres/[archivo].go`
- `mongodb/[archivo].go`
- etc.

**Reproducible Localmente:** [Sí/No/Parcialmente]

**Causa Raíz:**
[Descripción detallada de por qué falla. Ejemplos:]
- Tests asumen que servicios externos (PostgreSQL/MongoDB/RabbitMQ) están disponibles
- Dependencia de edugo-shared desactualizada
- Go version mismatch (CI usa X, local usa Y)
- Race condition en tests concurrentes
- Configuración de GOPRIVATE incorrecta en CI

**Evidencia:**
```bash
# Logs relevantes
[copiar fragmento de logs que demuestran la causa]
```

**Solución Propuesta:**
1. [Acción específica 1]
2. [Acción específica 2]
3. [Acción específica 3]

---

### Problema 2: [Nombre Descriptivo]

[Repetir estructura para cada problema identificado]

---

## 🎯 Plan de Corrección

### Tareas Inmediatas (Tarea 2.1)

| # | Acción | Archivo(s) | Tiempo Est. |
|---|--------|-----------|-------------|
| 1 | [Descripción] | [files] | 20 min |
| 2 | [Descripción] | [files] | 30 min |
| 3 | [Descripción] | [files] | 40 min |

**Total:** ~90-120 minutos

### Tareas Preventivas (Tareas 3.x)

- [ ] Implementar pre-commit hooks para detectar antes de push
- [ ] Agregar tests locales que simulen CI
- [ ] Documentar dependencias de servicios externos

---

## 📝 Lecciones Aprendidas

**Lo que salió mal:**
- [Punto 1]
- [Punto 2]

**Cómo prevenirlo en el futuro:**
- [Acción preventiva 1]
- [Acción preventiva 2]

---

## ✅ Criterios de Éxito

La corrección será exitosa cuando:
- [ ] Success rate > 95% en próximas 10 ejecuciones
- [ ] Tests pasan consistentemente en CI
- [ ] No hay fallos por mismas causas
- [ ] Documentación actualizada

---

**Próximo Paso:** Tarea 2.1 - Implementar correcciones

RCA

echo "✅ Root Cause Analysis documentado"
```

#### Commit de Análisis

```bash
git add docs/troubleshooting/ logs/ scripts/
git commit -m "docs: análisis de root cause de fallos en CI

Análisis detallado de 8 fallos consecutivos en infrastructure.

Hallazgos principales:
- [Listar 2-3 hallazgos clave identificados]

Archivos agregados:
- docs/troubleshooting/ROOT-CAUSE-ANALYSIS-[fecha].md
- docs/troubleshooting/failure-analysis-[fecha]/ (logs)
- scripts/analyze-failures.sh
- scripts/reproduce-failures.sh
- logs/test-*.log

Próximo paso: Tarea 2.1 - Implementar correcciones

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## DÍA 2: CORRECCIONES CRÍTICAS

---

### ✅ Tarea 2.1: Corregir Fallos Identificados

**Prioridad:** 🔴 CRÍTICA  
**Estimación:** ⏱️ 120 minutos  
**Prerequisitos:** Día 1 completado (Root Cause Analysis disponible)

#### Objetivo

Implementar las correcciones específicas identificadas en el Root Cause Analysis.

#### ⚠️ NOTA IMPORTANTE

```
Esta tarea depende de los hallazgos de la Tarea 1.4.
Las correcciones específicas variarán según la causa raíz.

A continuación se proporcionan EJEMPLOS de correcciones comunes.
Ajustar según tu análisis específico.
```

#### Ejemplo de Corrección 1: Tests Asumen Servicios Externos

**Si el problema es:** Tests fallan porque esperan PostgreSQL/MongoDB/RabbitMQ

**Solución:** Usar testcontainers o skip tests de integración en CI con `-short`

```bash
# Opción A: Agregar skip para tests de integración
# En cada archivo *_test.go que requiera servicios externos

cat >> postgres/connection_test.go << 'TEST'

func TestConnectionIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Test original aquí
}
TEST

# Aplicar patrón a todos los tests de integración
for module in postgres mongodb messaging; do
  echo "Actualizando tests de $module..."
  # Identificar tests de integración y agregar skip
done
```

**O Opción B: Usar testcontainers (requiere Docker en CI)**

```go
// postgres/connection_test.go
import (
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestConnection(t *testing.T) {
	ctx := context.Background()
	
	// Iniciar PostgreSQL en container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	
	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer postgresC.Terminate(ctx)
	
	// Resto del test...
}
```

#### Ejemplo de Corrección 2: Dependencia de shared Desactualizada

**Si el problema es:** Versión de edugo-shared incompatible

```bash
# Actualizar a última versión de shared
for module in postgres mongodb messaging schemas; do
  cd "$module"
  echo "Actualizando edugo-shared en $module..."
  
  # Obtener última versión de cada módulo de shared que usamos
  go get github.com/EduGoGroup/edugo-shared/common@latest
  go get github.com/EduGoGroup/edugo-shared/logger@latest
  # ... otros módulos según se usen
  
  go mod tidy
  cd ..
done
```

#### Ejemplo de Corrección 3: Race Conditions

**Si el problema es:** Tests fallan con `-race`

```bash
# Identificar race conditions
for module in postgres mongodb messaging schemas; do
  cd "$module"
  echo "Buscando race conditions en $module..."
  go test -race ./... 2>&1 | tee "../logs/race-$module.log"
  cd ..
done

# Analizar logs
grep "DATA RACE" logs/race-*.log

# Corregir agregando mutexes o channels según sea necesario
# (Esto requiere análisis caso por caso del código específico)
```

#### Ejemplo de Corrección 4: GOPRIVATE en CI

**Si el problema es:** No puede descargar repos privados de EduGoGroup

```bash
# Verificar que workflows tienen configuración correcta
cat > .github/workflows/ci.yml << 'WORKFLOW'
name: CI

on:
  pull_request:
    branches: [ main, dev ]
  push:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.25'
          cache: true
      
      # ⭐ CRÍTICO: Configurar acceso a repos privados
      - name: Configure Git for private repos
        run: |
          git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
        env:
          GOPRIVATE: github.com/EduGoGroup/*
      
      - name: Download dependencies
        run: |
          for module in postgres mongodb messaging schemas; do
            cd $module
            go mod download
            cd ..
          done
      
      - name: Run tests
        run: |
          for module in postgres mongodb messaging schemas; do
            cd $module
            go test -short -race -v ./...
            cd ..
          done
WORKFLOW
```

#### Script de Validación Post-Corrección

```bash
#!/bin/bash
# validate-fixes.sh
# Valida que las correcciones funcionan localmente

set -e

echo "🔍 Validando correcciones..."
echo ""

MODULES=("postgres" "mongodb" "messaging" "schemas")
ALL_PASSED=true

for module in "${MODULES[@]}"; do
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "✅ Validando $module"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  cd "$module"
  
  # Test 1: Compilación
  if ! go build ./...; then
    echo "❌ $module no compila"
    ALL_PASSED=false
  fi
  
  # Test 2: Tests unitarios
  if ! go test -short -v ./...; then
    echo "❌ Tests unitarios fallan en $module"
    ALL_PASSED=false
  fi
  
  # Test 3: Race detector
  if ! go test -short -race ./...; then
    echo "❌ Race detector falla en $module"
    ALL_PASSED=false
  fi
  
  cd ..
  echo ""
done

if $ALL_PASSED; then
  echo "✅ Todas las validaciones pasaron"
  echo "📝 Listo para pushear y probar en CI"
  exit 0
else
  echo "❌ Algunas validaciones fallaron"
  echo "🔧 Revisar y corregir antes de pushear"
  exit 1
fi
```

#### Ejecutar Validación

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Crear script
cat > scripts/validate-fixes.sh << 'SCRIPT'
[contenido del script de arriba]
SCRIPT

chmod +x scripts/validate-fixes.sh

# Ejecutar
./scripts/validate-fixes.sh
```

#### Criterios de Validación

- ✅ Todos los módulos compilan
- ✅ Tests unitarios pasan (`-short`)
- ✅ Race detector no encuentra problemas
- ✅ Script validate-fixes.sh pasa

#### Commit de Correcciones

```bash
git add .
git commit -m "fix: corregir fallos críticos de CI

Correcciones implementadas basadas en Root Cause Analysis.

Cambios principales:
- [Listar cambios específicos según lo que corregiste]
- Tests de integración skippeados con -short
- Dependencias de edugo-shared actualizadas
- Race conditions corregidas en [módulos]
- Configuración de GOPRIVATE en workflows

Validaciones:
- ✅ Compilación exitosa en todos los módulos
- ✅ Tests unitarios pasan
- ✅ Race detector limpio

Refs: docs/troubleshooting/ROOT-CAUSE-ANALYSIS-[fecha].md

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### ✅ Tarea 2.2: Migrar a Go 1.25

**Prioridad:** 🔴 Alta  
**Estimación:** ⏱️ 45 minutos  
**Prerequisitos:** Tarea 2.1 completada

#### Objetivo

Estandarizar infrastructure en Go 1.25, igual que shared.

#### Script de Migración

```bash
#!/bin/bash
# migrate-to-go-1.25.sh

set -e

echo "🚀 Migrando edugo-infrastructure a Go 1.25..."
echo ""

# 1. Actualizar go.mod de cada módulo
MODULES=("postgres" "mongodb" "messaging" "schemas")

for module in "${MODULES[@]}"; do
  if [ -f "$module/go.mod" ]; then
    echo "📝 Actualizando $module/go.mod..."
    cd "$module"
    
    # Actualizar directiva go
    sed -i '' 's/^go 1\.24/go 1.25/' go.mod
    sed -i '' 's/^go 1\.23/go 1.25/' go.mod
    
    # go mod tidy
    go mod tidy
    
    cd ..
  fi
done

# 2. Actualizar workflows
echo ""
echo "📝 Actualizando workflows..."

for workflow in .github/workflows/*.yml; do
  echo "  - Actualizando $workflow"
  
  # Actualizar go-version en setup-go
  sed -i '' 's/go-version: "1\.24"/go-version: "1.25"/' "$workflow"
  sed -i '' 's/go-version: "1\.23"/go-version: "1.25"/' "$workflow"
  sed -i '' "s/go-version: '1\.24'/go-version: '1.25'/" "$workflow"
  sed -i '' "s/go-version: '1\.23'/go-version: '1.25'/" "$workflow"
done

# 3. Actualizar README si menciona versión
if [ -f README.md ] && grep -q "Go 1\.24\|Go 1\.23" README.md; then
  echo ""
  echo "📝 Actualizando README.md..."
  sed -i '' 's/Go 1\.24/Go 1.25/g' README.md
  sed -i '' 's/Go 1\.23/Go 1.25/g' README.md
fi

echo ""
echo "✅ Migración completada"
echo ""
echo "Verificando cambios..."
git diff --stat
```

#### Ejecutar Migración

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Crear script
cat > scripts/migrate-to-go-1.25.sh << 'SCRIPT'
[contenido del script de arriba]
SCRIPT

chmod +x scripts/migrate-to-go-1.25.sh

# Ejecutar
./scripts/migrate-to-go-1.25.sh
```

#### Validar Cambios

```bash
# Verificar go.mod actualizados
for module in postgres mongodb messaging schemas; do
  echo "$module:"
  grep "^go " "$module/go.mod"
done

# Verificar workflows
grep "go-version" .github/workflows/*.yml

# Validar que todo compila con Go 1.25
for module in postgres mongodb messaging schemas; do
  cd "$module"
  echo "Compilando $module con Go 1.25..."
  go version
  go build ./...
  cd ..
done
```

#### Commit

```bash
git add .
git commit -m "chore: migrar a Go 1.25

Estandarización con shared y resto del ecosistema.

Cambios:
- go.mod: go 1.25 en todos los módulos
- Workflows: go-version: \"1.25\"
- README: Actualizado

Validaciones:
- ✅ go mod verify en todos los módulos
- ✅ go build exitoso
- ✅ Alineado con shared

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

### ✅ Tarea 2.3: Validar Workflows Localmente con act

**Prioridad:** 🟡 Media (Opcional)  
**Estimación:** ⏱️ 60 minutos  
**Prerequisitos:** Tareas 2.1 y 2.2 completadas

#### Objetivo

Validar workflows localmente antes de pushear para evitar más fallos en CI.

#### Instalación de act (si no está instalado)

```bash
# macOS
brew install act

# Verificar
act --version
```

#### Validar Workflows

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Listar workflows disponibles
act -l

# Validar ci.yml (dry-run primero)
act pull_request --dryrun

# Si dry-run pasa, ejecutar (puede tomar tiempo)
# act pull_request

# Validar sintaxis sin ejecutar
for workflow in .github/workflows/*.yml; do
  echo "Validando $workflow..."
  act -W "$workflow" --list || echo "❌ Error en $workflow"
done
```

#### Alternativa: GitHub API

```bash
# Validar workflows sin ejecutar localmente
gh api \
  --method GET \
  /repos/EduGoGroup/edugo-infrastructure/actions/workflows \
  --jq '.workflows[] | {name, path, state}'
```

#### Esta Tarea es Opcional

Si act causa problemas, está bien saltarla. La validación real será en GitHub.

---

### ✅ Tarea 2.4: Validar Tests de Todos los Módulos

**Prioridad:** 🔴 Alta  
**Estimación:** ⏱️ 60 minutos  
**Prerequisitos:** Tareas 2.1 y 2.2 completadas

#### Objetivo

Ejecutar suite completa de tests para asegurar que las correcciones funcionan.

#### Script de Testing Completo

```bash
#!/bin/bash
# test-all-modules.sh

set -e

echo "🧪 Ejecutando tests completos con Go 1.25..."
echo "Versión de Go: $(go version)"
echo ""

MODULES=("postgres" "mongodb" "messaging" "schemas")
SUCCESS=0
FAILED=0

mkdir -p logs/test-reports

for module in "${MODULES[@]}"; do
  if [ ! -d "$module" ]; then
    echo "⚠️  Módulo $module no encontrado"
    continue
  fi
  
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "🧪 Testing: $module"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  cd "$module"
  
  LOG_FILE="../logs/test-reports/$module-$(date +%Y%m%d-%H%M%S).log"
  
  # Tests con coverage
  if go test -short -v -race -cover -coverprofile=coverage.out ./... 2>&1 | tee "$LOG_FILE"; then
    echo "✅ Tests de $module pasaron"
    SUCCESS=$((SUCCESS + 1))
    
    # Mostrar cobertura
    if [ -f coverage.out ]; then
      COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $NF}')
      echo "📊 Cobertura: $COVERAGE"
      rm coverage.out
    fi
  else
    echo "❌ Tests de $module fallaron"
    FAILED=$((FAILED + 1))
  fi
  
  echo ""
  cd ..
done

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 RESUMEN"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ Exitosos: $SUCCESS"
echo "❌ Fallidos: $FAILED"
echo "📦 Total: ${#MODULES[@]}"
echo ""

if [ $FAILED -eq 0 ]; then
  echo "🎉 Todos los tests pasaron"
  echo "✅ Listo para push y CI"
  exit 0
else
  echo "⚠️  Algunos tests fallaron"
  echo "🔧 Revisar logs en logs/test-reports/"
  exit 1
fi
```

#### Ejecutar Tests

```bash
cd /Users/jhoanmedina/source/EduGo/repos-separados/edugo-infrastructure

# Crear script
cat > scripts/test-all-modules.sh << 'SCRIPT'
[contenido del script de arriba]
SCRIPT

chmod +x scripts/test-all-modules.sh

# Ejecutar
./scripts/test-all-modules.sh
```

#### Commit de Validación

```bash
git add logs/ scripts/
git commit -m "test: validar todos los módulos post-correcciones

Tests completos ejecutados con Go 1.25.

Resultados:
- ✅ Compilación exitosa en 4/4 módulos
- ✅ Tests unitarios pasando
- ✅ Race detector limpio
- 📊 Cobertura baseline documentada

Archivos:
- scripts/test-all-modules.sh
- logs/test-reports/ (logs detallados)

🤖 Generated with Claude Code

Co-Authored-By: Claude <noreply@anthropic.com>"
```

---

## DÍA 3: ESTANDARIZACIÓN

[Debido al límite de longitud, continuaré con resumen estructurado]

---

### Tareas Restantes del Sprint 1

#### Día 3: Estandarización (3-4h)
- **Tarea 3.1:** Alinear workflows con shared (90 min)
  - Copiar patrón de workflows de shared
  - Agregar conditional para evitar "fallos fantasma"
  - Estandarizar nombres
  
- **Tarea 3.2:** Implementar pre-commit hooks (60 min)
  - Copiar hooks de shared
  - Adaptar para estructura de infrastructure
  - Script setup-hooks.sh
  
- **Tarea 3.3:** Documentar configuración (45 min)
  - Crear WORKFLOWS.md
  - Actualizar README con badges
  - Documentar módulos

#### Día 4: Validación y Deploy (2-3h)
- **Tarea 4.1:** Testing exhaustivo en GitHub (60 min)
  - Push y ejecutar workflows
  - Validar al menos 3 ejecuciones exitosas
  
- **Tarea 4.2:** PR, review y merge (45 min)
  - Crear PR con template completo
  - Self-review checklist
  - Merge a dev
  
- **Tarea 4.3:** Validar success rate (30 min)
  - Ejecutar 5+ veces más
  - Confirmar success rate > 95%
  - Actualizar documentación

---

## 📊 Métricas de Éxito del Sprint 1

### Pre-Sprint 1
```yaml
success_rate: 20%
fallos_consecutivos: 8
go_version: "1.24 (inconsistente)"
workflows: 2 (básicos)
pre_commit_hooks: false
```

### Post-Sprint 1 (Objetivo)
```yaml
success_rate: 100%
fallos_consecutivos: 0
go_version: "1.25 (estandarizado)"
workflows: 2 (optimizados)
pre_commit_hooks: true
documentacion: "Completa"
```

---

## 🎯 Checkpoint Final

```bash
# Validar que todo está listo
✅ 8 fallos corregidos
✅ Success rate 20% → 100%
✅ Go 1.25 estandarizado
✅ Workflows alineados con shared
✅ Pre-commit hooks funcionando
✅ Documentación completa
✅ 5+ ejecuciones exitosas en GitHub
✅ PR mergeado a dev
```

---

**¡Sprint 1 Completado!**

**Siguiente:** Sprint 4 - Workflows Reusables (ver SPRINT-4-TASKS.md)

---

**Generado por:** Claude Code  
**Fecha:** 19 de Noviembre, 2025  
**Versión:** 1.0  
**Basado en:** Plan de shared v1.0

# Guía de Implementación: edugo-infrastructure v0.8.0

**Para:** edugo-worker  
**Fecha:** 2025-11-18  
**Prioridad:** ALTA  
**Cambios:** ✅ Nuevas collections disponibles + Breaking change en postgres

---

## 🎯 RESUMEN EJECUTIVO

**¡Buenas noticias!** Tus 3 collections MongoDB ya están implementadas en edugo-infrastructure v0.8.0:

✅ `material_summary`  
✅ `material_assessment_worker`  
✅ `material_event`

**Además:** El módulo `migrations/` fue eliminado (si lo usabas).

---

## 📦 NUEVAS COLLECTIONS DISPONIBLES

### 1. material_summary
- **Migración:** 007_create_material_summary.{up,down}.js
- **Seed:** material_summary.js (3 documentos)
- **Índices:** 4 (material_id unique)
- **Características:** Multi-idioma, versionado, tracking tokens

### 2. material_assessment_worker
- **Migración:** 008_create_material_assessment_worker.{up,down}.js
- **Seed:** material_assessment_worker.js (2 documentos)
- **Índices:** 5 (material_id unique)
- **Características:** 3-20 preguntas, dificultad, explicaciones

### 3. material_event
- **Migración:** 009_create_material_event.{up,down}.js
- **Seed:** material_event.js (5 documentos)
- **Índices:** 7 (1 con TTL 90 días)
- **Características:** Auditoría, retry tracking, auto-limpieza

---

## 🔧 PASOS DE IMPLEMENTACIÓN

### Paso 1: Actualizar Dependencia MongoDB

```bash
cd edugo-worker

# Actualizar mongodb a v0.6.0
go get github.com/EduGoGroup/edugo-infrastructure/mongodb@v0.6.0

# Si usabas migrations/, también actualizar postgres
go get github.com/EduGoGroup/edugo-infrastructure/postgres@v0.8.0

# Limpiar dependencias
go mod tidy
```

---

### Paso 2: Actualizar Import (SI USABAS migrations/)

**Solo si worker usa el módulo `migrations/` para testing:**

**ANTES:**
```go
import "github.com/EduGoGroup/edugo-infrastructure/migrations"
```

**DESPUÉS:**
```go
import pgtesting "github.com/EduGoGroup/edugo-infrastructure/postgres/testing"
```

**Nota:** Si NO usas migrations/, ignora este paso.

---

### Paso 3: Actualizar Código de Worker

**Ahora puedes usar las nuevas collections en tu código:**

#### Ejemplo 1: Guardar Resumen

```go
import (
    "go.mongodb.org/mongo-driver/mongo"
    "time"
)

func SaveSummary(db *mongo.Database, materialID string, summary string, keyPoints []string) error {
    collection := db.Collection("material_summary")
    
    doc := bson.M{
        "material_id":        materialID,
        "summary":            summary,
        "key_points":         keyPoints,
        "language":           "es",
        "word_count":         len(strings.Fields(summary)),
        "version":            1,
        "ai_model":           "gpt-4",
        "processing_time_ms": 3500,
        "created_at":         time.Now(),
        "updated_at":         time.Now(),
    }
    
    _, err := collection.InsertOne(context.Background(), doc)
    return err
}
```

#### Ejemplo 2: Guardar Assessment

```go
func SaveAssessment(db *mongo.Database, materialID string, questions []Question) error {
    collection := db.Collection("material_assessment_worker")
    
    doc := bson.M{
        "material_id":        materialID,
        "questions":          questions,
        "total_questions":    len(questions),
        "total_points":       calculateTotalPoints(questions),
        "version":            1,
        "ai_model":           "gpt-4",
        "processing_time_ms": 5000,
        "created_at":         time.Now(),
        "updated_at":         time.Now(),
    }
    
    _, err := collection.InsertOne(context.Background(), doc)
    return err
}
```

#### Ejemplo 3: Registrar Evento

```go
func LogEvent(db *mongo.Database, eventType string, materialID string, payload interface{}) error {
    collection := db.Collection("material_event")
    
    doc := bson.M{
        "event_type":  eventType,
        "material_id": materialID,
        "payload":     payload,
        "status":      "pending",
        "retry_count": 0,
        "created_at":  time.Now(),
        "updated_at":  time.Now(),
    }
    
    _, err := collection.InsertOne(context.Background(), doc)
    return err
}
```

---

### Paso 4: Actualizar Tests de Integración

**Si usas testcontainers:**

```go
func setupMongoDB(t *testing.T) *mongo.Database {
    // ... setup testcontainer ...
    
    // Aplicar migraciones de infrastructure
    cmd := exec.Command("go", "run", 
        "github.com/EduGoGroup/edugo-infrastructure/mongodb/migrate.go", 
        "up")
    cmd.Env = append(os.Environ(), 
        "MONGO_URI="+containerURI,
        "MONGO_DB=edugo_test")
    
    if err := cmd.Run(); err != nil {
        t.Fatalf("Error aplicando migraciones: %v", err)
    }
    
    return db
}
```

---

### Paso 5: Eliminar Scripts Locales (OPCIONAL)

Si tenías scripts locales de MongoDB que ahora están en infrastructure:

```bash
# Revisar qué puedes eliminar
ls -la scripts/mongodb/

# Eliminar redundantes (después de verificar)
rm scripts/mongodb/init_collections.js  # Ya en infrastructure
rm scripts/mongodb/seed_data.js         # Ya en infrastructure
```

**Mantén solo:** Scripts específicos de worker que NO estén en infrastructure.

---

### Paso 6: Verificar Todo Funciona

```bash
# 1. Compilar
go build ./...

# 2. Tests
go test ./...

# 3. Tests de integración (con migraciones)
make test-integration

# 4. Probar manualmente (opcional)
docker-compose up -d mongodb
cd ../edugo-infrastructure/mongodb
go run migrate.go up
make seed
```

---

### Paso 7: Commit y Push

```bash
git add .
git commit -m "chore: actualizar a edugo-infrastructure v0.8.0

- Actualizar mongodb a v0.6.0 (collections worker)
- Actualizar postgres a v0.8.0 (si usa testing helpers)
- Usar collections de infrastructure en lugar de scripts locales
- Eliminar scripts redundantes

Relacionado: edugo-infrastructure v0.8.0"

git push origin <tu-rama>
```

---

## ⚠️ IMPORTANTE

### Nombre de Collection: `material_assessment_worker`

La collection se llama **`material_assessment_worker`** (NO `material_assessment`) para evitar conflicto con la collection existente de api-admin.

**Actualiza tu código para usar este nombre:**

```go
// CORRECTO
collection := db.Collection("material_assessment_worker")

// INCORRECTO (conflicto con api-admin)
collection := db.Collection("material_assessment")
```

---

## ✅ CHECKLIST COMPLETO

### Actualización de Dependencias
- [ ] `go get mongodb@v0.6.0` ejecutado
- [ ] `go get postgres@v0.8.0` ejecutado (si usa)
- [ ] `go mod tidy` ejecutado

### Código
- [ ] Imports actualizados (si usa migrations/)
- [ ] Collections actualizadas a nombres correctos
- [ ] Scripts locales redundantes eliminados

### Testing
- [ ] `go build ./...` exitoso
- [ ] Tests unitarios: PASS
- [ ] Tests de integración: PASS

### Git
- [ ] Cambios commiteados
- [ ] Push realizado

---

## 📊 NUEVAS CAPACIDADES

Ahora puedes:
- ✅ Usar migraciones de infrastructure en testcontainers
- ✅ Compartir seeds entre proyectos
- ✅ Eliminar duplicación de scripts MongoDB
- ✅ Beneficiarte de TTL index automático (90 días)
- ✅ Usar validaciones JSON Schema estrictas

---

## ❓ FAQ

### ¿Debo cambiar el nombre de mis collections en código?
Sí, si usabas `material_assessment`, cámbialo a `material_assessment_worker`.

### ¿Las migraciones son retrocompatibles?
Sí, puedes ejecutarlas sobre BD existentes sin problemas.

### ¿El TTL afecta datos existentes?
Sí, material_event empezará a eliminar documentos >90 días automáticamente.

---

## 📞 SOPORTE

Cualquier problema, contacta al equipo de infrastructure.

---

**Generado por:** edugo-infrastructure  
**Versión:** v0.8.0  
**Fecha:** 2025-11-18

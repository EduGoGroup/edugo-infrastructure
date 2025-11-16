# PHASE2 BRIDGE - [Sprint Name]

## 📋 Resumen

**Sprint:** [Sprint-XX-Name]
**Archivo principal:** [ruta/al/archivo.go]
**Estado Fase 1:** ✅ COMPLETADO

---

## ✅ Completado en Fase 1

### Implementación

- [x] [Funcionalidad 1]
- [x] [Funcionalidad 2]
- [x] [Funcionalidad N]

### Tests Unitarios

- [x] Test para [función 1]
- [x] Test para [función 2]
- [x] Test para [función N]

### Código

```go
// Snippet de código relevante implementado
```

**Total de líneas:** XXX
**Cobertura de tests unitarios:** XX%

---

## ⏳ Pendiente para Fase 2

### Tests de Integración

1. **[Test case 1]**
   - Descripción: [...]
   - Requiere: PostgreSQL / RabbitMQ / etc.
   - Validar: [...]

2. **[Test case 2]**
   - Descripción: [...]
   - Requiere: [...]
   - Validar: [...]

### Edge Cases

1. **[Caso especial 1]**
   - Escenario: [...]
   - Validación: [...]

2. **[Caso especial 2]**
   - Escenario: [...]
   - Validación: [...]

---

## 🔧 Prerequisitos para Fase 2

### Servicios Requeridos

```bash
# [Servicio 1]
[comando para levantar]

# [Servicio 2]
[comando para levantar]
```

### Variables de Entorno

```bash
# Copiar .env.example
cp .env.example .env

# Variables necesarias:
DB_HOST=localhost
DB_PORT=5432
# ... etc
```

### Datos de Prueba

```bash
# Cargar seeds si es necesario
make seed
```

---

## 🧪 Tests de Integración a Implementar

### Archivo: `[nombre]_integration_test.go`

```go
// Pseudocódigo de tests a implementar

func TestXxxIntegration(t *testing.T) {
    // Setup Testcontainers
    // Ejecutar función
    // Validar resultado con BD real
}
```

### Casos de Prueba

1. **Happy path:** [...]
2. **Error handling:** [...]
3. **Rollback:** [...]
4. **Performance:** [...]

---

## 📝 Notas para Fase 2

- [Nota importante 1]
- [Nota importante 2]
- [Consideración especial]

---

## ✅ Checklist Fase 2

- [ ] Levantar servicios necesarios
- [ ] Configurar variables de entorno
- [ ] Implementar tests de integración
- [ ] Validar edge cases
- [ ] Actualizar documentación
- [ ] Commit y push

---

**Fase 1 completada:** [Fecha]
**Próximo paso:** Ejecutar PHASE2_PROMPT.txt

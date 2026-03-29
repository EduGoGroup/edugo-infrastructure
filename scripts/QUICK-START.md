# Quick Start - Auto Release Script

## 🚀 Uso Rápido en 3 Pasos

### Paso 1: Modifica el CHANGELOG

```bash
# Edita el CHANGELOG del módulo que quieres liberar
vim postgres/CHANGELOG.md
```

Asegúrate de que tenga este formato:

```markdown
## [Unreleased]

## [0.77.1] - 2026-03-29
### Changed
- Tu cambio aquí
```

**NO hagas commit todavía!**

### Paso 2: Ejecuta el Script

```bash
# Opción A: Usando Make (recomendado)
make auto-release

# Opción B: Con dry-run para ver qué haría
make auto-release-dry-run

# Opción C: Múltiples módulos automáticamente
make auto-release-all

# Opción D: Directamente con el script
./scripts/auto-release.sh
```

### Paso 3: Confirma y Listo

El script te mostrará un resumen y pedirá confirmación:

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RESUMEN DE RELEASES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Se crearán los siguientes releases:
  1. postgres/v0.77.1

¿Continuar con los commits, tags y push? (y/N): y
```

Escribe `y` y presiona Enter. El script:
- ✅ Hará commit del CHANGELOG
- ✅ Creará el tag `postgres/v0.77.1`
- ✅ Hará push al remote
- ✅ Activará GitHub Actions para crear el release

## 📋 Comandos Útiles

### Usando Make (Recomendado)

```bash
# Release interactivo
make auto-release

# Dry-run con verbose
make auto-release-dry-run

# Todos los módulos automáticamente
make auto-release-all

# Ver ayuda
make auto-release-help
```

### Usando el Script Directamente

```bash
# Ver ayuda completa
./scripts/auto-release.sh --help

# Dry-run con verbose (ver qué haría sin hacer cambios)
./scripts/auto-release.sh --dry-run --verbose

# Procesar solo módulos específicos
./scripts/auto-release.sh postgres mongodb

# Modo no interactivo (para scripts)
./scripts/auto-release.sh --all --yes
```

## ✅ Checklist Pre-Release

Antes de ejecutar el script, verifica:

- [ ] Has modificado el CHANGELOG.md del módulo
- [ ] La versión está en formato `[X.Y.Z]` (sin la `v`)
- [ ] La versión está después de `[Unreleased]`
- [ ] **NO has hecho commit** del CHANGELOG todavía
- [ ] El módulo pasa `make release-check` (el script lo verifica)

## 🎯 Ejemplo Completo

```bash
# 1. Verifica que no tienes cambios commiteados
git status
# Output: On branch main, nothing to commit

# 2. Edita el CHANGELOG
vim postgres/CHANGELOG.md
# Agrega tu nueva versión después de [Unreleased]

# 3. Verifica los cambios
git diff postgres/CHANGELOG.md

# 4. Ejecuta el script
make auto-release
# O directamente: ./scripts/auto-release.sh

# 5. El script detectará automáticamente:
#    - Módulo: postgres
#    - Versión: 0.77.1
#    - Tag: postgres/v0.77.1

# 6. Confirma cuando se te pregunte
# ¿Continuar? (y/N): y

# 7. ¡Listo! Verifica en GitHub Actions
# https://github.com/EduGoGroup/edugo-infrastructure/actions
```

## 🐛 Problemas Comunes

### "No hay CHANGELOGs modificados"

**Solución:** Asegúrate de haber modificado el CHANGELOG sin commitear:

```bash
git status  # Debe mostrar el CHANGELOG como modificado
```

### "No se pudo extraer versión"

**Solución:** Verifica el formato del CHANGELOG:

```bash
head -20 postgres/CHANGELOG.md
# Debe tener:
# ## [Unreleased]
# 
# ## [0.77.1] - 2026-03-29
```

### "El tag ya existe"

**Solución:** Incrementa la versión en el CHANGELOG:

```bash
git tag -l "postgres/*"  # Ver tags existentes
vim postgres/CHANGELOG.md  # Usar siguiente versión
```

## 📚 Documentación Completa

Para más detalles, consulta: [AUTO-RELEASE-README.md](AUTO-RELEASE-README.md)

---

**Tip:** Siempre usa `--dry-run` primero para ver qué hará el script antes de ejecutarlo realmente.
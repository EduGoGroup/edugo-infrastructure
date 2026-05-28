# Playground seeds

Datasets pequeños y focalizados para iterar la app carpeta a carpeta. Cada
playground se aplica **encima de L0** (el piso mínimo del sistema), nunca
encima de L1..L4 ni de demo.

## Reglas

- Cada playground vive en su propio subpaquete: `playground/<name>/<name>.go`.
- Expone una función `Apply(tx *gorm.DB) error` **idempotente**.
- Se registra en `playground.go` (`Available()` y el switch de `Apply()`).
- No participa de `ComputeFilesHash()` → cambiar un playground **no** requiere
  bump de `SchemaVersion`. Para iteración rápida.

## Disponibles

- `admin` — usuario `admin@edugo.local` / `12345678` con rol L0 super_admin y
  grant `*` (acceso total). Punto de partida para validar login y menú base.

## Cómo correrlo

Desde `edugo-dev-environment/migrator/`:

```bash
make docker-playground P=admin       # recrea BD docker + L0 + playground
make docker-playground-list          # lista los disponibles
```

El flag `--playground=<name>` del binario fuerza internamente:

- `FORCE_MIGRATION=true` (recreación completa de schemas).
- `SeedUpToLayer=l0` (solo el piso).
- `SeedDemo=false` (sin datos demo).

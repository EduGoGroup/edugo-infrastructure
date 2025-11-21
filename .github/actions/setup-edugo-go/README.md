# Setup EduGo Go Environment

Composite action para configurar el entorno Go estandar de EduGo.

---

## Caracteristicas

- Setup de Go con version configurable
- Configuracion automatica de GOPRIVATE
- Acceso a repos privados de EduGoGroup
- Cache de Go modules
- Verificacion de configuracion
- Output de version y cache hit

---

## Uso Basico

```yaml
steps:
  - uses: actions/checkout@v4

  - name: Setup Go
    uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
```

---

## Uso Avanzado

```yaml
steps:
  - name: Setup Go
    uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
    with:
      go-version: '1.25'
      cache: true
      cache-dependency-path: '**/go.sum'
      github-token: ${{ secrets.GITHUB_TOKEN }}
```

---

## Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `go-version` | No | `1.25` | Version de Go |
| `cache` | No | `true` | Habilitar cache |
| `cache-dependency-path` | No | `go.sum` | Path para cache |
| `github-token` | No | `github.token` | Token para repos privados |

---

## Outputs

| Output | Description |
|--------|-------------|
| `go-version` | Version de Go instalada |
| `cache-hit` | Si el cache fue encontrado (`true`/`false`) |

---

## Equivalencia

**Antes (15+ lineas):**
```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.25'
    cache: true

- name: Configurar repos privados
  run: |
    git config --global url."https://${{ secrets.GITHUB_TOKEN }}@github.com/".insteadOf "https://github.com/"
  env:
    GOPRIVATE: github.com/EduGoGroup/*

- name: Verificar
  run: |
    echo "Go: $(go version)"
    echo "GOPRIVATE: $GOPRIVATE"
```

**Despues (1 linea):**
```yaml
- uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main
```

---

**Reduccion:** ~93% menos codigo (15 lineas -> 1 linea)

---

**Mantenido por:** EduGo Team
**Ultima actualizacion:** 21 Nov 2025

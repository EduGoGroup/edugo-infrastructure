# EduGo Coverage Check

Composite action para validar cobertura de tests y generar reportes.

---

## Caracteristicas

- Calcula cobertura total del proyecto
- Valida contra umbral configurable
- Genera reportes HTML opcionales
- Outputs para integracion con PRs
- Detalles por paquete

---

## Uso Basico

```yaml
steps:
  - uses: actions/checkout@v4

  - name: Setup Go
    uses: EduGoGroup/edugo-infrastructure/.github/actions/setup-edugo-go@main

  - name: Run tests with coverage
    run: go test -coverprofile=coverage.out ./...

  - name: Check coverage
    uses: EduGoGroup/edugo-infrastructure/.github/actions/coverage-check@main
```

---

## Uso Avanzado

```yaml
steps:
  - name: Check coverage
    id: coverage
    uses: EduGoGroup/edugo-infrastructure/.github/actions/coverage-check@main
    with:
      coverage-file: 'coverage.out'
      threshold: '50'
      fail-on-threshold: true
      generate-html: true
      html-output: 'coverage-report.html'
      working-directory: './app'

  - name: Use coverage output
    run: |
      echo "Coverage: ${{ steps.coverage.outputs.coverage }}%"
      echo "Passed: ${{ steps.coverage.outputs.passed }}"
```

---

## Inputs

| Input | Required | Default | Description |
|-------|----------|---------|-------------|
| `coverage-file` | No | `coverage.out` | Archivo de cobertura |
| `threshold` | No | `33` | Umbral minimo (%) |
| `fail-on-threshold` | No | `true` | Fallar si bajo umbral |
| `generate-html` | No | `false` | Generar reporte HTML |
| `html-output` | No | `coverage.html` | Nombre archivo HTML |
| `working-directory` | No | `.` | Directorio de trabajo |

---

## Outputs

| Output | Description |
|--------|-------------|
| `coverage` | Porcentaje de cobertura |
| `passed` | Si paso el umbral (`true`/`false`) |
| `report` | Resumen del reporte |

---

## Ejemplo en PR

```yaml
- name: Check coverage
  id: coverage
  uses: EduGoGroup/edugo-infrastructure/.github/actions/coverage-check@main

- name: Comment PR
  if: github.event_name == 'pull_request'
  uses: actions/github-script@v7
  with:
    script: |
      github.rest.issues.createComment({
        issue_number: context.issue.number,
        owner: context.repo.owner,
        repo: context.repo.repo,
        body: 'Coverage: ${{ steps.coverage.outputs.coverage }}%'
      })
```

---

**Mantenido por:** EduGo Team
**Ultima actualizacion:** 21 Nov 2025

# Tareas: Sprint-01-Migrate-CLI

**Sprint:** Sprint-01-Migrate-CLI  
**Duración:** 1-2 horas

---

## TASK-001: Crear database/migrate.go

**Descripción:** Implementar CLI para gestionar migraciones

**Pasos:**

1. Crear archivo `database/migrate.go`

```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func main() {
    if len(os.Args) < 2 {
        printUsage()
        os.Exit(1)
    }

    command := os.Args[1]

    // Leer DATABASE_URL del entorno
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://edugo:changeme@localhost:5432/edugo_dev?sslmode=disable"
    }

    // Conectar a base de datos
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    defer db.Close()

    // Crear driver de migrate
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatalf("Error creating migrate driver: %v", err)
    }

    // Crear instancia de migrate
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres",
        driver,
    )
    if err != nil {
        log.Fatalf("Error creating migrate instance: %v", err)
    }

    // Ejecutar comando
    switch command {
    case "up":
        if err := m.Up(); err != nil && err != migrate.ErrNoChange {
            log.Fatalf("Error running migrations up: %v", err)
        }
        fmt.Println("✅ Migrations applied successfully")

    case "down":
        if err := m.Down(); err != nil && err != migrate.ErrNoChange {
            log.Fatalf("Error running migrations down: %v", err)
        }
        fmt.Println("✅ Migrations rolled back successfully")

    case "status":
        version, dirty, err := m.Version()
        if err != nil {
            fmt.Println("No migrations applied yet")
            return
        }
        status := "clean"
        if dirty {
            status = "dirty"
        }
        fmt.Printf("Current version: %d (%s)\n", version, status)

    case "create":
        if len(os.Args) < 3 {
            fmt.Println("Usage: go run migrate.go create <migration_name>")
            os.Exit(1)
        }
        createMigration(os.Args[2])

    default:
        fmt.Printf("Unknown command: %s\n", command)
        printUsage()
        os.Exit(1)
    }
}

func printUsage() {
    fmt.Println("Usage: go run migrate.go <command>")
    fmt.Println("")
    fmt.Println("Commands:")
    fmt.Println("  up       - Apply all pending migrations")
    fmt.Println("  down     - Rollback last migration")
    fmt.Println("  status   - Show current migration version")
    fmt.Println("  create   - Create new migration files")
    fmt.Println("")
    fmt.Println("Environment variables:")
    fmt.Println("  DATABASE_URL - PostgreSQL connection string")
    fmt.Println("                 Default: postgres://edugo:changeme@localhost:5432/edugo_dev?sslmode=disable")
}

func createMigration(name string) {
    // Encontrar siguiente número
    entries, err := os.ReadDir("migrations")
    if err != nil {
        log.Fatalf("Error reading migrations directory: %v", err)
    }

    nextNum := 1
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        // Extraer número del nombre (XXX_name.up.sql)
        var num int
        fmt.Sscanf(entry.Name(), "%d_", &num)
        if num >= nextNum {
            nextNum = num + 1
        }
    }

    // Crear archivos
    upFile := fmt.Sprintf("migrations/%03d_%s.up.sql", nextNum, name)
    downFile := fmt.Sprintf("migrations/%03d_%s.down.sql", nextNum, name)

    // Crear archivo UP
    if err := os.WriteFile(upFile, []byte("-- Add migration SQL here\n"), 0644); err != nil {
        log.Fatalf("Error creating UP file: %v", err)
    }

    // Crear archivo DOWN
    if err := os.WriteFile(downFile, []byte("-- Add rollback SQL here\n"), 0644); err != nil {
        log.Fatalf("Error creating DOWN file: %v", err)
    }

    fmt.Printf("✅ Created migration files:\n")
    fmt.Printf("   - %s\n", upFile)
    fmt.Printf("   - %s\n", downFile)
}
```

2. Actualizar `database/go.mod`

```go
module github.com/EduGoGroup/edugo-infrastructure/database

go 1.24

require (
    github.com/golang-migrate/migrate/v4 v4.17.0
    github.com/lib/pq v1.10.9
)
```

3. Ejecutar `go mod tidy` en database/

**Validación:**
```bash
cd database
go run migrate.go status
# Debe mostrar: No migrations applied yet (si es primera vez)

go run migrate.go up
# Debe aplicar migraciones 001-008

go run migrate.go status
# Debe mostrar: Current version: 8 (clean)
```

**Estimación:** 60-90 minutos

---

## TASK-002: Crear README de uso del CLI

**Descripción:** Documentar cómo usar migrate.go

**Pasos:**

1. Actualizar `database/README.md`

Agregar sección:

```markdown
## Uso del CLI de Migraciones

### Prerequisitos
- PostgreSQL 15+ corriendo
- Variable DATABASE_URL configurada (o usar default)

### Comandos

**Aplicar todas las migraciones:**
```bash
cd database
go run migrate.go up
```

**Revertir última migración:**
```bash
go run migrate.go down
```

**Ver estado actual:**
```bash
go run migrate.go status
```

**Crear nueva migración:**
```bash
go run migrate.go create add_new_table
# Crea: 009_add_new_table.up.sql y 009_add_new_table.down.sql
```

### Variables de Entorno

```bash
# Default (si no se especifica)
DATABASE_URL=postgres://edugo:changeme@localhost:5432/edugo_dev?sslmode=disable

# Para otro ambiente
export DATABASE_URL=postgres://user:pass@host:5432/database?sslmode=disable
go run migrate.go up
```

### En CI/CD

```yaml
- name: Run migrations
  env:
    DATABASE_URL: ${{ secrets.DATABASE_URL }}
  run: |
    cd database
    go run migrate.go up
```
```

**Estimación:** 15 minutos

---

## TASK-003: Validar con tests básicos

**Descripción:** Crear test básico del CLI

**Pasos:**

1. Crear `database/migrate_test.go`

```go
package main

import (
    "os"
    "testing"
)

func TestPrintUsage(t *testing.T) {
    // Capture stdout
    oldStdout := os.Stdout
    r, w, _ := os.Pipe()
    os.Stdout = w

    printUsage()

    w.Close()
    os.Stdout = oldStdout

    // Verificar que imprime algo
    buf := make([]byte, 1024)
    n, _ := r.Read(buf)
    
    if n == 0 {
        t.Error("printUsage should print something")
    }
}
```

2. Ejecutar test

```bash
cd database
go test -v
```

**Estimación:** 15 minutos

---

## ✅ Checklist de Completitud

- [ ] migrate.go creado
- [ ] go.mod actualizado con golang-migrate
- [ ] Comando `up` funciona
- [ ] Comando `down` funciona
- [ ] Comando `status` funciona
- [ ] Comando `create` funciona
- [ ] README.md actualizado
- [ ] Test básico creado
- [ ] Validación manual completada

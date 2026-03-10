# mock-generator processes

## Procesos propios del modulo

### 1. Leer un directorio de SQL

`pkg/parser/sql_parser.go` recorre archivos `*.sql` de un directorio y separa statements por `;`.

### 2. Parsear solo `INSERT`

El parser ignora statements que no puede leer o que no son `INSERT` y acumula filas por tabla.

### 3. Normalizar expresiones SQL

La evaluacion de expresiones convierte algunos casos a placeholders simples:

- `now()` -> `NOW()`
- `gen_random_uuid()` -> `UUID()`
- `NULL` -> `nil`

### 4. Generar paquete `dataset`

El generador crea:

- `helpers.go`
- `database.go`
- archivos por tabla
- loader de datos

### 5. Formatear codigo generado

Al final intenta ejecutar formateo sobre los archivos generados.

## Realidades que importan documentar

- Los defaults del CLI apuntan a `../../postgres/migrations/testing`, carpeta que hoy no existe.
- El mapping de tablas a entities en `pkg/types/mappings.go` cubre solo una porcion del dominio.
- El modulo hoy documenta mejor una intencion de generacion que una pipeline consolidada de uso diario.

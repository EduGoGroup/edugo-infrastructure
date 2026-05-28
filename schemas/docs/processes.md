# schemas processes

## Procesos propios del modulo

### 1. Cargar contratos embebidos

`validator.go` embebe `events/*.json` y carga todos los archivos `*.schema.json` al crear un `EventValidator`.

### 2. Derivar llave interna por nombre de archivo

La convencion actual deriva llaves como:

- `material-uploaded-v1.schema.json` -> `material.uploaded:1.0`
- `assessment-generated-v1.schema.json` -> `assessment.generated:1.0`

La convencion depende del nombre del archivo.

### 3. Validar eventos

La API publica permite tres caminos:

- `Validate(event interface{})`
- `ValidateWithType(loader, eventType, eventVersion)`
- `ValidateJSON(jsonBytes, eventType, eventVersion)`

### 4. Reportar errores de contrato

Si la validacion falla, el modulo acumula errores en un mensaje unico, orientado a debugging de payload.

### 5. Verificar comportamiento por tests y benchmarks

El modulo incluye tests de cobertura funcional y benchmarks para validacion repetida.

## Contratos observados

Actualmente existen 4 schemas:

- `assessment.generated` v1.0
- `material.deleted` v1.0
- `material.uploaded` v1.0
- `student.enrolled` v1.0

## Realidades que importan documentar

- La clave de lookup depende del nombre del archivo, no de un registro explicito aparte.
- El modulo es pequeno, pero altamente sensible a convenciones de naming.

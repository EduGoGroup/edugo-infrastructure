# 📖 Glosario de Términos - EduGo

Definición de términos del dominio educativo y técnico utilizados en el ecosistema EduGo.

---

## 📋 Índice Alfabético

- [A](#a) | [B](#b) | [C](#c) | [D](#d) | [E](#e) | [F](#f) | [G](#g) | [H](#h) | [I](#i) | [J](#j) | [K](#k) | [L](#l) | [M](#m) | [N](#n) | [O](#o) | [P](#p) | [Q](#q) | [R](#r) | [S](#s) | [T](#t) | [U](#u) | [V](#v) | [W](#w)

---

## A

### Academic Unit (Unidad Académica)
**Definición:** Entidad organizacional dentro de una escuela que puede representar diferentes niveles jerárquicos como grados, cursos, secciones, clubes o departamentos.

**Ejemplo:** "3° Medio A" es una unidad académica de tipo "class" que pertenece a "3° Medio" (tipo "grade").

**Tabla PostgreSQL:** `academic_units`

**Características:**
- Soporta jerarquía mediante `parent_unit_id`
- Tipos: school, grade, class, section, club, department
- Pertenece a una escuela específica

---

### Access Token
**Definición:** Token JWT de corta duración (15 minutos) usado para autenticar requests a las APIs.

**Características:**
- Contiene claims del usuario (id, email, role, school_ids)
- Se envía en header `Authorization: Bearer <token>`
- Al expirar, debe renovarse con refresh token

---

### Admin (Administrador)
**Definición:** Usuario con permisos completos para gestionar una escuela.

**Permisos:**
- Crear/editar unidades académicas
- Matricular/dar de baja usuarios
- Ver todos los reportes
- Gestionar configuración de la escuela

---

### Assessment
**Definición:** Evaluación o quiz generado automáticamente por IA a partir de un material educativo.

**Componentes:**
- Metadata en PostgreSQL (`assessment` table)
- Preguntas completas en MongoDB (`material_assessment_worker` collection)
- Referencia cruzada via `mongo_document_id`

**Estados:**
- `draft`: Creado pero no publicado
- `generated`: Generado por IA
- `published`: Disponible para estudiantes
- `archived`: Archivado, no visible
- `closed`: Cerrado para nuevos intentos

---

### Assessment Attempt (Intento de Assessment)
**Definición:** Registro de un estudiante tomando un assessment específico.

**Tabla PostgreSQL:** `assessment_attempt`

**Estados:**
- `in_progress`: Estudiante está respondiendo
- `submitted`: Enviado, pendiente de calificación
- `graded`: Calificado con score
- `abandoned`: Abandonado (timeout o cancelación)

---

### Assessment Attempt Answer (Respuesta de Intento)
**Definición:** Respuesta individual a una pregunta dentro de un intento.

**Tabla PostgreSQL:** `assessment_attempt_answer`

**Campos clave:**
- `question_index`: Índice de la pregunta (0-based)
- `student_answer`: Respuesta del estudiante
- `is_correct`: Resultado de la evaluación

---

## B

### Bearer Token
**Definición:** Esquema de autenticación donde el token se envía en el header HTTP.

**Formato:** `Authorization: Bearer eyJhbGciOiJIUzI1NiIs...`

---

### BSON
**Definición:** Binary JSON, formato de serialización usado por MongoDB.

**Uso en EduGo:** Tags `bson:"field_name"` en entities MongoDB.

---

## C

### Collection (MongoDB)
**Definición:** Equivalente a una tabla en bases de datos relacionales. Contenedor de documentos en MongoDB.

**Collections en EduGo:**
- `material_assessment_worker`
- `material_summary`
- `material_event`
- `schema_migrations`

---

### Consumer
**Definición:** Servicio que consume mensajes de una cola RabbitMQ.

**Consumers en EduGo:**
- Worker consume `material.uploaded` y `material.deleted`
- API Mobile consume `assessment.generated` y `student.enrolled`

---

### Coordinator (Coordinador)
**Definición:** Usuario que supervisa un conjunto de docentes y/o unidades académicas.

**Permisos:**
- Ver estadísticas de su área
- Gestionar membresías de su área
- Ver reportes de progreso

---

## D

### Dead Letter Queue (DLQ)
**Definición:** Cola especial donde van los mensajes que no pudieron procesarse exitosamente.

**Uso:** Debugging, reprocesamiento manual, análisis de errores.

---

### Document (MongoDB)
**Definición:** Registro individual en una collection de MongoDB, equivalente a una fila en SQL.

**Formato:** JSON/BSON con campos anidados y arrays.

---

## E

### Entity
**Definición:** Struct de Go que representa una tabla (PostgreSQL) o collection (MongoDB).

**Características:**
- Tags `db:"column"` para PostgreSQL
- Tags `bson:"field"` para MongoDB
- Método `TableName()` o `CollectionName()`

---

### Event (Evento)
**Definición:** Mensaje publicado en RabbitMQ para comunicación asíncrona entre servicios.

**Estructura base:**
```json
{
  "event_id": "uuid",
  "event_type": "material.uploaded",
  "event_version": "1.0",
  "timestamp": "ISO8601",
  "payload": {}
}
```

---

### Exchange (RabbitMQ)
**Definición:** Componente de RabbitMQ que recibe mensajes y los enruta a colas según reglas.

**Exchanges en EduGo:**
- `edugo.materials` (topic)
- `edugo.assessments` (topic)
- `edugo.students` (topic)

---

## F

### Foreign Key (FK)
**Definición:** Restricción de integridad referencial entre tablas.

**Ejemplo:** `materials.school_id` es FK a `schools.id`

---

## G

### Guardian (Apoderado)
**Definición:** Usuario que monitorea el progreso de uno o más estudiantes.

**Relación:** Tabla `guardian_relations` vincula guardian con estudiantes.

---

## H

### Hash (Password)
**Definición:** Representación cifrada de una contraseña usando bcrypt.

**Campo:** `users.password_hash`

---

## I

### Index (Índice)
**Definición:** Estructura de datos que mejora la velocidad de consultas.

**PostgreSQL:** `CREATE INDEX idx_name ON table(column)`
**MongoDB:** `db.collection.createIndex({field: 1})`

---

## J

### JSON Schema
**Definición:** Especificación para validar estructura de documentos JSON.

**Uso en EduGo:** Validar eventos antes de publicar/consumir.

**Ubicación:** `schemas/events/*.schema.json`

---

### JWT (JSON Web Token)
**Definición:** Estándar para tokens de autenticación.

**Partes:**
1. Header (algoritmo)
2. Payload (claims)
3. Signature (firma)

---

## K

### Key Points (Puntos Clave)
**Definición:** Lista de conceptos principales extraídos de un material por IA.

**Campo:** `material_summary.key_points` (MongoDB)

---

## L

### Latency (Latencia)
**Definición:** Tiempo de respuesta de una operación.

**Objetivos EduGo:**
- APIs: < 500ms
- Procesamiento IA: < 60s

---

## M

### Material
**Definición:** Archivo educativo subido por un docente (PDF, documento, etc).

**Tabla PostgreSQL:** `materials`

**Estados:**
- `uploaded`: Recién subido
- `processing`: Worker procesando
- `ready`: Listo con assessment generado
- `failed`: Error en procesamiento

---

### Membership (Membresía)
**Definición:** Relación entre un usuario, una escuela y opcionalmente una unidad académica.

**Tabla PostgreSQL:** `memberships`

**Roles posibles:**
- teacher, student, guardian
- coordinator, admin, assistant

---

### Migration (Migración)
**Definición:** Script que modifica el schema de base de datos de forma versionada.

**Archivos:**
- `XXX_name.up.sql` - Aplicar cambio
- `XXX_name.down.sql` - Revertir cambio

---

### MongoDB
**Definición:** Base de datos NoSQL orientada a documentos.

**Uso en EduGo:** Almacenar contenido de assessments (preguntas, opciones) y resúmenes.

---

## N

### Nullable
**Definición:** Campo que puede contener valor NULL.

**Go:** Usar punteros `*string`, `*time.Time`

---

## O

### ObjectId
**Definición:** Identificador único de 24 caracteres hexadecimales en MongoDB.

**Ejemplo:** `507f1f77bcf86cd799439011`

---

### OpenAI
**Definición:** Proveedor de IA usado para generar assessments y resúmenes.

**Modelos:** GPT-4, GPT-4-turbo

---

## P

### Payload
**Definición:** Datos principales de un evento o request.

---

### PostgreSQL
**Definición:** Base de datos relacional principal de EduGo.

**Uso:** Datos estructurados, transacciones ACID, relaciones.

---

### Progress (Progreso)
**Definición:** Registro del avance de un estudiante en materiales y assessments.

**Tabla PostgreSQL:** `progress`

---

### Publisher
**Definición:** Servicio que publica mensajes a RabbitMQ.

---

## Q

### Query
**Definición:** Consulta a base de datos.

---

### Question (Pregunta)
**Definición:** Elemento de un assessment con texto, opciones y respuesta correcta.

**Tipos:**
- `multiple_choice`: Opción múltiple
- `true_false`: Verdadero/Falso
- `open`: Respuesta abierta

---

### Queue (Cola)
**Definición:** Buffer de mensajes en RabbitMQ.

**Colas en EduGo:**
- `worker.materials.process`
- `api-mobile.assessments.ready`

---

## R

### RabbitMQ
**Definición:** Message broker para comunicación asíncrona.

**Puerto:** 5672 (AMQP), 15672 (Management UI)

---

### Refresh Token
**Definición:** Token de larga duración (7 días) para obtener nuevos access tokens.

---

### Role (Rol)
**Definición:** Tipo de usuario que determina permisos.

**Roles de usuario:** admin, teacher, student, guardian
**Roles de membresía:** teacher, student, guardian, coordinator, admin, assistant

---

### Routing Key
**Definición:** Clave usada por RabbitMQ para enrutar mensajes.

**Ejemplo:** `material.uploaded` → cola `worker.materials.process`

---

## S

### S3 (Amazon S3)
**Definición:** Servicio de almacenamiento de archivos en la nube.

**Uso:** Almacenar PDFs y materiales educativos.

---

### Schema
**Definición:** Estructura de datos (tablas, campos, tipos).

---

### School (Escuela)
**Definición:** Institución educativa registrada en EduGo.

**Tabla PostgreSQL:** `schools`

**Tiers de suscripción:** free, basic, premium, enterprise

---

### Score
**Definición:** Puntaje obtenido en un assessment (0-100).

---

### Seed
**Definición:** Datos iniciales para desarrollo/testing.

**Ubicación:** `seeds/postgres/`, `seeds/mongodb/`

---

### Soft Delete
**Definición:** Marcar registro como eliminado sin borrarlo físicamente.

**Campo:** `deleted_at` (timestamp o NULL)

---

### Student (Estudiante)
**Definición:** Usuario que consume contenido y rinde assessments.

---

### Subject (Materia)
**Definición:** Área de conocimiento (Matemáticas, Ciencias, etc).

**Campo:** `materials.subject`

---

### Summary (Resumen)
**Definición:** Texto resumido de un material generado por IA.

**Collection MongoDB:** `material_summary`

---

## T

### Teacher (Docente)
**Definición:** Usuario que sube materiales y ve progreso de estudiantes.

---

### Token
**Definición:** Cadena que representa autenticación o autorización.

---

### Transaction
**Definición:** Conjunto de operaciones que se ejecutan atómicamente.

---

## U

### Unit (Unidad)
**Definición:** División de contenido dentro de una materia.

**Tabla PostgreSQL:** `units`

---

### UUID
**Definición:** Identificador único universal de 128 bits.

**Formato:** `a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d`

---

## V

### Validator
**Definición:** Componente que valida datos contra reglas o schemas.

**Uso:** Validar eventos antes de publicar/consumir.

---

### Version (Versión)
**Definición:** Número que identifica una iteración de algo.

**Contextos:**
- Versión de evento: `1.0`
- Versión de migración: `001`, `002`
- Versión de material: `material_versions` table

---

## W

### Worker
**Definición:** Servicio que procesa tareas en background.

**Funciones:**
- Extraer texto de PDFs
- Generar assessments con IA
- Generar resúmenes con IA
- Actualizar estados en BD

---

### Workspace
**Definición:** Contexto de trabajo, generalmente una escuela.

---

## 📊 Relaciones entre Conceptos

```
School (Escuela)
├── Academic Units (Unidades Académicas)
│   └── [Jerarquía: grade → class → section]
├── Memberships (Membresías)
│   └── User + Role + Academic Unit
├── Materials (Materiales)
│   └── Assessment (Assessment)
│       └── Questions (Preguntas) [MongoDB]
└── Users (Usuarios)
    └── Progress (Progreso)
```

---

**Última actualización:** Diciembre 2024

# Diseno de Usuarios y Roles — EduGo Seed v2

**Fecha:** 2026-03-07
**Objetivo:** Redisenar desde cero los datos de prueba para cubrir escenarios multi-escuela,
multi-rol, flujo completo de evaluaciones (crear + tomar + calificar), y terminologia
por tipo de institucion.

---

## 1. Instituciones de Prueba

| ID | Nombre | Tipo Concepto | Ciudad | Suscripcion |
|----|--------|---------------|--------|-------------|
| SCH-01 | Colegio San Ignacio | high_school | Santiago | premium |
| SCH-02 | Taller CreArte | workshop | Valparaiso | basic |
| SCH-03 | Academia Global English | language_academy | Santiago | basic |

### Por que estas 3?
- **Colegio** — escenario clasico: grados, secciones, materias, evaluaciones formales
- **Taller** — terminologia diferente: modulos, grupos, facilitadores, ejercicios
- **Academia** — terminologia en ingles/mixto: levels, classes, teachers, tests

---

## 2. Terminologia por Tipo de Institucion

Ver documento separado: `02-TERMINOLOGIA.md`

---

## 3. Jerarquia Academica (Academic Units)

### Colegio San Ignacio (SCH-01)
```
Colegio San Ignacio (type=school)
  +-- 5to Basico (type=grade, level=5, year=2026)
  |     +-- 5to A (type=class)
  |     +-- 5to B (type=class)
  +-- 6to Basico (type=grade, level=6, year=2026)
        +-- 6to A (type=class)
```

### Taller CreArte (SCH-02)
```
Taller CreArte (type=school)
  +-- Modulo Pintura (type=grade)
  |     +-- Grupo Manana (type=class)
  +-- Modulo Escultura (type=grade)
        +-- Grupo Tarde (type=class)
```

### Academia Global English (SCH-03)
```
Academia Global English (type=school)
  +-- Level A2 (type=grade, level=A2)
  |     +-- Class Monday (type=class)
  +-- Level B1 (type=grade, level=B1)
        +-- Class Tuesday (type=class)
```

---

## 4. Usuarios de Prueba

**Password unificada:** `12345678`
**Hash bcrypt (cost=10):** `$2a$10$w9EyJdpR0T0leuTr9rso4O5xnOPdnVmVnkowe3MRJPEr94sRytzau`

### 4.1 Super Admin (acceso total)

| ID | Email | Nombre | Rol | Escuelas |
|----|-------|--------|-----|----------|
| U-01 | super@edugo.test | Santiago Ramirez | super_admin | TODAS (sin scope) |

**Permisos:** Todos (65). Ve el menu completo. Puede gestionar cualquier escuela.

---

### 4.2 Administradores de Escuela

| ID | Email | Nombre | Rol | Escuela |
|----|-------|--------|-----|---------|
| U-02 | admin.sanignacio@edugo.test | Carmen Valdes | school_admin | Colegio San Ignacio |
| U-03 | admin.crearte@edugo.test | Roberto Silva | school_admin | Taller CreArte |

**Permisos:** 26+ (control total de su institucion).
**Pantallas:** Dashboard admin, usuarios, unidades, memberships, materias, evaluaciones, materiales, progreso, estadisticas, guardian relations.

---

### 4.3 Coordinador Multi-Escuela

| ID | Email | Nombre | Roles | Escuelas |
|----|-------|--------|-------|----------|
| U-04 | coord.academico@edugo.test | Lucia Fernandez | school_coordinator (x2) | San Ignacio + CreArte |

**Escenario:** Coordinadora que trabaja en ambas instituciones. Cuando entra, elige contexto (switch-context). Ve pantallas diferentes segun la terminologia de cada institucion.

---

### 4.4 Profesores / Facilitadores

| ID | Email | Nombre | Roles | Escuelas | Unidades |
|----|-------|--------|-------|----------|----------|
| U-05 | prof.martinez@edugo.test | Maria Martinez | teacher (x2) | San Ignacio + Academia | 5to A (SI) + Class Monday (AcGE) |
| U-06 | prof.gonzalez@edugo.test | Pedro Gonzalez | teacher | San Ignacio | 5to B + 6to A |
| U-07 | facilitador.ruiz@edugo.test | Ana Ruiz | teacher | Taller CreArte | Grupo Manana + Grupo Tarde |

**Maria Martinez (U-05) — Caso multi-escuela:**
- En San Ignacio: es "Profesora de Matematicas" en 5to A
- En Academia Global: es "Teacher" en Class Monday (Level A2)
- Cuando cambia contexto, la terminologia del menu cambia

**Pedro Gonzalez (U-06) — Multi-clase:**
- En San Ignacio: ensena en 5to B Y 6to A (multiples memberships)

**Ana Ruiz (U-07) — Terminologia taller:**
- En CreArte: es "Facilitadora" (no "Profesora"), da "Talleres" (no "Materias")

**Permisos profesor:** materials:read/create/download, assessments:read/create/update/delete/publish, progress:read, subjects:read, assessments:grade, dashboard:view

---

### 4.5 Estudiantes / Participantes

| ID | Email | Nombre | Roles | Escuelas | Unidades |
|----|-------|--------|-------|----------|----------|
| U-08 | est.carlos@edugo.test | Carlos Mendoza | student (x2) | San Ignacio + CreArte | 5to A (SI) + Grupo Manana (CA) |
| U-09 | est.sofia@edugo.test | Sofia Herrera | student | San Ignacio | 5to A |
| U-10 | est.diego@edugo.test | Diego Vargas | student | San Ignacio | 5to B |
| U-11 | est.valentina@edugo.test | Valentina Rojas | student (x2) | San Ignacio + Academia | 6to A (SI) + Class Monday (AcGE) |
| U-12 | est.mateo@edugo.test | Mateo Fuentes | student | Taller CreArte | Grupo Manana |

**Carlos Mendoza (U-08) — Caso multi-escuela estudiante:**
- En San Ignacio: "Alumno" de 5to A, ve "Materias" y "Evaluaciones"
- En CreArte: "Participante" del Grupo Manana (Modulo Pintura), ve "Talleres" y "Ejercicios"
- Al hacer switch-context, la terminologia cambia completa

**Sofia Herrera (U-09) — Estudiante simple:**
- Solo en San Ignacio 5to A. Presenta evaluaciones creadas por Maria (U-05).

**Valentina Rojas (U-11) — Multi-escuela (colegio + academia):**
- En San Ignacio: "Alumna" de 6to A
- En Academia: "Student" de Class Monday, Level A2

**Permisos estudiante:** assessments:attempt, assessments:read, assessments:view_results, materials:read, materials:download, progress:read:own, progress:update, screens:read, dashboard:view

---

### 4.6 Tutores / Apoderados

| ID | Email | Nombre | Rol | Escuela | Hijos |
|----|-------|--------|-----|---------|-------|
| U-13 | tutor.mendoza@edugo.test | Ricardo Mendoza | guardian | San Ignacio + CreArte | Carlos (U-08) |
| U-14 | tutora.herrera@edugo.test | Patricia Herrera | guardian | San Ignacio | Sofia (U-09), Diego (U-10) |

**Ricardo Mendoza (U-13):**
- Padre de Carlos. Como Carlos esta en 2 escuelas, Ricardo tiene visibilidad en ambas.
- Ve progreso de Carlos en San Ignacio Y en CreArte.

**Patricia Herrera (U-14):**
- Madre de Sofia y tutora de Diego. Ambos en San Ignacio pero en clases diferentes.

---

## 5. Materias por Unidad

### Colegio San Ignacio
| Materia | Unidad | Profesor |
|---------|--------|----------|
| Matematicas | 5to A | Maria Martinez (U-05) |
| Ciencias Naturales | 5to A | Pedro Gonzalez (U-06) |
| Matematicas | 5to B | Pedro Gonzalez (U-06) |
| Historia | 6to A | Pedro Gonzalez (U-06) |

### Taller CreArte
| Taller | Grupo | Facilitador |
|--------|-------|-------------|
| Tecnicas de Pintura | Grupo Manana | Ana Ruiz (U-07) |
| Fundamentos de Escultura | Grupo Tarde | Ana Ruiz (U-07) |

### Academia Global English
| Course | Class | Teacher |
|--------|-------|---------|
| English Basics A2 | Class Monday | Maria Martinez (U-05) |

---

## 6. Evaluaciones de Prueba

### 6.1 Evaluaciones PUBLICADAS (con preguntas en MongoDB)

| ID | Titulo | Escuela | Profesor | Preguntas | Pass% | Status | Timed |
|----|--------|---------|----------|-----------|-------|--------|-------|
| ASS-01 | Examen Fracciones | San Ignacio | Maria (U-05) | 5 | 60% | published | Si (30min) |
| ASS-02 | Quiz Ciencias: Sistema Solar | San Ignacio | Pedro (U-06) | 4 | 50% | published | No |
| ASS-03 | Ejercicio Color y Forma | CreArte | Ana (U-07) | 3 | 70% | published | No |
| ASS-04 | English Grammar Test | Academia | Maria (U-05) | 4 | 60% | published | Si (20min) |

### 6.2 Evaluaciones en BORRADOR (sin publicar)

| ID | Titulo | Escuela | Profesor | Preguntas | Status |
|----|--------|---------|----------|-----------|--------|
| ASS-05 | Evaluacion Historia Chile | San Ignacio | Pedro (U-06) | 3 | draft |
| ASS-06 | Proyecto Final Escultura | CreArte | Ana (U-07) | 0 | draft |

### 6.3 Intentos de Estudiantes (Attempts)

| Estudiante | Evaluacion | Score | Max | % | Intento# | Status |
|-----------|-----------|-------|-----|---|----------|--------|
| Carlos (U-08) | Fracciones (ASS-01) | 80 | 100 | 80% | 1 | completed |
| Carlos (U-08) | Fracciones (ASS-01) | 92 | 100 | 92% | 2 | completed |
| Sofia (U-09) | Fracciones (ASS-01) | 68 | 100 | 68% | 1 | completed |
| Diego (U-10) | Ciencias (ASS-02) | 75 | 80 | 93.75% | 1 | completed |
| Carlos (U-08) | Color y Forma (ASS-03) | 60 | 100 | 60% | 1 | completed |
| Mateo (U-12) | Color y Forma (ASS-03) | 90 | 100 | 90% | 1 | completed |
| Valentina (U-11) | English Grammar (ASS-04) | 85 | 100 | 85% | 1 | completed |

**Nota:** Carlos tiene intentos en 2 escuelas diferentes (San Ignacio + CreArte).

---

## 7. Materiales Educativos

| ID | Titulo | Escuela | Profesor | Unidad | Status |
|----|--------|---------|----------|--------|--------|
| MAT-01 | Introduccion a las Fracciones | San Ignacio | Maria (U-05) | 5to A | ready |
| MAT-02 | El Sistema Solar | San Ignacio | Pedro (U-06) | 5to A | ready |
| MAT-03 | Historia de Chile: Independencia | San Ignacio | Pedro (U-06) | 6to A | ready |
| MAT-04 | Teoria del Color | CreArte | Ana (U-07) | Grp Manana | ready |
| MAT-05 | English Grammar Basics | Academia | Maria (U-05) | Class Monday | ready |

### Relacion Assessment <-> Material
| Assessment | Materiales |
|-----------|-----------|
| ASS-01 (Fracciones) | MAT-01 |
| ASS-02 (Ciencias) | MAT-02 |
| ASS-03 (Color y Forma) | MAT-04 |
| ASS-04 (English Grammar) | MAT-05 |
| ASS-05 (Historia) | MAT-03 |
| ASS-06 (Escultura) | (ninguno) |

---

## 8. Guardian Relations

| Tutor | Estudiante | Tipo | Primary | Status |
|-------|-----------|------|---------|--------|
| Ricardo (U-13) | Carlos (U-08) | parent | Si | active |
| Patricia (U-14) | Sofia (U-09) | parent | Si | active |
| Patricia (U-14) | Diego (U-10) | guardian | No | active |

---

## 9. Resumen de Pantallas por Rol

### Super Admin (U-01)
- Dashboard Superadmin
- Escuelas (CRUD), Usuarios (CRUD), Roles (CRUD), Permisos (CRUD)
- Templates/Instancias de pantalla, Tipos de concepto
- Auditoria
- Todo lo que tiene school_admin + mas

### School Admin (U-02, U-03)
- Dashboard Admin
- Usuarios, Unidades, Memberships, Materias
- Materiales, Evaluaciones (ver), Guardian Relations
- Progreso, Estadisticas
- Configuracion

### Coordinador (U-04)
- Similar a School Admin pero sin gestion de usuarios/escuelas
- Enfocado en contenido y academico

### Profesor/Facilitador (U-05, U-06, U-07)
- Dashboard Profesor
- Materiales (CRUD), Evaluaciones (CRUD + publicar/archivar)
- Preguntas de evaluacion (CRUD)
- Progreso de estudiantes, Estadisticas
- Materias (solo lectura)

### Estudiante/Participante (U-08 a U-12)
- Dashboard Estudiante
- Materiales (solo lectura + descarga)
- Evaluaciones disponibles (ver publicadas)
- Tomar evaluacion, Ver resultados
- Mi progreso

### Tutor/Apoderado (U-13, U-14)
- Dashboard Guardian
- Mis hijos (lista)
- Progreso por hijo
- Solicitudes de vinculacion

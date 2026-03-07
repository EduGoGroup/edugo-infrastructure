# Terminologia por Tipo de Institucion — EduGo

Cada tipo de institucion usa terminos diferentes para los mismos conceptos.
El sistema de `concept_types` + `concept_definitions` + `school_concepts` permite
que cada escuela tenga su propia terminologia.

---

## Tabla Comparativa Completa

| term_key | Colegio (high_school) | Taller (workshop) | Academia Idiomas (language_academy) | Instituto Tecnico (technical) | Escuela Primaria (primary) |
|----------|----------------------|-------------------|-------------------------------------|------------------------------|---------------------------|
| **ORGANIZACION** |
| org.name_singular | Colegio | Taller | Academia | Instituto | Escuela |
| org.name_plural | Colegios | Talleres | Academias | Institutos | Escuelas |
| **JERARQUIA** |
| unit.level1 | Ano | Modulo | Level | Semestre | Grado |
| unit.level1_plural | Anos | Modulos | Levels | Semestres | Grados |
| unit.level2 | Division | Grupo | Class | Seccion | Clase |
| unit.level2_plural | Divisiones | Grupos | Classes | Secciones | Clases |
| unit.period | Trimestre | Ciclo | Term | Cuatrimestre | Periodo |
| unit.period_plural | Trimestres | Ciclos | Terms | Cuatrimestres | Periodos |
| **MIEMBROS** |
| member.student | Alumno | Participante | Student | Aprendiz | Estudiante |
| member.student_plural | Alumnos | Participantes | Students | Aprendices | Estudiantes |
| member.teacher | Docente | Facilitador | Teacher | Instructor | Profesor |
| member.teacher_plural | Docentes | Facilitadores | Teachers | Instructores | Profesores |
| member.guardian | Tutor | Responsable | Parent | Representante | Acudiente |
| member.guardian_plural | Tutores | Responsables | Parents | Representantes | Acudientes |
| **CONTENIDO** |
| content.subject | Asignatura | Taller | Course | Modulo | Materia |
| content.subject_plural | Asignaturas | Talleres | Courses | Modulos | Materias |
| content.assessment | Examen | Ejercicio | Test | Prueba | Evaluacion |
| content.assessment_plural | Examenes | Ejercicios | Tests | Pruebas | Evaluaciones |
| content.material | Material | Recurso | Resource | Material | Material |
| content.material_plural | Materiales | Recursos | Resources | Materiales | Materiales |

---

## Impacto en la UI

Cuando un usuario entra a una institucion:
1. El frontend carga los `school_concepts` de la escuela activa
2. Los labels de la UI se reemplazan dinamicamente:
   - Menu lateral: "Materias" -> "Talleres" (en CreArte)
   - Encabezados: "Crear Evaluacion" -> "Crear Ejercicio"
   - Empty states: "No hay evaluaciones" -> "No hay ejercicios"
   - Filtros: "Profesores" -> "Facilitadores"

### Ejemplo: Carlos Mendoza (multi-escuela)

**En Colegio San Ignacio:**
```
Menu: Asignaturas | Examenes | Materiales
Dashboard: "Bienvenido, Alumno"
Lista: "Tus Examenes"
```

**En Taller CreArte:**
```
Menu: Talleres | Ejercicios | Recursos
Dashboard: "Bienvenido, Participante"
Lista: "Tus Ejercicios"
```

---

## Categorias de Terminos

| category | Descripcion | Ejemplo term_keys |
|----------|-------------|-------------------|
| org | Nombre de la organizacion | org.name_singular, org.name_plural |
| unit | Niveles de jerarquia academica | unit.level1, unit.level2, unit.period |
| member | Roles de personas | member.student, member.teacher, member.guardian |
| content | Tipos de contenido educativo | content.subject, content.assessment, content.material |

---

## Datos en BD

### Tabla: academic.concept_types
Define los TIPOS de institucion (plantillas base).

### Tabla: academic.concept_definitions
Define los TERMINOS por defecto para cada tipo (plantilla).

### Tabla: academic.school_concepts
Copia personalizable por escuela. Cuando se crea una escuela y se le asigna un concept_type,
se copian las definitions a school_concepts. La escuela puede luego personalizar sus terminos.

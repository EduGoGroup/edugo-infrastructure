# Especificación MongoDB - edugo-infrastructure

## Colecciones (3):

1. material_summary - Resúmenes IA
2. material_assessment - Quizzes IA
3. material_event - Log eventos

## Implementar en Sprint-03:

Ver: docs/isolated/04-Implementation/Sprint-03-MongoDB-Migrations/

### Archivos a crear:
- database/migrations/mongodb/001_create_material_summary.up.js
- database/migrations/mongodb/002_create_material_assessment.up.js
- database/migrations/mongodb/003_create_material_event.up.js
- Archivos .down.js correspondientes
- database/MONGODB_SCHEMA.md

### Campos principales por colección:

**material_summary:**
- material_id (UUID, unique)
- summary: {short, detailed, key_points[]}
- metadata: {word_count, difficulty_level, language}
- generated_at, model_version

**material_assessment:**
- material_id (UUID, unique)
- title, questions[], total_points
- questions: {question_id, text, type, points, options[], correct_answer}
- is_published, generated_at

**material_event:**
- material_id (UUID)
- event_type, status, message
- occurred_at
- TTL: 90 días

Ver /Analisys/00-Projects-Isolated/infrastructure/02-Requirements/MONGODB_SCHEMA.md para detalles completos.

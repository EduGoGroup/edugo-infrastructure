// Package common agrupa los helpers idempotentes compartidos por los
// playgrounds de la línea v2 (los que asumen que L0..L4 ya corrieron y reusan
// los roles L4 school_admin/teacher/student).
//
// Antes de este paquete, cada playground v2 redefinía ~9 funciones upsertXXX()
// casi idénticas (upsertSchool, upsertAcademicUnit, upsertSubject, upsertUser,
// upsertUserRole, upsertMembership, upsertSchoolMembership/upsertUnitMembership,
// upsertActivePeriod, upsertOffering, upsertEnrollment). Aquí viven las versiones
// canónicas, exportadas y parametrizadas, sin special-casing por playground.
//
// Estilo: mismo patrón que postgres/seeds/playground/common (v1, sobre L0) —
// structs *Spec con defaults y un SeedXxx(tx, spec) idempotente por PK. La
// diferencia con v1 es el modelo: v2 cubre el dominio academic completo
// (unidades, períodos, materias, ofertas, inscripciones), que v1 no tenía.
//
// Diferencias clave que motivan parametrizar TODO (en vez de leer constantes
// del paquete del playground, como hacían las copias):
//   - school/unit/period dejan de ser constantes implícitas; varios playgrounds
//     tienen N escuelas (n0n1_escuelas) o N unidades (multi_unidad).
//   - Membership cubre alcance COLEGIO (AcademicUnitID nil) y UNIDAD (no-nil)
//     con un único helper: el *uuid.UUID opcional decide el alcance.
//   - Offering acepta SectionLabel *string opcional (n1_inscripcion no la usa;
//     n4/n17 sí) y TeacherMembershipID *uuid.UUID opcional (oferta sin docente).
//
// Idempotencia: todos los Seed son OnConflict DoNothing por id (o por la PK
// compuesta en enrollment), idéntico a las copias originales.
package common

// Package fixtures contiene las piezas atómicas componibles que
// consumen los scenarios del paquete sibling scenarios/.
//
// Convenciones obligatorias para cada fixture nueva:
//
//   - Implementa la interface framework.Fixture (Manifest/Apply/Cleanup).
//   - Declara explícitamente Provides/Requires; el composer las resuelve
//     por orden topológico antes de tocar la base de datos.
//   - Todos los UUIDs que la fixture genere derivan del ApplyContext
//     (TenantPrefix/SchemaPrefix); cualquier intento de escribir en el
//     namespace del production seed (10000000-..., c1000000-...) debe
//     fallar a través de framework.AssertNotProductionNamespace.
//   - Idempotencia: aplicar la fixture N veces produce el mismo estado
//     final que aplicarla una sola vez (clause.OnConflict + UpsertBool).
//   - Para booleanos críticos usar framework.UpsertBool (lección F2·H5).
//   - Las constantes que el test consumirá se exportan a través del
//     campo Constants del FixtureManifest.
package fixtures

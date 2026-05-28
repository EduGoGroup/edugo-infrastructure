// Package scenarios reúne las recetas canónicas que combinan fixtures
// para reproducir casos de prueba focalizados (un rol, una pantalla,
// un sub-árbol del menú).
//
// Cada scenario:
//
//   - Implementa framework.Scenario (Manifest/BuildFixtures).
//   - Se registra en framework.RegisterScenario(...) desde la función
//     RegisterAll() del paquete (llamada una sola vez por el binario
//     seed_e2e o por los tests).
//   - Deriva su namespace del nombre del scenario vía
//     framework.Derive (hash determinístico).
//   - Documenta sus tags (rbac, menu, screen-config) para que la suite
//     Fase D pueda agruparlos.
//
// La fuente de verdad para los identificadores que necesitan los tests
// Kotlin es el JSON exportado por framework.ConstantsExport (ver
// seeds/e2e/exports/fixtures-constants.json).
package scenarios

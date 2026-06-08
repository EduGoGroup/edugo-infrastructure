package layers

// Poda SDUI material (2026-06-07): las ScreenInstances L3 `materials-list`
// y `material-form` (con sus slot_data canónicos) fueron ELIMINADAS.
//
// Razón — principio "nativa prevalece, SDUI solo guía mínima": las
// pantallas de material de la app KMP son NATIVAS (Compose) y NO consumen
// estos seeds SDUI. Los slot_data eran código muerto: nadie los renderiza.
// Al podarlos desaparece también el `api_prefix:"academic"` que vivía en
// el slot_data de `materials-list` (erróneo para una pantalla nativa).
//
// Lo que SE CONSERVA del recurso materials (no es código muerto):
//   - resource `materials` (l3_resources.go) — lo necesita el menú.
//   - permisos content.materials.* (l3_permissions.go + L4) — gate de menú.
//   - mapping `materials:list` (is_default) en l3_resource_screens.go —
//     ata el item de menú al screen_key nativo. Ahora vive SIN
//     ScreenInstance, igual que `material-detail` y las pantallas nativas
//     `join-requests-inbox` / `batch-enroll` / `enroll-one`: el resolver
//     solo necesita que el menú exponga el screen_key; el FE lo
//     intercepta con un composable nativo.
//
// Por eso L3 ya no siembra screen_instances: `applyL3Screens` fue
// eliminada y `l3.go::Apply` ya no la invoca.
//
// `material-detail` nunca tuvo ScreenInstance (dead desde Fase 6, ver
// system/l4/screen_instances.go). Este cambio alinea `materials-list` /
// `material-form` con ese mismo patrón.

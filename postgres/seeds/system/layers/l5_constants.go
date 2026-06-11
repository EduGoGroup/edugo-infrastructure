package layers

// L5_SEED_VERSION declara la versión semántica del contenido de L5.
// Bumpear en CADA cambio de dato visible en seeds/system/l5/.
//
// Historial:
//   - 1.0.0: clientes M2M iniciales (edugo-worker, edugo-api-learning)
//     con scope notifications.dispatch (plan 020 N5 F1.2).
const L5_SEED_VERSION = "1.0.0"

// L5_LAYER_NAME es el nombre canónico de la capa, usado por
// --seed-up-to-layer y por logs.
const L5_LAYER_NAME = "L5-m2m"

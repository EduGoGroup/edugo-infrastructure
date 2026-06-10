// Package l5 siembra clientes M2M para autenticación service JWT (plan 020 N5).
//
// Filas iniciales: edugo-worker y edugo-api-learning con scope
// notifications.dispatch. El secret en claro vive en push-secrets.env
// (desarrollo) o Secret Manager (cloud); en BD solo secret_hash.
package l5

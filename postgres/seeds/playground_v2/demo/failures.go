package demo

import (
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"gorm.io/gorm"
)

// Las filas de lo fallado del mundo de demo (HU-IA.7, decisión 10 de su ficha):
// tres, sobre intentos de E-R que ya existen, para que «Fallos de la IA» de la
// docente enseñe algo en cada estado que importa sin mover los conteos del
// fixture. Tener fila es además lo que impide que el barrido de learning tome
// esos intentos por huérfanos.
//
//   - Intento 1: `failed` de servicio, con los 5 reintentos agotados.
//   - Intento 2: `failed` de dato: la IA devolvió algo que no se pudo leer.
//   - Intento 3: `retrying` de servicio, con el próximo reintento en una hora.
//
// El cuarto intento queda sin fila. Las `failed` son inertes; la `retrying` (y
// el cuarto) se moverían si un worker corre contra demo: es lo esperado, y
// recrear demo los devuelve a su sitio. Corregir el intento en la revisión
// resuelve su fila, así que la escena de E-R sigue valiendo.

// AIFailureCount es el número de filas de assessment.ai_failure de demo.
const AIFailureCount = 3

// AIFailureID es el id de la fila n (1-based) de lo fallado.
func AIFailureID(n int) string { return newID(domAIFailure, n).String() }

// ReviewAttemptID es el id del intento i (1-based) de E-R.
func ReviewAttemptID(i int) string { return newID(domAttempt, i).String() }

// buildFailures arma las tres filas a partir del momento de aplicar. Es pura.
func buildFailures(w *world) []entities.AIFailure {
	created := w.day(-1).Add(16 * time.Hour)
	next := w.now.Add(time.Hour)
	return []entities.AIFailure{
		{
			ID: newID(domAIFailure, 1), TargetType: "attempt", TargetID: newID(domAttempt, 1),
			Status: "failed", CauseClass: "service", RetryCount: 5,
			LastError: ptr("llm: el servicio no respondió a tiempo (timeout tras 5 reintentos)"),
			CreatedAt: created, UpdatedAt: created.Add(5 * time.Hour),
		},
		{
			ID: newID(domAIFailure, 2), TargetType: "attempt", TargetID: newID(domAttempt, 2),
			Status: "failed", CauseClass: "data", RetryCount: 0,
			LastError: ptr("llm: la respuesta no es un JSON de revisión válido"),
			CreatedAt: created, UpdatedAt: created,
		},
		{
			ID: newID(domAIFailure, 3), TargetType: "attempt", TargetID: newID(domAttempt, 3),
			Status: "retrying", CauseClass: "service", RetryCount: 2, NextRetryAt: ptr(next),
			LastError: ptr("llm: el servicio respondió 503"),
			CreatedAt: created, UpdatedAt: created.Add(2 * time.Hour),
		},
	}
}

func writeFailures(tx *gorm.DB, w *world) error {
	rows := buildFailures(w)
	if err := insert(tx, &rows); err != nil {
		return fmt.Errorf("ai_failure: %w", err)
	}
	return nil
}

package demo

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// buildGrants arma el recorte del contrato que la P4 de R-026 permite en la
// base de demo, nunca en system/ (veredicto V5, firmado en R-028): filas `deny`
// en iam.user_grants para dos cuentas, sin tocar iam.role_grants.
//
//   - A la docente, messaging.*: quita «Mensajería» de su menú, que en
//     escritorio y web abre un cartel.
//   - A la tutora, los dos permisos de materiales: quita «Materiales» de su
//     menú, que sale siempre vacío (la lista de su acudida usa otro permiso).
//
// Son tres filas, no dos: la clave única es (usuario, patrón, efecto). Sin los
// `allow` de la V6. Que el menú respete un `deny` con comodín lo comprueba la
// sonda del arnés en /menu.
func buildGrants(w *world) {
	teacher := uuid.MustParse(TeacherUserID)
	guardian := uuid.MustParse(GuardianUserID)
	deny := func(n int, user uuid.UUID, pattern string) entities.UserGrant {
		return entities.UserGrant{
			ID: newID(domUserGrant, n), UserID: user, PermissionPattern: pattern,
			Effect: "deny", CreatedAt: w.day(periodStartDays),
		}
	}
	w.grants = []entities.UserGrant{
		deny(1, teacher, "messaging.*"),
		deny(2, guardian, "content.materials.read"),
		deny(3, guardian, "content.materials.download"),
	}
}

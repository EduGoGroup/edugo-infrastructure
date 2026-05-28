package scenarios

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/e2e/framework"
)

// RegisterAll registra los scenarios canónicos definidos en este
// paquete en el registry indicado. Si reg es nil se utiliza
// framework.DefaultRegistry.
//
// Scenarios actuales:
//   - observer_audits, teacher_grades_only, guardian_views_child:
//     casos focalizados de roles del producto.
//   - l0_minimal: valida la capa L0 del seed system (Fase 2 rebuild).
//   - l1_readonly: valida la capa L1 del seed system (Fase 3 rebuild).
//   - l2_form: valida la capa L2 del seed system (Fase 4 rebuild).
//   - l3_isolation: valida la capa L3 del seed system (Fase 5 rebuild).
//   - l4_full: valida la capa L4 del seed system completa + matriz
//     screens-por-rol derivada (Fase 6 rebuild).
//
// Nota (ADR-6): LegacyE2E se retiró en Fase 2 — el legacy ya no se
// aplica en runtime. Las fixtures legacy_* fueron eliminadas.
//
// La función falla en la primera duplicación detectada (errores
// "duplicate scenario"), de modo que es seguro invocarla una única vez
// por proceso. Los tests deben pasar un Registry fresco (NewRegistry)
// para evitar contaminación cruzada entre suites.
func RegisterAll(reg *framework.Registry) error {
	if reg == nil {
		reg = framework.DefaultRegistry
	}
	scenarios := []framework.Scenario{
		&ObserverAudits{},
		&TeacherGradesOnly{},
		&GuardianViewsChild{},
		&L0Minimal{},
		&L1ReadOnly{},
		&L2Form{},
		&L3Isolation{},
		&L4Full{},
		&SuperAdminGlobalFlow{},
	}
	for _, s := range scenarios {
		if err := reg.RegisterScenario(s); err != nil {
			return err
		}
	}
	return nil
}

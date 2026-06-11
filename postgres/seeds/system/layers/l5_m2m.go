// Capa L5 — Clientes M2M (service JWT, plan 020 N5)
// ===================================================
//
// Siembra auth.service_clients con los callers autorizados a invocar
// POST /api/v1/internal/notifications/dispatch. Depende de que la tabla
// exista (entity ServiceClient en AutoMigrate, SchemaVersion bump).
//
// Refs: docs/plans/020-n5-push-notificaciones/ (D15, D16).
package layers

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/l5"
	"gorm.io/gorm"
)

type l5Layer struct{}

// NewL5 construye una instancia de la capa L5.
// Se registra en system.Layers() tras NewL4().
func NewL5() *l5Layer { return &l5Layer{} }

func (l *l5Layer) Name() string { return L5_LAYER_NAME }

func (l *l5Layer) SeedVersion() string { return L5_SEED_VERSION }

func (l *l5Layer) Apply(tx *gorm.DB) error {
	return l5.ApplyServiceClients(tx)
}

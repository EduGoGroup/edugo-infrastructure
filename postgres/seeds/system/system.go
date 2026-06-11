package system

import (
	"database/sql"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system/layers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Layers retorna la lista ordenada de capas del seed system.
//
// Estado Fase 6 (en progreso):
//   - L0-minimal: dataset mínimo viable (17 filas).
//   - L1-readonly: rol viewer + escuela mínima para validar gating
//     de UI (6 filas adicionales, ADR-7; incluye membership). Acumulado: 23 filas.
//   - L2: segunda pantalla (announcement-form) + mapping form
//     (2 filas adicionales). Acumulado: 25 filas.
//   - L3: segundo recurso (materials) con CRUD parcial sin delete +
//     2 pantallas + 2 mappings (11 filas adicionales). Acumulado: 36 filas.
//   - L4-full: sistema completo reorganizado por dominio (resources,
//     roles_permissions, screen_templates, screen_instances,
//     resource_screens, concept_types). Implementación por batches
//     B1..B7 (phase-6-layer-l4/tasks.md). Stub en B0; datos reales
//     en seeds/system/l4/*.go.
//   - L5-m2m: clientes service JWT (edugo-worker, edugo-api-learning)
//     con scope notifications.dispatch (plan 020 N5).
//   - Layer_Legacy ([archivado pre-Fase-6] ) NO se aplica desde Fase 2
//     (ADR-6); el directorio se elimina del disco en el bloque C
//     de Fase 6.
func Layers() []Layer {
	return []Layer{
		layers.NewL0(),
		layers.NewL1(),
		layers.NewL2(),
		layers.NewL3(),
		layers.NewL4(),
		layers.NewL5(),
	}
}

// ApplySystem ejecuta las capas en orden hasta upTo (vacío = todas).
// Cada capa corre en su propia transacción.
func ApplySystem(db *sql.DB, upTo string) error {
	gdb, err := openGORM(db)
	if err != nil {
		return err
	}
	return ApplySystemGORM(gdb, upTo)
}

func ApplySystemGORM(gdb *gorm.DB, upTo string) error {
	list := Layers()
	target := resolveTargetIndex(list, upTo)
	if target < 0 {
		return fmt.Errorf("system.ApplySystem: layer %q not found", upTo)
	}
	for i := 0; i <= target; i++ {
		l := list[i]
		if err := l.Apply(gdb); err != nil {
			return &LayerError{Layer: l.Name(), Err: err}
		}
	}
	return nil
}

// resolveTargetIndex devuelve el índice hasta el cual aplicar.
// upTo == "" → última capa. Si upTo no se encuentra, retorna -1.
func resolveTargetIndex(list []Layer, upTo string) int {
	if upTo == "" {
		return len(list) - 1
	}
	for i, l := range list {
		if l.Name() == upTo {
			return i
		}
	}
	return -1
}

func openGORM(db *sql.DB) (*gorm.DB, error) {
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("system.openGORM: %w", err)
	}
	return gdb, nil
}

package system

import "gorm.io/gorm"

// Layer es la interfaz que toda capa del seed system debe implementar.
// Implementadores viven en seeds/system/<capa>/ (ej. legacy, l0, l1, ...).
type Layer interface {
	Name() string
	SeedVersion() string
	Apply(tx *gorm.DB) error
}

// LayerError envuelve un error originado en una capa específica.
type LayerError struct {
	Layer string
	Err   error
}

func (e *LayerError) Error() string {
	return "layer " + e.Layer + ": " + e.Err.Error()
}

func (e *LayerError) Unwrap() error { return e.Err }

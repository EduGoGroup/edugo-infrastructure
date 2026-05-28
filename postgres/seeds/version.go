package seeds

import (
	"crypto/sha256"
	"fmt"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/demo"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/system"
)

const DemoSeedHashVersion = demo.SeedVersion

// ComputeFilesHash itera system.Layers() para construir el hash dinámicamente.
// Cualquier nueva capa añadida a Layers() se incluye automáticamente.
func ComputeFilesHash() string {
	h := sha256.New()
	for _, l := range system.Layers() {
		h.Write([]byte(l.Name() + ":" + l.SeedVersion() + "\n"))
	}
	h.Write([]byte("demo:" + DemoSeedHashVersion + "\n"))
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

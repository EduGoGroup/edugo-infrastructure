package migrations

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"io/fs"
	"sort"
)

// SchemaVersion es la version actual de los scripts de migracion y seeds.
//
// OBLIGATORIO: Incrementar este valor cada vez que se modifique
// cualquier archivo en structure/*.sql o en seeds/.
// El migrador valida que esta version coincida con la registrada en BD.
const SchemaVersion = "1.1.4"

// ComputeFilesHash calcula un SHA256 de todos los archivos SQL embebidos
// en el paquete migrations. El hash cambia si cualquier archivo se modifica.
func ComputeFilesHash() string {
	return ComputeEmbedHash(Files, "structure")
}

// ComputeEmbedHash calcula el hash de un embed.FS para un directorio dado.
// Exportado para que el migrador pueda usarlo con seeds.Files tambien.
func ComputeEmbedHash(fsys embed.FS, dir string) string {
	h := sha256.New()

	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return "error"
	}

	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		content, err := fsys.ReadFile(dir + "/" + name)
		if err != nil {
			continue
		}
		h.Write([]byte(name))
		h.Write(content)
	}

	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

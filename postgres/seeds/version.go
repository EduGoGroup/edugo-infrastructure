package seeds

import (
	"crypto/sha256"
	"embed"
	"fmt"
	"hash"
	"io/fs"
	"sort"
)

// ComputeFilesHash calcula un SHA256 combinado de todos los archivos SQL
// de production/ y development/. El hash cambia si cualquier seed se modifica.
func ComputeFilesHash() string {
	h := sha256.New()
	hashDir(h, Files, "production")
	hashDir(h, Files, "development")
	return fmt.Sprintf("%x", h.Sum(nil))[:16]
}

func hashDir(h hash.Hash, fsys embed.FS, dir string) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return
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
		h.Write([]byte(dir + "/" + name))
		h.Write(content)
	}
}

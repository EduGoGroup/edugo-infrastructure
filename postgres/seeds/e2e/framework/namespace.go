package framework

import (
	"crypto/sha1" //nolint:gosec // SHA-1 truncado se usa como hash determinista de prefijo, no para criptografía.
	"encoding/hex"
	"fmt"
	"strings"
)

// LegacyScenarioName es el nombre reservado del scenario que reproduce
// bit-a-bit el seed E2E histórico (fase0..fase4). Su hash se fuerza a
// "00000000" para mantener paridad con los UUIDs e2e00000-... ya
// referenciados desde los .md del plan E2E previo (C-REQ-5.1, C-REQ-5.4).
const LegacyScenarioName = "legacy_e2e"

// LegacyHash es el hash forzado del scenario legacy. Tiene 5 caracteres
// hexadecimales para componer un primer segmento UUID canónico de 8
// caracteres ("e2e" + 5 hex). El sub-namespace fijo "00000" mantiene
// la paridad bit-a-bit con los UUIDs históricos del seed E2E vigente
// (ej. e2e00000-0000-0000-0000-000000000001).
const LegacyHash = "00000"

// scenarioHashLen es el largo del hash hexadecimal usado para derivar
// los prefijos de aislamiento. El valor está fijado para que "e2e" +
// hash forme exactamente el primer segmento UUID canónico (8 chars),
// permitiendo que las fixtures generen identificadores válidos contra
// columnas postgres con tipo `uuid`.
//
// Nota: design.md de Fase C menciona "4 bytes / 8 hex chars" pero ese
// patrón rompe el formato UUID (9 chars en el primer segmento). El
// framework toma los primeros 5 hex chars (~20 bits) como compromiso:
// suficiente entropía (~1M scenarios) sin desbordar el campo.
const scenarioHashLen = 5

// productionUUIDPrefixes enumera los rangos del production seed.
// Cualquier intento de escribir un UUID que empiece con uno de estos
// prefijos debe fallar con `forbidden namespace: cannot write into
// production seed range` (C-REQ-10.2).
var productionUUIDPrefixes = []string{
	"10000000-", // resources, roles, role_permissions, screen_*
	"c1000000-", // concept_types, concept_definitions
}

// developmentUUIDPrefix identifica el namespace de development/.
// Los scenarios E2E nunca deben escribir aquí.
const developmentUUIDPrefix = "00000000-"

// Derive calcula los prefijos de aislamiento para un scenario a partir
// de su nombre. El hash es SHA-1 truncado a `scenarioHashLen`
// caracteres hexadecimales.
//
// Ejemplo (hash 5 hex chars):
//
//	Derive("teacher_grades_only") -> ("E2E-A1B2C-", "e2ea1b2c-")
//	Derive("legacy_e2e")          -> ("E2E-",       "e2e00000-")
//
// El scenario legacy es un caso especial: el hash se fuerza a "00000"
// y el TenantPrefix se mantiene como "E2E-" sin segmento intermedio
// (C-REQ-5.1) para conservar paridad bit-a-bit con los códigos
// históricos del plan E2E previo (E2E-SCHOOL-01, etc.).
func Derive(scenarioName string) (tenantPrefix, schemaPrefix string) {
	if scenarioName == LegacyScenarioName {
		return "E2E-", "e2e00000-"
	}
	hash := scenarioHash(scenarioName)
	return "E2E-" + strings.ToUpper(hash) + "-", "e2e" + hash + "-"
}

// scenarioHash devuelve los primeros `scenarioHashLen` caracteres
// hexadecimales del SHA-1 del nombre del scenario.
func scenarioHash(scenarioName string) string {
	sum := sha1.Sum([]byte(scenarioName)) //nolint:gosec
	return hex.EncodeToString(sum[:])[:scenarioHashLen]
}

// AssertNotProductionNamespace falla si uuid cae en alguno de los
// rangos reservados al production seed o development. Las fixtures
// deben llamarla antes de cada INSERT con un UUID generado.
//
// Acepta tanto un UUID con guiones (ej. "10000000-0000-0000-0000-000000000001")
// como una string normalizada en minúscula. Compara por prefijo.
func AssertNotProductionNamespace(uuid string) error {
	if uuid == "" {
		return fmt.Errorf("forbidden namespace: empty UUID")
	}
	low := strings.ToLower(uuid)
	for _, p := range productionUUIDPrefixes {
		if strings.HasPrefix(low, p) {
			return fmt.Errorf("forbidden namespace: cannot write into production seed range (uuid=%s, range=%s...)", uuid, p)
		}
	}
	if strings.HasPrefix(low, developmentUUIDPrefix) {
		return fmt.Errorf("forbidden namespace: cannot write into development range (uuid=%s, range=%s...)", uuid, developmentUUIDPrefix)
	}
	return nil
}

// MakeUUID concatena el SchemaPrefix con un sufijo determinista. El
// sufijo se completa con ceros a la izquierda hasta los 36 chars
// canónicos de un UUID v4 textual (8-4-4-4-12).
//
// Ejemplo:
//
//	ctx.SchemaPrefix = "e2ea1b2c3d4-"
//	MakeUUID(ctx, "0000-0000-0000-000000000001")
//	  -> "e2ea1b2c3d4-0000-0000-0000-000000000001"
//
// La función NO valida la unicidad — los manifests son responsables de
// garantizar que dos fixtures del mismo scenario no compartan sufijos.
func MakeUUID(ctx *ApplyContext, suffix string) string {
	if ctx == nil {
		return suffix
	}
	return ctx.SchemaPrefix + suffix
}

// MakeCode produce un código visible con el TenantPrefix del scenario.
//
// Ejemplo:
//
//	MakeCode(ctx, "SCHOOL", "01") -> "E2E-A1B2C3D4-SCHOOL-01"
func MakeCode(ctx *ApplyContext, kind, index string) string {
	if ctx == nil {
		return kind + "-" + index
	}
	return ctx.TenantPrefix + kind + "-" + index
}

// MakeEmail aplica el patrón <role>-<fixture>-<hash>@edugo.test
// (C-REQ-2.4) usando el SchemaPrefix del scenario para extraer el hash
// (los 8 caracteres entre "e2e" y el primer guion del prefijo).
func MakeEmail(ctx *ApplyContext, role, fixture string) string {
	hash := schemaHashFromPrefix(ctx)
	return fmt.Sprintf("%s-%s-%s@edugo.test", role, fixture, hash)
}

// schemaHashFromPrefix extrae el hash de un SchemaPrefix derivado por
// Derive(). Sirve para componer emails, sufijos `ro_<hash>`, etc.
func schemaHashFromPrefix(ctx *ApplyContext) string {
	if ctx == nil {
		return LegacyHash
	}
	p := strings.TrimSuffix(strings.TrimPrefix(ctx.SchemaPrefix, "e2e"), "-")
	if p == "" {
		return LegacyHash
	}
	return p
}

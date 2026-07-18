package entities

import "testing"

func TestValidateSetting(t *testing.T) {
	cases := []struct {
		name    string
		key     string
		value   string
		wantErr bool
	}{
		{"enum válido llm.review.mode=api", SettingLLMReviewMode, "api", false},
		{"enum válido llm.review.mode=off", SettingLLMReviewMode, "off", false},
		{"enum inválido llm.review.mode=maybe", SettingLLMReviewMode, "maybe", true},
		{"enum válido llm.review.flow=teacher", SettingLLMReviewFlow, "teacher", false},
		{"enum válido llm.generation.mode=local", SettingLLMGenerationMode, "local", false},
		{"enum válido llm.pipeline.mode=on", SettingLLMPipelineMode, "on", false},
		{"enum válido llm.pipeline.mode=off", SettingLLMPipelineMode, "off", false},
		{"enum inválido llm.pipeline.mode=local", SettingLLMPipelineMode, "local", true},
		{"int válido import.max_questions=50", SettingImportMaxQuestions, "50", false},
		{"int con espacios import.max_questions", SettingImportMaxQuestions, "  200  ", false},
		{"int cero import.max_questions=0", SettingImportMaxQuestions, "0", true},
		{"int negativo import.max_questions=-1", SettingImportMaxQuestions, "-1", true},
		{"int no numérico import.max_json_bytes", SettingImportMaxJSONBytes, "abc", true},
		{"clave desconocida", "llm.unknown.key", "x", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateSetting(tc.key, tc.value)
			if tc.wantErr && err == nil {
				t.Fatalf("ValidateSetting(%q, %q) = nil, se esperaba error", tc.key, tc.value)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("ValidateSetting(%q, %q) = %v, se esperaba nil", tc.key, tc.value, err)
			}
		})
	}
}

func TestResolveDefaultHardDefaults(t *testing.T) {
	// Sin env seteada, cada clave cae a su default duro (D-039.2).
	want := map[string]string{
		SettingLLMGenerationMode:  "off",
		SettingLLMReviewMode:      "off",
		SettingLLMReviewFlow:      "teacher",
		SettingLLMPipelineMode:    "off",
		SettingImportMaxQuestions: "100",
		SettingImportMaxJSONBytes: "1048576",
	}
	for key, exp := range want {
		spec, ok := LookupSetting(key)
		if !ok {
			t.Fatalf("LookupSetting(%q) no encontró la clave en el catálogo", key)
		}
		// Aisla el test de un entorno con la env ya seteada.
		t.Setenv(spec.EnvDefault, "")
		if got := ResolveDefault(key); got != exp {
			t.Fatalf("ResolveDefault(%q) = %q, se esperaba el default duro %q", key, got, exp)
		}
	}
}

func TestResolveDefaultFromEnv(t *testing.T) {
	spec, _ := LookupSetting(SettingLLMReviewMode)
	t.Setenv(spec.EnvDefault, "api")
	if got := ResolveDefault(SettingLLMReviewMode); got != "api" {
		t.Fatalf("ResolveDefault con env válida = %q, se esperaba %q", got, "api")
	}

	// Una env inválida se ignora y cae al default duro.
	t.Setenv(spec.EnvDefault, "maybe")
	if got := ResolveDefault(SettingLLMReviewMode); got != spec.HardDefault {
		t.Fatalf("ResolveDefault con env inválida = %q, se esperaba el default duro %q", got, spec.HardDefault)
	}
}

func TestCatalogIsCopy(t *testing.T) {
	c := Catalog()
	if len(c) != len(schoolSettingCatalog) {
		t.Fatalf("Catalog() len = %d, se esperaba %d", len(c), len(schoolSettingCatalog))
	}
	c[0].Key = "mutado"
	if schoolSettingCatalog[0].Key == "mutado" {
		t.Fatal("Catalog() debe devolver una copia; mutó el catálogo interno")
	}
}

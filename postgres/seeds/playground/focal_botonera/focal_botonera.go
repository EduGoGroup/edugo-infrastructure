// Package focal_botonera es el playground 3 para validar la botonera SDUI y
// el sistema de permisos en 4 superficies (lista, tab Detalle con chips, tab
// Configuracion del master-detail, modal de pregunta). Estresa los frentes C
// (gating en chips de filtro), D (modal-actions declarativo) y F (overflow
// con priority/pin).
//
// Composicion (autosuficiente): depende de [focal_evaluacion_v2] para schools,
// units, assessments y questions (reutiliza school 62000000-...-0001 y unit
// 62000000-...-0002). Apply() encadena focal_evaluacion_v2.Apply() al inicio
// para soportar `P=focal-botonera` standalone; con `P=all` la doble
// aplicacion es no-op (OnConflict DoNothing).
//
// Rango UUID 63000000-... para usuarios/memberships; 13000000-... para roles.
// No colisiona con 60..., 61..., 62... de otros playgrounds.
package focal_botonera

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/common"
	focal_evaluacion_v2 "github.com/EduGoGroup/edugo-infrastructure/postgres/seeds/playground/focal_evaluacion_v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	// Credenciales del playground. Compartidas para simplicidad.
	ViewerEmail    = "botonera-viewer@edugo.local"
	AuthorEmail    = "botonera-author@edugo.local"
	PublisherEmail = "botonera-publisher@edugo.local"
	Password       = "12345678"

	// Tenant heredado de focal_evaluacion_v2.
	tenantSchoolID = "62000000-0000-0000-0000-000000000001"
	tenantUnitID   = "62000000-0000-0000-0000-000000000002"

	// Rango 63000000-... — usuarios + memberships del playground.
	viewerUserID    = "63000000-0000-0000-0000-000000000001"
	authorUserID    = "63000000-0000-0000-0000-000000000002"
	publisherUserID = "63000000-0000-0000-0000-000000000003"
	viewerMembID    = "63000000-0000-0000-0000-000000000011"
	authorMembID    = "63000000-0000-0000-0000-000000000012"
	publisherMembID = "63000000-0000-0000-0000-000000000013"

	// Rango 13000000-... — roles del playground.
	viewerRoleID    = "13000000-0000-0000-0000-000000000001"
	authorRoleID    = "13000000-0000-0000-0000-000000000002"
	publisherRoleID = "13000000-0000-0000-0000-000000000003"

	viewerRoleName    = "focal-viewer"
	authorRoleName    = "focal-author"
	publisherRoleName = "focal-publisher"
)

// Apply siembra el playground focal_botonera. Idempotente.
func Apply(tx *gorm.DB) error {
	if err := focal_evaluacion_v2.Apply(tx); err != nil {
		return fmt.Errorf("playground/focal_botonera: dependencia focal_evaluacion_v2: %w", err)
	}

	roles := []common.RoleSpec{
		{
			ID:          common.MustParseUUID(viewerRoleID),
			Name:        viewerRoleName,
			DisplayName: "Focal — Viewer",
			Description: "Solo lectura general (*.read). Playground focal-botonera.",
			Scope:       "school",
		},
		{
			ID:          common.MustParseUUID(authorRoleID),
			Name:        authorRoleName,
			DisplayName: "Focal — Author",
			Description: "Lee, crea y actualiza (no borra ni publica). Playground focal-botonera.",
			Scope:       "school",
		},
		{
			ID:          common.MustParseUUID(publisherRoleID),
			Name:        publisherRoleName,
			DisplayName: "Focal — Publisher",
			Description: "Catch-all (*). Playground focal-botonera.",
			Scope:       "school",
		},
	}
	for _, r := range roles {
		if err := common.SeedRole(tx, r); err != nil {
			return fmt.Errorf("playground/focal_botonera: roles: %w", err)
		}
	}

	grantSpecs := []struct {
		roleID   string
		patterns []string
	}{
		{viewerRoleID, []string{
			"content.assessments.read",
			"academic.announcements.read",
		}},
		{authorRoleID, []string{
			"content.assessments.read",
			"content.assessments.create",
			"content.assessments.update",
			"academic.announcements.read",
			"academic.announcements.create",
			"academic.announcements.update",
		}},
		{publisherRoleID, []string{
			"content.assessments.*",
			"academic.announcements.*",
		}},
	}
	for _, s := range grantSpecs {
		rid := common.MustParseUUID(s.roleID)
		for _, pattern := range s.patterns {
			if err := common.SeedRoleGrant(tx, rid, pattern); err != nil {
				return fmt.Errorf("playground/focal_botonera: role_grants: %w", err)
			}
		}
	}

	userSpecs := []common.UserSpec{
		{ID: common.MustParseUUID(viewerUserID), Email: ViewerEmail, Password: Password, FirstName: "Focal", LastName: "Viewer"},
		{ID: common.MustParseUUID(authorUserID), Email: AuthorEmail, Password: Password, FirstName: "Focal", LastName: "Author"},
		{ID: common.MustParseUUID(publisherUserID), Email: PublisherEmail, Password: Password, FirstName: "Focal", LastName: "Publisher"},
	}
	for _, u := range userSpecs {
		if err := common.SeedUser(tx, u); err != nil {
			return fmt.Errorf("playground/focal_botonera: users: %w", err)
		}
	}

	userRolePairs := [][2]string{
		{viewerUserID, viewerRoleID},
		{authorUserID, authorRoleID},
		{publisherUserID, publisherRoleID},
	}
	for _, p := range userRolePairs {
		if err := common.SeedUserRole(tx, common.MustParseUUID(p[0]), common.MustParseUUID(p[1])); err != nil {
			return fmt.Errorf("playground/focal_botonera: user_roles: %w", err)
		}
	}

	// Memberships al tenant heredado de focal_evaluacion_v2 (mismo
	// school+unit que aloja las assessments).
	auid := common.MustParseUUID(tenantUnitID)
	sid := common.MustParseUUID(tenantSchoolID)
	membershipSpecs := []common.MembershipSpec{
		{ID: common.MustParseUUID(viewerMembID), UserID: common.MustParseUUID(viewerUserID), SchoolID: sid, AcademicUnitID: &auid, Role: "teacher"},
		{ID: common.MustParseUUID(authorMembID), UserID: common.MustParseUUID(authorUserID), SchoolID: sid, AcademicUnitID: &auid, Role: "teacher"},
		{ID: common.MustParseUUID(publisherMembID), UserID: common.MustParseUUID(publisherUserID), SchoolID: sid, AcademicUnitID: &auid, Role: "admin"},
	}
	for _, m := range membershipSpecs {
		if err := common.SeedMembership(tx, m); err != nil {
			return fmt.Errorf("playground/focal_botonera: memberships: %w", err)
		}
	}

	if err := patchAssessmentsFormSlotData(tx); err != nil {
		return fmt.Errorf("playground/focal_botonera: assessments-form slot_data: %w", err)
	}
	if err := patchAssessmentQuestionsListSlotData(tx); err != nil {
		return fmt.Errorf("playground/focal_botonera: assessment-questions-list slot_data: %w", err)
	}
	if err := patchAssessmentQuestionFormSlotData(tx); err != nil {
		return fmt.Errorf("playground/focal_botonera: assessment-question-form slot_data: %w", err)
	}
	if err := upsertAnnouncements(tx); err != nil {
		return fmt.Errorf("playground/focal_botonera: announcements: %w", err)
	}
	return nil
}

// patchAssessmentsFormSlotData reemplaza slot_data del screen_instance
// `assessments-form` (master-detail-v1) con 7 actions declarativas que
// estresan los frentes:
//   - F (overflow): cada slot trae priority + pin. Threshold del template
//     resource-toolbar es 3; con publisher (que tiene todas las acciones)
//     al menos 3 caen a overflow.
//   - C/D (gating): cada permission usa un sufijo distinto para que cada
//     rol vea un subconjunto unico — viewer ve solo lo abierto, publisher
//     todo.
//
// El composer (compose.go) fusiona estos `actions_added` con los
// `default_actions` del template master-detail-v1 (save_new/save/delete/
// detail). Override por id: save se sobreescribe para inyectar priority/pin;
// detail (default "Detalle") se sobreescribe para "Help".
func patchAssessmentsFormSlotData(tx *gorm.DB) error {
	const slotDataJSON = `{
  "title": "Evaluación",
  "page_title": "Evaluación",
  "edit_title": "Editar evaluación",
  "fields": [
    {"key": "title", "label": "Título", "type": "text", "required": true},
    {"key": "description", "label": "Descripción", "type": "textarea"},
    {"key": "pass_threshold", "label": "Umbral de aprobación (%)", "type": "number"},
    {"key": "max_attempts", "label": "Intentos máximos", "type": "number"},
    {"key": "time_limit_minutes", "label": "Tiempo límite (min)", "type": "number"},
    {"key": "is_timed", "label": "Cronometrada", "type": "boolean"},
    {"key": "shuffle_questions", "label": "Mezclar preguntas", "type": "boolean"},
    {"key": "show_correct_answers", "label": "Mostrar respuestas correctas", "type": "boolean"},
    {"key": "available_from", "label": "Disponible desde", "type": "datetime"},
    {"key": "available_until", "label": "Disponible hasta", "type": "datetime"}
  ],
  "detail_config": {
    "screen_key": "assessment-questions-list",
    "modal_screen_key": "assessment-question-form",
    "parent_id_param": "assessmentId",
    "child_id_field": "id"
  },
  "actions_added": [
    {"id": "save",      "scope": "form-submit",      "label": "Guardar",   "icon": "save",         "permission": "content.assessments.update",  "condition": "edit-only",   "event_id": "submit-form", "style": "filled",      "order": 10, "priority": 1, "pin": true},
    {"id": "delete",    "scope": "form-submit",      "label": "Eliminar",  "icon": "trash",        "permission": "content.assessments.delete",  "condition": "edit-only",   "event_id": "delete",      "style": "destructive", "order": 20, "priority": 2},
    {"id": "publish",   "scope": "resource-toolbar", "label": "Publicar",  "icon": "check_circle", "permission": "content.assessments.publish", "condition": "edit-only",   "event_id": "publish",     "style": "icon",        "order": 30, "priority": 3},
    {"id": "archive",   "scope": "resource-toolbar", "label": "Archivar",  "icon": "archive",      "permission": "content.assessments.archive", "condition": "edit-only",   "event_id": "archive",     "style": "icon",        "order": 40, "priority": 4},
    {"id": "duplicate", "scope": "resource-toolbar", "label": "Duplicar",  "icon": "content_copy", "permission": "content.assessments.create",  "condition": "edit-only",   "event_id": "duplicate",   "style": "icon",        "order": 50, "priority": 5},
    {"id": "export",    "scope": "resource-toolbar", "label": "Exportar",  "icon": "download",     "permission": "content.assessments.read",    "condition": "edit-only",   "event_id": "export",      "style": "icon",        "order": 60, "priority": 6},
    {"id": "detail",    "scope": "resource-toolbar", "label": "Ayuda",     "icon": "help_outline",                                     "condition": "always",      "event_id": "view-help",   "style": "icon",        "order": 70, "priority": 7}
  ],
  "api_prefix": "learning"
}`
	return updateSlotData(tx, "assessments-form", slotDataJSON)
}

// patchAssessmentQuestionsListSlotData declara los slot labels que
// renderizan los 3 chips del template list-basic-v1. El permission gating
// de filter_processing se aplica desde AssessmentQuestionsListContract via
// FilterConfig.permission (frente C): viewer no ve ese chip.
func patchAssessmentQuestionsListSlotData(tx *gorm.DB) error {
	const slotDataJSON = `{
  "title": "Preguntas",
  "page_title": "Preguntas",
  "filter_all_label": "Todos",
  "filter_ready_label": "Opción múltiple",
  "filter_processing_label": "Solo activos",
  "columns": [
    {"key": "statement", "label": "Pregunta"},
    {"key": "kind", "label": "Tipo"},
    {"key": "score", "label": "Puntaje"}
  ],
  "actions": [
    {"id": "create", "scope": "header", "label": "Nuevo",    "icon": "plus",   "permission": "content.assessments.create"},
    {"id": "edit",   "scope": "row",    "label": "Editar",   "icon": "pencil", "permission": "content.assessments.update"},
    {"id": "delete", "scope": "row",    "label": "Eliminar", "icon": "trash",  "permission": "content.assessments.delete", "destructive": true}
  ],
  "api_prefix": "learning"
}`
	return updateSlotData(tx, "assessment-questions-list", slotDataJSON)
}

// patchAssessmentQuestionFormSlotData declara las acciones del modal
// "agregar pregunta" via scope=form-submit (alias canonico para
// modal-actions aceptado por MasterDetailItemModal.MODAL_ACTION_SCOPES).
// El modal cae al fallback hardcodeado si la zona form-submit esta vacia.
func patchAssessmentQuestionFormSlotData(tx *gorm.DB) error {
	const slotDataJSON = `{
  "title": "Pregunta",
  "page_title": "Pregunta",
  "edit_title": "Editar pregunta",
  "fields": [
    {"key": "question_text", "label": "Enunciado", "type": "textarea", "required": true},
    {"key": "question_type", "label": "Tipo", "type": "select", "required": true, "options": [
      {"value": "multiple_choice", "label": "Opción múltiple"},
      {"value": "true_false", "label": "Verdadero/Falso"},
      {"value": "short_answer", "label": "Respuesta corta"},
      {"value": "open_ended", "label": "Respuesta abierta"}
    ]},
    {"key": "points", "label": "Puntaje", "type": "number", "required": true},
    {"key": "correct_answer", "label": "Respuesta correcta", "type": "text"},
    {"key": "explanation", "label": "Explicación", "type": "textarea"},
    {"key": "difficulty", "label": "Dificultad", "type": "select", "options": [
      {"value": "easy", "label": "Fácil"},
      {"value": "medium", "label": "Media"},
      {"value": "hard", "label": "Difícil"}
    ]}
  ],
  "actions": [
    {"id": "cancel",   "scope": "form-submit", "label": "Cancelar", "icon": "close", "condition": "always",      "event_id": "cancel",      "style": "text",   "order": 5},
    {"id": "save_new", "scope": "form-submit", "label": "Guardar",  "icon": "save",  "permission": "content.assessments.create", "condition": "create-only", "event_id": "submit-form", "style": "filled", "order": 10, "priority": 1, "pin": true},
    {"id": "save",     "scope": "form-submit", "label": "Guardar",  "icon": "save",  "permission": "content.assessments.update", "condition": "edit-only",   "event_id": "submit-form", "style": "filled", "order": 10, "priority": 1, "pin": true}
  ],
  "api_prefix": "learning"
}`
	return updateSlotData(tx, "assessment-question-form", slotDataJSON)
}

// updateSlotData reemplaza por completo el campo `slot_data` del
// screen_instance identificado por screen_key. Asume que el row ya
// existe (sembrado en L4). Falla si no encuentra el row — pista
// temprana de que el playground se aplico sin la capa system.
func updateSlotData(tx *gorm.DB, screenKey, newSlotData string) error {
	res := tx.Model(&entities.ScreenInstance{}).
		Where("screen_key = ?", screenKey).
		Update("slot_data", json.RawMessage(newSlotData))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("screen_instance %q not found — apply L4 before this playground", screenKey)
	}
	return nil
}

// upsertAnnouncements siembra 3 anuncios en la escuela del playground
// (62000000-...-0001) para validar la botonera en una pantalla de tipo
// FORM (`announcement-form`, sembrada por L2). El maestro-detalle de
// evaluaciones ya cubre el shape master-detail; este seed cubre el otro
// shape principal.
//
// Autor: publisher del playground (63000000-...-0003). IDs deterministicos
// en rango 63000020-... para no colisionar con users/memberships del mismo
// playground. OnConflict DoNothing por id => idempotente.
func upsertAnnouncements(tx *gorm.DB) error {
	sid := common.MustParseUUID(tenantSchoolID)
	uid := common.MustParseUUID(tenantUnitID)
	aid := common.MustParseUUID(publisherUserID)

	now := time.Now().UTC()
	publishedNow := now
	publishedYesterday := now.AddDate(0, 0, -1)
	publishedLastWeek := now.AddDate(0, 0, -7)

	items := []entities.Announcement{
		{
			ID:             common.MustParseUUID("63000020-0000-0000-0000-000000000001"),
			SchoolID:       sid,
			AcademicUnitID: &uid,
			AuthorID:       aid,
			Title:          "Aviso fijado del playground botonera",
			Body:           "Anuncio pinned para validar la botonera del form simple (save / delete / cancel).",
			Scope:          "school",
			IsPinned:       true,
			PublishedAt:    &publishedLastWeek,
		},
		{
			ID:             common.MustParseUUID("63000020-0000-0000-0000-000000000002"),
			SchoolID:       sid,
			AcademicUnitID: &uid,
			AuthorID:       aid,
			Title:          "Anuncio regular de prueba",
			Body:           "Cuerpo de un anuncio normal sin pin. Sirve para abrir el form y probar edicion.",
			Scope:          "school",
			IsPinned:       false,
			PublishedAt:    &publishedYesterday,
		},
		{
			ID:             common.MustParseUUID("63000020-0000-0000-0000-000000000003"),
			SchoolID:       sid,
			AcademicUnitID: &uid,
			AuthorID:       aid,
			Title:          "Borrador reciente",
			Body:           "Anuncio publicado ahora mismo para validar el orden por created_at.",
			Scope:          "school",
			IsPinned:       false,
			PublishedAt:    &publishedNow,
		},
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&items).Error
}

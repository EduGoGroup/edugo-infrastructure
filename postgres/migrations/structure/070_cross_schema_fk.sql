-- ============================================================
-- 070: Cross-schema Foreign Keys
-- Todas las foreign keys que cruzan límites de schema
-- Se ejecuta después de que todas las tablas existen
-- ============================================================

-- auth.refresh_tokens -> auth.users
ALTER TABLE auth.refresh_tokens ADD CONSTRAINT fk_refresh_tokens_user
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE;

-- academic.memberships -> auth.users
ALTER TABLE academic.memberships ADD CONSTRAINT memberships_user_fkey
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE;

-- academic.guardian_relations -> auth.users
ALTER TABLE academic.guardian_relations ADD CONSTRAINT guardian_relations_guardian_fkey
    FOREIGN KEY (guardian_id) REFERENCES auth.users(id) ON DELETE CASCADE;
ALTER TABLE academic.guardian_relations ADD CONSTRAINT guardian_relations_student_fkey
    FOREIGN KEY (student_id) REFERENCES auth.users(id) ON DELETE CASCADE;

-- content.materials -> academic.schools, auth.users, academic.academic_units
ALTER TABLE content.materials ADD CONSTRAINT materials_school_fkey
    FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE;
ALTER TABLE content.materials ADD CONSTRAINT materials_teacher_fkey
    FOREIGN KEY (uploaded_by_teacher_id) REFERENCES auth.users(id) ON DELETE RESTRICT;
ALTER TABLE content.materials ADD CONSTRAINT materials_unit_fkey
    FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE SET NULL;

-- content.material_versions -> auth.users
ALTER TABLE content.material_versions ADD CONSTRAINT material_versions_created_by_fkey
    FOREIGN KEY (created_by) REFERENCES auth.users(id) ON DELETE SET NULL;

-- content.progress -> auth.users
ALTER TABLE content.progress ADD CONSTRAINT progress_user_fkey
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE;

-- assessment.assessment -> content.materials
ALTER TABLE assessment.assessment ADD CONSTRAINT assessment_material_fkey
    FOREIGN KEY (material_id) REFERENCES content.materials(id) ON DELETE CASCADE;

-- assessment.assessment_attempt -> auth.users
ALTER TABLE assessment.assessment_attempt ADD CONSTRAINT assessment_attempt_student_fkey
    FOREIGN KEY (student_id) REFERENCES auth.users(id) ON DELETE CASCADE;

-- iam.user_roles -> auth.users, academic.schools, academic.academic_units
ALTER TABLE iam.user_roles ADD CONSTRAINT fk_user_roles_user
    FOREIGN KEY (user_id) REFERENCES auth.users(id) ON DELETE CASCADE;
ALTER TABLE iam.user_roles ADD CONSTRAINT fk_user_roles_school
    FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE;
ALTER TABLE iam.user_roles ADD CONSTRAINT fk_user_roles_unit
    FOREIGN KEY (academic_unit_id) REFERENCES academic.academic_units(id) ON DELETE CASCADE;
ALTER TABLE iam.user_roles ADD CONSTRAINT fk_user_roles_granted_by
    FOREIGN KEY (granted_by) REFERENCES auth.users(id) ON DELETE SET NULL;

-- ui_config.screen_templates -> auth.users
ALTER TABLE ui_config.screen_templates ADD CONSTRAINT fk_screen_templates_created_by
    FOREIGN KEY (created_by) REFERENCES auth.users(id);

-- ui_config.screen_instances -> auth.users
ALTER TABLE ui_config.screen_instances ADD CONSTRAINT fk_screen_instances_created_by
    FOREIGN KEY (created_by) REFERENCES auth.users(id);

-- ui_config.resource_screens -> iam.resources
ALTER TABLE ui_config.resource_screens ADD CONSTRAINT fk_resource_screens_resource
    FOREIGN KEY (resource_id) REFERENCES iam.resources(id);

-- ui_config.screen_user_preferences -> auth.users
ALTER TABLE ui_config.screen_user_preferences ADD CONSTRAINT fk_screen_user_prefs_user
    FOREIGN KEY (user_id) REFERENCES auth.users(id);

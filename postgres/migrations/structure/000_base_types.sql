-- ====================================================================
-- TIPOS PERSONALIZADOS Y EXTENSIONES PARA POSTGRESQL
-- ====================================================================
-- Este archivo contiene las extensiones necesarias, definiciones de tipos ENUM
-- y tipos personalizados utilizados por múltiples tablas en el sistema EduGo.
-- Debe ejecutarse ANTES que cualquier otra migración.
-- VERSIÓN: postgres/v0.16.1
-- FECHA: 2026-02-13
-- ====================================================================

-- ====================================================================
-- EXTENSIÓN: uuid-ossp
-- DESCRIPCIÓN: Habilita funciones para generar UUIDs (uuid_generate_v4())
-- ====================================================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ====================================================================
-- TIPO: permission_scope
-- DESCRIPCIÓN: Define los alcances posibles de un permiso en el sistema RBAC
-- VALORES:
--   - system: Permiso a nivel de sistema (aplicable globalmente)
--   - school: Permiso a nivel de escuela (aplicable a una escuela específica)
--   - unit: Permiso a nivel de unidad académica (aplicable a una unidad específica)
-- ====================================================================
CREATE TYPE permission_scope AS ENUM ('system', 'school', 'unit');

COMMENT ON TYPE permission_scope IS 'Define los alcances posibles de un permiso en el sistema RBAC';

-- ============================================================
-- 094: academic.calendar_events
-- Schema: academic
-- Calendario escolar (feriados, examenes, reuniones)
-- ============================================================

CREATE TABLE academic.calendar_events (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    school_id uuid NOT NULL,
    title character varying(200) NOT NULL,
    description text,
    event_type character varying(30) NOT NULL,
    start_date date NOT NULL,
    end_date date,
    is_all_day boolean DEFAULT true NOT NULL,
    created_by uuid NOT NULL,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    CONSTRAINT calendar_events_pkey PRIMARY KEY (id),
    CONSTRAINT calendar_events_type_check CHECK (event_type IN ('holiday', 'exam', 'meeting', 'activity', 'deadline')),
    CONSTRAINT calendar_events_school_fkey FOREIGN KEY (school_id) REFERENCES academic.schools(id) ON DELETE CASCADE
);

CREATE INDEX idx_calendar_school ON academic.calendar_events(school_id);
CREATE INDEX idx_calendar_dates ON academic.calendar_events(start_date, end_date);

CREATE TRIGGER set_updated_at BEFORE UPDATE ON academic.calendar_events
    FOR EACH ROW EXECUTE FUNCTION public.update_updated_at_column();

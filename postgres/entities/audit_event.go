package entities

import (
	"time"

	"github.com/google/uuid"
)

// AuditEvent represents the audit.events table in PostgreSQL.
type AuditEvent struct {
	ID             uuid.UUID      `db:"id" gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()" validate:"required,uuid"`
	ActorID        uuid.UUID      `db:"actor_id" gorm:"column:actor_id;type:uuid;not null;index:idx_audit_events_actor" validate:"required,uuid"`
	ActorEmail     string         `db:"actor_email" gorm:"column:actor_email;not null;size:255" validate:"required,email"`
	ActorRole      string         `db:"actor_role" gorm:"column:actor_role;not null;size:100" validate:"required,max=100"`
	ActorIP        *string        `db:"actor_ip" gorm:"column:actor_ip;size:45" validate:"omitempty"`
	ActorUserAgent *string        `db:"actor_user_agent" gorm:"column:actor_user_agent" validate:"omitempty"`
	SchoolID       *uuid.UUID     `db:"school_id" gorm:"column:school_id;type:uuid;index:idx_audit_events_school" validate:"omitempty,uuid"`
	UnitID         *uuid.UUID     `db:"unit_id" gorm:"column:unit_id;type:uuid" validate:"omitempty,uuid"`
	ServiceName    string         `db:"service_name" gorm:"column:service_name;not null;size:50" validate:"required,min=2,max=50"`
	Action         string         `db:"action" gorm:"column:action;not null;size:100;index:idx_audit_events_action" validate:"required,max=100"`
	ResourceType   string         `db:"resource_type" gorm:"column:resource_type;not null;size:100;index:idx_audit_events_resource" validate:"required,max=100"`
	ResourceID     *string        `db:"resource_id" gorm:"column:resource_id;type:varchar(255);size:255;index:idx_audit_events_resource" validate:"omitempty"`
	PermissionUsed *string        `db:"permission_used" gorm:"column:permission_used;size:100" validate:"omitempty"`
	RequestMethod  *string        `db:"request_method" gorm:"column:request_method;size:10" validate:"omitempty"`
	RequestPath    *string        `db:"request_path" gorm:"column:request_path;size:500" validate:"omitempty"`
	RequestID      *string        `db:"request_id" gorm:"column:request_id;size:100" validate:"omitempty"`
	StatusCode     *int           `db:"status_code" gorm:"column:status_code" validate:"omitempty"`
	Changes        map[string]any `db:"changes" gorm:"column:changes;serializer:json;type:jsonb" validate:"-"`
	Metadata       map[string]any `db:"metadata" gorm:"column:metadata;serializer:json;type:jsonb" validate:"-"`
	ErrorMessage   *string        `db:"error_message" gorm:"column:error_message" validate:"omitempty"`
	CreatedAt      time.Time      `db:"created_at" gorm:"column:created_at;not null;autoCreateTime;index:idx_audit_events_actor;index:idx_audit_events_resource;index:idx_audit_events_action;index:idx_audit_events_school;index:idx_audit_events_created" validate:"-"`
	// NOTE: partial index idx_audit_events_severity (WHERE severity != 'info') must be created in post_gorm.sql
	Severity string `db:"severity" gorm:"column:severity;not null;size:20;default:'info';check:audit_events_severity_check,severity IN ('info','warning','critical')" validate:"required,oneof=info warning critical"`
	Category string `db:"category" gorm:"column:category;not null;size:50;default:'data';check:audit_events_category_check,category IN ('auth','data','config','admin')" validate:"required,oneof=auth data config admin"`
}

func (AuditEvent) TableName() string {
	return "audit.events"
}

// Migration: Create audit_logs collection
// Collection: audit_logs (Owner: infrastructure)
// Created by: edugo-infrastructure
// Used by: api-admin, api-mobile, worker
//
// Purpose: Stores audit trail of important system events and user actions for compliance and debugging

db.createCollection("audit_logs", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["event_type", "actor_id", "timestamp", "resource_type"],
      properties: {
        event_type: {
          bsonType: "string",
          enum: [
            "user.created", "user.updated", "user.deleted", "user.login", "user.logout",
            "school.created", "school.updated", "school.deleted",
            "material.uploaded", "material.updated", "material.deleted", "material.processed",
            "assessment.generated", "assessment.published", "assessment.archived",
            "attempt.started", "attempt.submitted", "attempt.graded",
            "membership.created", "membership.updated", "membership.deleted",
            "permission.granted", "permission.revoked",
            "system.backup", "system.restore", "system.migration"
          ],
          description: "Type of event that occurred"
        },
        actor_id: {
          bsonType: "string",
          description: "UUID of the user who performed the action (or 'system' for automated actions)"
        },
        actor_type: {
          bsonType: "string",
          enum: ["user", "system", "api", "worker"],
          description: "Type of actor that performed the action"
        },
        resource_type: {
          bsonType: "string",
          enum: ["user", "school", "academic_unit", "membership", "material", "assessment", "attempt", "system"],
          description: "Type of resource affected"
        },
        resource_id: {
          bsonType: "string",
          description: "ID of the affected resource"
        },
        action: {
          bsonType: "string",
          enum: ["create", "read", "update", "delete", "login", "logout", "upload", "process", "submit", "grade"],
          description: "Action performed"
        },
        details: {
          bsonType: "object",
          description: "Additional details about the event",
          properties: {
            ip_address: {
              bsonType: "string",
              description: "IP address of the actor"
            },
            user_agent: {
              bsonType: "string",
              description: "User agent string"
            },
            changes: {
              bsonType: "object",
              description: "What changed (before/after values)"
            },
            metadata: {
              bsonType: "object",
              description: "Additional metadata"
            },
            error: {
              bsonType: "object",
              description: "Error details if the action failed"
            }
          }
        },
        severity: {
          bsonType: "string",
          enum: ["info", "warning", "error", "critical"],
          description: "Severity level of the event"
        },
        timestamp: {
          bsonType: "date",
          description: "When the event occurred"
        },
        session_id: {
          bsonType: "string",
          description: "Session ID if applicable"
        },
        request_id: {
          bsonType: "string",
          description: "Request ID for tracking across services"
        }
      }
    }
  }
});

// Create indexes for efficient queries and log analysis
db.audit_logs.createIndex({ "timestamp": -1 }, { name: "idx_timestamp_desc" });
db.audit_logs.createIndex({ "event_type": 1, "timestamp": -1 }, { name: "idx_event_type_timestamp" });
db.audit_logs.createIndex({ "actor_id": 1, "timestamp": -1 }, { name: "idx_actor_timestamp" });
db.audit_logs.createIndex({ "resource_type": 1, "resource_id": 1 }, { name: "idx_resource" });
db.audit_logs.createIndex({ "severity": 1, "timestamp": -1 }, { name: "idx_severity_timestamp" });
db.audit_logs.createIndex({ "session_id": 1 }, { name: "idx_session_id" });
db.audit_logs.createIndex({ "request_id": 1 }, { name: "idx_request_id" });

// TTL index to automatically delete old logs after 90 days (configurable)
db.audit_logs.createIndex(
  { "timestamp": 1 },
  { name: "idx_ttl_90days", expireAfterSeconds: 7776000 }  // 90 days = 90 * 24 * 60 * 60
);

print("✅ Collection 'audit_logs' created successfully");

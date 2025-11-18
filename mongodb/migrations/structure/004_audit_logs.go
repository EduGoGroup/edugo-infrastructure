package structure

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateAuditLogs creates the audit_logs collection with schema validation
// Collection: audit_logs (Owner: infrastructure)
// Used by: api-mobile, worker
// Purpose: Stores comprehensive audit trail of all system events
func CreateAuditLogs(ctx context.Context, db *mongo.Database) error {
	collectionName := "audit_logs"

	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"event_type", "actor_id", "timestamp", "resource_type"},
			"properties": bson.M{
				"event_type": bson.M{
					"bsonType": "string",
					"enum": []string{
						"user.created", "user.updated", "user.deleted", "user.login", "user.logout",
						"school.created", "school.updated", "school.deleted",
						"material.uploaded", "material.updated", "material.deleted", "material.processed",
						"assessment.generated", "assessment.published", "assessment.archived",
						"attempt.started", "attempt.submitted", "attempt.graded",
						"membership.created", "membership.updated", "membership.deleted",
						"permission.granted", "permission.revoked",
						"system.backup", "system.restore", "system.migration",
					},
					"description": "Type of audit event",
				},
				"actor_id": bson.M{
					"bsonType":    "string",
					"description": "ID of the user/system performing the action",
				},
				"actor_type": bson.M{
					"bsonType":    "string",
					"enum":        []string{"user", "system", "api", "worker"},
					"description": "Type of actor",
				},
				"timestamp": bson.M{
					"bsonType": "date",
				},
				"resource_type": bson.M{
					"bsonType":    "string",
					"enum":        []string{"user", "school", "academic_unit", "membership", "material", "assessment", "attempt", "system"},
					"description": "Type of resource affected",
				},
				"resource_id": bson.M{
					"bsonType":    "string",
					"description": "ID of the resource affected",
				},
				"action": bson.M{
					"bsonType":    "string",
					"enum":        []string{"create", "read", "update", "delete", "login", "logout", "upload", "process", "submit", "grade"},
					"description": "Action performed",
				},
				"changes": bson.M{
					"bsonType":    "object",
					"description": "Details of changes made",
				},
				"metadata": bson.M{
					"bsonType":    "object",
					"description": "Additional context information",
				},
				"ip_address": bson.M{
					"bsonType":    "string",
					"description": "IP address of the actor",
				},
				"user_agent": bson.M{
					"bsonType":    "string",
					"description": "User agent string",
				},
				"severity": bson.M{
					"bsonType":    "string",
					"enum":        []string{"info", "warning", "error", "critical"},
					"description": "Severity level of the event",
				},
			},
		},
	}

	opts := options.CreateCollection().SetValidator(validator)
	err := db.CreateCollection(ctx, collectionName, opts)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	return nil
}

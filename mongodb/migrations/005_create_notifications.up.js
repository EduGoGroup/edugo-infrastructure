// Migration: Create notifications collection
// Collection: notifications (Owner: infrastructure)
// Created by: edugo-infrastructure
// Used by: api-admin, api-mobile, worker
//
// Purpose: Stores user notifications (in-app, push, email) for various system events

db.createCollection("notifications", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["user_id", "notification_type", "title", "is_read", "created_at"],
      properties: {
        user_id: {
          bsonType: "string",
          description: "UUID of the user receiving the notification"
        },
        notification_type: {
          bsonType: "string",
          enum: [
            "assessment.ready",
            "assessment.graded",
            "material.uploaded",
            "material.processed",
            "material.shared",
            "membership.added",
            "membership.removed",
            "deadline.approaching",
            "achievement.unlocked",
            "system.announcement",
            "system.maintenance"
          ],
          description: "Type of notification"
        },
        title: {
          bsonType: "string",
          description: "Notification title"
        },
        message: {
          bsonType: "string",
          description: "Notification message body"
        },
        priority: {
          bsonType: "string",
          enum: ["low", "medium", "high", "urgent"],
          description: "Priority level of the notification"
        },
        category: {
          bsonType: "string",
          enum: ["academic", "administrative", "social", "system"],
          description: "Category of the notification"
        },
        data: {
          bsonType: "object",
          description: "Additional data for the notification",
          properties: {
            resource_type: {
              bsonType: "string",
              description: "Type of related resource (material, assessment, etc.)"
            },
            resource_id: {
              bsonType: "string",
              description: "ID of the related resource"
            },
            action_url: {
              bsonType: "string",
              description: "Deep link or URL for action button"
            },
            action_label: {
              bsonType: "string",
              description: "Label for action button"
            },
            metadata: {
              bsonType: "object",
              description: "Additional metadata"
            }
          }
        },
        delivery: {
          bsonType: "object",
          description: "Delivery status across channels",
          properties: {
            in_app: {
              bsonType: "object",
              properties: {
                enabled: { bsonType: "bool" },
                delivered_at: { bsonType: "date" }
              }
            },
            push: {
              bsonType: "object",
              properties: {
                enabled: { bsonType: "bool" },
                sent_at: { bsonType: "date" },
                delivered_at: { bsonType: "date" },
                error: { bsonType: "string" }
              }
            },
            email: {
              bsonType: "object",
              properties: {
                enabled: { bsonType: "bool" },
                sent_at: { bsonType: "date" },
                delivered_at: { bsonType: "date" },
                error: { bsonType: "string" }
              }
            }
          }
        },
        is_read: {
          bsonType: "bool",
          description: "Whether the user has read the notification"
        },
        read_at: {
          bsonType: "date",
          description: "When the notification was read"
        },
        is_archived: {
          bsonType: "bool",
          description: "Whether the notification is archived"
        },
        archived_at: {
          bsonType: "date",
          description: "When the notification was archived"
        },
        expires_at: {
          bsonType: "date",
          description: "When the notification expires (optional)"
        },
        created_at: {
          bsonType: "date",
          description: "When the notification was created"
        }
      }
    }
  }
});

// Create indexes for efficient queries
db.notifications.createIndex({ "user_id": 1, "created_at": -1 }, { name: "idx_user_created_desc" });
db.notifications.createIndex({ "user_id": 1, "is_read": 1 }, { name: "idx_user_unread" });
db.notifications.createIndex({ "notification_type": 1 }, { name: "idx_notification_type" });
db.notifications.createIndex({ "priority": 1, "created_at": -1 }, { name: "idx_priority_created" });
db.notifications.createIndex({ "created_at": -1 }, { name: "idx_created_at_desc" });
db.notifications.createIndex({ "data.resource_type": 1, "data.resource_id": 1 }, { name: "idx_resource" });

// TTL index to automatically delete expired notifications
db.notifications.createIndex(
  { "expires_at": 1 },
  { name: "idx_ttl_expires", expireAfterSeconds: 0 }
);

// TTL index to delete read and archived notifications after 30 days
db.notifications.createIndex(
  { "archived_at": 1 },
  { name: "idx_ttl_archived_30days", expireAfterSeconds: 2592000 }  // 30 days
);

print("✅ Collection 'notifications' created successfully");

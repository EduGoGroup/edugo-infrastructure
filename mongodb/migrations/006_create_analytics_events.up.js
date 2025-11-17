// Migration: Create analytics_events collection
// Collection: analytics_events (Owner: infrastructure)
// Created by: edugo-infrastructure
// Used by: api-mobile, worker
//
// Purpose: Stores user behavior and analytics events for insights and reporting

db.createCollection("analytics_events", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["event_name", "timestamp"],
      properties: {
        event_name: {
          bsonType: "string",
          enum: [
            "page.view",
            "material.view",
            "material.download",
            "material.search",
            "assessment.start",
            "assessment.complete",
            "assessment.abandon",
            "question.answer",
            "question.skip",
            "video.play",
            "video.pause",
            "video.complete",
            "session.start",
            "session.end",
            "feature.click",
            "error.occurred",
            "search.performed",
            "filter.applied"
          ],
          description: "Name of the analytics event"
        },
        user_id: {
          bsonType: "string",
          description: "UUID of the user (null for anonymous events)"
        },
        session_id: {
          bsonType: "string",
          description: "Session ID to track user sessions"
        },
        timestamp: {
          bsonType: "date",
          description: "When the event occurred"
        },
        properties: {
          bsonType: "object",
          description: "Event-specific properties",
          properties: {
            page_path: {
              bsonType: "string",
              description: "Page path for navigation events"
            },
            page_title: {
              bsonType: "string",
              description: "Page title"
            },
            resource_id: {
              bsonType: "string",
              description: "ID of related resource (material, assessment, etc.)"
            },
            resource_type: {
              bsonType: "string",
              description: "Type of resource"
            },
            duration_seconds: {
              bsonType: "int",
              minimum: 0,
              description: "Duration of the event in seconds"
            },
            search_query: {
              bsonType: "string",
              description: "Search query text"
            },
            search_results_count: {
              bsonType: "int",
              description: "Number of search results"
            },
            button_label: {
              bsonType: "string",
              description: "Label of clicked button"
            },
            error_message: {
              bsonType: "string",
              description: "Error message if applicable"
            },
            custom_data: {
              bsonType: "object",
              description: "Additional custom properties"
            }
          }
        },
        device: {
          bsonType: "object",
          description: "Device and browser information",
          properties: {
            platform: {
              bsonType: "string",
              enum: ["web", "ios", "android"],
              description: "Platform/device type"
            },
            os: {
              bsonType: "string",
              description: "Operating system"
            },
            os_version: {
              bsonType: "string",
              description: "OS version"
            },
            browser: {
              bsonType: "string",
              description: "Browser name"
            },
            browser_version: {
              bsonType: "string",
              description: "Browser version"
            },
            device_type: {
              bsonType: "string",
              enum: ["mobile", "tablet", "desktop"],
              description: "Device type"
            },
            screen_resolution: {
              bsonType: "string",
              description: "Screen resolution (e.g., '1920x1080')"
            }
          }
        },
        location: {
          bsonType: "object",
          description: "Geographic information",
          properties: {
            ip_address: {
              bsonType: "string",
              description: "IP address (anonymized for privacy)"
            },
            country: {
              bsonType: "string",
              description: "Country code (ISO 3166-1 alpha-2)"
            },
            city: {
              bsonType: "string",
              description: "City name"
            },
            timezone: {
              bsonType: "string",
              description: "Timezone (e.g., 'America/Santiago')"
            }
          }
        },
        context: {
          bsonType: "object",
          description: "Additional context about the event",
          properties: {
            school_id: {
              bsonType: "string",
              description: "School ID if applicable"
            },
            academic_unit_id: {
              bsonType: "string",
              description: "Academic unit ID if applicable"
            },
            user_role: {
              bsonType: "string",
              enum: ["admin", "teacher", "student", "guardian"],
              description: "Role of the user"
            },
            ab_test_variant: {
              bsonType: "string",
              description: "A/B test variant if applicable"
            }
          }
        }
      }
    }
  }
});

// Create indexes for efficient querying and analytics
db.analytics_events.createIndex({ "timestamp": -1 }, { name: "idx_timestamp_desc" });
db.analytics_events.createIndex({ "event_name": 1, "timestamp": -1 }, { name: "idx_event_timestamp" });
db.analytics_events.createIndex({ "user_id": 1, "timestamp": -1 }, { name: "idx_user_timestamp" });
db.analytics_events.createIndex({ "session_id": 1, "timestamp": 1 }, { name: "idx_session_timeline" });
db.analytics_events.createIndex({ "properties.resource_type": 1, "properties.resource_id": 1 }, { name: "idx_resource" });
db.analytics_events.createIndex({ "context.school_id": 1, "timestamp": -1 }, { name: "idx_school_timestamp" });
db.analytics_events.createIndex({ "device.platform": 1, "timestamp": -1 }, { name: "idx_platform_timestamp" });

// Compound index for common analytics queries
db.analytics_events.createIndex(
  {
    "event_name": 1,
    "context.school_id": 1,
    "timestamp": -1
  },
  { name: "idx_event_school_timestamp" }
);

// TTL index to automatically delete old analytics events after 365 days
db.analytics_events.createIndex(
  { "timestamp": 1 },
  { name: "idx_ttl_365days", expireAfterSeconds: 31536000 }  // 365 days
);

print("✅ Collection 'analytics_events' created successfully");

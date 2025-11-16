// Seeds for audit_logs collection
// Execute with: mongosh edugo < audit_logs.js
// Or: mongosh --eval "$(cat audit_logs.js)"

db = db.getSiblingDB('edugo');

// Audit log 1 - User login
db.audit_logs.insertOne({
  event_type: "user.login",
  actor_id: "11111111-1111-1111-1111-111111111111",
  actor_type: "user",
  resource_type: "user",
  resource_id: "11111111-1111-1111-1111-111111111111",
  action: "login",
  details: {
    ip_address: "192.168.1.100",
    user_agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
    metadata: {
      login_method: "email_password",
      remember_me: true
    }
  },
  severity: "info",
  timestamp: new Date("2025-01-15T09:00:00Z"),
  session_id: "sess_abc123xyz",
  request_id: "req_001"
});

// Audit log 2 - Material uploaded
db.audit_logs.insertOne({
  event_type: "material.uploaded",
  actor_id: "22222222-2222-2222-2222-222222222222",
  actor_type: "user",
  resource_type: "material",
  resource_id: "66666666-6666-6666-6666-666666666666",
  action: "upload",
  details: {
    ip_address: "192.168.1.101",
    user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    changes: {
      file_name: "fisica_cuantica.pdf",
      file_size: 2048576,
      school_id: "55555555-5555-5555-5555-555555555555"
    }
  },
  severity: "info",
  timestamp: new Date("2025-01-15T10:00:00Z"),
  session_id: "sess_def456uvw",
  request_id: "req_002"
});

// Audit log 3 - Assessment published
db.audit_logs.insertOne({
  event_type: "assessment.published",
  actor_id: "22222222-2222-2222-2222-222222222222",
  actor_type: "user",
  resource_type: "assessment",
  resource_id: "99999999-9999-9999-9999-999999999999",
  action: "update",
  details: {
    ip_address: "192.168.1.101",
    user_agent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    changes: {
      status: {
        from: "generated",
        to: "published"
      }
    }
  },
  severity: "info",
  timestamp: new Date("2025-01-15T10:30:00Z"),
  session_id: "sess_def456uvw",
  request_id: "req_003"
});

// Audit log 4 - Failed login attempt
db.audit_logs.insertOne({
  event_type: "user.login",
  actor_id: "unknown",
  actor_type: "user",
  resource_type: "user",
  resource_id: "unknown",
  action: "login",
  details: {
    ip_address: "192.168.1.200",
    user_agent: "Mozilla/5.0 (X11; Linux x86_64)",
    error: {
      code: "INVALID_CREDENTIALS",
      message: "Invalid email or password"
    },
    metadata: {
      attempted_email: "test@example.com"
    }
  },
  severity: "warning",
  timestamp: new Date("2025-01-15T11:00:00Z"),
  request_id: "req_004"
});

// Audit log 5 - System backup
db.audit_logs.insertOne({
  event_type: "system.backup",
  actor_id: "system",
  actor_type: "system",
  resource_type: "system",
  action: "create",
  details: {
    metadata: {
      backup_type: "automated_daily",
      backup_size_mb: 1024,
      backup_location: "s3://edugo-backups/2025-01-15/"
    }
  },
  severity: "info",
  timestamp: new Date("2025-01-15T02:00:00Z"),
  request_id: "req_005"
});

print("✅ 5 audit_logs documents inserted successfully");

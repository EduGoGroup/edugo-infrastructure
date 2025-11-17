// Migration DOWN: Drop audit_logs collection

db.audit_logs.drop();

print("✅ Collection 'audit_logs' dropped successfully");

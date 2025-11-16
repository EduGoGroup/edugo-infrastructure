// Migration DOWN: Drop notifications collection

db.notifications.drop();

print("✅ Collection 'notifications' dropped successfully");

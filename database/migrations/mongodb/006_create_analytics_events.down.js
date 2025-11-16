// Migration DOWN: Drop analytics_events collection

db.analytics_events.drop();

print("✅ Collection 'analytics_events' dropped successfully");

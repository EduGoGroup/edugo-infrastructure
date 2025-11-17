// Seeds for analytics_events collection
// Execute with: mongosh edugo < analytics_events.js
// Or: mongosh --eval "$(cat analytics_events.js)"

db = db.getSiblingDB('edugo');

// Event 1 - Page view
db.analytics_events.insertOne({
  event_name: "page.view",
  user_id: "33333333-3333-3333-3333-333333333333",
  session_id: "sess_student_abc123",
  timestamp: new Date("2025-01-15T10:00:00Z"),
  properties: {
    page_path: "/materials",
    page_title: "Mis Materiales"
  },
  device: {
    platform: "web",
    os: "macOS",
    os_version: "14.0",
    browser: "Chrome",
    browser_version: "120.0",
    device_type: "desktop",
    screen_resolution: "1920x1080"
  },
  location: {
    country: "CL",
    city: "Santiago",
    timezone: "America/Santiago"
  },
  context: {
    school_id: "55555555-5555-5555-5555-555555555555",
    user_role: "student"
  }
});

// Event 2 - Material view
db.analytics_events.insertOne({
  event_name: "material.view",
  user_id: "33333333-3333-3333-3333-333333333333",
  session_id: "sess_student_abc123",
  timestamp: new Date("2025-01-15T10:01:00Z"),
  properties: {
    resource_id: "66666666-6666-6666-6666-666666666666",
    resource_type: "material",
    custom_data: {
      material_title: "Física Cuántica - Introducción",
      subject: "Física"
    }
  },
  device: {
    platform: "web",
    os: "macOS",
    os_version: "14.0",
    browser: "Chrome",
    browser_version: "120.0",
    device_type: "desktop",
    screen_resolution: "1920x1080"
  },
  location: {
    country: "CL",
    city: "Santiago",
    timezone: "America/Santiago"
  },
  context: {
    school_id: "55555555-5555-5555-5555-555555555555",
    user_role: "student"
  }
});

// Event 3 - Assessment start
db.analytics_events.insertOne({
  event_name: "assessment.start",
  user_id: "33333333-3333-3333-3333-333333333333",
  session_id: "sess_student_abc123",
  timestamp: new Date("2025-01-15T10:14:00Z"),
  properties: {
    resource_id: "99999999-9999-9999-9999-999999999999",
    resource_type: "assessment",
    custom_data: {
      questions_count: 2,
      subject: "Física"
    }
  },
  device: {
    platform: "web",
    os: "macOS",
    os_version: "14.0",
    browser: "Chrome",
    browser_version: "120.0",
    device_type: "desktop",
    screen_resolution: "1920x1080"
  },
  location: {
    country: "CL",
    city: "Santiago",
    timezone: "America/Santiago"
  },
  context: {
    school_id: "55555555-5555-5555-5555-555555555555",
    user_role: "student"
  }
});

// Event 4 - Assessment complete
db.analytics_events.insertOne({
  event_name: "assessment.complete",
  user_id: "33333333-3333-3333-3333-333333333333",
  session_id: "sess_student_abc123",
  timestamp: new Date("2025-01-15T10:16:15Z"),
  properties: {
    resource_id: "99999999-9999-9999-9999-999999999999",
    resource_type: "assessment",
    duration_seconds: 135,
    custom_data: {
      score: 100,
      questions_count: 2,
      correct_answers: 2
    }
  },
  device: {
    platform: "web",
    os: "macOS",
    os_version: "14.0",
    browser: "Chrome",
    browser_version: "120.0",
    device_type: "desktop",
    screen_resolution: "1920x1080"
  },
  location: {
    country: "CL",
    city: "Santiago",
    timezone: "America/Santiago"
  },
  context: {
    school_id: "55555555-5555-5555-5555-555555555555",
    user_role: "student"
  }
});

// Event 5 - Search performed
db.analytics_events.insertOne({
  event_name: "search.performed",
  user_id: "33333333-3333-3333-3333-333333333333",
  session_id: "sess_student_abc123",
  timestamp: new Date("2025-01-15T10:30:00Z"),
  properties: {
    search_query: "álgebra matrices",
    search_results_count: 3,
    custom_data: {
      filters_applied: {
        subject: "Matemáticas"
      }
    }
  },
  device: {
    platform: "web",
    os: "macOS",
    os_version: "14.0",
    browser: "Chrome",
    browser_version: "120.0",
    device_type: "desktop",
    screen_resolution: "1920x1080"
  },
  location: {
    country: "CL",
    city: "Santiago",
    timezone: "America/Santiago"
  },
  context: {
    school_id: "55555555-5555-5555-5555-555555555555",
    user_role: "student"
  }
});

// Event 6 - Mobile app session
db.analytics_events.insertOne({
  event_name: "session.start",
  user_id: "44444444-4444-4444-4444-444444444444",
  session_id: "sess_mobile_xyz789",
  timestamp: new Date("2025-01-15T11:00:00Z"),
  properties: {},
  device: {
    platform: "android",
    os: "Android",
    os_version: "13",
    device_type: "mobile",
    screen_resolution: "1080x2400"
  },
  location: {
    country: "CL",
    city: "Valparaíso",
    timezone: "America/Santiago"
  },
  context: {
    school_id: "55555555-5555-5555-5555-555555555555",
    user_role: "student"
  }
});

print("✅ 6 analytics_events documents inserted successfully");

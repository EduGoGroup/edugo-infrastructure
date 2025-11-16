// Seeds for notifications collection
// Execute with: mongosh --host localhost:27017/edugo < notifications.js

use edugo;

// Notification 1 - Assessment ready for student
db.notifications.insertOne({
  user_id: "33333333-3333-3333-3333-333333333333",
  notification_type: "assessment.ready",
  title: "Nuevo Assessment Disponible",
  message: "Tu profesor ha publicado un nuevo assessment de Física Cuántica. ¡Es hora de demostrar lo que has aprendido!",
  priority: "medium",
  category: "academic",
  data: {
    resource_type: "assessment",
    resource_id: "99999999-9999-9999-9999-999999999999",
    action_url: "/assessments/99999999-9999-9999-9999-999999999999",
    action_label: "Comenzar Assessment",
    metadata: {
      subject: "Física",
      teacher_name: "Prof. García"
    }
  },
  delivery: {
    in_app: {
      enabled: true,
      delivered_at: new Date("2025-01-15T10:30:00Z")
    },
    push: {
      enabled: true,
      sent_at: new Date("2025-01-15T10:30:01Z"),
      delivered_at: new Date("2025-01-15T10:30:02Z")
    },
    email: {
      enabled: false
    }
  },
  is_read: false,
  is_archived: false,
  created_at: new Date("2025-01-15T10:30:00Z")
});

// Notification 2 - Assessment graded
db.notifications.insertOne({
  user_id: "33333333-3333-3333-3333-333333333333",
  notification_type: "assessment.graded",
  title: "Assessment Calificado",
  message: "¡Felicitaciones! Has obtenido 100% en el assessment de Física Cuántica.",
  priority: "high",
  category: "academic",
  data: {
    resource_type: "attempt",
    resource_id: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    action_url: "/results/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
    action_label: "Ver Resultados",
    metadata: {
      score: 100,
      total_questions: 2
    }
  },
  delivery: {
    in_app: {
      enabled: true,
      delivered_at: new Date("2025-01-15T10:16:15Z")
    },
    push: {
      enabled: true,
      sent_at: new Date("2025-01-15T10:16:16Z"),
      delivered_at: new Date("2025-01-15T10:16:17Z")
    },
    email: {
      enabled: true,
      sent_at: new Date("2025-01-15T10:16:18Z"),
      delivered_at: new Date("2025-01-15T10:16:25Z")
    }
  },
  is_read: true,
  read_at: new Date("2025-01-15T10:20:00Z"),
  is_archived: false,
  created_at: new Date("2025-01-15T10:16:15Z")
});

// Notification 3 - Material uploaded
db.notifications.insertOne({
  user_id: "33333333-3333-3333-3333-333333333333",
  notification_type: "material.uploaded",
  title: "Nuevo Material Disponible",
  message: "El Prof. García ha subido un nuevo material: Álgebra Lineal - Matrices",
  priority: "low",
  category: "academic",
  data: {
    resource_type: "material",
    resource_id: "77777777-7777-7777-7777-777777777777",
    action_url: "/materials/77777777-7777-7777-7777-777777777777",
    action_label: "Ver Material",
    metadata: {
      subject: "Matemáticas",
      teacher_name: "Prof. García"
    }
  },
  delivery: {
    in_app: {
      enabled: true,
      delivered_at: new Date("2025-01-11T14:20:00Z")
    },
    push: {
      enabled: false
    },
    email: {
      enabled: false
    }
  },
  is_read: false,
  is_archived: false,
  created_at: new Date("2025-01-11T14:20:00Z")
});

// Notification 4 - System announcement
db.notifications.insertOne({
  user_id: "11111111-1111-1111-1111-111111111111",
  notification_type: "system.announcement",
  title: "Mantenimiento Programado",
  message: "El sistema estará en mantenimiento el sábado 18 de enero de 02:00 a 04:00 hrs.",
  priority: "urgent",
  category: "system",
  data: {
    metadata: {
      maintenance_start: "2025-01-18T02:00:00Z",
      maintenance_end: "2025-01-18T04:00:00Z",
      services_affected: ["assessments", "materials"]
    }
  },
  delivery: {
    in_app: {
      enabled: true,
      delivered_at: new Date("2025-01-15T08:00:00Z")
    },
    push: {
      enabled: true,
      sent_at: new Date("2025-01-15T08:00:01Z"),
      delivered_at: new Date("2025-01-15T08:00:02Z")
    },
    email: {
      enabled: true,
      sent_at: new Date("2025-01-15T08:00:03Z"),
      delivered_at: new Date("2025-01-15T08:00:10Z")
    }
  },
  is_read: true,
  read_at: new Date("2025-01-15T09:00:00Z"),
  is_archived: false,
  expires_at: new Date("2025-01-18T05:00:00Z"),
  created_at: new Date("2025-01-15T08:00:00Z")
});

print("✅ 4 notifications documents inserted successfully");

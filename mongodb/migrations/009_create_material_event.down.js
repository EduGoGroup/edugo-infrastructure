// 009_create_material_event.down.js
// Rollback: Eliminar collection material_event

print("🗑️  Dropping collection: material_event...");

db.material_event.drop();

print("✅ Collection material_event dropped successfully");

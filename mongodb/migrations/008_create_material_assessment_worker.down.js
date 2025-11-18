// 008_create_material_assessment_worker.down.js
// Rollback: Eliminar collection material_assessment_worker

print("🗑️  Dropping collection: material_assessment_worker...");

db.material_assessment_worker.drop();

print("✅ Collection material_assessment_worker dropped successfully");

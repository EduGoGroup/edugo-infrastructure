// 007_create_material_summary.down.js
// Rollback: Eliminar collection material_summary

print("🗑️  Dropping collection: material_summary...");

db.material_summary.drop();

print("✅ Collection material_summary dropped successfully");

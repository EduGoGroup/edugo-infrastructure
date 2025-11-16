// Migration DOWN: Drop assessment_attempt_result collection

db.assessment_attempt_result.drop();

print("✅ Collection 'assessment_attempt_result' dropped successfully");

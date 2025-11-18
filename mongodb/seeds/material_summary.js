// Seeds para material_summary
// Datos de prueba para resúmenes generados por IA

db = db.getSiblingDB("edugo");

print("🌱 Seeding material_summary...");

db.material_summary.insertMany([
  {
    material_id: "a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d",
    summary: "Este material cubre los fundamentos de la programación orientada a objetos en Java. Se explican conceptos clave como clases, objetos, herencia, polimorfismo y encapsulación con ejemplos prácticos.",
    key_points: [
      "Introducción a POO y sus principios fundamentales",
      "Clases y objetos: definición y uso",
      "Herencia y polimorfismo en Java",
      "Encapsulación y modificadores de acceso",
      "Ejemplos prácticos con código"
    ],
    language: "es",
    word_count: 42,
    version: 1,
    ai_model: "gpt-4",
    processing_time_ms: 3500,
    token_usage: {
      prompt_tokens: 850,
      completion_tokens: 180,
      total_tokens: 1030
    },
    metadata: {
      source_length: 5420,
      has_images: false
    },
    created_at: new Date("2025-11-15T10:30:00Z"),
    updated_at: new Date("2025-11-15T10:30:00Z")
  },
  {
    material_id: "f1a2b3c4-d5e6-4f5a-9b8c-7d6e5f4a3b2c",
    summary: "A comprehensive guide to React Hooks covering useState, useEffect, useContext, and custom hooks. Learn how to manage state and side effects in functional components effectively.",
    key_points: [
      "Introduction to React Hooks and their benefits",
      "useState for state management",
      "useEffect for side effects and lifecycle",
      "useContext for global state sharing",
      "Creating custom hooks for reusable logic"
    ],
    language: "en",
    word_count: 38,
    version: 1,
    ai_model: "gpt-4-turbo",
    processing_time_ms: 2800,
    token_usage: {
      prompt_tokens: 920,
      completion_tokens: 165,
      total_tokens: 1085
    },
    created_at: new Date("2025-11-16T14:20:00Z"),
    updated_at: new Date("2025-11-16T14:20:00Z")
  },
  {
    material_id: "b2c3d4e5-f6a7-4b5c-8d9e-0f1a2b3c4d5e",
    summary: "Material sobre estruturas de dados fundamentais: arrays, listas encadeadas, pilhas e filas. Inclui análise de complexidade e implementações práticas em Python.",
    key_points: [
      "Arrays e suas operações básicas",
      "Listas encadeadas: simples e duplas",
      "Pilhas (LIFO) e suas aplicações",
      "Filas (FIFO) e variantes",
      "Análise de complexidade temporal e espacial"
    ],
    language: "pt",
    word_count: 35,
    version: 1,
    ai_model: "gpt-4o",
    processing_time_ms: 3100,
    token_usage: {
      prompt_tokens: 780,
      completion_tokens: 155,
      total_tokens: 935
    },
    created_at: new Date("2025-11-17T09:45:00Z"),
    updated_at: new Date("2025-11-17T09:45:00Z")
  }
]);

print("✅ 3 material_summary documents inserted");

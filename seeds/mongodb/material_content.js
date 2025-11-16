// Seeds for material_content collection
// Execute with: mongosh --host localhost:27017/edugo < material_content.js

use edugo;

// Material content 1 (for Physics material)
db.material_content.insertOne({
  material_id: "66666666-6666-6666-6666-666666666666",
  content_type: "pdf_extracted",
  raw_text: "Física Cuántica - Introducción\n\nLa física cuántica es una de las ramas más fascinantes de la física moderna. En este documento exploraremos conceptos fundamentales como la dualidad onda-partícula y el principio de incertidumbre de Heisenberg.\n\n1. Dualidad Onda-Partícula\nLas partículas subatómicas exhiben propiedades tanto de ondas como de partículas. Este fenómeno fue observado por primera vez en experimentos con electrones.\n\n2. Principio de Incertidumbre\nHeisenberg estableció que no podemos conocer simultáneamente con precisión absoluta la posición y el momento de una partícula.",
  structured_content: {
    title: "Física Cuántica - Introducción",
    sections: [
      {
        section_index: 0,
        heading: "Introducción",
        content: "La física cuántica es una de las ramas más fascinantes de la física moderna.",
        page_number: 1
      },
      {
        section_index: 1,
        heading: "Dualidad Onda-Partícula",
        content: "Las partículas subatómicas exhiben propiedades tanto de ondas como de partículas.",
        page_number: 1
      },
      {
        section_index: 2,
        heading: "Principio de Incertidumbre",
        content: "Heisenberg estableció que no podemos conocer simultáneamente con precisión absoluta la posición y el momento de una partícula.",
        page_number: 2
      }
    ],
    summary: "Este material cubre conceptos fundamentales de física cuántica, incluyendo la dualidad onda-partícula y el principio de incertidumbre de Heisenberg.",
    key_concepts: [
      "Dualidad onda-partícula",
      "Principio de incertidumbre",
      "Física cuántica",
      "Heisenberg",
      "Partículas subatómicas"
    ]
  },
  processing_info: {
    processor_version: "v1.2.0",
    processed_at: new Date("2025-01-10T10:30:00Z"),
    processing_duration_ms: 1250,
    page_count: 2,
    word_count: 156
  },
  created_at: new Date("2025-01-10T10:30:00Z"),
  updated_at: new Date("2025-01-10T10:30:00Z")
});

// Material content 2 (for Algebra material)
db.material_content.insertOne({
  material_id: "77777777-7777-7777-7777-777777777777",
  content_type: "pdf_extracted",
  raw_text: "Álgebra Lineal - Matrices\n\nLas matrices son arreglos rectangulares de números que tienen múltiples aplicaciones en matemáticas y ciencias de la computación.\n\n1. Matriz Identidad\nUna matriz identidad es una matriz cuadrada que tiene 1s en la diagonal principal y 0s en el resto de las posiciones.\n\n2. Operaciones con Matrices\nLas matrices se pueden sumar, restar y multiplicar siguiendo reglas específicas.",
  structured_content: {
    title: "Álgebra Lineal - Matrices",
    sections: [
      {
        section_index: 0,
        heading: "Introducción a Matrices",
        content: "Las matrices son arreglos rectangulares de números que tienen múltiples aplicaciones.",
        page_number: 1
      },
      {
        section_index: 1,
        heading: "Matriz Identidad",
        content: "Una matriz identidad es una matriz cuadrada que tiene 1s en la diagonal principal y 0s en el resto.",
        page_number: 1
      },
      {
        section_index: 2,
        heading: "Operaciones con Matrices",
        content: "Las matrices se pueden sumar, restar y multiplicar siguiendo reglas específicas.",
        page_number: 2
      }
    ],
    summary: "Introducción a matrices, incluyendo la definición de matriz identidad y operaciones básicas.",
    key_concepts: [
      "Matrices",
      "Matriz identidad",
      "Álgebra lineal",
      "Operaciones matriciales"
    ]
  },
  processing_info: {
    processor_version: "v1.2.0",
    processed_at: new Date("2025-01-11T14:20:00Z"),
    processing_duration_ms: 980,
    page_count: 2,
    word_count: 98
  },
  created_at: new Date("2025-01-11T14:20:00Z"),
  updated_at: new Date("2025-01-11T14:20:00Z")
});

print("✅ 2 material_content documents inserted successfully");

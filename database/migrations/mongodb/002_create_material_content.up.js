// Migration: Create material_content collection
// Collection: material_content (Owner: infrastructure)
// Created by: edugo-infrastructure
// Used by: api-mobile, worker
//
// Purpose: Stores processed content from educational materials (extracted text, parsed structure)
// Related PostgreSQL table: materials (stores file metadata, references this via material_id)

db.createCollection("material_content", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["material_id", "content_type", "created_at", "updated_at"],
      properties: {
        material_id: {
          bsonType: "string",
          description: "UUID of the material in PostgreSQL (materials table)"
        },
        content_type: {
          bsonType: "string",
          enum: ["pdf_extracted", "video_transcript", "document_parsed", "slides_extracted"],
          description: "Type of content extraction"
        },
        raw_text: {
          bsonType: "string",
          description: "Raw extracted text from the material"
        },
        structured_content: {
          bsonType: "object",
          description: "Structured/parsed content with sections and metadata",
          properties: {
            title: {
              bsonType: "string",
              description: "Title of the material"
            },
            sections: {
              bsonType: "array",
              description: "Content sections",
              items: {
                bsonType: "object",
                properties: {
                  section_index: {
                    bsonType: "int",
                    minimum: 0
                  },
                  heading: {
                    bsonType: "string"
                  },
                  content: {
                    bsonType: "string"
                  },
                  page_number: {
                    bsonType: "int"
                  }
                }
              }
            },
            summary: {
              bsonType: "string",
              description: "AI-generated summary of the content"
            },
            key_concepts: {
              bsonType: "array",
              description: "Key concepts extracted from the material",
              items: {
                bsonType: "string"
              }
            }
          }
        },
        processing_info: {
          bsonType: "object",
          description: "Information about the processing",
          properties: {
            processor_version: {
              bsonType: "string",
              description: "Version of the content processor used"
            },
            processed_at: {
              bsonType: "date",
              description: "When the content was processed"
            },
            processing_duration_ms: {
              bsonType: "int",
              minimum: 0,
              description: "Processing duration in milliseconds"
            },
            page_count: {
              bsonType: "int",
              minimum: 0,
              description: "Number of pages processed"
            },
            word_count: {
              bsonType: "int",
              minimum: 0,
              description: "Word count in the content"
            }
          }
        },
        created_at: {
          bsonType: "date",
          description: "Timestamp when the content was created"
        },
        updated_at: {
          bsonType: "date",
          description: "Timestamp when the content was last updated"
        }
      }
    }
  }
});

// Create indexes for efficient queries
db.material_content.createIndex({ "material_id": 1 }, { name: "idx_material_id", unique: true });
db.material_content.createIndex({ "content_type": 1 }, { name: "idx_content_type" });
db.material_content.createIndex({ "created_at": -1 }, { name: "idx_created_at_desc" });
db.material_content.createIndex({ "processing_info.processed_at": -1 }, { name: "idx_processed_at_desc" });

// Text index for full-text search on raw_text and structured_content
db.material_content.createIndex(
  {
    "raw_text": "text",
    "structured_content.summary": "text",
    "structured_content.key_concepts": "text"
  },
  {
    name: "idx_fulltext_search",
    default_language: "spanish"
  }
);

print("✅ Collection 'material_content' created successfully");

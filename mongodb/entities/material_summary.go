package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MaterialSummary representa la collection 'material_summary' en MongoDB.
// Esta entity es el reflejo exacto del schema de BD.
//
// Seed: mongodb/seeds/material_summary.js
// Usada por: worker
//
// Nota: Contiene el resumen del material generado por IA.
type MaterialSummary struct {
	ID               bson.ObjectID    `bson:"_id,omitempty"`
	MaterialID       string           `bson:"material_id"`        // UUID del material en PostgreSQL
	Summary          string           `bson:"summary"`            // Resumen principal
	KeyPoints        []string         `bson:"key_points"`         // Puntos clave
	Language         string           `bson:"language"`           // es, en, pt, etc
	WordCount        int              `bson:"word_count"`         // Conteo de palabras del resumen
	Version          int              `bson:"version"`            // Versión del resumen
	AIModel          string           `bson:"ai_model"`           // Modelo IA usado
	ProcessingTimeMs int              `bson:"processing_time_ms"` // Tiempo de procesamiento en ms
	TokenUsage       *TokenUsage      `bson:"token_usage,omitempty"`
	Metadata         *SummaryMetadata `bson:"metadata,omitempty"`           // Metadata adicional (opcional)
	Sections         []DocumentSection `bson:"sections,omitempty" json:"sections,omitempty"` // Secciones extraídas del documento
	CreatedAt        time.Time        `bson:"created_at"`
	UpdatedAt        time.Time        `bson:"updated_at"`
}

// SummaryMetadata contiene metadata adicional extensible
type SummaryMetadata struct {
	SourceLength int  `bson:"source_length,omitempty"` // Longitud del material fuente
	HasImages    bool `bson:"has_images,omitempty"`    // Si el material tiene imágenes
}

// DocumentSection representa una sección extraída de un documento/material.
type DocumentSection struct {
	Index   int    `bson:"index" json:"index"`
	Title   string `bson:"title" json:"title"`
	Content string `bson:"content" json:"content"`
	Preview string `bson:"preview" json:"preview"` // Primeras ~200 palabras
}

// CollectionName retorna el nombre de la collection en MongoDB
func (MaterialSummary) CollectionName() string {
	return "material_summary"
}

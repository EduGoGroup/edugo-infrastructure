package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MaterialEvent representa la collection 'material_event' en MongoDB.
// Esta entity es el reflejo exacto del schema de BD.
//
// Seed: mongodb/seeds/material_event.js
// Usada por: worker
//
// Nota: Eventos de auditoría del procesamiento de materiales.
type MaterialEvent struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	EventType   string             `bson:"event_type"`   // material_uploaded, material_reprocess, assessment_attempt
	MaterialID  string             `bson:"material_id"`  // UUID del material en PostgreSQL
	UserID      string             `bson:"user_id"`      // UUID del usuario en PostgreSQL
	Payload     bson.M        `bson:"payload"`      // Datos del evento (flexible)
	Status      string             `bson:"status"`       // completed, processing, failed
	ErrorMsg    *string            `bson:"error_msg,omitempty"`
	StackTrace  *string            `bson:"stack_trace,omitempty"`
	RetryCount  int                `bson:"retry_count"`
	ProcessedAt *time.Time         `bson:"processed_at,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

// CollectionName retorna el nombre de la collection en MongoDB
func (MaterialEvent) CollectionName() string {
	return "material_event"
}

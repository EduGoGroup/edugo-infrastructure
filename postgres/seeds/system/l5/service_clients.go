package l5

import (
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ScopeNotificationsDispatch es el scope M2M requerido para POST /internal/notifications/dispatch.
const ScopeNotificationsDispatch = "notifications.dispatch"

// devServiceClientSecretHash es bcrypt(DefaultCost) de "change-me-dev-only"
// (valor en push-secrets.env / SERVICE_CLIENT_SECRET; nunca persiste en claro).
// Método: bcrypt vía golang.org/x/crypto/bcrypt (mismo algoritmo que auth.users.password_hash).
// Precomputado para determinismo del seed; en runtime la validación usa bcrypt.CompareHashAndPassword.
const devServiceClientSecretHash = "$2a$10$64yROwHnbokxE7wrWGCzzeMYfierzsRq53iZVI7l.qbSe74gf/C7K"

const (
	serviceClientWorkerID   = "c5000000-0000-0000-0000-000000000001"
	serviceClientLearningID = "c5000000-0000-0000-0000-000000000002"
)

// ApplyServiceClients inserta los clientes M2M iniciales (idempotente ON CONFLICT DO NOTHING).
func ApplyServiceClients(tx *gorm.DB) error {
	workerDesc := "Worker async — delega push al Notification Gateway"
	learningDesc := "API learning — productor de eventos de evaluación"

	rows := []entities.ServiceClient{
		{
			ID:          uuid.MustParse(serviceClientWorkerID),
			ClientID:    "edugo-worker",
			SecretHash:  devServiceClientSecretHash,
			Scopes:      pq.StringArray{ScopeNotificationsDispatch},
			IsActive:    true,
			Description: &workerDesc,
		},
		{
			ID:          uuid.MustParse(serviceClientLearningID),
			ClientID:    "edugo-api-learning",
			SecretHash:  devServiceClientSecretHash,
			Scopes:      pq.StringArray{ScopeNotificationsDispatch},
			IsActive:    true,
			Description: &learningDesc,
		},
	}

	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "client_id"}},
		DoNothing: true,
	}).Create(&rows).Error
}

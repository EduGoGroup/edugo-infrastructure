package dbutil

import (
	"fmt"
	"os"
)

// BuildDBURL construye una URL de conexión PostgreSQL desde variables de entorno.
// Soporta DATABASE_URL, DB_* y POSTGRES_* como fuentes.
func BuildDBURL() string {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	host := EnvFirst("DB_HOST", "POSTGRES_HOST", "localhost")
	port := EnvFirst("DB_PORT", "POSTGRES_PORT", "5432")
	name := EnvFirst("DB_NAME", "POSTGRES_DB", "edugo_dev")
	user := EnvFirst("DB_USER", "POSTGRES_USER", "edugo")
	password := EnvFirst("DB_PASSWORD", "POSTGRES_PASSWORD", "changeme")
	sslmode := EnvFirst("DB_SSL_MODE", "POSTGRES_SSLMODE", "disable")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslmode)
}

// EnvFirst retorna el valor de la primera variable de entorno no vacía,
// o fallback si ninguna está definida.
func EnvFirst(primary, secondary, fallback string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	if value := os.Getenv(secondary); value != "" {
		return value
	}
	return fallback
}

package mongodbutil

import (
	"fmt"
	"os"
)

// BuildMongoURI construye una URI de conexión MongoDB desde variables de entorno.
// Soporta MONGO_URI directo, o MONGO_HOST/PORT/USER/PASSWORD como alternativa.
func BuildMongoURI() string {
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return uri
	}

	host := EnvFirst("MONGO_HOST", "", "localhost")
	port := EnvFirst("MONGO_PORT", "", "27017")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")

	if user != "" && password != "" {
		return fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", user, password, host, port)
	}

	return fmt.Sprintf("mongodb://%s:%s", host, port)
}

// EnvFirst retorna el valor de la primera variable de entorno no vacía,
// o fallback si ninguna está definida. Si secondary es "", se omite.
func EnvFirst(primary, secondary, fallback string) string {
	if value := os.Getenv(primary); value != "" {
		return value
	}
	if secondary != "" {
		if value := os.Getenv(secondary); value != "" {
			return value
		}
	}
	return fallback
}

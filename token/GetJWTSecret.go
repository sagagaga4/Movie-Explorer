package token

import (
	"os"

	logger "SagiProjects.com/moviesite/Logger"
)

func GetJWTSecret(logger logger.Logger) string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-for-dev"
		logger.Info("JWT_SECRET not set, using default development secret")
	} else {
		logger.Info("Using JWT_SECRET from environment")
	}
	return jwtSecret
}

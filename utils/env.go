package utils

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	ApiPort    string
	DbHost     string
	DbUser     string
	DbPass     string
	DbName     string
	DbPort     string
	SigningKey []byte
)

func LoadEnv() {
	godotenv.Load(".env")
	ApiPort = getEnv("PORT", "8000")

	DbHost = getEnv("POSTGRES_HOST", "localhost")
	DbUser = getEnv("POSTGRES_USER", "user")
	DbPass = getEnv("POSTGRES_PASSWORD", "password")
	DbName = getEnv("POSTGRES_DB", "postgres")
	DbPort = getEnv("POSTGRES_PORT", "5432")

	secretKey := getEnv("SECRET_KEY", "secret")
	SigningKey = []byte(secretKey)

	herokuUrl := getEnv("DATABASE_URL", "")
	if herokuUrl != "" {
		formatHerokuDB(herokuUrl)
	}
}

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func formatHerokuDB(herokuUrl string) {
	t := strings.Split(herokuUrl, "://")
	if len(t) != 2 {
		return
	}
	userAnd := strings.Split(t[1], ":")
	if len(userAnd) != 3 {
		return
	}
	portAndDb := strings.Split(userAnd[2], "/")
	if len(portAndDb) != 2 {
		return
	}
	passAndHost := strings.Split(userAnd[1], "@")
	if len(passAndHost) != 2 {
		return
	}
	DbHost = passAndHost[1]
	DbUser = userAnd[0]
	DbPass = passAndHost[0]
	DbName = portAndDb[1]
	DbPort = portAndDb[0]
}

package common

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Cfg Config

type Config struct {
	DBDriver      string
	DBSource      string
	JwtSecretKey  string
	IsProduction  bool
	EncryptionKey string
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Did not find .env file.")
	}
	Cfg = Config{
		DBDriver:      loadEnv("DB_DRIVER"),
		DBSource:      loadEnv("DB_SOURCE"),
		JwtSecretKey:  loadEnv("JWT_SECRET_KEY"),
		EncryptionKey: loadEnv("ENCRYPTION_KEY"),
		IsProduction:  loadEnv("IS_PRODUCTION") == "true",
	}
}

func loadEnv(envVarName string) string {
	env, exists := os.LookupEnv(envVarName)
	if !exists {
		panic(fmt.Sprintf("Env variable %s cannot be loaded.", envVarName))
	}
	return env
}

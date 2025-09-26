// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"strconv"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func InitArgon2() {

}

type Configuration struct {
	JWTSecret  []byte
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	RedisAddr  string
	RedisPass  string
	RedisDB    int
	Argon2     *Argon2Params
}

var Config *Configuration

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY is not set")
	}

	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	Config = &Configuration{
		JWTSecret:  []byte(secret),
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		RedisAddr:  os.Getenv("REDIS_ADDRESS"),
		RedisPass:  os.Getenv("REDIS_PASSWORD"),
		RedisDB:    redisDB,
		Argon2: &Argon2Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
}

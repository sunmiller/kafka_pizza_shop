package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sunmiller/pizza-shop-eda/order-service/logger"
)

var env ConfigDto

type ConfigDto struct {
	Port              string
	DatabaseURL       string
	DatabaseName      string
	KafkaHost         string
	KafkaPort         string
	KafkaDefaultTopic string
	KafkaGroupID      string
}

func init() {
	if env.Port == "" {
		ConfigEnv()
	}
}

func ConfigEnv() {
	LoadEnvVariable()
	env = ConfigDto{
		Port:              os.Getenv("PORT"),
		DatabaseURL:       os.Getenv("MONGO_DB_URL"),
		DatabaseName:      os.Getenv("MONGO_DB_NAME"),
		KafkaHost:         os.Getenv("KAFKA_HOST"),
		KafkaPort:         os.Getenv("KAFKA_PORT"),
		KafkaDefaultTopic: os.Getenv("KAFKA_DEFAULT_TOPIC"),
		KafkaGroupID:      os.Getenv("KAFKA_GROUP_ID"),
	}

}
func LoadEnvVariable() {
	var envFile string

	if isRunningInDocker() {
		envFile = ".env" // default used in Docker (copied or mounted in Dockerfile/compose)
	} else {
		envFile = ".env.dev" // used locally during debugging
	}
	fmt.Printf("Found env file, %v", envFile)
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("❌ Failed to load %s: %v", envFile, err)
		}
		logger.Log(fmt.Sprintf("✅ Loaded environment from %s", envFile))
	} else if os.IsNotExist(err) {
		logger.Log(fmt.Sprintf("⚠️  %s not found, falling back to system environment variables", envFile))
	} else {
		log.Fatalf("❌ Error checking for %s: %v", envFile, err)
	}
}

func isRunningInDocker() bool {
	data, err := os.ReadFile("/proc/1/cgroup")
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "docker") || strings.Contains(string(data), "containerd")
}

func accessField(key string) (string, error) {
	v := reflect.ValueOf(env)
	t := v.Type()

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct got %v", t)
	}

	_, ok := t.FieldByName(key)
	if !ok {
		return "", fmt.Errorf("key %v does not exist", key)
	}

	fv := v.FieldByName(key)
	return fv.String(), nil
}

func GetEnvProperty(key string) string {
	logger.Log(fmt.Sprintf("you asked for key: %v", key))
	if env.Port == "" {
		ConfigEnv()
	}
	val, err := accessField(key)
	if err != nil {
		logger.Log(fmt.Sprintf("error loading .env file: %v", err))
	}
	if val != "" {
		logger.Log(fmt.Sprintf("I found value %v for key: %v", val, key))
	} else {
		logger.Log(fmt.Sprintf("Error Not found value %v for key: %v", val, key))
	}

	return val
}

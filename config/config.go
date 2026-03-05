package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	// only load environment for development
	// in production, env vars are set by the deployment environment
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file Found: ", err)
	}
}

func GetConnectionString() (string, string) {
	driver := os.Getenv("DB_DRIVER")
	user := os.Getenv("DB_USER")
	// password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	return driver, fmt.Sprintf("%s:@tcp(%s:%s)/%s", user, host, port, name)
}

func GetMainPort() string {
	// local server port
	appPort := os.Getenv("APP_PORT")
	return fmt.Sprintf(":%s", appPort)
}

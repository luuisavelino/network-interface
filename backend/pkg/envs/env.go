package envs

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Database  database
	JwtSecret string
)

type database struct {
	HOST        string
	PORT        string
	DB_NAME     string
	USER        string
	PASSWORD    string
	MAX_RETRIES int
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	maxRetries, err := strconv.Atoi(os.Getenv("MAX_RETRIES"))
	if err != nil {
		maxRetries = 5
	}

	Database = database{
		HOST:        os.Getenv("DB_HOST"),
		PORT:        os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
		USER:        os.Getenv("DB_USER"),
		PASSWORD:    os.Getenv("DB_PASSWORD"),
		MAX_RETRIES: maxRetries,
	}

	JwtSecret = os.Getenv("JWT_SECRET")

	if Database.HOST == "" {
		fmt.Println(errors.New("database host must not be empty"))
		panic("database host must not be empty")
	}

	if Database.DB_NAME == "" {
		fmt.Println(errors.New("database name must not be empty"))
		panic("database name must not be empty")
	}

	if Database.PORT == "" {
		fmt.Println(errors.New("database port must not be empty"))
		panic("database port must not be empty")
	}

	if Database.USER == "" {
		fmt.Println(errors.New("database user must not be empty"))
		panic("database user must not be empty")
	}

	if Database.PASSWORD == "" {
		fmt.Println(errors.New("database password must not be empty"))
		panic("database password must not be empty")
	}

	if JwtSecret == "" {
		fmt.Println(errors.New("jwt secret must not be empty"))
		panic("jwt secret must not be empty")
	}
}

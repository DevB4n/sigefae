package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	// Graph
	GraphClientSecret string
	GraphClientID     string
	GraphTenantID     string
	GraphUserEmail    string

	// Database
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBCharset   string
	DBParseTime string
	DBLoc       string
}
func Load(path string) (*Env, error) {
	if err := godotenv.Load(path); err != nil {
		return nil, err
	}

	cfg := &Env{
		// Graph
		GraphClientID:     os.Getenv("GRAPH_CLIENT_ID"),
		GraphClientSecret: os.Getenv("GRAPH_CLIENT_SECRET"),
		GraphTenantID:     os.Getenv("GRAPH_TENANT_ID"),
		GraphUserEmail:    os.Getenv("GRAPH_USER_EMAIL"),

		// DB
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBName:      os.Getenv("DB_NAME"),
		DBCharset:   os.Getenv("DB_CHARSET"),
		DBParseTime: os.Getenv("DB_PARSE_TIME"),
		DBLoc:       os.Getenv("DB_LOC"),
	}

	if cfg.GraphClientID == "" ||
		cfg.GraphClientSecret == "" ||
		cfg.GraphTenantID == "" ||
		cfg.GraphUserEmail == "" {
		return nil, fmt.Errorf("missing graph environment variables")
	}

	if cfg.DBHost == "" ||
		cfg.DBUser == "" ||
		cfg.DBPassword == "" ||
		cfg.DBName == "" {
		return nil, fmt.Errorf("missing database environment variables")
	}

	return cfg, nil
}
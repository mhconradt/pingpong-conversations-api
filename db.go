package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
)

const (
	DefaultUser     = "maxwellconradt"
	DefaultDatabase = "pingpong"
	DefaultSslMode  = "disable"
	DefaultHost     = "localhost"
	DefaultPassword = "password"
)

func LookupWithDefault(k, dv string) string {
	if v, found := os.LookupEnv(k); !found {
		return dv
	} else {
		return v
	}
}

func GetConnectionString() string {
	user := LookupWithDefault("DATABASE_USER", DefaultUser)
	database := LookupWithDefault("DATABASE_NAME", DefaultDatabase)
	host := LookupWithDefault("DATABASE_HOST", DefaultHost)
	password := LookupWithDefault("DATABASE_PASSWORD", DefaultPassword)
	sslMode := LookupWithDefault("DATABASE_SSL_MODE", DefaultSslMode)
	return fmt.Sprintf("database=%v host=%v user=%v sslmode=%v password=%v connect_timeout=5", database, host, user, sslMode, password)
}

func InitDB() (*gorm.DB, error) {
	cs := GetConnectionString()
	return gorm.Open("postgres", cs)
}

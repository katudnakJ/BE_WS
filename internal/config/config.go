// ❤ ❤ ยังไม่ได้ เติม API token นะจ๊ะ อุอิ ง่ายๆ ก็คือยังไม่สมบูรณ์ นั้นแหละ ถ้า .env แก้แล้วจะมาแก้ให้ เพราะมันดึงจาก .env เอาคราวๆประมาณนี้ไปก่อนนะ❤ ❤
package config

import (
	"os"
)

type Config struct {
	POSTGRESDB             string
	POSTGRESUSER           string
	POSTGRESPASSWORD       string
	POSTGRESPORT           string
	PGADMINDEFAULTEMAIL    string
	PGADMINDEFAULTPASSWORD string
	PGADMINPORT            string
}

func LoadConfig() Config {
	return Config{
		POSTGRESDB:             getEnv("POSTGRES_DB", "webservice"),
		POSTGRESUSER:           getEnv("POSTGRES_USER", "webservice_user"),
		POSTGRESPASSWORD:       getEnv("POSTGRES_PASSWORD", "your_strong_password"),
		POSTGRESPORT:           getEnv("POSTGRES_PORT", "5432"),
		PGADMINDEFAULTEMAIL:    getEnv("PGADMIN_DEFAULT_EMAIL", "nuttachot@hotmail.com"),
		PGADMINDEFAULTPASSWORD: getEnv("PGADMIN_DEFAULT_PASSWORD", "password"),
		PGADMINPORT:            getEnv("PGADMIN_PORT", "5050"),
	}
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

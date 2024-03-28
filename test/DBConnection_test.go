package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"os/exec"
	"testing"
	"time"
)

type DBConfig struct {
	host, port, user, password, name string
}

func NewDBConfig(host string, port string, user string, password string, name string) *DBConfig {
	return &DBConfig{host: host, port: port, user: user, password: password, name: name}
}

func PGStore(config *DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=America/New_York", config.host, config.user, config.password, config.name, config.port)
	//Best Performance Config GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	return db, nil
}

func TestDBEnv(t *testing.T) {
	// test cases for testing case
	testCases := []struct {
		name        string
		env         string
		expectedEnv string
	}{
		{name: "db_host", env: os.Getenv("PG_DB_HOST"), expectedEnv: "localhost"},
		{name: "db_host", env: os.Getenv("PG_DB_PORT"), expectedEnv: "5432"},
		{name: "db_user", env: os.Getenv("PG_DB_USER"), expectedEnv: "postgres"},
		{name: "db_password", env: os.Getenv("PG_DB_PASSWORD"), expectedEnv: "postgres"},
		{name: "db_name", env: os.Getenv("PG_DB_NAME"), expectedEnv: "postgres"},
	}

	for _, testCase := range testCases {
		// run testing every testCase
		t.Run(testCase.name, func(t *testing.T) {
			fmt.Println(testCase.env)
			require.Equal(t, testCase.expectedEnv, testCase.env)
		})
	}
}

func TestDBConnection(t *testing.T) {
	config := NewDBConfig(os.Getenv("PG_DB_HOST"), os.Getenv("PG_DB_PORT"), os.Getenv("PG_DB_USER"), os.Getenv("PG_DB_PASSWORD"), os.Getenv("PG_DB_NAME"))
	_, err := PGStore(config)
	require.NoError(t, err)
	err = exec.Command("echo", "Finish TDD").Run()
	require.NoError(t, err)
}

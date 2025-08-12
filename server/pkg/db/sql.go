package db

import (
	"fmt"
	"github.com/lzkking/harle/server/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func postgresClient(dbConfig *config.DatabaseConfig) *gorm.DB {
	dsn, err := dbConfig.DSN()
	if err != nil {
		panic(err)
	}
	dbClient, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      getGormLogger(dbConfig),
	})
	if err != nil {
		panic(err)
	}
	return dbClient
}

func mySQLClient(dbConfig *config.DatabaseConfig) *gorm.DB {
	dsn, err := dbConfig.DSN()
	if err != nil {
		panic(err)
	}
	dbClient, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      getGormLogger(dbConfig),
	})
	if err != nil {
		panic(err)
	}
	return dbClient
}

func newDBClient() *gorm.DB {
	dbConfig := config.GetServerConfig().DbConfig
	var dbClient *gorm.DB

	switch dbConfig.Dialect {
	case config.SQLITE:
		panic(fmt.Sprintf("SQLite is currently not supported"))
	case config.POSTGRES:
		dbClient = postgresClient(dbConfig)
	case config.MYSQL:
		dbClient = mySQLClient(dbConfig)
	default:
		panic(fmt.Sprintf("Unknown DB Dialect: '%s'", dbConfig.Dialect))
	}

	//迁移数据库
	var allDBModels = append(make([]interface{}, 0))
	var err error
	//	迁移所有的表
	for _, model := range allDBModels {
		err = dbClient.AutoMigrate(model)
		if err != nil {
			panic(err)
		}
	}

	sqlDB, err := dbClient.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)

	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	sqlDB.SetConnMaxLifetime(time.Hour)

	return dbClient
}

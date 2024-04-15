package database

import (
	"errors"
	"strconv"
	"strings"

	"blue-beetle/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var UserRepo UserRepository

type UserRepository struct {
	Database *gorm.DB
}

func (r *UserRepository) AutoMigrate() error {
	err := r.Database.AutoMigrate(Role{})
	if err != nil {
		return err
	}
	err = r.Database.AutoMigrate(User{})
	if err != nil {
		return err
	}
}

func (r *UserRepository) ConnectUserRepository(dbconfig config.DBConfig) error {
	err := checkDBConfig(&dbconfig)
	if err != nil {
		return errors.New("Database configuration not correct: " + err.Error())
	}
	if dbconfig.Type == "mysql" {
		connectionString := buildMySQLConnectionString(dbconfig)
		db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
		if err != nil {
			return err
		}
		r.Database = db
	} else if dbconfig.Type == "sqlite" {
		db, err := gorm.Open(sqlite.Open(dbconfig.FilePath), &gorm.Config{})
		if err != nil {
			return err
		}
		r.Database = db
	} else if dbconfig.Type == "postgresql" {
		connectionString := buildPostGresqlConnectionString(dbconfig)
		db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
		if err != nil {
			return err
		}
		r.Database = db
	} else if dbconfig.Type == "sql" {
		connectionString := buildSQLServerConnectionString(dbconfig)
		db, err := gorm.Open(sqlserver.Open(connectionString), &gorm.Config{})
		if err != nil {
			return err
		}
		r.Database = db
	}
	return nil
}

func checkDBConfig(dbconfig *config.DBConfig) error {
	if len(dbconfig.Type) == 0 {
		return errors.New("no database type defined in the database config file")
	}
	dbconfig.Type = strings.ToLower(dbconfig.Type)
	if dbconfig.Type == "sqlite" {
		if len(dbconfig.FilePath) == 0 {
			return errors.New("sqlite database types must define a file path")
		}
	} else if dbconfig.Type == "mysql" || dbconfig.Type == "postgresql" || dbconfig.Type == "sql" {
		if len(dbconfig.Username) == 0 {
			return errors.New("missing username for datbase")
		}
		if len(dbconfig.Password) == 0 {
			return errors.New("missing password for datbase")
		}
		if len(dbconfig.Server) == 0 {
			return errors.New("missing server for datbase")
		}
		if len(dbconfig.Protocol) == 0 && dbconfig.Type == "mysql" {
			return errors.New("missing protocol for datbase")
		}
		if dbconfig.Port == 0 {
			return errors.New("missing port for datbase")
		}
		if len(dbconfig.DBName) == 0 {
			return errors.New("missing databse name for datbase")
		}
	} else {
		return errors.New("unsupported database type")
	}
	return nil
}

func buildSQLServerConnectionString(dbconfig config.DBConfig) string {
	var con string
	// sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm
	con = "sqlserver://" + dbconfig.Username + ":" + dbconfig.Password +
		"@" + dbconfig.Server + ":" + strconv.Itoa(dbconfig.Port) + "?database=" + dbconfig.DBName
	if len(dbconfig.Options) > 0 {
		for i := 0; i < len(dbconfig.Options); i++ {
			if i > 0 {
				con += "&"
			}
			con += dbconfig.Options[i]
		}
	}
	return con
}

func buildMySQLConnectionString(dbconfig config.DBConfig) string {
	var con string
	// root:root@tcp(localhost:3306)/jwt_demo?parseTime=true
	con = dbconfig.Username + ":" + dbconfig.Password +
		"@" + dbconfig.Protocol + "(" + dbconfig.Server + ":" + strconv.Itoa(dbconfig.Port) + ")/" + dbconfig.DBName
	if len(dbconfig.Options) > 0 {
		for i := 0; i < len(dbconfig.Options); i++ {
			if i > 0 {
				con += "&"
			}
			con += dbconfig.Options[i]
		}
	}
	return con
}

func buildPostGresqlConnectionString(dbconfig config.DBConfig) string {
	var con string
	// host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai
	con = "host=" + dbconfig.Server + " user=" + dbconfig.Username + " password=" + dbconfig.Password +
		" port=" + strconv.Itoa(dbconfig.Port)
	if len(dbconfig.Options) > 0 {
		for i := 0; i < len(dbconfig.Options); i++ {
			con += " " + dbconfig.Options[i]
		}
	}
	return con
}

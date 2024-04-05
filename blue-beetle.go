package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"net/http"

	"blue-beetle/config"
	"blue-beetle/database"
)

func getConfig() (*config.SysConfig, error) {
	directory, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	config_path := filepath.Join(directory, "config", "config.yaml")
	_, err = os.Stat(config_path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("Could not find config file at: " + config_path)
		} else if errors.Is(err, os.ErrPermission) {
			return nil, errors.New("Do not have permissions to read config file at: " + config_path)
		} else {
			return nil, errors.New("Unknown error for config file at: " + config_path + " -> " + err.Error())
		}
	}
	return config.ProcessConfigYAMLFile(config_path)
}

func Migrate() {
	// database.Instance.AutoMigrate(&models.Page{})
	log.Println("Database Migration Completed!")
}

func main() {
	sconfig, err := getConfig()
	if err != nil {
		if strings.HasPrefix(err.Error(), "Could not find config file at") {
			sconfig = &config.SysConfig{Server: config.ServerConfig{Port: 8080}, Database: config.DBConfig{Type: "sqlite", FilePath: "wiki.db"}}
		} else {
			panic(err.Error())
		}
	}
	database.Connect(sconfig.Database)
	Migrate()

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(sconfig.Server.Port), nil))
}

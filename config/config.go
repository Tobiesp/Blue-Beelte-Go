package config

import (
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Type     string   `yaml:"type"`
	FilePath string   `yaml:"file-path"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	Server   string   `yaml:"server"`
	Protocol string   `yaml:"protocol"`
	Port     int      `yaml:"port"`
	DBName   string   `yaml:"db-name"`
	Options  []string `yaml:"options"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type SysConfig struct {
	Database DBConfig     `yaml:"database"`
	Server   ServerConfig `yaml:"server"`
}

func ProcessConfigYAMLFile(filePath string) (*SysConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("yamlFile.Get err " + err.Error())
	}
	yamlFile, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.New("yamlFile.Get err " + err.Error())
	}
	config, err := ProcessConfigYAML(string(yamlFile))
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return config, nil
}

func ProcessConfigYAML(yamlData string) (*SysConfig, error) {
	var config SysConfig
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		return nil, errors.New("Unmarshal: " + err.Error())
	}
	return &config, nil
}

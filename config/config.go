package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port         string `yaml:"port"`
		Env          int    `yaml:"env"`
		LogLevel     string `yaml:"log_level"`
		LogDir       string `yaml:"log_dir"`
		ShowSql      bool   `yaml:"show_sql"`
		NumOnceLimit int    `yaml:"num_once_limit"`
		ExecNoBackup bool   `yaml:"exec_no_backup"`
	} `yaml:"server"`

	DB struct {
		Address     string `yaml:"address"`
		Port        int    `yaml:"port"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		DBName      string `yaml:"db_name"`
		MaxIdleConn int    `yaml:"max_idle_conn"`
		MaxOpenConn int    `yaml:"max_open_conn"`
	} `yaml:"db"`
}

var Conf *Config

func InitConfig(path string) {
	if path == "" {
		path = getConfigPath()
	}

	conf, err := newConfig(path)
	if err != nil {
		log.Fatal("init config failed: ", err.Error())
	}
	Conf = conf
}

func newConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	err = yaml.NewDecoder(file).Decode(&config)

	return config, err
}

const (
	configPathEnv = "config_path"
	configPath    = "./config/config.yml"
)

func getConfigPath() string {
	path := os.Getenv(configPathEnv)
	if path != "" {
		return path
	}
	return configPath
}

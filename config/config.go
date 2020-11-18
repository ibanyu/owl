package config

import (
	"context"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"gitlab.pri.ibanyu.com/middleware/seaweed/xlog"
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
		AesKey       string `yaml:"aes_key"`
		AesIv        string `yaml:"aes_iv"`
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

	Login struct {
		LDAP struct {
			BaseDN  string `yaml:"base_dn"`
			Host    string `yaml:"host"`
			Port    int    `yaml:"port"`
			UseSSL  bool   `yaml:"use_ssl"`
			BindDN  string `yaml:"bind_dn"`
			BindPwd string `yaml:"bind_pwd"`
		} `yaml:"ldap"`

		LoginPath          string `yaml:"login_path"`
		TokenSecret        string `yaml:"token_secret"`
		TokenEffectiveHour int    `yaml:"token_effective_hour"`
	} `yaml:"login"`

	Role struct {
		From string `yaml:"from"`

		Conf struct {
			DBA              []string `yaml:"dba"`
			ReviewerRelation []struct {
				Reviewer string   `yaml:"reviewer"`
				Members  []string `yaml:"members"`
			} `yaml:"reviewer_relation"`
		} `yaml:"conf"`

		Net struct {
			ReviewerAPIAddress string `yaml:"reviewer_api_address"`
			ReviewerAPIToken   string `yaml:"reviewer_api_token"`
			DBAAPIAddress      string `yaml:"dba_api_address"`
			DBAAPIToken        string `yaml:"dba_api_token"`
			DBADepartmentID    int    `yaml:"dba_department_id"`
		} `yaml:"net"`
	} `yaml:"role"`
}

var Conf *Config

func InitConfig(path string) {
	if path == "" {
		path = getConfigPath()
	}

	conf, err := newConfig(path)
	if err != nil {
		xlog.Warnf(context.TODO(), "init config failed, err: %s, path: %s", err.Error(), path)
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

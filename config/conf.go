package config

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/vegacrypto/vegax_backend/tool"
	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Mysql struct {
		Port     string `yaml:"port"`
		Address  string `yaml:address`
		Username string `yaml:username`
		Password string `yaml:password`
		Db       string `yaml:db`
	}
	Log struct {
		Path  string `yaml:"path"`
		Level string `yaml:"level"`
		Node  string `yaml:"node"`
	}
}

var cfg *Config

func init() {
	env := "local"

	confFilePath := "./config/param-" + strings.ToLower(env) + ".yaml"

	if configFilePathFromEnv := os.Getenv("DALINK_GO_CONFIG_PATH"); configFilePathFromEnv != "" {
		confFilePath = configFilePathFromEnv
	}

	configFile, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		log.Fatal(err)
	}
	var data Config
	err2 := yaml.Unmarshal(configFile, &data)

	if err2 != nil {
		log.Fatal(err2)
	}
	cfg = &data

	writer2 := os.Stdout
	writer3, err := os.OpenFile(cfg.Log.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("create file log.txt failed: %v", err)
	}

	tool.Vlog.SetOutput(io.MultiWriter(writer2, writer3))
	if cfg.Log.Level == "debug" {
		tool.Vlog.SetLevel(logrus.DebugLevel)
	} else if cfg.Log.Level == "info" {
		tool.Vlog.SetLevel(logrus.InfoLevel)
	} else if cfg.Log.Level == "error" {
		tool.Vlog.SetLevel(logrus.ErrorLevel)
	} else if cfg.Log.Level == "warn" {
		tool.Vlog.SetLevel(logrus.WarnLevel)
	}
	tool.Vlog.Error("sys config done.")
}
func Get() *Config {
	return cfg
}

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

func NewConfig() *Server {
	return &Server{}
}

func LoadConfig() (*Server, error) {
	cfg := NewConfig()
	_ = godotenv.Load()
	appHome := os.Getenv("APP_HOME")
	if appHome == "" {
		return nil, errors.New("APP_HOME not set")
	}
	cfgPath := os.Getenv("CFG_PATH")
	if cfgPath == "" {
		return nil, errors.New("CFG_PATH not set")
	}
	cfgPath = filepath.Join(appHome, cfgPath)
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %s", err)
	}
	return cfg, nil
}

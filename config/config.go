package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
)

type Database struct {
	DBUsername string `json:"username"`
	DBPassword string `json:"password"`
	DBHost     string `json:"host"`
	DBPort     string `json:"port"`
	DBName     string `json:"database_name"`
}

type SecretManager struct {
	SecretString string
}

type PpnConfig struct {
	TimeStartPpn    int64 `json:"time_start_ppn"`
	TimeStartPpnNew int64 `json:"time_start_ppn_new"`
	TarifPpn        int64 `json:"tarif_ppn"`
	TarifPpnNew     int64 `json:"tarif_ppn_new"`
}

type Config struct {
	SourceDatabase  Database      `json:"source_database"`
	ServiceDatabase Database      `json:"service_database"`
	SecretManager   SecretManager `json:"secret_manager"`
	PpnConfig       PpnConfig     `json:"ppn_config"`
}

func LoadConfig(path string) (*Config, error) {
	bytes, err := os.ReadFile(filepath.Clean(path))
	log.Info("[config.LoadConfig]:: reading config file...")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return config, nil
}

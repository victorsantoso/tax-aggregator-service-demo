package config

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name           string
		args           args
		expectedConfig *Config
		expectedError  bool
	}{
		{
			name: "test failed when trying to read config file that doesn't exist",
			args: args{
				path: "../config/fail.json",
			},
			expectedConfig: nil,
			expectedError:  true,
		},
		{
			name: "test failed when format of json is not valid",
			args: args{
				path: "../config/fail.json",
			},
			expectedConfig: nil,
			expectedError:  true,
		},
		{
			name: "test passed when trying to open the file that exist",
			args: args{
				path: "../config/test.json",
			},
			expectedConfig: &Config{
				SourceDatabase: Database{
					DBUsername: "source-database-username",
					DBPassword: "source-database-password",
					DBHost:     "127.0.0.1",
					DBPort:     "3306",
					DBName:     "source-database-name",
				},
				ServiceDatabase: Database{
					DBUsername: "service-database-username",
					DBPassword: "service-database-password",
					DBHost:     "127.0.0.1",
					DBPort:     "5432",
					DBName:     "service-database-name",
				},
				SecretManager: SecretManager{
					SecretString: "",
				},
				PpnConfig: PpnConfig{
					TimeStartPpn:    1478624400,
					TarifPpn:        10,
					TimeStartPpnNew: 1648746000,
					TarifPpnNew:     11,
				},
			},
			expectedError: false,
		},
		{
			name: "test passed when trying to open config.json",
			args: args{
				path: "../config/config.json",
			},
			expectedConfig: &Config{
				SourceDatabase: Database{
					DBUsername: "root",
					DBPassword: "root",
					DBHost:     "127.0.0.1",
					DBPort:     "3306",
					DBName:     "source",
				},
				ServiceDatabase: Database{
					DBUsername: "postgres",
					DBPassword: "postgres",
					DBHost:     "127.0.0.1",
					DBPort:     "5432",
					DBName:     "tax-aggregator",
				},
				SecretManager: SecretManager{
					SecretString: "",
				},
				PpnConfig: PpnConfig{
					TimeStartPpn:    1478624400,
					TarifPpn:        10,
					TimeStartPpnNew: 1648746000,
					TarifPpnNew:     11,
				},
			},
			expectedError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := LoadConfig(tt.args.path)
			if (err != nil) != tt.expectedError {
				t.Errorf("LoadConfig() error = %v, expected error = %v", err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(result, tt.expectedConfig) {
				t.Errorf("LoadConfig() error = %v, expected config = %v", result, tt.expectedConfig)
			}
		})
	}
}

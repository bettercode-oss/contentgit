package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

var Config = struct {
	HttpPort   string
	DataSource struct {
		Host         string
		Port         string
		DatabaseName string
		UserName     string
		Password     string
	}
}{}

func InitConfig(path string) error {
	viper.AddConfigPath(path)

	if strings.EqualFold(os.Getenv("CONFIGOR_ENV"), "production") {
		viper.SetConfigName("production")
	} else {
		viper.SetConfigName("local")
	}
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	decodeHook := func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() == reflect.String {
			stringData := data.(string)
			if strings.HasPrefix(stringData, "${") && strings.HasSuffix(stringData, "}") {
				envVarValue := os.Getenv(strings.TrimPrefix(strings.TrimSuffix(stringData, "}"), "${"))
				if len(envVarValue) > 0 {
					return envVarValue, nil
				} else {
					return "", nil
				}
			}
		}
		return data, nil
	}
	if err := viper.Unmarshal(&Config, viper.DecodeHook(decodeHook)); err != nil {
		return err
	}

	return nil
}

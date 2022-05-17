package configuration

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var ErrInvalidConfig = errors.New("invalid config")
var ErrNameNotProvided = errors.New("config name not provided")
var ErrTypeNotProvided = errors.New("config type not provided")
var ErrPathNotProvided = errors.New("config path not provided")

const (
	Listen             = "Listen"
	JSONFile           = "JsonFile"
	CreateJSONEndpoint = "CreateJSONEndpoint"
	SQLFile            = "SQLFile"
)

var requiredConfig = []string{
	Listen,
	JSONFile,
	CreateJSONEndpoint,
	SQLFile,
}

type Configuration struct {
	Listen             string
	JSONFile           string
	CreateJSONEndpoint string
	SQLFile            string
}

// New is constructor for *Configuration object
func New(configName, configType, configPath string) (*Configuration, error) {
	if configName == "" {
		return nil, ErrNameNotProvided
	}

	if configType == "" {
		return nil, ErrTypeNotProvided
	}

	if configPath == "" {
		return nil, ErrPathNotProvided
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	homePath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	localConfigPath := fmt.Sprintf("%v/%v", homePath, "local_config")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath(localConfigPath)
	err = viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to read config file: %s.%s: %w", configName, configType, ErrInvalidConfig)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	config := &Configuration{
		Listen:             viper.GetString(Listen),
		JSONFile:           viper.GetString(JSONFile),
		CreateJSONEndpoint: viper.GetString(CreateJSONEndpoint),
		SQLFile:            viper.GetString(SQLFile),
	}

	if err := validateConfig(requiredConfig); err != nil {
		return nil, err
	}

	return config, nil
}

// validateConfig responsibility is to check if all required parameters are set within config file
func validateConfig(requiredParams []string) error {
	var missingParams []string

	for _, p := range requiredParams {
		if ok := viper.IsSet(p); !ok {
			missingParams = append(missingParams, p)
		}
	}

	if len(missingParams) > 0 {
		return fmt.Errorf("required parameters are missing: %q: %w", missingParams, ErrInvalidConfig)
	}

	return nil
}

package config

import "github.com/spf13/viper"

type Config struct {
	MongoConnectionString string `mapstructure:"mongo_connection_string"`
	MongoAuthMechanism    string `mapstructure:"mongo_auth_mechanism"`
	MongoAppName          string `mapstructure:"mongo_app_name"`
	MongoDatabaseName     string `mapstructure:"mongo_database_name"`
}

func Init(path string) (*Config, error) {
	vp := viper.New()
	vp.AddConfigPath(".")
	vp.SetConfigFile(path)
	vp.AutomaticEnv()
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = vp.UnmarshalExact(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

package config

import "github.com/spf13/viper"

type Config struct {
	MongoConnectionString  string `mapstructure:"mongo_connection_string"`
	MongoAuthMechanism     string `mapstructure:"mongo_auth_mechanism"`
	MongoAppName           string `mapstructure:"mongo_app_name"`
	MongoDatabaseName      string `mapstructure:"mongo_database_name"`
	ServerAddress          string `mapstructure:"server_address"`
	TerminationGracePeriod int    `mapstructure:"server_grace_period_in_seconds"`
	AutoAlertPeriod        int    `mapstructure:"auto_alert_period_in_minutes"`
	TwilioSID              string `mapstructure:"twilio_account_sid"`
	TwilioNumber           string `mapstructure:"twilio_number"`
	TwilioToken            string `mapstructure:"twilio_auth_token"`
	LogLevel               string `mapstructure:"log_level"`
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

package utils

import "github.com/spf13/viper"

// Config stores all configurations of the
// application, the values are read by viper
// from a config file or environment file
type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBUrl         string `mapstructure:"DB_URL"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// LoadConfig reads configuration values from a file path
// and returns the content as Config struct or an error
func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path)

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

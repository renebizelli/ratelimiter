package configs

import (
	pkg_utils "github.com/renebizelli/ratelimiter/pkg/utils"
	"github.com/spf13/viper"
)

type Config struct {
	WebServerPort                             string
	JWTSecret                                 string
	JWTExpires                                int
	RATELIMITER_IP_ON                         bool
	RATELIMITER_IP_MAX_REQUESTS               int
	RATELIMITER_IP_BLOCKED_SECONDS            int
	RATELIMITER_TOKEN_ON                      bool
	RATELIMITER_TOKEN_DEFAULT_MAX_REQUESTS    int
	RATELIMITER_TOKEN_DEFAULT_BLOCKED_SECONDS int
}

func LoadConfig(path string) *Config {

	var cfg *Config
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	pkg_utils.PanicIfError(err, "Load config file error")

	err = viper.Unmarshal(&cfg)
	pkg_utils.PanicIfError(err, "Unmarshal error")

	return cfg
}

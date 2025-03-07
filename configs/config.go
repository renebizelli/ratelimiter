package configs

import (
	"github.com/go-chi/jwtauth"
	pkg_utils "github.com/renebizelli/ratelimiter/pkg/utils"
	"github.com/spf13/viper"
)

type Config struct {
	WebServerPort string
	JWTToken      *jwtauth.JWTAuth
	JWTSecret     string
	JWTExpires    int64
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

	cfg.JWTToken = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	return cfg
}

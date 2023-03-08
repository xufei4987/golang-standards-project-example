package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang-standards-project-example/pkg/util/homedir"
	"log"
	"path/filepath"
	"strings"
)

type HttpServingInfo struct {
	Address string
}

const (
	// RecommendedHomeDir defines the default directory used to place all user service configurations.
	RecommendedHomeDir = ".user"

	// RecommendedEnvPrefix defines the ENV prefix used by all user service.
	RecommendedEnvPrefix = "EXAMPLE"
)

type Config struct {
	HttpServing *HttpServingInfo
	Mode        string
	Middlewares []string
	Healthz     bool
}

// NewConfig returns a Config struct with the default values.
func NewConfig() *Config {
	return &Config{
		Healthz:     true,
		Mode:        gin.DebugMode,
		Middlewares: []string{},
	}
}

// CompletedConfig is the completed configuration for GenericAPIServer.
type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// New returns a new instance of GenericAPIServer from the given config.
func (c CompletedConfig) New() (*GenericHttpServer, error) {
	// setMode before gin.New()
	gin.SetMode(c.Mode)

	s := &GenericHttpServer{
		HttpServingInfo: c.HttpServing,
		healthz:         c.Healthz,
		middlewares:     c.Middlewares,
		Engine:          gin.New(),
	}

	initGenericHttpServer(s)

	return s, nil
}

// LoadConfig reads in config file and ENV variables if set.
func LoadConfig(cfg string, defaultName string) {
	if cfg != "" {
		viper.SetConfigFile(cfg)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(filepath.Join(homedir.HomeDir(), RecommendedHomeDir))
		viper.SetConfigName(defaultName)
	}

	// Use config file from the flag.
	viper.SetConfigType("yaml")              // set the type of the configuration to yaml.
	viper.AutomaticEnv()                     // read in environment variables that match.
	viper.SetEnvPrefix(RecommendedEnvPrefix) // set ENVIRONMENT variables prefix.
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("WARNING: viper failed to discover and load the configuration file: %s\n", err.Error())
	}
}

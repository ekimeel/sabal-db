package env

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func GetConfigFile() string {
	file := os.Getenv(envConfigFile)
	if len(file) == 0 {
		log.Warnf(WarnDefault, envConfigFile, defaultConfigFile)
		return defaultConfigFile
	}

	return file
}

func GetLogLevel() string {
	level := os.Getenv(envLogLevel)
	if len(level) == 0 {
		log.Warnf(WarnDefault, envLogLevel, defaultLogLevel)
		return defaultLogLevel
	}
	return level
}

func SetConfig(config *Config) {
	GetEnv().Config = config
}

func GetEnv() *env {
	if impl == nil {
		impl = &env{}
	}
	return impl
}

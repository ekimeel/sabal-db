package env

const (
	/* defaults */
	defaultConfigFile = "./config.yml"
	defaultLogLevel   = "trace"

	/* os env variables */
	envConfigFile = "CONFIG_FILE"
	envLogLevel   = "LOG_LEVEL"

	/* logging */
	WarnDefault = "no [%s], using default [%s]"
)

var impl *env

type env struct {
	Config *Config
}

func (e *env) HasKafkaConsumer() bool {
	if e.Config == nil || e.Config.Kafka == nil || e.Config.Kafka.Consumer == nil || e.Config.Kafka.Consumer.Properties == nil {
		return false
	}

	return true
}

func (e *env) HasKafkaProducer() bool {
	if e.Config == nil || e.Config.Kafka == nil || e.Config.Kafka.Producer == nil {
		return false
	}

	return true
}

package env

const (
	defaultKafkaTimeout          = 5000
	defaultMaxFailedConnAttempts = 10
	defaultCacheSize             = 100000
	defaultCacheTTL              = 1800000
	defaultEager                 = true
	defaultPort                  = 6001
)

type Config struct {
	Name     string          `yaml:"name" json:"name"`
	Port     int             `yaml:"port,omitempty" json:"port,omitempty"`
	Node     *NodeConfig     `yaml:"node,omitempty" json:"node,omitempty"`
	Kafka    *KafkaConfig    `yaml:"kafka,omitempty" json:"kafka,omitempty"`
	Database *DatabaseConfig `yaml:"database,omitempty" json:"database,omitempty"`
}

type NodeConfig struct {
	Parallelism int          `yaml:"parallelism" json:"parallelism"`
	Cache       *CacheConfig `yaml:"cache,omitempty" json:"cache,omitempty"`
}

type CacheConfig struct {
	Size  uint64 `yaml:"size,omitempty" json:"size,omitempty"`
	TTL   uint64 `yaml:"ttl-ms,omitempty" json:"ttl-ms,omitempty"`
	Eager bool   `yaml:"eager,omitempty" json:"eager,omitempty"`
}

type KafkaConfig struct {
	Dialer   *KafkaDialerConfig   `yaml:"dialer,omitempty" json:"dialer,omitempty"`
	Consumer *KafkaConsumerConfig `yaml:"consumer,omitempty" json:"consumer,omitempty"`
	Producer *KafkaProducerConfig `yaml:"producer,omitempty" json:"producer,omitempty"`
}

type KafkaDialerConfig struct {
	Timeout               int `yaml:"timeout" json:"timeout"`
	MaxFailedConnAttempts int `yaml:"maxFailedConnAttempts,omitempty" json:"maxFailedConnAttempts,omitempty"`
}

type KafkaConsumerConfig struct {
	Topics     []string                `yaml:"topics,omitempty" json:"topics,omitempty"`
	Properties *map[string]interface{} `yaml:"properties,omitempty" json:"properties,omitempty"`
}

type KafkaProducerConfig struct {
	Topic      string                  `yaml:"topic,omitempty" json:"topic,omitempty"`
	Properties *map[string]interface{} `yaml:"properties,omitempty" json:"properties,omitempty"`
}

func (c *Config) SetDefaults() {

	if c.Port == 0 {
		c.Port = defaultPort
	}

	if c.Kafka != nil {
		if c.Kafka.Dialer == nil {
			c.Kafka.Dialer = &KafkaDialerConfig{
				Timeout:               defaultKafkaTimeout,
				MaxFailedConnAttempts: defaultMaxFailedConnAttempts,
			}
		}

		if c.Kafka.Dialer.MaxFailedConnAttempts == 0 {
			c.Kafka.Dialer.Timeout = defaultKafkaTimeout
		}

		if c.Kafka.Dialer.MaxFailedConnAttempts == 0 {
			c.Kafka.Dialer.MaxFailedConnAttempts = defaultMaxFailedConnAttempts
		}
	}

	if c.Node == nil {
		c.Node = &NodeConfig{
			Parallelism: 1,
			Cache:       &CacheConfig{Size: defaultCacheSize, TTL: defaultCacheTTL, Eager: defaultEager},
		}
	}

	if c.Node.Cache == nil {
		c.Node.Cache = &CacheConfig{Size: defaultCacheSize, TTL: defaultCacheTTL, Eager: defaultEager}
	}

	if c.Node.Parallelism == 0 {
		c.Node.Parallelism = 1
	}

}

func GetKafkaProducerConfig() *KafkaProducerConfig {
	return GetEnv().Config.Kafka.Producer
}

type DatabaseConfig struct {
	Path string `yaml:"path" json:"path"`
}

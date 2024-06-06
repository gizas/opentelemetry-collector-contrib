package elasticprocessor

type Config struct {
	AddSystemMetrics     bool `mapstructure:"add_system_metrics"`
	AddKubernetesMetrics bool `mapstructure:"add_k8s_metrics"`
}

func (c *Config) Validate() error {
	return nil
}

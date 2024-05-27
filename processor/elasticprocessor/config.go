package elasticprocessor

type Config struct {
	AddK8sMetrics bool `mapstructure:"add_k8s_metrics"`
}

func (c *Config) Validate() error {
	return nil
}

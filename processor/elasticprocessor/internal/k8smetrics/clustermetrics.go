package hostmetrics

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func addContainerMetrics(metrics pmetric.MetricSlice, resource pcommon.Resource, dataset string) error {
	var timestamp pcommon.Timestamp
	var cpu_limit_utilization, memory_usage_limit_pct float64

	// iterate all metrics in the current scope and generate the additional Elastic kubernetes integration metrics
	for i := 0; i < metrics.Len(); i++ {
		metric := metrics.At(i)
		if metric.Name() == "k8s.container.cpu.limit.utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			cpu_limit_utilization = dp.DoubleValue()
		} else if metric.Name() == "k8s.container.memory.limit.utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			memory_usage_limit_pct = dp.DoubleValue()
		}
	}

	addMetrics(metrics, resource, dataset,
		metric{
			dataType:    Gauge,
			name:        "kubernetes.container.cpu.usage.limit.pct",
			timestamp:   timestamp,
			doubleValue: &cpu_limit_utilization,
		},
		metric{
			dataType:    Gauge,
			name:        "kubernetes.container.memory.usage.limit.pct",
			timestamp:   timestamp,
			doubleValue: &memory_usage_limit_pct,
		},
	)

	return nil
}

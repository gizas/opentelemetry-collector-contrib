package hostmetrics

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func addClusterMetrics(metrics pmetric.MetricSlice, group string) error {
	var timestamp pcommon.Timestamp
	var node_allocatable_memory, node_allocatable_cpu int64

	// iterate all metrics in the current scope and generate the additional Elastic kubernetes integration metrics
	for i := 0; i < metrics.Len(); i++ {
		metric := metrics.At(i)
		if metric.Name() == "k8s.node.allocatable_cpu" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			node_allocatable_cpu = dp.IntValue()
		} else if metric.Name() == "k8s.node.allocatable_memory" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			node_allocatable_memory = dp.IntValue()
		}
	}

	addMetrics(metrics, group,
		metric{
			dataType:  Gauge,
			name:      "kubernetes.node.cpu.allocatable.cores",
			timestamp: timestamp,
			intValue:  &node_allocatable_cpu,
		},
		metric{
			dataType:  Gauge,
			name:      "kubernetes.node.memory.allocatable.bytes",
			timestamp: timestamp,
			intValue:  &node_allocatable_memory,
		},
	)

	return nil
}

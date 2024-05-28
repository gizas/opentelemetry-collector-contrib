package hostmetrics

import (
	"math"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func addNodeMetrics(metrics pmetric.MetricSlice, resource pcommon.Resource, dataset string) error {
	var timestamp pcommon.Timestamp
	var memory_usage, filesystem_capacity, filesystem_usage int64
	var cpu_usage float64

	for i := 0; i < metrics.Len(); i++ {
		metric := metrics.At(i)
		if metric.Name() == "k8s.node.cpu.usage" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			cpu_usage = dp.DoubleValue() * math.Pow10(9)
		} else if metric.Name() == "k8s.node.memory.usage" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			memory_usage = dp.IntValue()
		} else if metric.Name() == "k8s.node.filesystem.capacity" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			filesystem_capacity = dp.IntValue()
		} else if metric.Name() == "k8s.node.filesystem.usage" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			filesystem_usage = dp.IntValue()
		}

	}

	addMetrics(metrics, resource, dataset,
		metric{
			dataType:    Gauge,
			name:        "kubernetes.node.cpu.usage.nanocores",
			timestamp:   timestamp,
			doubleValue: &cpu_usage,
		},
		metric{
			dataType:  Gauge,
			name:      "kubernetes.node.memory.usage.bytes",
			timestamp: timestamp,
			intValue:  &memory_usage,
		},
		metric{
			dataType:  Gauge,
			name:      "kubernetes.node.fs.capacity.bytes",
			timestamp: timestamp,
			intValue:  &filesystem_capacity,
		},
		metric{
			dataType:  Gauge,
			name:      "kubernetes.node.fs.used.bytes",
			timestamp: timestamp,
			intValue:  &filesystem_usage,
		},
	)

	return nil
}

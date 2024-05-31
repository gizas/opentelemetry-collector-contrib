package hostmetrics

import (
	"math"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func addkubeletMetrics(metrics pmetric.MetricSlice, group string) error {
	var timestamp pcommon.Timestamp
	var total_transmited, total_received, memory_usage, filesystem_capacity, filesystem_usage int64
	var cpu_limit_utilization, container_cpu_limit_utilization, memory_usage_limit_pct, memory_limit_utilization, cpu_usage float64

	// iterate all metrics in the current scope and generate the additional Elastic kubernetes integration metrics

	//pod
	for i := 0; i < metrics.Len(); i++ {
		metric := metrics.At(i)
		// kubernetes.pod.cpu.usage.node.pct still needs to be implemented
		if metric.Name() == "k8s.pod.cpu_limit_utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			cpu_limit_utilization = dp.DoubleValue() * 100
		} else if metric.Name() == "k8s.pod.memory_limit_utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			memory_limit_utilization = dp.DoubleValue() * 100
		} else if metric.Name() == "k8s.pod.network.io" {
			dataPoints := metric.Sum().DataPoints()
			for j := 0; j < dataPoints.Len(); j++ {
				dp := dataPoints.At(j)
				if timestamp == 0 {
					timestamp = dp.Timestamp()
				}

				value := dp.IntValue()
				if direction, ok := dp.Attributes().Get("direction"); ok {
					switch direction.Str() {
					case "receive":
						total_received += value
					case "transmit":
						total_transmited += value
					}
				}
			}
			//node
		} else if metric.Name() == "k8s.node.cpu.usage" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp() * 100
			}
			cpu_usage = dp.DoubleValue() * math.Pow10(9)
		} else if metric.Name() == "k8s.node.memory.usage" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp() * 100
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
				timestamp = dp.Timestamp() * 100
			}
			filesystem_usage = dp.IntValue()
			// container
		} else if metric.Name() == "k8s.container.cpu_limit_utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp() * 100
			}
			container_cpu_limit_utilization = dp.DoubleValue()
		} else if metric.Name() == "k8s.container.memory_limit_utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			memory_usage_limit_pct = dp.DoubleValue() * 100
		}

	}

	addMetrics(metrics, group,
		metric{
			dataType:    Gauge,
			name:        "kubernetes.pod.cpu.usage.limit.pct",
			timestamp:   timestamp,
			doubleValue: &cpu_limit_utilization,
		},
		metric{
			dataType:    Gauge,
			name:        "kubernetes.pod.memory.usage.limit.pct",
			timestamp:   timestamp,
			doubleValue: &memory_limit_utilization,
		},
		metric{
			dataType:  Sum,
			name:      "kubernetes.pod.network.tx.bytes",
			timestamp: timestamp,
			intValue:  &total_transmited,
		},
		metric{
			dataType:  Sum,
			name:      "kubernetes.pod.network.rx.bytes",
			timestamp: timestamp,
			intValue:  &total_received,
		}, metric{
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
		metric{
			dataType:    Gauge,
			name:        "kubernetes.container.cpu.usage.limit.pct",
			timestamp:   timestamp,
			doubleValue: &container_cpu_limit_utilization,
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

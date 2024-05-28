package hostmetrics

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func addPodMetrics(metrics pmetric.MetricSlice, resource pcommon.Resource, dataset string) error {
	var timestamp pcommon.Timestamp
	var total_transmited, total_received int64
	var cpu_limit_utilization, memory_limit_utilization float64

	// iterate all metrics in the current scope and generate the additional Elastic kubernetes integration metrics
	for i := 0; i < metrics.Len(); i++ {
		metric := metrics.At(i)
		if metric.Name() == "k8s.pod.cpu_limit_utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			cpu_limit_utilization = dp.DoubleValue()
		} else if metric.Name() == "k8s.pod.memory_limit_utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			memory_limit_utilization = dp.DoubleValue()
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
		}

	}

	addMetrics(metrics, resource, dataset,
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
		},
	)

	return nil
}

package hostmetrics

import (
	"fmt"
	"math"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

func addkubeletMetrics(metrics pmetric.MetricSlice, group string) error {
	var timestamp pcommon.Timestamp
	var total_transmited, total_received, node_memory_usage, filesystem_capacity, filesystem_usage, node_cpu_utilization, pod_memory_available, pod_memory_usage int64
	var cpu_limit_utilization, container_cpu_limit_utilization, memory_usage_limit_pct, memory_limit_utilization, node_cpu_usage, node_cpu_available, pod_cpu_usage, pod_cpu_usage_node, pod_memory_usage_node float64

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
			cpu_limit_utilization = dp.DoubleValue()
		} else if metric.Name() == "k8s.pod.cpu.usage" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			pod_cpu_usage = dp.DoubleValue()
		} else if metric.Name() == "k8s.pod.memory_limit_utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			memory_limit_utilization = dp.DoubleValue()
		} else if metric.Name() == "k8s.pod.memory.available" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			pod_memory_available = dp.IntValue()
		} else if metric.Name() == "k8s.pod.memory.usage" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			pod_memory_usage = dp.IntValue()

			if (pod_memory_available) > 0 {
				pod_memory_usage_node = float64(pod_memory_usage) / (float64(pod_memory_available) + float64(pod_memory_usage))
				name, _ := dp.Attributes().Get("k8s.pod.name")
				return fmt.Errorf("pas %v , %v, s%v , %v", name, pod_memory_usage_node, pod_memory_usage, pod_memory_available)
			}
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
				timestamp = dp.Timestamp()
			}
			node_cpu_usage = dp.DoubleValue() * math.Pow10(9)
		} else if metric.Name() == "k8s.node.cpu.utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			node_cpu_utilization = dp.IntValue()
		} else if metric.Name() == "k8s.node.memory.usage" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			node_memory_usage = dp.IntValue()
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
				timestamp = dp.Timestamp()
			}
			container_cpu_limit_utilization = dp.DoubleValue()
		} else if metric.Name() == "k8s.container.memory_limit_utilization" {
			dp := metric.Gauge().DataPoints().At(0)
			if timestamp == 0 {
				timestamp = dp.Timestamp()
			}
			memory_usage_limit_pct = dp.DoubleValue()
		}

	}

	if (node_cpu_utilization) > 0 {
		node_cpu_available = (node_cpu_usage * float64(1-node_cpu_utilization)) / float64(node_cpu_utilization)
	}
	if (node_cpu_available) > 0 {
		pod_cpu_usage_node = pod_cpu_usage / node_cpu_available
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
			name:        "kubernetes.pod.cpu.usage.node.pct",
			timestamp:   timestamp,
			doubleValue: &pod_cpu_usage_node,
		},
		metric{
			dataType:    Gauge,
			name:        "kubernetes.pod.memory.usage.node.pct",
			timestamp:   timestamp,
			doubleValue: &pod_memory_usage_node,
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
			doubleValue: &node_cpu_usage,
		},
		metric{
			dataType:  Gauge,
			name:      "kubernetes.node.memory.usage.bytes",
			timestamp: timestamp,
			intValue:  &node_memory_usage,
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

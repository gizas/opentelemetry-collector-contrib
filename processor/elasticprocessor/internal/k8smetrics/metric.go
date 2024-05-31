package hostmetrics

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type dataType int

const (
	Gauge dataType = iota
	Sum
)

type metric struct {
	dataType       dataType
	name           string
	timestamp      pcommon.Timestamp
	startTimestamp pcommon.Timestamp
	intValue       *int64
	doubleValue    *float64
	attributes     *pcommon.Map
}

func addMetrics(ms pmetric.MetricSlice, group string, metrics ...metric) {
	ms.EnsureCapacity(ms.Len() + len(metrics))

	for _, metric := range metrics {
		m := ms.AppendEmpty()
		m.SetName(metric.name)

		var dp pmetric.NumberDataPoint
		switch metric.dataType {
		case Gauge:
			dp = m.SetEmptyGauge().DataPoints().AppendEmpty()
		case Sum:
			dp = m.SetEmptySum().DataPoints().AppendEmpty()
		}

		if metric.intValue != nil {
			dp.SetIntValue(*metric.intValue)
		} else if metric.doubleValue != nil {
			dp.SetDoubleValue(*metric.doubleValue)
		}

		dp.SetTimestamp(metric.timestamp)
		if metric.startTimestamp != 0 {
			dp.SetStartTimestamp(metric.startTimestamp)
		}

		// Calculate datastream attribute as an attribute to each datapoint
		dataset := addDatastream(metric.name)

		if metric.attributes != nil {
			ip, ipok := metric.attributes.Get("k8s.pod.ip")
			if dataset == "kubernetes.node" {
				slice, ok := metric.attributes.Get("k8s.node.name")
				if ok {
					dp.Attributes().PutStr("k8s.node.name", slice.AsString())
				}

			} else {
				metric.attributes.CopyTo(dp.Attributes())
				if ipok {
					dp.Attributes().PutStr("k8s.pod.ip", ip.AsString())
				}
			}
		}
		dp.Attributes().PutStr("event.module", "kubernetes")
		dp.Attributes().PutStr("data_stream.dataset", dataset)
	}
}

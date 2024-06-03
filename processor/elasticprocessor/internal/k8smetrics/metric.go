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
			if dataset == "kubernetes.node" {
				slice, ok := metric.attributes.Get("k8s.node.name")
				slice_uid, _ := metric.attributes.Get("k8s.node.uid")
				// slice_pod, _ := metric.attributes.Get("k8s.pod.name")
				// slice_pod_uid, _ := metric.attributes.Get("k8s.pod.uid")
				if ok {
					dp.Attributes().PutStr("k8s.node.name", slice.AsString())
					dp.Attributes().PutStr("k8s.node.uid", slice_uid.AsString())
					// dp.Attributes().PutStr("k8s.pod.name", slice_pod.AsString())
					// dp.Attributes().PutStr("k8s.pod.uid", slice_pod_uid.AsString())
				}

			} else {
				metric.attributes.CopyTo(dp.Attributes())
			}
		}
		dp.Attributes().PutStr("event.module", "kubernetes")
		dp.Attributes().PutStr("data_stream.dataset", dataset)
	}
}

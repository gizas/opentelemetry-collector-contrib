package metadata

import (
	"go.opentelemetry.io/collector/component"
)

var (
	Type = component.MustNewType("elastic")
)

const (
	TracesStability  = component.StabilityLevelBeta
	MetricsStability = component.StabilityLevelBeta
	LogsStability    = component.StabilityLevelBeta
)

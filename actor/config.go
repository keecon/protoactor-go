package actor

import (
	"time"

	"go.opentelemetry.io/otel/metric"
)

type Config struct {
	DeadLetterThrottleInterval  time.Duration      // throttle deadletter logging after this interval
	DeadLetterThrottleCount     int32              // throttle deadletter logging after this count
	DeadLetterRequestLogging    bool               // do not log dead-letters with sender
	DeveloperSupervisionLogging bool               // console log and promote supervision logs to Warning level
	DiagnosticsSerializer       func(Actor) string // extract diagnostics from actor and return as string
	MetricsProvider             metric.MeterProvider
}

func defaultConfig() *Config {
	return &Config{
		MetricsProvider:             nil,
		DeadLetterThrottleInterval:  1 * time.Second,
		DeadLetterThrottleCount:     3,
		DeadLetterRequestLogging:    true,
		DeveloperSupervisionLogging: false,
		DiagnosticsSerializer: func(actor Actor) string {
			return ""
		},
	}
}

func NewConfig() *Config {
	return defaultConfig()
}

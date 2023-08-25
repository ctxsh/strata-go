package strata

import "time"

const (
	DefaultTimeout                  = 5 * time.Second
	DefaultMaxAge     time.Duration = 10 * time.Minute
	DefaultAgeBuckets uint32        = 5
)

package config

import "time"

type Config struct {
	RateCheckerInterval time.Duration
	NumberOfChecks      uint
}

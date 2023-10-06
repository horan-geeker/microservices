package jwt

import (
	"time"
)

// Options contains configuration items related to API server features.
type Options struct {
	Realm      string        `json:"realm"       mapstructure:"realm"`
	Key        string        `json:"key"         mapstructure:"key"`
	Timeout    time.Duration `json:"timeout"     mapstructure:"timeout"`
	MaxRefresh time.Duration `json:"maxRefresh" mapstructure:"max-refresh"`
}

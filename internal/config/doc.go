// Package config defines the Config struct used throughout logslice to carry
// runtime settings parsed from CLI flags. It also provides a Validate method
// that enforces business rules (e.g. valid URL scheme, coherent time range)
// before any network or I/O operations begin.
package config

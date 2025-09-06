package config

import (
	"errors"
)

var (
	errInvalidConfigInstance = errors.New("invalid config manager instance")
)

const (
	InstanceViper int = iota
)

// NewConfigManagerFactory creates a new config manager instance based on the specified type
func NewConfigManagerFactory(instance int) (ConfigManager, error) {
	switch instance {
	case InstanceViper:
		return NewViperConfigManager()
	default:
		return nil, errInvalidConfigInstance
	}
}

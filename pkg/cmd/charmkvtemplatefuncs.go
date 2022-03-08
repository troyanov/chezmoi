package cmd

import (
	"time"

	"github.com/charmbracelet/charm/kv"
)

type charmKVConfig struct {
	Name          string        `mapstructure:"name"`
	RefreshPeriod time.Duration `mapstructure:"refreshPeriod"`
	kv            *kv.KV
}

func (c *Config) charmKVTemplateFunc(key string) string {
	// FIXME add name as second argument
	// FIXME support encryption
	// FIXME check other BadgerDB options

	if c.CharmKV.kv == nil {
		var err error
		c.CharmKV.kv, err = kv.OpenWithDefaults(c.CharmKV.Name)
		if err != nil {
			panic(err)
		}

		// FIXME call c.CharmKV.Sync() if last refresh was more than
		// refreshPeriod ago
	}

	value, err := c.CharmKV.kv.Get([]byte(key))
	if err != nil {
		panic(err)
	}

	return string(value)
}

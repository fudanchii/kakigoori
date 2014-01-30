package event

import (
	"log"
	"github.com/mitchellh/mapstructure"
)

type Config map[string]interface{}

func (cfg *Config) get(key string) Config {
	if value := (*cfg)[key]; value != nil {
		return value.(map[string]interface{})
	}
	log.Printf("ERR> %s is nil", key)
	return nil
}

func (cfg *Config) Decode(key string, val interface{}) error {
	if cfgstruct := cfg.get(key); cfgstruct != nil {
		return mapstructure.Decode(cfgstruct, val)
	}
	return cfg
}

func (cfg *Config) Error() string {
	return "Config error."
}

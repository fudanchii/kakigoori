package event

type Config map[string]interface{}

func (cfg *Config) Get(key string) Config {
	return (*cfg)[key].(map[string]interface{})
}

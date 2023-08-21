//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DNS: "",
	},
	Redis: RedisConfig{
		Addr: "",
	},
}

//go:build k8s

package config

var Config = config{
	DB: DBConfig{
		DNS: "root:root@tcp(localhost:13316)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:6380",
	},
}

//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DNS: "root:Cloudwalk#galaxy@tcp(10.178.16.5:33091)/webook",
	},
	Redis: RedisConfig{
		Addr: "",
	},
}

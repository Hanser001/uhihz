package main

import "zhihu/boot"

func main() {
	boot.ViperSet("./config/config.yaml")
	boot.LoggerSet("./docs/all.log")
	boot.MysqlSet()
	boot.RedisSet()
	boot.ServerSet()
}

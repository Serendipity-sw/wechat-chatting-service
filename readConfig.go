package main

import "github.com/guotie/config"

func readConfig() {
	logsDir = config.GetStringDefault("logsDir", "./logs")
}

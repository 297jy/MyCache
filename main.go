package main

import (
	"fmt"
	"gomemory/config"
	"gomemory/lib/logger"
	"gomemory/lib/utils"
	"gomemory/server"
	"gomemory/tcp"
	"os"
)

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/
`

var defaultProperties = &config.ServerProperties{
	Bind:  "0.0.0.0",
	Port:  6399,
	RunID: utils.RandString(40),
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	print(banner)
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "GoMemory",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	configFilename := os.Getenv("CONFIG")

	if configFilename == "" {
		if fileExists("mycache.conf") {
			config.SetupCacheConfig("mycache.conf")
		} else {
			config.Properties = defaultProperties
		}
	} else {
		config.SetupCacheConfig(configFilename)
	}
	err := tcp.ListenAndServerWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, server.MakeHandler())

	if err != nil {
		logger.Error(err)
	}
}

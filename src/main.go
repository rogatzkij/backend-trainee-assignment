package main

import (
	"github.com/rs/zerolog/log"

	"os"
	"os/signal"
)

func main() {
	config, err := InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("не удалось прочитать настройки")
	}

	controller, err := NewConnector(config.ConnectorType)
	if err != nil {
		log.Fatal().Err(err).Msg("не создать коннектор")
	}

	service := NewService(config, controller)

	service.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	service.Stop()
}

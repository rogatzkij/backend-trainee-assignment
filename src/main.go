package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"

	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	config, err := initConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Не удалось прочитать настройки")
	}

	router := initRouter()
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: router,
	}

	go func() {
		log.Info().Str("Host", config.Host).Int("Port", config.Port).Msg("Сервис запускается")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Не удалось запустить сервис")
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Info().Msg("Сервис закрывается...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Warn().Err(err).Msg("Ошибка закрытия сервиса")
	}
	log.Info().Msg("Сервис закрыт")
}

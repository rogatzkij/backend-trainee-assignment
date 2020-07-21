package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Service struct {
	config    *Config
	connector Connector
	server    http.Server
}

func (s *Service) Start() {
	go func() {
		log.Info().Str("Host", s.config.Host).Int("Port", s.config.Port).Msg("Сервис запущен")
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("Не удалось запустить сервис")
		}
	}()
}

func (s *Service) Stop() {
	log.Info().Msg("Сервис закрывается...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Warn().Err(err).Msg("Ошибка закрытия сервиса")
	}
	log.Info().Msg("Сервис закрыт")
}

func NewService(config *Config, controller Connector) *Service {
	service := &Service{
		config:    config,
		connector: controller,
	}

	service.server = http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: service.initRouter(),
	}

	return service
}

type Config struct {
	Port          int    `default:"9000"`
	Host          string `default:"localhost"`
	ConnectorType string `default:"postgres"`
}

func InitConfig() (*Config, error) {
	config := &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *Service) initRouter() *mux.Router {
	router := mux.NewRouter()

	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/add", s.createUser).Methods(http.MethodPost)

	chatRouter := router.PathPrefix("/chats").Subrouter()
	chatRouter.HandleFunc("/add", s.createChat).Methods(http.MethodPost)
	chatRouter.HandleFunc("/get", s.getChats).Methods(http.MethodPost)

	messagesRouter := router.PathPrefix("/messages").Subrouter()
	messagesRouter.HandleFunc("/add", s.sendMessage).Methods(http.MethodPost)
	messagesRouter.HandleFunc("/get", s.getMessages).Methods(http.MethodPost)

	return router
}

// Добавить нового пользователя
func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Создать новый чат между пользователями
func (s *Service) createChat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Отправить сообщение в чат от лица пользователя
func (s *Service) sendMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Получить список чатов конкретного пользователя
func (s *Service) getChats(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Получить список сообщений в конкретном чате
func (s *Service) getMessages(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

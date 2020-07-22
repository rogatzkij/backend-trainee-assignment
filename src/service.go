package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"io/ioutil"
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
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: service.initRouter(),
	}

	return service
}

type Config struct {
	Port          int    `default:"9000"`
	Host          string `default:""`
	ConnectorType string `split_words:"true" default:"postgres"`
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

	router.HandleFunc("/", s.createUser).Methods(http.MethodPost)

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

type ErrorCodeType int

const (
	AlreadyExist ErrorCodeType = iota
	NotExist
)

type ErrorResponse struct {
	ErrorCode   ErrorCodeType `json:"code"`
	Description string        `json:"description"`
}

// Добавить нового пользователя
func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось прочитать тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	requestBody := struct {
		Username string `json:"username"`
	}{}

	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Warn().Err(err).Msg("Не удалось анмаршалить тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isExist, err := s.connector.checkUsername(requestBody.Username)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось проверить пользователя")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isExist {
		responseBody := ErrorResponse{
			ErrorCode:   AlreadyExist,
			Description: fmt.Sprintf("Пользователь уже %s существует", requestBody.Username),
		}

		body, err = json.Marshal(responseBody)
		if err != nil {
			log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	user, err := s.connector.createUser(requestBody.Username)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось создать пользователя")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody := struct {
		ID uint64 `json:"id"`
	}{
		ID: user.ID,
	}

	body, err = json.Marshal(responseBody)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

// Создать новый чат между пользователями
func (s *Service) createChat(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		Name  string   `json:"name"`
		Users []uint64 `json:"users"`
	}{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось прочитать тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Warn().Err(err).Msg("Не удалось анмаршалить тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Проверяем существование чата
	isExist, err := s.connector.checkChartName(requestBody.Name)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось проверить пользователя")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isExist {
		responseBody := ErrorResponse{
			ErrorCode:   AlreadyExist,
			Description: fmt.Sprintf("Чат уже %s существует", requestBody.Name),
		}

		body, err = json.Marshal(responseBody)
		if err != nil {
			log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	// Проверяем существование пользователей
	for _, userID := range requestBody.Users {
		isExist, err := s.connector.checkUserID(userID)
		if err != nil {
			log.Warn().Err(err).Msg("Не удалось проверить пользователя")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !isExist {
			responseBody := ErrorResponse{
				ErrorCode:   NotExist,
				Description: fmt.Sprintf("Пользователь c id %d не существует", userID),
			}

			body, err = json.Marshal(responseBody)
			if err != nil {
				log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(body)
			return
		}
	}

	// Создаем чат
	_, err = s.connector.createChart(requestBody.Name, requestBody.Users)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось создать чат")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Отправить сообщение в чат от лица пользователя
func (s *Service) sendMessage(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		ChatID uint64 `json:"chat"`
		UserID uint64 `json:"author"`
		Text   string `json:"text"`
	}{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось прочитать тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Warn().Err(err).Msg("Не удалось анмаршалить тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Проверяем существование чата
	isExist, err := s.connector.checkChartID(requestBody.ChatID)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось проверить пользователя")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isExist {
		responseBody := ErrorResponse{
			ErrorCode:   NotExist,
			Description: fmt.Sprintf("Чат c id %d не существует", requestBody.ChatID),
		}

		body, err = json.Marshal(responseBody)
		if err != nil {
			log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	// Проверяем существование пользователей
	isExist, err = s.connector.checkUserID(requestBody.UserID)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось проверить пользователя")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isExist {
		responseBody := ErrorResponse{
			ErrorCode:   NotExist,
			Description: fmt.Sprintf("Пользователь c id %d не существует", requestBody.UserID),
		}

		body, err = json.Marshal(responseBody)
		if err != nil {
			log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	// Отправляем сообщение
	_, err = s.connector.sendMessage(requestBody.ChatID, requestBody.UserID, requestBody.Text)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось создать чат")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Получить список чатов конкретного пользователя
func (s *Service) getChats(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		UserID uint64 `json:"user"`
	}{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось прочитать тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Warn().Err(err).Msg("Не удалось анмаршалить тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Проверяем существование пользователей
	isExist, err := s.connector.checkUserID(requestBody.UserID)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось проверить пользователя")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !isExist {
		responseBody := ErrorResponse{
			ErrorCode:   NotExist,
			Description: fmt.Sprintf("Пользователь c id %d не существует", requestBody.UserID),
		}

		body, err = json.Marshal(responseBody)
		if err != nil {
			log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	// Отправляем сообщение
	chats, err := s.connector.getCharts(requestBody.UserID)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось создать чат")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody := struct {
		Chats []Chat `json:"chats"`
	}{
		Chats: chats,
	}

	body, err = json.Marshal(responseBody)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// Получить список сообщений в конкретном чате
func (s *Service) getMessages(w http.ResponseWriter, r *http.Request) {
	requestBody := struct {
		ChatID uint64 `json:"chat"`
	}{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось прочитать тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Warn().Err(err).Msg("Не удалось анмаршалить тело")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Проверяем существование чата
	isExist, err := s.connector.checkChartID(requestBody.ChatID)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось проверить пользователя")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isExist {
		responseBody := ErrorResponse{
			ErrorCode:   NotExist,
			Description: fmt.Sprintf("Чат c id %d не существует", requestBody.ChatID),
		}

		body, err = json.Marshal(responseBody)
		if err != nil {
			log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
		return
	}

	// Получаем сообщения
	messages, err := s.connector.getMessages(requestBody.ChatID)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось создать чат")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody := struct {
		Messages []Message `json:"messages"`
	}{
		Messages: messages,
	}

	body, err = json.Marshal(responseBody)
	if err != nil {
		log.Warn().Err(err).Msg("Не удалось замаршалить ответ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

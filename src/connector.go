package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"strings"
)

type Connector interface {
	createUser(username string) (User, error)
	createChart(name string, users []string) (Chat, error)
	getCharts(username string) ([]Chat, error)
	sendMessage(chatID uint64, authorID uint64, text string) (Message, error)
}

func NewConnector(controllerType string) (Connector, error) {
	switch strings.ToLower(controllerType) {
	case "postgressql", "postgres_sql", "postgres":
		config, err := initConfigPostgresSQL()
		if err != nil {
			return nil, err
		}

		return &ConnectorPostgresSQL{config: config}, nil
	default:
		return nil, fmt.Errorf("неизвестный коннектор %s", controllerType)
	}
}

type ConfigPostgresSQL struct {
	Login    string `default:"admin"`
	Password string `default:"admin"`
	Host     string `default:"db"`
	Port     string `default:"8080"`
}

func initConfigPostgresSQL() (*ConfigPostgresSQL, error) {
	config := &ConfigPostgresSQL{}
	err := envconfig.Process("PostgresSQL", config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

type ConnectorPostgresSQL struct {
	config *ConfigPostgresSQL
}

func (cp *ConnectorPostgresSQL) createUser(username string) (User, error) {
	panic("implement me")
}

func (cp *ConnectorPostgresSQL) createChart(name string, users []string) (Chat, error) {
	panic("implement me")
}

func (cp *ConnectorPostgresSQL) getCharts(username string) ([]Chat, error) {
	panic("implement me")
}

func (cp *ConnectorPostgresSQL) sendMessage(chatID uint64, authorID uint64, text string) (Message, error) {
	panic("implement me")
}

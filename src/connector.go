package main

import (
	"fmt"
	"strings"
)

// Интерфейс описывает работу с хранилищем данных
type Connector interface {
	createUser(username string) (User, error)
	checkUsername(username string) (bool, error)
	checkUserID(user uint64) (bool, error)
	createChart(name string, users []uint64) (Chat, error)
	checkChartName(name string) (bool, error)
	checkChartID(chat uint64) (bool, error)
	getCharts(user uint64) ([]Chat, error)
	sendMessage(chatID uint64, authorID uint64, text string) (Message, error)
	getMessages(chatID uint64) ([]Message, error)
}

func NewConnector(controllerType string) (Connector, error) {
	switch strings.ToLower(controllerType) {
	case "mysql":
		config, err := initConfigMySQL()
		if err != nil {
			return nil, err
		}

		return &ConnectorMySQL{config: config}, nil
	default:
		return nil, fmt.Errorf("неизвестный коннектор %s", controllerType)
	}
}

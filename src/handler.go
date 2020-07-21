package main

import "net/http"

// Добавить нового пользователя
func createUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Создать новый чат между пользователями
func createChat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Отправить сообщение в чат от лица пользователя
func sendMessage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Получить список чатов конкретного пользователя
func getChats(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Получить список сообщений в конкретном чате
func getMessages(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

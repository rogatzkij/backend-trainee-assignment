package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func initRouter() *mux.Router {
	router := mux.NewRouter()

	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/add", createUser).Methods(http.MethodPost)

	chatRouter := router.PathPrefix("/chats").Subrouter()
	chatRouter.HandleFunc("/add", createChat).Methods(http.MethodPost)
	chatRouter.HandleFunc("/get", getChats).Methods(http.MethodPost)

	messagesRouter := router.PathPrefix("/messages").Subrouter()
	messagesRouter.HandleFunc("/add", sendMessage).Methods(http.MethodPost)
	messagesRouter.HandleFunc("/get", getMessages).Methods(http.MethodPost)

	return router
}
